-- name: CreateFollow :one
INSERT INTO follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFollows :many
SELECT 
    * 
FROM 
    follows
WHERE
    user_id = $1
ORDER BY 
    created_at DESC;