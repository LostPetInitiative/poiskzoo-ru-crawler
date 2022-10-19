package crawler

import (
	"fmt"
	"log"
	"net/url"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
)

const poiskZooBaseURL string = "https://poiskzoo.ru"

type PoiskZooCatalogPage struct {
}

func (p *PoiskZooCatalogPage) getLatestKnownCards() []types.CardID {
	var pageNum = 0
	for {
		effectiveUrlStr := fmt.Sprintf("%s/poteryashka/page-%d", poiskZooBaseURL, pageNum)
		effectiveUrl, err := url.Parse(effectiveUrlStr)
		if err != nil {
			log.Fatalf("Unable to parse URL: %s (%v)", effectiveUrlStr, effectiveUrl)
		}
		pageNum += 1
	}
}
