package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version"
)

func SetUserAgentHeader(header http.Header) {
	header.Set("User-Agent", fmt.Sprintf("LostPetInitiative:poiskzoo-crawler / %s:%.6s (https://kashtanka.pet/)", version.AppVersion, version.GitCommit))
}

var httpClient http.Client = http.Client{}

// Performs the HTTP GET request over the specified targetURL and returns the response body as a string
func HttpGet(targetUrl *url.URL, acceptHeader string) ([]byte, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
