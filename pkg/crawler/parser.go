package crawler

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
)

func ExtractCardUrlsFromCatalogPage(htmlContent string) []string {
	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse HTML file: %v", err))
	}

	nodes, err := htmlquery.QueryAll(doc, "//div[contains(@class, 'pzplitkadiv')]//div[contains(@class, 'pzplitkalink')]/a")
	if err != nil {
		panic(`not a valid XPath expression.`)
	}

	log.Printf("nodes count %d\n", len(nodes))

	res := make([]string, len(nodes))

	for i, n := range nodes {
		for _, a := range n.Attr {
			if a.Key == "href" {
				res[i] = a.Val
				break
			}
		}
	}

	return res
}
