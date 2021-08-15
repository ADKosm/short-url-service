package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var ctx = context.Background()

type AddResponse struct {
	ShortUrl string `json:"ShortUrl"`
}

type AddRequest struct {
	Url string `json:"url"`
}

type Handler struct {
	redisClient *redis.Client
}

func genShortUrl() string {
	alphabet := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	rand.Shuffle(len(alphabet), func(i, j int) {
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
	})
	id := string(alphabet[:8])
	return id
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var req AddRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := genShortUrl()
	err = h.redisClient.Set(ctx, shortUrl, req.Url, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := AddResponse{
		ShortUrl: shortUrl,
	}
	rawRes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(rawRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) handleRedirectToSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortKey, ok := vars["shortKey"]
	if !ok {
		http.Error(w, "cannot parse short key from path", http.StatusBadRequest)
		return
	}

	url, err := h.redisClient.Get(ctx, shortKey).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/not-found.html")
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/index.html")
}


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	rand.Seed(time.Now().Unix())

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	handler := &Handler{redisClient: rdb}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/api/add", handler.handleAdd)
	r.HandleFunc("/not-found", handler.handleNotFound)
	r.HandleFunc("/{shortKey:\\w{8}}", handler.handleRedirectToSite)
	r.HandleFunc("/", handler.handleRoot)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
