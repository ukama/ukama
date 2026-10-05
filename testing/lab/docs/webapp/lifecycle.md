# Plan, SIM and customer lifecycle — Patch 12

Apply after lab Patch 11. This patch changes ukama-lab only. It does not change
console source or apply the old console companion archives from Patches 4/5.

## Local execution

Rebuild the C CLI and run `npm --prefix adapters/webapp ci` followed by
`npm --prefix adapters/webapp run build`. Use the existing local-stack setup in
`commerce.md` and your saved owner session. Run individual files under
`scenarios/webapp/p0/lifecycle/` using the same `ukama-lab run` options as existing
commerce scenarios. WB031 and WB042 in `expanded/` are strengthened too.

Use a disposable local organization with a fully visible pool (at most 100 rows).
WB121 requires the sole eligible SIM option to be the run-owned Factory ICCID;
foreign eligible inventory aborts before allocation. WB126 requires an actually
empty pool. These are explicit prerequisites; the lab never deletes foreign
inventory or fabricates pool data to meet them.

## Actions and checks

`web_commerce` adds `rename_plan` (owned name plus `-renamed`) and `allocate_auto`.
Both use existing world package/UE/network references. Rename retains the
immutable creation name for cleanup, while later UI checks resolve the new name.
Reload verifies persistence and all four commercial edit fields must be read-only.
Auto-allocation observes its mutation response only to record ownership; ICCID,
plan, allowance and validity acceptance come from the rendered drawer.

`name_pending` and `name_failure` require a prior observed name query and affect
only the exact planned name at that unique endpoint. `pool_failure` affects the
previously observed SIM type and limit at the unique SimPoolOverview endpoint.
`clear_fault` removes the route; pending reads release on clear/close. Fault checks
must report `*:applied`; navigation resets applied proof. These scenarios are
classified controlled UI evidence even if a report is labelled live.

Pool reconciliation compares unique rendered ICCID/status rows with Available,
Assigned and Faulty KPI counts. The footer must prove the full inventory is
visible; filtered/capped/duplicate/error states cannot pass. Status filtering is
an ordinary UI click. `SIM option present` checks the exact ICCID and option value
in the allocation dialog. `Customer plan` checks the exact named row across all
three customer lenses after reload. WB031 checks supported validity options as
well as the existing field validation errors.

## Coverage and retained gaps

**75/132 fully automated (56.82%), 61/86 P0; live verified 0/132.** The gate remains
unmet. Fixture success qualifies lab mechanics, not the original React/Apollo app.
The table records every requirement assigned to this patch.

| Requirement | State and boundary |
|---|---|
| WEB-PLAN-007 | Partial: wb-124/125 propose fail-closed behavior for pending/failed name reads. The unchanged source permits submission in both states. The agreed policy still needs confirmation; no complete credit. |
| WEB-PLAN-008 | Partial: wb-031 adds the closed supported validity choices to name/price/allowance errors. Unsupported validity tampering is not injected into app state. |
| WEB-PLAN-011 | wb-120 renames an owned plan, reloads, verifies terms/price and all four read-only edit fields; creation ownership name is retained separately from the acknowledged UI name. |
| WEB-SIM-002 | wb-122/123/126 reconcile exact unique visible ICCID/status rows against all three KPI counts. Capped or filtered lists fail reconciliation. Worker fixtures also exercise a nonzero faulty count. |
| WEB-SIM-003 | wb-122 checks the allocated ICCID is absent and the other owned ICCID remains selectable, then verifies Assigned/Available after reload. |
| WEB-SIM-004 | Gap: server malformed-CSV/duplicate-import policy is not agreed; no unsafe duplicate Factory upload or fabricated inventory is used. |
| WEB-SIM-005 | Partial: wb-122 covers available/assigned filter membership. Original SIM pool has no search or pagination. |
| WEB-SIM-006 | Partial: wb-126 requires an actually empty inventory and checks the blocked Add customer allocation entry with its visible guidance. Existing no-SIM drawer allocation with an exhausted pool remains a gap. |
| WEB-SIM-007 | wb-123 injects an exact observed endpoint SimPoolOverview section error, proves the injection applied, rejects a successful empty state and checks restored-read recovery. Controlled evidence only. |
| WEB-CUSTOMER-001 | Partial: UI-created names persist; the source renders phone, not email, in list/detail. Name/email correlation receipts establish ownership only, not visible email acceptance. |
| WEB-CUSTOMER-003 | wb-121 uses Auto-assign from pool only when the sole eligible option is the run-owned ICCID; checks ICCID, current plan, allowance, validity and reload. |
| WEB-CUSTOMER-004 | Gap: source retries addSubscriber after allocation failure. Safe qualification of the combined flow needs staged ownership of every subscriber/possible SIM; separate create/allocate is not claimed as this test. |
| WEB-CUSTOMER-005 | Gap: source flattens a single SIM and has no independently selectable multi-SIM drawer. |
| WEB-CUSTOMER-006 | Partial: wb-023 verifies visible service off/on. UE traffic denial and recovery across the service toggle remain unimplemented. |
| WEB-CUSTOMER-007 | Partial: usage/allowance/validity checks remain; no remaining-data display exists and the drawer clamps unknown cycle usage to zero. |
| WEB-CUSTOMER-008 | Partial: existing search/filter/sort membership; no source pagination control. |
| WEB-CUSTOMER-009 | wb-042 now checks the owned customer row and exact world-derived plan name in Business, Network and Customer lists after reload. |
| WEB-CUSTOMER-010 | Partial: wb-127 adds persisted no-SIM/No plan drawer actions to missing-plan creation guidance. Assigned SIM without plan and suspended variants remain gaps. |

The name-read cases intentionally express proposed fail-closed behavior. The
unchanged console only blocks a name marked taken and permits pending/error
states. These tests expose that behavior; they do not establish an agreed policy.
No artificial app state, API-only acceptance assertion, or console repair is used.
Private traces and ownership receipts may contain customer/ICCID data. Preserve
the existing artifact handling and cleanup rules in `commerce.md`.
