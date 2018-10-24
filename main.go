package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const resources = "./resources"

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

	// Register routes
	r.HandleFunc("/theme", ThemeHandler).Methods("GET").Headers("domain")
	r.HandleFunc("/health", Health).Methods("GET")

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

func ThemeHandler(w http.ResponseWriter, r *http.Request) {

	domain := r.Header.Get("Domain")

	for _, d := range getAllDomainResourceNames() {
		if strings.Compare(d, domain) == 0 {
			http.FileServer(http.Dir(resources + domain)) //wtf is this doing here
		}
	}

}
