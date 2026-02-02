package main

import (
	"github.com/KevinHaeusler/gator/internal/config"
	"github.com/KevinHaeusler/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
