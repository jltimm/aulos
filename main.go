package main

import (
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
	postgres.Initialize()
	crawler.Crawl(secrets.GetSearchURL(0, 1))
}
