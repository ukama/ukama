# epcemu

EPC emulator for virtual tower-node PCRF testing.

This app is starter-managed and control-plane only. It does not generate or
forward user-plane traffic yet.

## Responsibilities

- Resolve local service ports from `/etc/services`
- Read `init-network` status and UE CIDR
- Check PCRF readiness
- Accept clean UE attach/detach requests
- Translate UE attach/detach into PCRF `/v1/session` calls
- Keep UE state in memory
- Best-effort detach all attached UEs on shutdown

## Endpoints

```text
GET    /v1/ping
GET    /v1/version
GET    /v1/status
POST   /v1/ue/attach
DELETE /v1/ue/detach
GET    /v1/ue/:imsi
GET    /v1/ues
```

## /etc/services

```text
epcemu        18092/tcp
pcrf          18090/tcp
init-network 18026/tcp
```

## Attach

```bash
curl -X POST localhost:18092/v1/ue/attach \
  -H 'Content-Type: application/json' \
  -d '{"imsi":"001010000000001","ip":"192.168.8.2","apn":"internet"}'
```

## Detach

```bash
curl -X DELETE localhost:18092/v1/ue/detach \
  -H 'Content-Type: application/json' \
  -d '{"imsi":"001010000000001"}'
```

## TODO

- Add UE container runner
- Add media container
- Add EPC/user-plane interface wiring into OVS/br0
- Add full UE traffic path through PCRF-controlled OVS flows
