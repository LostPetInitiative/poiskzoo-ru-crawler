package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/crawler"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/geocoding"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

type LocationJSON struct {
	Address    string  `json:"Address"`
	Lat        float64 `json:"Lat"`
	Lon        float64 `json:"Lon"`
	Provenance string  `json:"CoordsProvenance"`
}

type ContactInfoJSON struct {
	Comment string   `json:"Comment"`
	Tel     []string `json:"Tel"`
	Website []string `json:"Website"`
	Email   []string `json:"Email"`
	Name    string   `json:"Name"`
}

type EncodedImageJSON struct {
	Type string `json:"type"`
	// base64 encoded byte[]
	Data string `json:"data"`
}

type CardJSON struct {
	Uid                 string             `json:"uid"`
	Species             string             `json:"animal"`
	Location            *LocationJSON      `json:"location"`
	EventTime           time.Time          `json:"event_time"`
	EventTimeProvenance string             `json:"event_time_provenance"`
	EventType           string             `json:"card_type"`
	ContactInfo         *ContactInfoJSON   `json:"contact_info"`
	Images              []EncodedImageJSON `json:"images"`
	ProvenanceURL       string             `json:"provenance_url"`
}

func (c *CardJSON) JsonSerialize() string {
	encoded, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panicf("Failed to JSON encode %v: %v", c, err)
	}
	return string(encoded)
}

func NewCardJSON(
	card *crawler.PetCard,
	geoCoords *geocoding.GeoCoords,
	geoCoordsProvenance string,
	imageData []byte,
	imageMime string) *CardJSON {

	var emptyStrSlice []string = make([]string, 0)

	return &CardJSON{
		Uid:     fmt.Sprintf("poiskzooru_%d", card.ID),
		Species: card.Species.String(),
		Location: &LocationJSON{
			Address:    fmt.Sprintf("%s, %s", card.City, card.Address),
			Lat:        geoCoords.Lat,
			Lon:        geoCoords.Lon,
			Provenance: geoCoordsProvenance,
		},
		EventTime:           card.EventTime,
		EventTimeProvenance: "Указано на сайте poiskzoo.ru",
		EventType:           card.EventType.String(),
		ContactInfo: &ContactInfoJSON{
			Comment: card.Comment,
			Tel:     emptyStrSlice,
			Website: emptyStrSlice,
			Email:   emptyStrSlice,
			Name:    "",
		},
		Images:        []EncodedImageJSON{*EncodeImage(imageData, imageMime)},
		ProvenanceURL: fmt.Sprintf("https://poiskzoo.ru/%d", card.ID),
	}

}

func EncodeImage(data []byte, mimeType string) *EncodedImageJSON {
	var typeString string
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		typeString = "jpg"
	case "image/png":
		typeString = "png"
	default:
		log.Panicf("Unsupported image mime type: %s", mimeType)
	}

	return &EncodedImageJSON{
		Data: utils.Base64Encode(data),
		Type: typeString,
	}
}
