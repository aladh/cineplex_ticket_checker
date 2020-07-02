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

	movies := strings.Split(flag.Arg(0), ",")

	availableChan := spawnWorkers(&movies, theatreIDs)

	if len(availableChan) > 0 {
		for movie := range availableChan {
			log.Printf("Tickets to %s are available\n", movie)
		}

		log.Fatalln("Go buy tickets!")
	}
}

func spawnWorkers(movies *[]string, theatreIDs *string) chan string {
	availableChan := make(chan string, len(*movies))
	wg := sync.WaitGroup{}

	for _, movie := range *movies {
		wg.Add(1)

		go func(movie string, theatreIDs *string, availableChan chan<- string) {
			defer wg.Done()

			if isAvailable(&movie, theatreIDs) {
				availableChan <- movie
			}
		}(movie, theatreIDs, availableChan)
	}

	wg.Wait()
	close(availableChan)

	return availableChan
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
