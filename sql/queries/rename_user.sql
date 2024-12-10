-- name: RenameUser :one
UPDATE users
  SET username = $2
  WHERE id = $1
RETURNING *;
