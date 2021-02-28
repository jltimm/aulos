package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./crawler"
	"./io/postgres"
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
	flag.Parse()
	postgres.Initialize()
	crawler.Crawl(secrets.GetSearchURL(0, 1), *limitPtr)
}
