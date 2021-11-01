package checker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const baseURL = "https://www.cineplex.com/Movie/"
const availableMessage = "Check out showtimes for this movie"

func FindAvailableMovies(movies []string, theatreIDs string) <-chan string {
	availableMovies := make(chan string)
	go findAvailableMoviesAsync(movies, theatreIDs, availableMovies)
	return availableMovies
}

func MovieUrl(movie string) string {
	return baseURL + movie
}

func findAvailableMoviesAsync(movies []string, theatreIDs string, availableChan chan<- string) {
	wg := sync.WaitGroup{}
	theatreIDsRegex := regexp.MustCompile(strings.ReplaceAll(theatreIDs, ",", "|"))

	for _, movie := range movies {
		wg.Add(1)

		go func(movie string) {
			defer wg.Done()

			available, err := isAvailable(movie, theatreIDsRegex)
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

func isAvailable(movie string, theatreIDsRegex *regexp.Regexp) (bool, error) {
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
		return false, fmt.Errorf("received bad response %s for movie %s", res.Status, movie)
	}

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %s\n", err)
	}

	html := string(respBytes)

	// Reduce false positives by checking for this
	if !strings.Contains(html, availableMessage) {
		return false, nil
	}

	return theatreIDsRegex.MatchString(html), nil
}
