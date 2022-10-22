package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

type Nominatim struct {
	baseUrl *url.URL
	mutex   *sync.Mutex
}

func (n *Nominatim) Geocode(toponym string) (*GeoCoords, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	requestFullURLstr := fmt.Sprintf("%s?q=%s&format=jsonv2", n.baseUrl, url.QueryEscape(toponym))
	requestFullURL, err := url.Parse(requestFullURLstr)
	if err != nil {
		return nil, err
	}

	resp, err := utils.HttpGet(requestFullURL, types.JsonMimeType)
	if err != nil {
		return nil, err
	}
	body := resp.Body

	var foundToponyms []FoundToponymJSON
	err = json.Unmarshal(body, &foundToponyms)
	if err != nil {
		return nil, err
	}

	if len(foundToponyms) > 0 {
		first := foundToponyms[0]
		parsedLat, err := strconv.ParseFloat(first.Lat, 64)
		if err != nil {
			return nil, err
		}
		parsedLon, err := strconv.ParseFloat(first.Lon, 64)
		if err != nil {
			return nil, err
		}
		return &GeoCoords{
			Lat: parsedLat,
			Lon: parsedLon,
		}, nil
	}
	return nil, errors.New("Geocoder failed to find any coordinates")
}

func NewNominatim(serviceUrl *url.URL) *Nominatim {
	return &Nominatim{serviceUrl, &sync.Mutex{}}
}

const openStreetMapsNominatimURL string = "https://nominatim.openstreetmap.org/search.php"

// Constructs a Nominatim API client that uses OpenStreetMap(OSM) free public instance
func NewOpenStreetMapsNominatim() *Nominatim {
	url, err := url.Parse(openStreetMapsNominatimURL)
	if err != nil {
		panic("Failed to parse openStreetMapsNominatimURL")
	}
	return &Nominatim{url, &sync.Mutex{}}
}

type FoundToponymJSON struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
