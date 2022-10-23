package geocoding

import (
	"log"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

type cacheRes struct {
	fst *GeoCoords
	snd error
}

type LRUCacheDecorator struct {
	target *Geocoder
	cache  *utils.LRUCache[string, cacheRes]
}

func NewLRUCacheDecorator(target *Geocoder, cacheCapacity int) *LRUCacheDecorator {
	return &LRUCacheDecorator{
		target: target,
		cache:  utils.NewLRUCache[string, cacheRes](cacheCapacity),
	}
}

func (c *LRUCacheDecorator) Geocode(toponym string) (*GeoCoords, error) {
	cached, exists := c.cache.Get(toponym)
	if exists {
		log.Printf("Cache hit geocoding \"%s\"\n", toponym)
		return cached.fst, cached.snd
	}

	lookupRes, err := (*c.target).Geocode(toponym)

	c.cache.Set(toponym, cacheRes{lookupRes, err})

	return lookupRes, err
}
