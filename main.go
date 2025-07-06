package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dr-check/blog-aggregator/internal/config"
	"github.com/dr-check/blog-aggregator/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

//databaseURL = "postgres://postgres:postgres@localhost:5432/gator"

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
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error fetching user ID: %w", err)
		}
		return handler(s, cmd, currentUser)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("unable to fetch feeds", err)
		return
	}
	log.Println("new feed fetched")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("unable to update feed %s fetch: %v", feed.Name, err)
		return
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("unable to log feed %s: %v", feed.Name, err)
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("unable to create post: %v", err)
			continue
		}
	}
	log.Printf("feed %s added, %v new posts found", feed.Name, len(feedData.Channel.Item))
}
