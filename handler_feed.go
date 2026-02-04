package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KevinHaeusler/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: <name> <url>")
	}
	ctx := context.Background()

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	printFeed(feed)
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	_ = cmd
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}
	for _, feed := range feeds {
		printFeedRow(feed)
	}
	return nil
}

func printFeedRow(row database.GetFeedsRow) {
	fmt.Printf("* ID:      %v\n", row.FeedID)
	fmt.Printf("* Name:    %s\n", row.FeedName)
	fmt.Printf("* URL:     %s\n", row.FeedUrl)
	fmt.Printf("* User Name:  %v\n", row.UserName)
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:      %v\n", feed.ID)
	fmt.Printf("* Name:    %s\n", feed.Name)
	fmt.Printf("* URL:     %s\n", feed.Url)
	fmt.Printf("* User ID:  %v\n", feed.UserID)
}
