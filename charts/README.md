# Install chart

First install will take some time.

```bash
helm dependency update
helm -n app-ggwd upgrade --install ggwd ggwd --values values.yaml
```

```bash
kubectl create namespace test-ggwd
helm -n test-ggwd upgrade --install ggwd ggwd --values values_test.yaml
kubectl delete namespace test-ggwd
```

minimum `values.yaml`:

```yaml
---
fullnameOverride: ggwd

worker:
  #command: ["sleep", "infinity"]
  config:
    shoutrrrTargets:
      - 'discord://...@...'
    watchdogs:
      CompetitiveMaps:
        enabled: true
        notifyAbout:
          - divisionName: "Master I"
            gameMode: "standardDuels"  # "standardDuels", "noMoveDuels", "nmpzDuels"
      UserStats:
        enabled: true
        observeUsers:
          - ...
    cookie: '...'

cacheBrowser:
  enabled: true
  ingress:
    enabled: true
    host: ggwd-cache.local

temporal:
  enabled: true
  web:
    enabled: true
    ingress:
      enabled: true
      hosts:
        - ggwd-temporal.local

gluetun:
  enabled: true
  secretEnv:
    WIREGUARD_PRIVATE_KEY: "..."
    VPN_SERVICE_PROVIDER: "nordvpn"
    VPN_TYPE: "wireguard"
    SERVER_COUNTRIES: "Poland"
```
