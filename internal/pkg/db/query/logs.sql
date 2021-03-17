
-- name: InsertLog :one
INSERT INTO logs (
  inserts,
  updates,
  errors,
  timestamp_started,
  timestamp_finished

) VALUES (
  $1, $2, $3, $4, $5
) 
RETURNING *;

   