# Blog Aggregator
This blog aggregator, henceforth nicknamed Gator, is a multi-user CLI application for logging and browsing RSS feeds. It's designed for local, single-device use but supports multiple user profiles which can represent different users or even different interests if you wished to separate tech blogs from entertainment blogs for example.

## Prerequisites
Before you begin, ensure you have met the following requirements:

You have installed the latest version of Go
You have a PostgreSQL database set up and running

##Installing Gator
To install Gator, follow these steps:

```
go install github.com/dr-check/blog-aggregator@latest
```

## Configuration

Create a JSON configuration file named config.json in the same directory as Gator. The file should have the following structure:

```json
{
  "db_url": "your_postgresql_connection_string",
  "current_user_name": "username_goes_here"
}
```

Replace your_postgresql_connection_string with your actual PostgreSQL connection string.

## Using Gator

Gator provides several commands for managing users, feeds, and posts:

*`login`: Log in as a user
*`register`: Register a new user
*`reset`: Reset the database
*`users`: List all users
*`agg`: Fetch RSS feeds
*`addfeed`: Add a new RSS feed (requires login)
*`feeds`: List all feeds
*`follow`: Follow a feed (requires login)
*`following`: List followed feeds (requires login)
*`unfollow`: Unfollow a feed (requires login)
*`browse`: Browse posts

To use a command, run:

`go run . [command]`

For commands that require login, ensure you're logged in first using the login command.

## Note on Security
This application doesn't include user-based authorization. Anyone with database credentials can act as any user. This design is intentional for learning purposes, focusing on SQL, CLIs, and long-running services.
