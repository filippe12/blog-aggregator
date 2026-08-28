# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows you to register users, follow RSS feeds, periodically fetch new posts, and browse posts from the feeds you follow.

## Requirements

Before running Gator, make sure you have the following installed:

* **Go** — Gator is written in Go and uses `go install` to install the CLI.
* **PostgreSQL** — Gator stores users, feeds, feed follows, and posts in a PostgreSQL database.

You can verify your Go installation with:

```bash
go version
```

and your PostgreSQL installation with:

```bash
psql --version
```

## Installing Gator

From the root of the project, install the CLI with:

```bash
go install ./cmd/gator
```

This builds the `gator` executable and installs it into your Go binary directory.

Make sure your Go binary directory is in your `PATH`. You can check it with:

```bash
go env GOPATH
```

For example, if your Go binary directory is not already in your `PATH`, you can add it with:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

After installation, verify that Gator is available:

```bash
gator
```

## Configuration

Gator uses a configuration file in your home directory:

```text
~/.gatorconfig.json
```

The configuration file should contain your PostgreSQL connection string and the currently logged-in user.

For example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the PostgreSQL username, password, database name, and other connection details with the values for your local PostgreSQL installation.

Before using Gator, make sure the database exists and the database migrations have been applied.

## Basic Usage

### Register a user

```bash
gator register filip
```

Registering a user also makes that user the currently active user.

### Login as a user

```bash
gator login filip
```

You can see all registered users with:

```bash
gator users
```

The currently logged-in user is marked with `*`.

### Add an RSS feed

```bash
gator addfeed "Hacker News" "https://hnrss.org/frontpage"
```

Adding a feed also automatically follows it for the current user.

### List feeds

```bash
gator feeds
```

### Follow an existing feed

```bash
gator follow "https://hnrss.org/frontpage"
```

### See feeds you follow

```bash
gator following
```

### Unfollow a feed

```bash
gator unfollow "https://hnrss.org/frontpage"
```

### Fetch posts

The `agg` command periodically fetches RSS feeds.

The argument is a Go duration:

```bash
gator agg 10s
```

This will fetch feeds every 10 seconds.

Other examples:

```bash
gator agg 1m
gator agg 30s
gator agg 500ms
```

The aggregator runs continuously until you stop it with `Ctrl+C`.

### Browse posts

After the aggregator has fetched some posts, you can browse posts from feeds you follow:

```bash
gator browse
```

By default, Gator displays 2 posts. You can provide a different limit:

```bash
gator browse 10
```

### Reset the database

To delete all users, feeds, follows, and posts:

```bash
gator reset
```

**Warning:** this permanently deletes the data from the Gator database.
