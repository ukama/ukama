import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTimeISO: { input: any; output: any; }
};

export type AddDraftInput = {
  lastSaved: Scalars['Float']['input'];
  name: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type AddMemberInputDto = {
  role: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type AddNetworkInputDto = {
  budget: Scalars['Float']['input'];
  countries?: InputMaybe<Array<Scalars['String']['input']>>;
  name: Scalars['String']['input'];
  networks?: InputMaybe<Array<Scalars['String']['input']>>;
  org: Scalars['String']['input'];
};

export type AddNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
};

export type AddNodeToSiteInput = {
  networkId: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export type AddPackageInputDto = {
  amount: Scalars['Float']['input'];
  dataUnit: Scalars['String']['input'];
  dataVolume: Scalars['Int']['input'];
  duration: Scalars['Int']['input'];
  name: Scalars['String']['input'];
};

export type AddPackageSimResDto = {
  __typename?: 'AddPackageSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type AddPackageToSimInputDto = {
  package_id: Scalars['String']['input'];
  sim_id: Scalars['String']['input'];
  start_date: Scalars['DateTimeISO']['input'];
};

export type AddSiteInputDto = {
  site: Scalars['String']['input'];
};

export type AllocateSimApiDto = {
  __typename?: 'AllocateSimAPIDto';
  activationsCount: Scalars['String']['output'];
  allocated_at: Scalars['String']['output'];
  deactivationsCount: Scalars['String']['output'];
  firstActivatedOn?: Maybe<Scalars['String']['output']>;
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi?: Maybe<Scalars['String']['output']>;
  is_physical: Scalars['Boolean']['output'];
  lastActivatedOn?: Maybe<Scalars['String']['output']>;
  msisdn: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
  package: SimAllocatePackageDto;
  status: Scalars['String']['output'];
  subscriber_id: Scalars['String']['output'];
  sync_status: Scalars['String']['output'];
  traffic_policy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type AllocateSimInputDto = {
  iccid?: InputMaybe<Scalars['String']['input']>;
  network_id: Scalars['String']['input'];
  package_id: Scalars['String']['input'];
  sim_type: Scalars['String']['input'];
  subscriber_id: Scalars['String']['input'];
  traffic_policy: Scalars['Float']['input'];
};

export type AppChangeLog = {
  __typename?: 'AppChangeLog';
  date: Scalars['Float']['output'];
  version: Scalars['String']['output'];
};

export type AppChangeLogs = {
  __typename?: 'AppChangeLogs';
  logs: Array<AppChangeLog>;
  type: NodeTypeEnum;
};

export type AttachNodeInput = {
  anodel: Scalars['String']['input'];
  anoder: Scalars['String']['input'];
  parentNode: Scalars['String']['input'];
};

export type AttachedNodes = {
  __typename?: 'AttachedNodes';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  site: NodeSite;
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type CBooleanResponse = {
  __typename?: 'CBooleanResponse';
  success: Scalars['Boolean']['output'];
};

export type CoverageInput = {
  height: Scalars['Float']['input'];
  lat: Scalars['Float']['input'];
  lng: Scalars['Float']['input'];
  mode: Scalars['String']['input'];
};

export type DefaultMarkupHistoryDto = {
  __typename?: 'DefaultMarkupHistoryDto';
  Markup: Scalars['Float']['output'];
  createdAt: Scalars['String']['output'];
  deletedAt: Scalars['String']['output'];
};

export type DefaultMarkupHistoryResDto = {
  __typename?: 'DefaultMarkupHistoryResDto';
  markupRates?: Maybe<Array<DefaultMarkupHistoryDto>>;
};

export type DefaultMarkupInputDto = {
  markup: Scalars['Float']['input'];
};

export type DefaultMarkupResDto = {
  __typename?: 'DefaultMarkupResDto';
  markup: Scalars['Float']['output'];
};

export type DeleteDraftRes = {
  __typename?: 'DeleteDraftRes';
  id: Scalars['String']['output'];
};

export type DeleteInvitationResDto = {
  __typename?: 'DeleteInvitationResDto';
  id: Scalars['String']['output'];
};

export type DeleteLinkRes = {
  __typename?: 'DeleteLinkRes';
  id: Scalars['String']['output'];
};

export type DeleteNode = {
  __typename?: 'DeleteNode';
  id: Scalars['String']['output'];
};

export type DeleteSimInputDto = {
  simId: Scalars['String']['input'];
};

export type DeleteSimResDto = {
  __typename?: 'DeleteSimResDto';
  simId?: Maybe<Scalars['String']['output']>;
};

export type DeleteSiteRes = {
  __typename?: 'DeleteSiteRes';
  id: Scalars['String']['output'];
};

export type Draft = {
  __typename?: 'Draft';
  events: Array<Event>;
  id: Scalars['ID']['output'];
  lastSaved: Scalars['Float']['output'];
  links: Array<Link>;
  name: Scalars['String']['output'];
  sites: Array<Site>;
  userId: Scalars['String']['output'];
};

export type Event = {
  __typename?: 'Event';
  createdAt: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  operation: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type GetInvitationByOrgResDto = {
  __typename?: 'GetInvitationByOrgResDto';
  invitations: Array<InvitationDto>;
};

export type GetNodesInput = {
  isFree: Scalars['Boolean']['input'];
};

export type GetPackagesForSimInputDto = {
  sim_id: Scalars['String']['input'];
};

export type GetSimBySubscriberInputDto = {
  subscriberId: Scalars['String']['input'];
};

export type GetSimInputDto = {
  simId: Scalars['String']['input'];
};

export type GetSimPackagesDtoApi = {
  __typename?: 'GetSimPackagesDtoAPI';
  packages: Array<SimToPackagesDto>;
  sim_id: Scalars['String']['output'];
};

export type IdResponse = {
  __typename?: 'IdResponse';
  uuid: Scalars['String']['output'];
};

export type InvitationDto = {
  __typename?: 'InvitationDto';
  email: Scalars['String']['output'];
  expiresAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  link: Scalars['String']['output'];
  name: Scalars['String']['output'];
  org: Scalars['String']['output'];
  role: Scalars['String']['output'];
  status: Scalars['String']['output'];
  userId: Scalars['String']['output'];
};

export type Link = {
  __typename?: 'Link';
  id: Scalars['String']['output'];
  siteA: Scalars['String']['output'];
  siteB: Scalars['String']['output'];
};

export type LinkInput = {
  lastSaved: Scalars['Float']['input'];
  siteA: Scalars['String']['input'];
  siteB: Scalars['String']['input'];
};

export type Location = {
  __typename?: 'Location';
  address: Scalars['String']['output'];
  id: Scalars['String']['output'];
  lat: Scalars['String']['output'];
  lng: Scalars['String']['output'];
};

export type LocationInput = {
  address: Scalars['String']['input'];
  lastSaved: Scalars['Float']['input'];
  lat: Scalars['String']['input'];
  lng: Scalars['String']['input'];
};

export type MemberDto = {
  __typename?: 'MemberDto';
  isDeactivated: Scalars['Boolean']['output'];
  memberSince?: Maybe<Scalars['String']['output']>;
  orgId: Scalars['String']['output'];
  role: Scalars['String']['output'];
  userId: Scalars['String']['output'];
};

export type MembersResDto = {
  __typename?: 'MembersResDto';
  members: Array<MemberDto>;
};

export type Mutation = {
  __typename?: 'Mutation';
  AddSiteToNetwork: SiteDto;
  addDraft: Draft;
  addLink: Draft;
  addMember: MemberDto;
  addNetwork: NetworkDto;
  addNode: Node;
  addNodeToSite: CBooleanResponse;
  addPackage: PackageDto;
  addPackageToSim: AddPackageSimResDto;
  addSite: Draft;
  addSubscriber: SubscriberDto;
  allocateSim: AllocateSimApiDto;
  attachNode: CBooleanResponse;
  coverage: Site;
  defaultMarkup: CBooleanResponse;
  deleteDraft: DeleteDraftRes;
  deleteInvitation: DeleteInvitationResDto;
  deleteLink: DeleteLinkRes;
  deleteNodeFromOrg: DeleteNode;
  deletePackage: IdResponse;
  deleteSim: DeleteSimResDto;
  deleteSite: DeleteSiteRes;
  deleteSubscriber: CBooleanResponse;
  detachhNode: CBooleanResponse;
  releaseNodeFromSite: CBooleanResponse;
  removeMember: CBooleanResponse;
  removePackageForSim: RemovePackageFromSimResDto;
  sendInvitation: SendInvitationResDto;
  setActivePackageForSim: SetActivePackageForSimResDto;
  toggleSimStatus: SimStatusResDto;
  updateDraftName: Draft;
  updateEvent: Event;
  updateFirstVisit: UserFistVisitResDto;
  updateInvitation: UpdateInvitationResDto;
  updateLocation: Location;
  updateMember: CBooleanResponse;
  updateNode: Node;
  updateNodeState: Node;
  updatePackage: PackageDto;
  updateSite: Draft;
  updateSubscriber: CBooleanResponse;
  uploadSims: UploadSimsResDto;
};


export type MutationAddSiteToNetworkArgs = {
  data: AddSiteInputDto;
  networkId: Scalars['String']['input'];
};


export type MutationAddDraftArgs = {
  data: AddDraftInput;
};


export type MutationAddLinkArgs = {
  data: LinkInput;
  draftId: Scalars['String']['input'];
};


export type MutationAddMemberArgs = {
  data: AddMemberInputDto;
};


export type MutationAddNetworkArgs = {
  data: AddNetworkInputDto;
};


export type MutationAddNodeArgs = {
  data: AddNodeInput;
};


export type MutationAddNodeToSiteArgs = {
  data: AddNodeToSiteInput;
};


export type MutationAddPackageArgs = {
  data: AddPackageInputDto;
};


export type MutationAddPackageToSimArgs = {
  data: AddPackageToSimInputDto;
};


export type MutationAddSiteArgs = {
  data: SiteInput;
  draftId: Scalars['String']['input'];
};


export type MutationAddSubscriberArgs = {
  data: SubscriberInputDto;
};


export type MutationAllocateSimArgs = {
  data: AllocateSimInputDto;
};


export type MutationAttachNodeArgs = {
  data: AttachNodeInput;
};


export type MutationCoverageArgs = {
  data: CoverageInput;
  siteId: Scalars['String']['input'];
};


export type MutationDefaultMarkupArgs = {
  data: DefaultMarkupInputDto;
};


export type MutationDeleteDraftArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteInvitationArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteLinkArgs = {
  draftId: Scalars['String']['input'];
  lastSaved: Scalars['Float']['input'];
  linkId: Scalars['String']['input'];
};


export type MutationDeleteNodeFromOrgArgs = {
  data: NodeInput;
};


export type MutationDeletePackageArgs = {
  packageId: Scalars['String']['input'];
};


export type MutationDeleteSimArgs = {
  data: DeleteSimInputDto;
};


export type MutationDeleteSiteArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type MutationDetachhNodeArgs = {
  data: NodeInput;
};


export type MutationReleaseNodeFromSiteArgs = {
  data: NodeInput;
};


export type MutationRemoveMemberArgs = {
  id: Scalars['String']['input'];
};


export type MutationRemovePackageForSimArgs = {
  data: RemovePackageFormSimInputDto;
};


export type MutationSendInvitationArgs = {
  data: SendInvitationInputDto;
};


export type MutationSetActivePackageForSimArgs = {
  data: SetActivePackageForSimInputDto;
};


export type MutationToggleSimStatusArgs = {
  data: ToggleSimStatusInputDto;
};


export type MutationUpdateDraftNameArgs = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};


export type MutationUpdateEventArgs = {
  data: UpdateEventInput;
  eventId: Scalars['String']['input'];
};


export type MutationUpdateFirstVisitArgs = {
  data: UserFistVisitInputDto;
};


export type MutationUpdateInvitationArgs = {
  data: UpateInvitationInputDto;
  id: Scalars['String']['input'];
};


export type MutationUpdateLocationArgs = {
  data: LocationInput;
  draftId: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
};


export type MutationUpdateMemberArgs = {
  data: UpdateMemberInputDto;
  memberId: Scalars['String']['input'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeInput;
};


export type MutationUpdateNodeStateArgs = {
  data: UpdateNodeStateInput;
};


export type MutationUpdatePackageArgs = {
  data: UpdatePackageInputDto;
  packageId: Scalars['String']['input'];
};


export type MutationUpdateSiteArgs = {
  data: SiteInput;
  draftId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};


export type MutationUpdateSubscriberArgs = {
  data: UpdateSubscriberInputDto;
  subscriberId: Scalars['String']['input'];
};


export type MutationUploadSimsArgs = {
  data: UploadSimsInputDto;
};

export type NetworkDto = {
  __typename?: 'NetworkDto';
  budget: Scalars['Float']['output'];
  countries: Array<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networks: Array<Scalars['String']['output']>;
  orgId: Scalars['String']['output'];
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
  orgId: Scalars['String']['output'];
};

export type Node = {
  __typename?: 'Node';
  attached: Array<AttachedNodes>;
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  site: NodeSite;
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type NodeApp = {
  __typename?: 'NodeApp';
  cpu: Scalars['String']['output'];
  date: Scalars['Float']['output'];
  memory: Scalars['String']['output'];
  name: Scalars['String']['output'];
  notes: Scalars['String']['output'];
  version: Scalars['String']['output'];
};

export type NodeApps = {
  __typename?: 'NodeApps';
  apps: Array<NodeApp>;
  type: NodeTypeEnum;
};

export type NodeAppsChangeLogInput = {
  type: NodeTypeEnum;
};

export type NodeInput = {
  id: Scalars['String']['input'];
};

export type NodeLocation = {
  __typename?: 'NodeLocation';
  id: Scalars['String']['output'];
  lat: Scalars['String']['output'];
  lng: Scalars['String']['output'];
  state: NodeStatusEnum;
};

export type NodeSite = {
  __typename?: 'NodeSite';
  addedAt?: Maybe<Scalars['String']['output']>;
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
};

export type NodeStatus = {
  __typename?: 'NodeStatus';
  connectivity: Scalars['String']['output'];
  state: Scalars['String']['output'];
};

/** Node status enums */
export enum NodeStatusEnum {
  Active = 'ACTIVE',
  Configured = 'CONFIGURED',
  Faulty = 'FAULTY',
  Maintenance = 'MAINTENANCE',
  Onboarded = 'ONBOARDED',
  Undefined = 'UNDEFINED'
}

/** Node type enums */
export enum NodeTypeEnum {
  Anode = 'anode',
  Hnode = 'hnode',
  Tnode = 'tnode'
}

export type Nodes = {
  __typename?: 'Nodes';
  nodes: Array<Node>;
};

export type NodesInput = {
  networkId: Scalars['String']['input'];
  nodeFilterState: NodeStatusEnum;
};

export type NodesLocation = {
  __typename?: 'NodesLocation';
  networkId: Scalars['String']['output'];
  nodes: Array<NodeLocation>;
};

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgsResDto = {
  __typename?: 'OrgsResDto';
  memberOf: Array<OrgDto>;
  ownerOf: Array<OrgDto>;
  user: Scalars['String']['output'];
};

export type PackageDto = {
  __typename?: 'PackageDto';
  active: Scalars['Boolean']['output'];
  amount: Scalars['Float']['output'];
  apn: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  dataUnit: Scalars['String']['output'];
  dataVolume: Scalars['String']['output'];
  deletedAt: Scalars['String']['output'];
  dlbr: Scalars['String']['output'];
  duration: Scalars['String']['output'];
  flatrate: Scalars['Boolean']['output'];
  from: Scalars['String']['output'];
  markup: PackageMarkupApiDto;
  messageUnit: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  provider: Scalars['String']['output'];
  rate: PackageRateApiDto;
  simType: Scalars['String']['output'];
  smsVolume: Scalars['String']['output'];
  to: Scalars['String']['output'];
  type: Scalars['String']['output'];
  ulbr: Scalars['String']['output'];
  updatedAt: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
  voiceUnit: Scalars['String']['output'];
  voiceVolume: Scalars['String']['output'];
};

export type PackageMarkupApiDto = {
  __typename?: 'PackageMarkupAPIDto';
  baserate: Scalars['String']['output'];
  markup: Scalars['Float']['output'];
};

export type PackageRateApiDto = {
  __typename?: 'PackageRateAPIDto';
  amount: Scalars['Float']['output'];
  data: Scalars['Float']['output'];
  sms_mo: Scalars['String']['output'];
  sms_mt: Scalars['Float']['output'];
};

export type PackagesResDto = {
  __typename?: 'PackagesResDto';
  packages: Array<PackageDto>;
};

export type Query = {
  __typename?: 'Query';
  getAllSites: SitesResDto;
  getAppsChangeLog: AppChangeLogs;
  getDataUsage: SimDataUsage;
  getDefaultMarkup: DefaultMarkupResDto;
  getDefaultMarkupHistory: DefaultMarkupHistoryResDto;
  getDraft: Draft;
  getDrafts: Array<Draft>;
  getInvitation: InvitationDto;
  getInvitationsByOrg: GetInvitationByOrgResDto;
  getMember: MemberDto;
  getMembers: MembersResDto;
  getNetwork: NetworkDto;
  getNetworks: NetworksResDto;
  getNode: Node;
  getNodeApps: NodeApps;
  getNodeLocation: NodeLocation;
  getNodes: Nodes;
  getNodesByNetwork: Nodes;
  getNodesForSite: Nodes;
  getNodesLocation: NodesLocation;
  getOrg: OrgDto;
  getOrgs: OrgsResDto;
  getPackage: PackageDto;
  getPackages: PackagesResDto;
  getPackagesForSim: GetSimPackagesDtoApi;
  getSim: SimDto;
  getSimPoolStats: SimPoolStatsDto;
  getSims: SimsResDto;
  getSimsBySubscriber: SubscriberToSimsDto;
  getSingleSite: SiteDto;
  getSubscriber: SubscriberDto;
  getSubscriberMetricsByNetwork: SubscriberMetricsByNetworkDto;
  getSubscribersByNetwork: SubscribersResDto;
  getUser: UserResDto;
  whoami: WhoamiDto;
};


export type QueryGetAllSitesArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetAppsChangeLogArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetDataUsageArgs = {
  simId: Scalars['String']['input'];
};


export type QueryGetDraftArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetDraftsArgs = {
  userId: Scalars['String']['input'];
};


export type QueryGetInvitationArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetMemberArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  data: NodeInput;
};


export type QueryGetNodeAppsArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetNodeLocationArgs = {
  data: NodeInput;
};


export type QueryGetNodesArgs = {
  data: GetNodesInput;
};


export type QueryGetNodesByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodesForSiteArgs = {
  siteId: Scalars['String']['input'];
};


export type QueryGetNodesLocationArgs = {
  data: NodesInput;
};


export type QueryGetPackageArgs = {
  packageId: Scalars['String']['input'];
};


export type QueryGetPackagesForSimArgs = {
  data: GetPackagesForSimInputDto;
};


export type QueryGetSimArgs = {
  data: GetSimInputDto;
};


export type QueryGetSimPoolStatsArgs = {
  type: Scalars['String']['input'];
};


export type QueryGetSimsArgs = {
  type: Scalars['String']['input'];
};


export type QueryGetSimsBySubscriberArgs = {
  data: GetSimBySubscriberInputDto;
};


export type QueryGetSingleSiteArgs = {
  networkId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};


export type QueryGetSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type QueryGetSubscriberMetricsByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetSubscribersByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetUserArgs = {
  userId: Scalars['String']['input'];
};

export type RemovePackageFormSimInputDto = {
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type RemovePackageFromSimResDto = {
  __typename?: 'RemovePackageFromSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type SendInvitationInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  role: Scalars['String']['input'];
};

export type SendInvitationResDto = {
  __typename?: 'SendInvitationResDto';
  id: Scalars['String']['output'];
  message: Scalars['String']['output'];
};

export type SetActivePackageForSimInputDto = {
  package_id: Scalars['String']['input'];
  sim_id: Scalars['String']['input'];
};

export type SetActivePackageForSimResDto = {
  __typename?: 'SetActivePackageForSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type SimAllocatePackageDto = {
  __typename?: 'SimAllocatePackageDto';
  endDate?: Maybe<Scalars['String']['output']>;
  id?: Maybe<Scalars['String']['output']>;
  isActive?: Maybe<Scalars['Boolean']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  startDate?: Maybe<Scalars['String']['output']>;
};

export type SimDataUsage = {
  __typename?: 'SimDataUsage';
  usage: Scalars['String']['output'];
};

export type SimDto = {
  __typename?: 'SimDto';
  activationCode?: Maybe<Scalars['String']['output']>;
  createdAt?: Maybe<Scalars['String']['output']>;
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isAllocated: Scalars['String']['output'];
  isPhysical: Scalars['String']['output'];
  msisdn: Scalars['String']['output'];
  qrCode: Scalars['String']['output'];
  simType: Scalars['String']['output'];
  smapAddress: Scalars['String']['output'];
};

export type SimPoolStatsDto = {
  __typename?: 'SimPoolStatsDto';
  available: Scalars['Float']['output'];
  consumed: Scalars['Float']['output'];
  esim: Scalars['Float']['output'];
  failed: Scalars['Float']['output'];
  physical: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type SimStatusResDto = {
  __typename?: 'SimStatusResDto';
  simId?: Maybe<Scalars['String']['output']>;
};

export type SimToPackagesDto = {
  __typename?: 'SimToPackagesDto';
  end_date: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_active: Scalars['Boolean']['output'];
  package_id: Scalars['String']['output'];
  start_date: Scalars['String']['output'];
};

export type SimsResDto = {
  __typename?: 'SimsResDto';
  sim: Array<SimDto>;
};

export type Site = {
  __typename?: 'Site';
  apOption: Scalars['String']['output'];
  draftId: Scalars['String']['output'];
  east: Scalars['Float']['output'];
  height: Scalars['Float']['output'];
  id: Scalars['String']['output'];
  isSetlite: Scalars['Boolean']['output'];
  location: Location;
  name: Scalars['String']['output'];
  north: Scalars['Float']['output'];
  populationCovered: Scalars['Float']['output'];
  populationUrl: Scalars['String']['output'];
  solarUptime: Scalars['Float']['output'];
  south: Scalars['Float']['output'];
  status: Scalars['String']['output'];
  totalBoxesCovered: Scalars['Float']['output'];
  url: Scalars['String']['output'];
  west: Scalars['Float']['output'];
};

export type SiteDto = {
  __typename?: 'SiteDto';
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
};

export type SiteInput = {
  address: Scalars['String']['input'];
  apOption: Scalars['String']['input'];
  height: Scalars['Float']['input'];
  isSetlite: Scalars['Boolean']['input'];
  lastSaved: Scalars['Float']['input'];
  lat: Scalars['String']['input'];
  lng: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
  siteName: Scalars['String']['input'];
  solarUptime: Scalars['Float']['input'];
};

export type SitesResDto = {
  __typename?: 'SitesResDto';
  networkId: Scalars['String']['output'];
  sites: Array<SiteDto>;
};

export type SubscriberDto = {
  __typename?: 'SubscriberDto';
  address: Scalars['String']['output'];
  dob: Scalars['String']['output'];
  email: Scalars['String']['output'];
  firstName: Scalars['String']['output'];
  gender: Scalars['String']['output'];
  idSerial: Scalars['String']['output'];
  lastName: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  proofOfIdentification: Scalars['String']['output'];
  sim?: Maybe<Array<SubscriberSimDto>>;
  uuid: Scalars['String']['output'];
};

export type SubscriberInputDto = {
  address?: InputMaybe<Scalars['String']['input']>;
  dob?: InputMaybe<Scalars['String']['input']>;
  email: Scalars['String']['input'];
  first_name?: InputMaybe<Scalars['String']['input']>;
  gender?: InputMaybe<Scalars['String']['input']>;
  id_serial?: InputMaybe<Scalars['String']['input']>;
  last_name?: InputMaybe<Scalars['String']['input']>;
  network_id: Scalars['String']['input'];
  org_id: Scalars['String']['input'];
  phone?: InputMaybe<Scalars['String']['input']>;
  proof_of_identification?: InputMaybe<Scalars['String']['input']>;
};

export type SubscriberMetricsByNetworkDto = {
  __typename?: 'SubscriberMetricsByNetworkDto';
  active: Scalars['Float']['output'];
  inactive: Scalars['Float']['output'];
  terminated: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type SubscriberSimDto = {
  __typename?: 'SubscriberSimDto';
  activationsCount: Scalars['String']['output'];
  allocatedAt: Scalars['String']['output'];
  deactivationsCount: Scalars['String']['output'];
  firstActivatedOn?: Maybe<Scalars['String']['output']>;
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical?: Maybe<Scalars['Boolean']['output']>;
  lastActivatedOn?: Maybe<Scalars['String']['output']>;
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  package?: Maybe<Scalars['String']['output']>;
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type SubscriberToSimsDto = {
  __typename?: 'SubscriberToSimsDto';
  sims: Array<SimDto>;
  subscriber_id: Scalars['String']['output'];
};

export type SubscribersResDto = {
  __typename?: 'SubscribersResDto';
  subscribers: Array<SubscriberDto>;
};

export type ToggleSimStatusInputDto = {
  sim_id: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type UpateInvitationInputDto = {
  status: Scalars['String']['input'];
};

export type UpdateEventInput = {
  operation: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type UpdateInvitationResDto = {
  __typename?: 'UpdateInvitationResDto';
  id: Scalars['String']['output'];
};

export type UpdateMemberInputDto = {
  isDeactivated: Scalars['Boolean']['input'];
  role: Scalars['String']['input'];
};

export type UpdateNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type UpdateNodeStateInput = {
  id: Scalars['String']['input'];
  state: NodeStatusEnum;
};

export type UpdatePackageInputDto = {
  active: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
};

export type UpdateSubscriberInputDto = {
  address?: InputMaybe<Scalars['String']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  first_name?: InputMaybe<Scalars['String']['input']>;
  id_serial?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  proof_of_identification?: InputMaybe<Scalars['String']['input']>;
};

export type UploadSimsInputDto = {
  data: Scalars['String']['input'];
  simType: Scalars['String']['input'];
};

export type UploadSimsResDto = {
  __typename?: 'UploadSimsResDto';
  iccid: Array<Scalars['String']['output']>;
};

export type UserFistVisitInputDto = {
  email: Scalars['String']['input'];
  firstVisit: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type UserFistVisitResDto = {
  __typename?: 'UserFistVisitResDto';
  firstVisit: Scalars['Boolean']['output'];
};

export type UserResDto = {
  __typename?: 'UserResDto';
  authId: Scalars['String']['output'];
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  registeredSince: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type WhoamiDto = {
  __typename?: 'WhoamiDto';
  memberOf: Array<OrgDto>;
  ownerOf: Array<OrgDto>;
  user: UserResDto;
};

export type NodeFragment = { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } };

export type NetworkSiteFragment = { __typename?: 'SiteDto', id: string, name: string, networkId: string, isDeactivated: string, createdAt: string };

export type GetNodeQueryVariables = Exact<{
  data: NodeInput;
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetAllSitesQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetAllSitesQuery = { __typename?: 'Query', getAllSites: { __typename?: 'SitesResDto', sites: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, isDeactivated: string, createdAt: string }> } };

export type GetSingleSiteQueryVariables = Exact<{
  siteId: Scalars['String']['input'];
  networkId: Scalars['String']['input'];
}>;


export type GetSingleSiteQuery = { __typename?: 'Query', getSingleSite: { __typename?: 'SiteDto', id: string, name: string, networkId: string, isDeactivated: string, createdAt: string } };

export type AddSiteToNetworkMutationVariables = Exact<{
  data: AddSiteInputDto;
  networkId: Scalars['String']['input'];
}>;


export type AddSiteToNetworkMutation = { __typename?: 'Mutation', AddSiteToNetwork: { __typename?: 'SiteDto', id: string, name: string, networkId: string, isDeactivated: string, createdAt: string } };

export type AddNodeMutationVariables = Exact<{
  data: AddNodeInput;
}>;


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesQueryVariables = Exact<{
  data: GetNodesInput;
}>;


export type GetNodesQuery = { __typename?: 'Query', getNodes: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type GetNodesByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetNodesByNetworkQuery = { __typename?: 'Query', getNodesByNetwork: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type GetNodesForSiteQueryVariables = Exact<{
  siteId: Scalars['String']['input'];
}>;


export type GetNodesForSiteQuery = { __typename?: 'Query', getNodesForSite: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type DeleteNodeMutationVariables = Exact<{
  data: NodeInput;
}>;


export type DeleteNodeMutation = { __typename?: 'Mutation', deleteNodeFromOrg: { __typename?: 'DeleteNode', id: string } };

export type AttachNodeMutationVariables = Exact<{
  data: AttachNodeInput;
}>;


export type AttachNodeMutation = { __typename?: 'Mutation', attachNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type DetachhNodeMutationVariables = Exact<{
  data: NodeInput;
}>;


export type DetachhNodeMutation = { __typename?: 'Mutation', detachhNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type ReleaseNodeFromSiteMutationVariables = Exact<{
  data: NodeInput;
}>;


export type ReleaseNodeFromSiteMutation = { __typename?: 'Mutation', releaseNodeFromSite: { __typename?: 'CBooleanResponse', success: boolean } };

export type AddNodeToSiteMutationVariables = Exact<{
  data: AddNodeToSiteInput;
}>;


export type AddNodeToSiteMutation = { __typename?: 'Mutation', addNodeToSite: { __typename?: 'CBooleanResponse', success: boolean } };

export type UpdateNodeStateMutationVariables = Exact<{
  data: UpdateNodeStateInput;
}>;


export type UpdateNodeStateMutation = { __typename?: 'Mutation', updateNodeState: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type UpdateNodeMutationVariables = Exact<{
  data: UpdateNodeInput;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, orgId: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodeAppsQueryVariables = Exact<{
  data: NodeAppsChangeLogInput;
}>;


export type GetNodeAppsQuery = { __typename?: 'Query', getNodeApps: { __typename?: 'NodeApps', type: NodeTypeEnum, apps: Array<{ __typename?: 'NodeApp', name: string, date: number, version: string, cpu: string, memory: string, notes: string }> } };

export type GetNodesLocationQueryVariables = Exact<{
  data: NodesInput;
}>;


export type GetNodesLocationQuery = { __typename?: 'Query', getNodesLocation: { __typename?: 'NodesLocation', networkId: string, nodes: Array<{ __typename?: 'NodeLocation', id: string, lat: string, lng: string, state: NodeStatusEnum }> } };

export type GetNodeLocationQueryVariables = Exact<{
  data: NodeInput;
}>;


export type GetNodeLocationQuery = { __typename?: 'Query', getNodeLocation: { __typename?: 'NodeLocation', id: string, lat: string, lng: string, state: NodeStatusEnum } };

export type MemberFragment = { __typename?: 'MemberDto', role: string, orgId: string, userId: string, isDeactivated: boolean, memberSince?: string | null };

export type GetMembersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetMembersQuery = { __typename?: 'Query', getMembers: { __typename?: 'MembersResDto', members: Array<{ __typename?: 'MemberDto', role: string, orgId: string, userId: string, isDeactivated: boolean, memberSince?: string | null }> } };

export type GetMemberQueryVariables = Exact<{
  memberId: Scalars['String']['input'];
}>;


export type GetMemberQuery = { __typename?: 'Query', getMember: { __typename?: 'MemberDto', role: string, orgId: string, userId: string, isDeactivated: boolean, memberSince?: string | null } };

export type AddMemberMutationVariables = Exact<{
  data: AddMemberInputDto;
}>;


export type AddMemberMutation = { __typename?: 'Mutation', addMember: { __typename?: 'MemberDto', role: string, orgId: string, userId: string, isDeactivated: boolean, memberSince?: string | null } };

export type RemoveMemberMutationVariables = Exact<{
  memberId: Scalars['String']['input'];
}>;


export type RemoveMemberMutation = { __typename?: 'Mutation', removeMember: { __typename?: 'CBooleanResponse', success: boolean } };

export type UpdateMemberMutationVariables = Exact<{
  memberId: Scalars['String']['input'];
  data: UpdateMemberInputDto;
}>;


export type UpdateMemberMutation = { __typename?: 'Mutation', updateMember: { __typename?: 'CBooleanResponse', success: boolean } };

export type OrgFragment = { __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string };

export type GetOrgsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgsQuery = { __typename?: 'Query', getOrgs: { __typename?: 'OrgsResDto', user: string, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }> } };

export type GetOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgQuery = { __typename?: 'Query', getOrg: { __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string } };

export type PackageRateFragment = { __typename?: 'PackageDto', rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number } };

export type PackageMarkupFragment = { __typename?: 'PackageDto', markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type SimPackagesFragment = { __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean };

export type PackageFragment = { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type GetPackagesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { __typename?: 'Query', getPackages: { __typename?: 'PackagesResDto', packages: Array<{ __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } }> } };

export type GetPackageQueryVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type GetPackageQuery = { __typename?: 'Query', getPackage: { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type GetSimsBySubscriberQueryVariables = Exact<{
  data: GetSimBySubscriberInputDto;
}>;


export type GetSimsBySubscriberQuery = { __typename?: 'Query', getSimsBySubscriber: { __typename?: 'SubscriberToSimsDto', sims: Array<{ __typename?: 'SimDto', activationCode?: string | null, createdAt?: string | null, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string }> } };

export type AddPackageMutationVariables = Exact<{
  data: AddPackageInputDto;
}>;


export type AddPackageMutation = { __typename?: 'Mutation', addPackage: { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type RemovePackageForSimMutationVariables = Exact<{
  data: RemovePackageFormSimInputDto;
}>;


export type RemovePackageForSimMutation = { __typename?: 'Mutation', removePackageForSim: { __typename?: 'RemovePackageFromSimResDto', packageId?: string | null } };

export type DeletePackageMutationVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type DeletePackageMutation = { __typename?: 'Mutation', deletePackage: { __typename?: 'IdResponse', uuid: string } };

export type AddPackageToSimMutationVariables = Exact<{
  data: AddPackageToSimInputDto;
}>;


export type AddPackageToSimMutation = { __typename?: 'Mutation', addPackageToSim: { __typename?: 'AddPackageSimResDto', packageId?: string | null } };

export type SetActivePackageForSimMutationVariables = Exact<{
  data: SetActivePackageForSimInputDto;
}>;


export type SetActivePackageForSimMutation = { __typename?: 'Mutation', setActivePackageForSim: { __typename?: 'SetActivePackageForSimResDto', packageId?: string | null } };

export type GetPackagesForSimQueryVariables = Exact<{
  data: GetPackagesForSimInputDto;
}>;


export type GetPackagesForSimQuery = { __typename?: 'Query', getPackagesForSim: { __typename?: 'GetSimPackagesDtoAPI', sim_id: string, packages: Array<{ __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean }> } };

export type DeleteSimMutationVariables = Exact<{
  data: DeleteSimInputDto;
}>;


export type DeleteSimMutation = { __typename?: 'Mutation', deleteSim: { __typename?: 'DeleteSimResDto', simId?: string | null } };

export type UpdatePacakgeMutationVariables = Exact<{
  packageId: Scalars['String']['input'];
  data: UpdatePackageInputDto;
}>;


export type UpdatePacakgeMutation = { __typename?: 'Mutation', updatePackage: { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type GetSimpoolStatsQueryVariables = Exact<{
  type: Scalars['String']['input'];
}>;


export type GetSimpoolStatsQuery = { __typename?: 'Query', getSimPoolStats: { __typename?: 'SimPoolStatsDto', total: number, available: number, consumed: number, failed: number, physical: number, esim: number } };

export type UploadSimsMutationVariables = Exact<{
  data: UploadSimsInputDto;
}>;


export type UploadSimsMutation = { __typename?: 'Mutation', uploadSims: { __typename?: 'UploadSimsResDto', iccid: Array<string> } };

export type SimPoolFragment = { __typename?: 'SimDto', activationCode?: string | null, createdAt?: string | null, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string };

export type SimAllocationPackageFragment = { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null };

export type SimAllocationFragment = { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, org_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } };

export type AllocateSimMutationVariables = Exact<{
  data: AllocateSimInputDto;
}>;


export type AllocateSimMutation = { __typename?: 'Mutation', allocateSim: { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, org_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } } };

export type ToggleSimStatusMutationVariables = Exact<{
  data: ToggleSimStatusInputDto;
}>;


export type ToggleSimStatusMutation = { __typename?: 'Mutation', toggleSimStatus: { __typename?: 'SimStatusResDto', simId?: string | null } };

export type GetSimQueryVariables = Exact<{
  data: GetSimInputDto;
}>;


export type GetSimQuery = { __typename?: 'Query', getSim: { __typename?: 'SimDto', activationCode?: string | null, createdAt?: string | null, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string } };

export type GetSimsQueryVariables = Exact<{
  type: Scalars['String']['input'];
}>;


export type GetSimsQuery = { __typename?: 'Query', getSims: { __typename?: 'SimsResDto', sim: Array<{ __typename?: 'SimDto', activationCode?: string | null, createdAt?: string | null, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string }> } };

export type SubscriberSimFragment = { __typename?: 'SubscriberDto', sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> | null };

export type SubscriberFragment = { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> | null };

export type AddSubscriberMutationVariables = Exact<{
  data: SubscriberInputDto;
}>;


export type AddSubscriberMutation = { __typename?: 'Mutation', addSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> | null } };

export type GetSubscriberQueryVariables = Exact<{
  subscriberId: Scalars['String']['input'];
}>;


export type GetSubscriberQuery = { __typename?: 'Query', getSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> | null } };

export type UpdateSubscriberMutationVariables = Exact<{
  subscriberId: Scalars['String']['input'];
  data: UpdateSubscriberInputDto;
}>;


export type UpdateSubscriberMutation = { __typename?: 'Mutation', updateSubscriber: { __typename?: 'CBooleanResponse', success: boolean } };

export type DeleteSubscriberMutationVariables = Exact<{
  subscriberId: Scalars['String']['input'];
}>;


export type DeleteSubscriberMutation = { __typename?: 'Mutation', deleteSubscriber: { __typename?: 'CBooleanResponse', success: boolean } };

export type GetSubscribersByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscribersByNetworkQuery = { __typename?: 'Query', getSubscribersByNetwork: { __typename?: 'SubscribersResDto', subscribers: Array<{ __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> | null }> } };

export type GetSubscriberMetricsByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscriberMetricsByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };

export type UserFragment = { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string };

export type WhoamiQueryVariables = Exact<{ [key: string]: never; }>;


export type WhoamiQuery = { __typename?: 'Query', whoami: { __typename?: 'WhoamiDto', user: { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string }, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }> } };

export type GetUserQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetUserQuery = { __typename?: 'Query', getUser: { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string } };

export type UNetworkFragment = { __typename?: 'NetworkDto', id: string, name: string, orgId: string, budget: number, isDeactivated: string, createdAt: string, countries: Array<string>, networks: Array<string> };

export type GetNetworksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { __typename?: 'Query', getNetworks: { __typename?: 'NetworksResDto', orgId: string, networks: Array<{ __typename?: 'NetworkDto', id: string, name: string, orgId: string, budget: number, isDeactivated: string, createdAt: string, countries: Array<string>, networks: Array<string> }> } };

export type AddNetworkMutationVariables = Exact<{
  data: AddNetworkInputDto;
}>;


export type AddNetworkMutation = { __typename?: 'Mutation', addNetwork: { __typename?: 'NetworkDto', id: string, name: string, orgId: string, budget: number, isDeactivated: string, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type LocationFragment = { __typename?: 'Location', id: string, lat: string, lng: string, address: string };

export type LinkFragment = { __typename?: 'Link', id: string, siteA: string, siteB: string };

export type SiteFragment = { __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } };

export type EventFragment = { __typename?: 'Event', id: string, value: string, operation: string, createdAt: string };

export type DraftFragment = { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> };

export type AddDraftMutationVariables = Exact<{
  data: AddDraftInput;
}>;


export type AddDraftMutation = { __typename?: 'Mutation', addDraft: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateDraftNameMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  name: Scalars['String']['input'];
}>;


export type UpdateDraftNameMutation = { __typename?: 'Mutation', updateDraftName: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type GetDraftsQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetDraftsQuery = { __typename?: 'Query', getDrafts: Array<{ __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> }> };

export type GetDraftQueryVariables = Exact<{
  draftId: Scalars['String']['input'];
}>;


export type GetDraftQuery = { __typename?: 'Query', getDraft: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type AddSiteMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  data: SiteInput;
}>;


export type AddSiteMutation = { __typename?: 'Mutation', addSite: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateSiteMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  data: SiteInput;
}>;


export type UpdateSiteMutation = { __typename?: 'Mutation', updateSite: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateLocationMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
  data: LocationInput;
}>;


export type UpdateLocationMutation = { __typename?: 'Mutation', updateLocation: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } };

export type DeleteDraftMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
}>;


export type DeleteDraftMutation = { __typename?: 'Mutation', deleteDraft: { __typename?: 'DeleteDraftRes', id: string } };

export type DeleteSiteMutationVariables = Exact<{
  siteId: Scalars['String']['input'];
}>;


export type DeleteSiteMutation = { __typename?: 'Mutation', deleteSite: { __typename?: 'DeleteSiteRes', id: string } };

export type DeleteLinkMutationVariables = Exact<{
  lastSaved: Scalars['Float']['input'];
  draftId: Scalars['String']['input'];
  linkId: Scalars['String']['input'];
}>;


export type DeleteLinkMutation = { __typename?: 'Mutation', deleteLink: { __typename?: 'DeleteLinkRes', id: string } };

export type AddLinkMutationVariables = Exact<{
  data: LinkInput;
  draftId: Scalars['String']['input'];
}>;


export type AddLinkMutation = { __typename?: 'Mutation', addLink: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type CoverageMutationVariables = Exact<{
  siteId: Scalars['String']['input'];
  data: CoverageInput;
}>;


export type CoverageMutation = { __typename?: 'Mutation', coverage: { __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } } };

export const NodeFragmentDoc = gql`
    fragment node on Node {
  id
  name
  orgId
  type
  attached {
    id
    name
    orgId
    type
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
  site {
    nodeId
    siteId
    networkId
    addedAt
  }
  status {
    connectivity
    state
  }
}
    `;
export const NetworkSiteFragmentDoc = gql`
    fragment networkSite on SiteDto {
  id
  name
  networkId
  isDeactivated
  createdAt
}
    `;
export const MemberFragmentDoc = gql`
    fragment member on MemberDto {
  role
  orgId
  userId
  isDeactivated
  memberSince
}
    `;
export const OrgFragmentDoc = gql`
    fragment Org on OrgDto {
  id
  name
  owner
  certificate
  isDeactivated
  createdAt
}
    `;
export const SimPackagesFragmentDoc = gql`
    fragment SimPackages on SimToPackagesDto {
  id
  package_id
  start_date
  end_date
  is_active
}
    `;
export const PackageRateFragmentDoc = gql`
    fragment PackageRate on PackageDto {
  rate {
    sms_mo
    sms_mt
    data
    amount
  }
}
    `;
export const PackageMarkupFragmentDoc = gql`
    fragment PackageMarkup on PackageDto {
  markup {
    baserate
    markup
  }
}
    `;
export const PackageFragmentDoc = gql`
    fragment Package on PackageDto {
  uuid
  name
  orgId
  active
  duration
  simType
  createdAt
  deletedAt
  updatedAt
  smsVolume
  dataVolume
  voiceVolume
  ulbr
  dlbr
  type
  dataUnit
  voiceUnit
  messageUnit
  flatrate
  currency
  from
  to
  country
  provider
  apn
  ownerId
  amount
  ...PackageRate
  ...PackageMarkup
}
    ${PackageRateFragmentDoc}
${PackageMarkupFragmentDoc}`;
export const SimPoolFragmentDoc = gql`
    fragment SimPool on SimDto {
  activationCode
  createdAt
  iccid
  id
  isAllocated
  isPhysical
  msisdn
  qrCode
  simType
  smapAddress
}
    `;
export const SimAllocationPackageFragmentDoc = gql`
    fragment SimAllocationPackage on SimAllocatePackageDto {
  id
  packageId
  startDate
  endDate
  isActive
}
    `;
export const SimAllocationFragmentDoc = gql`
    fragment SimAllocation on AllocateSimAPIDto {
  id
  subscriber_id
  network_id
  org_id
  package {
    ...SimAllocationPackage
  }
  iccid
  msisdn
  imsi
  type
  status
  is_physical
  traffic_policy
  firstActivatedOn
  lastActivatedOn
  activationsCount
  deactivationsCount
  allocated_at
  sync_status
}
    ${SimAllocationPackageFragmentDoc}`;
export const SubscriberSimFragmentDoc = gql`
    fragment SubscriberSim on SubscriberDto {
  sim {
    id
    subscriberId
    networkId
    orgId
    iccid
    msisdn
    imsi
    type
    status
    firstActivatedOn
    lastActivatedOn
    activationsCount
    deactivationsCount
    allocatedAt
    isPhysical
    package
  }
}
    `;
export const SubscriberFragmentDoc = gql`
    fragment Subscriber on SubscriberDto {
  uuid
  address
  dob
  email
  firstName
  lastName
  gender
  idSerial
  networkId
  orgId
  phone
  proofOfIdentification
  ...SubscriberSim
}
    ${SubscriberSimFragmentDoc}`;
export const UserFragmentDoc = gql`
    fragment User on UserResDto {
  name
  uuid
  email
  phone
  authId
  isDeactivated
  registeredSince
}
    `;
export const UNetworkFragmentDoc = gql`
    fragment UNetwork on NetworkDto {
  id
  name
  orgId
  budget
  isDeactivated
  createdAt
  countries
  networks
}
    `;
export const LinkFragmentDoc = gql`
    fragment link on Link {
  id
  siteA
  siteB
}
    `;
export const LocationFragmentDoc = gql`
    fragment location on Location {
  id
  lat
  lng
  address
}
    `;
export const SiteFragmentDoc = gql`
    fragment site on Site {
  id
  url
  east
  name
  west
  north
  south
  status
  height
  apOption
  isSetlite
  solarUptime
  populationUrl
  populationCovered
  totalBoxesCovered
  location {
    ...location
  }
}
    ${LocationFragmentDoc}`;
export const EventFragmentDoc = gql`
    fragment event on Event {
  id
  value
  operation
  createdAt
}
    `;
export const DraftFragmentDoc = gql`
    fragment draft on Draft {
  id
  name
  userId
  lastSaved
  links {
    ...link
  }
  sites {
    ...site
  }
  events {
    ...event
  }
}
    ${LinkFragmentDoc}
${SiteFragmentDoc}
${EventFragmentDoc}`;
export const GetNodeDocument = gql`
    query getNode($data: NodeInput!) {
  getNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodeQuery__
 *
 * To run a query within a React component, call `useGetNodeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeQuery(baseOptions: Apollo.QueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
      }
export function useGetNodeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export type GetNodeQueryHookResult = ReturnType<typeof useGetNodeQuery>;
export type GetNodeLazyQueryHookResult = ReturnType<typeof useGetNodeLazyQuery>;
export type GetNodeQueryResult = Apollo.QueryResult<GetNodeQuery, GetNodeQueryVariables>;
export const GetAllSitesDocument = gql`
    query getAllSites($networkId: String!) {
  getAllSites(networkId: $networkId) {
    sites {
      ...networkSite
    }
  }
}
    ${NetworkSiteFragmentDoc}`;

/**
 * __useGetAllSitesQuery__
 *
 * To run a query within a React component, call `useGetAllSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetAllSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetAllSitesQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetAllSitesQuery(baseOptions: Apollo.QueryHookOptions<GetAllSitesQuery, GetAllSitesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetAllSitesQuery, GetAllSitesQueryVariables>(GetAllSitesDocument, options);
      }
export function useGetAllSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetAllSitesQuery, GetAllSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetAllSitesQuery, GetAllSitesQueryVariables>(GetAllSitesDocument, options);
        }
export type GetAllSitesQueryHookResult = ReturnType<typeof useGetAllSitesQuery>;
export type GetAllSitesLazyQueryHookResult = ReturnType<typeof useGetAllSitesLazyQuery>;
export type GetAllSitesQueryResult = Apollo.QueryResult<GetAllSitesQuery, GetAllSitesQueryVariables>;
export const GetSingleSiteDocument = gql`
    query getSingleSite($siteId: String!, $networkId: String!) {
  getSingleSite(siteId: $siteId, networkId: $networkId) {
    ...networkSite
  }
}
    ${NetworkSiteFragmentDoc}`;

/**
 * __useGetSingleSiteQuery__
 *
 * To run a query within a React component, call `useGetSingleSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSingleSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSingleSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSingleSiteQuery(baseOptions: Apollo.QueryHookOptions<GetSingleSiteQuery, GetSingleSiteQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSingleSiteQuery, GetSingleSiteQueryVariables>(GetSingleSiteDocument, options);
      }
export function useGetSingleSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSingleSiteQuery, GetSingleSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSingleSiteQuery, GetSingleSiteQueryVariables>(GetSingleSiteDocument, options);
        }
export type GetSingleSiteQueryHookResult = ReturnType<typeof useGetSingleSiteQuery>;
export type GetSingleSiteLazyQueryHookResult = ReturnType<typeof useGetSingleSiteLazyQuery>;
export type GetSingleSiteQueryResult = Apollo.QueryResult<GetSingleSiteQuery, GetSingleSiteQueryVariables>;
export const AddSiteToNetworkDocument = gql`
    mutation AddSiteToNetwork($data: AddSiteInputDto!, $networkId: String!) {
  AddSiteToNetwork(data: $data, networkId: $networkId) {
    ...networkSite
  }
}
    ${NetworkSiteFragmentDoc}`;
export type AddSiteToNetworkMutationFn = Apollo.MutationFunction<AddSiteToNetworkMutation, AddSiteToNetworkMutationVariables>;

/**
 * __useAddSiteToNetworkMutation__
 *
 * To run a mutation, you first call `useAddSiteToNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSiteToNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSiteToNetworkMutation, { data, loading, error }] = useAddSiteToNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useAddSiteToNetworkMutation(baseOptions?: Apollo.MutationHookOptions<AddSiteToNetworkMutation, AddSiteToNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSiteToNetworkMutation, AddSiteToNetworkMutationVariables>(AddSiteToNetworkDocument, options);
      }
export type AddSiteToNetworkMutationHookResult = ReturnType<typeof useAddSiteToNetworkMutation>;
export type AddSiteToNetworkMutationResult = Apollo.MutationResult<AddSiteToNetworkMutation>;
export type AddSiteToNetworkMutationOptions = Apollo.BaseMutationOptions<AddSiteToNetworkMutation, AddSiteToNetworkMutationVariables>;
export const AddNodeDocument = gql`
    mutation addNode($data: AddNodeInput!) {
  addNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type AddNodeMutationFn = Apollo.MutationFunction<AddNodeMutation, AddNodeMutationVariables>;

/**
 * __useAddNodeMutation__
 *
 * To run a mutation, you first call `useAddNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeMutation, { data, loading, error }] = useAddNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeMutation, AddNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeMutation, AddNodeMutationVariables>(AddNodeDocument, options);
      }
export type AddNodeMutationHookResult = ReturnType<typeof useAddNodeMutation>;
export type AddNodeMutationResult = Apollo.MutationResult<AddNodeMutation>;
export type AddNodeMutationOptions = Apollo.BaseMutationOptions<AddNodeMutation, AddNodeMutationVariables>;
export const GetNodesDocument = gql`
    query getNodes($data: GetNodesInput!) {
  getNodes(data: $data) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesQuery__
 *
 * To run a query within a React component, call `useGetNodesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodesQuery(baseOptions: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
      }
export function useGetNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<typeof useGetNodesLazyQuery>;
export type GetNodesQueryResult = Apollo.QueryResult<GetNodesQuery, GetNodesQueryVariables>;
export const GetNodesByNetworkDocument = gql`
    query getNodesByNetwork($networkId: String!) {
  getNodesByNetwork(networkId: $networkId) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesByNetworkQuery__
 *
 * To run a query within a React component, call `useGetNodesByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetNodesByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetNodesByNetworkQuery, GetNodesByNetworkQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesByNetworkQuery, GetNodesByNetworkQueryVariables>(GetNodesByNetworkDocument, options);
      }
export function useGetNodesByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesByNetworkQuery, GetNodesByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesByNetworkQuery, GetNodesByNetworkQueryVariables>(GetNodesByNetworkDocument, options);
        }
export type GetNodesByNetworkQueryHookResult = ReturnType<typeof useGetNodesByNetworkQuery>;
export type GetNodesByNetworkLazyQueryHookResult = ReturnType<typeof useGetNodesByNetworkLazyQuery>;
export type GetNodesByNetworkQueryResult = Apollo.QueryResult<GetNodesByNetworkQuery, GetNodesByNetworkQueryVariables>;
export const GetNodesForSiteDocument = gql`
    query getNodesForSite($siteId: String!) {
  getNodesForSite(siteId: $siteId) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesForSiteQuery__
 *
 * To run a query within a React component, call `useGetNodesForSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesForSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesForSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetNodesForSiteQuery(baseOptions: Apollo.QueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
      }
export function useGetNodesForSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
export type GetNodesForSiteQueryHookResult = ReturnType<typeof useGetNodesForSiteQuery>;
export type GetNodesForSiteLazyQueryHookResult = ReturnType<typeof useGetNodesForSiteLazyQuery>;
export type GetNodesForSiteQueryResult = Apollo.QueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
export const DeleteNodeDocument = gql`
    mutation deleteNode($data: NodeInput!) {
  deleteNodeFromOrg(data: $data) {
    id
  }
}
    `;
export type DeleteNodeMutationFn = Apollo.MutationFunction<DeleteNodeMutation, DeleteNodeMutationVariables>;

/**
 * __useDeleteNodeMutation__
 *
 * To run a mutation, you first call `useDeleteNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteNodeMutation, { data, loading, error }] = useDeleteNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteNodeMutation(baseOptions?: Apollo.MutationHookOptions<DeleteNodeMutation, DeleteNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteNodeMutation, DeleteNodeMutationVariables>(DeleteNodeDocument, options);
      }
export type DeleteNodeMutationHookResult = ReturnType<typeof useDeleteNodeMutation>;
export type DeleteNodeMutationResult = Apollo.MutationResult<DeleteNodeMutation>;
export type DeleteNodeMutationOptions = Apollo.BaseMutationOptions<DeleteNodeMutation, DeleteNodeMutationVariables>;
export const AttachNodeDocument = gql`
    mutation attachNode($data: AttachNodeInput!) {
  attachNode(data: $data) {
    success
  }
}
    `;
export type AttachNodeMutationFn = Apollo.MutationFunction<AttachNodeMutation, AttachNodeMutationVariables>;

/**
 * __useAttachNodeMutation__
 *
 * To run a mutation, you first call `useAttachNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAttachNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [attachNodeMutation, { data, loading, error }] = useAttachNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAttachNodeMutation(baseOptions?: Apollo.MutationHookOptions<AttachNodeMutation, AttachNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AttachNodeMutation, AttachNodeMutationVariables>(AttachNodeDocument, options);
      }
export type AttachNodeMutationHookResult = ReturnType<typeof useAttachNodeMutation>;
export type AttachNodeMutationResult = Apollo.MutationResult<AttachNodeMutation>;
export type AttachNodeMutationOptions = Apollo.BaseMutationOptions<AttachNodeMutation, AttachNodeMutationVariables>;
export const DetachhNodeDocument = gql`
    mutation detachhNode($data: NodeInput!) {
  detachhNode(data: $data) {
    success
  }
}
    `;
export type DetachhNodeMutationFn = Apollo.MutationFunction<DetachhNodeMutation, DetachhNodeMutationVariables>;

/**
 * __useDetachhNodeMutation__
 *
 * To run a mutation, you first call `useDetachhNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDetachhNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [detachhNodeMutation, { data, loading, error }] = useDetachhNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDetachhNodeMutation(baseOptions?: Apollo.MutationHookOptions<DetachhNodeMutation, DetachhNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DetachhNodeMutation, DetachhNodeMutationVariables>(DetachhNodeDocument, options);
      }
export type DetachhNodeMutationHookResult = ReturnType<typeof useDetachhNodeMutation>;
export type DetachhNodeMutationResult = Apollo.MutationResult<DetachhNodeMutation>;
export type DetachhNodeMutationOptions = Apollo.BaseMutationOptions<DetachhNodeMutation, DetachhNodeMutationVariables>;
export const ReleaseNodeFromSiteDocument = gql`
    mutation releaseNodeFromSite($data: NodeInput!) {
  releaseNodeFromSite(data: $data) {
    success
  }
}
    `;
export type ReleaseNodeFromSiteMutationFn = Apollo.MutationFunction<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;

/**
 * __useReleaseNodeFromSiteMutation__
 *
 * To run a mutation, you first call `useReleaseNodeFromSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useReleaseNodeFromSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [releaseNodeFromSiteMutation, { data, loading, error }] = useReleaseNodeFromSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useReleaseNodeFromSiteMutation(baseOptions?: Apollo.MutationHookOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>(ReleaseNodeFromSiteDocument, options);
      }
export type ReleaseNodeFromSiteMutationHookResult = ReturnType<typeof useReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationResult = Apollo.MutationResult<ReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationOptions = Apollo.BaseMutationOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;
export const AddNodeToSiteDocument = gql`
    mutation addNodeToSite($data: AddNodeToSiteInput!) {
  addNodeToSite(data: $data) {
    success
  }
}
    `;
export type AddNodeToSiteMutationFn = Apollo.MutationFunction<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;

/**
 * __useAddNodeToSiteMutation__
 *
 * To run a mutation, you first call `useAddNodeToSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeToSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeToSiteMutation, { data, loading, error }] = useAddNodeToSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeToSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>(AddNodeToSiteDocument, options);
      }
export type AddNodeToSiteMutationHookResult = ReturnType<typeof useAddNodeToSiteMutation>;
export type AddNodeToSiteMutationResult = Apollo.MutationResult<AddNodeToSiteMutation>;
export type AddNodeToSiteMutationOptions = Apollo.BaseMutationOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;
export const UpdateNodeStateDocument = gql`
    mutation updateNodeState($data: UpdateNodeStateInput!) {
  updateNodeState(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeStateMutationFn = Apollo.MutationFunction<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;

/**
 * __useUpdateNodeStateMutation__
 *
 * To run a mutation, you first call `useUpdateNodeStateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeStateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeStateMutation, { data, loading, error }] = useUpdateNodeStateMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeStateMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>(UpdateNodeStateDocument, options);
      }
export type UpdateNodeStateMutationHookResult = ReturnType<typeof useUpdateNodeStateMutation>;
export type UpdateNodeStateMutationResult = Apollo.MutationResult<UpdateNodeStateMutation>;
export type UpdateNodeStateMutationOptions = Apollo.BaseMutationOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;
export const UpdateNodeDocument = gql`
    mutation UpdateNode($data: UpdateNodeInput!) {
  updateNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeMutationFn = Apollo.MutationFunction<UpdateNodeMutation, UpdateNodeMutationVariables>;

/**
 * __useUpdateNodeMutation__
 *
 * To run a mutation, you first call `useUpdateNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeMutation, { data, loading, error }] = useUpdateNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeMutation, UpdateNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeMutation, UpdateNodeMutationVariables>(UpdateNodeDocument, options);
      }
export type UpdateNodeMutationHookResult = ReturnType<typeof useUpdateNodeMutation>;
export type UpdateNodeMutationResult = Apollo.MutationResult<UpdateNodeMutation>;
export type UpdateNodeMutationOptions = Apollo.BaseMutationOptions<UpdateNodeMutation, UpdateNodeMutationVariables>;
export const GetNodeAppsDocument = gql`
    query getNodeApps($data: NodeAppsChangeLogInput!) {
  getNodeApps(data: $data) {
    apps {
      name
      date
      version
      cpu
      memory
      notes
    }
    type
  }
}
    `;

/**
 * __useGetNodeAppsQuery__
 *
 * To run a query within a React component, call `useGetNodeAppsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeAppsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeAppsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeAppsQuery(baseOptions: Apollo.QueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
      }
export function useGetNodeAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
export type GetNodeAppsQueryHookResult = ReturnType<typeof useGetNodeAppsQuery>;
export type GetNodeAppsLazyQueryHookResult = ReturnType<typeof useGetNodeAppsLazyQuery>;
export type GetNodeAppsQueryResult = Apollo.QueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export const GetNodesLocationDocument = gql`
    query GetNodesLocation($data: NodesInput!) {
  getNodesLocation(data: $data) {
    networkId
    nodes {
      id
      lat
      lng
      state
    }
  }
}
    `;

/**
 * __useGetNodesLocationQuery__
 *
 * To run a query within a React component, call `useGetNodesLocationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesLocationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesLocationQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodesLocationQuery(baseOptions: Apollo.QueryHookOptions<GetNodesLocationQuery, GetNodesLocationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesLocationQuery, GetNodesLocationQueryVariables>(GetNodesLocationDocument, options);
      }
export function useGetNodesLocationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesLocationQuery, GetNodesLocationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesLocationQuery, GetNodesLocationQueryVariables>(GetNodesLocationDocument, options);
        }
export type GetNodesLocationQueryHookResult = ReturnType<typeof useGetNodesLocationQuery>;
export type GetNodesLocationLazyQueryHookResult = ReturnType<typeof useGetNodesLocationLazyQuery>;
export type GetNodesLocationQueryResult = Apollo.QueryResult<GetNodesLocationQuery, GetNodesLocationQueryVariables>;
export const GetNodeLocationDocument = gql`
    query GetNodeLocation($data: NodeInput!) {
  getNodeLocation(data: $data) {
    id
    lat
    lng
    state
  }
}
    `;

/**
 * __useGetNodeLocationQuery__
 *
 * To run a query within a React component, call `useGetNodeLocationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeLocationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeLocationQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeLocationQuery(baseOptions: Apollo.QueryHookOptions<GetNodeLocationQuery, GetNodeLocationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeLocationQuery, GetNodeLocationQueryVariables>(GetNodeLocationDocument, options);
      }
export function useGetNodeLocationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeLocationQuery, GetNodeLocationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeLocationQuery, GetNodeLocationQueryVariables>(GetNodeLocationDocument, options);
        }
export type GetNodeLocationQueryHookResult = ReturnType<typeof useGetNodeLocationQuery>;
export type GetNodeLocationLazyQueryHookResult = ReturnType<typeof useGetNodeLocationLazyQuery>;
export type GetNodeLocationQueryResult = Apollo.QueryResult<GetNodeLocationQuery, GetNodeLocationQueryVariables>;
export const GetMembersDocument = gql`
    query GetMembers {
  getMembers {
    members {
      ...member
    }
  }
}
    ${MemberFragmentDoc}`;

/**
 * __useGetMembersQuery__
 *
 * To run a query within a React component, call `useGetMembersQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMembersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMembersQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetMembersQuery(baseOptions?: Apollo.QueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
      }
export function useGetMembersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
        }
export type GetMembersQueryHookResult = ReturnType<typeof useGetMembersQuery>;
export type GetMembersLazyQueryHookResult = ReturnType<typeof useGetMembersLazyQuery>;
export type GetMembersQueryResult = Apollo.QueryResult<GetMembersQuery, GetMembersQueryVariables>;
export const GetMemberDocument = gql`
    query GetMember($memberId: String!) {
  getMember(id: $memberId) {
    ...member
  }
}
    ${MemberFragmentDoc}`;

/**
 * __useGetMemberQuery__
 *
 * To run a query within a React component, call `useGetMemberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMemberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemberQuery({
 *   variables: {
 *      memberId: // value for 'memberId'
 *   },
 * });
 */
export function useGetMemberQuery(baseOptions: Apollo.QueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
      }
export function useGetMemberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
        }
export type GetMemberQueryHookResult = ReturnType<typeof useGetMemberQuery>;
export type GetMemberLazyQueryHookResult = ReturnType<typeof useGetMemberLazyQuery>;
export type GetMemberQueryResult = Apollo.QueryResult<GetMemberQuery, GetMemberQueryVariables>;
export const AddMemberDocument = gql`
    mutation addMember($data: AddMemberInputDto!) {
  addMember(data: $data) {
    ...member
  }
}
    ${MemberFragmentDoc}`;
export type AddMemberMutationFn = Apollo.MutationFunction<AddMemberMutation, AddMemberMutationVariables>;

/**
 * __useAddMemberMutation__
 *
 * To run a mutation, you first call `useAddMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addMemberMutation, { data, loading, error }] = useAddMemberMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddMemberMutation(baseOptions?: Apollo.MutationHookOptions<AddMemberMutation, AddMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddMemberMutation, AddMemberMutationVariables>(AddMemberDocument, options);
      }
export type AddMemberMutationHookResult = ReturnType<typeof useAddMemberMutation>;
export type AddMemberMutationResult = Apollo.MutationResult<AddMemberMutation>;
export type AddMemberMutationOptions = Apollo.BaseMutationOptions<AddMemberMutation, AddMemberMutationVariables>;
export const RemoveMemberDocument = gql`
    mutation removeMember($memberId: String!) {
  removeMember(id: $memberId) {
    success
  }
}
    `;
export type RemoveMemberMutationFn = Apollo.MutationFunction<RemoveMemberMutation, RemoveMemberMutationVariables>;

/**
 * __useRemoveMemberMutation__
 *
 * To run a mutation, you first call `useRemoveMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeMemberMutation, { data, loading, error }] = useRemoveMemberMutation({
 *   variables: {
 *      memberId: // value for 'memberId'
 *   },
 * });
 */
export function useRemoveMemberMutation(baseOptions?: Apollo.MutationHookOptions<RemoveMemberMutation, RemoveMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveMemberMutation, RemoveMemberMutationVariables>(RemoveMemberDocument, options);
      }
export type RemoveMemberMutationHookResult = ReturnType<typeof useRemoveMemberMutation>;
export type RemoveMemberMutationResult = Apollo.MutationResult<RemoveMemberMutation>;
export type RemoveMemberMutationOptions = Apollo.BaseMutationOptions<RemoveMemberMutation, RemoveMemberMutationVariables>;
export const UpdateMemberDocument = gql`
    mutation updateMember($memberId: String!, $data: UpdateMemberInputDto!) {
  updateMember(memberId: $memberId, data: $data) {
    success
  }
}
    `;
export type UpdateMemberMutationFn = Apollo.MutationFunction<UpdateMemberMutation, UpdateMemberMutationVariables>;

/**
 * __useUpdateMemberMutation__
 *
 * To run a mutation, you first call `useUpdateMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateMemberMutation, { data, loading, error }] = useUpdateMemberMutation({
 *   variables: {
 *      memberId: // value for 'memberId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateMemberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateMemberMutation, UpdateMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateMemberMutation, UpdateMemberMutationVariables>(UpdateMemberDocument, options);
      }
export type UpdateMemberMutationHookResult = ReturnType<typeof useUpdateMemberMutation>;
export type UpdateMemberMutationResult = Apollo.MutationResult<UpdateMemberMutation>;
export type UpdateMemberMutationOptions = Apollo.BaseMutationOptions<UpdateMemberMutation, UpdateMemberMutationVariables>;
export const GetOrgsDocument = gql`
    query getOrgs {
  getOrgs {
    user
    ownerOf {
      ...Org
    }
    memberOf {
      ...Org
    }
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgsQuery__
 *
 * To run a query within a React component, call `useGetOrgsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgsQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
      }
export function useGetOrgsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
export type GetOrgsQueryHookResult = ReturnType<typeof useGetOrgsQuery>;
export type GetOrgsLazyQueryHookResult = ReturnType<typeof useGetOrgsLazyQuery>;
export type GetOrgsQueryResult = Apollo.QueryResult<GetOrgsQuery, GetOrgsQueryVariables>;
export const GetOrgDocument = gql`
    query getOrg {
  getOrg {
    ...Org
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgQuery__
 *
 * To run a query within a React component, call `useGetOrgQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
      }
export function useGetOrgLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
export type GetOrgQueryHookResult = ReturnType<typeof useGetOrgQuery>;
export type GetOrgLazyQueryHookResult = ReturnType<typeof useGetOrgLazyQuery>;
export type GetOrgQueryResult = Apollo.QueryResult<GetOrgQuery, GetOrgQueryVariables>;
export const GetPackagesDocument = gql`
    query getPackages {
  getPackages {
    packages {
      ...Package
    }
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackagesQuery__
 *
 * To run a query within a React component, call `useGetPackagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetPackagesQuery(baseOptions?: Apollo.QueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
      }
export function useGetPackagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
export type GetPackagesQueryHookResult = ReturnType<typeof useGetPackagesQuery>;
export type GetPackagesLazyQueryHookResult = ReturnType<typeof useGetPackagesLazyQuery>;
export type GetPackagesQueryResult = Apollo.QueryResult<GetPackagesQuery, GetPackagesQueryVariables>;
export const GetPackageDocument = gql`
    query getPackage($packageId: String!) {
  getPackage(packageId: $packageId) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackageQuery__
 *
 * To run a query within a React component, call `useGetPackageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackageQuery({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useGetPackageQuery(baseOptions: Apollo.QueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
      }
export function useGetPackageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
export type GetPackageQueryHookResult = ReturnType<typeof useGetPackageQuery>;
export type GetPackageLazyQueryHookResult = ReturnType<typeof useGetPackageLazyQuery>;
export type GetPackageQueryResult = Apollo.QueryResult<GetPackageQuery, GetPackageQueryVariables>;
export const GetSimsBySubscriberDocument = gql`
    query getSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    sims {
      ...SimPool
    }
  }
}
    ${SimPoolFragmentDoc}`;

/**
 * __useGetSimsBySubscriberQuery__
 *
 * To run a query within a React component, call `useGetSimsBySubscriberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsBySubscriberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsBySubscriberQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsBySubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
      }
export function useGetSimsBySubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
export type GetSimsBySubscriberQueryHookResult = ReturnType<typeof useGetSimsBySubscriberQuery>;
export type GetSimsBySubscriberLazyQueryHookResult = ReturnType<typeof useGetSimsBySubscriberLazyQuery>;
export type GetSimsBySubscriberQueryResult = Apollo.QueryResult<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>;
export const AddPackageDocument = gql`
    mutation addPackage($data: AddPackageInputDto!) {
  addPackage(data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type AddPackageMutationFn = Apollo.MutationFunction<AddPackageMutation, AddPackageMutationVariables>;

/**
 * __useAddPackageMutation__
 *
 * To run a mutation, you first call `useAddPackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackageMutation, { data, loading, error }] = useAddPackageMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackageMutation(baseOptions?: Apollo.MutationHookOptions<AddPackageMutation, AddPackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackageMutation, AddPackageMutationVariables>(AddPackageDocument, options);
      }
export type AddPackageMutationHookResult = ReturnType<typeof useAddPackageMutation>;
export type AddPackageMutationResult = Apollo.MutationResult<AddPackageMutation>;
export type AddPackageMutationOptions = Apollo.BaseMutationOptions<AddPackageMutation, AddPackageMutationVariables>;
export const RemovePackageForSimDocument = gql`
    mutation removePackageForSim($data: RemovePackageFormSimInputDto!) {
  removePackageForSim(data: $data) {
    packageId
  }
}
    `;
export type RemovePackageForSimMutationFn = Apollo.MutationFunction<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;

/**
 * __useRemovePackageForSimMutation__
 *
 * To run a mutation, you first call `useRemovePackageForSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemovePackageForSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removePackageForSimMutation, { data, loading, error }] = useRemovePackageForSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRemovePackageForSimMutation(baseOptions?: Apollo.MutationHookOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>(RemovePackageForSimDocument, options);
      }
export type RemovePackageForSimMutationHookResult = ReturnType<typeof useRemovePackageForSimMutation>;
export type RemovePackageForSimMutationResult = Apollo.MutationResult<RemovePackageForSimMutation>;
export type RemovePackageForSimMutationOptions = Apollo.BaseMutationOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;
export const DeletePackageDocument = gql`
    mutation deletePackage($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}
    `;
export type DeletePackageMutationFn = Apollo.MutationFunction<DeletePackageMutation, DeletePackageMutationVariables>;

/**
 * __useDeletePackageMutation__
 *
 * To run a mutation, you first call `useDeletePackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeletePackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deletePackageMutation, { data, loading, error }] = useDeletePackageMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useDeletePackageMutation(baseOptions?: Apollo.MutationHookOptions<DeletePackageMutation, DeletePackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeletePackageMutation, DeletePackageMutationVariables>(DeletePackageDocument, options);
      }
export type DeletePackageMutationHookResult = ReturnType<typeof useDeletePackageMutation>;
export type DeletePackageMutationResult = Apollo.MutationResult<DeletePackageMutation>;
export type DeletePackageMutationOptions = Apollo.BaseMutationOptions<DeletePackageMutation, DeletePackageMutationVariables>;
export const AddPackageToSimDocument = gql`
    mutation addPackageToSim($data: AddPackageToSimInputDto!) {
  addPackageToSim(data: $data) {
    packageId
  }
}
    `;
export type AddPackageToSimMutationFn = Apollo.MutationFunction<AddPackageToSimMutation, AddPackageToSimMutationVariables>;

/**
 * __useAddPackageToSimMutation__
 *
 * To run a mutation, you first call `useAddPackageToSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackageToSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackageToSimMutation, { data, loading, error }] = useAddPackageToSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackageToSimMutation(baseOptions?: Apollo.MutationHookOptions<AddPackageToSimMutation, AddPackageToSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackageToSimMutation, AddPackageToSimMutationVariables>(AddPackageToSimDocument, options);
      }
export type AddPackageToSimMutationHookResult = ReturnType<typeof useAddPackageToSimMutation>;
export type AddPackageToSimMutationResult = Apollo.MutationResult<AddPackageToSimMutation>;
export type AddPackageToSimMutationOptions = Apollo.BaseMutationOptions<AddPackageToSimMutation, AddPackageToSimMutationVariables>;
export const SetActivePackageForSimDocument = gql`
    mutation setActivePackageForSim($data: SetActivePackageForSimInputDto!) {
  setActivePackageForSim(data: $data) {
    packageId
  }
}
    `;
export type SetActivePackageForSimMutationFn = Apollo.MutationFunction<SetActivePackageForSimMutation, SetActivePackageForSimMutationVariables>;

/**
 * __useSetActivePackageForSimMutation__
 *
 * To run a mutation, you first call `useSetActivePackageForSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetActivePackageForSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setActivePackageForSimMutation, { data, loading, error }] = useSetActivePackageForSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSetActivePackageForSimMutation(baseOptions?: Apollo.MutationHookOptions<SetActivePackageForSimMutation, SetActivePackageForSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetActivePackageForSimMutation, SetActivePackageForSimMutationVariables>(SetActivePackageForSimDocument, options);
      }
export type SetActivePackageForSimMutationHookResult = ReturnType<typeof useSetActivePackageForSimMutation>;
export type SetActivePackageForSimMutationResult = Apollo.MutationResult<SetActivePackageForSimMutation>;
export type SetActivePackageForSimMutationOptions = Apollo.BaseMutationOptions<SetActivePackageForSimMutation, SetActivePackageForSimMutationVariables>;
export const GetPackagesForSimDocument = gql`
    query getPackagesForSim($data: GetPackagesForSimInputDto!) {
  getPackagesForSim(data: $data) {
    sim_id
    packages {
      ...SimPackages
    }
  }
}
    ${SimPackagesFragmentDoc}`;

/**
 * __useGetPackagesForSimQuery__
 *
 * To run a query within a React component, call `useGetPackagesForSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesForSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesForSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPackagesForSimQuery(baseOptions: Apollo.QueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
      }
export function useGetPackagesForSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
export type GetPackagesForSimQueryHookResult = ReturnType<typeof useGetPackagesForSimQuery>;
export type GetPackagesForSimLazyQueryHookResult = ReturnType<typeof useGetPackagesForSimLazyQuery>;
export type GetPackagesForSimQueryResult = Apollo.QueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export const DeleteSimDocument = gql`
    mutation deleteSim($data: DeleteSimInputDto!) {
  deleteSim(data: $data) {
    simId
  }
}
    `;
export type DeleteSimMutationFn = Apollo.MutationFunction<DeleteSimMutation, DeleteSimMutationVariables>;

/**
 * __useDeleteSimMutation__
 *
 * To run a mutation, you first call `useDeleteSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSimMutation, { data, loading, error }] = useDeleteSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteSimMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSimMutation, DeleteSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSimMutation, DeleteSimMutationVariables>(DeleteSimDocument, options);
      }
export type DeleteSimMutationHookResult = ReturnType<typeof useDeleteSimMutation>;
export type DeleteSimMutationResult = Apollo.MutationResult<DeleteSimMutation>;
export type DeleteSimMutationOptions = Apollo.BaseMutationOptions<DeleteSimMutation, DeleteSimMutationVariables>;
export const UpdatePacakgeDocument = gql`
    mutation updatePacakge($packageId: String!, $data: UpdatePackageInputDto!) {
  updatePackage(packageId: $packageId, data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type UpdatePacakgeMutationFn = Apollo.MutationFunction<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;

/**
 * __useUpdatePacakgeMutation__
 *
 * To run a mutation, you first call `useUpdatePacakgeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePacakgeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePacakgeMutation, { data, loading, error }] = useUpdatePacakgeMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePacakgeMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>(UpdatePacakgeDocument, options);
      }
export type UpdatePacakgeMutationHookResult = ReturnType<typeof useUpdatePacakgeMutation>;
export type UpdatePacakgeMutationResult = Apollo.MutationResult<UpdatePacakgeMutation>;
export type UpdatePacakgeMutationOptions = Apollo.BaseMutationOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;
export const GetSimpoolStatsDocument = gql`
    query getSimpoolStats($type: String!) {
  getSimPoolStats(type: $type) {
    total
    available
    consumed
    failed
    physical
    esim
  }
}
    `;

/**
 * __useGetSimpoolStatsQuery__
 *
 * To run a query within a React component, call `useGetSimpoolStatsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimpoolStatsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimpoolStatsQuery({
 *   variables: {
 *      type: // value for 'type'
 *   },
 * });
 */
export function useGetSimpoolStatsQuery(baseOptions: Apollo.QueryHookOptions<GetSimpoolStatsQuery, GetSimpoolStatsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimpoolStatsQuery, GetSimpoolStatsQueryVariables>(GetSimpoolStatsDocument, options);
      }
export function useGetSimpoolStatsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimpoolStatsQuery, GetSimpoolStatsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimpoolStatsQuery, GetSimpoolStatsQueryVariables>(GetSimpoolStatsDocument, options);
        }
export type GetSimpoolStatsQueryHookResult = ReturnType<typeof useGetSimpoolStatsQuery>;
export type GetSimpoolStatsLazyQueryHookResult = ReturnType<typeof useGetSimpoolStatsLazyQuery>;
export type GetSimpoolStatsQueryResult = Apollo.QueryResult<GetSimpoolStatsQuery, GetSimpoolStatsQueryVariables>;
export const UploadSimsDocument = gql`
    mutation uploadSims($data: UploadSimsInputDto!) {
  uploadSims(data: $data) {
    iccid
  }
}
    `;
export type UploadSimsMutationFn = Apollo.MutationFunction<UploadSimsMutation, UploadSimsMutationVariables>;

/**
 * __useUploadSimsMutation__
 *
 * To run a mutation, you first call `useUploadSimsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUploadSimsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [uploadSimsMutation, { data, loading, error }] = useUploadSimsMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUploadSimsMutation(baseOptions?: Apollo.MutationHookOptions<UploadSimsMutation, UploadSimsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UploadSimsMutation, UploadSimsMutationVariables>(UploadSimsDocument, options);
      }
export type UploadSimsMutationHookResult = ReturnType<typeof useUploadSimsMutation>;
export type UploadSimsMutationResult = Apollo.MutationResult<UploadSimsMutation>;
export type UploadSimsMutationOptions = Apollo.BaseMutationOptions<UploadSimsMutation, UploadSimsMutationVariables>;
export const AllocateSimDocument = gql`
    mutation allocateSim($data: AllocateSimInputDto!) {
  allocateSim(data: $data) {
    ...SimAllocation
  }
}
    ${SimAllocationFragmentDoc}`;
export type AllocateSimMutationFn = Apollo.MutationFunction<AllocateSimMutation, AllocateSimMutationVariables>;

/**
 * __useAllocateSimMutation__
 *
 * To run a mutation, you first call `useAllocateSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAllocateSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [allocateSimMutation, { data, loading, error }] = useAllocateSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAllocateSimMutation(baseOptions?: Apollo.MutationHookOptions<AllocateSimMutation, AllocateSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AllocateSimMutation, AllocateSimMutationVariables>(AllocateSimDocument, options);
      }
export type AllocateSimMutationHookResult = ReturnType<typeof useAllocateSimMutation>;
export type AllocateSimMutationResult = Apollo.MutationResult<AllocateSimMutation>;
export type AllocateSimMutationOptions = Apollo.BaseMutationOptions<AllocateSimMutation, AllocateSimMutationVariables>;
export const ToggleSimStatusDocument = gql`
    mutation toggleSimStatus($data: ToggleSimStatusInputDto!) {
  toggleSimStatus(data: $data) {
    simId
  }
}
    `;
export type ToggleSimStatusMutationFn = Apollo.MutationFunction<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;

/**
 * __useToggleSimStatusMutation__
 *
 * To run a mutation, you first call `useToggleSimStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleSimStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleSimStatusMutation, { data, loading, error }] = useToggleSimStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleSimStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>(ToggleSimStatusDocument, options);
      }
export type ToggleSimStatusMutationHookResult = ReturnType<typeof useToggleSimStatusMutation>;
export type ToggleSimStatusMutationResult = Apollo.MutationResult<ToggleSimStatusMutation>;
export type ToggleSimStatusMutationOptions = Apollo.BaseMutationOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;
export const GetSimDocument = gql`
    query getSim($data: GetSimInputDto!) {
  getSim(data: $data) {
    ...SimPool
  }
}
    ${SimPoolFragmentDoc}`;

/**
 * __useGetSimQuery__
 *
 * To run a query within a React component, call `useGetSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimQuery(baseOptions: Apollo.QueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
      }
export function useGetSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
export type GetSimQueryHookResult = ReturnType<typeof useGetSimQuery>;
export type GetSimLazyQueryHookResult = ReturnType<typeof useGetSimLazyQuery>;
export type GetSimQueryResult = Apollo.QueryResult<GetSimQuery, GetSimQueryVariables>;
export const GetSimsDocument = gql`
    query getSims($type: String!) {
  getSims(type: $type) {
    sim {
      ...SimPool
    }
  }
}
    ${SimPoolFragmentDoc}`;

/**
 * __useGetSimsQuery__
 *
 * To run a query within a React component, call `useGetSimsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsQuery({
 *   variables: {
 *      type: // value for 'type'
 *   },
 * });
 */
export function useGetSimsQuery(baseOptions: Apollo.QueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
      }
export function useGetSimsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
export type GetSimsQueryHookResult = ReturnType<typeof useGetSimsQuery>;
export type GetSimsLazyQueryHookResult = ReturnType<typeof useGetSimsLazyQuery>;
export type GetSimsQueryResult = Apollo.QueryResult<GetSimsQuery, GetSimsQueryVariables>;
export const AddSubscriberDocument = gql`
    mutation addSubscriber($data: SubscriberInputDto!) {
  addSubscriber(data: $data) {
    ...Subscriber
  }
}
    ${SubscriberFragmentDoc}`;
export type AddSubscriberMutationFn = Apollo.MutationFunction<AddSubscriberMutation, AddSubscriberMutationVariables>;

/**
 * __useAddSubscriberMutation__
 *
 * To run a mutation, you first call `useAddSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSubscriberMutation, { data, loading, error }] = useAddSubscriberMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<AddSubscriberMutation, AddSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSubscriberMutation, AddSubscriberMutationVariables>(AddSubscriberDocument, options);
      }
export type AddSubscriberMutationHookResult = ReturnType<typeof useAddSubscriberMutation>;
export type AddSubscriberMutationResult = Apollo.MutationResult<AddSubscriberMutation>;
export type AddSubscriberMutationOptions = Apollo.BaseMutationOptions<AddSubscriberMutation, AddSubscriberMutationVariables>;
export const GetSubscriberDocument = gql`
    query getSubscriber($subscriberId: String!) {
  getSubscriber(subscriberId: $subscriberId) {
    ...Subscriber
  }
}
    ${SubscriberFragmentDoc}`;

/**
 * __useGetSubscriberQuery__
 *
 * To run a query within a React component, call `useGetSubscriberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscriberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscriberQuery({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *   },
 * });
 */
export function useGetSubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
      }
export function useGetSubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
        }
export type GetSubscriberQueryHookResult = ReturnType<typeof useGetSubscriberQuery>;
export type GetSubscriberLazyQueryHookResult = ReturnType<typeof useGetSubscriberLazyQuery>;
export type GetSubscriberQueryResult = Apollo.QueryResult<GetSubscriberQuery, GetSubscriberQueryVariables>;
export const UpdateSubscriberDocument = gql`
    mutation updateSubscriber($subscriberId: String!, $data: UpdateSubscriberInputDto!) {
  updateSubscriber(subscriberId: $subscriberId, data: $data) {
    success
  }
}
    `;
export type UpdateSubscriberMutationFn = Apollo.MutationFunction<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>;

/**
 * __useUpdateSubscriberMutation__
 *
 * To run a mutation, you first call `useUpdateSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSubscriberMutation, { data, loading, error }] = useUpdateSubscriberMutation({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>(UpdateSubscriberDocument, options);
      }
export type UpdateSubscriberMutationHookResult = ReturnType<typeof useUpdateSubscriberMutation>;
export type UpdateSubscriberMutationResult = Apollo.MutationResult<UpdateSubscriberMutation>;
export type UpdateSubscriberMutationOptions = Apollo.BaseMutationOptions<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>;
export const DeleteSubscriberDocument = gql`
    mutation deleteSubscriber($subscriberId: String!) {
  deleteSubscriber(subscriberId: $subscriberId) {
    success
  }
}
    `;
export type DeleteSubscriberMutationFn = Apollo.MutationFunction<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>;

/**
 * __useDeleteSubscriberMutation__
 *
 * To run a mutation, you first call `useDeleteSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSubscriberMutation, { data, loading, error }] = useDeleteSubscriberMutation({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *   },
 * });
 */
export function useDeleteSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>(DeleteSubscriberDocument, options);
      }
export type DeleteSubscriberMutationHookResult = ReturnType<typeof useDeleteSubscriberMutation>;
export type DeleteSubscriberMutationResult = Apollo.MutationResult<DeleteSubscriberMutation>;
export type DeleteSubscriberMutationOptions = Apollo.BaseMutationOptions<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>;
export const GetSubscribersByNetworkDocument = gql`
    query getSubscribersByNetwork($networkId: String!) {
  getSubscribersByNetwork(networkId: $networkId) {
    subscribers {
      ...Subscriber
    }
  }
}
    ${SubscriberFragmentDoc}`;

/**
 * __useGetSubscribersByNetworkQuery__
 *
 * To run a query within a React component, call `useGetSubscribersByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscribersByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscribersByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSubscribersByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
      }
export function useGetSubscribersByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
        }
export type GetSubscribersByNetworkQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkQuery>;
export type GetSubscribersByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkLazyQuery>;
export type GetSubscribersByNetworkQueryResult = Apollo.QueryResult<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>;
export const GetSubscriberMetricsByNetworkDocument = gql`
    query getSubscriberMetricsByNetwork($networkId: String!) {
  getSubscriberMetricsByNetwork(networkId: $networkId) {
    total
    active
    inactive
    terminated
  }
}
    `;

/**
 * __useGetSubscriberMetricsByNetworkQuery__
 *
 * To run a query within a React component, call `useGetSubscriberMetricsByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscriberMetricsByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscriberMetricsByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSubscriberMetricsByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
      }
export function useGetSubscriberMetricsByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
        }
export type GetSubscriberMetricsByNetworkQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkQuery>;
export type GetSubscriberMetricsByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkLazyQuery>;
export type GetSubscriberMetricsByNetworkQueryResult = Apollo.QueryResult<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>;
export const WhoamiDocument = gql`
    query Whoami {
  whoami {
    user {
      ...User
    }
    ownerOf {
      ...Org
    }
    memberOf {
      ...Org
    }
  }
}
    ${UserFragmentDoc}
${OrgFragmentDoc}`;

/**
 * __useWhoamiQuery__
 *
 * To run a query within a React component, call `useWhoamiQuery` and pass it any options that fit your needs.
 * When your component renders, `useWhoamiQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useWhoamiQuery({
 *   variables: {
 *   },
 * });
 */
export function useWhoamiQuery(baseOptions?: Apollo.QueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
      }
export function useWhoamiLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
        }
export type WhoamiQueryHookResult = ReturnType<typeof useWhoamiQuery>;
export type WhoamiLazyQueryHookResult = ReturnType<typeof useWhoamiLazyQuery>;
export type WhoamiQueryResult = Apollo.QueryResult<WhoamiQuery, WhoamiQueryVariables>;
export const GetUserDocument = gql`
    query GetUser($userId: String!) {
  getUser(userId: $userId) {
    ...User
  }
}
    ${UserFragmentDoc}`;

/**
 * __useGetUserQuery__
 *
 * To run a query within a React component, call `useGetUserQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetUserQuery(baseOptions: Apollo.QueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
      }
export function useGetUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export type GetUserQueryHookResult = ReturnType<typeof useGetUserQuery>;
export type GetUserLazyQueryHookResult = ReturnType<typeof useGetUserLazyQuery>;
export type GetUserQueryResult = Apollo.QueryResult<GetUserQuery, GetUserQueryVariables>;
export const GetNetworksDocument = gql`
    query getNetworks {
  getNetworks {
    orgId
    networks {
      ...UNetwork
    }
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworksQuery__
 *
 * To run a query within a React component, call `useGetNetworksQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworksQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworksQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNetworksQuery(baseOptions?: Apollo.QueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
      }
export function useGetNetworksLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
export type GetNetworksQueryHookResult = ReturnType<typeof useGetNetworksQuery>;
export type GetNetworksLazyQueryHookResult = ReturnType<typeof useGetNetworksLazyQuery>;
export type GetNetworksQueryResult = Apollo.QueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export const AddNetworkDocument = gql`
    mutation AddNetwork($data: AddNetworkInputDto!) {
  addNetwork(data: $data) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;
export type AddNetworkMutationFn = Apollo.MutationFunction<AddNetworkMutation, AddNetworkMutationVariables>;

/**
 * __useAddNetworkMutation__
 *
 * To run a mutation, you first call `useAddNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNetworkMutation, { data, loading, error }] = useAddNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNetworkMutation(baseOptions?: Apollo.MutationHookOptions<AddNetworkMutation, AddNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNetworkMutation, AddNetworkMutationVariables>(AddNetworkDocument, options);
      }
export type AddNetworkMutationHookResult = ReturnType<typeof useAddNetworkMutation>;
export type AddNetworkMutationResult = Apollo.MutationResult<AddNetworkMutation>;
export type AddNetworkMutationOptions = Apollo.BaseMutationOptions<AddNetworkMutation, AddNetworkMutationVariables>;
export const AddDraftDocument = gql`
    mutation AddDraft($data: AddDraftInput!) {
  addDraft(data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddDraftMutationFn = Apollo.MutationFunction<AddDraftMutation, AddDraftMutationVariables>;

/**
 * __useAddDraftMutation__
 *
 * To run a mutation, you first call `useAddDraftMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddDraftMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addDraftMutation, { data, loading, error }] = useAddDraftMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddDraftMutation(baseOptions?: Apollo.MutationHookOptions<AddDraftMutation, AddDraftMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddDraftMutation, AddDraftMutationVariables>(AddDraftDocument, options);
      }
export type AddDraftMutationHookResult = ReturnType<typeof useAddDraftMutation>;
export type AddDraftMutationResult = Apollo.MutationResult<AddDraftMutation>;
export type AddDraftMutationOptions = Apollo.BaseMutationOptions<AddDraftMutation, AddDraftMutationVariables>;
export const UpdateDraftNameDocument = gql`
    mutation UpdateDraftName($draftId: String!, $name: String!) {
  updateDraftName(id: $draftId, name: $name) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type UpdateDraftNameMutationFn = Apollo.MutationFunction<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>;

/**
 * __useUpdateDraftNameMutation__
 *
 * To run a mutation, you first call `useUpdateDraftNameMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateDraftNameMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateDraftNameMutation, { data, loading, error }] = useUpdateDraftNameMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      name: // value for 'name'
 *   },
 * });
 */
export function useUpdateDraftNameMutation(baseOptions?: Apollo.MutationHookOptions<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>(UpdateDraftNameDocument, options);
      }
export type UpdateDraftNameMutationHookResult = ReturnType<typeof useUpdateDraftNameMutation>;
export type UpdateDraftNameMutationResult = Apollo.MutationResult<UpdateDraftNameMutation>;
export type UpdateDraftNameMutationOptions = Apollo.BaseMutationOptions<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>;
export const GetDraftsDocument = gql`
    query GetDrafts($userId: String!) {
  getDrafts(userId: $userId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;

/**
 * __useGetDraftsQuery__
 *
 * To run a query within a React component, call `useGetDraftsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDraftsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDraftsQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetDraftsQuery(baseOptions: Apollo.QueryHookOptions<GetDraftsQuery, GetDraftsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDraftsQuery, GetDraftsQueryVariables>(GetDraftsDocument, options);
      }
export function useGetDraftsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDraftsQuery, GetDraftsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDraftsQuery, GetDraftsQueryVariables>(GetDraftsDocument, options);
        }
export type GetDraftsQueryHookResult = ReturnType<typeof useGetDraftsQuery>;
export type GetDraftsLazyQueryHookResult = ReturnType<typeof useGetDraftsLazyQuery>;
export type GetDraftsQueryResult = Apollo.QueryResult<GetDraftsQuery, GetDraftsQueryVariables>;
export const GetDraftDocument = gql`
    query GetDraft($draftId: String!) {
  getDraft(id: $draftId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;

/**
 * __useGetDraftQuery__
 *
 * To run a query within a React component, call `useGetDraftQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDraftQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDraftQuery({
 *   variables: {
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useGetDraftQuery(baseOptions: Apollo.QueryHookOptions<GetDraftQuery, GetDraftQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDraftQuery, GetDraftQueryVariables>(GetDraftDocument, options);
      }
export function useGetDraftLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDraftQuery, GetDraftQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDraftQuery, GetDraftQueryVariables>(GetDraftDocument, options);
        }
export type GetDraftQueryHookResult = ReturnType<typeof useGetDraftQuery>;
export type GetDraftLazyQueryHookResult = ReturnType<typeof useGetDraftLazyQuery>;
export type GetDraftQueryResult = Apollo.QueryResult<GetDraftQuery, GetDraftQueryVariables>;
export const AddSiteDocument = gql`
    mutation addSite($draftId: String!, $data: SiteInput!) {
  addSite(draftId: $draftId, data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddSiteMutationFn = Apollo.MutationFunction<AddSiteMutation, AddSiteMutationVariables>;

/**
 * __useAddSiteMutation__
 *
 * To run a mutation, you first call `useAddSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSiteMutation, { data, loading, error }] = useAddSiteMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddSiteMutation, AddSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSiteMutation, AddSiteMutationVariables>(AddSiteDocument, options);
      }
export type AddSiteMutationHookResult = ReturnType<typeof useAddSiteMutation>;
export type AddSiteMutationResult = Apollo.MutationResult<AddSiteMutation>;
export type AddSiteMutationOptions = Apollo.BaseMutationOptions<AddSiteMutation, AddSiteMutationVariables>;
export const UpdateSiteDocument = gql`
    mutation UpdateSite($draftId: String!, $siteId: String!, $data: SiteInput!) {
  updateSite(draftId: $draftId, siteId: $siteId, data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type UpdateSiteMutationFn = Apollo.MutationFunction<UpdateSiteMutation, UpdateSiteMutationVariables>;

/**
 * __useUpdateSiteMutation__
 *
 * To run a mutation, you first call `useUpdateSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSiteMutation, { data, loading, error }] = useUpdateSiteMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      siteId: // value for 'siteId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSiteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSiteMutation, UpdateSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSiteMutation, UpdateSiteMutationVariables>(UpdateSiteDocument, options);
      }
export type UpdateSiteMutationHookResult = ReturnType<typeof useUpdateSiteMutation>;
export type UpdateSiteMutationResult = Apollo.MutationResult<UpdateSiteMutation>;
export type UpdateSiteMutationOptions = Apollo.BaseMutationOptions<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const UpdateLocationDocument = gql`
    mutation UpdateLocation($draftId: String!, $locationId: String!, $data: LocationInput!) {
  updateLocation(draftId: $draftId, locationId: $locationId, data: $data) {
    ...location
  }
}
    ${LocationFragmentDoc}`;
export type UpdateLocationMutationFn = Apollo.MutationFunction<UpdateLocationMutation, UpdateLocationMutationVariables>;

/**
 * __useUpdateLocationMutation__
 *
 * To run a mutation, you first call `useUpdateLocationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateLocationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateLocationMutation, { data, loading, error }] = useUpdateLocationMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      locationId: // value for 'locationId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateLocationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateLocationMutation, UpdateLocationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateLocationMutation, UpdateLocationMutationVariables>(UpdateLocationDocument, options);
      }
export type UpdateLocationMutationHookResult = ReturnType<typeof useUpdateLocationMutation>;
export type UpdateLocationMutationResult = Apollo.MutationResult<UpdateLocationMutation>;
export type UpdateLocationMutationOptions = Apollo.BaseMutationOptions<UpdateLocationMutation, UpdateLocationMutationVariables>;
export const DeleteDraftDocument = gql`
    mutation DeleteDraft($draftId: String!) {
  deleteDraft(id: $draftId) {
    id
  }
}
    `;
export type DeleteDraftMutationFn = Apollo.MutationFunction<DeleteDraftMutation, DeleteDraftMutationVariables>;

/**
 * __useDeleteDraftMutation__
 *
 * To run a mutation, you first call `useDeleteDraftMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteDraftMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteDraftMutation, { data, loading, error }] = useDeleteDraftMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useDeleteDraftMutation(baseOptions?: Apollo.MutationHookOptions<DeleteDraftMutation, DeleteDraftMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteDraftMutation, DeleteDraftMutationVariables>(DeleteDraftDocument, options);
      }
export type DeleteDraftMutationHookResult = ReturnType<typeof useDeleteDraftMutation>;
export type DeleteDraftMutationResult = Apollo.MutationResult<DeleteDraftMutation>;
export type DeleteDraftMutationOptions = Apollo.BaseMutationOptions<DeleteDraftMutation, DeleteDraftMutationVariables>;
export const DeleteSiteDocument = gql`
    mutation DeleteSite($siteId: String!) {
  deleteSite(id: $siteId) {
    id
  }
}
    `;
export type DeleteSiteMutationFn = Apollo.MutationFunction<DeleteSiteMutation, DeleteSiteMutationVariables>;

/**
 * __useDeleteSiteMutation__
 *
 * To run a mutation, you first call `useDeleteSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSiteMutation, { data, loading, error }] = useDeleteSiteMutation({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useDeleteSiteMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSiteMutation, DeleteSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSiteMutation, DeleteSiteMutationVariables>(DeleteSiteDocument, options);
      }
export type DeleteSiteMutationHookResult = ReturnType<typeof useDeleteSiteMutation>;
export type DeleteSiteMutationResult = Apollo.MutationResult<DeleteSiteMutation>;
export type DeleteSiteMutationOptions = Apollo.BaseMutationOptions<DeleteSiteMutation, DeleteSiteMutationVariables>;
export const DeleteLinkDocument = gql`
    mutation DeleteLink($lastSaved: Float!, $draftId: String!, $linkId: String!) {
  deleteLink(lastSaved: $lastSaved, draftId: $draftId, linkId: $linkId) {
    id
  }
}
    `;
export type DeleteLinkMutationFn = Apollo.MutationFunction<DeleteLinkMutation, DeleteLinkMutationVariables>;

/**
 * __useDeleteLinkMutation__
 *
 * To run a mutation, you first call `useDeleteLinkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteLinkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteLinkMutation, { data, loading, error }] = useDeleteLinkMutation({
 *   variables: {
 *      lastSaved: // value for 'lastSaved'
 *      draftId: // value for 'draftId'
 *      linkId: // value for 'linkId'
 *   },
 * });
 */
export function useDeleteLinkMutation(baseOptions?: Apollo.MutationHookOptions<DeleteLinkMutation, DeleteLinkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteLinkMutation, DeleteLinkMutationVariables>(DeleteLinkDocument, options);
      }
export type DeleteLinkMutationHookResult = ReturnType<typeof useDeleteLinkMutation>;
export type DeleteLinkMutationResult = Apollo.MutationResult<DeleteLinkMutation>;
export type DeleteLinkMutationOptions = Apollo.BaseMutationOptions<DeleteLinkMutation, DeleteLinkMutationVariables>;
export const AddLinkDocument = gql`
    mutation AddLink($data: LinkInput!, $draftId: String!) {
  addLink(data: $data, draftId: $draftId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddLinkMutationFn = Apollo.MutationFunction<AddLinkMutation, AddLinkMutationVariables>;

/**
 * __useAddLinkMutation__
 *
 * To run a mutation, you first call `useAddLinkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddLinkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addLinkMutation, { data, loading, error }] = useAddLinkMutation({
 *   variables: {
 *      data: // value for 'data'
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useAddLinkMutation(baseOptions?: Apollo.MutationHookOptions<AddLinkMutation, AddLinkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddLinkMutation, AddLinkMutationVariables>(AddLinkDocument, options);
      }
export type AddLinkMutationHookResult = ReturnType<typeof useAddLinkMutation>;
export type AddLinkMutationResult = Apollo.MutationResult<AddLinkMutation>;
export type AddLinkMutationOptions = Apollo.BaseMutationOptions<AddLinkMutation, AddLinkMutationVariables>;
export const CoverageDocument = gql`
    mutation Coverage($siteId: String!, $data: CoverageInput!) {
  coverage(data: $data, siteId: $siteId) {
    ...site
  }
}
    ${SiteFragmentDoc}`;
export type CoverageMutationFn = Apollo.MutationFunction<CoverageMutation, CoverageMutationVariables>;

/**
 * __useCoverageMutation__
 *
 * To run a mutation, you first call `useCoverageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCoverageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [coverageMutation, { data, loading, error }] = useCoverageMutation({
 *   variables: {
 *      siteId: // value for 'siteId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useCoverageMutation(baseOptions?: Apollo.MutationHookOptions<CoverageMutation, CoverageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CoverageMutation, CoverageMutationVariables>(CoverageDocument, options);
      }
export type CoverageMutationHookResult = ReturnType<typeof useCoverageMutation>;
export type CoverageMutationResult = Apollo.MutationResult<CoverageMutation>;
export type CoverageMutationOptions = Apollo.BaseMutationOptions<CoverageMutation, CoverageMutationVariables>;