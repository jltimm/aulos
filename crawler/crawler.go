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
)

var client = auth.GetConfig().Client(context.Background())

func crawlArtists(url string) {
	count := postgres.GetArtistCount()
	if count >= 2000 {
		log.Printf("%d rows in table: skipping artist crawl\n", count)
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var response common.Response
		json.Unmarshal([]byte(body), &response)
		postgres.InsertArtist(response.Artists.Items)
		if response.Artists.Next != "" {
			// Sleep for five seconds as not to hammer the API
			time.Sleep(1000 * time.Millisecond)
			Crawl(response.Artists.Next)
		}
	} else if resp.StatusCode == 404 {
		var response common.NotFound
		json.Unmarshal(body, &response)
		log.Println(response.Error.Message)
	} else {
		// Need to retry
	}
}

// Crawl grabs the top 10,000 artists on Spotify
func Crawl(url string) {
	crawlArtists(url)
}
