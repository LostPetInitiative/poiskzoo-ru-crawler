package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

// OSM Nominatim Usage Policy (aka Geocoding Policy)
// Requirements
// * No heavy uses (an absolute maximum of 1 request per second).
// * Provide a valid HTTP Referer or User-Agent identifying the application (stock User-Agents as set by http libraries will not do).
// * Clearly display attribution as suitable for your medium.
// * Data is provided under the ODbL license which requires to share alike (although small extractions are likely to be covered by fair usage / fair dealing).

// Bulk Geocoding
//
// * limit your requests to a single thread
// * limited to 1 machine only, no distributed scripts (including multiple Amazon EC2 instances or similar)
// * Results must be cached on your side. Clients sending repeatedly the same query may be classified as faulty and blocked.

// see https://operations.osmfoundation.org/policies/nominatim/

type Nominatim struct {
	baseUrl       *url.URL
	mutex         *sync.Mutex
	latestRequest *time.Time
	// setting this to 0 disables throttling
	minIntervalBetweenRequest time.Duration
}

func (n *Nominatim) Geocode(toponym string) (*GeoCoords, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	elapsed := time.Now().UTC().Sub(*n.latestRequest)
	toWait := n.minIntervalBetweenRequest - elapsed
	// log.Printf("Time to wait %v\n", toWait)
	if toWait > 0 {
		time.Sleep(toWait)
	}
	now := time.Now().UTC()
	n.latestRequest = &now

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

func NewNominatim(serviceUrl *url.URL, minIntervalBetweenRequest time.Duration) *Nominatim {
	zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	return &Nominatim{
		baseUrl:                   serviceUrl,
		mutex:                     &sync.Mutex{},
		latestRequest:             &zero,
		minIntervalBetweenRequest: minIntervalBetweenRequest,
	}
}

const openStreetMapsNominatimURL string = "https://nominatim.openstreetmap.org/search.php"

// Constructs a Nominatim API client that uses OpenStreetMap(OSM) free public instance
func NewOpenStreetMapsNominatim() *Nominatim {
	url, err := url.Parse(openStreetMapsNominatimURL)
	if err != nil {
		panic("Failed to parse openStreetMapsNominatimURL")
	}
	return NewNominatim(url, time.Duration(1e9))
}

type FoundToponymJSON struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
