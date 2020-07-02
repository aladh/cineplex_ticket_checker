package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const baseURL = "https://www.cineplex.com/Movie/"

func main() {
	theatreIDs := flag.String("t", "", "A comma-separated list of theatre IDs to look for")
	flag.Parse()

	movie := flag.Arg(0)
	availableMovies := make(chan string, 1)
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func(availableMovies chan<- string) {
		defer wg.Done()

		if isAvailable(&movie, theatreIDs) {
			availableMovies <- movie
		}
	}(availableMovies)

	wg.Wait()
	close(availableMovies)

	if len(availableMovies) > 0 {
		for movie := range availableMovies {
			log.Printf("Tickets to %s are available\n", movie)
		}

		log.Fatalln("Go buy tickets!")
	}
}

func isAvailable(movie *string, theatreIDs *string) bool {
	url := baseURL + *movie

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error making request: %s\n", err)
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
	}

	return regexp.MustCompile(strings.ReplaceAll(*theatreIDs, ",", "|")).MatchString(string(html))
}
