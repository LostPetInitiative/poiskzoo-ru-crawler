package crawler

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

const poiskZooBaseURL string = "https://poiskzoo.ru"

func GetCardCatalogPage(pageNum int) ([]types.CardID, error) {
	effectiveUrlStr := fmt.Sprintf("%s/poteryashka/page-%d", poiskZooBaseURL, pageNum)
	effectiveUrl, err := url.Parse(effectiveUrlStr)
	if err != nil {
		log.Fatalf("Unable to parse URL: %s (%v)", effectiveUrlStr, effectiveUrl)
	}

	body, err := utils.HttpGet(effectiveUrl, "text/html")
	if err != nil {
		return nil, err
	}

	parsedNode := ParseHtmlContent(string(body))
	var urls []string = ExtractCardUrlsFromDocument(parsedNode)

	var result []types.CardID = make([]types.CardID, len(urls))
	for i, url := range urls {
		// urls are like "/bijsk/propala-koshka/162257"
		lastIdx := strings.LastIndex(url, "/")
		if lastIdx == -1 {
			panic(fmt.Sprintf("card URL in not in supported format: %q", url))
		}
		cardIdStr := url[lastIdx+1:]
		cardID, err := strconv.ParseInt(cardIdStr, 10, 32)
		if err != nil {
			return nil, err
		}
		result[i] = types.CardID(cardID)
	}
	return result, nil

}
