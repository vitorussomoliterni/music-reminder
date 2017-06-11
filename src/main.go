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
	userAgent           string = "Music_Reminder_Bot/0.1 ( https://github.com/vitorussomoliterni/music-reminder/ )"
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

	artist := getBestArtistMatch(artists)

	fmt.Println("Best match found:", artist.Name)
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
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Fatal error: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func getBestArtistMatch(artists []Artist) Artist {
	for _, a := range artists {
		if a.SearchScore == 100 {
			return a
		}
	}

	return Artist{}
}

type Artist struct {
	Name        string `xml:"name"`
	StillActive bool   `xml:"life-span>ended"`
	ID          string `xml:"id,attr"`
	SearchScore int32  `xml:"score,attr"`
}

type ArtistList struct {
	Artists []Artist `xml:"artist-list>artist"`
}
