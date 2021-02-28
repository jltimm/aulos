package crawler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../auth"
	"../common"
	"../io/postgres"
	"../secrets"
)

var client = auth.GetConfig().Client(context.Background())

func getResponseBodyAndStatus(url string) ([]byte, int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode
}

func crawlArtists(url string) {
	body, statusCode := getResponseBodyAndStatus(url)
	if statusCode == 200 {
		var response common.Response
		json.Unmarshal([]byte(body), &response)
		postgres.InsertArtists(response.Artists.Items)
	} else if statusCode == 404 {
		var response common.NotFound
		json.Unmarshal(body, &response)
		log.Println(response.Error.Message)
	} else {
		// Need to retry
	}
}

func crawlRecommendedArtists() {
	artists := postgres.GetArtistsWithNullRecommended()
	numArtists := len(artists)
	if numArtists == 0 {
		return
	}
	for i := 0; i < numArtists; i++ {
		id := artists[i]
		body, statusCode := getResponseBodyAndStatus(secrets.GetRecommendedURL(id))
		if statusCode == 200 {
			var response common.RecommendedResponse
			json.Unmarshal([]byte(body), &response)
			postgres.InsertArtists(response.Artists)
			postgres.UpdateRecommended(id, response.Artists)
		} else {
			log.Println(string(body))
		}
		time.Sleep(500 * time.Millisecond)
	}
	// crawl through the artists that were just added
	crawlRecommendedArtists()
}

func createMap(artists []common.Artist) map[string]common.Array {
	artistsMap := make(map[string]common.Array)
	for i := 0; i < len(artists); i++ {
		artist := artists[i]
		artistsMap[artist.ID] = artist.Recommended
	}
	return artistsMap
}

func updateRecommendedArtists() {
	artists := postgres.GetAllArtists()
	artistsMap := createMap(artists)
	for i := 0; i < len(artists); i++ {
		artist := artists[i]
		for j := 0; j < len(artist.Recommended); j++ {
			artistToUpdateRecommended := artistsMap[artist.Recommended[j]]
			if !artistToUpdateRecommended.Contains(artist.ID) {
				artistToUpdateRecommended = append(artistToUpdateRecommended, artist.ID)
			}
			artistsMap[artist.Recommended[j]] = artistToUpdateRecommended
		}
	}
	postgres.UpdateMissingRecommended(artistsMap)
}

// Crawl grabs the top 5,000 artists on Spotify
func Crawl(url string) {
	crawlArtists(url)
	crawlRecommendedArtists()
	updateRecommendedArtists()
}
