# Patch 9: staged onboarding

This patch changes **ukama-lab only**, after lab Patches 1–8. Use the original
console source; ignore the console companion archives from Patches 4 and 5.
The new mode builds and readies one owned virtual node trio before running the
wizard. YAML then drives each browser step instead of hiding network/site
creation in setup. No console changes, API-created acceptance resources, login
bypass, or repaired wizard URLs are included.

## Local prerequisites

Use the local console, auth app, BFF, Factory, Warehouse and virtual runtime in
one isolated test organization. Use an owner who has completed the welcome
screen. Keep unrelated unassigned node bundles out of this organization: the
original wizard auto-detects a ready bundle. A different selected tower fails
before submission. Supply access/spectrum configuration for the runtime tower
part number, and valid switch, backhaul and power component options.

Rebuild the C CLI normally, install/build the worker, and capture authentication
interactively. Captured files contain credentials and must remain private.

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
node adapters/webapp/dist/cli.js auth \
  --base-url http://localhost:3000 --out .auth/owner.json --landing configure
```

Use the existing `--landing dashboard` default if that is your current landing.
Configure capture requires an actual visible configure heading at the console
origin; it neither creates an account nor acknowledges welcome for you.
The scenarios use `flow=add-network`, so an existing organization's first-network
self-guard does not skip the network form. They do not prove fresh-account
welcome/activation policy, which remains part of live qualification.

Set these independent expectations from your local test inventory. Use the
**visible option descriptions**, not component IDs. Coordinates are the known
Factory/runtime tower coordinates, formatted as displayed on Site Detail.
Do not derive expected values from the screen under test.

| Variable | Value |
|---|---|
| `ULAB_ONBOARD_SWITCH` | Switch option description |
| `ULAB_ONBOARD_BACKHAUL` | Backhaul option description |
| `ULAB_ONBOARD_POWER` | Power option description |
| `ULAB_ONBOARD_COORDINATES` | Expected `latitude, longitude` text |

For example, with your actual local endpoints exported:

```sh
./bin/ukama-lab run scenarios/webapp/p0/onboarding/wb-073-site-persistence.yaml \
  --repo "$UKAMA_REPO" --scripts ./scripts \
  --warehouse-url "$WAREHOUSE_URL" --factory-url "$FACTORY_URL" \
  --asr-url "$ASR_URL" --bff "$BFF_URL" \
  --out runs --run-id onboarding-persistence \
  --webapp-worker ./utils/webapp-worker.sh
