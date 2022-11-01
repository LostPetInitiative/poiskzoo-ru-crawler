package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/crawler"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

type DirectoryCardStorage struct {
	cardsDir string
}

func NewDirectoryCardStorage(cardDir string) *DirectoryCardStorage {
	return &DirectoryCardStorage{cardsDir: cardDir}
}

func (d *DirectoryCardStorage) getCardDir(card types.CardID) string {
	return path.Join(d.cardsDir, fmt.Sprintf("%d", card))
}

func (d *DirectoryCardStorage) IsCardExist(card types.CardID) bool {
	_, err := os.Stat(d.getCardDir(card))
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}

func (d *DirectoryCardStorage) SaveCard(petCard *crawler.PetCard, jsonCard *crawler.CardJSON, fetchedImage *utils.HttpFetchResult) {
	card := petCard.ID
	log.Printf("%d:\tDumping card to disk...\n", card)
	cardDir := d.getCardDir(card)

	err := os.Mkdir(cardDir, 0644)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		log.Panicf("%d:\tFailed to create card dir: %v", card, d.cardsDir)
	}

	// replacing embedded base64 image with file reference
	var imageFileName string

	if fetchedImage != nil {
		var imageFileExt string
		switch strings.ToLower(fetchedImage.ContentType) {
		case "image/jpeg":
			imageFileExt = "jpg"
		default:
			imageFileExt = strings.TrimPrefix(fetchedImage.ContentType, "image/")
		}
		imageFileName = fmt.Sprintf("image.%s", imageFileExt)

		jsonCard.Images = []crawler.EncodedImageJSON{{Type: "file", Data: imageFileName}}

	}
	var serialized string = jsonCard.JsonSerialize()

	cardFilePath := path.Join(cardDir, "card.json")
	err = os.WriteFile(cardFilePath, []byte(serialized), 0644)
	if err != nil {
		log.Panicf("%d:\t%v\n", card, err)
	} else {
		log.Printf("%d:\tJSON card saved to disk\n", card)
	}
	if fetchedImage != nil {
		imageFilePath := path.Join(cardDir, imageFileName)
		err = os.WriteFile(imageFilePath, fetchedImage.Body, 0644)
		if err != nil {
			log.Panicf("%d:\t%v\n", card, err)
		} else {
			log.Printf("%d:\timage file saved to disk\n", card)
		}
	}
}
