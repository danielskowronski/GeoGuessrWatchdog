-- 00002_fetch_status.sql

-- +goose Up
CREATE TABLE IF NOT EXISTS fetch_status (
  fetch_type TEXT PRIMARY KEY,
  last_success TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- +goose Down
-- intentionally empty
