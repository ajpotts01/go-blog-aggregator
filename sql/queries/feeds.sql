-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, last_fetched_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT
    *
FROM
    feeds
ORDER BY
    last_fetched_at NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE
    feeds
SET
    last_fetched_at = now()::timestamp(0),
    updated_at = now()::timestamp(0)
WHERE
    id = $1;