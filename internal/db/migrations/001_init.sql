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

CREATE TABLE IF NOT EXISTS user_progress_history (
  id BIGSERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,

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

-- TODO: table for notifier last sent notifications