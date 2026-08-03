# Patch 2: reusable console checks

Prerequisite: apply Patch 1 (`modern BFF access`) first.

This patch adds reusable BFF-backed checks for console acceptance scenarios.
It does not add or change smoke/P0 scenario YAML and it does not add new
action events. Those remain for later patches.

## Added checks

- `list_count_equals`
- `entity_fields_equal`
- `entity_reconciles`
- `node_status_equals`
- `software_status_equals`
- `software_count_equals`
- `node_operation_status_equals`
- `site_operation_status_equals`
- `kpi_state_equals`
- `kpi_timeseries`
- `package_fields_equal`

`package_catalog_equals` remains accepted as a compatibility alias for
existing scenarios, but the canonical name is now `package_fields_equal`.

## BFF operations used

- `getNetworks`, `getNetwork`
- `getSites`, `getSite`
- `getNodes`, `getNode`
- `getSubscribersByNetwork`, `getSubscriber`
- `getPackages`, `getPackage`
- `getSimsByNetwork`
- `getSoftwares`
- `getNodeOperationStatus`
- `getSiteOperationStatus`
- `getKpiValues`
- `getKpiTimeSeries`

The checks do not call analytics services directly.

## Important semantics

- `list_count_equals` is scoped to exactly one selected network, except when
  counting networks themselves.
- `node_status_equals` validates connectivity and lifecycle state as separate
  fields.
- `software_status_equals` validates the console-facing `getSoftwares` state.
  Runtime version and health should still use `node_version_equals` and
  `node_health_ok`.
- `kpi_state_equals` distinguishes missing/unavailable, no-data, present,
  zero/non-zero, and fresh/stale values.
- `kpi_timeseries` validates bucket order, requested bounds, bucket count,
  aggregate value, and computed timestamps without testing chart pixels.

## Install as drop-in replacements

From the ukama-lab repository root after Patch 1:

```sh
cp -a patch-2-reusable-console-checks/docs/. docs/
cp -a patch-2-reusable-console-checks/inc/. inc/
cp -a patch-2-reusable-console-checks/src/. src/
```

Alternatively, apply the separately supplied unified patch with:

```sh
patch -p1 < ukama-lab-patch-2-reusable-console-checks.patch
```

## Validation performed

- Modified C sources pass strict syntax compilation with `-Wall -Wextra
  -Werror -Wdeclaration-after-statement`.
- All 101 existing smoke/P0 YAML files parse unchanged.
- A parser fixture containing every new check name parses successfully.
- The generated unified patch applies cleanly on top of Patch 1.
- No retired console-view GraphQL operation was reintroduced.

A complete project link was not possible from the standalone uploaded archive
because its Makefile depends on the external `nodes/ukamaOS` build tree. A live
udev run was also not used while the current `/get-user` authentication path is
blocked by the Init/Registry lookup error discussed separately.

Suggested commit:

```text
Add reusable BFF checks for console acceptance scenarios
```
