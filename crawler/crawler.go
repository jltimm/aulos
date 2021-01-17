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

// Crawl grabs the top 10,000 artists on Spotify
func Crawl(url string) {
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var response common.Response
		json.Unmarshal([]byte(body), &response)
		go postgres.Insert(response.Artists.Items)
		if response.Artists.Next != "" {
			// Sleep for five seconds as not to hammer the API
			time.Sleep(5000 * time.Millisecond)
			Crawl(response.Artists.Next)
		}
	} else {
		// Need to retry
	}
}
