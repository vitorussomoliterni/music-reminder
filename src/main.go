package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Variables used for testing
const (
	artist              string = "nofx"
	baseArtistSearchURL string = "https://musicbrainz.org/ws/2/artist?query="
)

func main() {
	artistCleanedName := cleanArtistName(artist)

	searchResult, err := getHTTPResponse(baseArtistSearchURL + artistCleanedName)

	if err != nil {
		fmt.Println(err)
	}

	artists, err := getArtistList(searchResult)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Artists:", artists)
}

func cleanArtistName(artist string) string {
	return strings.Replace(artist, " ", "%20", -1)
}

func getArtistList(xmlResponse []byte) ([]Artist, error) {
	artistList := ArtistList{}
	err := xml.Unmarshal(xmlResponse, &artistList)

	if err != nil {
		return nil, err
	}

	return artistList.Artists, nil
}

func getHTTPResponse(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

type Artist struct {
	Name        string `xml:"name"`
	StillActive bool   `xml:"life-span>ended"`
	ID          string `xml:"id,attr"`
	SearchScore int32  `xml:"score,attr"`
}

type ArtistList struct {
	Artists []Artist `xml:"artist-list>artist"`
	Count   int32    `xml:"count,attr"`
}
