# Cineplex Ticker Checker
Sends you an email when tickets are on sale for a movie

### Configuration
Deployed as a Cron triggered AWS Lambda

- Set environment variable EMAIL_ADDRESS with email to notify
- Set environment variable THEATRE_IDS with JSON array of theatre IDs (ex. \["1", "2"]) you care about
- Fill in MOVIES constant in index.js with a list of movie paths to check
    - Example: "deadpool-2" from https://www.cineplex.com/Movie/deadpool-2 