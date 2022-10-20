package crawler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

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

// Returns relative URL from cards found on the catalog page
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

// Parses time in HH:mm format as Duration since midnight
func parseTime(timeStr string) time.Duration {
	if timeStr[2] != ':' && len(timeStr) != 5 {
		panic(fmt.Sprintf("Time is supposed to be in HH:MM format, but instead got %s", timeStr))
	}
	hours, err := strconv.ParseInt(timeStr[0:2], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("Can't parse hours in time string %s", timeStr))
	}
	minutes, err := strconv.ParseInt(timeStr[3:5], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("Can't parse minutes in time string %s", timeStr))
	}
	return time.Duration((hours*60 + minutes) * 60 * 1e9)
}

// today - is midnight of some date (UTC)
func ExtractEventDateFromCardPage(doc *html.Node, today time.Time) time.Time {
	node := htmlquery.FindOne(doc, "//span[contains(@class, 'bd_item_date')]")
	text := strings.TrimSpace(node.FirstChild.Data)
	lowerText := strings.ToLower(text)
	switch {
	case strings.HasPrefix(lowerText, "вчера в"):
		timeStr := strings.TrimSpace(strings.TrimPrefix(lowerText, "вчера в"))
		timeDur := parseTime(timeStr) - time.Duration(24*60*60*1e9)
		return today.Add(timeDur)
	case strings.HasPrefix(lowerText, "сегодня в"):
		timeStr := strings.TrimSpace(strings.TrimPrefix(lowerText, "сегодня в"))
		timeDur := parseTime(timeStr)
		return today.Add(timeDur)
	default:
		parts := strings.Split(lowerText, " ")
		if len(parts) != 3 {
			panic(fmt.Sprintf("Exected date to be in format \"DD MMM YYYY\" but got \"%s\"", lowerText))
		}
		day, err := strconv.ParseInt(parts[0], 10, 8)
		if err != nil {
			panic(fmt.Sprintf("Expected day to be integer, but got \"%s\"", parts[0]))
		}
		year, err := strconv.ParseInt(parts[2], 10, 16)
		if err != nil {
			panic(fmt.Sprintf("Expected year to be integer, but got \"%s\"", parts[2]))
		}
		var month time.Month = -1
		months := []string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}
		for i, m := range months {
			if parts[1] == m {
				month = time.Month(i + 1)
				break
			}
		}
		if month == -1 {
			panic(fmt.Sprintf("Expected month to be one of %v, but got \"%s\"", months, parts[1]))
		}
		return time.Date(int(year), month, int(day), 0, 0, 0, 0, time.UTC)
	}
}

func ExtractCommentFromCardPage(doc *html.Node) string {
	node := htmlquery.FindOne(doc, "//div[@itemprop='description']/br")
	textNode := node.NextSibling
	return strings.TrimSpace(textNode.Data)
}

func ExtractAnimalSexSpecFromCardPage(doc *html.Node) types.Sex {
	sexNode := htmlquery.FindOne(doc, "//strong[contains(text(), 'Пол животного')]")

	if sexNode == nil {
		panic("Can't find pet sex specification element on the page")
	}

	var curNode *html.Node = sexNode

	for {
		sib := curNode.NextSibling
		if sib.Type == html.ElementNode && sib.Data == "strong" {
			break
		}
		if sib == nil {
			break
		}
		if sib.Type == html.TextNode {
			trimmed := strings.ToLower(strings.TrimSpace(sib.Data))
			switch trimmed {
			case "---":
				return types.UndefinedSex
			case "самка":
				return types.Female
			case "самец":
				return types.Male
			}
		}
		curNode = sib
	}
	panic("Can't find animal sex specification on pet card page")
}

func ExtractSmallPhotoUrlFromCardPage(doc *html.Node) *url.URL {
	photoNode := htmlquery.FindOne(doc, "//img[contains(@class, 'bd_image_small2')]")
	if photoNode == nil {
		panic("Could not find photo node")
	}

	for _, attr := range photoNode.Attr {
		if attr.Key == "src" {
			urlText := attr.Val
			const mediumPrefix string = "https://poiskzoo.ru/images/board/medium"
			const smallPrefix string = "https://poiskzoo.ru/images/board/small"
			if strings.HasPrefix(urlText, mediumPrefix) {
				suffix := strings.TrimPrefix(urlText, mediumPrefix)
				smallUrlText := fmt.Sprintf("%s%s?v=0053", smallPrefix, suffix)
				result, err := url.Parse(smallUrlText)
				if err != nil {
					panic(fmt.Sprintf("Failed to parse url %s", smallUrlText))
				}
				return result
			} else {
				panic(fmt.Sprintf("Unsupported image prefix in image url: %s", urlText))
			}
		}
	}
	panic("Image node does not contain src attribute")
}
