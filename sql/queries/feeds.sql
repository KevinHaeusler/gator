-- noinspection SqlResolveForFile

-- name: CreateFeed :one
INSERT INTO feeds ( created_at, updated_at, name, url, user_id)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5

       )
RETURNING *;

-- name: GetFeeds :many
SELECT
    feeds.id  AS feed_id,
    feeds.name  AS feed_name,
    feeds.url   AS feed_url,
    users.name  AS user_name
FROM feeds
         INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;