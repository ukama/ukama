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
  ID: { input: string | number; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
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
  network_name: Scalars['String']['input'];
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

export type AddOrgInputDto = {
  certificate: Scalars['String']['input'];
  name: Scalars['String']['input'];
  owner_uuid: Scalars['String']['input'];
};

export type AddPackageInputDto = {
  amount: Scalars['Float']['input'];
  dataUnit: Scalars['String']['input'];
  dataVolume: Scalars['Int']['input'];
  duration: Scalars['Int']['input'];
  name: Scalars['String']['input'];
};

export type AddSiteInputDto = {
  site: Scalars['String']['input'];
};

export type AllocateSimInputDto = {
  iccid: Scalars['String']['input'];
  networkId: Scalars['String']['input'];
  packageId: Scalars['String']['input'];
  simType: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
};

export type AttachNodeInput = {
  anodel: Scalars['String']['input'];
  anoder: Scalars['String']['input'];
  parentNode: Scalars['String']['input'];
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

export type GetNodes = {
  __typename?: 'GetNodes';
  nodes: Array<Node>;
};

export type GetNodesInput = {
  isFree: Scalars['Boolean']['input'];
};

export type GetSimInputDto = {
  simId: Scalars['String']['input'];
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

export type MemberInputDto = {
  memberId: Scalars['Boolean']['input'];
  orgName: Scalars['Boolean']['input'];
};

export type MemberObj = {
  __typename?: 'MemberObj';
  isDeactivated: Scalars['Boolean']['output'];
  memberSince?: Maybe<Scalars['String']['output']>;
  orgId: Scalars['String']['output'];
  role: Scalars['String']['output'];
  user: UserResDto;
  userId: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addDraft: Draft;
  addLink: Draft;
  addMember: MemberObj;
  addNetwork: NetworkDto;
  addNode: Node;
  addNodeToSite: CBooleanResponse;
  addOrg: OrgDto;
  addPackage: PackageDto;
  addSite: Draft;
  addSubscriber: SubscriberDto;
  allocateSim: SimDto;
  attachNode: CBooleanResponse;
  coverage: Site;
  defaultMarkup: CBooleanResponse;
  deleteDraft: DeleteDraftRes;
  deleteInvitation: DeleteInvitationResDto;
  deleteLink: DeleteLinkRes;
  deleteNodeFromOrg: DeleteNode;
  deletePackage: IdResponse;
  deleteSite: DeleteSiteRes;
  deleteSubscriber: CBooleanResponse;
  detachhNode: CBooleanResponse;
  getSim: SetActivePackageForSimResDto;
  releaseNodeFromSite: CBooleanResponse;
  removeMember: CBooleanResponse;
  sendInvitation: SendInvitationResDto;
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


export type MutationAddOrgArgs = {
  data: AddOrgInputDto;
};


export type MutationAddPackageArgs = {
  data: AddPackageInputDto;
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


export type MutationDeleteSiteArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type MutationDetachhNodeArgs = {
  data: NodeInput;
};


export type MutationGetSimArgs = {
  data: SetActivePackageForSimInputDto;
};


export type MutationReleaseNodeFromSiteArgs = {
  data: NodeInput;
};


export type MutationRemoveMemberArgs = {
  data: MemberInputDto;
};


export type MutationSendInvitationArgs = {
  data: SendInvitationInputDto;
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
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
  orgId: Scalars['String']['output'];
};

export type Node = {
  __typename?: 'Node';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type NodeInput = {
  id: Scalars['String']['input'];
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

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgMembersResDto = {
  __typename?: 'OrgMembersResDto';
  members: Array<MemberObj>;
  org: Scalars['String']['output'];
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
  addSite: SiteDto;
  getDataUsage: SimDataUsage;
  getDefaultMarkup: DefaultMarkupResDto;
  getDefaultMarkupHistory: DefaultMarkupHistoryResDto;
  getDraft: Draft;
  getDrafts: Array<Draft>;
  getInvitation: InvitationDto;
  getInvitationsByOrg: GetInvitationByOrgResDto;
  getNetwork: NetworkDto;
  getNetworks: NetworksResDto;
  getNode: Node;
  getNodes: GetNodes;
  getOrg: OrgDto;
  getOrgMembers: OrgMembersResDto;
  getOrgs: OrgsResDto;
  getPackage: PackageDto;
  getPackages: PackagesResDto;
  getSim: SimDto;
  getSimPoolStats: SimPoolStatsDto;
  getSims: SimsResDto;
  getSite: SiteDto;
  getSites: SitesResDto;
  getSubscriber: SubscriberDto;
  getSubscriberMetricsByNetwork: SubscriberMetricsByNetworkDto;
  getSubscribersByNetwork: SubscribersResDto;
  getUser: UserResDto;
  whoami: WhoamiDto;
};


export type QueryAddSiteArgs = {
  data: AddSiteInputDto;
  networkId: Scalars['String']['input'];
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


export type QueryGetNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  data: NodeInput;
};


export type QueryGetNodesArgs = {
  data: GetNodesInput;
};


export type QueryGetPackageArgs = {
  packageId: Scalars['String']['input'];
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


export type QueryGetSiteArgs = {
  networkId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};


export type QueryGetSitesArgs = {
  networkId: Scalars['String']['input'];
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
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type SetActivePackageForSimResDto = {
  __typename?: 'SetActivePackageForSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type SimDataUsage = {
  __typename?: 'SimDataUsage';
  usage: Scalars['String']['output'];
};

export type SimDto = {
  __typename?: 'SimDto';
  activationCode: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
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
  sim: Array<SubscriberSimDto>;
  uuid: Scalars['String']['output'];
};

export type SubscriberInputDto = {
  address: Scalars['String']['input'];
  dob: Scalars['String']['input'];
  email: Scalars['String']['input'];
  first_name: Scalars['String']['input'];
  gender: Scalars['String']['input'];
  id_serial: Scalars['String']['input'];
  last_name: Scalars['String']['input'];
  network_id: Scalars['String']['input'];
  org_id: Scalars['String']['input'];
  phone: Scalars['String']['input'];
  proof_of_identification: Scalars['String']['input'];
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

export type SubscribersResDto = {
  __typename?: 'SubscribersResDto';
  subscribers: Array<SubscriberDto>;
};

export type ToggleSimStatusInputDto = {
  simId: Scalars['String']['input'];
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
  orgName: Scalars['String']['input'];
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

export type NodeFragment = { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } };

export type GetNodeQueryVariables = Exact<{
  data: NodeInput;
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesQueryVariables = Exact<{
  data: GetNodesInput;
}>;


export type GetNodesQuery = { __typename?: 'Query', getNodes: { __typename?: 'GetNodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

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

export type AddNodeMutationVariables = Exact<{
  data: AddNodeInput;
}>;


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

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


export type UpdateNodeStateMutation = { __typename?: 'Mutation', updateNodeState: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type UpdateNodeMutationVariables = Exact<{
  data: UpdateNodeInput;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'Node', id: string, name: string, orgId: string, type: NodeTypeEnum, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type OrgFragment = { __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string };

export type OrgUserFragment = { __typename?: 'UserResDto', name: string, email: string, uuid: string, phone: string, isDeactivated: boolean, registeredSince: string };

export type MemberFragment = { __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null };

export type GetOrgMemberQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgMemberQuery = { __typename?: 'Query', getOrgMembers: { __typename?: 'OrgMembersResDto', org: string, members: Array<{ __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null, user: { __typename?: 'UserResDto', name: string, email: string, uuid: string, phone: string, isDeactivated: boolean, registeredSince: string } }> } };

export type AddMemberMutationVariables = Exact<{
  data: AddMemberInputDto;
}>;


export type AddMemberMutation = { __typename?: 'Mutation', addMember: { __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null } };

export type GetOrgsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgsQuery = { __typename?: 'Query', getOrgs: { __typename?: 'OrgsResDto', user: string, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string }> } };

export type GetOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgQuery = { __typename?: 'Query', getOrg: { __typename?: 'OrgDto', id: string, name: string, owner: string, certificate: string, isDeactivated: boolean, createdAt: string } };

export type PackageRateFragment = { __typename?: 'PackageDto', rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number } };

export type PackageMarkupFragment = { __typename?: 'PackageDto', markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type PackageFragment = { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type GetPackagesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { __typename?: 'Query', getPackages: { __typename?: 'PackagesResDto', packages: Array<{ __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } }> } };

export type GetPackageQueryVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type GetPackageQuery = { __typename?: 'Query', getPackage: { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type AddPackageMutationVariables = Exact<{
  data: AddPackageInputDto;
}>;


export type AddPackageMutation = { __typename?: 'Mutation', addPackage: { __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type DeletePacakgeMutationVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type DeletePacakgeMutation = { __typename?: 'Mutation', deletePackage: { __typename?: 'IdResponse', uuid: string } };

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

export type SimPoolFragment = { __typename?: 'SimDto', activationCode: string, createdAt: string, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string };

export type GetSimsQueryVariables = Exact<{
  type: Scalars['String']['input'];
}>;


export type GetSimsQuery = { __typename?: 'Query', getSims: { __typename?: 'SimsResDto', sim: Array<{ __typename?: 'SimDto', activationCode: string, createdAt: string, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string }> } };

export type SubscriberSimFragment = { __typename?: 'SubscriberDto', sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> };

export type SubscriberFragment = { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> };

export type AddSubscriberMutationVariables = Exact<{
  data: SubscriberInputDto;
}>;


export type AddSubscriberMutation = { __typename?: 'Mutation', addSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> } };

export type GetSubscriberQueryVariables = Exact<{
  subscriberId: Scalars['String']['input'];
}>;


export type GetSubscriberQuery = { __typename?: 'Query', getSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> } };

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


export type GetSubscribersByNetworkQuery = { __typename?: 'Query', getSubscribersByNetwork: { __typename?: 'SubscribersResDto', subscribers: Array<{ __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> }> } };

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

export type GetNetworksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { __typename?: 'Query', getNetworks: { __typename?: 'NetworksResDto', orgId: string, networks: Array<{ __typename?: 'NetworkDto', id: string, name: string, orgId: string, isDeactivated: string, createdAt: string }> } };

export type GetSitesQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSitesQuery = { __typename?: 'Query', getSites: { __typename?: 'SitesResDto', networkId: string, sites: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, isDeactivated: string, createdAt: string }> } };

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
  status {
    connectivity
    state
  }
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
export const OrgUserFragmentDoc = gql`
    fragment OrgUser on UserResDto {
  name
  email
  uuid
  phone
  isDeactivated
  registeredSince
}
    `;
export const MemberFragmentDoc = gql`
    fragment Member on MemberObj {
  uuid
  userId
  orgId
  role
  isDeactivated
  memberSince
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
export const GetOrgMemberDocument = gql`
    query getOrgMember {
  getOrgMembers {
    org
    members {
      ...Member
      user {
        ...OrgUser
      }
    }
  }
}
    ${MemberFragmentDoc}
${OrgUserFragmentDoc}`;

/**
 * __useGetOrgMemberQuery__
 *
 * To run a query within a React component, call `useGetOrgMemberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgMemberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgMemberQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgMemberQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgMemberQuery, GetOrgMemberQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgMemberQuery, GetOrgMemberQueryVariables>(GetOrgMemberDocument, options);
      }
export function useGetOrgMemberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgMemberQuery, GetOrgMemberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgMemberQuery, GetOrgMemberQueryVariables>(GetOrgMemberDocument, options);
        }
export type GetOrgMemberQueryHookResult = ReturnType<typeof useGetOrgMemberQuery>;
export type GetOrgMemberLazyQueryHookResult = ReturnType<typeof useGetOrgMemberLazyQuery>;
export type GetOrgMemberQueryResult = Apollo.QueryResult<GetOrgMemberQuery, GetOrgMemberQueryVariables>;
export const AddMemberDocument = gql`
    mutation addMember($data: AddMemberInputDto!) {
  addMember(data: $data) {
    ...Member
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
export const DeletePacakgeDocument = gql`
    mutation deletePacakge($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}
    `;
export type DeletePacakgeMutationFn = Apollo.MutationFunction<DeletePacakgeMutation, DeletePacakgeMutationVariables>;

/**
 * __useDeletePacakgeMutation__
 *
 * To run a mutation, you first call `useDeletePacakgeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeletePacakgeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deletePacakgeMutation, { data, loading, error }] = useDeletePacakgeMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useDeletePacakgeMutation(baseOptions?: Apollo.MutationHookOptions<DeletePacakgeMutation, DeletePacakgeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeletePacakgeMutation, DeletePacakgeMutationVariables>(DeletePacakgeDocument, options);
      }
export type DeletePacakgeMutationHookResult = ReturnType<typeof useDeletePacakgeMutation>;
export type DeletePacakgeMutationResult = Apollo.MutationResult<DeletePacakgeMutation>;
export type DeletePacakgeMutationOptions = Apollo.BaseMutationOptions<DeletePacakgeMutation, DeletePacakgeMutationVariables>;
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
      id
      name
      orgId
      isDeactivated
      createdAt
    }
  }
}
    `;

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
export const GetSitesDocument = gql`
    query getSites($networkId: String!) {
  getSites(networkId: $networkId) {
    networkId
    sites {
      id
      name
      networkId
      isDeactivated
      createdAt
    }
  }
}
    `;

/**
 * __useGetSitesQuery__
 *
 * To run a query within a React component, call `useGetSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSitesQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSitesQuery(baseOptions: Apollo.QueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
      }
export function useGetSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
export type GetSitesQueryHookResult = ReturnType<typeof useGetSitesQuery>;
export type GetSitesLazyQueryHookResult = ReturnType<typeof useGetSitesLazyQuery>;
export type GetSitesQueryResult = Apollo.QueryResult<GetSitesQuery, GetSitesQueryVariables>;
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