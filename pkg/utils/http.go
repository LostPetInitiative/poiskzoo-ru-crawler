package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version"
	"golang.org/x/net/html/charset"
)

func SetUserAgentHeader(header http.Header) {
	header.Set("User-Agent", fmt.Sprintf("LostPetInitiative:poiskzoo-crawler / %s:%.6s (https://kashtanka.pet/)", version.AppVersion, version.GitCommit))
}

var httpClient http.Client = http.Client{}

type HttpFetchResult struct {
	Body        []byte
	ContentType string
}

// Performs the HTTP GET request over the specified targetURL and returns the response body as a string
func HttpGet(targetUrl *url.URL, acceptHeader string) (*HttpFetchResult, error) {
	req, err := http.NewRequest("GET", targetUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", acceptHeader)
	SetUserAgentHeader(req.Header)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var contentType string = resp.Header.Get(http.CanonicalHeaderKey("content-type"))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpFetchResult{body, contentType}, nil
}

// Performs the HTTP POST request to the specified targetUrl. Returns successful HTTP code, if returned err is nil
func HttpPost(targetUrl *url.URL, contentTypeHeader string, body []byte) (*int, error) {
	req, err := http.NewRequest("POST", targetUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentTypeHeader)
	SetUserAgentHeader(req.Header)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		// successful
		return &resp.StatusCode, nil
	} else {
		return nil, fmt.Errorf("not successful HTTP status: %s", resp.Status)
	}
}

// Performs the HTTP GET request over the specified targetURL, recodes the response to UTF-8
func HttpGetHtml(targetUrl *url.URL) (*HttpFetchResult, error) {
	resp, err := HttpGet(targetUrl, types.HtmlMimeType)
	if err != nil {
		return nil, err
	}

	var declaredEncoding string
	if strings.HasPrefix(resp.ContentType, "text/html; charset=") {
		declaredEncoding = strings.TrimPrefix(resp.ContentType, "text/html; charset=")
	}
	// log.Printf("HTML fetch: declared HTML page encoding: %s\n", declaredEncoding)

	// determinedEnc, encName, certain := charset.DetermineEncoding(resp.Body, declaredEncoding)
	_, encName, _ := charset.DetermineEncoding(resp.Body, declaredEncoding)
	// log.Printf("HTML fetch: detected enc %v (%v); certain: %v\n", determinedEnc, encName, certain)

	bodyReader := bytes.NewReader(resp.Body)
	if declaredEncoding == "" {
		declaredEncoding = encName
	}
	// log.Printf("HTML fetch: treating encoding as %v. ReEncoding body from it into UTF-8\n", declaredEncoding)
	reEncodedBodyReader, err := charset.NewReaderLabel(declaredEncoding, bodyReader)
	if err != nil {
		log.Panicln(err)
	}
	reEncodedBody, err := io.ReadAll(reEncodedBodyReader)
	if err != nil {
		log.Panicln(err)
	}

	return &HttpFetchResult{
		ContentType: resp.ContentType,
		Body:        reEncodedBody,
	}, nil

}
