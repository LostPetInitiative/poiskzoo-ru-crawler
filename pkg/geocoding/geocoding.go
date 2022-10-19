package geocoding

type GeoCoords struct {
	Lat, Lon float64
}

type Geocoder interface {
	// if error is nil, GeoCoords must be not nil
	Geocode(toponym string) (*GeoCoords, error)
}
