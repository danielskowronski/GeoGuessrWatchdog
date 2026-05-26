TAG ?= dev
PLATFORM ?= linux/arm64,linux/amd64
IMAGE_NAME ?= ghcr.io/danielskowronski/geoguessrwatchdog
RELEASE_VERSION := $(shell tr -d '[:space:]' < VERSION)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: sqlc local-run local-trigger go-mod build push build-and-push release-prepare 

go-mod:
	go mod download
	go mod tidy

sqlc:
	cd internal/db && sqlc generate

local-run:
	cd dev && docker compose up --build -d

local-trigger:
	cd dev && docker compose run --rm worker trigger-workflow $(WORKFLOW)

api-direct-run:
	(cd cmd/api && GGWD_CONFIG_PATH=../../dev/api_config_direct.yaml go run .)

build: sqlc go-mod
	docker buildx build --platform $(PLATFORM) -t $(IMAGE_NAME):$(TAG) --build-arg RELEASE_VERSION=$(RELEASE_VERSION) --build-arg DATE=$(DATE) .

push:
	docker push $(IMAGE_NAME):$(TAG)

build-and-push: build push


release-prepare:
	@echo "Preparing release for version $(RELEASE_VERSION)"
	docker manifest inspect $(IMAGE_NAME):$(RELEASE_VERSION) > /dev/null && (echo "Image with tag $(RELEASE_VERSION) already exists. Please update the VERSION file." && exit 1) || echo "No existing image with tag $(RELEASE_VERSION), proceeding..."
	git tag -l | grep $(RELEASE_VERSION) > /dev/null && (echo "Git tag $(RELEASE_VERSION) already exists. Please update the VERSION file." && exit 1) || echo "No existing git tag $(RELEASE_VERSION), proceeding..."
	sed -i '' 's/appVersion: ".*"/appVersion: "$(RELEASE_VERSION)"/' charts/ggwd/Chart.yaml
	$(MAKE) build TAG=$(RELEASE_VERSION) VERSION=$(RELEASE_VERSION) DATE=$(DATE)
	git status

release-finish:
	git tag --points-at HEAD | grep -q . && (echo "HEAD already has a git tag. Please commit any changes before running this command." && exit 1) || echo "No existing git tag on HEAD, proceeding..."
	git tag $(RELEASE_VERSION)
	git push
	git push --tags
	$(MAKE) push TAG=$(RELEASE_VERSION)