```

Use a new run ID for each execution. `commerce.md` documents the shared local
stack/runtime setup, credentials and cleanup prerequisites; `matrix.md` documents
serial browser runs. No endpoint in this patch requires AWS.

## Scenarios

All twelve files are under `scenarios/webapp/p0/onboarding/`.

| Case | Checks |
|---|---|
| WB-070 | Network name from the world, form submission, returned network selected, exactly one mutation |
| WB-071 | Too short/long, uppercase, space and underscore errors; zero mutations |
| WB-072 | Re-select the newly created, now existing network through the radio list; retain its ID |
| WB-073 | Network selection through install/settings; expected tower; selected components; completed site, three saved node IDs and coordinates; reload persistence |
| WB-074 | Each visible checking/creating/confirming stage becomes active in order and completes |
| WB-075 | Controlled missing amplifier; checklist, blocked form, zero site submissions; recovery |
| WB-076 | Controlled offline tower; same readiness/recovery checks |
| WB-077 | Controlled initializing controller; same readiness/recovery checks |
| WB-078 | Controlled missing location; same readiness/recovery checks |
| WB-079 | Forward one site mutation, lose its response, let the app reconcile; no duplicate; saved state survives reload |
| WB-080 | Intercept and reject one site request before forwarding; visible error/retry; one explicit retry reaches the backend; saved state survives reload |
| WB-081 | Original SIM guidance step, absent upload input, and Finish setup navigation; **gap observation only** |

WB-075–080 are controlled response tests. The matrix classifies them as
`controlled_ui` even when an operator labels the manifest `--evidence live`.
They earn no live verified credit. Real missing/offline/unready, transport loss
and backend rejection qualification is still required in Patch 17. A passing
fixture is never evidence that the original React/Apollo application passed.

## Contract and ownership

`webapp.session_mode: onboarding` requires saved auth, no `auth_origin`, at most
one planned network/site/trio, and no commerce/UE/SIM resources. It initializes
an isolated context without claiming authenticated dashboard readiness. Only
`web_onboard` browser events and `web_onboard_equals` checks are accepted in this
mode. Normal authenticated and `auth_test` scenarios remain unchanged.

Events use `view: configure` and a closed `action`:

| Action | Fields / behavior |
|---|---|
| `open` | `value`: overview, network, add_network, select_network, install, sims, complete |
| `fill` | `label`: Network name or Site name; `value` or matching `value_from`: network_name/site_name |
| `validate_name` | Invalid Network name/Site name only; click validation; block any attempted backend create and fail |
| `click` | `value`: Get started, Next, Name site, Installed, Finish setup, Go to Console, Skip for now, Check now |
| `choose_network` | Select the planned network through its visible radio and Continue |
| `select_component` | `label`: Switch, Backhaul, Power; `value`: visible option label |
| `submit_network`, `submit_site` | Save parent intent and worker receipt before submission; pending acknowledgement, no early binding |
| `finish_network`, `finish_site` | Wait for visible next page and bind the matching original creation receipt |
| `retry_site` | Only the known adapter-intercepted rejection; never retry an unknown remote outcome |
| `reload`, `site_detail` | Browser reload or visible Sites navigation and owned card selection |
| `arm_fault`, `clear_fault` | One controlled fault; clear readiness injection before checking recovery |

The fault values are `missing_amplifier`, `offline_tower`, `unready_controller`,
`missing_location`, `transport`, and `rejection`. The readiness faults change
only the planned trio's browser response, never backend/node state. The timeout
forwards the browser request once and aborts its response; the adapter does not
reissue the mutation or use an API to decide acceptance. A rejection is consumed
before any forwarding, making the explicit retry unambiguous.

Checks require `label`, `requirement`, optional `subject`, and exactly one of
`expected` or `expected_property`. World properties are `network_name`,
`network_id`, `site_name`, `tower_id`, `node_ids`. Available labels are Heading,
Path, Step, Network ID, Selected network, Field error/value, Button state,
Readiness/count, Tower, Component, Error, Progress order/complete, Mutation
count, Fault consumed, Site name/nodes, Coordinates, Saved component, SIM
guidance, and SIM upload present. Refer to the scenarios for exact subjects.
Expected booleans and counts are quoted strings. Assertions poll the visible UI;
passive mutation counts supplement those UI checks.

Creation responses establish **cleanup ownership only**. The parent verifies
kind, ref, name and original submit command ID before binding. A prepared receipt
means no submit; an intercepted rejection means no backend forwarding. A lost or
unidentified submitted result remains uncertain, fails the run and is never
retried or guessed away. Teardown uses the existing owned-resource journal and
BFF cleanup path. No pre-existing network is claimed or deleted: WB-072 selects
the network created earlier in that same run.

## Original console gaps and qualification

- WEB-ONBOARD-004 remains partial. Settings expose the tower and saved detail
  exposes all three IDs, but the discovery checklist does not display every
  node identity. No complete automation credit is claimed.
- WEB-ONBOARD-012 remains partial. `/configure/sims` contains instructions and
  Finish setup, with no inventory upload control. An absence assertion documents
  the gap; it cannot satisfy the upload requirement.
- `useCreateSite.ts` marks confirm done without first making it active. The
  strict progress observer is expected to fail that source behavior. No timer
  monkeypatch, hidden-state check or fabricated progress is used to pass it.
- Select-network and install do not preserve the original node-pool `nid`.
  The wizard can choose another ready bundle. A mismatched settings URL/read-only
  tower fails before creating a site; lab does not inject a corrective `nid`.
- Completion navigation is covered; fresh-account self-guard behavior, activation
  policy and service-side fault qualification require the live local stack.

Complete automation is now **54/132 (40.91%)**, including **46/86 P0**. Verification
remains **0/132**. Two of Patch 9's twelve requirements receive no completion
credit. The denominator and the 100% P0 / 90% overall browser gate are unchanged.
