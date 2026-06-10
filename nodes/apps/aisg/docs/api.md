# aisgd API

## Lifecycle

- `GET /v1/ping`: plain `OK`, HTTP 200.
- `GET /v1/version`: plain version string, HTTP 200.
- `GET /v1/status`: JSON status.
- `POST /v1/reconcile`: reconcile controller/device state.

## Device

- `POST /v1/device/scan`
- `GET /v1/device`
- `GET /v1/device/info`
- `GET /v1/device/alarms`
- `POST /v1/device/alarms/clear`
- `POST /v1/device/alarms/subscribe`
- `POST /v1/device/self-test`
- `POST /v1/device/config`
- `POST /v1/device/calibrate`
- `GET /v1/device/tilt`
- `POST /v1/device/tilt`
- `GET /v1/device/data/:field`
- `POST /v1/device/reset`
