package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/higor-tavares/shortener/src/url"
	"encoding/json"
	"flag"
    "github.com/grafana/pyroscope-go"
	"os"
)

var (
	port *int
	logEnabled *bool
	baseUrl string
)

type Headers map[string]string
type Redirector struct {
	stats chan string
}
func init() {
	port = flag.Int("p",9999, "port")
	logEnabled = flag.Bool("l", true, "log enabled/disabled")
	flag.Parse()
	baseUrl = fmt.Sprintf("http://localhost:%d", *port)
}

func main() {
	
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "shortener",
		ServerAddress:   "http://localhost:4040",
		Logger:          pyroscope.StandardLogger,
		Tags:            map[string]string{"hostname": os.Getenv("HOSTNAME")},
	
		ProfileTypes: []pyroscope.ProfileType{
		  pyroscope.ProfileCPU,
		  pyroscope.ProfileAllocObjects,
		  pyroscope.ProfileAllocSpace,
		  pyroscope.ProfileInuseObjects,
		  pyroscope.ProfileInuseSpace,
		  pyroscope.ProfileGoroutines,
		},
	})

	stats := make(chan string)
	defer close(stats)
	go registerStats(stats)
	url.SetUpRepository(url.NewMemoryRepository())
	http.HandleFunc("/api/short", Shortener)
	http.Handle("/r/", &Redirector{stats})
	http.HandleFunc("/api/stats/", Statistics)
	logMessage("Starting the server on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
func Statistics(w http.ResponseWriter, r *http.Request) {
	searchAndExec(
		w,
		r,
		func (url *url.Url) {
			json, err := json.Marshal(url.Stats())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			sendJson(w, string(json))
	})
}

func (redirector *Redirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	searchAndExec(
		w, 
		r,
		func (url *url.Url) {
			http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
			redirector.stats <- url.ID
	})
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
	logMessage("URL: %s -> SHORT URL :%s.", url.Destination, shortUrl)
	respondWith(w, status, Headers{
		"Location":shortUrl,
		"Link": fmt.Sprintf("<%s/api/stats/%s>;rel=\"stats\"", baseUrl, url.ID),
	})
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
		logMessage("Click on %s registered successfuly!", id)
	}
}

func sendJson(w http.ResponseWriter, resposta string) {
    respondWith(w, http.StatusOK, Headers{
        "Content-Type": "application/json",
    })
    fmt.Fprintf(w, resposta)
}

func searchAndExec(
	w http.ResponseWriter, 
	r *http.Request,
	exec func(*url.Url)) {
		path := strings.Split(r.URL.Path,"/")
		id := path[len(path)-1]
		if url := url.Search(id); url != nil {
			exec(url)
		} else {
			http.NotFound(w, r)
		}
}

func logMessage(format string, values ...interface{}) {
	if *logEnabled { 
		log.Printf(fmt.Sprintf("%s\n",format), values...)
	}
}