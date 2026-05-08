# Site Controller HTTP API

The first implementation exposes a small HTTP API. API Gateway can wrap this later.

## Health

```http
GET /v1/ping
GET /v1/version
```

## Site state

```http
GET /v1/sites/{site_id}/state
GET /v1/sites/{site_id}
```

Returns desired intent, derived state, component snapshot, and static port map.

## Static port map

```http
PUT /v1/sites/{site_id}/ports
GET /v1/sites/{site_id}/ports
```

Request:

```json
{
  "ports": [
    {"port":1,"role":"tower","node_id":"emu-tnode","cnode_id":"emu-cnode","class":"critical","policy":"protected"},
    {"port":2,"role":"cnode","node_id":"emu-cnode","cnode_id":"emu-cnode","class":"critical","policy":"never_off_remote"},
    {"port":3,"role":"amplifier","node_id":"emu-anode","cnode_id":"emu-cnode","class":"critical","policy":"protected"}
  ]
}
```

## Apply switch policy

```http
PUT /v1/sites/{site_id}/switch-policy
POST /v1/sites/{site_id}/switch-policy
```

Generates switch.d policy JSON from static port map and sends it to CNode `switch.d`.

## Site/service/radio commands

```http
POST /v1/sites/{site_id}/on
POST /v1/sites/{site_id}/off
POST /v1/sites/{site_id}/service/on
POST /v1/sites/{site_id}/service/off
POST /v1/sites/{site_id}/radio/on
POST /v1/sites/{site_id}/radio/off
```

Optional body:

```json
{"reason":"maintenance","requestedBy":"operator"}
```

## Power cycle

```http
POST /v1/sites/{site_id}/nodes/{role}/power-cycle
```

Example:

```http
POST /v1/sites/emu-site/nodes/amplifier/power-cycle
```

CNode power-cycle is rejected because the CNode controls the switch and is powered by PoE.
