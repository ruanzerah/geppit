// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rename_user.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const renameUser = `-- name: RenameUser :one
UPDATE users
  SET username = $2
  WHERE id = $1
RETURNING id, username, email, hash, balance, created_at, updated_at
`

type RenameUserParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (q *Queries) RenameUser(ctx context.Context, arg RenameUserParams) (User, error) {
	row := q.db.QueryRow(ctx, renameUser, arg.ID, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Hash,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}