package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const baseURL = "https://www.cineplex.com/Movie/"

func main() {
	theatreIDs := flag.String("t", "", "A comma-separated list of theatre IDs to look for")
	flag.Parse()

	movie := flag.Arg(0)
	url := baseURL + movie

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error making request: %s\n", err)
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
	}

	if isAvailable(theatreIDs, &html) {
		log.Fatalf("Tickets to %s are available\n", movie)
	}
}

func isAvailable(theatreIDs *string, html *[]byte) bool {
	return regexp.MustCompile(strings.ReplaceAll(*theatreIDs, ",", "|")).MatchString(string(*html))
}
