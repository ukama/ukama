# Local customer and commerce journeys — Patch 6

Apply **ukama-lab only**, after lab Patches 1–5. Ignore the console companion
archives from Patches 4 and 5. This patch requires no console source changes.
It replaces their test-ID dependency with locators from the original app.

This document describes the Patch 6 foundation. For the current extensions and
remaining gaps, read [lifecycle.md](lifecycle.md) (Patch 12) and
[payments.md](payments.md) (Patch 13).

## Execution and fixture ownership

The C runner creates networks through the UI as before. If the world contains
SIMs/UEs, a bounded local Factory/Warehouse job prepares a fresh run-specific
batch, exports its private CSV and assigns its ICCIDs/IMSIs to the world.
The parent reloads that CSV after the bounded child exits. Playwright uploads
it with the console file picker and reads every expected ICCID's Available
status from the pool screen. No BFF upload or hidden acceptance query is used.

Plans and customers are explicit phase events, followed by separate allocation.
The customer form deliberately leaves its optional plan empty, so creation
and allocation have separate durable ownership receipts. Allocation selects the
exact ICCID, never auto-assigns a foreign pool SIM. Every action is scoped to
an explicitly opened view and the selected world network/customer.

`start_ues` waits for ASR registration, ensures local media, starts the selected
UEs and waits for attachment. `traffic` uses the existing local runtime script.
Neither is an acceptance assertion. The scenario then checks the app's visible
usage values. Node/runtime provisioning and teardown keep their existing roles;
operator actions run in Playwright. Execution is on the host, not AWS.

Before each creation/payment the C runner saves intent. The worker records
prepared/submitted/identified state in a private `creation-N.json`. Passive UI
mutation metadata supplies only ownership/correlation IDs. Prices, allowance,
validity, status, usage and receipt contents are accepted only from rendered UI.
No browser request client, response interception, injected state or API fallback
is used. An unidentified submission is never retried or guessed into ownership.

Cleanup deletes only identified run-owned SIM allocations, subscribers, plans,
nodes, sites and networks, in dependency order. SIM packages are removed and
service is switched off before deletion. Cash payment ledger entries remain
`retained_ledger`; Factory/pool inventory remains `retained_factory_pool` in the
journal. These are intentional retained records, not refunds or deletions.
Use a disposable local test organization. UI timeouts with an identified ID
still allow teardown; an uncertain creation fails the run and requires manual
reconciliation from the journal. Browser evidence and CSV credentials stay private.

## Writing scenarios

Use existing world package declarations. A network package is referenced as
`REF__net-001`; an organization package is `REF__org`. Use either `duration_days`
or `duration_minutes`, not both. Valid browser durations are:

| Frontend days | Canonical backend minutes | Visible validity |
|---|---|---|
| 1 | 1440 | 1 day |
| 7 | 10080 | 1 week |
| 30 | 43200 | 1 month |

The lab translates minutes to the form's day choice. Reloaded cards, read-only
edit fields and allocated package date spans check the reverse path. Package
spans expose whole calendar days; this does not claim sub-day timestamp precision.
Minute-only durations such as 5 or 90 cannot be created through this form and
are rejected by v2 validation. Existing BFF scenarios can continue testing them.

`data_mb` uses binary units: 1024 MB = 1 GB = 1073741824 bytes. `create_plan`
requires `unit: MB|GB`; GB requires an exact multiple of 1024 MB. The operator's
organization supplies currency/country in the actual app. Included scenarios
expect USD formatting; configure that test organization or adjust their explicit
price expectations. The worker does not alter organization settings.

```yaml
packages:
  - ref: weekly
    name: weekly
    data_mb: 1024
    duration_minutes: 10080
    amount: 10
    currency: USD
    country: USA
    scope: network
phases:
  - name: create-plan
    events:
      - type: web_open
        view: business_data_plans
        networks: net-001
      - type: web_commerce
        view: business_data_plans
        networks: net-001
        action: create_plan
        package: weekly__net-001
        unit: GB
  - name: persisted-plan
    events:
      - type: web_reload
    checks:
      - type: web_commerce_equals
        view: business_data_plans
        package: weekly__net-001
        label: Plan terms
        expected: "1 GB data · 1 week validity"
        requirement: WEB-PLAN-003
```

This is a phase fragment; the supplied examples contain the complete v2 header,
world, browser profile and UI setup. `sims_per_network: 1` creates one customer/SIM
without node runtime. `ues_per_site: 1` gives a SIM an actual local site/IP for
traffic. Patch 6 supports one SIM per customer, at most 100 SIMs total, and one
cash top-up per UE per run. A second command ID cannot repeat the same mutation
reference. Replaying the original command returns its cached response.

| Action | View | Required selectors/fields |
|---|---|---|
| create_plan | business_data_plans | networks, package, unit |
| edit_plan | business_data_plans | networks, package; opens the existing edit form |
| create_customer / open_customer / close_customer | customer_customers | networks, ues |
| allocate_sim / top_up / cancel_top_up / open_receipt | customer_customers | networks, ues, package |
| activate_sim / deactivate_sim | customer_customers | networks, ues |
| close_dialog | either commerce view | networks |

Actions do not navigate implicitly. Open the view first. Open the customer drawer
before allocation, SIM toggling, top-up or receipts. Use a distinct top-up plan
so the package row is unambiguous. `edit_plan` enables inspection; changing an
existing plan's name is not implemented by this patch. Reload is always explicit.

