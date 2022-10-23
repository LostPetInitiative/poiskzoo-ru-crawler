package crawler

import (
	"os"
	"testing"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/geocoding"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

func TestGetCardCatalogPage(t *testing.T) {
	cards, err := GetCardCatalogPage(0)
	if err != nil {
		t.Errorf("Got error while getting card catalog: %v", err)
		t.FailNow()
	}
	t.Logf("Got %d cards", len(cards))
	if len(cards) == 0 {
		t.Error("Got empty set of cards")
		t.FailNow()
	}
}

func TestFullCardDownload(t *testing.T) {
	card, err := GetPetCard(types.CardID(164971))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	image, err := utils.HttpGet(card.ImagesURL, types.AnyMimeType)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	jsonCard := NewCardJSON(card,
		&geocoding.GeoCoords{Lat: 10.0, Lon: 20.0},
		"hardcoded",
		image.Body,
		image.ContentType)
	serialized := jsonCard.JsonSerialize()

	expectedBytes, err := os.ReadFile("./testdata/164971.json")
	expected := string(expectedBytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected != serialized {
		for i := 0; i < len(expected) && i < len(serialized); i++ {
			if expected[i] != serialized[i] {
				t.Errorf("Expected != actual. Diff is at byte idx: %d\n", i)
				t.Errorf("Actual: %v\n", serialized)
				t.FailNow()
			}
		}
	}

}
