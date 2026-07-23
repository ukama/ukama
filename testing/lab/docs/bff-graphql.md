# BFF GraphQL operations

The BFF client is based on resolver/test definitions from Console BFF:

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

Queries:

- `getPackage(packageId: String!)`
- `getPackages(networkId: String)`
- `isPackageNameAvailable(name: String!)`
- `packagesDashboard(networkId: String)`
- `getSimsUsageByNetwork(networkId: String!)`
- `getPackagesForSim(data: GetPackagesForSimInputDto!)`
- `getPayments(data: GetPaymentsInputDto!)`
- `getKpiValues(data: KpiValuesInput!)`
- `getPerformanceReport(data: PerformanceReportInput!)`
- `getNode(data: NodeInput!) { id status { connectivity state } }`
- `networkOverview(networkId: String!)`
- `siteView(siteId: String!)`

The P0 data-package scenarios communicate only with these BFF GraphQL
operations. Initial allocation billing is observed indirectly through
`packagesDashboard` revenue and attachment count. Subsequent cash sales use
`addPayment`; entitlement creation, ordering, and transitions are observed
through `getPackagesForSim`. No scenario calls package, payment, subscriber,
or analytics backend services directly.
