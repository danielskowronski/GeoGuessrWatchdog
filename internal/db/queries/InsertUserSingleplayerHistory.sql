-- name: InsertUserSingleplayerHistory :one
INSERT INTO user_singleplayer_history (
  user_id,
  fetch_id,
  games_played,
  rounds_played
) VALUES (
  $1, $2, $3, $4
)
RETURNING id;