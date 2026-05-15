-- name: InsertUserProgressHistory :one
INSERT INTO user_progress_history (
  user_id,
  fetch_id,
  division_name,
  division_number,
  rating_overall,
  rating_moving,
  rating_nomove,
  rating_nmpz,
  guessed_first,
  best_countries,
  worst_countries
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id;