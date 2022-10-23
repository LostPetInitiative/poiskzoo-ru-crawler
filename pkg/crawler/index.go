package crawler

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

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

	resp, err := utils.HttpGetHtml(effectiveUrl)
	if err != nil {
		return nil, err
	}
	body := resp.Body

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

type PetCard struct {
	ID        types.CardID
	Species   types.Species
	SexSpec   types.Sex
	City      string
	Address   string
	EventTime time.Time
	EventType types.EventType
	Comment   string
	ImagesURL *url.URL
}

func GetPetCard(card types.CardID) (*PetCard, error) {
	cardUrl, err := url.Parse(fmt.Sprintf("%s/%d", poiskZooBaseURL, card))
	if err != nil {
		return nil, err
	}
	resp, err := utils.HttpGetHtml(cardUrl)
	if err != nil {
		return nil, err
	}

	// TODO: account for content-type header? seams that this header is misleading for poiskzoo.ru. The true encoding is utf-8, while content encoding headers says "windows-1251"
	// TODO: parse encoding
	parsed := ParseHtmlContent(string(resp.Body))

	cityWithAddress := ExtractAddressFromCardPage(parsed)

	nowUtc := time.Now().UTC()
	today := time.Date(nowUtc.Year(), nowUtc.Month(), nowUtc.Day(), 0, 0, 0, 0, time.UTC)

	return &PetCard{
		ID:        card,
		Species:   ExtractSpeciesFromCardPage(parsed),
		SexSpec:   ExtractAnimalSexSpecFromCardPage(parsed),
		City:      cityWithAddress.City,
		Address:   cityWithAddress.Address,
		EventTime: ExtractEventDateFromCardPage(parsed, today),
		EventType: ExtractCardTypeFromCardPage(parsed),
		Comment:   ExtractCommentFromCardPage(parsed),
		ImagesURL: ExtractSmallPhotoUrlFromCardPage(parsed),
	}, nil

}
