# lifecycle.d tests

See ../README.md for commands and contracts.

- test_fsm.c: indefinite READY; completed/pending/failed decisions; retry
  generation; reboot; stale state; startup gate recovery; checkpoint queue.
- component_test.py: real lifecycle HTTP service with starter/config/notify
  peers; notification outage; SIGKILL; readiness-gated repeats; new boot ID.
- configd_integration_test.py: the actual config.d drop 1 and lifecycle drop 2
  executables, with starter/notify peers.
- platform/: local test adapters for libusys. HTTP, JSON and persistence use
  the actual libraries/implementations. Never install these adapters on nodes.

The old stubs/ headers are retained from the supplied tree but are not used by
these targets; syntax checks now compile against the actual dependency headers.
