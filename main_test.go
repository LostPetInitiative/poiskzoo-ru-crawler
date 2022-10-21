package main

import (
	"io/ioutil"
	"testing"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/crawler"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/geocoding"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/storage"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

func TestFullCardDownload(t *testing.T) {
	card, err := crawler.GetPetCard(types.CardID(164971))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	image, err := utils.HttpGet(card.ImagesURL, types.HtmlAcceptAnyMimeType)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	jsonCard := storage.NewCardJSON(card,
		&geocoding.GeoCoords{Lat: 10.0, Lon: 20.0},
		"hardcoded",
		image.Body,
		image.ContentType)
	serialized := jsonCard.JsonSerialize()

	expectedBytes, err := ioutil.ReadFile("./testdata/164971.json")
	expected := string(expectedBytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected != serialized {
		for i := 0; i < len(expected) && i < len(serialized); i++ {
			if expected[i] != serialized[i] {
				t.Errorf("Expected != actual. Diff is at byte idx: %d", i)
				t.FailNow()
			}
		}
	}

}
