package crawler

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/geocoding"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

// TODO: inject implementation?
var nominatim geocoding.Geocoder = geocoding.NewOpenStreetMapsNominatim()
var cachedNominatim geocoding.Geocoder = geocoding.NewLRUCacheDecorator(&nominatim, 128)

// Download card, save it to disk, post it to HTTP (kafka REST API) if notification url is not nil
func DoCardJob(card types.CardID, cardsDir string, notificationUrl *url.URL) {
	cardJobFailurePrinter := func() {
		if a := recover(); a != nil {
			log.Printf("%d:\tPanic during fetching of card %v", card, a)
			panic(a)
		}
	}
	defer cardJobFailurePrinter()

	log.Printf("%d:\tFetching card...\n", card)
	fetchedCard, err := GetPetCard(card)
	if err != nil {
		log.Panicf("%d:\tFailed to download card: %v\n", card, err)
	}
	log.Printf("%d:\tDownloaded card\n", card)
	var fetchedImage *utils.HttpFetchResult
	if fetchedCard.ImagesURL != nil {
		fetchedImage, err = utils.HttpGet(fetchedCard.ImagesURL, "*/*")
		if err != nil {
			log.Panicf("%d:\tFailed to download image for card: %v\n", card, err)
		}
		log.Printf("%d:\tDownloaded image for card (%d bytes; mime %s)\n", card, len(fetchedImage.Body), fetchedImage.ContentType)
	}

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
	if fetchedImage != nil {
		imageBytes = fetchedImage.Body
		imageMime = fetchedImage.ContentType
	}

	jsonCard := NewCardJSON(fetchedCard,
		geoCoords,
		"Геокодер OSM Moninatim",
		imageBytes,
		imageMime)
	serialized := jsonCard.JsonSerialize()

	if notificationUrl != nil {
		// doing notification
		log.Printf("%d:\tSending snapshot to pipeline...\n\n", card)
		_, err = utils.HttpPost(notificationUrl, types.JsonMimeType, []byte(serialized))
		if err != nil {
			log.Panicf("%d:Failed to notify pipeline %v\t\n", card, err)
		} else {
			log.Printf("%d:\tSuccessfully notified the pipeline\n", card)
		}
	} else {
		log.Printf("%d:\tSkipped pipeline notification, as no notification URL is set\n", card)
	}

	log.Printf("%d:\tDumping card to disk...\n", card)

	cardDir := path.Join(cardsDir, fmt.Sprintf("%d", card))
	err = os.Mkdir(cardDir, 0644)
	if err != nil {
		log.Panicf("%d:\tFailed to create card dir: %v", card, cardsDir)
	}

	// replacing embedded base64 image with file reference
	var imageFileName string
	if imageBytes != nil {
		var imageFileExt string
		switch strings.ToLower(fetchedImage.ContentType) {
		case "image/jpeg":
			imageFileExt = "jpg"
		default:
			imageFileExt = strings.TrimPrefix(fetchedImage.ContentType, "image/")
		}
		imageFileName = fmt.Sprintf("image.%s", imageFileExt)

		jsonCard.Images = []EncodedImageJSON{{Type: "file", Data: imageFileName}}
		serialized = jsonCard.JsonSerialize()
	}

	cardFilePath := path.Join(cardDir, "card.json")
	err = os.WriteFile(cardFilePath, []byte(serialized), 0644)
	if err != nil {
		log.Panicf("%d:\t%v\n", card, err)
	} else {
		log.Printf("%d:\tJSON card saved to disk\n", card)
	}
	if imageBytes != nil {
		imageFilePath := path.Join(cardDir, imageFileName)
		err = os.WriteFile(imageFilePath, fetchedImage.Body, 0644)
		if err != nil {
			log.Panicf("%d:\t%v\n", card, err)
		} else {
			log.Printf("%d:\timage file saved to disk\n", card)
		}
	}
}
