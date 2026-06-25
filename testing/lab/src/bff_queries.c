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
" addPackage(data: $data) { uuid name dataVolume dataUnit duration amount } }";

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

const char *BFF_GET_SIM_PACKAGES =
"query GetPackagesForSim($data: GetPackagesForSimInputDto!) {"
" getPackagesForSim(data: $data) { sim_id packages { "
" package_id is_active } } }";

const char *BFF_GET_NODE_STATE =
"query GetNodeState($id: String!) {"
" getNodeState(id: $id) { state connectivity } }";

const char *BFF_NETWORK_OVERVIEW =
"query NetworkOverview($networkId: String!) {"
" networkOverview(networkId: $networkId) {"
" network { network { id name } error } sites { sites { id name } error } } }";

const char *BFF_SITE_VIEW =
"query SiteView($siteId: String!) {"
" siteView(siteId: $siteId) { site { site { id name } error } "
" nodes { nodes { id name status { state connectivity } } error } } }";

const char *BFF_GET_NETWORKS =
"query GetNetworks { getNetworks { networks { id name } } }";

const char *BFF_GET_SITES =
"query GetSites($networkId: String!) { getSites(networkId: $networkId) {"
" sites { id name } } }";

const char *BFF_GET_NODES_FOR_SITE =
"query GetNodesForSite($siteId: String!) { getNodesForSite(siteId: $siteId) {"
" nodes { id name type } } }";

const char *BFF_GET_COMPONENTS_BY_USER_ID =
"query GetComponentsByUserId($data: ComponentTypeInputDto!) {"
" getComponentsByUserId(data: $data) {"
" components { id category type description partNumber } } }";

const char *BFF_GET_NODES =
"query GetNodes($data: NodesFilterInput!) {"
" getNodes(data: $data) { nodes { id name type latitude longitude "
" site { nodeId siteId networkId addedAt } "
" status { state connectivity } } } }";
