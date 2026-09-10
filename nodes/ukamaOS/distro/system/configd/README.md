# config.d — drop 1

Persist configuration decisions and expose their status to lifecycle.d.
This drop changes config.d only. Deploy it with lifecycle drop 2 and backend
drop 3: the previous lifecycle implementation reads configuration progress
from `/v1/ready` and does not understand this new status contract.

## NOCONFIG command

Backend sends this on assignment through the existing node-feeder route,
`POST configd/v1/config`:

```json
{"mode":"NOCONFIG","requestId":"assignment-123"}
```

The request ID identifies the one lifetime assignment decision. Reuse it for
every retry. IDs contain 1–95 ASCII letters, digits, `.`, `_`, `-`, or `:`.
The command accepts exactly these two fields. No assignment removal or
reassignment is implemented.

| HTTP | Meaning |
| --- | --- |
| 200 | Decision and completion are durable; body is the current status snapshot |
| 400 | Invalid command or request identity |
| 409 | Different NOCONFIG identity, active upload, or an existing CONFIG decision |
| 503 | Persistent state cannot safely be read/written |

200 acknowledges config.d completion, not node Operational. Backend drop 3
must wait for the lifecycle OPERATIONAL event and retry the same command every
60 seconds until confirmed. A concurrent request can advance the response
snapshot; the retained GET status is authoritative.

## Lifecycle interface

`GET /v1/config/status` returns HTTP 200 with structured facts, including when
a configuration transaction has failed:

```json
{
  "schemaVersion": 1,
  "mode": "NOCONFIG",
  "phase": "completed",
  "requestId": "assignment-123",
  "generation": 1,
  "revision": 0,
  "error": ""
}
```

| Field | Values/meaning |
| --- | --- |
| mode | `NONE`, `NOCONFIG`, `CONFIG` |
| phase | `awaiting`, `pending`, `completed`, `failed` |
| requestId | Stable decision identity; empty before assignment |
| generation | Durable counter; increases for every accepted NOCONFIG delivery or new CONFIG session |
| revision | Existing CONFIG timestamp; zero for NOCONFIG |
| error | Recovery/storage error, or empty; phase `failed` alone also indicates failure |

A missing record returns NONE/awaiting with generation 0. A corrupt or
unreadable record returns failed and is not silently replaced. State files
are bound to the node ID.

Lifecycle drop 2 polls this endpoint. There is no new config.d-to-lifecycle
HTTP command and no backend call to lifecycle's old `/v1/configure` endpoint.
The retained decision means lifecycle cannot miss an immediate completion:

- On node startup, consume the restored/current record through the same
  READY → CONFIGURING → OPERATIONAL reducer path.
- A completed record must still pass through CONFIGURING logically; lifecycle
  must not require observing the brief pending phase in a separate poll.
- While already Operational, a higher generation of the same completed decision
  requests a fresh OPERATIONAL confirmation after checking current app readiness.
  It must not force a backward state transition.
- Polling an unchanged generation does not repeatedly emit OPERATIONAL.
- A config.d-only restart retains the generation. A new node boot must consume
  the record again even when its generation is unchanged. Lifecycle-only crash
  recovery must reconcile its own checkpoint with this retained record.
- Lifecycle owns readiness gating and all node state/events. Config.d does
  not call notify.d or publish OPERATIONAL itself.

`GET /v1/ready` now reports daemon availability only:
`{"ready":true,"reason":"ready"}`. Configuration failure is visible through
`/v1/config/status`; it is not encoded as daemon unavailability.

## Persistence and recovery

Default: `/ukama/configs/configd/state/configuration.status`.
Override with `--state-file /absolute/persistent/path/configuration.status`.
The directory must survive daemon/container replacement and node reboot.
It is outside the application's `active`/`archive` configuration directories.

The small versioned text file contains node identity, mode, phase, generation,
request ID, and CONFIG revision/application names. Treat it as daemon-owned;
use the HTTP status endpoint to inspect it.

Each update writes a same-directory temporary file, flushes and fsyncs it,
renames it atomically, and fsyncs the parent directory before exposing success.
A lifetime advisory lock prevents two config.d processes sharing one state file.
There is no database or background worker.

First NOCONFIG persists pending intent, then successful completion. On restart,
a pending NOCONFIG can safely complete because it has no app side effects.
Repeated completed NOCONFIG only persists a new confirmation generation.
A transient persistence failure reports failure; a later command reloads and
reconciles the durable file before retrying. It never guesses whether a failed
rename succeeded.

Corrupt or mismatched-node state requires operator diagnosis/restoration; the
backend retry cannot overwrite it. Interrupted CONFIG reports failure and
requires explicit CONFIG retry. It cannot be converted to successful NOCONFIG.

## Existing CONFIG uploads

The existing file-envelope route remains `POST /v1/config`, without a `mode`
field. File payload structure and starter restart calls are preserved.

The integration persists intent before staging, tracks affected apps, and
persists terminal status after application. Calls are serialized against each
other and NOCONFIG. Old envelopes without requestId use `config-<timestamp>`
in the new status interface. Supplied IDs follow the same identifier rules.

Before recording completed CONFIG, config.d verifies the active revision and
flushes the configuration filesystem. On startup it verifies that each recorded
app's active path still resolves to that revision's archive directory. This
checks revision presence, not a content hash. App readiness remains lifecycle's
responsibility; restoration does not restart apps again.

This is persistence integration, not a rewrite of CONFIG uploads. Existing
file-count/duplicate handling, upload expiry, commit metadata handling, and
multi-app rollback behavior remain outside this drop. Existing HTTP 201 upload
responses are not a guarantee of complete configuration. No existing active
configuration is inferred to be an assignment when the new status file is absent.

## Build and tests

Build in the existing Ukama tree with `make`. New source files are picked up
by the existing wildcard; pthread flags are explicit.

Storage tests require only a C compiler, Linux/libc and pthreads:

```sh
make -C test unit
```

They cover assignment, concurrent duplicate delivery, exclusive writer locking,
process crash, each fsync failure boundary, transient-error retry, corrupt state,
node identity, interrupted CONFIG and active-revision verification.

Test a built daemon with real HTTP requests:

```sh
CONFIGD_BIN=/absolute/path/config.d CONFIGD_PORT=8080 make -C test component
```

Set CONFIGD_PORT to the actual service-map port. Tests launch the daemon with
debug node identity and an isolated temporary state file; they kill/restart
that test process. They do not require the backend or lifecycle.

For local testing without libusys, `make -C test component-build` builds all
production sources against real Ulfius/Jansson using test-only platform adapters
under `test/platform`. This is not a production binary. These adapters choose
port 18080 by default. Production repository build and on-node integration
remain separate gates.

## Changed-file manifest

Existing files changed:

- `Makefile`
- `README.md`
- `inc/config.h`
- `inc/configd.h`
- `inc/web_service.h`
- `src/configd.c`
- `src/main.c`
- `src/network.c`
- `src/web_service.c`

New implementation files:

- `inc/state_store.h`
- `src/state_store.c`

New tests: `test/Makefile`, `test/test_state_store.c`,
`test/component_test.py`, and the test-only platform adapters in `test/platform/`.

No files in lifecycle.d, backend, starter.d, notify.d, or lookout.d are changed.
