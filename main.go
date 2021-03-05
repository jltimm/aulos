package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./crawler"
	"./io/filecreator"
	"./io/postgres"
	"./path"
	"./secrets"
)

func cleanup() {
	fmt.Println("Shutting down...")
	postgres.Close()
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
	limitPtr := flag.Int("limit", 5000, "The maximum number of artists to crawl")
	popularityCutoffPtr := flag.Int("popularityCutoff", 65, "Popularity cutoff for artists. Set to -1 to ignore.")
	flag.Parse()
	postgres.Initialize()
	postgres.SetPopularityCutoff(*popularityCutoffPtr)
	crawler.Crawl(secrets.GetSearchURL(0, 1), *limitPtr)
	keyMap, distance, path := path.FloydWarshall()
	filecreator.CreateFile(keyMap, distance, path)
}
