package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aladh/cineplex_ticket_checker/checker"
	"github.com/aladh/cineplex_ticket_checker/webhook"
)

var theatreIDs string
var webhookURL string

func main() {
	flag.StringVar(&theatreIDs, "t", "", "A comma-separated list of theatre IDs to look for")
	flag.StringVar(&webhookURL, "w", "", "A URL to send webhooks to when movies are available")
	flag.Parse()

	movies := strings.Split(flag.Arg(0), ",")
	availableChan := make(chan string, len(movies))

	go checker.FindAvailableMovies(movies, theatreIDs, availableChan)
	sendWebhooks(availableChan)
}

func sendWebhooks(availableChan <-chan string) {
	for movie := range availableChan {
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
