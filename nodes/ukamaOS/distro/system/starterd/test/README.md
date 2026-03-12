# starter.d smoke tests

This bundle contains self-contained smoke tests for:

1. `starter.d` app lifecycle with a mock WIMC server and a mock `example_app`
2. `init-starter` slot switch behavior on child exit code `77`

## Files

- `smoke_example_app.sh`
- `smoke_init_switch.sh`
- `run_all.sh`
- `mock_wimc.py`
- `common.sh`

## Requirements

These scripts assume the following are available on the machine where you run them:

- `bash`
- `python3`
- `curl`
- `tar`
- `mktemp`

And the binaries under test:

- `starter.d`
- `init-starter`

## Usage

Run the example app smoke test:

```bash
STARTERD_BIN=/path/to/starter.d ./smoke_example_app.sh
```

Run the init-starter slot-switch smoke test:

```bash
INIT_STARTER_BIN=/path/to/init-starter ./smoke_init_switch.sh
```

Run both:

```bash
STARTERD_BIN=/path/to/starter.d \
INIT_STARTER_BIN=/path/to/init-starter \
./run_all.sh
```

## What `smoke_example_app.sh` does

- creates a temp root
- creates two mock app packages: `example_app-v1.tar.gz` and `example_app-v2.tar.gz`
- starts a mock WIMC HTTP server that serves `/v1/apps/<app>/<tag>/pkg`
- writes a manifest that loads `example_app` in the `boot` space
- starts `starter.d` with temp env vars
- validates:
  - `starter.d` responds on `/v1/ping`
  - `example_app` boots and reports version `v1`
  - `/v1/terminate` stops the app
  - `/v1/update` updates the app to `v2` and restarts it

## What `smoke_init_switch.sh` does

- creates a temp `init-starter` root with `slots/A` and `slots/B`
- creates mock `starter.d` slot binaries
- slot `A` touches the ready file, then exits `77`
- slot `B` touches the ready file and runs briefly
- starts `init-starter`
- validates that:
  - `init-starter` flips `current` from `slots/A` to `slots/B`
  - `prev` is set to `slots/A`
  - the new slot is executed

## Notes

- The `init-starter` smoke test validates the `77` handoff behavior directly.
- It does **not** invoke the real `starter.d` self-update endpoint to populate `next`; it isolates and validates the bootstrapping/switch logic itself.
- If your `starter.d` build still requires service-registry lookup before honoring `STARTERD_HTTP_PORT` / `STARTERD_WIMC_PORT`, make sure your build can start with the env overrides used here.
