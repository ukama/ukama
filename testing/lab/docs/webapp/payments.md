# Payment and receipt qualification — Patch 13

Apply after lab Patch 12. This patch changes ukama-lab only. Use the local owner
session, organization currency/country and Factory/runtime prerequisites in
`commerce.md`. No console companion patch is required or included.

## Run locally

Rebuild the C CLI, run `npm --prefix adapters/webapp ci`, then
`npm --prefix adapters/webapp run build`. Install **Poppler `pdftotext` and
`pdftoppm` on PATH** for receipt download validation (`poppler-utils` on Debian/
Ubuntu or `poppler` on macOS). These tools are only needed by WB130.

Use the usual local-stack `ukama-lab run` options with files under
`scenarios/webapp/p0/payments/`. Saved browser sessions must match the configured
local stack. All assertions remain browser-visible; backend response listeners
record mutation ownership rather than supplying expected UI values.

| Scenario | Behavior |
|---|---|
| WB130 | Cash top-up, exact receipt correlation after reload, date window, real UI download, PDF extraction and rendering |
| WB131 | Controlled rejection, visible error, no entitlement after reload, one explicit successful retry |
| WB132 | Real rapid double-click, in-flight payment ownership, exactly one displayed entitlement and correlated receipt after reload |
| WB133 | Persisted no-SIM customer has no top-up action; complements WB041's empty-plan guard |

## Downloaded receipts

`web_commerce action: download_receipt` requires the owned customer, purchased
package and identified payment. It snapshots the rendered receipt, clicks
Download, checks `receipt-<first-eight-payment-ID-chars>.pdf`, bounds the retained
file to 5 MiB and validates its PDF signature. Poppler extracts UTF-8 text and
renders the first page to a PNG. The lab retains private `receipt-N.pdf`, `.txt`
and `.png` files with mode 0600 and registers them in the artifact results.

Comparison covers the payment ID, short receipt number, paid date, cash method,
displayed payer/SIM, organization, status, plan, line amount, total, quantity and
footer. Metadata, billed-to, total and payment ID must also occur in the correct
text sections, so a correct line-item amount cannot hide an incorrect total.
The dialog must still contain the same values after download. `Receipt PDF`
returns `matched` only for that exact current receipt. Rendering produces review
evidence; this is not a pixel-perfect layout assertion or OCR check.

`Receipt date window` parses the displayed UTC minute and compares it with the
lab's local submit/identified-completion timestamps, allowing two minutes on
either side for minute precision and clock skew. Synchronize the local host and
payment service clocks. A missing, malformed or old date cannot pass. This is
an independent host-clock check, not a date copied from a GraphQL response.

## Controlled rejection and ownership

`failed_top_up` has a durable payment intent. Its temporary browser route blocks
all writes while submitting, fulfills only the exact owned payment request at
the single passively observed commerce endpoint, and passes reads through.
Wrong endpoint, SIM, plan, amount, currency or unexpected mutation is blocked.
The exact matching request receives a GraphQL rejection; no ledger identity is
fabricated. The private receipt records `intercepted`, and C recovery records
`not_submitted`. A visible error and an applied-rejection check are required.
The route is removed in all outcomes. WB131 reloads and checks zero new package
entitlements, then explicitly retries once. Only that successful retry creates
a payment. This scenario is always classified as controlled UI evidence.

`top_up_rapid` dispatches Playwright's real double-click with zero inter-click
delay. It does not replay a worker command or invoke a JavaScript click handler.
Matching in-flight requests must all yield identified responses before completion;
a short settling interval avoids dropping an immediate second request. Distinct
returned payment IDs are retained in the private receipt and make ownership
ambiguous. The C journal reports unresolved creation and retains resources for
manual reconciliation. Later responses cannot erase that ambiguity. Unknown or
late responses fail closed. The lab never automatically retries a payment or
deletes financial ledger entries.

## Coverage and limits

**77/132 fully automated (58.33%), 62/86 P0; live verified 0/132.** The gate remains
unmet. Fixture passes establish lab behavior, not original-console readiness.

| Requirement | State and remaining boundary |
|---|---|
| WEB-PAY-003 | Partial. Amount, plan, method, completed status, payment ID, persistence and paid-date checks exist. The source hardcodes `Walk-in customer` and shows a currency symbol without an ISO code, so correct payer and unambiguous currency reconciliation remain gaps. |
| WEB-PAY-004 | Partial. No selected plan and no SIM are covered. A selected plan becoming stale/deleted is untested; app state is not modified to manufacture it. |
| WEB-PAY-005 | Partial. A visible rejected submission, no entitlement and explicit recovery are covered. An actual backend ledger row returned with `failed`/`pending` status remains untested. The source treats every `onCompleted` payment as successful. |
| WEB-PAY-006 | Implemented for rapid clicks with owned response correlation and independent visible entitlement/receipt checks. Command replay protection remains a separate test. |
| WEB-PAY-007 | Implemented for actual PDF-to-visible-dialog equality. It does not imply the dialog itself has the right payer or currency. |

The test fixture generates a small real PDF with known values; it does not run
the original jsPDF export or prove its layout. The sample render was visually
reviewed for legibility. The unchanged source, live services, Firefox and WebKit
still require qualification. Downloads and traces may contain customer and
payment data; retain the existing private artifact handling rules.
