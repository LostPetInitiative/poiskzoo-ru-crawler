package geocoding

import (
	"math"
	"testing"
)

func TestOSMGeocoder(t *testing.T) {
	var coder Geocoder = NewOpenStreetMapsNominatim()

	result, err := coder.Geocode("Таруса, пл. Ленина")

	var expected GeoCoords = GeoCoords{
		Lat: 54.7291584,
		Lon: 37.1807652,
	}

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	latDiff := math.Abs(expected.Lat - result.Lat)
	if latDiff > 1e-3 {
		t.Logf("Lat diff is %v which is too high", latDiff)
		t.Fail()
	}

	lonDiff := math.Abs(expected.Lon - result.Lon)
	if lonDiff > 1e-3 {
		t.Logf("Lon diff is %v which is too high", lonDiff)
		t.Fail()
	}
}
