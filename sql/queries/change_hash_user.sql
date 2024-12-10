-- name: ChangeHashUser :exec
UPDATE users
  SET hash = $2
  WHERE id = $1;

-- Scary query? idk.
