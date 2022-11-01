package crawler

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/geocoding"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

// TODO: inject implementation?
var nominatim geocoding.Geocoder = geocoding.NewOpenStreetMapsNominatim()
var cachedNominatim geocoding.Geocoder = geocoding.NewLRUCacheDecorator(&nominatim, 128)

type LocalCardStorage interface {
	IsCardExist(card types.CardID) bool
	SaveCard(petCard *PetCard, jsonCard *CardJSON, fetchedImage *utils.HttpFetchResult)
}

type Crawler struct {
	cardStorage     *LocalCardStorage
	notificationUrl *url.URL
}

func NewCrawler(localStorage *LocalCardStorage, notificationUrl *url.URL) *Crawler {
	return &Crawler{
		cardStorage:     localStorage,
		notificationUrl: notificationUrl,
	}
}

// Download card, save it to disk, post it to HTTP (kafka REST API) if notification url is not nil
func (c *Crawler) DoCardJob(card types.CardID) {
	cardJobFailurePrinter := func() {
		if a := recover(); a != nil {
			log.Printf("%d:\tPanic during fetching of card %v", card, a)
			panic(a)
		}
	}
	defer cardJobFailurePrinter()

	// workaround for paid promotion
	// TODO: do something smarter
	if (*c.cardStorage).IsCardExist(card) {
		log.Printf("%d:\t Card dir already exists. Consider it as processed. skipping it\n", card)
		return
	}

	log.Printf("%d:\tFetching card...\n", card)
	fetchedCard, err := GetPetCard(card)
	if err != nil {
		log.Panicf("%d:\tFailed to download card: %v\n", card, err)
	}
	log.Printf("%d:\tDownloaded card\n", card)
	var fetchedImage *utils.HttpFetchResult = DownloadImage(fetchedCard.ImagesURL, fmt.Sprintf("%d:\t", card))

	var locationSpecFormats []string = []string{
		fmt.Sprintf("Россия, г. %s, %s", fetchedCard.City, fetchedCard.Address),
		fmt.Sprintf("Россия, %s, %s", fetchedCard.City, fetchedCard.Address),
		fmt.Sprintf("Россия, г. %s", fetchedCard.City),
		fmt.Sprintf("Россия, %s", fetchedCard.City),
		fmt.Sprintf("г. %s, %s", fetchedCard.City, fetchedCard.Address),
		fmt.Sprintf("%s, %s", fetchedCard.City, fetchedCard.Address),
		fmt.Sprintf("г. %s", fetchedCard.City),
		fetchedCard.City,
	}
	var geoCoords *geocoding.GeoCoords
	for _, locationSpec := range locationSpecFormats {
		log.Printf("%d:\tTrying to geocode \"%s\"...\n", card, locationSpec)
		coords, err := cachedNominatim.Geocode(locationSpec)
		if err == nil {
			log.Printf("%d:\tSuccessfully geocoded \"%s\" as lat:%f lon:%f\n", card, locationSpec, coords.Lat, coords.Lon)
			geoCoords = coords
			break
		}
	}

	var imageBytes []byte
	var imageMime string
	if fetchedImage != nil && strings.Contains(fetchedImage.ContentType, "image") {
		imageBytes = fetchedImage.Body
		imageMime = fetchedImage.ContentType
	} else {
		fetchedImage = nil // fetch data from image URL is not an image
	}

	jsonCard := NewCardJSON(fetchedCard,
		geoCoords,
		"Геокодер OSM Moninatim",
		imageBytes,
		imageMime)
	serialized := jsonCard.JsonSerialize()

	if c.notificationUrl != nil {
		// doing notification
		log.Printf("%d:\tSending snapshot to pipeline...\n\n", card)
		_, err = utils.HttpPost(c.notificationUrl, types.JsonMimeType, []byte(serialized))
		if err != nil {
			log.Panicf("%d:Failed to notify pipeline %v\t\n", card, err)
		} else {
			log.Printf("%d:\tSuccessfully notified the pipeline\n", card)
		}
	} else {
		log.Printf("%d:\tSkipped pipeline notification, as no notification URL is set\n", card)
	}

	(*c.cardStorage).SaveCard(fetchedCard, jsonCard, fetchedImage)
}

func DownloadImage(imageURL *url.URL, logPrefix string) *utils.HttpFetchResult {
	var fetchedImage *utils.HttpFetchResult
	var err error
	if imageURL != nil {
		log.Printf("%sDownloading image %v\n", logPrefix, *imageURL)
		fetchedImage, err = utils.HttpGet(imageURL, "*/*")
		if err != nil {
			log.Panicf("%sFailed to download image for card: %v\n", logPrefix, err)
		}
		log.Printf("%sDownloaded image (%d bytes; mime %s)\n", logPrefix, len(fetchedImage.Body), fetchedImage.ContentType)
	}
	return fetchedImage
}
