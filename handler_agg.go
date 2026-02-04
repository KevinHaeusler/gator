package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KevinHaeusler/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("failed to get feeds: %v", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("failed to mark feed fetched: %v", err)
		return
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("failed to fetch feed: %v", err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Done! %s Feed: %v posts found\n", feed.Name, len(feedData.Channel.Item))
}
