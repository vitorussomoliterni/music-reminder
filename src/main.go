package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	bandName := "nofx"

	uri := "https://musicbrainz.org/ws/2/artist?query="
	artistData := getHTTPResponse(uri + bandName)

	artistList := ArtistList{}
	err := xml.Unmarshal(artistData, &artistList)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Result:", artistList)
}

func getHTTPResponse(url string) []byte {
	formattedURL := strings.Replace(url, " ", "%20", -1)

	resp, err := http.Get(formattedURL)
	if err != nil {
		fmt.Println("Error getting a reponse:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error status not OK:", err)
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return data
}

type Artist struct {
	Name  string `xml:"name"`
	Ended bool   `xml:"life-span>ended"`
	ID    string `xml:"id,attr"`
	Score int32  `xml:"score,attr"`
}

type ArtistList struct {
	Artists []Artist `xml:"artist-list>artist"`
	Count   int32    `xml:"count,attr"`
}
