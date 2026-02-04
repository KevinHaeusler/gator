# Gator - RSS Feed Aggregator CLI

Gator is a command-line tool for aggregating and managing RSS feeds. It allows you to follow multiple feeds, browse
posts, and manage your feed subscriptions.

## Prerequisites

Before running Gator, you need to have the following installed:

- **Go** (version 1.19 or higher) - [Download Go](https://go.dev/dl/)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install the Gator CLI using `go install`:

## Configure gator

`gator` reads its config from a JSON file in your home directory:

- **Path:** `~/.gatorconfig.json`

Create it with contents like:

```json
{ "db_url": "postgres:// : @ :/?sslmode=disable", "current_user_name": "" }
```




## Useful commands

A few commands you can try:

- `gator register <username>`  
  Creates a user (and typically sets you as the current user).

- `gator login <username>`  
  Sets the current user in `~/.gatorconfig.json`.

- `gator users`  
  Lists users; marks the currently selected one.

- `gator addfeed <name> <url>`  
  Adds a feed for the current user (requires being logged in).

- `gator feeds`  
  Lists all feeds.

- `gator follow <feed-url>` / `gator unfollow <feed-url>`  
  Follow or unfollow an existing feed (requires being logged in).

- `gator following`  
  Shows which feeds the current user follows.

- `gator agg <time_between_reqs>`  
  Periodically fetches feeds and stores posts, e.g.
  ```bash
  gator agg 30s
  ```

- `gator browse [limit]`  
  Shows posts for the current user. Optionally set a limit:
  ```bash
  gator browse 10
  ```
  

## Troubleshooting

- If you see config errors, confirm `~/.gatorconfig.json` exists and contains a valid `db_url`.
- If database queries fail, confirm Postgres is running and the schema/migrations have been applied.
- If `gator` isn’t found after `go install`, confirm `$(go env GOPATH)/bin` is on your `PATH`.