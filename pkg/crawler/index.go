package crawler

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

const poiskZooBaseURL string = "https://poiskzoo.ru"

func GetCardCatalogPage(pageNum int) ([]Card, error) {
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
	var cards []Card = ExtractCardsFromCatalogDocument(parsedNode)

	return cards, nil
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
