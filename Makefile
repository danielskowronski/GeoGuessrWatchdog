.PHONY: sqlc local-run local-trigger

sqlc:
	cd internal/db && sqlc generate

local-run:
	cd dev && docker compose up --build -d

local-trigger:
	cd dev && docker compose run --rm worker trigger-workflow $(WORKFLOW)