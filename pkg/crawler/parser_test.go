package crawler

import (
	"log"
	"os"
	"testing"

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

func TestExtractSpeciesFromPetCard(t *testing.T) {
	testCases := []struct {
		path     string
		expected types.Species
	}{
		{"./testdata/164931.html.dump", types.Dog},
		{"./testdata/164921.html.dump", types.Cat},
		{"./testdata/164923.html.dump", types.Cat},
		{"./testdata/164929.html.dump", types.Dog},
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
