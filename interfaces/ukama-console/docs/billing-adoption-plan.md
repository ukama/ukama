# Billing hookup — adopting the legacy console flow on ukama-console

Status: plan / not yet implemented
Scope agreed: Phases A + B + C (full Stripe pay flow), using the **direct `GetPayments`** integration pattern.

## Context

The legacy console (`interfaces/console/src/app/manage/billing/page.tsx`) already
implements a complete invoice **read + pay + PDF** flow. The new console
(`interfaces/ukama-console`) renders Billing at
`src/app/(dashboard)/business/manage/_components/BillingScreen.tsx`, currently
wired only for reading invoices/balance through the `commerceView` composite,
with the Pay / Update payment / PDF actions still stubbed as toasts.

Both apps talk to the **same** backend-for-frontend, `systems/console-bff`. That
BFF already exposes everything this flow needs — no Go / `systems/billing`
changes are required for any phase below.

### What already exists (no work needed)

- **console-bff `payment` module**: `getPayment`, `getPayments`, `updatePayment`
  (mutation), `processPayment` (mutation), `getToken`.
- **console-bff `report` module**: `getGeneratedPdfReport` (returns
  `{ contentType, filename, downloadUrl }`, where `downloadUrl` is base64).
- **console-bff `billing` module**: `getReports` / `getReport` with full
  `rawReport` (amounts, `fileUrl`, subscriptions, fees).
- **ukama-console GraphQL ops**: `src/client/graphql/billing.graphql` already
  contains, lifted verbatim from the legacy console, the `payment` fragment,
  `UpdatePayment`, `ProcessPayment`, `GetPayment`, `GetPayments`,
  `GetReports`/`GetReport`, and `getGeneratedPdfReport`. Hooks are generated in
  `src/client/graphql/billing.generated.ts`:
  `useGetPaymentsQuery`, `useUpdatePaymentMutation`, `useProcessPaymentMutation`,
  `useGetReportsQuery`, `useGetGeneratedPdfReportLazyQuery`.

### The real gaps

1. No on-demand PDF wiring (current button uses `rawReport.fileUrl` directly).
2. No payment-state load (`getPayments`) and no per-invoice `extra` (Stripe
   client secret) plumbing.
3. **Stripe is not installed** — `@stripe/stripe-js` and
   `@stripe/react-stripe-js` are absent; there is no `STRIPE_PK` env var and no
   ported `StripePaymentDialog`.

## How the legacy pay flow works (the behaviour we are adopting)

1. Load invoices (`getReports`, `report_type: 'invoice'`) and payments
   (`getPayments`, `type: 'invoice'`) in parallel.
2. Match each invoice to its payment by `payment.itemId === report.id`.
3. On Pay:
   - If the matched payment already has `extra` (the Stripe **client secret**),
     open the Stripe dialog immediately.
   - Otherwise call `updatePayment({ id, paymentMethod: 'stripe', payerEmail,
     payerName })`. The payment service provisions a Stripe PaymentIntent and
     returns it in `extra`; refetch payments, read the new `extra`, open dialog.
4. `StripePaymentDialog` mounts `<Elements clientSecret={extra}>` +
   `<PaymentElement>` and calls `stripe.confirmPayment({ redirect: 'if_required' })`.
   On `paymentIntent.status === 'succeeded'`, refetch reports + payments, close.
5. PDF: `getGeneratedPdfReport(id)` → base64 `downloadUrl` →
   `base64ToBlob` → temporary `<a download>` click.

`payerEmail` / `payerName` in the new app come from `useAuth()`
(`src/lib/auth/context.tsx`, returns `AuthUser` with `name`, `email`,
`currency`).

---

## Phase A — On-demand PDF (small, no new deps)

Replace the provider-URL PDF button with the legacy on-demand route (more
reliable than depending on `rawReport.fileUrl` being populated).

Changes
- Add a `base64ToBlob` helper (ukama-console has none today). Suggested:
  `src/lib/download.ts` exporting `base64ToBlob(dataUrl, mime)` and a small
  `downloadBlob(blob, filename)`.
- In `BillingScreen.tsx`: use `useGetGeneratedPdfReportLazyQuery` with
  `onCompleted` → build blob from `downloadUrl` → download as `filename`;
  `onError` → toast. Track a `downloadingId` for per-row spinner/disabled state.
