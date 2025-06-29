package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dr-check/blog-aggregator/internal/config"

	"github.com/dr-check/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

const databaseURL = "postgres://postgres:postgres@localhost:5432/gator"

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	jsonConfig, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read from file: %v", err)
	}

	db, err := sql.Open("postgres", jsonConfig.DbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	defer db.Close()

	dbQueries := database.New(db)

	newState := &state{
		db:  dbQueries,
		cfg: jsonConfig,
	}

	cmds := commands{
		availableCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

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
