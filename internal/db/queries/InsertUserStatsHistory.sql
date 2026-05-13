-- name: InsertUserStatsHistory :one
INSERT INTO user_stats_history (
  user_id,
  ranked_team_moving_games,
  ranked_team_moving_wins,
  ranked_team_moving_ratio,
  ranked_team_nomove_games,
  ranked_team_nomove_wins,
  ranked_team_nomove_ratio,
  ranked_team_nmpz_games,
  ranked_team_nmpz_wins,
  ranked_team_nmpz_ratio,
  ranked_solo_moving_games,
  ranked_solo_moving_wins,
  ranked_solo_moving_ratio,
  ranked_solo_nomove_games,
  ranked_solo_nomove_wins,
  ranked_solo_nomove_ratio,
  ranked_solo_nmpz_games,
  ranked_solo_nmpz_wins,
  ranked_solo_nmpz_ratio,
  unranked_solo_moving_games,
  unranked_solo_moving_wins,
  unranked_solo_moving_ratio,
  unranked_solo_nomove_games,
  unranked_solo_nomove_wins,
  unranked_solo_nomove_ratio,
  unranked_solo_nmpz_games,
  unranked_solo_nmpz_wins,
  unranked_solo_nmpz_ratio 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
)
RETURNING id;