-- name: CreatePost :one
/*
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url VARCHAR(150) NOT NULL, -- is this enough?
    description VARCHAR(250),
    published_at TIMESTAMP,
    feed_id UUID,
    UNIQUE(url)
*/
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT
    *
FROM
    posts P
    INNER JOIN feeds FD ON P.feed_id = FD.id
    INNER JOIN follows FW ON FD.id = FW.feed_id
WHERE
    FW.user_id = $1
ORDER BY
    published_at DESC;