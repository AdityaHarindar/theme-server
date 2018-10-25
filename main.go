package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const resources = "./resources/"

var domains []string

func init() {
	domains = getAllDomainResourceNames()
}

func getAllDomainResourceNames() []string {

	//Get all file names as string array
	out, err := exec.Command("ls", resources).Output()
	if err != nil {
		panic("excuse me, what?")
	}

	//File names returned by 'ls' are \n separated
	fNames := strings.Split(string(out), "\n")

	return fNames
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

	// Register routes for all domains in resources
	r.HandleFunc("/health", Health)
	for _, d := range domains {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(resources + d)))
	}

	http.Handle("/", r)

	func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Alive"))
}
