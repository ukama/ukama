# BFF GraphQL operations

The BFF client is based on resolver/test definitions from Console BFF:

Mutations:

- `addNetwork(data: AddNetworkInputDto!)`
- `addSite(data: AddSiteInputDto!)`
- `addNode(data: AddNodeInput!)`
- `addNodeToSite(data: AddNodeToSiteInput!)`
- `addPackage(data: AddPackageInputDto!)`
- `addSubscriber(data: SubscriberInputDto!)`
- `allocateSim(data: AllocateSimInputDto!)`

Queries:

- `getDataUsage(simId: String!)`
- `getPackagesForSim(data: GetPackagesForSimInputDto!)`
- `getNodeState(nodeId: String!)`
- `networkOverview(networkId: String!)`
- `siteView(siteId: String!)`
