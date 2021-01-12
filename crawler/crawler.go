package crawler

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"../auth"
	"../secrets"
)

// Crawl grabs the top 10,000 artists on Spotify
func Crawl(offset int, limit int) {
	url := secrets.GetSearchURL(offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	client := auth.GetConfig().Client(context.Background())
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		log.Println(string([]byte(body)))
	} else {
		// Need to retry
	}
}
