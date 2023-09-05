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
)

func init() {
	port = 8888
	baseUrl = fmt.Sprintf("http://localhost:%d",port)
}

func main() {
	http.HanleFunc("/api/short", Sortener)
	http.HanleFunc("/r", Redirector)
	log.Fatal(
		http.ListenAndServe(fmt.Sprintf("%d",port),nil)
	)
}