# Authentication, session and role acceptance — Patch 8

These scenarios change **ukama-lab only** and target the original uploaded
console. They drive real browser navigation and visible account state. The
local two-origin fixture tests qualify the lab implementation, not the live
React/Apollo/auth stack. All inventory verification remains `not_run`.

## Preparation

Use an isolated local test organization with a working console, auth app and
BFF. The organization for ordinary and welcome accounts must already have a
network and site so the activation guard does not divert to `/configure`.
No networks, nodes, SIMs or customers are provisioned by this suite.

Capture each account through the normal interactive sign-in flow:

```sh
npm --prefix adapters/webapp run build
node adapters/webapp/dist/cli.js auth \
  --base-url http://localhost:3000 --out .auth/owner.json
node adapters/webapp/dist/cli.js auth \
  --base-url http://localhost:3000 --landing welcome --out .auth/welcome.json
node adapters/webapp/dist/cli.js auth \
  --base-url http://localhost:3000 --landing unauthorized --out .auth/no-org.json
```

`--landing dashboard|welcome|unauthorized` defaults to dashboard. The capture
command requires a terminal/display, verifies the console origin and visible
landing, and writes private state. It never creates sessions or signs tokens.
Leave the welcome screen unacknowledged when capturing welcome state.

Set these variables from independently known test-account records, not values
extracted from a saved token or returned by the page being tested:

| Variable | Value |
|---|---|
| `ULAB_WEBAPP_AUTH_ORIGIN` | Exact auth origin, e.g. `http://localhost:3001`; distinct from the console; no path or trailing slash. |
| `ULAB_AUTH_ORG` | Expected displayed organization name. |
| `ULAB_AUTH_COUNTRY` | Expected formatted country label, e.g. `United States`. |
| `ULAB_AUTH_CURRENCY` | Expected displayed currency, e.g. `USD ($)`. |
| `ULAB_AUTH_OWNER_NAME`, `ULAB_AUTH_OWNER_EMAIL` | Expected owner identity. Email must be verified for WB-050. |
| `ULAB_AUTH_OWNER_STATE` | Path to a real owner storage-state file with session and gateway cookies. |
| `ULAB_AUTH_ADMIN_STATE`, `ULAB_AUTH_NETWORK_OWNER_STATE`, `ULAB_AUTH_VENDOR_STATE`, `ULAB_AUTH_MEMBER_STATE` | Separate real account snapshots. Set matching `ULAB_AUTH_<ROLE>_NAME` and `_EMAIL` for each. |
| `ULAB_AUTH_NO_ORG_STATE`, `ULAB_AUTH_NO_ROLE_STATE` | Real independently prepared unresolved-organization and unresolved-role accounts, captured at unauthorized. No usable cached gateway token. |
| `ULAB_AUTH_LOGOUT_STATE` | A **separate sign-in session** for the expected owner identity. Do not copy the normal owner snapshot: the logout service may invalidate that server session. |
| `ULAB_AUTH_EXPIRED_STATE` | An independently expired/revoked real owner session captured before expiry; retain its original cached token to exercise stale-access recovery. |
| `ULAB_AUTH_WELCOME_STATE`, `ULAB_AUTH_WELCOME_NAME` | Disposable first-visit **Member** in the activated test organization. |

Paths resolve from the lab working directory. Keep `.auth` private and outside
version control. Only variables referenced by the selected YAML are required.
The state files must come from real sign-ins; dummy or fabricated cookies belong
only in `adapters/webapp/test/session-fixture.mjs`.

The expired-session precondition is established at the auth service or by waiting
for that real session to expire, outside this scenario. It is not simulated by
clearing the session cookie or forging its value. The lab deliberately does not
silently repair an expired snapshot. WB-057 fails if stale dashboard access
persists. WB-054/055 fail if unresolved accounts reach the dashboard.

## Scenarios

| ID | Acceptance |
|---|---|
| WB-050 | Owner account, organization, country, currency and verified email in all three Settings lenses; reload and account menu. |
| WB-051 | Remove only the isolated gateway cookie; recover through the valid session. |
| WB-052 | Malformed gateway cookie recovery. |
| WB-053 | Preserve token payload but damage its signature; require observed `/api/auth/refresh`, then correct persisted identity. |
| WB-054 / 055 | Unresolved organization / role: unauthorized title, logout/support links, no dashboard or expected account data, persistent denial. |
| WB-056 | Account-menu logout, actual auth-app `/user/logout` handoff, auth landing and denied protected deep links afterward. |
| WB-057 | Independently expired/revoked session: blocked protected routes and absent stale account/organization data. |
| WB-058 | Welcome country/org/member identity, Continue, persisted removal of the gate after revisiting `/welcome` and reloading. |
| WB-059 | Set only the isolated token expiry claim to a past value; gateway refresh must restore the expected account. |
| WB-060–063 | Admin, Network owner, Vendor and Member identities/role labels, account menus and reloads. |
| WB-064 | Empty isolated browser context; protected deep links reach the configured auth origin without dashboard data. |
| WB-065 | Configurable route/control policy probe; partial permission coverage only. |
| WB-066 | Billing absent in all three sidebars; direct route returns HTTP 404 **and** visible 404 page; Settings remains reachable. |

