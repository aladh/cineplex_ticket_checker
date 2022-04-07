# Cineplex Ticket Checker
Sends a webhook when tickets are on sale for a movie.

### Configuration

This repo is mirrored to GitLab, and runs a scheduled job to periodically check for new movies.

The CLI can be used as follows:

```shell
cineplex_ticket_checker -w https://webook.com example-movie1,example-movie2
```

`w` is a URL to send webhooks to when movies are available

`example-movie1` is the path portion of a movie URL (ex. https://www.cineplex.com/Movie/example-movie1)
