package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mainmod/core/gateways"
	"net/http"
	"os"
)

type AddResponse struct {
	ShortUrl string `json:"ShortUrl"`
}

type AddRequest struct {
	Url string `json:"url"`
}

type Handler struct {
	redisClient  gateways.RedisGateway
	urlGenerator ShortGenerator
}

func NewHandler(
	redisClient gateways.RedisGateway,
	urlGenerator ShortGenerator,
) *Handler {

	return &Handler{
		redisClient:  redisClient,
		urlGenerator: urlGenerator,
	}
}

func (h *Handler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	var req AddRequest
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := h.urlGenerator.GenShortURL()

	err = h.redisClient.Write(ctx, shortUrl, req.Url)
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

func (h *Handler) HandleRedirectToSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	shortKey, ok := vars["shortKey"]
	if !ok {
		http.Error(w, "cannot parse short key from path", http.StatusBadRequest)
		return
	}

	url, err := h.redisClient.Read(ctx, shortKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/html/not-found.html")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "html/text")
	_, err = w.Write(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/index.html")
}
