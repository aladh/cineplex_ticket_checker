# Cineplex Ticket Checker

Checks if tickets are available for one or more movies.

## Building

Ensure that `go` is installed before attempting to build this project. 

This project can be built by running `go build` in the root directory, which will produce an executable named
`cineplex_ticket_checker`.

## Usage

This repo is mirrored to GitLab and runs a scheduled job to periodically check for new movies. See the
`check-availability` job in `.gitlab-ci.yml` for CI configuration.

The CLI can be used as follows:

```shell
cineplex_ticket_checker -w https://webook.com example-movie1,example-movie2
```
- Movie names are taken from the name portion of a movie URL (ex. https://www.cineplex.com/Movie/example-movie1), and
  should be comma-separated.
- Webhooks can optionally be sent when movies are available. The webhook URL is specified with `-w`. The webhook payload
  is compatible with Discord.
