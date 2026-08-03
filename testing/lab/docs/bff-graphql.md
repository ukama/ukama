# BFF GraphQL operations

Ukama Lab talks only to Console BFF. The client uses the same direct GraphQL
operations that back the current console screens; it does not use the retired
composite dashboard views.

Mutations:

- `addNetwork(data: AddNetworkInputDto!)`
- `addSite(data: AddSiteInputDto!)`
- `addNode(data: AddNodeInput!)`
- `addNodeToSite(data: AddNodeToSiteInput!)`
- `addPackage(data: AddPackageInputDto!)`
- `updatePackage(packageId: String!, data: UpdatePackageInputDto!)`
- `addSubscriber(data: SubscriberInputDto!)`
- `allocateSim(data: AllocateSimInputDto!)`
- `addPayment(data: AddPaymentInputDto!)`
- `updateSoftware(data: UpdateSoftwareInputDto!)`

Direct resource queries:

- `getNetwork(networkId: String!)`
- `getNetworks`
- `getSites(data: SitesInputDto!)`
- `getSite(siteId: String!)`
- `getNodes(data: NodesFilterInput!)`
- `getNode(data: NodeInput!)`
- `getNodesForSite(siteId: String!)`
- `getSubscribersByNetwork(networkId: String!)`
- `getSubscriber(subscriberId: String!)`
- `getSims(data: ListSimsInput!)`
- `getSimsByNetwork(networkId: String!)`
- `getPackage(packageId: String!)`
- `getPackages(networkId: String)`
- `isPackageNameAvailable(name: String!)`
- `getPackagesForSim(data: GetPackagesForSimInputDto!)`
- `getSimsUsageByNetwork(networkId: String!)`
- `getPayments(data: GetPaymentsInputDto!)`
- `getApps(data: GetAppsInputDto!)`

Analytics queries:

- `getKpiValues(data: KpiValuesInput!)`
- `getPerformanceReport(data: PerformanceReportInput!)`

Release queries:

- `getReleaseCatalog(data: GetReleaseCatalogInput!)`

## Two SIM list queries — `getSims` is preferred

`getSims(data: {networkId, status})` proxies subscriber `/v1/sim` with query
filters. `getSimsByNetwork` proxies `/v1/sims/networks/{network_id}`, which is
marked *"Deprecated: Use ListSims with networkId as filtering param instead"*
in both `systems/subscriber/api-gateway/pkg/rest/router.go` and
`sim-manager/pkg/db/sim_repo.go`. Both return the same rows.

**Prefer `getSims` for new call sites.** `getSimsByNetwork` is retained and
still used by `bff_get_package_metrics`, so both queries live in
`src/bff_queries.c`. `list_count_equals target: sims` uses `getSims`.

`status` is a non-null `String` on `ListSimsInput` and must be sent. The empty
string parses to `SimStatusUnknown` (`ukama/sim_status.go`), which the
sim-manager repo treats as "no status filter" — every SIM in the network.

The two queries expose **different package types**: `getSims` returns
`SimPackage` with camelCase `packageId`/`isActive`, while `getSimsByNetwork`
returned `SimPackageDto` with snake_case `package_id`/`is_active`. Any code
reading a SIM's package must match the query it used.

Subscriber lists still use `getSubscribersByNetwork(networkId:)` — that route
(`/v1/subscribers/networks/{network_id}`) is not deprecated.

The existing scenario check names are kept temporarily for compatibility, but
their implementations now read direct resources, payments, KPI values and
performance reports. Scenario/check renaming is handled separately from this
BFF-access patch.
