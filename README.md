# GeoGuessrWatchdog

**work in progress**

Periodically check GeoGuessr API for various things (like current map in Competitive Duels and when was it updated) and send notifications (e.g. to Discord) when state changes. Written as exercise in Temporal and Go.

## Dev

### Code generation

```bash
brew install sqlc

# TODO: makefile?

cd internal/db && sqlc generate; cd -
```

### Local dev env

[dev/README.md](./dev/README.md)
