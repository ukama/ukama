# Patch 14 — analytics observations and reconciliation

Lab-only; apply after Patch 13. The console is unchanged. All acceptance values
come from rendered UI. Service payloads and React internals are never an oracle.
The worker requires the exact preceding `web_open`, including its bound detail
path and selected network, before every analytics action/assertion.

## Scenarios

| Case | Observation | Prerequisite |
| --- | --- | --- |
| WB140 | Node rail units/missing value, threshold legend, disconnected SVG subpaths, sampled tooltip, missing-sample tooltip and selected range | Controlled metric profile below |
| WB141 | Site component selection, battery tooltip and failed backhaul metric | Controlled metric profile below |
| WB142 | Two selected app identities, CPU, RSS, disk values and missing value | Controlled app profile below |
| WB143 | Revenue tooltip date and amount | Controlled revenue series below |
| WB144 | Two customers with a $25 and $10 purchase; each plan's sold/revenue/share and business customer/revenue totals | Fresh owned network, local commerce services |
| WB145 | Explicit zero revenue/purchases/package sales in a fresh network | Local analytics collector caught up |

WB144 creates both customers/plans/SIMs through the existing browser lifecycle,
uses identified payment ownership, and retains the existing cleanup rules.
WB145's zero expectation is not evidence that an unavailable service is zero.

## Metric profile and limits

WB140–143 are read-only qualification examples requiring an independently
prepared local telemetry/analytics environment. They **do not seed telemetry**,
freeze the browser clock, insert application state or fabricate live data.
Prepare observations for the run-owned node/site before the check, or adapt the
expectations to an independently recorded local profile. Their fixed sample
clock is illustrative and must fall within the selected live date window.
These cases remain partial automation until a reproducible local producer is
connected; a passing DOM fixture is never product verification.

The UTC example has samples at 2026-10-05 12:00, 12:15, 12:30, 12:45, 13:00.
The middle sample is missing; the others are real. At 12:15 the temperature is
55.50 °C and battery charge 55.50%; decimal precision is two places. Thresholds
are normal below 60, high up to 80, critical at/above 80. Memory is unavailable;
backhaul's range read fails. The corresponding controller metric label is
`Backhaul downlink`. The node rail's current temperature is 55.50 °C.
The resource dialog profile uses apps `metrics` and `controller`, CPU 12.5%,
RSS 64.0 MB (binary units), disk read 2.0 KB and disk write unavailable. Revenue
at Oct 5 is $25.5 after conversion from 2550 cents. A live app with an all-or-none
resource object may not expose a single missing disk field: that remains a gap.

`chart_hover` takes a percentage along the **rendered SVG plot**, not a sample
index. For nonuniform timestamps, calculate the percentage independently from
the known time domain. It refuses ambiguous clipping geometry. It moves the
pointer away before sampling and requires a new same-page/card hover before a
tooltip check. `Chart segments` counts rendered SVG move subpaths; combine it
with a missing-point hover, since a segment count alone cannot locate a gap.

## Commands

`web_interact` adds `chart_hover` (label=chart title, value=0..100),
`chart_range` (label=chart title, value=Day|Week|Month), and `metric_select`:
label=tab|section|component|app|dialog, value=the visible choice (dialog=Close).
These are closed interactions; no arbitrary selector, script, URL or request.

`web_ui_equals` adds `Chart tooltip`, `Chart x axis`, `Chart y axis`,
`Chart legend`, `Chart segments`, `Chart state`, `Chart range`, `Metric value`,
`App resource`, and `App identity`. Use subject=chart/field title except for
App identity. Axis/legend values join visible labels with `|`. States are data,
empty or error; absent/hidden/loading/ambiguous content cannot count as zero.
Tooltip strings normalize whitespace only. Exact formatted expectations encode
units and rounding; there is no permissive numeric tolerance or unit stripping.

## Remaining requirements (frozen 132 denominator)

| IDs | Work added / remaining gap |
| --- | --- |
| WEB-SITE-007/008 | Component/error, units, legend/tooltip observers; deterministic live metric producer and broader components are still required. |
| WEB-NODE-006 | Temperature/memory observers; every exposed node type/category and reproducible missing-data stimulus remain. |
| WEB-NODE-008 | App identity/resource values observed; the source has a value dialog and no app resource time-series charts. |
| WEB-BIZ-002 | Two customer/plan totals; site totals, cross-network isolation and full rollup reconciliation remain. |
| WEB-BIZ-003 | Nontrivial $25/$10 distribution with 71%/29% shares; source has no per-customer revenue breakdown. |
| WEB-BIZ-004 | Existing WB039 checks rolling selectors/current purchases; revenue trend and package bars remain hardcoded to 30 days in source. Historical exclusion is not tested. |
| WEB-BIZ-005 | Visible SVG axes/tooltip observer; no full revenue legend contract or live series producer. |
| WEB-BIZ-006 | Explicit zero case added. Home KPI errors and missing values both degrade to dashes; separate failed-read presentation remains incomplete. |
| WEB-BIZ-007 | Two correct plan rows/totals; date and cross-network attribution remain partial. |
| WEB-UI-003 | Exact unit/rounding assertions for selected metrics/resources; all display families and agreed tolerance remain open. |
| WEB-UI-004 | Online sample gap probes; MetricChartCard's `off` branch still replaces every sample with zero. |
| WEB-UI-005 | UTC tooltip and selected range observed. Worker timezone is fixed UTC; non-UTC/DST and historical filter boundaries remain. |

No assigned requirement is promoted to complete by this patch. Coverage stays
77/132 implemented, 62/86 P0, 0 live verified. Patch 15 handles members/support;
Patch 16/17 retain the unclosed gaps. Do not change the frozen denominator.
