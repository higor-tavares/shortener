package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/higor-tavares/shortener/src/url"
)

var (
	port int
	baseUrl string
	stats chan string
)

type Headers map[string]string

func init() {
	port = 9999
	baseUrl = fmt.Sprintf("http://localhost:%d",port)
}

func main() {
	stats = make(chan string)
	defer close(stats)
	go registerStats(stats)
	url.SetUpRepository(url.NewMemoryRepository())
	http.HandleFunc("/api/short", Shortener)
	http.HandleFunc("/r/", Redirector)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
func Redirector(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path,"/")
	id := path[len(path)-1]
	if url := url.Search(id); url != nil {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
		stats <- id
	} else {
		http.NotFound(w, r)
	}
}
func Shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWith(w, http.StatusMethodNotAllowed, Headers{
			"Allow":"POST",
		})
		return
	}
	url, isNew, err := url.CreateIfNotExists(extractUrl(r))
	if err != nil {
		respondWith(w, http.StatusBadRequest, nil)
		return
	}
	var status int
	if isNew {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}
	shortUrl := fmt.Sprintf("%s/r/%s", baseUrl, url.ID)
	respondWith(w, status, Headers{"Location":shortUrl})
}

func respondWith(
	w http.ResponseWriter,
	status int,
	headers Headers) {
	for k,v := range headers {
		w.Header().Set(k,v)
	}
	w.WriteHeader(status)
}

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}
func registerStats(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		fmt.Printf("Click on %s registered successfuly\n", id)
	}
}