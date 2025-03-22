-- name: UpdateUserInfo :exec
UPDATE users
SET hashed_password = $2, email = $3
FROM refresh_tokens
WHERE users.id = refresh_tokens.user_id
and refresh_tokens.token = $1;