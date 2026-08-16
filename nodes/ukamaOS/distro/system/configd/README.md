# config.d

Config service running on the ukama node to configure apps.

## Lifecycle readiness

`GET /v1/ready` reports both service readiness and configuration progress:

| Phase | HTTP | `ready` | `reason` |
| --- | ---: | --- | --- |
| Waiting for configuration | 200 | `true` | `awaiting_configuration` |
| Applying configuration | 202 | `false` | `configuration_in_progress` |
| Configuration applied | 200 | `true` | `configuration_applied` |
| Configuration failed | 503 | `false` | `configuration_failed` |

When an incoming configuration includes `requestId`, the same value is
returned by `/v1/ready`. This lets `lifecycle.d` correlate configuration
progress with its active `/v1/configure` command.

After activating configuration symlinks, `config.d` calls the atomic
`starter.d` restart endpoint with the target application:

```json
{
  "space": "services",
  "name": "app-name"
}
```

The existing `POST /v1/config` payload remains compatible. `requestId` is an
optional lifecycle-correlation field.

## Build

```
make
```

## Test
