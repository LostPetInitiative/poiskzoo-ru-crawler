package utils_test

import (
	"testing"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

func TestValueCanBeExtracted(t *testing.T) {
	cache := utils.NewLRUCache[string, int](10)
	cache.Set("a", 3)
	retreived, exist := cache.Get("a")
	if !exist {
		t.Error("exist must be true")
		t.Fail()
	}
	if *retreived != 3 {
		t.Errorf("expected %d, but got %d", 3, *retreived)
		t.Fail()
	}
}

func TestExtractionOfMissingValue(t *testing.T) {
	cache := utils.NewLRUCache[string, int](10)
	cache.Set("a", 3)
	retreived, exist := cache.Get("b")
	if exist {
		t.Error("exist must be false")
		t.Fail()
	}
	if retreived != nil {
		t.Error("extracted must be nil")
		t.Fail()
	}
}

func TestCapacityIsRespected(t *testing.T) {
	cache := utils.NewLRUCache[string, int](3)
	cache.Set("a", 3)
	cache.Set("b", 4)
	cache.Set("c", 5)
	cache.Set("d", 6)
	retreived, exist := cache.Get("b")
	if !exist {
		t.Error("exist must be true")
		t.Fail()
	}
	if *retreived != 4 {
		t.Error("unexpected extracted value")
		t.Fail()
	}

	// purged due to capacity
	retreived, exist = cache.Get("a")
	if exist {
		t.Error("exist must be false")
		t.Fail()
	}
	if retreived != nil {
		t.Error("unexpected extracted value")
		t.Fail()
	}
}

func TestAccessIsAccounted(t *testing.T) {
	cache := utils.NewLRUCache[string, int](3)
	cache.Set("a", 3)
	cache.Set("b", 4)
	cache.Set("c", 5)
	cache.Get("a") // now "b" is the least resently used
	cache.Set("d", 6)
	retreived, exist := cache.Get("a")
	if !exist {
		t.Error("exist must be true")
		t.Fail()
	}
	if *retreived != 3 {
		t.Error("unexpected extracted value")
		t.Fail()
	}

	// purged due to capacity
	retreived, exist = cache.Get("b")
	if exist {
		t.Error("exist must be false")
		t.Fail()
	}
	if retreived != nil {
		t.Error("unexpected extracted value")
		t.Fail()
	}
}