`web_commerce_equals` requires a WEB-* requirement and exactly one `expected`
or `expected_property`. Values compare exactly after whitespace normalization.
Identity properties are `iccid`, `plan_name`, or `payment_id`, resolved by the C
parent from the planned world/identified mutation, not copied from the assertion
DOM. Reports include the package/UE selectors and expected-property name.

| View | Labels | Scope |
|---|---|---|
| business_data_plans | Plan terms, Plan price, Plan scope | package |
| business_data_plans, edit open | Validity, Price, Data volume, Unit | package |
| business_sim_pool | Pool status | ues |
| business_packages | Performance price, Performance sold, Performance revenue, Performance share | package |
| customer_customers, drawer open | ICCID, SIM status, Phone, Active plan, Cycle usage, Total usage | ues |
| customer_customers, drawer open | Package count, Package status, Package dates, Package days | ues + package |
| customer_customers, receipt open | Receipt total, Receipt payment ID, Receipt method, Receipt status, Receipt empty | ues |
| customer_customers, receipt open | Receipt plan | ues + package |

Missing/ambiguous/hidden values fail. Polling an assertion never reloads or submits
an action. A toast or mutation ID cannot make an entitlement or receipt check pass.

## Included journeys and host commands

| Scenario | Main assertions |
|---|---|
| wb-020 | Daily/weekly/monthly plan cards, MB/GB, price/scope, reload and edit validity |
| wb-021 | Factory SIM import, customer, exact allocation, weekly span, reload, no fabricated receipt |
| wb-022 | Cancel top-up, one cash payment, one queued package, receipt correlation, Business revenue/purchase totals, per-plan sales, reload persistence |
| wb-023 | Visible SIM service off/on |
| wb-024 | Site-bound UE attachment and traffic, explicit reload, cycle and total usage |
| wb-025 / wb-026 | Daily/monthly allocation spans and allowance after reload |
| wb-027 | Complete UI network/site/customer journey, local UE traffic, top-up, receipt and Business totals (36 visible checks) |

Build the lab using its normal host build, then:

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
npm --prefix adapters/webapp run install:browser
./utils/webapp-worker.sh auth --base-url http://localhost:3000 --out .auth/owner.json
```

Capture owner auth as described in `worker.md`. Use the same organization for
console/auth/BFF, Factory and teardown credentials described in `provisioning.md`.
Edit `webapp.base_url`/`auth_state` in your local profile copies as needed.

```sh
./bin/ukama-lab run scenarios/webapp/p0/commerce/wb-022-cash-topup-receipt.yaml \
  --repo /absolute/path/to/ukama --scripts ./scripts \
  --bff http://localhost:8080/graphql \
  --warehouse-url http://localhost:8070 --factory-url http://localhost:8071 \
  --asr-url http://localhost:8072 \
  --webapp-worker ./utils/webapp-worker.sh --out ./runs --run-id commerce-001
```

The ports above are examples; use your actual local service endpoints and the
console's configured SIM type (`--sim-type` must agree). Use a new run ID each
time. A top-up's ledger record survives cleanup; rerunning is a new purchase.

For wb-024 and wb-027, independently determine the expected displayed traffic totals for
the local runtime/accounting setup, including its accounting overhead. Set:

```sh
export ULAB_WEBAPP_EXPECTED_CYCLE_USAGE='64 MB of 1 GB used this cycle'
export ULAB_WEBAPP_EXPECTED_TOTAL_USAGE='64 MB'
```

These illustrative expectations apply only where the 64 MiB workload is accounted
as 64 MiB. Never populate them by reading the same screen being tested. The v2
loader refuses missing environment expectations; there is no loose fallback.

The cash scenarios independently expect one $25 sale in the new test network:
Business Home Revenue = $25, Customers = 1; Revenue Purchases = 1,
Avg purchase = $25, Paid customers = 1; Packages revenue = $25 and sold = 1.
The named top-up plan row must show price/revenue $25, sold 1 and share 100%.
The initial allocation has no recorded payment and contributes no sale revenue.
Receipt formatting is $25.00; dashboard currency formatting is $25.
These are literal expectations derived from the purchase, never copied across
screens. Analytics checks allow 900 seconds and do not silently reload.
The supplied scenarios use the app's default rolling date span; date-range
switching and cross-network isolation remain Patch 7 work.

## Retained limits

No live console/auth/BFF/Factory/Podman stack was available during delivery.
Controlled Chromium/C fixture passes qualify the lab infrastructure only.
All 132 requirement verification states remain `not_run`. The 90–100% product
coverage goal has not been achieved. The inventory keeps partial requirements
uncredited; no new duplicate denominator entries were added.

The original wizard loses the tower ID, so site/UE scenarios may fail before
site creation. Original operation locks/polling/optimistic toggles may also fail
their scenarios. The lab never repairs URLs, injects state or changes the console.

The pool list has a cap and no search/pagination control. A missing imported row
fails instead of being verified by API. Customer email is not displayed in the
drawer; receipt payer identity is hardcoded Walk-in customer. Full email/payer
reconciliation therefore stays a gap. Patches 7, 12 and 13 add invalid-form and
scope probes, guarded auto-assignment, rapid double-clicks, controlled payment
rejection and downloaded receipt/date checks. Multiple SIM selection,
expiry/remaining-allowance semantics, terminal failed-payment records and actual
UE traffic denial during service-off remain gaps. See the later guides for exact
boundaries. Command replay and app rapid-click protection are tested separately.
