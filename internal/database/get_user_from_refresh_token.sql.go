// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get_user_from_refresh_token.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens WHERE token=$1
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}
