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

The existing scenario check names are kept temporarily for compatibility, but
their implementations now read direct resources, payments, KPI values and
performance reports. Scenario/check renaming is handled separately from this
BFF-access patch.
