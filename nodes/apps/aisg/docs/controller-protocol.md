# Controller protocol

The controller contract is newline-delimited JSON over a UNIX domain socket.

Default socket:

```text
/var/run/aisg-ctrl.sock
```

Request:

```json
{
  "id": "req-001",
  "type": "get_tilt",
  "payload": {}
}
```

Response:

```json
{
  "id": "req-001",
  "ok": true,
  "code": "OK",
  "reason": "",
  "payload": {}
}
```

V1 message types:

- `ping`
- `get_status`
- `scan`
- `get_info`
- `get_alarm_status`
- `clear_active_alarms`
- `alarm_subscribe`
- `self_test`
- `send_configuration_data`
- `calibrate`
- `get_tilt`
- `set_tilt`
- `get_device_data`
- `reset_software`
