# Operations and recovery — Patch 11

Apply after lab Patch 10 and rebuild the C CLI and browser worker. This patch
changes **ukama-lab only**. Use the unchanged uploaded console; ignore console
companions from Patches 4/5. The existing local runtime, BFF teardown credentials,
UI-created networks/sites and node ownership journal remain required.

## Scenarios

All files are under `scenarios/webapp/p0/operations/`.

| File | Acceptance boundary |
|---|---|
| `wb-110-observed-site-restart.yaml` | Arm before submission; observe tower and amplifier go offline and recover while the controller stays visibly online; require cellular Off after reload. |
| `wb-111-switch-port-isolation.yaml` | Exact four port identities and complete state vectors; toggle each off/on separately; require Port 9 Off to survive reload and explicitly restore it. |
| `wb-112-node-status-read-failure.yaml` | Controlled status read error; Restart node disabled with a visible reason; restore reads and recover availability. |
| `wb-113-site-status-read-failure.yaml` | Same for Restart site, Radio, Cellular and every port, including reasons and restoration. |
| `wb-114-optimistic-timeout.yaml` | Controlled idle status reads; busy/disabled through the ten-second optimistic fallback; unlock, restore real reads, submit one explicit retry and recover. |
| `wb-115-distinct-session-recovery.yaml` | Primary confirmation; a distinct saved session starts a controller update; refocus primary, block stale confirmation, cancel/reload and observe completion/version. |
| `wb-116-operation-prerequisites.yaml` | Offline tower disables RF/service; offline controller disables all ports; reconnect restores availability. No operator mutation is submitted. |
| `wb-117-busy-controller-ports.yaml` | A UI-requested controller update locks every port with a reason; completion unlocks every port. |

The original seven journeys remain available. WB013 still tests busy tower,
amplifier and controller locks; WB014 still covers interrupted software updates
and explicit retry. WB110 adds the uninterrupted **rendered** observation that
WB011's post-restart snapshots lacked.

## Local prerequisites

Prepare fresh local trios and capture the primary saved owner session as before.
WB115 needs a second independently authenticated session authorized for the same
organization. Save its Playwright storage state privately and set:

```sh
export ULAB_WEBAPP_PEER_AUTH_STATE=/absolute/private/path/peer.json
export ULAB_SOFTWARE_CURRENT_VERSION=your-installed-version
export ULAB_SOFTWARE_TARGET_VERSION=your-different-published-version
export ULAB_CONTROLLER_RESTART_REASON='expected controller update reason'
```

The last three variables retain their meanings from `operations.md`. They are
independent expectations supplied before execution. Neither software releases
nor sessions are created through hidden APIs by this patch. A different cookie
proves a distinct session; it is not proof of a different user or an agreed role
policy. Primary and peer must see the same test network. WB115's update must
last long enough to observe both refocus and reload while busy.

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
./bin/ukama-lab run scenarios/webapp/p0/operations/wb-110-observed-site-restart.yaml \
  --repo /path/to/ukama --scripts scripts --out runs \
  --bff http://localhost:8080/graphql --factory-url http://localhost:8082
