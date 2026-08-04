/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

const char *BFF_ADD_NETWORK =
"mutation AddNetwork($data: AddNetworkInputDto!) {"
" addNetwork(data: $data) { id name countries networks budget } }";

const char *BFF_ADD_SITE =
"mutation AddSite($data: AddSiteInputDto!) {"
" addSite(data: $data) { id name networkId backhaulId powerId "
" accessId spectrumId switchId isDeactivated latitude longitude "
" installDate createdAt location } }";

const char *BFF_ADD_NODE =
"mutation AddNode($data: AddNodeInput!) {"
" addNode(data: $data) { id name orgId type status { state connectivity } } }";

const char *BFF_ADD_NODE_TO_SITE =
"mutation AddNodeToSite($data: AddNodeToSiteInput!) {"
" addNodeToSite(data: $data) { success } }";

const char *BFF_ADD_PACKAGE =
"mutation AddPackage($data: AddPackageInputDto!) {"
" addPackage(data: $data) { uuid name active dataVolume dataUnit duration "
" amount currency country } }";

const char *BFF_UPDATE_PACKAGE =
"mutation UpdatePackage($packageId: String!, $data: UpdatePackageInputDto!) {"
" updatePackage(packageId: $packageId, data: $data) {"
" uuid name active duration amount currency country } }";

const char *BFF_GET_PACKAGE =
"query GetPackage($packageId: String!) {"
" getPackage(packageId: $packageId) { uuid name active duration dataVolume "
" dataUnit amount currency country networkId } }";

const char *BFF_GET_PACKAGES =
"query GetPackages($networkId: String) {"
" getPackages(networkId: $networkId) { packages { uuid name active "
" networkId } } }";

const char *BFF_PACKAGE_NAME_AVAILABLE =
"query IsPackageNameAvailable($name: String!) {"
" isPackageNameAvailable(name: $name) { isAvailable name } }";

const char *BFF_GET_NETWORK =
"query GetNetwork($networkId: String!) {"
" getNetwork(networkId: $networkId) { id name } }";

const char *BFF_GET_NETWORKS =
"query GetNetworks { getNetworks { networks { id name } } }";

const char *BFF_GET_SITES =
"query GetSites($data: SitesInputDto!) { getSites(data: $data) {"
" sites { id name networkId latitude longitude isDeactivated "
" installDate createdAt location } } }";

const char *BFF_GET_SITE =
"query GetSite($siteId: String!) { getSite(siteId: $siteId) {"
" id name networkId latitude longitude isDeactivated "
" installDate createdAt location } }";

const char *BFF_GET_NODES =
"query GetNodes($data: NodesFilterInput!) {"
" getNodes(data: $data) { nodes { id name type latitude longitude "
" site { nodeId siteId networkId addedAt } "
" status { state connectivity } } } }";

const char *BFF_GET_NODES_FOR_SITE =
"query GetNodesForSite($siteId: String!) { getNodesForSite(siteId: $siteId) {"
" nodes { id name type latitude longitude status { state connectivity } } } }";

const char *BFF_GET_SUBSCRIBERS_BY_NETWORK =
"query GetSubscribersByNetwork($networkId: String!) {"
" getSubscribersByNetwork(networkId: $networkId) { subscribers {"
" uuid name email phone networkId sim { id status networkId package {"
" package_id is_active } } } } }";

const char *BFF_GET_SUBSCRIBER =
"query GetSubscriber($subscriberId: String!) {"
" getSubscriber(subscriberId: $subscriberId) {"
" uuid name email phone networkId sim { id status networkId package {"
" package_id is_active } } } }";

/*
 * Network-scoped SIM list, deprecated upstream but retained.
 *
 * Proxies subscriber /v1/sims/networks/{network_id}, marked "Deprecated: Use
 * ListSims with networkId as filtering param instead" in both
 * systems/subscriber/api-gateway/pkg/rest/router.go and
 * sim-manager/pkg/db/sim_repo.go. Kept because it returns SubscriberSimDto,
 * whose package is SimPackageDto with snake_case package_id/is_active --
 * a different GraphQL type from what getSims below returns. Prefer
 * BFF_GET_SIMS for new call sites.
 */
const char *BFF_GET_SIMS_BY_NETWORK =
"query GetSimsByNetwork($networkId: String!) {"
" getSimsByNetwork(networkId: $networkId) { sims {"
" id subscriberId networkId status package { package_id is_active } } } }";

/*
 * Network-scoped SIM list, preferred.
 *
 * Proxies subscriber /v1/sim with query filters, the supported path. Returns
 * the same rows as BFF_GET_SIMS_BY_NETWORK above.
 *
 * status is a non-null String on ListSimsInput and must be sent. An empty
 * string parses to SimStatusUnknown (ukama/sim_status.go ParseSimStatus),
 * which the repo treats as "no status filter", i.e. every SIM.
 *
 * Note the package field names: getSims returns SimPackage (camelCase
 * packageId/isActive), whereas getSimsByNetwork returned SimPackageDto
 * (snake_case package_id/is_active). They are different GraphQL types.
 */
const char *BFF_GET_SIMS =
"query GetSims($data: ListSimsInput!) {"
" getSims(data: $data) { sims {"
" id subscriberId networkId status package { packageId isActive } } } }";

