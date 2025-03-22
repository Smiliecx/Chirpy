-- name: UserOwnsChirp :one
SELECT EXISTS (
  SELECT 1
  FROM chirps
  WHERE chirps.id = $1 and chirps.user_id = $2
);