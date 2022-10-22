package main

import (
	"container/heap"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/crawler"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version"
)

// main loop
// 1. take latest set of "known card ids"
// 2. crawl latest cards until we intersect with "latest known"
// 3. push jobs for downloading corresponding images
// 4. update "latest known card ids"
// 5. notify pipeline

const CARDS_DIR_ENVVAR = "CARDS_DIR"
const PIPELINE_NOTIFICATION_URL = "PIPELINE_URL"

type void struct{}

var voidVal void

const workerCount = 2
const maxKnownCardsCount = 10000
const defaultPollInterval time.Duration = 5 * 60 * 1e9

func main() {
	log.SetFlags(log.LUTC | log.Ltime)

	log.Printf("Starting up...\tVersion: %s\tGit commit: %.6s\n", version.AppVersion, version.GitCommit)

	cardsDir, ok := os.LookupEnv(CARDS_DIR_ENVVAR)
	if !ok {
		log.Printf("%s env var is not set, using default \"./db\"\n", CARDS_DIR_ENVVAR)
		cardsDir = "./db"
	} else {
		log.Printf("%s env var is set to %s, using it as cards directory\n", CARDS_DIR_ENVVAR, cardsDir)
	}

	pipelineNotificationUrlStr, ok := os.LookupEnv(PIPELINE_NOTIFICATION_URL)
	var pipelineNotificationUrl *url.URL = nil
	var err error
	if !ok {
		log.Printf("%s env var is not set, will not do pipeline notification\n", PIPELINE_NOTIFICATION_URL)
	} else {
		log.Printf("%s env var is set to %s, using it to notify pipeline\n", PIPELINE_NOTIFICATION_URL, pipelineNotificationUrlStr)
		pipelineNotificationUrl, err = url.Parse(pipelineNotificationUrlStr)
		if err != nil {
			log.Panicf("Failed to parse pipeline notification URL: %v", err)
		}
	}

	// reading card dirs

	cardDirContent, err := ioutil.ReadDir(cardsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Printf("Creating non existing dir %s", cardsDir)
			err = os.Mkdir(cardsDir, os.FileMode(0644))
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}

	var knownIDsHeap *utils.CardIDHeap = &utils.CardIDHeap{}
	for _, cardDirEntry := range cardDirContent {
		parsedID, parsedOk := strconv.ParseInt(cardDirEntry.Name(), 10, 32)
		if !cardDirEntry.IsDir() || parsedOk != nil {
			continue
		}
		heap.Push(knownIDsHeap, types.CardID(parsedID))
	}

	foundKnownIdsCount := len(*knownIDsHeap)
	log.Printf("Found %d stored cards\n", foundKnownIdsCount)

	for {
		startTime := time.Now().UTC()

		foundKnownIdsCount := len(*knownIDsHeap)
		log.Printf("Considering %d already downloaded cards\n", foundKnownIdsCount)

		if foundKnownIdsCount > maxKnownCardsCount {
			log.Printf("Will use only latest %d of known cards out of %d", maxKnownCardsCount, foundKnownIdsCount)
			*knownIDsHeap = (*knownIDsHeap)[:maxKnownCardsCount]
		}

		// log.Printf("Cards: %v\n", *knownIDsHeap)

		var knownCardsIdSet map[types.CardID]void = make(map[types.CardID]void)
		for _, v := range *knownIDsHeap {
			knownCardsIdSet[v] = voidVal
		}

		// fetching catalog
		var newDetectedCardIDs []types.CardID = nil
		if len(knownCardsIdSet) == 0 {
			// fetching only the first page
			log.Println("The card storage is empty. Fetching the first catalog page page...")
			newDetectedCardIDs, err = crawler.GetCardCatalogPage(0)
			if err != nil {
				log.Panicf("Failed to get catalog page: %v\n", err)
			}
		} else {
			// looking for
			log.Println("Fetching the catalog pages util we find the known card")
			var pageNum int = 1
		pagesLoop:
			for {
				log.Printf("Fetching catalog page %d...\n", pageNum)
				pageNewDetectedCardIDs, err := crawler.GetCardCatalogPage(0)
				if err != nil {
					log.Panicf("Failed to get catalog page: %v\n", err)
				}
				log.Printf("Got %d cards for page %d of the catalog\n", len(pageNewDetectedCardIDs), pageNum)

				if newDetectedCardIDs == nil {
					newDetectedCardIDs = pageNewDetectedCardIDs
				} else {
					newDetectedCardIDs = append(newDetectedCardIDs, pageNewDetectedCardIDs...)
				}

				// analyzing pageNewDetectedCardIDs for intersection with known IDS
				for _, newCardID := range pageNewDetectedCardIDs {
					if _, exists := knownCardsIdSet[newCardID]; exists {
						log.Printf("Found already known card %d at page %d\n", newCardID, pageNum)
						break pagesLoop
					}
				}

				pageNum += 1
			}
		}

		// finding what exactly cards are new (not previously downloaded)
		var newCardsIDs []types.CardID = make([]types.CardID, 0, len(newDetectedCardIDs))
		for _, newCardIdCandidate := range newDetectedCardIDs {
			if _, alreadyDownloaded := knownCardsIdSet[newCardIdCandidate]; !alreadyDownloaded {
				newCardsIDs = append(newCardsIDs, newCardIdCandidate)
				heap.Push(knownIDsHeap, newCardIdCandidate)
			}
		}
		log.Printf("%d new cards to download\n", len(newCardsIDs))

		var cardsJobQueue chan types.CardID = make(chan types.CardID)
		var workersWG sync.WaitGroup
		workersWG.Add(workerCount)

		runWorker := func() {
			for card := range cardsJobQueue {
				crawler.DoCardJob(card, cardsDir, pipelineNotificationUrl)
			}
			workersWG.Done()
		}

		for i := 0; i < workerCount; i++ {
			go runWorker()
		}

		for _, newCardID := range newCardsIDs {
			cardsJobQueue <- newCardID
		}
		close(cardsJobQueue)

		workersWG.Wait()
		log.Printf("All %d new cards are fetched\n", len(newCardsIDs))

		endTime := time.Now().UTC()
		elapsed := endTime.Sub(startTime)
		toWait := defaultPollInterval - elapsed
		if toWait > 0 {
			log.Printf("Sleeping for %v...", toWait)
			time.Sleep(toWait)
		}
	}
}
