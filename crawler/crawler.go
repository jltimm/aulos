package crawler

import (
	"io/ioutil"
	"log"
	"net/http"

	"../secrets"
)

// Crawl grabs the top 10,000 artists on Spotify
func Crawl(offset int, limit int) {
	url := secrets.GetSearchURL(offset, limit)
	var bearer = "Bearer " + secrets.GetToken()
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
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
