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
