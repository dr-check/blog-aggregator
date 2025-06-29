package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dr-check/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	jsonConfig, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read from file: %v", err)
	}

	newState := &state{
		cfg: jsonConfig,
	}

	cmds := commands{
		availableCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: CLI <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(newState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
