-- name: CreateUser :one
INSERT INTO users (
  username, email, hash
) VALUES ( $1, $2, $3 )
RETURNING *;
