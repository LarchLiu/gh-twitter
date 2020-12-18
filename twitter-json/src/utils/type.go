package utils

import twitterscraper "github.com/n0madic/twitter-scraper"

// Twitter represents a struct of twitter.
type Twitter struct {
	Profile twitterscraper.Profile
	Tweets  []twitterscraper.Tweet
}
