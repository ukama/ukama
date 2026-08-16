Ukama light-weight service mesh and reverse/forward proxy control and data
plane.

The local service API exposes two separate health signals:

- `GET /v1/ready` returns `200 {"ready":true}` when mesh.d is serving its
  local API. This endpoint participates in node readiness.
- `GET /v1/status` reports the remote websocket connection as `connected`,
  `reason`, and `changedAt`. Remote connectivity does not block node
  readiness.
