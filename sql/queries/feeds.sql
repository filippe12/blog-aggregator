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

-- name: PrettyListFeeds :many
SELECT feeds.name, feeds.url, users.name AS username
FROM feeds
JOIN users
ON feeds.user_id = users.id;


-- name: DeleteFeeds :exec
DELETE FROM feeds WHERE true;
