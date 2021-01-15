package crawler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../auth"
)

// Item is part of Spotify's JSON response and contains
// the Spotify ID, name of the artist, and the popularity
type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

// Artists is a part of Spotify's JSON response and contains
// the list of artists, the next page, current offset, and previous page
type Artists struct {
	Items    []Item `json:"items"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
}

// Response is part of Spotify's JSON response and contains
// the data that we need
type Response struct {
	Artists Artists `json:"artists"`
}

var client = auth.GetConfig().Client(context.Background())

// Crawl grabs the top 10,000 artists on Spotify
func Crawl(url string) {
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var response Response
		json.Unmarshal([]byte(body), &response)
		if response.Artists.Next != "" {
			// Sleep for five seconds as not to hammer the API
			time.Sleep(5000 * time.Millisecond)
			Crawl(response.Artists.Next)
		}
	} else {
		// Need to retry
	}
}
