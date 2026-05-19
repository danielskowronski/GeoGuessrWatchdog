TAG ?= dev
PLATFORM ?= linux/amd64,linux/arm64

.PHONY: sqlc local-run local-trigger go-mod build push build-and-push

go-mod:
	go mod download
	go mod tidy

sqlc:
	cd internal/db && sqlc generate

local-run:
	cd dev && docker compose up --build -d

local-trigger:
	cd dev && docker compose run --rm worker trigger-workflow $(WORKFLOW)

build: sqlc go-mod
	docker buildx build --platform $(PLATFORM) -t ghcr.io/danielskowronski/geoguessrwatchdog:$(TAG) .

push:
	docker push ghcr.io/danielskowronski/geoguessrwatchdog:$(TAG)

build-and-push: build push
