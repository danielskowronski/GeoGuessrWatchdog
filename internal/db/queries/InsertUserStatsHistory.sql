-- name: InsertUserStatsHistory :one
INSERT INTO user_stats_history (
  user_id,
  ranked_team_moving_games,
  ranked_team_moving_wins,
  ranked_team_nomove_games,
  ranked_team_nomove_wins,
  ranked_team_nmpz_games,
  ranked_team_nmpz_wins,
  ranked_solo_moving_games,
  ranked_solo_moving_wins,
  ranked_solo_nomove_games,
  ranked_solo_nomove_wins,
  ranked_solo_nmpz_games,
  ranked_solo_nmpz_wins,
  unranked_solo_moving_games,
  unranked_solo_moving_wins,
  unranked_solo_nomove_games,
  unranked_solo_nomove_wins,
  unranked_solo_nmpz_games,
  unranked_solo_nmpz_wins
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
)
RETURNING id;