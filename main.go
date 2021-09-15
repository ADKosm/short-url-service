package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "plain/text")
}

type HTTPHandler struct {
	storage Storage
}

var alphabet = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HTTPHandler) handlePostUrl(rw http.ResponseWriter, r *http.Request) {
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.storage.PutURL(RedirectURL(data.Url))
	if err != nil {
		internalError(rw, err)
		return
	}
	//  http://my.site.com/bdfhfd

	response := PutResponseData{
		Key: string(key),
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *HTTPHandler) handleGetUrl(rw http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
	url, err := h.storage.GetURL(Key(key))
	switch {
	case err == nil:
		http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
	case errors.Is(err, ErrNotFound):
		http.NotFound(rw, r)
	default:
		internalError(rw, err)
	}
}

func internalError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	body, _ := json.Marshal(map[string]string{"error": err.Error()})
	_, _ = rw.Write(body)
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	handler := &HTTPHandler{
		storage: NewStorage(),
	}

	r.HandleFunc("/", handleRoot).Methods("GET", "POST")
	r.HandleFunc("/{shortUrl:\\w{5}}", handler.handleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.handlePostUrl)

	return &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
