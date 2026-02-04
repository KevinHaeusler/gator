package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	_ = cmd
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	fmt.Println("Database has been reset!")
	return nil
}
