package checker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const baseURL = "https://www.cineplex.com/movie/"
const availabilityIndicator = `"hasShowtimes":true`

func FindAvailableMovies(movies []string) <-chan string {
	availableMovies := make(chan string)
	go findAvailableMoviesAsync(movies, availableMovies)
	return availableMovies
}

func MovieUrl(movie string) string {
	return baseURL + movie
}

func findAvailableMoviesAsync(movies []string, availableChan chan<- string) {
	wg := sync.WaitGroup{}

	for _, movie := range movies {
		wg.Add(1)

		go func(movie string) {
			defer wg.Done()

			available, err := isAvailable(movie)
			if err != nil {
				log.Fatalf("error checking movie %s: %s\n", movie, err)
				return
			}

			if available {
				availableChan <- movie
			}
		}(movie)
	}

	wg.Wait()
	close(availableChan)
}

func isAvailable(movie string) (bool, error) {
	log.Printf("Checking %s\n", movie)

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// Don't follow redirects because movies that are not found redirect to the cineplex home page
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(MovieUrl(movie))
	if err != nil {
		return false, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("error closing response body: %s\n", err)
		}
	}()

	if res.StatusCode != 200 {
		log.Printf("received bad response %s for movie %s\n", res.Status, movie)
		return false, nil
	}

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %s\n", err)
	}

	html := string(respBytes)

	return strings.Contains(html, availabilityIndicator), nil
}
