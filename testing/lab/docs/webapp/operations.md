# Browser operations and recovery — Patch 5

Apply the lab patch after Patches 1–4 and the console companion after its Patch 4
companion. Rebuild both. The existing host setup, owner authentication, local
virtual nodes, endpoint configuration and owned-resource cleanup are described
in [provisioning.md](provisioning.md). No AWS execution is introduced.

## Scenarios

All files below are in `scenarios/webapp/p0/operations/`. Every scenario creates
its own network/site through the UI and uses a local tower/amplifier/controller
trio. Operator actions and acceptance observations go through the app. Host
runtime controls inject disconnection/reconnection; backend teardown is still
limited to resources owned by that run.

| File | Browser journey |
|---|---|
| `wb-010-node-restart.yaml` | Online/Operational → cancel → confirm once → visible progress and disabled control → release → reload and recovery. |
| `wb-011-site-operations.yaml` | Missing/wrong/exact site name, cancellation, RF and service off/on with reported state, site restart and recovery snapshots for each node. |
| `wb-012-software-update.yaml` | Installed/target version and available state → update → updating/locked → up to date with exact target installed. |
| `wb-013-node-and-site-locks.yaml` | Update controller, tower and amplifier separately. Each blocks site restart; the controller also blocks RF/service with expected reasons. Observe unlock after each update. |
| `wb-014-software-failure-retry.yaml` | Disconnect during updating → Offline with Operational lifecycle → reconnect → Update failed with unchanged installed version → explicit retry → exact installed target. |
| `wb-015-stale-confirmation.yaml` | Open confirmation in primary tab, update through secondary tab, refocus primary and require disabled confirmation, cancel/reload and recover progress/completion. |
| `wb-016-offline-recovery.yaml` | Offline restart restriction and reason, lifecycle preservation, reconnect and availability. |

The software scenarios require the `example` application to start at
`ULAB_SOFTWARE_CURRENT_VERSION`, with a different published desired release
`ULAB_SOFTWARE_TARGET_VERSION` already available to the runtime. Prepare that
release as a host fixture before the browser scenario. This app has no release
promotion screen; the worker never promotes a release through an API. Each
software scenario needs fresh nodes/reset installed versions; do not run the
suite repeatedly against a reused already-updated trio.

For wb-013, set `ULAB_CONTROLLER_RESTART_REASON`, `ULAB_CONTROLLER_RF_REASON`
and `ULAB_CONTROLLER_SERVICE_REASON` to the expected case-sensitive visible
reason text, or a stable meaningful substring specified by the tested backend
build. Include the controller/operation identity where its contract supplies it.
These are test expectations supplied before execution, never derived from the
page during a check. The worker preserves the full actual reason in evidence.
The placeholders in automated fixture tests are not a claim about the live BFF's
wording.

Updates/restarts must last long enough to observe their in-progress states.
For wb-014, the local update fixture must produce a terminal failed update when
interrupted. An update that completes before disconnection causes this scenario
to fail; a success state is not silently accepted as a failure/retry test.

After configuring the endpoint/auth/runtime values in `provisioning.md` and
setting the software/reason expectations above:

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
./bin/ukama-lab lint scenarios/webapp
./bin/ukama-lab run scenarios/webapp/p0/operations/wb-010-node-restart.yaml \
  --repo /path/to/ukama --scripts scripts --out runs \
  --bff http://localhost:8080/graphql --factory-url http://localhost:8082
