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

func serve(cfg *apiConfig) {
	
	mux := http.NewServeMux()
	handlerFileServer := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(handlerFileServer))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", cfg.handlerUser)
	mux.HandleFunc("POST /api/chirps", cfg.handlerPostChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", getBearerTokenMiddleware(cfg.handlerRefresh))
	mux.HandleFunc("POST /api/revoke", getBearerTokenMiddleware(cfg.handlerRevokeRefreshToken))
	mux.HandleFunc("GET /admin/metrics", cfg.handlerNumberRequests)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	s := &http.Server {
		Addr:		port,
		Handler:	mux,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
	}
	
	log.Println("Starting to serve on port", port)
	log.Fatal(s.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
}
