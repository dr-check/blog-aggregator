package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dr-check/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please enter a proper command: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("unable to set username: %w", err)
	}

	fmt.Println("Username successfully changed!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:		%v\n", feed.ID)
	fmt.Printf(" * Created At:	%v\n", feed.CreatedAt)
	fmt.Printf(" * Updated At:	%v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:	%v\n", feed.Name)
	fmt.Printf(" * URL:	%v\n", feed.Url)
	fmt.Printf(" * User:	%v\n", feed.UserID)
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please enter a proper command: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	fmt.Println("user created successfully:")
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("please enter a proper command: %s", cmd.Name)
	}

	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("unable to delete users: %w", err)
	}
	fmt.Println("successfully deleted users")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("please enter a proper command: %v", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}
	for i := 0; i < len(users); i++ {
		if users[i].Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", users[i].Name)
		} else {
			fmt.Printf("* %v\n", users[i].Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("please enter a proper command: %v", cmd.Name)
	}
	rFeed, err := fetchFeed(context.TODO(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("please enter a proper command: %s <name> <URL>", cmd.Name)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetching user ID: %w", err)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("feed created successfully")
	printFeed(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("please enter a proper command: %v", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}
	for i := 0; i < len(feeds); i++ {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("unable to locate feed: %w", err)
		}
		fmt.Printf("* Name: %v\n", feeds[i].Name)
		fmt.Printf("* URL: %v\n", feeds[i].Url)
		fmt.Printf("* User ID: %v\n", user.Name)
	}
	return nil
}
