package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type httpServer struct {
	Activities *Activities
}

type IDDocument struct {
	ID uint64 `json:"id"`
}

type ActivityDocument struct {
	Activity Activity `json:"activity"`
}

func NewHTTPServer(addr string) *http.Server {
	server := &httpServer{
		Activities: &Activities{},
	}
	r := mux.NewRouter()
	r.HandleFunc("/", server.handlePost).Methods("POST")
	r.HandleFunc("/", server.handleGet).Methods("GET")

	var handler http.Handler = r
	handler = logRequestHandler(handler)

	return &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second, // introduced in Go 1.8
		Addr:         addr,
		Handler:      r,
	}
}

func (s *httpServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var decoder = schema.NewDecoder()
	var req IDDocument
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		log.Println("Err in GET parameters: ", err)
	}
	activity, err := s.Activities.Retrieve(req.ID)

	if err == ErrIDNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	res := ActivityDocument{Activity: activity}
	json.NewEncoder(w).Encode(res)
}

func (s *httpServer) handlePost(w http.ResponseWriter, r *http.Request) {
	var req ActivityDocument
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// valid
	id := s.Activities.Insert(req.Activity)
	res := IDDocument{ID: id}
	json.NewEncoder(w).Encode(res)
}
