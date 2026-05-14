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

### Discord bot

Assuming channel on server

1. Go to https://discord.com/developers/applications/
2. Create bot
3. In *Install* disable public link
4. In *Bot* set:
   1. *Public* to disabled
   2. *Permissions* to *send messages*
   3. Click *Token* to get token
5. In *OAuth* set:
   1. *Scopes* to *bot*
   2. *Permissions* to *send messages*
   3. *Type* to *guild* (this is server)
   4. Open generated link and install on desired server
6. Obtain channel ID (right click on channel and select relevant option)

set overrides via env:

```bash
export GGWD_notifierApi_target=discord
export GGWD_notifierApi_discord_botToken='LONG_TOKEN_HERE'
export GGWD_notifierApi_discord_receivers=NUMERICAL_ID_1,NUMERICAL_ID_2
```
