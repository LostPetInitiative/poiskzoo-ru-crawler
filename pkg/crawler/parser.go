package crawler

import (
	"fmt"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func ParseHtmlContent(htmlContent string) *html.Node {
	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse HTML file: %v", err))
	}

	return doc
}

func ExtractCardUrlsFromDocument(doc *html.Node) []string {
	nodes, err := htmlquery.QueryAll(doc, "//div[contains(@class, 'pzplitkadiv')]//div[contains(@class, 'pzplitkalink')]/a")
	if err != nil {
		panic(`not a valid XPath expression.`)
	}

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

func ExtractSpeciesFromCardPage(doc *html.Node) types.Species {
	node := htmlquery.FindOne(doc, "//h1[contains(@class, 'con_heading')]")
	dataText := strings.ToLower(node.FirstChild.Data)
	switch {
	case strings.Contains(dataText, "собака"), strings.Contains(dataText, "пес"):
		return types.Dog
	case strings.Contains(dataText, "кот"), strings.Contains(dataText, "кошка"):
		return types.Cat
	default:
		panic(fmt.Sprintf("Can't extract species type"))
	}

}
