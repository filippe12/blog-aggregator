package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/filippe12/blog-aggregator/internal/database"
	"github.com/filippe12/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handler[cmd.name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required for login command")
	}
	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user %s does not exist yet, to register user run register <username>", username)
	}

	s.cfg.SetUser(username)
	fmt.Println("username set to:", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required for register command")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("registered user", cmd.arguments[0])
	fmt.Println("debug:", user)

	s.cfg.SetUser(cmd.arguments[0])
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}
	}
	return nil
}

func handlerAgg(_ *state, _ command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Print(feed)
	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		log.Fatal(err)
	}

	return nil
}
