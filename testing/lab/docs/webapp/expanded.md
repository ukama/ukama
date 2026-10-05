# Patch 7: browser interactions and coverage evidence

This patch changes **ukama-lab only**, on top of lab Patch 6. Continue to ignore
both console companion archives from Patches 4 and 5. All selectors target the
original uploaded console. The adapter now opens customer rows using the table's
`tbody > tr`: clickable rows in the app override the row role with `role="button"`.

## New scenarios

All examples are in `scenarios/webapp/p0/expanded/`. They use the same host-only
console/auth/BFF/Factory/runtime setup and cleanup rules as `commerce.md`.

| Scenario | Browser acceptance |
| --- | --- |
| wb-030 | Duplicate plan name shows an error, disables submission, and leaves one plan after cancel/reload |
| wb-031 | Empty name and nonpositive price/volume show field errors and block submission |
| wb-032 | Two networks isolate customer rows and network plans; an organization plan is available in both |
| wb-033 | Customer name search, no-match state, descending sort, SIM filter and no-plan filter |
| wb-034 | Explicit empty customer state in Business, Network and Customer views |
| wb-035 | Missing-plan prerequisite blocks the customer dialog with the expected message |
| wb-036 | Invalid email shows a field error and disables customer creation |
| wb-037 | Keyboard opens the plan dialog; Escape closes it and restores opener focus |
| wb-038 | A 390 × 844 viewport opens mobile navigation, follows a link and closes the drawer |
| wb-039 | A recent purchase has the expected total under all three rolling date ranges |
| wb-040 | Clearing the browser's cookies redirects to the explicitly configured auth origin |
| wb-041 | Top-up submission stays disabled until a plan is selected |
| wb-042 | Customer identity and page routes persist across three views and reloads |
| wb-043 | Keyboard command palette opens, searches and navigates to Data plans |

`wb-040` requires `ULAB_WEBAPP_AUTH_ORIGIN` (origin only, including scheme and
port). It clears cookies in this run's isolated browser context, never the saved
auth file. It requires zero created world resources so cleanup cannot depend on
the cleared session. This verifies missing-session redirection; it does not claim
logout-provider or expired-token recovery coverage.

The date-range scenario creates a real cash purchase, as in Patch 6. Ledger
records and Factory/pool inventory are retained and recorded in the private
resource journal. Run only against your designated test organization.

## `web_interact`

Every event declares `view` and `action`. Optional `networks: net-NNN` verifies the
current selected network on scoped pages. The worker maps semantic labels to
source-derived controls; YAML never contains CSS selectors or JavaScript.

| Action | Arguments |
| --- | --- |
| `search` | `value` on one of the three customer lists |
| `sort` | `label`: Customer, Data usage, or Last seen |
| `filter` | `label`: SIM or Active plan; `value`: exact menu item, including All |
| `open_form` | `label`: Create plan or Add customer; optional `value: keyboard` |
| `fill`, `select` | `label` names a mapped dialog field; `value` is text or native option label |
| `cancel` | Closes the current dialog using Cancel |
| `press` | `value`: Tab, Shift+Tab, or Escape; optional button `label` to focus first |
| `viewport` | `value`: desktop (1440 × 1000) or narrow (390 × 844) |
| `date_range` | `value`: Last 24h, Last 7 days, or Last 30 days |
| `go_back` | Browser Back |
| `palette` | Control+K opens the command palette |
| `palette_choose` | `value`: result label; focus and Enter navigate |
| `mobile_open`, `mobile_link` | Open navigation, or follow the exact link in `value` |
| `open_allocate`, `open_topup` | Open the visible customer drawer's corresponding dialog |
| `clear_session` | Only on `view: session` with zero world networks |

Fields supported by `fill`/`select`: Data plan name, Price, Data volume, Unit,
Validity, Network, First name, Last name, Email, Data plan, SIM. Controls retain
their actual native input/select behavior. Actions fail if the field is absent,
ambiguous, disabled or on the wrong page. There is no generic submit action.

For `fill`, `select` and `search`, substitute `value_from: plan_name`,
`customer_name` or `network_name` for `value`, with a matching `package`, `ues`
or `networks` reference. The C runner resolves only identities already created
and bound to this world. UI actions cannot return resource-creation bindings.
Continue using journaled `web_commerce` actions for mutations.

## `web_ui_equals`

Declare `view`, `label`, `expected` and `requirement: WEB-*`. Text comparison is
exact after whitespace normalization. Boolean expectations are quoted strings
`'true'` / `'false'`. A missing field does not become an empty string or zero.

| Label | Additional selector |
| --- | --- |
| Path; Selected network; Dialog open; Dialog title; Mobile navigation open; Date range; Content fits viewport | None |
| Field error; Field value; Field readonly | `subject`: mapped field label |
| Text visible | `subject`: exact visible text |
| Button enabled; Button visible; Focus on button | `subject`: exact button label in the dialog, or main content |
| Customer present | `ues`: world customer identity; customer list view |
| Customer order | Customer list; expected visible names joined with `\|` |
| Plan option present | `package`: world plan; Add customer form |
| Plan count | `package`: world plan; Business Data plans |
| Auth origin; Dashboard visible | Only `view: session` |

For text expectations, `expected_property: plan_name|customer_name|network_name`
can replace `expected` with the corresponding world reference. Loading skeletons
are not counted as customer rows. Absent customers are accepted only in a ready
table or an explicit empty/no-match state. Checks wait for route transitions and
use the existing command/scenario deadlines.

## Remaining scope

The unchanged inventory has **132 requirements**, including 86 P0 requirements.
Patch 7 raises fully implemented automation from 29 to **36** (27.27%); 32 of 86 P0
requirements are automated. Live verification remains **0** here. The original
90–100% coverage target and the Patch 7 release milestone are **not met**.

The seven newly complete entries are WEB-AUTH-001, WEB-SHELL-001/008/009 and
WEB-PLAN-006/009/010. Other new scenarios deliberately retain partial credit only:
customer pagination has no control in the uploaded DataTable; cross-network
counts/KPIs, historical date exclusions/trends, broad keyboard/focus coverage,
role permissions, error/retry states, analytics charts and several mutation
recovery paths remain unimplemented. `coverage.json` retains every requirement
and records specific partial gaps. A fixture pass never changes verification.

See `matrix.md` for serial browser runs, report generation and the release gate.
