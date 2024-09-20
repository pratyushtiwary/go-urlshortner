package main

import (
	"log"
	"net/http"
	"time"
	"urlshortner/shortner"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/shorten", shortner.ShortenHandler).Methods("POST")
	router.HandleFunc("/go/{name}", shortner.RedirectHandler).Methods("GET")

	server := http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on port 8080...")
	log.Fatal(server.ListenAndServe())

}
