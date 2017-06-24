package queryartist

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Variables used for testing
const (
	artist              string = "the men"
	baseArtistSearchURL string = "https://musicbrainz.org/ws/2/artist?query="
	userAgent           string = "Musichino_Bot/0.1 ( https://github.com/vitorussomoliterni/musichino/ )"
)

// GetArtist fetches a list of artist names from musicbrainz.
func GetArtist(artist string) []Artist {
	artistCleanedName := cleanArtistName(artist)

	searchResult, err := getHTTPResponse(baseArtistSearchURL + artistCleanedName)

	if err != nil {
		fmt.Println(err)
	}

	artists, err := getArtistList(searchResult)

	if err != nil {
		fmt.Println(err)
	}

	bestMatches := getBestArtistMatches(artists)
	fmt.Println("Best matches found:")
	for _, a := range bestMatches {
		fmt.Println(a.friendlyString())
	}

	return artists
}

func cleanArtistName(artist string) string {
	return strings.Replace(artist, " ", "%20", -1)
}

func getArtistList(xmlResponse []byte) ([]Artist, error) {
	artistList := artistList{}
	err := xml.Unmarshal(xmlResponse, &artistList)
	if err != nil {
		return nil, err
	}

	if len(artistList.Artists) == 0 {
		return nil, fmt.Errorf("no artist found")
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

func getBestArtistMatches(artists []Artist) []Artist {
	results := make([]Artist, 2)

	for i := range artists {
		if artists[i].SearchScore == 100 {
			results = append(results, artists[i])
		}
	}

	if len(results) > 1 { // If more than one result has a 100% score, return all results
		return artists
	}

	return results
}

func (artist *Artist) friendlyString() string {
	result := artist.Name

	if len(artist.Area) > 0 {
		result += " | " + artist.Area
	}
	if len(artist.ActivityBegin) > 0 {
		result += " (" + artist.ActivityBegin
	}
	if artist.ActvityEnded && len(artist.ActivityEnd) > 0 {
		result += " - " + artist.ActivityEnd
	}
	if len(artist.ActivityBegin) > 0 {
		result += ")"
	}
	if len(artist.Disambiguation) > 0 {
		result += " [" + artist.Disambiguation + "]"
	}

	return result
}

// Artist struct that maps the xml results to a go struct.
type Artist struct {
	Name           string `xml:"name"`
	Area           string `xml:"area>name"`
	ActivityBegin  string `xml:"life-span>begin"`
	ActivityEnd    string `xml:"life-span>end"`
	ActvityEnded   bool   `xml:"life-span>ended"`
	ID             string `xml:"id,attr"`
	SearchScore    int32  `xml:"score,attr"`
	Disambiguation string `xml:"disambiguation"`
}

type artistList struct {
	Artists []Artist `xml:"artist-list>artist"`
}
