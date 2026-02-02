package main

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

	stateGator := state{
		cfg: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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
