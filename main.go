package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const baseURL = "https://www.cineplex.com/Movie/"

func main() {
	theatreID := flag.String("t", "", "ID of the theatre to check")
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

	available := strings.Contains(string(html), *theatreID)

	if available {
		log.Fatalf("%s is available at %s\n", movie, *theatreID)
	}
}
