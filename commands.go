package main

import (
	"context"
	"database/sql"
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
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist")

		} else {
			return fmt.Errorf("error checking for existing user: %w", err)
		}
	} else {
		err := s.cfg.SetUser(cmd.args[0])
		if err != nil {
			return fmt.Errorf("failed to set user: %w", err)
		}
		fmt.Printf("User %s has been set!", cmd.args[0])
	}

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
		if err == sql.ErrNoRows {
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
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	fmt.Println("Database has been reset!")
	return nil
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
