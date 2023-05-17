package main

import (
	"log"

	"github.com/evandejesus/activity-tracker/internal/server"
)

func main() {
	url := "localhost:8080"
	s := server.NewHTTPServer(url)
	log.Printf("serving @ %s\n", url)
	s.ListenAndServe()
}
