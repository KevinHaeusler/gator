package main

import (
	"fmt"
	"log"

	"github.com/KevinHaeusler/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	if err := cfg.SetUser("kevin"); err != nil {
		log.Fatalf("set user: %v", err)
	}

	cfg2, err := config.Read()
	if err != nil {
		log.Fatalf("read config again: %v", err)
	}

	fmt.Printf("%+v\n", cfg2)
}
