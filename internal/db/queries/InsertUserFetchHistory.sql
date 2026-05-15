-- name: InsertUserFetchHistory :one
INSERT INTO user_fetch_history (
  user_id
) VALUES (
  $1
)
RETURNING id;