-- 00001_init.sql

-- +goose Up
CREATE TABLE IF NOT EXISTS division_info (
  id BIGSERIAL PRIMARY KEY,
  division_name TEXT NOT NULL,
  game_mode TEXT NOT NULL,
  map_id TEXT NOT NULL,
  map_name TEXT NOT NULL,
  last_changed TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS division_info_selector_idx
  ON division_info (division_name, game_mode);

CREATE TABLE IF NOT EXISTS division_history (
  id BIGSERIAL PRIMARY KEY,
  division_id BIGINT NOT NULL REFERENCES division_info(id),
  map_id TEXT NOT NULL,
  map_name TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS division_history_division_id_idx
  ON division_history (division_id);

CREATE TABLE IF NOT EXISTS map_info (
  id BIGSERIAL PRIMARY KEY,
  api_id TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,

  bounds_min_lat DOUBLE PRECISION,
  bounds_min_lon DOUBLE PRECISION,
  bounds_max_lat DOUBLE PRECISION,
  bounds_max_lon DOUBLE PRECISION,
  max_err_distance BIGINT,

  api_updated TIMESTAMPTZ NOT NULL DEFAULT now(),

  location_count BIGINT NOT NULL DEFAULT 0,

  last_changed TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS map_info_map_id_idx
  ON map_info (api_id);

CREATE TABLE IF NOT EXISTS map_history (
  id BIGSERIAL PRIMARY KEY,
  map_id BIGINT NOT NULL REFERENCES map_info(id),

  bounds_min_lat DOUBLE PRECISION,
  bounds_min_lon DOUBLE PRECISION,
  bounds_max_lat DOUBLE PRECISION,
  bounds_max_lon DOUBLE PRECISION,
  max_err_distance BIGINT,

  api_updated TIMESTAMPTZ NOT NULL DEFAULT now(),

  location_count BIGINT NOT NULL DEFAULT 0,

  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS map_history_map_id_idx
  ON map_history (map_id);

CREATE TABLE IF NOT EXISTS user_fetch_history (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS user_fetch_history_user_id_idx
  ON user_fetch_history (user_id);

CREATE TABLE IF NOT EXISTS user_progress_history (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  fetch_id BIGINT NOT NULL REFERENCES user_fetch_history(id),

  division_name TEXT NOT NULL,
  division_number INT NOT NULL,
  rating_overall INT NOT NULL,
  rating_moving INT NOT NULL,
  rating_nomove INT NOT NULL,
  rating_nmpz INT NOT NULL,
  guessed_first DOUBLE PRECISION NOT NULL,
  best_countries TEXT NOT NULL,
  worst_countries TEXT NOT NULL,

  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS user_progress_history_user_id_idx
  ON user_progress_history (user_id);

CREATE TABLE IF NOT EXISTS user_stats_history (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  fetch_id BIGINT NOT NULL REFERENCES user_fetch_history(id),

  ranked_team_moving_games BIGINT NOT NULL,
  ranked_team_moving_wins  BIGINT NOT NULL,
  ranked_team_nomove_games BIGINT NOT NULL,
  ranked_team_nomove_wins  BIGINT NOT NULL,
  ranked_team_nmpz_games BIGINT NOT NULL,
  ranked_team_nmpz_wins  BIGINT NOT NULL,

  ranked_solo_moving_games BIGINT NOT NULL,
  ranked_solo_moving_wins  BIGINT NOT NULL,
  ranked_solo_nomove_games BIGINT NOT NULL,
  ranked_solo_nomove_wins  BIGINT NOT NULL,
  ranked_solo_nmpz_games BIGINT NOT NULL,
  ranked_solo_nmpz_wins  BIGINT NOT NULL,

  unranked_solo_moving_games BIGINT NOT NULL,
  unranked_solo_moving_wins  BIGINT NOT NULL,
  unranked_solo_nomove_games BIGINT NOT NULL,
  unranked_solo_nomove_wins  BIGINT NOT NULL,
  unranked_solo_nmpz_games BIGINT NOT NULL,
  unranked_solo_nmpz_wins  BIGINT NOT NULL,

  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS user_stats_history_user_id_idx
  ON user_stats_history (user_id);

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
  ush.unranked_solo_nmpz_wins

FROM user_fetch_history ufh
LEFT JOIN user_progress_history uph
  ON uph.fetch_id = ufh.id
LEFT JOIN user_stats_history ush
  ON ush.fetch_id = ufh.id;

CREATE OR REPLACE VIEW division_map_info AS
SELECT
  di.division_name as division_name,
  di.game_mode as game_mode,
  di.map_id as map_id,
  di.map_name as map_name,
  mi.bounds_min_lat as map_bounds_min_lat,
  mi.bounds_min_lon as map_bounds_min_lon,
  mi.bounds_max_lat as map_bounds_max_lat,
  mi.bounds_max_lon as map_bounds_max_lon,
  mi.max_err_distance as map_max_err_distance,
  mi.api_updated as map_api_updated,
  mi.location_count as map_location_count
FROM division_info di
JOIN map_info mi
  ON mi.api_id = di.map_id;

CREATE TABLE IF NOT EXISTS division_map_last_notifications (
  id BIGSERIAL PRIMARY KEY,

  division_name TEXT NOT NULL,
  game_mode TEXT NOT NULL,

  map_id TEXT NOT NULL,
  map_bounds_min_lat DOUBLE PRECISION,
  map_bounds_min_lon DOUBLE PRECISION,
  map_bounds_max_lat DOUBLE PRECISION,
  map_bounds_max_lon DOUBLE PRECISION,
  map_max_err_distance BIGINT,
  map_api_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
  location_count BIGINT NOT NULL DEFAULT 0,

  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS division_map_last_notifications_selector_idx
  ON division_map_last_notifications (division_name, game_mode);

-- +goose Down
-- intentionally empty
