package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	_ = s
	_ = cmd

	url := "https://www.wagslane.dev/index.xml"

	var feed RSSFeed
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	fmt.Println(feed)
	return nil
}
