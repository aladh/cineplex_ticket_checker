package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "https://www.cineplex.com/Movie/"

func main() {
	url := baseURL + flag.Arg(0)

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error making request: %s\n", err)
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
	}

	log.Printf("%s\n", html)
}
