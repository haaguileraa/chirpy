package main

import (
	"log"
	"net/http"
	"time"
)

const (
	filepathRoot = "."
	pathPrefix = "/app"
	port = ":8080"
)

func serve() {

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix(pathPrefix, http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/healthz", readinessEndpoint())

	s := &http.Server {
		Addr:		port,
		Handler:	mux,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
	}
	
	log.Println("Starting to serve on port", port)
	log.Fatal(s.ListenAndServe())
}

func readinessEndpoint() http.Handler {
	serve := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	return http.HandlerFunc(serve)
}