- Keep `rawReport.fileUrl` as a fallback only if `getGeneratedPdfReport` errors
  (optional).

Acceptance criteria
- Clicking PDF on any invoice downloads a correctly named `.pdf`.
- The clicked row shows a loading state until the download resolves.
- Errors surface as a toast; no unhandled promise rejections.

## Phase B — Payment state via direct `GetPayments` (small/medium)

Adopt the legacy data shape using the chosen **direct query** pattern (no BFF
change).

Changes
- In `BillingScreen.tsx`, add `useGetPaymentsQuery({ variables: { data: { type:
  'invoice' } }, fetchPolicy: 'network-only' })` alongside the existing
  `useBillingOverviewQuery`.
- Build a lookup `paymentByItemId = Map(payment.itemId -> payment)`.
- Derive real per-invoice pay state (e.g. `status`, `failureReason`) and gate
  the Pay button on having a matching payment record.
- Note: invoices continue to render via `commerceView`; `getPayments` only
  supplies payment records. (Amounts already shipped via `commerceView` in the
  previous change — `rawReport.totalAmountCents` + balance `outstandingAmount`.)

Acceptance criteria
- Each unpaid invoice maps to its payment record; rows with no payment record do
  not show an actionable Pay button.
- No duplicate/over-fetch: `getPayments` runs once per network scope.

## Phase C — Stripe pay flow (substantial; new deps + env + UI)

Changes
1. **Dependencies**: add `@stripe/stripe-js` and `@stripe/react-stripe-js` to
   `interfaces/ukama-console/package.json`.
2. **Env**: add `NEXT_PUBLIC_STRIPE_PK` to the zod schema in `src/env.ts` (and to
   `src/lib/runtime-env.ts` if runtime injection is needed for Docker/CI, mirroring
   the existing `pick()` pattern). Document the var in the deployment env.
3. **Dialog**: port `StripePaymentDialog` into the new design system at
   `src/app/(dashboard)/business/manage/_components/StripePaymentDialog.tsx`
   (or `src/components/` if shared). Reuse the legacy logic:
   `loadStripe(env.NEXT_PUBLIC_STRIPE_PK)`, `<Elements clientSecret={extra}>`,
   `<PaymentElement>`, `confirmPayment({ redirect: 'if_required' })`. Restyle
   buttons/typography to ukama-console tokens; show the amount via
   `useCurrency().money(Number(totalAmountCents)/100)` instead of a hardcoded `$`.
4. **Orchestration in `BillingScreen.tsx`**: port `handleAddPayment(billId)`:
   - find matching payment; if `extra` present → open dialog;
   - else `updatePayment({ data: { id, paymentMethod: 'stripe', payerEmail:
     user.email, payerName: user.name } })` → refetch payments → read new
     `extra` → open dialog;
   - on success → refetch `useBillingOverviewQuery` + `getPayments`, close dialog,
     clear state, success toast.
   Wire it to "Pay now" (balance card), the per-row pay action, and decide
   whether "Update payment" reuses the same path.

Acceptance criteria
- With a valid `STRIPE_PK`, clicking Pay opens the Stripe dialog pre-loaded with
  the invoice's client secret and correct amount.
- A successful test-card payment marks the invoice paid after refetch and closes
  the dialog.
- Card errors / cancellations surface cleanly and leave the invoice unpaid.
- Build + typecheck pass for `ukama-console` (`tsc --noEmit`) and the BFF is
  unchanged.

---

## Integration decision (recorded)

- **Payment data**: direct `GetPayments` query in the component (mirrors legacy
  console exactly; zero BFF change). The alternative — extending `commerceView`
  with a `payments`/`extra` section to keep a single composite query — was
  considered and deferred.

## Out of scope / not required

- No changes to `systems/billing` (Go) or to `console-bff` for any phase.
- No codegen schema change (all operations already exist in `billing.graphql`);
  run `pnpm codegen` against a live gateway only if hooks need refreshing.

## Suggested sequencing

A and B first (self-contained, no new dependencies), then C as its own change
since it introduces Stripe deps, env configuration, and a payment UI to
re-style.
