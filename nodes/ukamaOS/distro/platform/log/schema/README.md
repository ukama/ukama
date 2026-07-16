# Ukama structured logging

This directory defines the versioned record contract emitted by `usys_log`.
The contract is intentionally small. Applications emit a common envelope and
add event-specific fields only when they are useful.

## Record format

Each record is one UTF-8 JSON object followed by one newline. There is no
outer JSON array.

Required fields are defined in `log-v1.schema.json`. Common field names and
initial events are listed in `fields.yaml` and `events.yaml`.

## Ownership

- Applications own `component`, `event`, `msg`, and event-specific fields.
- `usys_log` owns serialization, timestamps, process metadata, and source
  location.
- `starterd` will later add trusted application and capture metadata.
- `rlog.d` will later add node, boot, and canonical sequence metadata.

## Naming

Event and field names use lowercase snake case. Application names do not
belong in event names. For example, use `dependency_connect_failed` with
`dependency=pcrf`, not `pcrf_connect_failed`.

## Sensitive data

Never log credentials, tokens, authorization headers, private keys, Ki, OPc,
or authentication vectors. Use a stable redacted reference for subscriber
identifiers in normal operation.
