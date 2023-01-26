package crawler

import (
	"log"
	"net/url"
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

	extractedUrls := ExtractCardsFromCatalogDocument(ParseHtmlContent(catalogHtml))
	const expectedCount int = 52
	if len(extractedUrls) != expectedCount {
		t.Logf("Expected to extract %d card IDs but extracted %d", expectedCount, len(extractedUrls))
		t.Fail()
	}

	cardMap := make(map[int]Card)
	for _, card := range extractedUrls {
		cardMap[int(card.Id)] = card
	}

	card, exists := cardMap[164833]
	if !exists {
		t.Logf("Did not find expected URL in the result set")
		t.Fail()
	}
	if card.HasPaidPromotion {
		t.Log("Card 164833 is promoted")
		t.Fail()
	}

	card, exists = cardMap[164656]
	if !exists {
		t.Log("Did not find expected URL in the result set")
		t.Fail()
	}
	if !card.HasPaidPromotion {
		t.Log("Card 164656 is not promoted")
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
		{"./testdata/168308.html.dump", types.Bird},
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
		expected types.EventType
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

		var extractedType types.EventType = ExtractCardTypeFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedType != testCase.expected {
			t.Logf("Wrong card type extracted for %s. Expected %v, but got %v", testCase.path, testCase.expected, extractedType)
			t.Fail()
		}
	}
}

func TestExtractAnimalSexSpecFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path     string
		expected types.Sex
	}{
		{"./testdata/164921.html.dump", types.Male},
		{"./testdata/164923.html.dump", types.Female},
		{"./testdata/164929.html.dump", types.Female},
		{"./testdata/164931.html.dump", types.Female},
		{"./testdata/164978.html.dump", types.UndefinedSex},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedSexSpec types.Sex = ExtractAnimalSexSpecFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedSexSpec != testCase.expected {
			t.Logf("Wrong card type extracted for %s. Expected %v, but got %v", testCase.path, testCase.expected, extractedSexSpec)
			t.Fail()
		}
	}
}

func TestSmallPhotoUrlFromCardPage(t *testing.T) {
	testCases := []struct {
		path, expected string
	}{
		{"./testdata/164921.html.dump", "https://poiskzoo.ru/images/board/small/propala-koshka-164921-propala-koshka-kot-g-orenburg.jpg?v=0053"},
		{"./testdata/164923.html.dump", "https://poiskzoo.ru/images/board/small/naydena-koshka-164923-naydena-koshka-britanka-g-orehovo-zuevo.jpg?v=0053"},
		{"./testdata/164929.html.dump", "https://poiskzoo.ru/images/board/small/naydena-sobaka-169c9d1281e68d3ccdd52712bc493686-naydena-sobaka-umnaya-ryzhaya-vzroslaya-devochka-g-vladivostok.jpg?v=0053"},
		{"./testdata/164931.html.dump", "https://poiskzoo.ru/images/board/small/propala-sobaka-164931-propala-sobaka-toy-pudel-g-surgut.jpg?v=0053"},
		{"./testdata/164978.html.dump", "https://poiskzoo.ru/images/board/small/naydena-koshka-0076a8bb193a494482fd60268126d9ee-naydena-koshka-ili-kot-g-stavropol.jpg?v=0053"},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedPhotoUrl *url.URL = ExtractSmallPhotoUrlFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedPhotoUrl.String() != testCase.expected {
			t.Logf("Wrong card type extracted for %s. Expected %v, but got %v", testCase.path, testCase.expected, extractedPhotoUrl)
			t.Fail()
		}
	}
}

func TestExtractAddressFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path, city, address string
	}{
		{"./testdata/164921.html.dump", "Оренбург", "Центральный"},
		{"./testdata/164923.html.dump", "Орехово-Зуево", "Демихово"},
		{"./testdata/164929.html.dump", "Владивосток", "Владивосток, район Арт-пляжа."},
		{"./testdata/164931.html.dump", "Сургут", "г. Сургут, пр. Пролетарский 8/1-8/2"},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extracted *CityAndAddress = ExtractAddressFromCardPage(ParseHtmlContent(catalogHtml))
		if extracted.City != testCase.city {
			t.Logf("Wrong city extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.city, extracted.City)
			t.Fail()
		}

		if extracted.Address != testCase.address {
			t.Logf("Wrong address extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.city, extracted.City)
			t.Fail()
		}
	}
}

func TestExtractEventTimeFromPetCardPage(t *testing.T) {
	today := time.Date(2022, 10, 17, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		path string
		time time.Time
	}{
		{"./testdata/164793.html.dump", time.Date(2022, 10, 13, 0, 0, 0, 0, time.UTC)},
		{"./testdata/164921.html.dump", time.Date(2022, 10, 16, 22, 27, 0, 0, time.UTC)},
		{"./testdata/164923.html.dump", time.Date(2022, 10, 17, 0, 10, 0, 0, time.UTC)},
		{"./testdata/164929.html.dump", time.Date(2022, 10, 17, 6, 52, 0, 0, time.UTC)},
		{"./testdata/164931.html.dump", time.Date(2022, 10, 17, 7, 45, 0, 0, time.UTC)},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedEventTime time.Time = ExtractEventDateFromCardPage(ParseHtmlContent(catalogHtml), today)
		if extractedEventTime != testCase.time {
			t.Logf("Wrong event time extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.time, extractedEventTime)
			t.Fail()
		}
	}
}

func TestExtractCommentFromPetCardPage(t *testing.T) {
	testCases := []struct {
		path, comment string
	}{
		{"./testdata/164921.html.dump", "Бенгальский кот пропал в приделах улиц советская цвиллинга рыбаковская, кот длинный, окрас леопардовый"},
		{"./testdata/164923.html.dump", "Найдена британская Кошечка. Кто потерял?"},
		{"./testdata/164929.html.dump", "15 октября прибилась поздно вечером эта девочка. Воспитанная, умная, видно что домашняя, не боится детей и машин. Хозяева, отзовитесь"},
		{"./testdata/164931.html.dump", "Очень срочно  Сегодня утром, 17. 10. 22 г. в 6. 10-6. 20, в р-не пр. Пролетарского 8/1-8/2 потерялась маленькая собачка - той-пудель рыжего окраса. В холке очень маленькая - 22 см. Собака взрослая, хоть и выглядит, как щенок. Напугал волкодав, погнался за моей собакой. Возможно, укусил, так как моя собака заскулила очень сильно. Возможно, просто испугалась - я не увидела. У нее уже был сердечный приступ, может от перенесенного страха нуждаться в помощи ветеринара. Прошу оказать помощь в поиске. Собачка контактная, может пойти на зов. Зовут Нэсси. На заднем бедре есть клеймо - SLN 853. Если даже просто увидите - дайте знать.. Людмила"},
	}

	for _, testCase := range testCases {
		fileContent, err := os.ReadFile(testCase.path)
		if err != nil {
			log.Fatal(err)
		}
		catalogHtml := string(fileContent)

		var extractedAddress string = ExtractCommentFromCardPage(ParseHtmlContent(catalogHtml))
		if extractedAddress != testCase.comment {
			t.Logf("Wrong comment extracted for %s. Expected \"%v\", but got \"%v\"", testCase.path, testCase.comment, extractedAddress)
			t.Fail()
		}
	}
}
