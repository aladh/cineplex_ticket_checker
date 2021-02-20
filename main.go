package main

import (
	"flag"
	"fmt"
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

	availableChan := checkMovies(movies, theatreIDs)

	if len(availableChan) > 0 {
		log.Fatalln("Go buy tickets!")
	}
}

func checkMovies(movies []string, theatreIDs *string) chan string {
	availableChan := make(chan string, len(movies))
	wg := sync.WaitGroup{}

	for _, movie := range movies {
		wg.Add(1)

		go func(movie string, theatreIDs *string, availableChan chan<- string) {
			defer wg.Done()

			available, err := isAvailable(&movie, theatreIDs)
			if err != nil {
				log.Printf("error checking availability: %s\n", err)
				return
			}

			if available {
				log.Printf("Tickets to %s are available\n", movie)
				availableChan <- movie
			}
		}(movie, theatreIDs, availableChan)
	}

	wg.Wait()
	close(availableChan)

	return availableChan
}

func isAvailable(movie *string, theatreIDs *string) (bool, error) {
	log.Printf("Checking %s\n", *movie)

	url := baseURL + *movie

	res, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing response body: %w", err)
		}
	}()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
	}

	if strings.Contains(string(html), "January 1, 0001") {
		return false, fmt.Errorf("failed to find movie %s", *movie)
	}

	return regexp.MustCompile(strings.ReplaceAll(*theatreIDs, ",", "|")).MatchString(string(html)), nil
}
