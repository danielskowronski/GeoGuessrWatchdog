# GeoGuessrWatchdog

**work in progress**

Periodically check GeoGuessr API for various things (like current map in Competitive Duels and when was it updated) and send notifications (e.g. to Discord) when state changes. Written as exercise in Temporal and Go.

## Dev

### Requirements

```bash
brew install sqlc make

docker buildx create --name local --driver docker-container --use
```

### Code generation

```bash
make sqlc
```

### Local dev env

Set up according to [dev/README.md](./dev/README.md)

```bash
make local-run

make local-trigger WORKFLOW=FetchDivisionsMaps
make local-trigger WORKFLOW=FetchUserStatsAndProgress
```

Endpoints:

- [Temporal](http://localhost:8080/)
- [cache](http://localhost:10001/)
- [DB](postgres://app:app@localhost:5432/app)
- [API](http://localhost:8099/docs)
- [Frontend](http://localhost:8099/)

### Build and push image

All take params:

- `TAG` - tag of Docker image, defaults to `dev`
- `PLATFORM` - list of buildx platforms to use, defaults to `linux/arm64,linux/amd64`

```bash
make build
make push
make build-and-push

# quick build
make build-and-push PLATFORM=linux/arm64 TAG=stage
```
