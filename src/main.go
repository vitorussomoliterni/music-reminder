package main

import (
	"fmt"

	"github.com/vitorussomoliterni/musichino/services/queryartist"
)

// Variables used for testing
const artist string = "the men"

func main() {
	artistBestMatches := queryartist.GetArtist(artist)
	fmt.Println("Best matches found:")
	for _, a := range artistBestMatches {
		fmt.Println(a.friendlyString())
	}
}
