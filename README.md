# GeoGuessrWatchdog

**work in progress**

Periodically check GeoGuessr API for various things (like current map in Competitive Duels and when was it updated) and send notifications (e.g. to Discord) when state changes. Written as exercise in Temporal and Go.

## Dev

### Requirements

```bash
brew install sqlc make
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
