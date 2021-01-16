package main

import (
	"./crawler"
	"./io/postgres"
	"./secrets"
)

func main() {
	postgres.Initialize()
	crawler.Crawl(secrets.GetSearchURL(0, 50))
}
