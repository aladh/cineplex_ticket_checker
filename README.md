# Cineplex Ticket Checker
Sends a webhook when tickets are on sale for a movie.

### Configuration

This repo is mirrored to GitLab, and runs a scheduled job to periodically check for new movies.

The CLI can be used as follows:

```shell
cineplex_ticket_checker -t 1234 -w https://webook.com example-movie1,example-movie2
```

`t` is a comma-separated list of theatre IDs to check for

`w` is a URL to send webhooks to when movies are available

`example-movie1` is the path portion of a movie URL (ex. https://www.cineplex.com/Movie/example-movie1)