WB-056 invalidates its real logout session. WB-058 permanently acknowledges
welcome for its account. Prepare fresh sessions/accounts before repeating them
or running another browser. Do not reuse a consumed welcome account and count
the resulting failure as infrastructure noise. The normal saved-state files are
never overwritten by the worker. Gateway-token faults affect only its isolated
browser context; cookie scope and flags and the real session are preserved.
WB-059 edits the source-defined token expiry field solely to inject a fault,
never to derive expected identity or mint a valid signature.

Run one scenario using the usual rebuilt lab CLI and local endpoint flags:

```sh
./bin/ukama-lab run scenarios/webapp/p0/auth/wb-050-account-settings.yaml \
  --repo /path/to/ukama --out runs --run-id auth-owner-001 \
  --webapp-worker ./utils/webapp-worker.sh
```

The endpoint/runtime setup and CLI build instructions remain in `runner.md`.
To use the matrix/report tool, see `matrix.md`; explicitly select these scenarios
and prepare the consumed-account preconditions for each browser run. Defaults
must not be interpreted as a precondition reset. Evidence must identify the
actual console/backend build and browser. A fixture run receives no live credit.

## Role policy remains a gap

`src/lib/roles.ts` defines Owner, Admin, Network owner, Vendor and Member. The
uploaded sidebar is static across roles; labels alone establish no permission
policy. No console permissions were changed by this patch.

WB-065 uses `ULAB_AUTH_POLICY_STATE` and these eight independent expectations:

| Route group | Expected surface | Expected control |
|---|---|---|
| Members | `ULAB_POLICY_MEMBERS_SURFACE` | `ULAB_POLICY_MEMBERS_CONTROL` (Invite member) |
| Data plans | `ULAB_POLICY_PLANS_SURFACE` | `ULAB_POLICY_PLANS_CONTROL` (Create plan) |
| SIM pool | `ULAB_POLICY_SIMS_SURFACE` | `ULAB_POLICY_SIMS_CONTROL` (Upload SIMs) |
| Customers | `ULAB_POLICY_CUSTOMERS_SURFACE` | `ULAB_POLICY_CUSTOMERS_CONTROL` (Add customer) |

Surfaces are `dashboard`, `unauthorized`, `auth`, `welcome` or `not_found`;
control states are `enabled`, `disabled` or `absent`. Supply an approved policy,
then repeat the probe for each agreed role. Do not infer the policy from current
UI visibility. The probe does not execute permission-sensitive mutations and
cannot complete WEB-SHELL-010. Allowed and forbidden action execution plus the
full agreed route set remain planned, with closure in Patches 16–17. Inviting a
recipient, for example, is outside Patch 8.

## DSL and enforcement

The existing profile defaults to `session_mode: authenticated`; normal resource
journeys still require a verified dashboard at initialization. Only the explicit
`auth_test` profile allows `auth_state: none` and initialization without a
verified account. Its acknowledgement is `initialized: true`,
`authenticated: false`, `session_mode: auth_test`; the C parent verifies all
three. It permits no world resources, provisioning, runtime events or ordinary
web actions/checks. Auth assertions earn credit only after explicit browser
steps; initialization itself earns none.

```yaml
webapp:
  session_mode: auth_test
  auth_origin: ${ULAB_WEBAPP_AUTH_ORIGIN}
  auth_state: none
# Other required profile/world fields are shown in the supplied scenarios.
phases:
  - name: missing-session
    events:
      - type: web_session
        view: session
        action: navigate
        value: /business/settings
    checks:
      - type: web_session_equals
        view: session
        label: Surface
        expected: auth
        requirement: WEB-AUTH-001
```

Events: `navigate` (closed console-route list), `reload`, `settings_tab`
(My account/Organization/Preferences), `open_account`, `logout`, `ack_welcome`,
`drop_token`, `invalidate_token`, `reject_token`, `expire_token`. Only navigate
and settings_tab take a value. Arbitrary URLs, selectors and JavaScript are
rejected by the C contract and the worker.

Checks use exact normalized visible text. `subject` is required only for
Settings field, Welcome field, Sensitive text absent, Nav visible and Control
state; named fields/controls are closed lists. Other labels cover Surface,
Path, Auth origin, Dashboard visible, Organization, Account name/details,
Welcome title/error, Access blocked, Logout handoff, Refresh observed,
Navigation stable, Unauthorized title/logout/support and Document status.
Navigation must settle for 500 ms, with a bounded document-request count. This
is a convergence check, not a guarantee against every possible delayed loop.
Auth HTTP error pages do not count as successful auth landings. Billing checks
require an actual document 404; a page merely displaying “404” is insufficient.

State values and raw token claims are never emitted as command results. Account
identity is necessarily present in private assertion evidence. Failure screenshots
and Playwright traces can contain sensitive page/session data; use the existing
private-artifact handling described in `worker.md`. Reports are acceptance
records, not a signed security attestation. Matrix environment hashes bind the
selected paths/expectations, not the mutable contents of credential files.
