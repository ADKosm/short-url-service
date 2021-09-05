package main

import (
	"github.com/gorilla/mux"
	"log"
	"mainmod/core/gateways"
	"mainmod/core/handlers"
	"math/rand"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	rand.Seed(time.Now().Unix())

	redisClient := gateways.NewGateway()
	generator := handlers.NewGenerator()

	handler := handlers.NewHandler(redisClient, generator)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/api/add", handler.HandleAdd)
	r.HandleFunc("/not-found", handler.HandleNotFound)
	r.HandleFunc("/{shortKey:\\w{8}}", handler.HandleRedirectToSite)
	r.HandleFunc("/", handler.HandleRoot)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