```

Replace the scenario filename to run another journey. Use your configured ports.

## Writing operation steps

`web_action` always names the current detail view, one network, and its matching
node or site reference. It never navigates to repair a stale or wrong scope.
A preceding `web_open` must have selected that entity. C resolves identities
and checks network membership; the worker checks the visible selection and URL.

| Action | View | Extra fields |
|---|---|---|
| `open_restart`, `confirm_restart`, `cancel_dialog` | Node or site detail | None. Site `open_restart` requires an explicitly opened Site actions menu. |
| `open_site_actions`, `close_site_actions` | Site detail | None. |
| `fill_confirmation` | Site detail | Exactly one of `value` (including empty) or `value_from: site_name`. The latter uses the planned world name, not page text. |
| `set_radio`, `set_service` | Site detail | `value: "on"` or `"off"`; menu already open. An already-matching state fails rather than crediting a mutation that never occurred. |
| `open_software` | Node detail | None. Clicks the visible Software tab. |
| `update_software`, `retry_update` | Node detail | `app` and `tag`. The visible target must equal `tag` before clicking. Retry additionally requires the rendered Update failed state and Retry update control. |

```yaml
- name: start_update
  events:
    - type: web_action
      view: network_node_detail
      networks: net-001
      nodes: controller-site-001-001
      action: update_software
      app: example
      tag: ${ULAB_SOFTWARE_TARGET_VERSION}
  checks:
    - type: web_field_equals
      view: network_node_detail
      label: Software status
      app: example
      expected: Updating…
      requirement: WEB-OPS-009
    - type: web_action_available
      view: network_node_detail
      label: Restart node
      available: false
      requirement: WEB-OPS-009
```

`web_tab` takes only `tab: primary|secondary` and optional `timeout_seconds`.
The secondary page is created once in the authenticated context. Each page
retains its own navigation scope and dialogs; a new page requires `web_open`.
Switching uses browser focus, without injected visibility events. This covers
two tabs of one authenticated account; distinct users/roles remain separate work.

## Visible assertions and evidence

Existing checks add these semantic operation labels:

- Availability: Restart node, Site actions, Restart site, Confirm restart,
  Cancel, Radio, Cellular, Update Now and Retry update. Software buttons require
  `app`. Hidden/missing controls never satisfy `available: false`.
- Text: Dialog (`open`/`closed`), Confirmation, Restart progress, Restart reason,
  Radio state, Cellular state, Radio reason, Cellular reason, Lifecycle and
  Notification. Software status, Current version and Target version require
  `app`. A missing target/status is not an empty or successful value.
- `match: contains` is allowed only on nonempty literal reason checks. Other
  text, versions, units and counts remain exact after whitespace normalization.
  C independently verifies the match; reports retain `match`, `app` when used,
  expected text and the complete actual text. Restart progress removes only
  the button's parenthesized elapsed timer.

Actions click once. Retrying a command ID replays its recorded result, not its
mutation. Observation polling cannot submit, navigate, reload, dismiss a dialog
or inject responses. A failed action/assertion terminates browser work and
preserves screenshot/trace/diagnostics; explicit `retry_update` is a new step
following an observed product failure, not a worker transport retry.

## Console companion and limits

The console now treats pending/failed lock reads as unverified, disables affected
controls, retries reads while visible, and refreshes on focus. Restart and
software share the node lock. Node detail polls every 30 seconds while visible,
or every 3 seconds while busy/on Software, and refetches when the lock clears.
RF/service show Changing… until health reports confirm the requested state;
unknown/error reads show Unknown. A success acknowledgement alone does not set
On/Off. Pending confirmation expires after 90 seconds with a visible notice.
Mutation `success: false` is treated as a rejection.

The following inventory entries retain **no full automation credit**:

- WEB-OPS-004: wb-011 checks recovery snapshots, not uninterrupted controller
  uptime throughout site restart. Continuous UI observation is still needed.
- WEB-OPS-008: switch-port controls have only local optimistic state; the app
  must expose reported per-port state before UI-only isolation can be proved.
- WEB-OPS-011: wb-014 covers interrupted-update failure and explicit retry;
  deterministic lease-expiry/timeout coverage is still needed.
- WEB-OPS-013: the fail-closed policy has unit tests and the browser adapter has
  a controlled error-state test. Real React/Apollo error/recovery integration
  against a configured fault profile still needs qualification.

The inventory also retains distinct-session/user coverage as a gap on
WEB-OPS-012: wb-015 currently covers two same-account tabs. No fixture result
counts as verified product coverage. See [patch-05-validation.md](patch-05-validation.md).
