// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: post.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts ( title, content, category, image, created_at, updated_at) VALUES ( $1, $2, $3, $4, $5, $6) RETURNING id, title, category, content, image, created_at, updated_at
`

type CreatePostParams struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.queryRow(ctx, q.createPostStmt, createPost,
		arg.Title,
		arg.Content,
		arg.Category,
		arg.Image,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.Content,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deletePostStmt, deletePost, id)
	return err
}

const getPostById = `-- name: GetPostById :one
SELECT id, title, category, content, image, created_at, updated_at FROM posts WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPostById(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.queryRow(ctx, q.getPostByIdStmt, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.Content,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPosts = `-- name: ListPosts :many
SELECT id, title, category, content, image, created_at, updated_at FROM posts ORDER BY id LIMIT $1 OFFSET $2
`

type ListPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error) {
	rows, err := q.query(ctx, q.listPostsStmt, listPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Content,
			&i.Image,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updatePost = `-- name: UpdatePost :one
UPDATE posts set title = coalesce($1, title), category = coalesce($2, category), content = coalesce($3, content) , image = coalesce($4, image), updated_at = coalesce($5, updated_at ) WHERE id = $6 RETURNING id, title, category, content, image, created_at, updated_at
`

type UpdatePostParams struct {
	Title     sql.NullString `json:"title"`
	Category  sql.NullString `json:"category"`
	Content   sql.NullString `json:"content"`
	Image     sql.NullString `json:"image"`
	UpdatedAt sql.NullTime   `json:"updated_at "`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.queryRow(ctx, q.updatePostStmt, updatePost,
		arg.Title,
		arg.Category,
		arg.Content,
		arg.Image,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.Content,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
