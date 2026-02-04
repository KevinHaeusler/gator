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
