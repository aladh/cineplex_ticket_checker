package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ali-l/cineplex_ticket_checker/checker"
	"github.com/ali-l/cineplex_ticket_checker/webhook"
)

func main() {
	theatreIDs := flag.String("t", "", "A comma-separated list of theatre IDs to look for")
	webhookURL := flag.String("w", "", "A URL to send webhooks to when movies are available")
	flag.Parse()
	movies := strings.Split(flag.Arg(0), ",")

	availableMovies, err := checker.AvailableMovies(movies, *theatreIDs)
	if err != nil {
		log.Fatalf("error checking availability: %s\n", err)
	}

	for _, movie := range availableMovies {
		message := fmt.Sprintf("Tickets to %s are available: %s", movie, checker.MovieUrl(movie))
		log.Println(message)

		if len(*webhookURL) > 0 {
			err := webhook.Send(*webhookURL, message)
			if err != nil {
				log.Fatalf("error sending webhook: %s", err)
			}
		}
	}
}
