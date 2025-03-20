-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens WHERE token=$1;