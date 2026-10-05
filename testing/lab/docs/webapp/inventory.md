# Patch 10 — navigation, inventory and maps

Ukama-lab only. Apply after lab Patch 9. The original uploaded console remains
unchanged; ignore the Patch 4/5 console companions. Execution stays on the local
host. This patch adds 18 scenarios, WB-090–107, under
`scenarios/webapp/p0/inventory/`.

## Before a live run

Use the existing owner session, local BFF/Factory/Warehouse/Podman configuration,
worker installation and owned-resource cleanup described in `worker.md`,
`commerce.md` and `matrix.md`. The scenarios create their networks and sites
through the browser. Existing console onboarding defects can therefore stop
setup before an inventory assertion. The lab does not repair those defects.

WB-092 and WB-099 require independently known site location text; WB-092 also
requires coordinates in the console's decimal `latitude, longitude` format:

```sh
export ULAB_WEBAPP_SITE_LOCATION='your known site location'
export ULAB_WEBAPP_SITE_COORDINATES='your latitude, your longitude'
```

These must describe the Factory/runtime sites configured for this run. The
scenario checks their displayed values; it does not set coordinates or copy
expectations from the console's data response. The current scenario uses one
location expectation for its test sites. Adjust the declarative assertions if
your prepared sites differ.

WB-103 uses three independent percentage expectations derived from known
site-time history for the target network and the selected windows:

```sh
export ULAB_WEBAPP_UPTIME_DAY='98.5%'
export ULAB_WEBAPP_UPTIME_WEEK='97.2%'
export ULAB_WEBAPP_UPTIME_MONTH='95.0%'
```

Those are examples, not defaults. Prepare a reproducible local history and use
its expected values; missing or unprepared analytics are a qualification gap.
The adapter never reads an API value to supply its own expectation.

WB-106 extends the existing local UE journey: UI plan/customer/SIM allocation,
UE attach and 64 MB traffic, then Network Home on two networks. Supply the
existing `ULAB_WEBAPP_EXPECTED_CYCLE_USAGE` and
`ULAB_WEBAPP_EXPECTED_TOTAL_USAGE`, plus `ULAB_WEBAPP_NET1_ACTIVE`,
`ULAB_WEBAPP_NET1_VOLUME`, `ULAB_WEBAPP_NET2_ACTIVE` and
`ULAB_WEBAPP_NET2_VOLUME`. These values follow the known traffic/accounting
contract and scenario scope. Missing KPI readings must not be called zero.
Browser fixtures use `1`, `64 MB`, `0`, `0 B`; this is fixture qualification,
not evidence those exact accounting values hold on your stack.

Rebuild the C CLI and worker. Example with your existing local runner options:

```sh
npm --prefix adapters/webapp run build
bin/ukama-lab run scenarios/webapp/p0/inventory/wb-090-site-membership.yaml \
  --repo /path/to/ukama --bff http://localhost:8080/graphql \
  --factory-url http://localhost:8081 --warehouse-url http://localhost:8082 \
  --out runs/webapp --run-id inventory-001 \
  --webapp-worker "$PWD/utils/webapp-worker.sh"
```

Use your actual endpoints and required runtime options. The matrix runner in
`matrix.md` supports the same suite with selected browsers and evidence reports.

## Scenarios and scope

| Scenario | Assertions |
|---|---|
| WB-090 | Complete site names and header total in two networks, then reload |
| WB-091 | Exact owned-name search, explicit no-match state, clear and restored membership |
| WB-092 | Site list location/count and detail name/location/coordinates/all three node IDs |
| WB-093/094 | Network switch from site/node detail lands on new-network list without corrective navigation |
| WB-095/096 | Sticky observation of foreign site/node identities during network switching |
| WB-097 | Browser Back and direct site/node URLs, each followed by reload, retain the expected world identity and network |
| WB-098 | Controlled stale stored network selection self-heals to a visible menu option |
| WB-099 | Map count, site selection/name/location/color token, exact detail link, second site and Clear |
| WB-100 | Tower outage, explicit reload, Degraded status, preserved node IDs, reconnect and Online |
| WB-101 | Complete node membership/count across two networks and reload |
| WB-102 | Controlled missing KPI and registry data show unavailable, then recover |
| WB-103 | Uptime percentage for three explicit date windows |
| WB-104 | Explicit empty site/node states and zero totals |
| WB-105 | Distinct Sites online KPI scopes through outage, switch, reload and recovery |
| WB-106 | UI-created customer and real local UE traffic reconciled on both networks |
| WB-107 | Connectivity remains Online while Radio and Cellular change independently |

