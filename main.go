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

var theatreIDsRegex *regexp.Regexp
var movies []string

func init() {
	theatreIDs := flag.String("t", "", "A comma-separated list of theatre IDs to look for")
	flag.Parse()

	theatreIDsRegex = regexp.MustCompile(strings.ReplaceAll(*theatreIDs, ",", "|"))

	movies = strings.Split(flag.Arg(0), ",")
}

func main() {
	availableChan := make(chan string, len(movies))
	wg := sync.WaitGroup{}

	for _, movie := range movies {
		wg.Add(1)

		go func(movie string, availableChan chan<- string) {
			defer wg.Done()

			available, err := isAvailable(movie)
			if err != nil {
				log.Printf("error checking availability: %s\n", err)
				return
			}

			if available {
				log.Printf("Tickets to %s are available\n", movie)
				availableChan <- movie
			}
		}(movie, availableChan)
	}

	wg.Wait()
	close(availableChan)

	if len(availableChan) > 0 {
		log.Fatalln("Go buy tickets!")
	}
}

func isAvailable(movie string) (bool, error) {
	log.Printf("Checking %s\n", movie)

	url := baseURL + movie

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing response body: %w", err)
		}
	}()

	if res.StatusCode != 200 {
		return false, fmt.Errorf("failed to find movie %s", movie)
	}

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response body: %s\n", err)
	}

	return theatreIDsRegex.MatchString(string(html)), nil
}
