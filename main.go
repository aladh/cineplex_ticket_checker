package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aladh/cineplex_ticket_checker/checker"
	"github.com/aladh/cineplex_ticket_checker/webhook"
)

var webhookURL string

func init() {
	flag.StringVar(&webhookURL, "w", "", "A URL to send webhooks to when movies are available")
	flag.Parse()
}

func main() {
	movies := strings.Split(flag.Arg(0), ",")
	availableMovies := checker.FindAvailableMovies(movies)
	notify(availableMovies)
}

func notify(availableMovies <-chan string) {
	for movie := range availableMovies {
		message := fmt.Sprintf("Tickets to %s are available: %s", movie, checker.MovieUrl(movie))
		log.Println(message)

		if len(webhookURL) > 0 {
			err := webhook.Send(webhookURL, message)
			if err != nil {
				log.Fatalf("error sending webhook: %s", err)
			}
		}
	}
}
