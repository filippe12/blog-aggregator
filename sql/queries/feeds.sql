-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)

RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: PrettyListFeeds :many
SELECT feeds.name, feeds.url, users.name AS username
FROM feeds
JOIN users
ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    updated_at = $2,
    last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT feeds.*
FROM feed_follows
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY feeds.last_fetched_at NULLS FIRST
LIMIT 1;

-- name: DeleteFeeds :exec
DELETE FROM feeds WHERE true;
