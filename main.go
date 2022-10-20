package main

import (
	"container/heap"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"

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

type void struct{}

var voidVal void

func main() {
	log.SetFlags(log.LUTC | log.Ltime)

	cardsDir, ok := os.LookupEnv(CARDS_DIR_ENVVAR)
	if !ok {
		log.Printf("%s env var is not set, using default \"./db\"\n", CARDS_DIR_ENVVAR)
		cardsDir = "./db"
	} else {
		log.Printf("%s env var is set to %s, using it as cards directory\n", CARDS_DIR_ENVVAR, cardsDir)
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

	var knownIDS map[int]void = make(map[int]void, 0)
	log.Printf("Found %d stored cards", len(knownIDS))

	//	var knownCards []CardID

	// TODO: load from disk
	//	knownCards = make([]CardID, 0) // empty for now
	log.Printf("Starting up...\nVersion: %s\nGit commit: %.6s\n", version.AppVersion, version.GitCommit)
}
