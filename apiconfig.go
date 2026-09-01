package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
}

func newApiConfig() *apiConfig {
	return &apiConfig {}
} 

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) handlerNumberRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	metrics := getFormattedMetrics(int(cfg.fileserverHits.Load()))
	w.Write([]byte(metrics))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset"))
}
