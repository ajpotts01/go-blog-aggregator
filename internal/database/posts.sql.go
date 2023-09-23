// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
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
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      uuid.NullUUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsByUser = `-- name: GetPostsByUser :many
SELECT
    p.id, p.created_at, p.updated_at, title, p.url, description, published_at, p.feed_id, fd.id, fd.created_at, fd.updated_at, name, fd.url, fd.user_id, last_fetched_at, fw.id, fw.created_at, fw.updated_at, fw.feed_id, fw.user_id
FROM
    posts P
    INNER JOIN feeds FD ON P.feed_id = FD.id
    INNER JOIN follows FW ON FD.id = FW.feed_id
WHERE
    FW.user_id = $1
ORDER BY
    published_at DESC
`

type GetPostsByUserRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Url           string
	Description   sql.NullString
	PublishedAt   sql.NullTime
	FeedID        uuid.NullUUID
	ID_2          uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	Name          string
	Url_2         string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	ID_3          uuid.UUID
	CreatedAt_3   time.Time
	UpdatedAt_3   time.Time
	FeedID_2      uuid.UUID
	UserID_2      uuid.UUID
}

func (q *Queries) GetPostsByUser(ctx context.Context, userID uuid.UUID) ([]GetPostsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByUserRow
	for rows.Next() {
		var i GetPostsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
			&i.Url_2,
			&i.UserID,
			&i.LastFetchedAt,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.FeedID_2,
			&i.UserID_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
