package crawler

import (
	"testing"
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
