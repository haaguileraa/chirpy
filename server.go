package main

import (
	"log"
	"net/http"
	"time"
)

const (
	filepathRoot = "./"
	port = ":8080"
)

func serve() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	s := &http.Server {
		Addr:		port,
		Handler:	mux,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
	}
	
	log.Println("Starting to serve on port", port)
	log.Fatal(s.ListenAndServe())
}
