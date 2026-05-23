# init-network.d

Starter-managed virtual tower-node network initializer.

It prepares OVS/br0 before PCRF and EPC-emu (in virtual node, EPC in real) starts:

- checks required tools
- starts `ovsdb-server` and `ovs-vswitchd` when needed
- creates/configures `br0`
- configures PCRF OpenFlow management socket
- enables IPv4 forwarding and optional NAT
- installs default UE-subnet drop flows
- exposes `/v1/ping`, `/v1/version`, and `/v1/status`

*** This app is intended for tower-node virtual images only. ***

## Endpoints

```text
GET /v1/ping
GET /v1/version
GET /v1/status
```

`/v1/ping` returns `200` only after OVS setup is ready. It returns `503` while setup is incomplete or failed.

## Manifest entry

Place first in tower-node `boot` space:

```json
{
  "name": "init-network",
  "tag": "latest",
  "cmd": "sbin/init-network.d",
  "argv": [
    "init-network.d",
    "--config",
    "/ukama/configs/init-network/config.toml"
  ]
}
```

## TODO

- Add EPC-emu user-plane interface wiring into `br0`.
- Add media/forwarder interface wiring into `br0`.
- Add policy routing tables after finalizing the EPC-emu traffic contract.
- Add stale PCRF flow cleanup/reconciliation once PCRF cookie ownership is explicit.
