package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const themeResources = "./resources/themes/"

type domain string
type theme string

// Mapping
var M = map[domain]theme{
	"localhost": "1",
}

func main() {

	// Define a router
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Register routes for all domains based on domain -> theme mapping
	r.HandleFunc("/health", Health)
	for d, t := range M {
		r.PathPrefix("/").
			Handler(http.FileServer(http.Dir(themeResources + t))).
			Host(string(d)).
			Methods("GET")
	}

	// Register a default route that returns the default theme
	r.PathPrefix("/").
		Handler(http.FileServer(http.Dir(themeResources + "default"))).
		Methods("GET")

	http.Handle("/", r)

	func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

}

func Health(w http.ResponseWriter, r *http.Request) {

	h, _ := os.Hostname()
	resp := struct {
		Status   string `json:"status"`
		Hostname string `json:"hostname"`
	}{"Alive", h}

	resBody, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBody)
}