const char *BFF_INVENTORY_OVERVIEW =
"query InventoryOverview { inventoryView { components { total byCategory {"
" category count } } simStock { total available consumed } } }";

const char *BFF_SIM_POOL_OVERVIEW =
"query SimPoolOverview($simType: String!, $limit: Int!) {"
" simPoolView(simType: $simType) { stats { total available consumed } "
" sims(limit: $limit) { sims { id } } } }";

const char *BFF_ADD_SUBSCRIBER =
"mutation AddSubscriber($data: SubscriberInputDto!) {"
" addSubscriber(data: $data) { uuid email name networkId phone } }";

const char *BFF_ALLOCATE_SIM =
"mutation AllocateSim($data: AllocateSimInputDto!) {"
" allocateSim(data: $data) { id subscriber_id network_id iccid imsi status "
" package { packageId isActive startDate endDate } } }";

const char *BFF_GET_DATA_USAGE =
"query GetDataUsage($data: SimUsageInputDto!) {"
" getDataUsage(data: $data) { simId usage } }";

const char *BFF_GET_SIMS_USAGE_BY_NETWORK =
"query GetSimsUsageByNetwork($networkId: String!) {"
" getSimsUsageByNetwork(networkId: $networkId) { simId usage } }";

const char *BFF_GET_SIM_PACKAGES =
"query GetPackagesForSim($data: GetPackagesForSimInputDto!) {"
" getPackagesForSim(data: $data) { sim_id packages { "
" id package_id start_date end_date is_active } } }";

const char *BFF_ADD_PAYMENT =
"mutation RecordCashPackageSale($data: AddPaymentInputDto!) {"
" addPayment(data: $data) { id itemId itemType amount currency paymentMethod "
" status paidAt payerEmail payerPhone metadata } }";

const char *BFF_GET_PAYMENTS =
"query GetPayments($data: GetPaymentsInputDto!) {"
" getPayments(data: $data) { payments { id itemId itemType amount currency "
" paymentMethod status paidAt payerEmail payerPhone metadata } } }";

const char *BFF_GET_KPI_VALUES =
"query GetKpiValues($data: KpiValuesInput!) {"
" getKpiValues(data: $data) { values { kpi value span op from to unit symbol "
" isPartial computedAt scope { key value } trend { direction changePct "
" changeAbs prevValue hasPrevious } } } }";

const char *BFF_GET_PERFORMANCE_REPORT =
"query GetPerformanceReport($data: PerformanceReportInput!) {"
" getPerformanceReport(data: $data) { report span rows { entityId status "
" attributes { key value } cells { column value unit symbol format } } } }";

const char *BFF_GET_NODE =
"query GetNode($data: NodeInput!) {"
" getNode(data: $data) { id name type site { nodeId siteId networkId addedAt } "
" status { connectivity state } } }";

const char *BFF_GET_RELEASE_CATALOG =
"query GetReleaseCatalog($name: String!, $type: String!) {"
" getReleaseCatalog(data: { name: $name type: $type }) {"
" releases { name type version available chunked desired uploadedAt } } }";

const char *BFF_PROMOTE_RELEASE =
"mutation PromoteRelease($name: String!, $type: String!, "
"$version: String!) {"
" promoteRelease(data: { name: $name type: $type version: $version }) {"
" desiredVersion message name } }";

const char *BFF_UPDATE_SOFTWARE =
"mutation UpdateSoftware($data: UpdateSoftwareInputDto!) {"
" updateSoftware(data: $data) { message } }";

const char *BFF_GET_APPS =
"query GetApps($data: GetAppsInputDto!) {"
" getApps(data: $data) { apps { name version tag status } } }";

const char *BFF_GET_SOFTWARES =
"query GetSoftwares($data: GetSoftwaresInput!) {"
" getSoftwares(data: $data) { software { id releaseDate nodeId status "
" currentVersion desiredVersion name createdAt updatedAt } } }";

const char *BFF_GET_NODE_OPERATION_STATUS =
"query GetNodeOperationStatus($nodeId: String!) {"
" getNodeOperationStatus(data: { nodeId: $nodeId }) { nodeId type busy "
" operation { id type status requestedBy startedAt leaseExpiresAt } } }";

const char *BFF_GET_SITE_OPERATION_STATUS =
"query GetSiteOperationStatus($siteId: String!) {"
" getSiteOperationStatus(data: { siteId: $siteId }) { siteId busy degraded "
" nodes { nodeId type busy operation { id type status requestedBy startedAt "
" leaseExpiresAt } } actions { restartSite { available reason } "
" rf { available reason } service { available reason } } } }";

const char *BFF_GET_KPI_TIMESERIES =
"query GetKpiTimeSeries($data: KpiTimeSeriesInput!) {"
" getKpiTimeSeries(data: $data) { values { kpi value span op from to unit "
" symbol isPartial computedAt scope { key value } trend { direction "
" changePct changeAbs prevValue hasPrevious } } } }";

const char *BFF_GET_COMPONENTS_BY_USER_ID =
"query GetComponentsByUserId($data: ComponentTypeInputDto!) {"
" getComponentsByUserId(data: $data) {"
" components { id category type description partNumber } } }";
