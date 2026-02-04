package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/KevinHaeusler/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided for login command")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist")

		}

		return fmt.Errorf("error checking for existing user: %w", err)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User %s has been set!", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided for register command")
	}

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user, err := s.db.CreateUser(context.Background(), createUserParams)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			fmt.Printf("User %s has been registered!", user.Name)
			err = s.cfg.SetUser(cmd.args[0])
			if err != nil {
				return fmt.Errorf("failed to set user: %w", err)
			}
		} else {
			return fmt.Errorf("error checking for existing user: %w", err)
		}
	} else {
		return fmt.Errorf("user already exists: %s", user.Name)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	_ = cmd
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	fmt.Println("Database has been reset!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	_ = cmd
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("- %s (current)\n", user.Name)
		} else {
			fmt.Printf("- %s\n", user.Name)
		}
	}
	return nil
}

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

func (c *commands) run(s *state, cmd command) error {
	if cmd.name == "" {
		return fmt.Errorf("no command name provided")
	}
	handler, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
