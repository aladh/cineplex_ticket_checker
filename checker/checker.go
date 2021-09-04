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

var theatreIDsRegex *regexp.Regexp

func AvailableMovies(movies []string, theatreIDs string) ([]string, error) {
	availableChan := make(chan string, len(movies))
	errorChan := make(chan error, len(movies))
	wg := sync.WaitGroup{}
	theatreIDsRegex = regexp.MustCompile(strings.ReplaceAll(theatreIDs, ",", "|"))

	for _, movie := range movies {
		wg.Add(1)

		go func(movie string, availableChan chan<- string) {
			defer wg.Done()

			available, err := isAvailable(movie)
			if err != nil {
				errorChan <- err
				return
			}

			if available {
				availableChan <- movie
			}
		}(movie, availableChan)
	}

	wg.Wait()
	close(availableChan)
	close(errorChan)

	if len(errorChan) > 0 {
		return nil, fmt.Errorf("%d errors ocurred: %s\n", len(errorChan), <-errorChan)
	}

	if len(availableChan) > 0 {
		return toSlice(availableChan), nil
	}

	return nil, nil
}

func MovieUrl(movie string) string {
	return baseURL + movie
}

func isAvailable(movie string) (bool, error) {
	log.Printf("Checking %s\n", movie)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects because movies that are not found redirect to the cineplex home page
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(MovieUrl(movie))
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

func toSlice(c chan string) []string {
	s := make([]string, 0)

	for i := range c {
		s = append(s, i)
	}

	return s
}
