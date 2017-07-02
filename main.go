package main

import (
	"fmt"
	"musichino/src/services"
)

// Variables used for testing
const artist string = "the men"

func main() {
	artistBestMatches := services.GetArtist(artist)
	fmt.Println("Best matches found:")
	for _, a := range artistBestMatches {
		fmt.Println(a.FriendlyString())
	}
}
