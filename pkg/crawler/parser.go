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
		panic("Can't extract species type")
	}
}

func ExtractCardTypeFromCardPage(doc *html.Node) types.CardType {
	node := htmlquery.FindOne(doc, "//h1[contains(@class, 'con_heading')]")
	dataText := strings.ToLower(node.FirstChild.Data)
	switch {
	case strings.Contains(dataText, "найден"):
		return types.Found
	case strings.Contains(dataText, "пропал"):
		return types.Lost
	default:
		panic("Can't extract card type")
	}
}

func ExtractAddressFromCardPage(doc *html.Node) string {
	node := htmlquery.FindOne(doc, "//h1[contains(@class, 'con_heading')]")
	dataText := node.FirstChild.Data
	words := strings.Split(dataText, " ")
	if len(words) < 1 {
		panic("Heading does not contain enough data (city name at the end?)")
	}
	lastWord := words[len(words)-1]

	regionNode := htmlquery.FindOne(doc, "//strong[contains(text(), 'Район где')]")
	if regionNode == nil {
		regionNode = htmlquery.FindOne(doc, "//strong[contains(text(), 'Адрес где')]")
	}

	if regionNode == nil {
		panic("Can't find address/region element on the page")
	}

	text := make([]string, 1)
	text[0] = lastWord // this usually contains the City

	var curNode *html.Node = regionNode

	for {
		sib := curNode.NextSibling
		if sib.Type == html.ElementNode && sib.Data == "strong" {
			break
		}
		if sib == nil {
			break
		}
		if sib.Type == html.TextNode {
			trimmed := strings.TrimSpace(sib.Data)
			if len(trimmed) > 0 {
				text = append(text, trimmed)
			}
		}
		curNode = sib
	}
	textJoined := strings.Join(text, ", ")

	return textJoined
}
