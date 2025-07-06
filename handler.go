package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
	fmt.Printf(" * LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
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
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("please enter a proper command: %v <time_between_reqs>", cmd.Name)
	}

	timeReq, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("parsing feeds every %s...", timeReq)

	ticker := time.NewTicker(timeReq)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("please enter a proper command: %s <name> <URL>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	_, err2 := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err2 != nil {
		return fmt.Errorf("error following feed: %w", err2)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please enter a proper command: %s <URL>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error locating feed: %w", err)
	}

	_, err2 := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err2 != nil {
		return fmt.Errorf("error following feed: %w", err2)
	}

	fmt.Println(feed.Name)
	fmt.Println(user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("please enter a proper command: %v", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}
	for i := 0; i < len(follows); i++ {
		fmt.Printf("* %v\n", follows[i].FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please enter a proper command: %s <URL>", cmd.Name)
	}

	url := cmd.Args[0]

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  url,
	})
	if err != nil {
		return fmt.Errorf("error locating feed: %w", err)
	}

	fmt.Println("successfully unfollowed feed")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("unable to get posts for user: %w", err)
	}

	fmt.Printf("found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("     %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
