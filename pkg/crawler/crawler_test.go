package crawler

import (
	"testing"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

type issue13StorageStub struct {
}

func (s *issue13StorageStub) IsCardExist(card types.CardID) bool {
	return false
}

func (s *issue13StorageStub) SaveCard(petCard *PetCard, jsonCard *CardJSON, fetchedImage *utils.HttpFetchResult) {
	// there must be no image here
	if fetchedImage != nil {
		panic("Image must be nil")
	}
}

func TestIssue13(t *testing.T) {
	// attempt to follow image link for photo download redirects to poiskzoo main page.
	// this must be handled as absence of the image
	var storage LocalCardStorage = &issue13StorageStub{}
	crawler := NewCrawler(&storage, nil)

	crawler.DoCardJob(types.CardID(165457))
}
