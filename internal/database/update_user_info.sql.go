// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: update_user_info.sql

package database

import (
	"context"
)

const updateUserInfo = `-- name: UpdateUserInfo :exec
UPDATE users
SET hashed_password = $2, email = $3
FROM refresh_tokens
WHERE users.id = refresh_tokens.user_id
and refresh_tokens.token = $1
`

type UpdateUserInfoParams struct {
	Token          string
	HashedPassword string
	Email          string
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateUserInfo, arg.Token, arg.HashedPassword, arg.Email)
	return err
}
