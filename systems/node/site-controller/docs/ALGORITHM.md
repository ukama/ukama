# Site Controller Algorithm

## Ownership

- Registry owns site/node topology.
- site-controller owns site intent, static switch port map, generated switch policy, and derived site state.
- node-controller owns internal node command dispatch over gRPC.
- switch.d owns local switch observation and local policy enforcement.
- Node-GW remains node-facing for health/state/event reports.

## First boot

1. CNode boots and starts switch.d.
2. switch.d loads `/ukama/configs/switch/policy.json` if present.
3. If the policy is missing or invalid, switch.d remains observe-only: reads/report ports, blocks destructive writes.
4. site-controller starts or receives CNode health later.
5. site-controller reads its static site port map.
6. site-controller generates switch.d policy JSON.
7. site-controller calls node-controller over gRPC to send `PUT /switch/v1/ports/policy` to CNode switch.d.
8. switch.d validates, atomically stores, and enforces the policy.
9. site-controller stores/derives site state and emits site-level events.

## Site on

1. Store desired site/service/radio = on.
2. Apply switch policy.
3. Turn required protected ports on idempotently.
4. Turn amplifier radio on.
5. Turn tower service on.
6. Derive access as available only when power/service/radio are healthy/on/running.

## Site off

1. Store desired site/service/radio = off.
2. Apply switch policy if possible.
3. Turn radio off.
4. Turn service off.
5. Do not turn PoE off.

## Radio off

1. Store desired radio = off.
2. Turn amplifier radio off.
3. Leave tower service running.
4. Leave PoE on.
5. Derive access unavailable with reason `radio_off`.

## Service off

1. Store desired service = off.
2. Turn tower service off.
3. Leave radio and PoE unchanged.
4. Derive access unavailable with reason `service_off`.

## Power-cycle

1. Look up target role in static port map.
2. Send `POST /switch/v1/ports/{port}/poe/cycle` through node-controller to CNode switch.d (site-controller does not reject by role; **switch.d** on the node enforces policy, e.g. `never_off_remote` for the CNode port).
3. Keep desired site/service/radio unchanged and reconcile after node health returns.

## Edge cases

- Missing port map: do not control switch, mark degraded.
- Invalid port map: do not push policy.
- CNode unreachable: mark control degraded; running access may continue but control is degraded.
- switch.d policy missing: push policy again.
- Wrong static map: operator must correct it; discovery will be added later.
- CNode port: operator maps it as `never_off_remote`; **switch.d** rejects destructive actions for that port at the node.
