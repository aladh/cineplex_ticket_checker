package main

import (
	"flag"
	"log"
	"strings"

	"github.com/ali-l/cineplex_ticket_checker/checker"
)

var movies []string
var theatreIDs *string

func init() {
	theatreIDs = flag.String("t", "", "A comma-separated list of theatre IDs to look for")
	flag.Parse()

	movies = strings.Split(flag.Arg(0), ",")
}

func main() {
	availableMovies, err := checker.AvailableMovies(movies, *theatreIDs)
	if err != nil {
		log.Fatalf("error checking availability: %s\n", err)
	}

	for _, movie := range availableMovies {
		log.Printf("Tickets to %s are available\n", movie)
	}

	if len(availableMovies) > 0 {
		log.Fatalln("Go buy tickets!")
	}
}
