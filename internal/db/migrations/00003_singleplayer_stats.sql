-- 00003_singleplayer_stats.sql

-- +goose Up
CREATE TABLE IF NOT EXISTS user_singleplayer_history (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  fetch_id BIGINT NOT NULL REFERENCES user_fetch_history(id),

  games_played BIGINT NOT NULL,
  rounds_played BIGINT NOT NULL,

  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP VIEW IF EXISTS user_fetch_combined_history;
CREATE OR REPLACE VIEW user_fetch_combined_history AS
SELECT
  ufh.id AS fetch_id,
  ufh.user_id,
  ufh.timestamp AS fetch_timestamp,

  uph.division_name,
  uph.division_number,
  uph.rating_overall,
  uph.rating_moving,
  uph.rating_nomove,
  uph.rating_nmpz,
  uph.guessed_first,
  uph.best_countries,
  uph.worst_countries,

  ush.ranked_team_moving_games,
  ush.ranked_team_moving_wins,
  ush.ranked_team_nomove_games,
  ush.ranked_team_nomove_wins,
  ush.ranked_team_nmpz_games,
  ush.ranked_team_nmpz_wins,

  ush.ranked_solo_moving_games,
  ush.ranked_solo_moving_wins,
  ush.ranked_solo_nomove_games,
  ush.ranked_solo_nomove_wins,
  ush.ranked_solo_nmpz_games,
  ush.ranked_solo_nmpz_wins,

  ush.unranked_solo_moving_games,
  ush.unranked_solo_moving_wins,
  ush.unranked_solo_nomove_games,
  ush.unranked_solo_nomove_wins,
  ush.unranked_solo_nmpz_games,
  ush.unranked_solo_nmpz_wins,

  u1h.games_played AS singleplayer_games_played,
  u1h.rounds_played AS singleplayer_rounds_played

FROM user_fetch_history ufh
LEFT JOIN user_progress_history uph
  ON uph.fetch_id = ufh.id
LEFT JOIN user_stats_history ush
  ON ush.fetch_id = ufh.id
LEFT JOIN user_singleplayer_history u1h 
  ON u1h.fetch_id = ufh.id;

-- +goose Down
-- intentionally empty
