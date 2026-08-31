Client for init system

initClient is:

1. Registration of the system to the local Init and global Init using <name:address:port>. Rcv: UUID
   $(INIT_SYSTEM): init.ukama.com
2. Query init for specific system info
3. Send periodic health update, restart, update to init system (using UUID)
4. De-register itself.
5. handle GRPC from services within the System.

Work required:
1. Need to add support for maintaining URL's too in init.

## Registration ownership

Init stores one record per `{org, system}`, so several instances of the same
system's initClient address the same record over its lifetime. Each instance
therefore tracks which registration it owns:

* On start-up the instance registers unless init already holds a record that
  matches its configuration *and* carries the registration id cached in
  `ENV_INIT_CLIENT_TEMP_FILE`. Registering returns a new id, which becomes this
  instance's proof of ownership.
* On shutdown the instance de-registers only when init still holds the id it
  registered. An instance that has been replaced finds a different id and leaves
  the record alone, so a rolling restart cannot remove the registration of the
  instance that took over from it.
* Every `ENV_INIT_RECONCILE_PERIOD` seconds (default 30) the instance re-checks
  the record. A missing record is registered again and a drifted one is
  refreshed, so a registration lost for any reason is restored within one
  period rather than staying gone until the container is replaced. A record
  owned by another instance is left untouched.

## Environment variables

| Variable | Required | Notes |
| --- | --- | --- |
| `ENV_SYSTEM_ORG` | yes | Organization the system belongs to |
| `ENV_SYSTEM_NAME` | yes | System name, the key of the init record |
| `ENV_SYSTEM_PORT` | yes | api-gw port |
| `ENV_SYSTEM_CERT` | yes | Certificate registered for the system |
| `ENV_SYSTEM_DNS` / `ENV_SYSTEM_ADDR` | one of | api-gw name to resolve, or a literal address |
| `ENV_SYSTEM_API_GW_URL` | no | Stable api-gw URL; preferred by consumers over the raw IP |
| `ENV_SYSTEM_DNS_NODE_GW` / `ENV_SYSTEM_NODE_GW_ADDR` | no | node-gw name to resolve, or a literal address |
| `ENV_SYSTEM_NODE_GW_PORT` | with node-gw | Required when any node-gw variable is set |
| `ENV_INIT_SYSTEM_ADDR` / `ENV_INIT_SYSTEM_PORT` | yes | Local init api-gw |
| `ENV_GLOBAL_INIT_ENABLE` | no | `true` also registers with the global init |
| `ENV_GLOBAL_INIT_SYSTEM_ADDR` / `_PORT` | with global | Required when global init is enabled |
| `ENV_INIT_CLIENT_ADDR` / `_PORT` | yes | Bind address and port for the lookup endpoint |
| `ENV_INIT_CLIENT_TEMP_FILE` | yes | Where the registration ids are cached |
| `ENV_DNS_REFRESH_TIME_PERIOD` | no | Seconds between address re-resolutions (default 10) |
| `ENV_INIT_RECONCILE_PERIOD` | no | Seconds between reconcile passes (default 30) |
| `ENV_DNS_SERVER` | no | `true` reads the nameserver from /etc/resolv.conf |
| `ENV_INIT_CLIENT_LOG_LEVEL` | no | `DEBUG`, `INFO` or `ERROR` |
