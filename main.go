package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlePost).Methods("POST")
	r.HandleFunc("/", handleGet).Methods("GET")
	srv := &http.Server{Addr: ":8080", Handler: r}
	fmt.Println("fuck you")
	srv.ListenAndServe()
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get\n")
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post\n")
}
