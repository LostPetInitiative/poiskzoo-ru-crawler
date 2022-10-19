package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version"
)

type Nominatim struct {
	baseUrl    *url.URL
	mutex      *sync.Mutex
	httpClient *http.Client
}

func (n *Nominatim) Geocode(toponym string) (*GeoCoords, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	requestFullURL := fmt.Sprintf("%s?q=%s&format=jsonv2", n.baseUrl, url.QueryEscape(toponym))

	req, err := http.NewRequest("GET", requestFullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Set("User-Agent", fmt.Sprintf("LostPetInitiative:poiskzoo-crawler / %s:%.6s (https://kashtanka.pet/)", version.AppVersion, version.GitCommit))
	resp, err := n.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
	return &Nominatim{serviceUrl, &sync.Mutex{}, &http.Client{}}
}

const openStreetMapsNominatimURL string = "https://nominatim.openstreetmap.org/search.php"

// Constructs a Nominatim API client that uses OpenStreetMap(OSM) free public instance
func NewOpenStreetMapsNominatim() *Nominatim {
	url, err := url.Parse(openStreetMapsNominatimURL)
	if err != nil {
		panic("Failed to parse openStreetMapsNominatimURL")
	}
	return &Nominatim{url, &sync.Mutex{}, &http.Client{}}
}

type FoundToponymJSON struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
