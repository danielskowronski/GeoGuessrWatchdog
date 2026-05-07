# Prototype / PoC

## Goals

- check client IP
- query just divisions list and unmarshal results to list of flattened `DivisionModeMapInfo` struct (division name, game mode, map ID, map name)
- have option to use HTTP_PROXY but only on some parts of code (GeoGuessr via proxy, other APIs directly)
- package VPN provider which exposes HTTP proxy
- run in container with Docker Compose to allow easy porting to k8s

## Setup

1. Create `app.env` with `GGWD_NCFA=` set to `_ncfa` cookie value obtained from browser session after logging in to GeoGuessr
2. Create `vpn.env` with env vars for gluetun
3. Execute `docker compose up --build fetcher` 

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

command:

```bash
docker compose up --build fetcher 
```

example output:

```
[+] up 3/3
 ✔ Image prototype-fetcher       Built     0.7s
 ✔ Container prototype-gluetun-1 Running   0.0s
 ✔ Container prototype-fetcher-1 Recreated 0.1s
Attaching to fetcher-1
Container prototype-gluetun-1 Waiting 
Container prototype-gluetun-1 Healthy 
fetcher-1  | Real IP:  ...
fetcher-1  | Proxy IP: ...
fetcher-1  | DivisionModeMapInfo{DivisionName=Champion         GameMode=standard MapID=676340ae2f718dbabdf30331 MapName="An Informed World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Champion         GameMode=noMove   MapID=67101e71077612dd6ac8de79 MapName="Supersonic"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Champion         GameMode=nmpz     MapID=643dbc7ccc47d3a344307998 MapName="An Arbitrary Rural World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_I         GameMode=standard MapID=676340ae2f718dbabdf30331 MapName="An Informed World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_I         GameMode=noMove   MapID=65c86935d327035509fd616f MapName="A Rainbolt World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_I         GameMode=nmpz     MapID=643dbc7ccc47d3a344307998 MapName="An Arbitrary Rural World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_II        GameMode=standard MapID=676340ae2f718dbabdf30331 MapName="An Informed World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_II        GameMode=noMove   MapID=65c86935d327035509fd616f MapName="A Rainbolt World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_II        GameMode=nmpz     MapID=643dbc7ccc47d3a344307998 MapName="An Arbitrary Rural World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_III       GameMode=standard MapID=676340ae2f718dbabdf30331 MapName="An Informed World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_III       GameMode=noMove   MapID=65c86935d327035509fd616f MapName="A Rainbolt World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_III       GameMode=nmpz     MapID=643dbc7ccc47d3a344307998 MapName="An Arbitrary Rural World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_IV        GameMode=nmpz     MapID=643dbc7ccc47d3a344307998 MapName="An Arbitrary Rural World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_IV        GameMode=standard MapID=676340ae2f718dbabdf30331 MapName="An Informed World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Master_IV        GameMode=noMove   MapID=65c86935d327035509fd616f MapName="A Rainbolt World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_I           GameMode=standard MapID=698ef135b77fb917bcdf8425 MapName="A Modified Moving World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_I           GameMode=noMove   MapID=62a44b22040f04bd36e8a914 MapName="A Community World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_II          GameMode=standard MapID=698ef135b77fb917bcdf8425 MapName="A Modified Moving World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_II          GameMode=noMove   MapID=62a44b22040f04bd36e8a914 MapName="A Community World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_III         GameMode=standard MapID=698ef135b77fb917bcdf8425 MapName="A Modified Moving World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_III         GameMode=noMove   MapID=62a44b22040f04bd36e8a914 MapName="A Community World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_IV          GameMode=standard MapID=698ef135b77fb917bcdf8425 MapName="A Modified Moving World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Gold_IV          GameMode=noMove   MapID=62a44b22040f04bd36e8a914 MapName="A Community World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_I         GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_I         GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_II        GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_II        GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_III       GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_III       GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_IV        GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Silver_IV        GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_I         GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_I         GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_II        GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_II        GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_III       GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_III       GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_IV        GameMode=noMove   MapID=6983611e411dbe3f3b2a8c5b MapName="A Figsy World"}
fetcher-1  | DivisionModeMapInfo{DivisionName=Bronze_IV        GameMode=standard MapID=66ca0ac26929fff14a1371b6 MapName="An Extensive Urban World"}
fetcher-1 exited with code 0
```