```

Rebuild `bin/ukama-lab` with the repository's usual C build first. Substitute your
existing host paths and service ports. Run each scenario independently with
fresh/reset resources. Controller software updates must not finish before the
busy-port checks. Restarts must expose the tower/amplifier offline interval in
the console's own polling; an unobserved interval fails, it never receives
snapshot-only credit.

## Extended contract

These `web_action` values require the current, already-opened detail identity
and network, exactly like the earlier operations. They do not navigate to repair
wrong scope. `web_field_equals` remains an observation, without mutation/reload.

| Action | View | Fields and meaning |
|---|---|---|
| `open_nodes` | Site | Click the Node component tile. |
| `open_ports` | Site | Click the Switch component tile. |
| `set_port` | Site | `value: "1:off"`, `"1:on"` etc.; only ports 1, 2, 3, 9 and on/off. Click once; already-matching state fails. |
| `watch_restart` | Site | No YAML node list. C resolves exactly one owned tower/amplifier/controller trio; the worker observes its rendered cards. |
| `watch_timeout` | Node | Arm before restart. Observe actual UI time; no clock fast-forward. |
| `status_fault` | Node or site | `value: read_error`, `idle` or `none`. Controlled response injection, described below. |

New fields: `Port states`, `Port 1 reason` (also 2/3/9), `Restart observation`,
`Optimistic timeout` and `Status fault`. New availability labels: `Port 1`,
`Port 2`, `Port 3`, `Port 9`. Missing, hidden or duplicate controls never count
as disabled. Complete vectors compare visible text and checkbox state:

```text
1:Tnode PoE=On;2:Cnode PoE=On;3:Anode PoE=On;9:Uplink SFP=Off
```

`web_tab` retains `primary` and `secondary`. First use of `tab: peer` additionally
requires `auth_state: "${ULAB_WEBAPP_PEER_AUTH_STATE}"`. Later switches use only
`tab: peer`. This creates a separate BrowserContext, validates a different
unexpired `ukama_session` for the console, checks the authenticated dashboard,
and verifies the session remains distinct after navigation. Reusing the primary
session or replacing an existing peer's auth state fails. All pages participate
in bounded cleanup. Failures retain `trace.zip` and `trace-peer.zip`; these are
private raw artifacts, not shareable sanitized reports.

## Observation and fault boundaries

Restart observation arms only when the identified trio is visibly online and
no operation is running. It observes DOM mutations plus 50 ms samples. Controller
non-online state, missing/ambiguous identity, or lost document visibility/path
is sticky failure. Completion requires visible busy, tower and amplifier offline
transitions, then their online recovery and no busy label. Navigation destroys
the watch and cannot silently rearm it. This establishes UI continuity during
the observed restart, **not physical uptime between backend health reports**.

Timeout observation requires visible busy and disabled controls before a later
unlock. An unlock earlier than the source's 10-second fallback, allowing 250 ms
observation tolerance, is sticky failure. A permanently busy UI times out. The
legitimate retry is a separate explicit action after fault removal and reload;
failed assertions and lost transport responses cannot trigger it.

Read faults target only the exact entity and named `GetNodeOperationStatus` or
`GetSiteOperationStatus` query at the gateway endpoint passively observed from
that page. Endpoint discovery reads request metadata only; no response value
becomes an acceptance expectation. An unseen/ambiguous endpoint fails setup.
Cross-origin gateways work; other entities, origins, mutations, batches and
unrelated requests pass through unchanged. `read_error` supplies a GraphQL error;
`idle` supplies an idle status for the optimistic fallback. `none` removes the
route. Reload is an explicit scenario event, never an assertion repair.

Every injected phase first requires `Status fault` to be `read_error:applied`
or `idle:applied`; merely arming a route is insufficient. No browser check uses
BFF data as its oracle. The matrix classifies any `status_fault` scenario as
`controlled_ui`, even if its input manifest falsely calls it live.
Navigation resets application proof; a different detail path passes through
and reports wrong scope. Responses settling from an earlier document cannot
credit the new document.

## Remaining gaps and credit

Three assigned requirements now have complete browser automation: WEB-OPS-004,
WEB-OPS-012 and WEB-OPS-013. Automation totals are **69/132 (52.27%)**, including
**56/86 P0**. All live verification remains **0/132**; the coverage gate is unmet.

- **WEB-SITE-006 / WEB-OPS-008 remain partial.** Original SwitchPortRow initializes
  `enabled=true` and retains optimistic component state, with no reported
  per-port state. WB111 exposes the reload defect. Four UI toggles cannot prove
  hardware port isolation, even when the controlled fixture passes.
- **WEB-OPS-011 remains partial.** WB114 exercises the console's optimistic timer,
  not backend lease expiration or terminal backend timeout. WB014 plus WB114
  still cannot establish that missing backend behavior.
- **Read failures remain a product defect.** Original operation hooks derive
  false busy and available fallbacks after read errors. WB112/113 assert the
  desired fail-closed behavior; the lab does not change React/Apollo policy.
  Their controlled execution cannot earn live qualification credit.

An implementation count describes executable assertions, not passing product
behavior. Original wizard discovery defects can stop setup before these checks.
No console modifications, repaired app state, API operator mutations or live
coverage claims are included. Next is Patch 12: plan, SIM and customer lifecycle.