## Inventory command contract

`web_inventory` events accept `view`, `networks`, `action`, optional `sites` or
`nodes`, and `value` only where allowed. Views are restricted to Network Home,
Sites, Nodes, site detail and node detail. References must belong to the declared
world network. No raw selector, arbitrary URL, arbitrary JavaScript, arbitrary
request or generic mutation is accepted.

| Action | Inputs and behavior |
|---|---|
| `switch` | New list/home view and destination network; click actual network switcher |
| `switch_watch` | Same, but requires at least two networks with site/node identities; arm DOM observer before clicking |
| `search` | Sites view; one `sites` reference or literal `value`, including empty to clear |
| `direct` | Known view and matching entity for detail; same-origin URL from resolved world ID; no network correction |
| `back` | Browser Back only; destination is asserted separately |
| `map_select` | Home with one site; click visible pins and identify the site by its rendered popup |
| `map_open` | Home with the same selected site; click Open site exactly once |
| `map_clear` | Home; click Clear |
| `stale_selection` | Home; change only networkId in isolated `uk-ui-prefs`, reload; controlled evidence |
| `mask_home` | Home with `value: kpis`, `registry` or `reset`; controlled read-response masking |

`web_inventory_equals` checks accept `view`, `networks`, `label`, optional
matching `sites`/`nodes`, `requirement`, timeout, and exactly one expectation.
Identity labels require the specific `expected_property`: `path`,
`network_name`, `site_names`, `site_name`, `node_ids`, `site_node_ids`, or `node_id`.
The C runner resolves those values before dispatch. Whole memberships are sorted
arrays that retain duplicates; the parent independently compares returned JSON.
Reports preserve full expected and observed values. Literal fields include
header totals, site location/coordinates/count/status, map count/color/location,
empty/clear states, scope leaks and the applied fault. See the shipped scenarios
for exact mappings.

Checks do not navigate or reload to make a failed state pass. The one exception
is `Selection valid`, a read-only inspection that opens the network menu and
closes it with Escape, without choosing a network. Mutations remain on the
existing UI provisioning/operation/commerce paths and their ownership journal.

## Evidence and remaining gaps

- `switch_watch` observes rendered DOM at mutation callbacks and animation
  frames. A captured foreign identity remains a failure even after convergence.
  It cannot prove every paint, canvas content, anonymous KPI value or other
  unvisited view. WEB-SHELL-005 remains partial.
- Map checks use public rendered Leaflet pins/popups and selection content.
  They do not inspect Leaflet internals, React state, GraphQL values or app store.
  Geographic pin placement and offline/degraded colors remain unverified.
  WEB-NET-006 stays partial; the online SVG fill token alone is not a full visual
  color/position qualification.
- Original Sites has search only. Original Nodes has no search/filter/sort/
  pagination controls. WEB-SITE-005 and WEB-NODE-004 remain partial. Automatic
  node type ordering does not count as a user sorting control.
- The site list has name/location/node count, but no coordinates or individual
  node IDs. Detail checks validate those values independently. Complete
  list-to-detail equivalence for WEB-SITE-002 remains partial.
- `stale_selection` simulates an invalid preference without deleting any server
  network. `mask_home` changes only `GetKpiValues`/`SitesList` read responses,
  requires proof the selected fault was applied, and never retries a mutation.
  WB-098/102 cannot receive live credit even in a manifest labeled `live`.
  Genuine removed-network and service-side failure qualification remains open.
- PageHeader uses a `.pagetitle` div, not a heading role. The provisioning
  adapter now waits for that actual source element after site creation.

The fixed inventory remains 132 requirements, 86 P0. Complete automation is
66/132 (50%), including 54/86 P0; live verified remains 0/132. All five partial
requirements stay in the denominator. Patches 11–16 and live qualification in
17 remain. This patch is not a completed coverage gate.
