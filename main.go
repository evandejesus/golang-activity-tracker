package main

import "github.com/evandejesus/activity-tracker/internal/server"

func main() {
	s := server.NewHTTPServer("localhost:8080")
	s.ListenAndServe()
}
