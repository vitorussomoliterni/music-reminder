package main

import (
	"fmt"

	"github.com/vitorussomoliterni/musichino/services"
)

// Variables used for testing
const artist string = "the men"

func main() {
	artistBestMatches := services.GetArtist(artist)
	fmt.Println("Best matches found:")
	for _, a := range artistBestMatches {
		fmt.Println(a.friendlyString())
	}
}
