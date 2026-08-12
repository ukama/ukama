# lifecycle.d tests

## Unit test

The FSM and state store tests use no Ukama platform or vendor libraries:

```bash
make unit
```

They validate:

- STARTING and CHECKING_IN timing.
- Stable READY behavior.
- No-config timeout to OPERATIONAL.
- Configuration progress, success and failure.
- Configure command idempotency.
- Starter fault and recovery.
- Starter-unavailable timeout.
- Boot-scoped state persistence.

## Component test

`component_test.py` starts lightweight fake starter and notify HTTP services,
then runs the real lifecycle daemon as a child process.

```bash
make component LIFECYCLED_BIN=../lifecycle.d
```

It verifies the complete event order without requiring a backend, starter,
notify, configuration service, or a 60-second test delay.
