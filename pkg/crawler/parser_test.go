package crawler

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
)

func TestExtractCardUrlsFromCatalogPage(t *testing.T) {
	fileContent, err := os.ReadFile("./testdata/catalog.html.dump")
	if err != nil {
		log.Fatal(err)
	}
	catalogHtml := string(fileContent)

	extractedUrls := ExtractCardUrlsFromDocument(ParseHtmlContent(catalogHtml))
	const expectedCount int = 52
	if len(extractedUrls) != expectedCount {
		t.Logf("Expected to extract %d card IDs but extracted %d", expectedCount, len(extractedUrls))
		t.Fail()
	}

	urlMap := make(map[string]int)
	for i, url := range extractedUrls {
		urlMap[url] = i
	}

	_, exists := urlMap["/stavropol/najdena-sobaka/164833"]
	if !exists {
		t.Logf("Did not find expected URL in the result set")
		t.Fail()
	}
}

func TestExtractSpeciesFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path     string
		expected types.Species
	}{
		{"./testdata/164921.html.dump", types.Cat},
		{"./testdata/164923.html.dump", types.Cat},
		{"./testdata/164929.html.dump", types.Dog},
		{"./testdata/164931.html.dump", types.Dog},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedSpecies types.Species = ExtractSpeciesFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedSpecies != testCase.expected {
			t.Logf("Wrong species extracted for %s. Expected %v, but got %v", testCase.path, testCase.expected, extractedSpecies)
			t.Fail()
		}
	}
}

func TestExtractCardTypeFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path     string
		expected types.CardType
	}{
		{"./testdata/164921.html.dump", types.Lost},
		{"./testdata/164923.html.dump", types.Found},
		{"./testdata/164929.html.dump", types.Found},
		{"./testdata/164931.html.dump", types.Lost},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedType types.CardType = ExtractCardTypeFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedType != testCase.expected {
			t.Logf("Wrong card type extracted for %s. Expected %v, but got %v", testCase.path, testCase.expected, extractedType)
			t.Fail()
		}
	}
}

func TestExtractAddressFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path, address string
	}{
		{"./testdata/164921.html.dump", "Оренбург, Центральный"},
		{"./testdata/164923.html.dump", "Орехово-Зуево, Демихово"},
		{"./testdata/164929.html.dump", "Владивосток, Владивосток, район Арт-пляжа."},
		{"./testdata/164931.html.dump", "Сургут, г. Сургут, пр. Пролетарский 8/1-8/2"},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedAddress string = ExtractAddressFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedAddress != testCase.address {
			t.Logf("Wrong card type extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.address, extractedAddress)
			t.Fail()
		}
	}
}

func TestExtractEventTimeFromPetCardPage(t *testing.T) {
	utcLoc, err := time.LoadLocation("UTC")
	if err != nil {
		panic("Can't load UTC location")
	}
	today := time.Date(2022, 10, 17, 0, 0, 0, 0, utcLoc)

	testCases := []struct {
		path string
		time time.Time
	}{
		{"./testdata/164793.html.dump", time.Date(2022, 10, 13, 0, 0, 0, 0, utcLoc)},
		{"./testdata/164921.html.dump", time.Date(2022, 10, 16, 22, 27, 0, 0, utcLoc)},
		{"./testdata/164923.html.dump", time.Date(2022, 10, 17, 0, 10, 0, 0, utcLoc)},
		{"./testdata/164929.html.dump", time.Date(2022, 10, 17, 6, 52, 0, 0, utcLoc)},
		{"./testdata/164931.html.dump", time.Date(2022, 10, 17, 7, 45, 0, 0, utcLoc)},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedEventTime time.Time = ExtractEventDate(ParseHtmlContent(catalogHtml), today)
		if extractedEventTime != testCase.time {
			t.Logf("Wrong event time extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.time, extractedEventTime)
			t.Fail()
		}
	}
}
