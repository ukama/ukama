# lifecycle.d

`lifecycle.d` is the sole owner of the node lifecycle state. It converts
starter readiness, backend commands, and local timers into node state
transitions.

It does not start applications, apply configuration, control the radio, or
send periodic health reports.

## V1 states

```text
STARTING -> CHECKING_IN -> READY -> CONFIGURING -> OPERATIONAL
                                  \-> FAULTY
```

`READY` is stable. The node remains READY until an explicit CONFIGURING
command is received.

`UNKNOWN` and online/offline availability are backend concepts. They are not
local lifecycle states.

## Inputs

- `POST /v1/check-in` from `starter.d` after boot applications start.
- `GET starter.d/v1/status` for aggregate and per-app readiness.
- `POST /v1/configure` from the node's backend command path.
- Monotonic local timers.

`lifecycle.d` never polls individual managed applications.

## Output

Every state transition is posted to:

```text
POST notify.d/v1/event/lifecycle
```

`starter.d` must stop producing node state events when Patch 2 activates this
daemon as the node-state authority.

## HTTP API

```text
GET  /v1/ping
GET  /v1/version
GET  /v1/ready
GET  /v1/status
POST /v1/check-in
GET  /v1/gate
POST /v1/configure
```

`/v1/ready` reports whether the lifecycle daemon itself is healthy. It does
not mean the node lifecycle state is READY.

Start check-in:

```json
{
  "bootId": "optional-current-kernel-boot-id",
  "bootResult": "ready"
}
```

Start configuration:

```json
{
  "requestId": "required-idempotency-id",
  "assignmentId": "optional-backend-assignment-id"
}
```

## Configuration

Environment variables:

```text
LIFECYCLED_HTTP_ADDR
LIFECYCLED_HTTP_PORT
LIFECYCLED_STARTER_HOST
LIFECYCLED_STARTER_PORT
LIFECYCLED_NOTIFY_HOST
LIFECYCLED_NOTIFY_PORT
LIFECYCLED_STATE_FILE
LIFECYCLED_CHECKIN_TIMEOUT_SEC
LIFECYCLED_CONFIG_TIMEOUT_SEC
LIFECYCLED_STARTER_UNAVAILABLE_TIMEOUT_SEC
LIFECYCLED_POLL_INTERVAL_MS
LIFECYCLED_REQUEST_TIMEOUT_SEC
LIFECYCLED_LOG_LEVEL
```

Production check-in and no-configuration defaults are both 60 seconds. Tests
override them with one-second values.

## Tests

The pure FSM unit tests need only a C compiler:

```bash
make -C test unit
```

The component test starts the real daemon with fake starter and notify HTTP
services:

```bash
make
make -C test component LIFECYCLED_BIN=../lifecycle.d
```
