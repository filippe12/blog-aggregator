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

func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("not enough arguments, run: addfeed <name> <url>")
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedEntry, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedEntry.ID,
	})

	fmt.Println(feedEntry)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	prettyFeedEntries, err := s.db.PrettyListFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, entry := range prettyFeedEntries {
		fmt.Println(entry)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments, run: follow <url>")
	}
	url := cmd.arguments[0]
	username := s.cfg.CurrentUserName

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing(s *state, _ command) error {
	username := s.cfg.CurrentUserName
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), username)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollows {
		fmt.Println("*", feedFollow.FeedName)
	}

	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := s.db.DeleteFeeds(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := s.db.DeleteFeedFollows(context.Background()); err != nil {
		log.Fatal(err)
	}

	return nil
}
