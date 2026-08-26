package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/filippe12/blog-aggregator/internal/config"
	"github.com/filippe12/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg := config.Read()
	s := state{cfg: &cfg}
	commands := commands{
		handler: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddfeed)
	commands.register("feeds", handlerFeeds)
	commands.register("reset", handlerReset)

	cliArguments := os.Args
	if len(cliArguments) < 2 {
		log.Fatal("not enough arguments")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	command := command{name: cliArguments[1], arguments: cliArguments[2:]}
	err = commands.run(&s, command)
	if err != nil {
		log.Fatal(err)
	}
}
