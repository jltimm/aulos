package main

import (
	"./crawler"
	"./secrets"
)

func main() {
	crawler.Crawl(secrets.GetSearchURL(0, 50))
}
