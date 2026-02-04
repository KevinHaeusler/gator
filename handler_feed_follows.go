package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KevinHaeusler/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: <feed-url>")

	}
	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Printf("User %s following feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	_ = cmd
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	fmt.Printf("Following:\n")
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}
	for _, feedFollow := range following {
		fmt.Printf("- %s\n", feedFollow.FeedName)
	}
	return nil
}
