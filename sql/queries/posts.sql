-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT
    P.*,
    FD.name as feed_name,
    FD.url as feed_url
FROM
    posts P
    INNER JOIN feeds FD ON P.feed_id = FD.id
    INNER JOIN follows FW ON FD.id = FW.feed_id
WHERE
    FW.user_id = $1
ORDER BY
    published_at DESC
LIMIT
    $2;