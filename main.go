package main

import (
	"database/sql"

	"github.com/KevinHaeusler/gator/internal/database"
	_ "github.com/lib/pq"
)
import (
	"fmt"
	"os"

	"github.com/KevinHaeusler/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		_ = db.Close()
	}()

	dbQueries := database.New(db)
	stateGator := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&stateGator, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
