-- name: UpsertFetchStatus :one
INSERT INTO fetch_status (
  fetch_type,
  last_success
) VALUES ($1, now())
ON CONFLICT (fetch_type) DO UPDATE SET
  last_success = EXCLUDED.last_success
RETURNING last_success;