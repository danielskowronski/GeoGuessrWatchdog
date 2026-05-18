# GeoGuessrWatchdog

**work in progress**

Periodically check GeoGuessr API for various things (like current map in Competitive Duels and when was it updated) and send notifications (e.g. to Discord) when state changes. Written as exercise in Temporal and Go.

## Local tests

```bash
brew install sqlc

# TODO: makefile?

cd internal/db && sqlc generate; cd -
```


```bash
docker compose up --build -d
docker compose run --rm worker help
docker compose run --rm worker trigger-workflow FetchUserStatsAndProgress
```

example `ggwd.env`:

```
GGWD_geoguessrApi_cookie=...
```

example `vpn.env` for NordVPN:

```
# 1. Go to https://my.nordaccount.com/pl/dashboard/nordvpn/access-tokens/ and get API token
# 2. Run `curl https://api.nordvpn.com/v1/users/services/credentials -u 'token:API_TOKEN'`
# 3. Extract `nordlynx_private_key` and set to WIREGUARD_PRIVATE_KEY

WIREGUARD_PRIVATE_KEY="..."
VPN_SERVICE_PROVIDER="nordvpn"
VPN_TYPE="wireguard"
SERVER_COUNTRIES="Poland"
```

example `postgres_app.env`:

```
POSTGRES_USER=app
POSTGRES_PASSWORD=app
POSTGRES_DB=app
```

example `postgres_temporal.env`:

```
POSTGRES_USER=temporal

POSTGRES_PASSWORD=temporal
POSTGRES_PWD=temporal

POSTGRES_DB=temporal
DB=temporal
```

also update `config_dockercompose.yaml`

### External notifications

Follow https://containrrr.dev/shoutrrr/v0.8/services/discord/ (or other target in Shoutrrr) to get URI. Set it in env var `GGWD_notifierApi_shoutrrr` (multiple URLs - comma separated).

Each URI can also have following params:

- `ggwd_prepend` - text added before each message body; useful to tag some people
- `ggwd_prepend` - text added after each message body; useful to tag some people

Those `ggwd_*` URI params are not valid Shoutrrr parameters and are stripped before passing to library. Example:

```
GGWD_notifierApi_shoutrrr="discord://TOKEN@WEBHOOK?ggwd_append=%20appended&ggwd_prepend=prepended%20"
```
