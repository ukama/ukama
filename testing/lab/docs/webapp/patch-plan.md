# Web-app patch sequence

All execution is local to the host. The existing BFF scenario suite stays a
separate test surface and denominator. Each delivery includes changed files,
a unified patch, validation results, and its own commit message.

| Patch | Deliverable | Proposed commit subject | Exit evidence |
|---|---|---|---|
| 1 | v2 foundation contract, offline lint, inventory, WIP examples | Add web-app scenario contract and offline validation | Real C parser/validator tests and legacy compatibility; no browser execution claim. |
| 2 | Local Playwright worker and authenticated browser session | Add local Playwright worker for console-app actions | Chromium opens the target app, establishes/loads valid auth, and checks a visible field. |
| 3 | Runner lifecycle, entity binding, journaling, cleanup, reporting | Integrate web-app steps with lab scenario lifecycle | Ordered runtime/browser execution and failure/cancellation cleanup with one result. |
| 4 | UI provisioning and first operational journeys | Add web-app provisioning and outage recovery scenarios | UI-created network/sites, node list/detail, 2/2 to 1/2 to 2/2 without hidden reload. |
| 5 | Operations, locks, recovery, software update | Cover console-app operations and recovery states | Correct controls, reasons, progress, terminal states, and retry behavior. |
| 6 | Plans, SIMs, customers, UE usage, payments, receipts | Add web-app customer and commerce journeys | Day/minute and data-unit conversion, allocation, usage, payment/entitlement/receipt reconciliation. |
| 7 | Error states, scope, accessibility, browser coverage, coverage report | Expand web-app coverage and release reporting | Agreed critical requirements fully covered, at least 90% of agreed functional inventory, gaps retained. |

Patches 1–4 are implemented. Patch 3 executes zero-fixture authenticated browser
scenarios through the C CLI and verifies transport, ordering, resource ownership,
cancellation and cleanup against controlled local fixtures. Patch 4 connects UI network/site creation, local virtual-node setup, fault
injection and owned-resource teardown. The two operational examples are active.
See `provisioning.md`, `runner.md` and `worker.md`.

The live target-app milestones still need a configured console/auth/BFF/runtime
stack on the host. Fixture passes do not establish those milestones or product
coverage. Later commit subjects remain proposed; each commit body describes
its actual changes and validation.

## Worker contract and provisioning extension

- The C runner owns scenario parsing, phase order, world expectations, deadlines,
  and the final result. It launches one persistent local browser worker per
  scenario and exchanges versioned JSON messages over dedicated pipes.
- Each command has run ID, command ID, action, resolved inputs, and deadline.
  Each response has matching IDs, status, expected/actual detail, duration,
  observed entity bindings, and artifact paths. Logs never corrupt the protocol.
- Browser operations stay serial within a scenario. Playwright locators and
  assertions handle UI readiness; a retry of observation must not resubmit a
  mutation. Command-ID handling must distinguish a lost response from permission
  to repeat a payment or creation.
- Real authentication state stays outside source control and shareable reports.
  A fresh browser context isolates cookies/storage between scenarios. A local
  profile controls console, auth, BFF, factory, and runtime endpoints.
- Setup through the UI drives real dialogs/wizards, journals resources as they
  appear, verifies visible state, and binds actual IDs to lab references.
  Passive browser response metadata may identify a resource for cleanup; it is
  not a substitute for an acceptance assertion.
- On failure, retain screenshot, browser trace, failed assertion, page URL,
  and relevant sanitized console/network diagnostics before cleanup. Browser
  context traces alone do not contain the lab assertions; preserve both artifacts.
- Cancellation terminates the worker and invokes bounded cleanup of journaled
  resources. UI failures cannot disable backend cleanup of test-owned fixtures.
- Implement short valid network names (the supplied form permits 3..40
  lowercase letters/numbers/hyphens) independently from internal run IDs.
- Keep one local worker initially. Later concurrency requires isolated browser
  contexts, accounts where needed, network data, inventory, and runtime nodes.

## Coverage accounting

Freeze the agreed requirement inventory before publishing a percentage.
Automation coverage counts requirements with all necessary browser handlers and
assertions implemented. Verified coverage additionally requires every necessary
assertion to pass on the identified app/backend build and browser profile.
Contract-only, skipped, blocked, xfail, and BFF-only cases earn no verified credit.
Keep real-backend journeys distinct from controlled-response UI-state tests.
