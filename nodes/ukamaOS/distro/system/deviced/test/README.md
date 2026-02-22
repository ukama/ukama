# device.d tests (Unity)

These tests run **device.d as-is** and mock only its external services
    (**node.d**, **notify.d**, **fem.d**) by spinning up lightweight Ulfius HTTP
    servers on the **same ports device.d expects** (derived from `usys_find_service_port`).

Because `actions_tower.c` invokes absolute-path shell commands and `actions_restart.c`
    may call `reboot()`, the test runner launches device.d with a small
    `LD_PRELOAD` shim (`libdeviced_testwrap.so`) to make those operations safe and
    deterministic during tests.

## Run

```bash
make -C test UNITY_DIR=/path/to/Unity run
```

Where `UNITY_DIR` contains `src/Unity.c` and `src/Unity.h`.

Notes:
- The tests rely on your usual Ukama build environment (Ulfius/Jansson/CURL/usys available,
    and `usys_services` configured so `usys_find_service_port(...)` works).
- The runner tries to add `loclahost` -> `127.0.0.1` to `/etc/hosts` if writable
    (workaround for a host typo in this snapshot).
