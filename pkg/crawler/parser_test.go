package crawler

import (
	"log"
	"os"
	"testing"
)

func TestExtractCardUrlsFromCatalogPage(t *testing.T) {
	fileContent, err := os.ReadFile("./testdata/catalog.html.dump")
	if err != nil {
		log.Fatal(err)
	}
	catalogHtml := string(fileContent)

	extractedUrls := ExtractCardUrlsFromCatalogPage(catalogHtml)
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
