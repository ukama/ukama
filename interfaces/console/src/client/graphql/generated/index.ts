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
};

export type AddMemberInputDto = {
  role: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type AddNetworkInputDto = {
  budget?: InputMaybe<Scalars['Float']['input']>;
  countries?: InputMaybe<Array<Scalars['String']['input']>>;
  isDefault?: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  networks?: InputMaybe<Array<Scalars['String']['input']>>;
};

export type AddNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type AddNodeToSiteInput = {
  networkId: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export type AddPackagSimResDto = {
  __typename?: 'AddPackagSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type AddPackageInputDto = {
  amount: Scalars['Float']['input'];
  country: Scalars['String']['input'];
  currency: Scalars['String']['input'];
  dataUnit: Scalars['String']['input'];
  dataVolume: Scalars['Int']['input'];
  duration: Scalars['Int']['input'];
  name: Scalars['String']['input'];
};

export type AddPackagesSimResDto = {
  __typename?: 'AddPackagesSimResDto';
  packages: Array<AddPackagSimResDto>;
};

export type AddPackagesToSimInputDto = {
  packages: Array<PackagesToSimInputDto>;
  sim_id: Scalars['String']['input'];
};

export type AddSiteInputDto = {
  access_id: Scalars['String']['input'];
  backhaul_id: Scalars['String']['input'];
  install_date: Scalars['String']['input'];
  latitude: Scalars['String']['input'];
  location: Scalars['String']['input'];
  longitude: Scalars['String']['input'];
  name: Scalars['String']['input'];
  network_id: Scalars['String']['input'];
  power_id: Scalars['String']['input'];
  spectrum_id: Scalars['String']['input'];
  switch_id: Scalars['String']['input'];
};

export type AllocateSimApiDto = {
  __typename?: 'AllocateSimAPIDto';
  allocated_at: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi?: Maybe<Scalars['String']['output']>;
  is_physical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
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
  latitude: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
  site: NodeSite;
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type CBooleanResponse = {
  __typename?: 'CBooleanResponse';
  success: Scalars['Boolean']['output'];
};

export enum Component_Type {
  Access = 'access',
  All = 'all',
  Backhaul = 'backhaul',
  Power = 'power',
  Spectrum = 'spectrum',
  Switch = 'switch'
}

export type Component = {
  __typename?: 'Component';
  componentId?: Maybe<Scalars['String']['output']>;
  componentName?: Maybe<Scalars['String']['output']>;
  elementType: Scalars['String']['output'];
};

export type ComponentDto = {
  __typename?: 'ComponentDto';
  category: Scalars['String']['output'];
  datasheetUrl: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imageUrl: Scalars['String']['output'];
  inventoryId: Scalars['String']['output'];
  managed: Scalars['String']['output'];
  manufacturer: Scalars['String']['output'];
  partNumber: Scalars['String']['output'];
  specification: Scalars['String']['output'];
  type: Scalars['String']['output'];
  userId: Scalars['String']['output'];
  warranty: Scalars['Float']['output'];
};

export type ComponentTypeInputDto = {
  category: Component_Type;
};

export type ComponentsResDto = {
  __typename?: 'ComponentsResDto';
  components: Array<ComponentDto>;
};

export type CountriesRes = {
  __typename?: 'CountriesRes';
  countries: Array<CountryDto>;
};

export type CountryDto = {
  __typename?: 'CountryDto';
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type CreateInvitationInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  role: Role_Type;
};

export type CurrencyRes = {
  __typename?: 'CurrencyRes';
  code: Scalars['String']['output'];
  image: Scalars['String']['output'];
  symbol: Scalars['String']['output'];
};

export type CustomerDto = {
  __typename?: 'CustomerDto';
  addressLine1?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  email?: Maybe<Scalars['String']['output']>;
  externalId: Scalars['String']['output'];
  legalName?: Maybe<Scalars['String']['output']>;
  legalNumber?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  phone?: Maybe<Scalars['String']['output']>;
  timezone?: Maybe<Scalars['String']['output']>;
  vatRate: Scalars['Float']['output'];
};

export type DataPlan = {
  __typename?: 'DataPlan';
  elementType: Scalars['String']['output'];
  planId: Scalars['String']['output'];
  planName: Scalars['String']['output'];
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

export type DeleteInvitationResDto = {
  __typename?: 'DeleteInvitationResDto';
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

export type FeeDto = {
  __typename?: 'FeeDto';
  eventsCount: Scalars['String']['output'];
  item: ItemResDto;
  taxesAmountCents: Scalars['String']['output'];
  taxesPreciseAmount: Scalars['String']['output'];
  totalAmountCents: Scalars['String']['output'];
  totalAmountCurrency: Scalars['String']['output'];
  units: Scalars['Float']['output'];
};

export type GetNodeLatestMetricInput = {
  nodeId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type GetNodesByStateInput = {
  connectivity: NodeConnectivityEnum;
  state: NodeStateEnum;
};

export type GetPackagesForSimInputDto = {
  sim_id: Scalars['String']['input'];
};

export type GetPaymentsInputDto = {
  paymentMethod?: InputMaybe<Scalars['String']['input']>;
  status?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
};

export type GetPdfReportUrlDto = {
  __typename?: 'GetPdfReportUrlDto';
  contentType: Scalars['String']['output'];
  downloadUrl: Scalars['String']['output'];
  filename: Scalars['String']['output'];
  id: Scalars['String']['output'];
};

export type GetReportDto = {
  __typename?: 'GetReportDto';
  report: ReportDto;
};

export type GetReportsDto = {
  __typename?: 'GetReportsDto';
  reports: Array<ReportDto>;
};

export type GetReportsInputDto = {
  count?: InputMaybe<Scalars['Float']['input']>;
  isPaid?: InputMaybe<Scalars['Boolean']['input']>;
  networkId?: InputMaybe<Scalars['String']['input']>;
  ownerId?: InputMaybe<Scalars['String']['input']>;
  ownerType: Scalars['String']['input'];
  report_type: Scalars['String']['input'];
  sort?: InputMaybe<Scalars['Boolean']['input']>;
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

export type GetSimsInput = {
  status: Sim_Status;
  type: Sim_Types;
};

export enum Invitation_Status {
  InviteAccepted = 'INVITE_ACCEPTED',
  InviteDeclined = 'INVITE_DECLINED',
  InvitePending = 'INVITE_PENDING'
}

export type IdResponse = {
  __typename?: 'IdResponse';
  uuid: Scalars['String']['output'];
};

export type InvitationDto = {
  __typename?: 'InvitationDto';
  email: Scalars['String']['output'];
  expireAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  link: Scalars['String']['output'];
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
  status: Invitation_Status;
  userId: Scalars['String']['output'];
};

export type InvitationsResDto = {
  __typename?: 'InvitationsResDto';
  invitations: Array<InvitationDto>;
};

export type ItemResDto = {
  __typename?: 'ItemResDto';
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type ListSimsInput = {
  networkId: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type MemberDto = {
  __typename?: 'MemberDto';
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  memberId: Scalars['String']['output'];
  memberSince?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
  userId: Scalars['String']['output'];
};

export type Members = {
  __typename?: 'Members';
  activeMembers: Scalars['String']['output'];
  inactiveMembers: Scalars['String']['output'];
  totalMembers: Scalars['String']['output'];
};

export type MembersResDto = {
  __typename?: 'MembersResDto';
  members: Array<MemberDto>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addMember: MemberDto;
  addNetwork: NetworkDto;
  addNode: Node;
  addNodeToSite: CBooleanResponse;
  addPackage: PackageDto;
  addPackagesToSim: AddPackagesSimResDto;
  addSite: SiteDto;
  addSubscriber: SubscriberDto;
  allocateSim: AllocateSimApiDto;
  attachNode: CBooleanResponse;
  createInvitation: InvitationDto;
  defaultMarkup: CBooleanResponse;
  deleteInvitation: DeleteInvitationResDto;
  deleteNodeFromOrg: DeleteNode;
  deletePackage: IdResponse;
  deleteSim: DeleteSimResDto;
  deleteSubscriber: CBooleanResponse;
  detachhNode: CBooleanResponse;
  processPayment: ProcessPaymentDto;
  releaseNodeFromSite: CBooleanResponse;
  removeMember: CBooleanResponse;
  removePackageForSim: RemovePackageFromSimResDto;
  restartNode: CBooleanResponse;
  restartNodes: CBooleanResponse;
  restartSite: CBooleanResponse;
  setDefaultNetwork: CBooleanResponse;
  toggleInternetSwitch: CBooleanResponse;
  toggleRFStatus: CBooleanResponse;
  toggleSimStatus: SimStatusResDto;
  updateFirstVisit: UserFistVisitResDto;
  updateInvitation: UpdateInvitationResDto;
  updateMember: CBooleanResponse;
  updateNode: Node;
  updateNodeState: Node;
  updateNotification: UpdateNotificationResDto;
  updatePackage: PackageDto;
  updatePayment: PaymentDto;
  updateSite: SiteDto;
  updateSubscriber: CBooleanResponse;
  uploadSims: UploadSimsResDto;
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


export type MutationAddPackagesToSimArgs = {
  data: AddPackagesToSimInputDto;
};


export type MutationAddSiteArgs = {
  data: AddSiteInputDto;
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


export type MutationCreateInvitationArgs = {
  data: CreateInvitationInputDto;
};


export type MutationDefaultMarkupArgs = {
  data: DefaultMarkupInputDto;
};


export type MutationDeleteInvitationArgs = {
  id: Scalars['String']['input'];
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


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type MutationDetachhNodeArgs = {
  data: NodeInput;
};


export type MutationProcessPaymentArgs = {
  data: ProcessPaymentInputDto;
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


export type MutationRestartNodeArgs = {
  data: RestartNodeInputDto;
};


export type MutationRestartNodesArgs = {
  data: RestartNodesInputDto;
};


export type MutationRestartSiteArgs = {
  data: RestartSiteInputDto;
};


export type MutationSetDefaultNetworkArgs = {
  data: SetDefaultNetworkInputDto;
};


export type MutationToggleInternetSwitchArgs = {
  data: ToggleInternetSwitchInputDto;
};


export type MutationToggleRfStatusArgs = {
  data: ToggleRfStatusInputDto;
};


export type MutationToggleSimStatusArgs = {
  data: ToggleSimStatusInputDto;
};


export type MutationUpdateFirstVisitArgs = {
  data: UserFistVisitInputDto;
};


export type MutationUpdateInvitationArgs = {
  data: UpateInvitationInputDto;
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


export type MutationUpdateNotificationArgs = {
  id: Scalars['String']['input'];
  isRead: Scalars['Boolean']['input'];
};


export type MutationUpdatePackageArgs = {
  data: UpdatePackageInputDto;
  packageId: Scalars['String']['input'];
};


export type MutationUpdatePaymentArgs = {
  data: UpdatePaymentInputDto;
};


export type MutationUpdateSiteArgs = {
  data: UpdateSiteInputDto;
  siteId: Scalars['String']['input'];
};


export type MutationUpdateSubscriberArgs = {
  data: UpdateSubscriberInputDto;
  subscriberId: Scalars['String']['input'];
};


export type MutationUploadSimsArgs = {
  data: UploadSimsInputDto;
};

export enum Notification_Scope {
  ScopeInvalid = 'SCOPE_INVALID',
  ScopeNetwork = 'SCOPE_NETWORK',
  ScopeNetworks = 'SCOPE_NETWORKS',
  ScopeNode = 'SCOPE_NODE',
  ScopeOrg = 'SCOPE_ORG',
  ScopeOwner = 'SCOPE_OWNER',
  ScopeSite = 'SCOPE_SITE',
  ScopeSites = 'SCOPE_SITES',
  ScopeSubscriber = 'SCOPE_SUBSCRIBER',
  ScopeSubscribers = 'SCOPE_SUBSCRIBERS',
  ScopeUser = 'SCOPE_USER',
  ScopeUsers = 'SCOPE_USERS'
}

export enum Notification_Type {
  TypeActionableCritical = 'TYPE_ACTIONABLE_CRITICAL',
  TypeActionableError = 'TYPE_ACTIONABLE_ERROR',
  TypeActionableInfo = 'TYPE_ACTIONABLE_INFO',
  TypeActionableWarning = 'TYPE_ACTIONABLE_WARNING',
  TypeCritical = 'TYPE_CRITICAL',
  TypeError = 'TYPE_ERROR',
  TypeInfo = 'TYPE_INFO',
  TypeInvalid = 'TYPE_INVALID',
  TypeWarning = 'TYPE_WARNING'
}

export type Network = {
  __typename?: 'Network';
  elementType: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  networkName: Scalars['String']['output'];
  sites: Array<Site>;
  subscribers?: Maybe<Subscribers>;
};

export type NetworkDto = {
  __typename?: 'NetworkDto';
  budget: Scalars['Float']['output'];
  countries: Array<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  isDefault: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  networks: Array<Scalars['String']['output']>;
  overdraft: Scalars['Float']['output'];
  paymentLinks: Scalars['Boolean']['output'];
  trafficPolicy: Scalars['Float']['output'];
};

export type NetworkStats = {
  __typename?: 'NetworkStats';
  activeSubscriber: Scalars['Float']['output'];
  averageSignalStrength: Scalars['Float']['output'];
  averageThroughput: Scalars['Float']['output'];
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
};

export type Node = {
  __typename?: 'Node';
  attached: Array<AttachedNodes>;
  id: Scalars['String']['output'];
  latitude: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
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

/** Node connectivity enums */
export enum NodeConnectivityEnum {
  Offline = 'Offline',
  Online = 'Online',
  Unknown = 'Unknown'
}

export type NodeInput = {
  id: Scalars['String']['input'];
};

export type NodeLatestMetric = {
  __typename?: 'NodeLatestMetric';
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type NodeSite = {
  __typename?: 'NodeSite';
  addedAt?: Maybe<Scalars['String']['output']>;
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
};

/** Node state enums */
export enum NodeStateEnum {
  Configured = 'Configured',
  Faulty = 'Faulty',
  Operational = 'Operational',
  Unknown = 'Unknown'
}

export type NodeStateRes = {
  __typename?: 'NodeStateRes';
  createdAt: Scalars['String']['output'];
  currentState: NodeStateEnum;
  id: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  previousState?: Maybe<NodeStateEnum>;
  previousStateId?: Maybe<Scalars['String']['output']>;
};

export type NodeStatus = {
  __typename?: 'NodeStatus';
  connectivity: Scalars['String']['output'];
  state: Scalars['String']['output'];
};

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

export type NodesFilterInput = {
  connectivity?: InputMaybe<Scalars['String']['input']>;
  id?: InputMaybe<Scalars['String']['input']>;
  networkId?: InputMaybe<Scalars['String']['input']>;
  siteId?: InputMaybe<Scalars['String']['input']>;
  state?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
};

export type NotificationResDto = {
  __typename?: 'NotificationResDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  resourceId: Scalars['String']['output'];
  scope: Notification_Scope;
  subscriberId: Scalars['String']['output'];
  title: Scalars['String']['output'];
  type: Notification_Type;
  userId: Scalars['String']['output'];
};

export type NotificationsDto = {
  __typename?: 'NotificationsDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  scope: Notification_Scope;
  title: Scalars['String']['output'];
  type: Notification_Type;
};

export type NotificationsResDto = {
  __typename?: 'NotificationsResDto';
  notifications: Array<NotificationsDto>;
};

export type Org = {
  __typename?: 'Org';
  country: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  dataplans: Array<DataPlan>;
  elementType: Scalars['String']['output'];
  members?: Maybe<Members>;
  networks: Array<Network>;
  orgId: Scalars['String']['output'];
  orgName: Scalars['String']['output'];
  ownerEmail: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  ownerName: Scalars['String']['output'];
  sims?: Maybe<Sims>;
};

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgTreeRes = {
  __typename?: 'OrgTreeRes';
  org: Org;
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
  dataVolume: Scalars['Float']['output'];
  deletedAt: Scalars['String']['output'];
  dlbr: Scalars['String']['output'];
  duration: Scalars['Float']['output'];
  flatrate: Scalars['Boolean']['output'];
  from: Scalars['String']['output'];
  markup: PackageMarkupApiDto;
  messageUnit: Scalars['String']['output'];
  name: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  provider: Scalars['String']['output'];
  rate: PackageRateApiDto;
  simType: Scalars['String']['output'];
  smsVolume: Scalars['Float']['output'];
  to: Scalars['String']['output'];
  type: Scalars['String']['output'];
  ulbr: Scalars['String']['output'];
  updatedAt: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
  voiceUnit: Scalars['String']['output'];
  voiceVolume: Scalars['Float']['output'];
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

export type PackagesToSimInputDto = {
  package_id: Scalars['String']['input'];
  start_date: Scalars['String']['input'];
};

export type PaymentDto = {
  __typename?: 'PaymentDto';
  amount: Scalars['String']['output'];
  correspondent: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  depositedAmount: Scalars['String']['output'];
  description: Scalars['String']['output'];
  extra: Scalars['String']['output'];
  failureReason: Scalars['String']['output'];
  id: Scalars['String']['output'];
  itemId: Scalars['String']['output'];
  itemType: Scalars['String']['output'];
  paidAt: Scalars['String']['output'];
  payerEmail: Scalars['String']['output'];
  payerName: Scalars['String']['output'];
  payerPhone: Scalars['String']['output'];
  paymentMethod: Scalars['String']['output'];
  status: Scalars['String']['output'];
};

export type PaymentsDto = {
  __typename?: 'PaymentsDto';
  payments: Array<PaymentDto>;
};

export type ProcessPaymentDto = {
  __typename?: 'ProcessPaymentDto';
  payment: PaymentDto;
};

export type ProcessPaymentInputDto = {
  correspondent?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['String']['input'];
  token: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  example?: Maybe<Scalars['String']['output']>;
  getAppsChangeLog: AppChangeLogs;
  getComponentById: ComponentDto;
  getComponentsByUserId: ComponentsResDto;
  getCountries: CountriesRes;
  getCurrencySymbol: CurrencyRes;
  getDataUsage: SimDataUsage;
  getDataUsages: SimDataUsages;
  getDefaultMarkup: DefaultMarkupResDto;
  getDefaultMarkupHistory: DefaultMarkupHistoryResDto;
  getGeneratedPdfReport: GetPdfReportUrlDto;
  getInvitation: InvitationDto;
  getInvitations: InvitationsResDto;
  getInvitationsByEmail: InvitationsResDto;
  getMember: MemberDto;
  getMemberByUserId: MemberDto;
  getMembers: MembersResDto;
  getNetwork: NetworkDto;
  getNetworkStats: NetworkStats;
  getNetworks: NetworksResDto;
  getNode: Node;
  getNodeApps: NodeApps;
  getNodeLatestMetric: NodeLatestMetric;
  getNodeState: NodeStateRes;
  getNodes: Nodes;
  getNodesByNetwork: Nodes;
  getNodesByState: Nodes;
  getNodesForSite: Nodes;
  getNodesLocation: Nodes;
  getNotification: NotificationResDto;
  getNotifications: NotificationsResDto;
  getOrg: OrgDto;
  getOrgTree: OrgTreeRes;
  getOrgs: OrgsResDto;
  getPackage: PackageDto;
  getPackages: PackagesResDto;
  getPackagesForSim: GetSimPackagesDtoApi;
  getPayment: PaymentDto;
  getPayments: PaymentsDto;
  getReport: GetReportDto;
  getReportPdf: GetReportDto;
  getReports: GetReportsDto;
  getSim: SimDto;
  getSimPoolStats: SimPoolStatsDto;
  getSims: SimsResDto;
  getSimsByNetwork: SubscriberSimsResDto;
  getSimsBySubscriber: SubscriberToSimsDto;
  getSimsFromPool: SimsPoolResDto;
  getSite: SiteDto;
  getSites: SitesResDto;
  getSubscriber: SubscriberDto;
  getSubscriberMetricsByNetwork: SubscriberMetricsByNetworkDto;
  getSubscribersByNetwork: SubscribersResDto;
  getTimezones: TimezoneRes;
  getToken: TokenResDto;
  getUser: UserResDto;
  whoami: WhoamiDto;
};


export type QueryGetAppsChangeLogArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetComponentByIdArgs = {
  componentId: Scalars['String']['input'];
};


export type QueryGetComponentsByUserIdArgs = {
  data: ComponentTypeInputDto;
};


export type QueryGetCurrencySymbolArgs = {
  code: Scalars['String']['input'];
};


export type QueryGetDataUsageArgs = {
  data: SimUsageInputDto;
};


export type QueryGetDataUsagesArgs = {
  data: SimUsagesInputDto;
};


export type QueryGetGeneratedPdfReportArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetInvitationArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetInvitationsByEmailArgs = {
  email: Scalars['String']['input'];
};


export type QueryGetMemberArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetMemberByUserIdArgs = {
  userId: Scalars['String']['input'];
};


export type QueryGetNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNetworkStatsArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  data: NodeInput;
};


export type QueryGetNodeAppsArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetNodeLatestMetricArgs = {
  data: GetNodeLatestMetricInput;
};


export type QueryGetNodeStateArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetNodesArgs = {
  data: NodesFilterInput;
};


export type QueryGetNodesByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodesByStateArgs = {
  data: GetNodesByStateInput;
};


export type QueryGetNodesForSiteArgs = {
  siteId: Scalars['String']['input'];
};


export type QueryGetNotificationArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetPackageArgs = {
  packageId: Scalars['String']['input'];
};


export type QueryGetPackagesForSimArgs = {
  data: GetPackagesForSimInputDto;
};


export type QueryGetPaymentArgs = {
  paymentId: Scalars['String']['input'];
};


export type QueryGetPaymentsArgs = {
  data: GetPaymentsInputDto;
};


export type QueryGetReportArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetReportPdfArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetReportsArgs = {
  data: GetReportsInputDto;
};


export type QueryGetSimArgs = {
  data: GetSimInputDto;
};


export type QueryGetSimPoolStatsArgs = {
  data: GetSimsInput;
};


export type QueryGetSimsArgs = {
  data: ListSimsInput;
};


export type QueryGetSimsByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetSimsBySubscriberArgs = {
  data: GetSimBySubscriberInputDto;
};


export type QueryGetSimsFromPoolArgs = {
  data: GetSimsInput;
};


export type QueryGetSiteArgs = {
  siteId: Scalars['String']['input'];
};


export type QueryGetSitesArgs = {
  data: SitesInputDto;
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


export type QueryGetTokenArgs = {
  paymentId: Scalars['String']['input'];
};


export type QueryGetUserArgs = {
  userId: Scalars['String']['input'];
};

export enum Role_Type {
  RoleAdmin = 'ROLE_ADMIN',
  RoleInvalid = 'ROLE_INVALID',
  RoleNetworkOwner = 'ROLE_NETWORK_OWNER',
  RoleOwner = 'ROLE_OWNER',
  RoleUser = 'ROLE_USER',
  RoleVendor = 'ROLE_VENDOR'
}

export type RawReportDto = {
  __typename?: 'RawReportDto';
  currency: Scalars['String']['output'];
  customer: CustomerDto;
  fees: Array<FeeDto>;
  feesAmountCents: Scalars['String']['output'];
  fileUrl: Scalars['String']['output'];
  invoiceType: Scalars['String']['output'];
  issuingDate: Scalars['String']['output'];
  paymentDueDate: Scalars['String']['output'];
  paymentOverdue: Scalars['Boolean']['output'];
  paymentStatus: Scalars['String']['output'];
  status: Scalars['String']['output'];
  subTotalExcludingTaxesAmountCents: Scalars['String']['output'];
  subTotalIncludingTaxesAmountCents: Scalars['String']['output'];
  subscriptions: Array<SubscriptionDto>;
  taxesAmountCents: Scalars['String']['output'];
  totalAmountCents: Scalars['String']['output'];
  vatAmountCents: Scalars['String']['output'];
  vatAmountCurrency?: Maybe<Scalars['String']['output']>;
};

export type RemovePackageFormSimInputDto = {
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type RemovePackageFromSimResDto = {
  __typename?: 'RemovePackageFromSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type ReportDto = {
  __typename?: 'ReportDto';
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isPaid: Scalars['Boolean']['output'];
  networkId: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  ownerType: Scalars['String']['output'];
  period: Scalars['String']['output'];
  rawReport: RawReportDto;
  type: Scalars['String']['output'];
};

export type RestartNodeInputDto = {
  nodeId: Scalars['String']['input'];
};

export type RestartNodesInputDto = {
  networkId: Scalars['String']['input'];
  nodeIds: Array<Scalars['String']['input']>;
};

export type RestartSiteInputDto = {
  networkId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export enum Sim_Status {
  All = 'ALL',
  Assigned = 'ASSIGNED',
  Unassigned = 'UNASSIGNED'
}

export enum Sim_Types {
  OperatorData = 'operator_data',
  Test = 'test',
  UkamaData = 'ukama_data',
  Unknown = 'unknown'
}

export type SetDefaultNetworkInputDto = {
  id: Scalars['String']['input'];
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
  simId: Scalars['String']['output'];
  usage: Scalars['String']['output'];
};

export type SimDataUsages = {
  __typename?: 'SimDataUsages';
  usages: Array<SimDataUsage>;
};

export type SimDto = {
  __typename?: 'SimDto';
  activationsCount: Scalars['String']['output'];
  allocatedAt: Scalars['String']['output'];
  deactivationsCount: Scalars['String']['output'];
  firstActivatedOn: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  lastActivatedOn: Scalars['String']['output'];
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  package?: Maybe<SimPackage>;
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  syncStatus: Scalars['String']['output'];
  trafficPolicy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type SimPackage = {
  __typename?: 'SimPackage';
  asExpired: Scalars['Boolean']['output'];
  defaultDuration: Scalars['String']['output'];
  endDate: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isActive: Scalars['Boolean']['output'];
  packageId: Scalars['String']['output'];
  startDate: Scalars['String']['output'];
};

export type SimPackageDto = {
  __typename?: 'SimPackageDto';
  created_at: Scalars['String']['output'];
  end_date: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_active: Scalars['Boolean']['output'];
  package_id: Scalars['String']['output'];
  start_date: Scalars['String']['output'];
  updated_at: Scalars['String']['output'];
};

export type SimPoolResDto = {
  __typename?: 'SimPoolResDto';
  activationCode: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  deletedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isAllocated: Scalars['Boolean']['output'];
  isFailed: Scalars['Boolean']['output'];
  isPhysical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  qrCode: Scalars['String']['output'];
  simType: Scalars['String']['output'];
  smApAddress: Scalars['String']['output'];
  updatedAt: Scalars['String']['output'];
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

export type SimUsageInputDto = {
  iccid: Scalars['String']['input'];
  simId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type SimUsagesInputDto = {
  networkId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type Sims = {
  __typename?: 'Sims';
  availableSims: Scalars['String']['output'];
  consumed: Scalars['String']['output'];
  totalSims: Scalars['String']['output'];
};

export type SimsPoolResDto = {
  __typename?: 'SimsPoolResDto';
  sims: Array<SimPoolResDto>;
};

export type SimsResDto = {
  __typename?: 'SimsResDto';
  sims: Array<SimDto>;
};

export type Site = {
  __typename?: 'Site';
  components: Array<Component>;
  elementType: Scalars['String']['output'];
  siteId: Scalars['String']['output'];
  siteName: Scalars['String']['output'];
};

export type SiteDto = {
  __typename?: 'SiteDto';
  accessId: Scalars['String']['output'];
  backhaulId: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  installDate: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  latitude: Scalars['String']['output'];
  location: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  powerId: Scalars['String']['output'];
  spectrumId: Scalars['String']['output'];
  switchId: Scalars['String']['output'];
};

export type SitesInputDto = {
  networkId?: InputMaybe<Scalars['String']['input']>;
};

export type SitesResDto = {
  __typename?: 'SitesResDto';
  sites: Array<SiteDto>;
};

export type SubscriberDto = {
  __typename?: 'SubscriberDto';
  address: Scalars['String']['output'];
  dob: Scalars['String']['output'];
  email: Scalars['String']['output'];
  gender: Scalars['String']['output'];
  idSerial: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  proofOfIdentification: Scalars['String']['output'];
  sim?: Maybe<Array<SubscriberSimDto>>;
  uuid: Scalars['String']['output'];
};

export type SubscriberInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  network_id: Scalars['String']['input'];
  phone?: InputMaybe<Scalars['String']['input']>;
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
  allocatedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical?: Maybe<Scalars['Boolean']['output']>;
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  package?: Maybe<SimPackageDto>;
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  sync_status?: Maybe<Scalars['String']['output']>;
  type: Scalars['String']['output'];
};

export type SubscriberSimsDto = {
  __typename?: 'SubscriberSimsDto';
  allocatedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  syncStatus: Scalars['String']['output'];
  trafficPolicy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type SubscriberSimsResDto = {
  __typename?: 'SubscriberSimsResDto';
  sims: Array<SubscriberSimDto>;
};

export type SubscriberToSimsDto = {
  __typename?: 'SubscriberToSimsDto';
  sims: Array<SubscriberSimsDto>;
  subscriberId: Scalars['String']['output'];
};

export type Subscribers = {
  __typename?: 'Subscribers';
  activeSubscribers: Scalars['String']['output'];
  inactiveSubscribers: Scalars['String']['output'];
  totalSubscribers: Scalars['String']['output'];
};

export type SubscribersResDto = {
  __typename?: 'SubscribersResDto';
  subscribers: Array<SubscriberDto>;
};

export type SubscriptionDto = {
  __typename?: 'SubscriptionDto';
  canceledAt?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  externalCustomerId: Scalars['String']['output'];
  externalId: Scalars['String']['output'];
  name?: Maybe<Scalars['String']['output']>;
  planCode: Scalars['String']['output'];
  startedAt: Scalars['String']['output'];
  status: Scalars['String']['output'];
  terminatedAt?: Maybe<Scalars['String']['output']>;
};

export type TimezoneDto = {
  __typename?: 'TimezoneDto';
  abbr: Scalars['String']['output'];
  isdst: Scalars['Boolean']['output'];
  offset: Scalars['Float']['output'];
  text: Scalars['String']['output'];
  utc: Array<Scalars['String']['output']>;
  value: Scalars['String']['output'];
};

export type TimezoneRes = {
  __typename?: 'TimezoneRes';
  timezones: Array<TimezoneDto>;
};

export type ToggleInternetSwitchInputDto = {
  port: Scalars['Float']['input'];
  siteId: Scalars['String']['input'];
  status: Scalars['Boolean']['input'];
};

export type ToggleRfStatusInputDto = {
  nodeId: Scalars['String']['input'];
  status: Scalars['Boolean']['input'];
};

export type ToggleSimStatusInputDto = {
  sim_id: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type TokenResDto = {
  __typename?: 'TokenResDto';
  token: Scalars['String']['output'];
};

export type UpateInvitationInputDto = {
  email: Scalars['String']['input'];
  id: Scalars['String']['input'];
  status: Invitation_Status;
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
  state: NodeStateEnum;
};

export type UpdateNotificationResDto = {
  __typename?: 'UpdateNotificationResDto';
  id: Scalars['String']['output'];
};

export type UpdatePackageInputDto = {
  active: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
};

export type UpdatePaymentInputDto = {
  country?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['String']['input'];
  payerEmail?: InputMaybe<Scalars['String']['input']>;
  payerName?: InputMaybe<Scalars['String']['input']>;
  paymentMethod?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateSiteInputDto = {
  name: Scalars['String']['input'];
};

export type UpdateSubscriberInputDto = {
  address?: InputMaybe<Scalars['String']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  id_serial?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  proof_of_identification?: InputMaybe<Scalars['String']['input']>;
};

export type UploadSimsInputDto = {
  data: Scalars['String']['input'];
  simType: Sim_Types;
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

export type NodeFragment = { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } };

export type GetNodeQueryVariables = Exact<{
  data: NodeInput;
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesQueryVariables = Exact<{
  data: NodesFilterInput;
}>;


export type GetNodesQuery = { __typename?: 'Query', getNodes: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

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


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

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


export type UpdateNodeStateMutation = { __typename?: 'Mutation', updateNodeState: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesForSiteQueryVariables = Exact<{
  siteId: Scalars['String']['input'];
}>;


export type GetNodesForSiteQuery = { __typename?: 'Query', getNodesForSite: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type UpdateNodeMutationVariables = Exact<{
  data: UpdateNodeInput;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodeAppsQueryVariables = Exact<{
  data: NodeAppsChangeLogInput;
}>;


export type GetNodeAppsQuery = { __typename?: 'Query', getNodeApps: { __typename?: 'NodeApps', type: NodeTypeEnum, apps: Array<{ __typename?: 'NodeApp', name: string, date: number, version: string, cpu: string, memory: string, notes: string }> } };

export type GetNodeStateQueryVariables = Exact<{
  getNodeStateId: Scalars['String']['input'];
}>;


export type GetNodeStateQuery = { __typename?: 'Query', getNodeState: { __typename?: 'NodeStateRes', id: string, nodeId: string, previousStateId?: string | null, previousState?: NodeStateEnum | null, currentState: NodeStateEnum, createdAt: string } };

export type RestartNodeMutationVariables = Exact<{
  data: RestartNodeInputDto;
}>;


export type RestartNodeMutation = { __typename?: 'Mutation', restartNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleInternetSwitchMutationVariables = Exact<{
  data: ToggleInternetSwitchInputDto;
}>;


export type ToggleInternetSwitchMutation = { __typename?: 'Mutation', toggleInternetSwitch: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleRfStatusMutationVariables = Exact<{
  data: ToggleRfStatusInputDto;
}>;


export type ToggleRfStatusMutation = { __typename?: 'Mutation', toggleRFStatus: { __typename?: 'CBooleanResponse', success: boolean } };

export type MemberFragment = { __typename?: 'MemberDto', role: string, userId: string, isDeactivated: boolean, memberSince?: string | null, id: string };

export type GetMembersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetMembersQuery = { __typename?: 'Query', getMembers: { __typename?: 'MembersResDto', members: Array<{ __typename?: 'MemberDto', name: string, email: string, role: string, userId: string, isDeactivated: boolean, memberSince?: string | null, id: string }> } };

export type GetMemberQueryVariables = Exact<{
  memberId: Scalars['String']['input'];
}>;


export type GetMemberQuery = { __typename?: 'Query', getMember: { __typename?: 'MemberDto', role: string, userId: string, isDeactivated: boolean, memberSince?: string | null, id: string } };

export type AddMemberMutationVariables = Exact<{
  data: AddMemberInputDto;
}>;


export type AddMemberMutation = { __typename?: 'Mutation', addMember: { __typename?: 'MemberDto', role: string, userId: string, isDeactivated: boolean, memberSince?: string | null, id: string } };

export type RemoveMemberMutationVariables = Exact<{
  memberId: Scalars['String']['input'];
}>;


export type RemoveMemberMutation = { __typename?: 'Mutation', removeMember: { __typename?: 'CBooleanResponse', success: boolean } };

export type UpdateMemberMutationVariables = Exact<{
  memberId: Scalars['String']['input'];
  data: UpdateMemberInputDto;
}>;


export type UpdateMemberMutation = { __typename?: 'Mutation', updateMember: { __typename?: 'CBooleanResponse', success: boolean } };

export type GetMemberByUserIdQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetMemberByUserIdQuery = { __typename?: 'Query', getMemberByUserId: { __typename?: 'MemberDto', userId: string, name: string, email: string, memberId: string, isDeactivated: boolean, role: string, memberSince?: string | null } };

export type OrgFragment = { __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean };

export type GetOrgsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgsQuery = { __typename?: 'Query', getOrgs: { __typename?: 'OrgsResDto', user: string, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }> } };

export type GetOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgQuery = { __typename?: 'Query', getOrg: { __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean } };

export type PackageRateFragment = { __typename?: 'PackageDto', rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number } };

export type PackageMarkupFragment = { __typename?: 'PackageDto', markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type SimPackagesFragment = { __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean };

export type SubscriberSimsFragment = { __typename?: 'SubscriberToSimsDto', subscriberId: string, sims: Array<{ __typename?: 'SubscriberSimsDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> };

export type PackageFragment = { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type GetPackagesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { __typename?: 'Query', getPackages: { __typename?: 'PackagesResDto', packages: Array<{ __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } }> } };

export type GetPackageQueryVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type GetPackageQuery = { __typename?: 'Query', getPackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type GetSimsBySubscriberQueryVariables = Exact<{
  data: GetSimBySubscriberInputDto;
}>;


export type GetSimsBySubscriberQuery = { __typename?: 'Query', getSimsBySubscriber: { __typename?: 'SubscriberToSimsDto', subscriberId: string, sims: Array<{ __typename?: 'SubscriberSimsDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> } };

export type AddPackageMutationVariables = Exact<{
  data: AddPackageInputDto;
}>;


export type AddPackageMutation = { __typename?: 'Mutation', addPackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type RemovePackageForSimMutationVariables = Exact<{
  data: RemovePackageFormSimInputDto;
}>;


export type RemovePackageForSimMutation = { __typename?: 'Mutation', removePackageForSim: { __typename?: 'RemovePackageFromSimResDto', packageId?: string | null } };

export type DeletePackageMutationVariables = Exact<{
  packageId: Scalars['String']['input'];
}>;


export type DeletePackageMutation = { __typename?: 'Mutation', deletePackage: { __typename?: 'IdResponse', uuid: string } };

export type GetPackagesForSimQueryVariables = Exact<{
  data: GetPackagesForSimInputDto;
}>;


export type GetPackagesForSimQuery = { __typename?: 'Query', getPackagesForSim: { __typename?: 'GetSimPackagesDtoAPI', sim_id: string, packages: Array<{ __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean }> } };

export type AddPackagesToSimMutationVariables = Exact<{
  data: AddPackagesToSimInputDto;
}>;


export type AddPackagesToSimMutation = { __typename?: 'Mutation', addPackagesToSim: { __typename?: 'AddPackagesSimResDto', packages: Array<{ __typename?: 'AddPackagSimResDto', packageId?: string | null }> } };

export type DeleteSimMutationVariables = Exact<{
  data: DeleteSimInputDto;
}>;


export type DeleteSimMutation = { __typename?: 'Mutation', deleteSim: { __typename?: 'DeleteSimResDto', simId?: string | null } };

export type UpdatePacakgeMutationVariables = Exact<{
  packageId: Scalars['String']['input'];
  data: UpdatePackageInputDto;
}>;


export type UpdatePacakgeMutation = { __typename?: 'Mutation', updatePackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type PaymentFragment = { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string };

export type UpdatePaymentMutationVariables = Exact<{
  data: UpdatePaymentInputDto;
}>;


export type UpdatePaymentMutation = { __typename?: 'Mutation', updatePayment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } };

export type ProcessPaymentMutationVariables = Exact<{
  data: ProcessPaymentInputDto;
}>;


export type ProcessPaymentMutation = { __typename?: 'Mutation', processPayment: { __typename?: 'ProcessPaymentDto', payment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } } };

export type GetPaymentQueryVariables = Exact<{
  paymentId: Scalars['String']['input'];
}>;


export type GetPaymentQuery = { __typename?: 'Query', getPayment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string } };

export type GetPaymentsQueryVariables = Exact<{
  data: GetPaymentsInputDto;
}>;


export type GetPaymentsQuery = { __typename?: 'Query', getPayments: { __typename?: 'PaymentsDto', payments: Array<{ __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string }> } };

export type CustomerFragment = { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string };

export type SubscriptionFragment = { __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null };

export type FeeFragment = { __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } };

export type RawReportFragment = { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> };

export type GetReportsQueryVariables = Exact<{
  data: GetReportsInputDto;
}>;


export type GetReportsQuery = { __typename?: 'Query', getReports: { __typename?: 'GetReportsDto', reports: Array<{ __typename?: 'ReportDto', id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> } }> } };

export type GetReportQueryVariables = Exact<{
  id: Scalars['String']['input'];
}>;


export type GetReportQuery = { __typename?: 'Query', getReport: { __typename?: 'GetReportDto', report: { __typename?: 'ReportDto', id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> } } } };

export type GetSimPoolStatsQueryVariables = Exact<{
  data: GetSimsInput;
}>;


export type GetSimPoolStatsQuery = { __typename?: 'Query', getSimPoolStats: { __typename?: 'SimPoolStatsDto', total: number, available: number, consumed: number, failed: number, esim: number, physical: number } };

export type GetSimsFromPoolQueryVariables = Exact<{
  data: GetSimsInput;
}>;


export type GetSimsFromPoolQuery = { __typename?: 'Query', getSimsFromPool: { __typename?: 'SimsPoolResDto', sims: Array<{ __typename?: 'SimPoolResDto', id: string, qrCode: string, iccid: string, msisdn: string, isAllocated: boolean, isFailed: boolean, simType: string, smApAddress: string, activationCode: string, createdAt: string, deletedAt: string, updatedAt: string, isPhysical: boolean }> } };

export type UploadSimsMutationVariables = Exact<{
  data: UploadSimsInputDto;
}>;


export type UploadSimsMutation = { __typename?: 'Mutation', uploadSims: { __typename?: 'UploadSimsResDto', iccid: Array<string> } };

export type SimPackageFragment = { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean };

export type SimFragment = { __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null };

export type SimAllocationPackageFragment = { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null };

export type SimAllocationFragment = { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } };

export type AllocateSimMutationVariables = Exact<{
  data: AllocateSimInputDto;
}>;


export type AllocateSimMutation = { __typename?: 'Mutation', allocateSim: { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } } };

export type ToggleSimStatusMutationVariables = Exact<{
  data: ToggleSimStatusInputDto;
}>;


export type ToggleSimStatusMutation = { __typename?: 'Mutation', toggleSimStatus: { __typename?: 'SimStatusResDto', simId?: string | null } };

export type GetSimQueryVariables = Exact<{
  data: GetSimInputDto;
}>;


export type GetSimQuery = { __typename?: 'Query', getSim: { __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null } };

export type GetSimsQueryVariables = Exact<{
  data: ListSimsInput;
}>;


export type GetSimsQuery = { __typename?: 'Query', getSims: { __typename?: 'SimsResDto', sims: Array<{ __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null }> } };

export type SubscriberSimFragment = { __typename?: 'SubscriberDto', sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status?: string | null, isPhysical?: boolean | null, package?: { __typename?: 'SimPackageDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null };

export type SubscriberFragment = { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status?: string | null, isPhysical?: boolean | null, package?: { __typename?: 'SimPackageDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null };

export type AddSubscriberMutationVariables = Exact<{
  data: SubscriberInputDto;
}>;


export type AddSubscriberMutation = { __typename?: 'Mutation', addSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status?: string | null, isPhysical?: boolean | null, package?: { __typename?: 'SimPackageDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null } };

export type GetSubscriberQueryVariables = Exact<{
  subscriberId: Scalars['String']['input'];
}>;


export type GetSubscriberQuery = { __typename?: 'Query', getSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status?: string | null, isPhysical?: boolean | null, package?: { __typename?: 'SimPackageDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null } };

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


export type GetSubscribersByNetworkQuery = { __typename?: 'Query', getSubscribersByNetwork: { __typename?: 'SubscribersResDto', subscribers: Array<{ __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status?: string | null, isPhysical?: boolean | null, package?: { __typename?: 'SimPackageDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null }> } };

export type GetSubscriberMetricsByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscriberMetricsByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };

export type GetGeneratedPdfReportQueryVariables = Exact<{
  Id: Scalars['String']['input'];
}>;


export type GetGeneratedPdfReportQuery = { __typename?: 'Query', getGeneratedPdfReport: { __typename?: 'GetPdfReportUrlDto', contentType: string, filename: string, downloadUrl: string } };

export type UserFragment = { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string };

export type WhoamiQueryVariables = Exact<{ [key: string]: never; }>;


export type WhoamiQuery = { __typename?: 'Query', whoami: { __typename?: 'WhoamiDto', user: { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string }, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }> } };

export type GetUserQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetUserQuery = { __typename?: 'Query', getUser: { __typename?: 'UserResDto', name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string } };

export type UNetworkFragment = { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> };

export type GetNetworksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { __typename?: 'Query', getNetworks: { __typename?: 'NetworksResDto', networks: Array<{ __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> }> } };

export type GetNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetNetworkQuery = { __typename?: 'Query', getNetwork: { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type AddNetworkMutationVariables = Exact<{
  data: AddNetworkInputDto;
}>;


export type AddNetworkMutation = { __typename?: 'Mutation', addNetwork: { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type SetDefaultNetworkMutationVariables = Exact<{
  data: SetDefaultNetworkInputDto;
}>;


export type SetDefaultNetworkMutation = { __typename?: 'Mutation', setDefaultNetwork: { __typename?: 'CBooleanResponse', success: boolean } };

export type USiteFragment = { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string };

export type GetSiteQueryVariables = Exact<{
  siteId: Scalars['String']['input'];
}>;


export type GetSiteQuery = { __typename?: 'Query', getSite: { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type AddSiteMutationVariables = Exact<{
  data: AddSiteInputDto;
}>;


export type AddSiteMutation = { __typename?: 'Mutation', addSite: { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type GetSitesQueryVariables = Exact<{
  data: SitesInputDto;
}>;


export type GetSitesQuery = { __typename?: 'Query', getSites: { __typename?: 'SitesResDto', sites: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string }> } };

export type UpdateSiteMutationVariables = Exact<{
  siteId: Scalars['String']['input'];
  data: UpdateSiteInputDto;
}>;


export type UpdateSiteMutation = { __typename?: 'Mutation', updateSite: { __typename?: 'SiteDto', name: string } };

export type UComponentFragment = { __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string };

export type GetComponentByIdQueryVariables = Exact<{
  componentId: Scalars['String']['input'];
}>;


export type GetComponentByIdQuery = { __typename?: 'Query', getComponentById: { __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string } };

export type GetComponentsByUserIdQueryVariables = Exact<{
  data: ComponentTypeInputDto;
}>;


export type GetComponentsByUserIdQuery = { __typename?: 'Query', getComponentsByUserId: { __typename?: 'ComponentsResDto', components: Array<{ __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string }> } };

export type InvitationFragment = { __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Invitation_Status };

export type CreateInvitationMutationVariables = Exact<{
  data: CreateInvitationInputDto;
}>;


export type CreateInvitationMutation = { __typename?: 'Mutation', createInvitation: { __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Invitation_Status } };

export type GetInvitationsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetInvitationsQuery = { __typename?: 'Query', getInvitations: { __typename?: 'InvitationsResDto', invitations: Array<{ __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Invitation_Status }> } };

export type DeleteInvitationMutationVariables = Exact<{
  deleteInvitationId: Scalars['String']['input'];
}>;


export type DeleteInvitationMutation = { __typename?: 'Mutation', deleteInvitation: { __typename?: 'DeleteInvitationResDto', id: string } };

export type UpdateInvitationMutationVariables = Exact<{
  data: UpateInvitationInputDto;
}>;


export type UpdateInvitationMutation = { __typename?: 'Mutation', updateInvitation: { __typename?: 'UpdateInvitationResDto', id: string } };

export type GetInvitationsByEmailQueryVariables = Exact<{
  email: Scalars['String']['input'];
}>;


export type GetInvitationsByEmailQuery = { __typename?: 'Query', getInvitationsByEmail: { __typename?: 'InvitationsResDto', invitations: Array<{ __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Invitation_Status }> } };

export type GetCountriesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCountriesQuery = { __typename?: 'Query', getCountries: { __typename?: 'CountriesRes', countries: Array<{ __typename?: 'CountryDto', name: string, code: string }> } };

export type GetCurrencySymbolQueryVariables = Exact<{
  code: Scalars['String']['input'];
}>;


export type GetCurrencySymbolQuery = { __typename?: 'Query', getCurrencySymbol: { __typename?: 'CurrencyRes', code: string, symbol: string, image: string } };

export type GetTimezonesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetTimezonesQuery = { __typename?: 'Query', getTimezones: { __typename?: 'TimezoneRes', timezones: Array<{ __typename?: 'TimezoneDto', value: string, abbr: string, offset: number, isdst: boolean, text: string, utc: Array<string> }> } };

export type UpdateNotificationMutationVariables = Exact<{
  isRead: Scalars['Boolean']['input'];
  updateNotificationId: Scalars['String']['input'];
}>;


export type UpdateNotificationMutation = { __typename?: 'Mutation', updateNotification: { __typename?: 'UpdateNotificationResDto', id: string } };

export type GetDataUsagesQueryVariables = Exact<{
  data: SimUsagesInputDto;
}>;


export type GetDataUsagesQuery = { __typename?: 'Query', getDataUsages: { __typename?: 'SimDataUsages', usages: Array<{ __typename?: 'SimDataUsage', usage: string, simId: string }> } };

export const NodeFragmentDoc = gql`
    fragment node on Node {
  id
  name
  latitude
  longitude
  type
  attached {
    id
    name
    latitude
    longitude
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
export const MemberFragmentDoc = gql`
    fragment member on MemberDto {
  role
  userId
  id: memberId
  isDeactivated
  memberSince
}
    `;
export const OrgFragmentDoc = gql`
    fragment Org on OrgDto {
  id
  name
  owner
  country
  currency
  createdAt
  certificate
  isDeactivated
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
export const SubscriberSimsFragmentDoc = gql`
    fragment SubscriberSims on SubscriberToSimsDto {
  subscriberId
  sims {
    id
    subscriberId
    networkId
    iccid
    msisdn
    imsi
    type
    status
    allocatedAt
    isPhysical
  }
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
export const PaymentFragmentDoc = gql`
    fragment payment on PaymentDto {
  id
  itemId
  itemType
  amount
  currency
  paymentMethod
  depositedAmount
  paidAt
  payerName
  payerEmail
  payerPhone
  correspondent
  country
  description
  status
  failureReason
  extra
  createdAt
}
    `;
export const CustomerFragmentDoc = gql`
    fragment customer on CustomerDto {
  externalId
  name
  email
  addressLine1
  legalName
  legalNumber
  phone
  currency
  timezone
  vatRate
  createdAt
}
    `;
export const SubscriptionFragmentDoc = gql`
    fragment subscription on SubscriptionDto {
  externalCustomerId
  externalId
  planCode
  name
  status
  createdAt
  startedAt
  canceledAt
  terminatedAt
}
    `;
export const FeeFragmentDoc = gql`
    fragment fee on FeeDto {
  taxesAmountCents
  taxesPreciseAmount
  totalAmountCents
  totalAmountCurrency
  eventsCount
  units
  item {
    type
    code
    name
  }
}
    `;
export const RawReportFragmentDoc = gql`
    fragment rawReport on RawReportDto {
  issuingDate
  paymentDueDate
  paymentOverdue
  invoiceType
  status
  paymentStatus
  feesAmountCents
  taxesAmountCents
  subTotalExcludingTaxesAmountCents
  subTotalIncludingTaxesAmountCents
  vatAmountCents
  vatAmountCurrency
  totalAmountCents
  currency
  fileUrl
  customer {
    ...customer
  }
  subscriptions {
    ...subscription
  }
  fees {
    ...fee
  }
}
    ${CustomerFragmentDoc}
${SubscriptionFragmentDoc}
${FeeFragmentDoc}`;
export const SimPackageFragmentDoc = gql`
    fragment SimPackage on SimPackage {
  id
  packageId
  startDate
  endDate
  defaultDuration
  isActive
  asExpired
}
    `;
export const SimFragmentDoc = gql`
    fragment Sim on SimDto {
  id
  subscriberId
  networkId
  iccid
  msisdn
  imsi
  type
  status
  isPhysical
  trafficPolicy
  firstActivatedOn
  lastActivatedOn
  activationsCount
  deactivationsCount
  allocatedAt
  syncStatus
  package {
    ...SimPackage
  }
}
    ${SimPackageFragmentDoc}`;
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
    iccid
    msisdn
    imsi
    type
    status
    allocatedAt
    sync_status
    isPhysical
    package {
      id
      package_id
      start_date
      end_date
      is_active
      created_at
      updated_at
    }
  }
}
    `;
export const SubscriberFragmentDoc = gql`
    fragment Subscriber on SubscriberDto {
  uuid
  address
  dob
  email
  name
  gender
  idSerial
  networkId
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
  isDefault
  budget
  overdraft
  trafficPolicy
  isDeactivated
  paymentLinks
  createdAt
  countries
  networks
}
    `;
export const USiteFragmentDoc = gql`
    fragment USite on SiteDto {
  id
  name
  networkId
  backhaulId
  powerId
  accessId
  spectrumId
  switchId
  isDeactivated
  latitude
  longitude
  installDate
  createdAt
  location
}
    `;
export const UComponentFragmentDoc = gql`
    fragment UComponent on ComponentDto {
  id
  inventoryId
  type
  userId
  description
  category
  datasheetUrl
  imageUrl
  partNumber
  manufacturer
  managed
  warranty
  specification
}
    `;
export const InvitationFragmentDoc = gql`
    fragment Invitation on InvitationDto {
  email
  expireAt
  id
  name
  role
  link
  userId
  status
}
    `;
export const GetNodeDocument = gql`
    query GetNode($data: NodeInput!) {
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
export function useGetNodeQuery(baseOptions: Apollo.QueryHookOptions<GetNodeQuery, GetNodeQueryVariables> & ({ variables: GetNodeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
      }
export function useGetNodeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export type GetNodeQueryHookResult = ReturnType<typeof useGetNodeQuery>;
export type GetNodeLazyQueryHookResult = ReturnType<typeof useGetNodeLazyQuery>;
export type GetNodeSuspenseQueryHookResult = ReturnType<typeof useGetNodeSuspenseQuery>;
export type GetNodeQueryResult = Apollo.QueryResult<GetNodeQuery, GetNodeQueryVariables>;
export const GetNodesDocument = gql`
    query GetNodes($data: NodesFilterInput!) {
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
export function useGetNodesQuery(baseOptions: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables> & ({ variables: GetNodesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
      }
export function useGetNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<typeof useGetNodesLazyQuery>;
export type GetNodesSuspenseQueryHookResult = ReturnType<typeof useGetNodesSuspenseQuery>;
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
export function useGetNodesForSiteQuery(baseOptions: Apollo.QueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables> & ({ variables: GetNodesForSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
      }
export function useGetNodesForSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
export type GetNodesForSiteQueryHookResult = ReturnType<typeof useGetNodesForSiteQuery>;
export type GetNodesForSiteLazyQueryHookResult = ReturnType<typeof useGetNodesForSiteLazyQuery>;
export type GetNodesForSiteSuspenseQueryHookResult = ReturnType<typeof useGetNodesForSiteSuspenseQuery>;
export type GetNodesForSiteQueryResult = Apollo.QueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
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
export function useGetNodeAppsQuery(baseOptions: Apollo.QueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables> & ({ variables: GetNodeAppsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
      }
export function useGetNodeAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
export type GetNodeAppsQueryHookResult = ReturnType<typeof useGetNodeAppsQuery>;
export type GetNodeAppsLazyQueryHookResult = ReturnType<typeof useGetNodeAppsLazyQuery>;
export type GetNodeAppsSuspenseQueryHookResult = ReturnType<typeof useGetNodeAppsSuspenseQuery>;
export type GetNodeAppsQueryResult = Apollo.QueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export const GetNodeStateDocument = gql`
    query GetNodeState($getNodeStateId: String!) {
  getNodeState(id: $getNodeStateId) {
    id
    nodeId
    previousStateId
    previousState
    currentState
    createdAt
  }
}
    `;

/**
 * __useGetNodeStateQuery__
 *
 * To run a query within a React component, call `useGetNodeStateQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeStateQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeStateQuery({
 *   variables: {
 *      getNodeStateId: // value for 'getNodeStateId'
 *   },
 * });
 */
export function useGetNodeStateQuery(baseOptions: Apollo.QueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables> & ({ variables: GetNodeStateQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
      }
export function useGetNodeStateLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
export type GetNodeStateQueryHookResult = ReturnType<typeof useGetNodeStateQuery>;
export type GetNodeStateLazyQueryHookResult = ReturnType<typeof useGetNodeStateLazyQuery>;
export type GetNodeStateSuspenseQueryHookResult = ReturnType<typeof useGetNodeStateSuspenseQuery>;
export type GetNodeStateQueryResult = Apollo.QueryResult<GetNodeStateQuery, GetNodeStateQueryVariables>;
export const RestartNodeDocument = gql`
    mutation RestartNode($data: RestartNodeInputDto!) {
  restartNode(data: $data) {
    success
  }
}
    `;
export type RestartNodeMutationFn = Apollo.MutationFunction<RestartNodeMutation, RestartNodeMutationVariables>;

/**
 * __useRestartNodeMutation__
 *
 * To run a mutation, you first call `useRestartNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestartNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restartNodeMutation, { data, loading, error }] = useRestartNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRestartNodeMutation(baseOptions?: Apollo.MutationHookOptions<RestartNodeMutation, RestartNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RestartNodeMutation, RestartNodeMutationVariables>(RestartNodeDocument, options);
      }
export type RestartNodeMutationHookResult = ReturnType<typeof useRestartNodeMutation>;
export type RestartNodeMutationResult = Apollo.MutationResult<RestartNodeMutation>;
export type RestartNodeMutationOptions = Apollo.BaseMutationOptions<RestartNodeMutation, RestartNodeMutationVariables>;
export const ToggleInternetSwitchDocument = gql`
    mutation ToggleInternetSwitch($data: ToggleInternetSwitchInputDto!) {
  toggleInternetSwitch(data: $data) {
    success
  }
}
    `;
export type ToggleInternetSwitchMutationFn = Apollo.MutationFunction<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;

/**
 * __useToggleInternetSwitchMutation__
 *
 * To run a mutation, you first call `useToggleInternetSwitchMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleInternetSwitchMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleInternetSwitchMutation, { data, loading, error }] = useToggleInternetSwitchMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleInternetSwitchMutation(baseOptions?: Apollo.MutationHookOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>(ToggleInternetSwitchDocument, options);
      }
export type ToggleInternetSwitchMutationHookResult = ReturnType<typeof useToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationResult = Apollo.MutationResult<ToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationOptions = Apollo.BaseMutationOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;
export const ToggleRfStatusDocument = gql`
    mutation ToggleRFStatus($data: ToggleRFStatusInputDto!) {
  toggleRFStatus(data: $data) {
    success
  }
}
    `;
export type ToggleRfStatusMutationFn = Apollo.MutationFunction<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;

/**
 * __useToggleRfStatusMutation__
 *
 * To run a mutation, you first call `useToggleRfStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleRfStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleRfStatusMutation, { data, loading, error }] = useToggleRfStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleRfStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>(ToggleRfStatusDocument, options);
      }
export type ToggleRfStatusMutationHookResult = ReturnType<typeof useToggleRfStatusMutation>;
export type ToggleRfStatusMutationResult = Apollo.MutationResult<ToggleRfStatusMutation>;
export type ToggleRfStatusMutationOptions = Apollo.BaseMutationOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;
export const GetMembersDocument = gql`
    query GetMembers {
  getMembers {
    members {
      ...member
      name
      email
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
export function useGetMembersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
        }
export type GetMembersQueryHookResult = ReturnType<typeof useGetMembersQuery>;
export type GetMembersLazyQueryHookResult = ReturnType<typeof useGetMembersLazyQuery>;
export type GetMembersSuspenseQueryHookResult = ReturnType<typeof useGetMembersSuspenseQuery>;
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
export function useGetMemberQuery(baseOptions: Apollo.QueryHookOptions<GetMemberQuery, GetMemberQueryVariables> & ({ variables: GetMemberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
      }
export function useGetMemberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
        }
export function useGetMemberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
        }
export type GetMemberQueryHookResult = ReturnType<typeof useGetMemberQuery>;
export type GetMemberLazyQueryHookResult = ReturnType<typeof useGetMemberLazyQuery>;
export type GetMemberSuspenseQueryHookResult = ReturnType<typeof useGetMemberSuspenseQuery>;
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
export const GetMemberByUserIdDocument = gql`
    query GetMemberByUserId($userId: String!) {
  getMemberByUserId(userId: $userId) {
    userId
    name
    email
    memberId
    isDeactivated
    role
    memberSince
  }
}
    `;

/**
 * __useGetMemberByUserIdQuery__
 *
 * To run a query within a React component, call `useGetMemberByUserIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMemberByUserIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemberByUserIdQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetMemberByUserIdQuery(baseOptions: Apollo.QueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables> & ({ variables: GetMemberByUserIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
      }
export function useGetMemberByUserIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
        }
export function useGetMemberByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
        }
export type GetMemberByUserIdQueryHookResult = ReturnType<typeof useGetMemberByUserIdQuery>;
export type GetMemberByUserIdLazyQueryHookResult = ReturnType<typeof useGetMemberByUserIdLazyQuery>;
export type GetMemberByUserIdSuspenseQueryHookResult = ReturnType<typeof useGetMemberByUserIdSuspenseQuery>;
export type GetMemberByUserIdQueryResult = Apollo.QueryResult<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>;
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
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
export type GetOrgsQueryHookResult = ReturnType<typeof useGetOrgsQuery>;
export type GetOrgsLazyQueryHookResult = ReturnType<typeof useGetOrgsLazyQuery>;
export type GetOrgsSuspenseQueryHookResult = ReturnType<typeof useGetOrgsSuspenseQuery>;
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
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
export type GetOrgQueryHookResult = ReturnType<typeof useGetOrgQuery>;
export type GetOrgLazyQueryHookResult = ReturnType<typeof useGetOrgLazyQuery>;
export type GetOrgSuspenseQueryHookResult = ReturnType<typeof useGetOrgSuspenseQuery>;
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
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
export type GetPackagesQueryHookResult = ReturnType<typeof useGetPackagesQuery>;
export type GetPackagesLazyQueryHookResult = ReturnType<typeof useGetPackagesLazyQuery>;
export type GetPackagesSuspenseQueryHookResult = ReturnType<typeof useGetPackagesSuspenseQuery>;
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
export function useGetPackageQuery(baseOptions: Apollo.QueryHookOptions<GetPackageQuery, GetPackageQueryVariables> & ({ variables: GetPackageQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
      }
export function useGetPackageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
export type GetPackageQueryHookResult = ReturnType<typeof useGetPackageQuery>;
export type GetPackageLazyQueryHookResult = ReturnType<typeof useGetPackageLazyQuery>;
export type GetPackageSuspenseQueryHookResult = ReturnType<typeof useGetPackageSuspenseQuery>;
export type GetPackageQueryResult = Apollo.QueryResult<GetPackageQuery, GetPackageQueryVariables>;
export const GetSimsBySubscriberDocument = gql`
    query getSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    ...SubscriberSims
  }
}
    ${SubscriberSimsFragmentDoc}`;

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
export function useGetSimsBySubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables> & ({ variables: GetSimsBySubscriberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
      }
export function useGetSimsBySubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
export type GetSimsBySubscriberQueryHookResult = ReturnType<typeof useGetSimsBySubscriberQuery>;
export type GetSimsBySubscriberLazyQueryHookResult = ReturnType<typeof useGetSimsBySubscriberLazyQuery>;
export type GetSimsBySubscriberSuspenseQueryHookResult = ReturnType<typeof useGetSimsBySubscriberSuspenseQuery>;
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
export function useGetPackagesForSimQuery(baseOptions: Apollo.QueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables> & ({ variables: GetPackagesForSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
      }
export function useGetPackagesForSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
export type GetPackagesForSimQueryHookResult = ReturnType<typeof useGetPackagesForSimQuery>;
export type GetPackagesForSimLazyQueryHookResult = ReturnType<typeof useGetPackagesForSimLazyQuery>;
export type GetPackagesForSimSuspenseQueryHookResult = ReturnType<typeof useGetPackagesForSimSuspenseQuery>;
export type GetPackagesForSimQueryResult = Apollo.QueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export const AddPackagesToSimDocument = gql`
    mutation addPackagesToSim($data: AddPackagesToSimInputDto!) {
  addPackagesToSim(data: $data) {
    packages {
      packageId
    }
  }
}
    `;
export type AddPackagesToSimMutationFn = Apollo.MutationFunction<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;

/**
 * __useAddPackagesToSimMutation__
 *
 * To run a mutation, you first call `useAddPackagesToSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackagesToSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackagesToSimMutation, { data, loading, error }] = useAddPackagesToSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackagesToSimMutation(baseOptions?: Apollo.MutationHookOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>(AddPackagesToSimDocument, options);
      }
export type AddPackagesToSimMutationHookResult = ReturnType<typeof useAddPackagesToSimMutation>;
export type AddPackagesToSimMutationResult = Apollo.MutationResult<AddPackagesToSimMutation>;
export type AddPackagesToSimMutationOptions = Apollo.BaseMutationOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;
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
export const UpdatePaymentDocument = gql`
    mutation UpdatePayment($data: UpdatePaymentInputDto!) {
  updatePayment(data: $data) {
    id
    itemId
    itemType
    amount
    currency
    paymentMethod
    depositedAmount
    paidAt
    payerName
    payerEmail
    payerPhone
    correspondent
    country
    description
    status
    failureReason
    createdAt
  }
}
    `;
export type UpdatePaymentMutationFn = Apollo.MutationFunction<UpdatePaymentMutation, UpdatePaymentMutationVariables>;

/**
 * __useUpdatePaymentMutation__
 *
 * To run a mutation, you first call `useUpdatePaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePaymentMutation, { data, loading, error }] = useUpdatePaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePaymentMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePaymentMutation, UpdatePaymentMutationVariables>(UpdatePaymentDocument, options);
      }
export type UpdatePaymentMutationHookResult = ReturnType<typeof useUpdatePaymentMutation>;
export type UpdatePaymentMutationResult = Apollo.MutationResult<UpdatePaymentMutation>;
export type UpdatePaymentMutationOptions = Apollo.BaseMutationOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>;
export const ProcessPaymentDocument = gql`
    mutation ProcessPayment($data: ProcessPaymentInputDto!) {
  processPayment(data: $data) {
    payment {
      id
      itemId
      itemType
      amount
      currency
      paymentMethod
      depositedAmount
      paidAt
      payerName
      payerEmail
      payerPhone
      correspondent
      country
      description
      status
      failureReason
      createdAt
    }
  }
}
    `;
export type ProcessPaymentMutationFn = Apollo.MutationFunction<ProcessPaymentMutation, ProcessPaymentMutationVariables>;

/**
 * __useProcessPaymentMutation__
 *
 * To run a mutation, you first call `useProcessPaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useProcessPaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [processPaymentMutation, { data, loading, error }] = useProcessPaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useProcessPaymentMutation(baseOptions?: Apollo.MutationHookOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ProcessPaymentMutation, ProcessPaymentMutationVariables>(ProcessPaymentDocument, options);
      }
export type ProcessPaymentMutationHookResult = ReturnType<typeof useProcessPaymentMutation>;
export type ProcessPaymentMutationResult = Apollo.MutationResult<ProcessPaymentMutation>;
export type ProcessPaymentMutationOptions = Apollo.BaseMutationOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>;
export const GetPaymentDocument = gql`
    query GetPayment($paymentId: String!) {
  getPayment(paymentId: $paymentId) {
    ...payment
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentQuery__
 *
 * To run a query within a React component, call `useGetPaymentQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentQuery({
 *   variables: {
 *      paymentId: // value for 'paymentId'
 *   },
 * });
 */
export function useGetPaymentQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables> & ({ variables: GetPaymentQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
      }
export function useGetPaymentLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
export type GetPaymentQueryHookResult = ReturnType<typeof useGetPaymentQuery>;
export type GetPaymentLazyQueryHookResult = ReturnType<typeof useGetPaymentLazyQuery>;
export type GetPaymentSuspenseQueryHookResult = ReturnType<typeof useGetPaymentSuspenseQuery>;
export type GetPaymentQueryResult = Apollo.QueryResult<GetPaymentQuery, GetPaymentQueryVariables>;
export const GetPaymentsDocument = gql`
    query GetPayments($data: GetPaymentsInputDto!) {
  getPayments(data: $data) {
    payments {
      ...payment
    }
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentsQuery__
 *
 * To run a query within a React component, call `useGetPaymentsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPaymentsQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables> & ({ variables: GetPaymentsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
      }
export function useGetPaymentsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
export type GetPaymentsQueryHookResult = ReturnType<typeof useGetPaymentsQuery>;
export type GetPaymentsLazyQueryHookResult = ReturnType<typeof useGetPaymentsLazyQuery>;
export type GetPaymentsSuspenseQueryHookResult = ReturnType<typeof useGetPaymentsSuspenseQuery>;
export type GetPaymentsQueryResult = Apollo.QueryResult<GetPaymentsQuery, GetPaymentsQueryVariables>;
export const GetReportsDocument = gql`
    query GetReports($data: GetReportsInputDto!) {
  getReports(data: $data) {
    reports {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportsQuery__
 *
 * To run a query within a React component, call `useGetReportsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetReportsQuery(baseOptions: Apollo.QueryHookOptions<GetReportsQuery, GetReportsQueryVariables> & ({ variables: GetReportsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
      }
export function useGetReportsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
export type GetReportsQueryHookResult = ReturnType<typeof useGetReportsQuery>;
export type GetReportsLazyQueryHookResult = ReturnType<typeof useGetReportsLazyQuery>;
export type GetReportsSuspenseQueryHookResult = ReturnType<typeof useGetReportsSuspenseQuery>;
export type GetReportsQueryResult = Apollo.QueryResult<GetReportsQuery, GetReportsQueryVariables>;
export const GetReportDocument = gql`
    query GetReport($id: String!) {
  getReport(id: $id) {
    report {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportQuery__
 *
 * To run a query within a React component, call `useGetReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetReportQuery(baseOptions: Apollo.QueryHookOptions<GetReportQuery, GetReportQueryVariables> & ({ variables: GetReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
      }
export function useGetReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
export type GetReportQueryHookResult = ReturnType<typeof useGetReportQuery>;
export type GetReportLazyQueryHookResult = ReturnType<typeof useGetReportLazyQuery>;
export type GetReportSuspenseQueryHookResult = ReturnType<typeof useGetReportSuspenseQuery>;
export type GetReportQueryResult = Apollo.QueryResult<GetReportQuery, GetReportQueryVariables>;
export const GetSimPoolStatsDocument = gql`
    query GetSimPoolStats($data: GetSimsInput!) {
  getSimPoolStats(data: $data) {
    total
    available
    consumed
    failed
    esim
    physical
  }
}
    `;

/**
 * __useGetSimPoolStatsQuery__
 *
 * To run a query within a React component, call `useGetSimPoolStatsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimPoolStatsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimPoolStatsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimPoolStatsQuery(baseOptions: Apollo.QueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables> & ({ variables: GetSimPoolStatsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
      }
export function useGetSimPoolStatsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
export type GetSimPoolStatsQueryHookResult = ReturnType<typeof useGetSimPoolStatsQuery>;
export type GetSimPoolStatsLazyQueryHookResult = ReturnType<typeof useGetSimPoolStatsLazyQuery>;
export type GetSimPoolStatsSuspenseQueryHookResult = ReturnType<typeof useGetSimPoolStatsSuspenseQuery>;
export type GetSimPoolStatsQueryResult = Apollo.QueryResult<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>;
export const GetSimsFromPoolDocument = gql`
    query GetSimsFromPool($data: GetSimsInput!) {
  getSimsFromPool(data: $data) {
    sims {
      id
      qrCode
      iccid
      msisdn
      isAllocated
      isFailed
      simType
      smApAddress
      activationCode
      createdAt
      deletedAt
      updatedAt
      isPhysical
    }
  }
}
    `;

/**
 * __useGetSimsFromPoolQuery__
 *
 * To run a query within a React component, call `useGetSimsFromPoolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsFromPoolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsFromPoolQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsFromPoolQuery(baseOptions: Apollo.QueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables> & ({ variables: GetSimsFromPoolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
      }
export function useGetSimsFromPoolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
export type GetSimsFromPoolQueryHookResult = ReturnType<typeof useGetSimsFromPoolQuery>;
export type GetSimsFromPoolLazyQueryHookResult = ReturnType<typeof useGetSimsFromPoolLazyQuery>;
export type GetSimsFromPoolSuspenseQueryHookResult = ReturnType<typeof useGetSimsFromPoolSuspenseQuery>;
export type GetSimsFromPoolQueryResult = Apollo.QueryResult<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>;
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
    ...Sim
  }
}
    ${SimFragmentDoc}`;

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
export function useGetSimQuery(baseOptions: Apollo.QueryHookOptions<GetSimQuery, GetSimQueryVariables> & ({ variables: GetSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
      }
export function useGetSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
export type GetSimQueryHookResult = ReturnType<typeof useGetSimQuery>;
export type GetSimLazyQueryHookResult = ReturnType<typeof useGetSimLazyQuery>;
export type GetSimSuspenseQueryHookResult = ReturnType<typeof useGetSimSuspenseQuery>;
export type GetSimQueryResult = Apollo.QueryResult<GetSimQuery, GetSimQueryVariables>;
export const GetSimsDocument = gql`
    query GetSims($data: ListSimsInput!) {
  getSims(data: $data) {
    sims {
      ...Sim
    }
  }
}
    ${SimFragmentDoc}`;

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
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsQuery(baseOptions: Apollo.QueryHookOptions<GetSimsQuery, GetSimsQueryVariables> & ({ variables: GetSimsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
      }
export function useGetSimsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
export type GetSimsQueryHookResult = ReturnType<typeof useGetSimsQuery>;
export type GetSimsLazyQueryHookResult = ReturnType<typeof useGetSimsLazyQuery>;
export type GetSimsSuspenseQueryHookResult = ReturnType<typeof useGetSimsSuspenseQuery>;
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
export function useGetSubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables> & ({ variables: GetSubscriberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
      }
export function useGetSubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
        }
export function useGetSubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
        }
export type GetSubscriberQueryHookResult = ReturnType<typeof useGetSubscriberQuery>;
export type GetSubscriberLazyQueryHookResult = ReturnType<typeof useGetSubscriberLazyQuery>;
export type GetSubscriberSuspenseQueryHookResult = ReturnType<typeof useGetSubscriberSuspenseQuery>;
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
export function useGetSubscribersByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables> & ({ variables: GetSubscribersByNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
      }
export function useGetSubscribersByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
        }
export function useGetSubscribersByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
        }
export type GetSubscribersByNetworkQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkQuery>;
export type GetSubscribersByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkLazyQuery>;
export type GetSubscribersByNetworkSuspenseQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkSuspenseQuery>;
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
export function useGetSubscriberMetricsByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables> & ({ variables: GetSubscriberMetricsByNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
      }
export function useGetSubscriberMetricsByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
        }
export function useGetSubscriberMetricsByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
        }
export type GetSubscriberMetricsByNetworkQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkQuery>;
export type GetSubscriberMetricsByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkLazyQuery>;
export type GetSubscriberMetricsByNetworkSuspenseQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkSuspenseQuery>;
export type GetSubscriberMetricsByNetworkQueryResult = Apollo.QueryResult<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>;
export const GetGeneratedPdfReportDocument = gql`
    query getGeneratedPdfReport($Id: String!) {
  getGeneratedPdfReport(id: $Id) {
    contentType
    filename
    downloadUrl
  }
}
    `;

/**
 * __useGetGeneratedPdfReportQuery__
 *
 * To run a query within a React component, call `useGetGeneratedPdfReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetGeneratedPdfReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetGeneratedPdfReportQuery({
 *   variables: {
 *      Id: // value for 'Id'
 *   },
 * });
 */
export function useGetGeneratedPdfReportQuery(baseOptions: Apollo.QueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables> & ({ variables: GetGeneratedPdfReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
      }
export function useGetGeneratedPdfReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
export type GetGeneratedPdfReportQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportQuery>;
export type GetGeneratedPdfReportLazyQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportLazyQuery>;
export type GetGeneratedPdfReportSuspenseQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportSuspenseQuery>;
export type GetGeneratedPdfReportQueryResult = Apollo.QueryResult<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>;
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
export function useWhoamiSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
        }
export type WhoamiQueryHookResult = ReturnType<typeof useWhoamiQuery>;
export type WhoamiLazyQueryHookResult = ReturnType<typeof useWhoamiLazyQuery>;
export type WhoamiSuspenseQueryHookResult = ReturnType<typeof useWhoamiSuspenseQuery>;
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
export function useGetUserQuery(baseOptions: Apollo.QueryHookOptions<GetUserQuery, GetUserQueryVariables> & ({ variables: GetUserQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
      }
export function useGetUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export function useGetUserSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export type GetUserQueryHookResult = ReturnType<typeof useGetUserQuery>;
export type GetUserLazyQueryHookResult = ReturnType<typeof useGetUserLazyQuery>;
export type GetUserSuspenseQueryHookResult = ReturnType<typeof useGetUserSuspenseQuery>;
export type GetUserQueryResult = Apollo.QueryResult<GetUserQuery, GetUserQueryVariables>;
export const GetNetworksDocument = gql`
    query getNetworks {
  getNetworks {
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
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
export type GetNetworksQueryHookResult = ReturnType<typeof useGetNetworksQuery>;
export type GetNetworksLazyQueryHookResult = ReturnType<typeof useGetNetworksLazyQuery>;
export type GetNetworksSuspenseQueryHookResult = ReturnType<typeof useGetNetworksSuspenseQuery>;
export type GetNetworksQueryResult = Apollo.QueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export const GetNetworkDocument = gql`
    query getNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworkQuery__
 *
 * To run a query within a React component, call `useGetNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables> & ({ variables: GetNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
      }
export function useGetNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
export type GetNetworkQueryHookResult = ReturnType<typeof useGetNetworkQuery>;
export type GetNetworkLazyQueryHookResult = ReturnType<typeof useGetNetworkLazyQuery>;
export type GetNetworkSuspenseQueryHookResult = ReturnType<typeof useGetNetworkSuspenseQuery>;
export type GetNetworkQueryResult = Apollo.QueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
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
export const SetDefaultNetworkDocument = gql`
    mutation SetDefaultNetwork($data: SetDefaultNetworkInputDto!) {
  setDefaultNetwork(data: $data) {
    success
  }
}
    `;
export type SetDefaultNetworkMutationFn = Apollo.MutationFunction<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;

/**
 * __useSetDefaultNetworkMutation__
 *
 * To run a mutation, you first call `useSetDefaultNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetDefaultNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setDefaultNetworkMutation, { data, loading, error }] = useSetDefaultNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSetDefaultNetworkMutation(baseOptions?: Apollo.MutationHookOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>(SetDefaultNetworkDocument, options);
      }
export type SetDefaultNetworkMutationHookResult = ReturnType<typeof useSetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationResult = Apollo.MutationResult<SetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationOptions = Apollo.BaseMutationOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;
export const GetSiteDocument = gql`
    query getSite($siteId: String!) {
  getSite(siteId: $siteId) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;

/**
 * __useGetSiteQuery__
 *
 * To run a query within a React component, call `useGetSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetSiteQuery(baseOptions: Apollo.QueryHookOptions<GetSiteQuery, GetSiteQueryVariables> & ({ variables: GetSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
      }
export function useGetSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
export type GetSiteQueryHookResult = ReturnType<typeof useGetSiteQuery>;
export type GetSiteLazyQueryHookResult = ReturnType<typeof useGetSiteLazyQuery>;
export type GetSiteSuspenseQueryHookResult = ReturnType<typeof useGetSiteSuspenseQuery>;
export type GetSiteQueryResult = Apollo.QueryResult<GetSiteQuery, GetSiteQueryVariables>;
export const AddSiteDocument = gql`
    mutation addSite($data: AddSiteInputDto!) {
  addSite(data: $data) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;
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
export const GetSitesDocument = gql`
    query GetSites($data: SitesInputDto!) {
  getSites(data: $data) {
    sites {
      ...USite
    }
  }
}
    ${USiteFragmentDoc}`;

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
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSitesQuery(baseOptions: Apollo.QueryHookOptions<GetSitesQuery, GetSitesQueryVariables> & ({ variables: GetSitesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
      }
export function useGetSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
export type GetSitesQueryHookResult = ReturnType<typeof useGetSitesQuery>;
export type GetSitesLazyQueryHookResult = ReturnType<typeof useGetSitesLazyQuery>;
export type GetSitesSuspenseQueryHookResult = ReturnType<typeof useGetSitesSuspenseQuery>;
export type GetSitesQueryResult = Apollo.QueryResult<GetSitesQuery, GetSitesQueryVariables>;
export const UpdateSiteDocument = gql`
    mutation updateSite($siteId: String!, $data: UpdateSiteInputDto!) {
  updateSite(siteId: $siteId, data: $data) {
    name
  }
}
    `;
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
export const GetComponentByIdDocument = gql`
    query getComponentById($componentId: String!) {
  getComponentById(componentId: $componentId) {
    ...UComponent
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentByIdQuery__
 *
 * To run a query within a React component, call `useGetComponentByIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentByIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentByIdQuery({
 *   variables: {
 *      componentId: // value for 'componentId'
 *   },
 * });
 */
export function useGetComponentByIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables> & ({ variables: GetComponentByIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
      }
export function useGetComponentByIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
export type GetComponentByIdQueryHookResult = ReturnType<typeof useGetComponentByIdQuery>;
export type GetComponentByIdLazyQueryHookResult = ReturnType<typeof useGetComponentByIdLazyQuery>;
export type GetComponentByIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentByIdSuspenseQuery>;
export type GetComponentByIdQueryResult = Apollo.QueryResult<GetComponentByIdQuery, GetComponentByIdQueryVariables>;
export const GetComponentsByUserIdDocument = gql`
    query GetComponentsByUserId($data: ComponentTypeInputDto!) {
  getComponentsByUserId(data: $data) {
    components {
      ...UComponent
    }
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentsByUserIdQuery__
 *
 * To run a query within a React component, call `useGetComponentsByUserIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentsByUserIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentsByUserIdQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetComponentsByUserIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables> & ({ variables: GetComponentsByUserIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
      }
export function useGetComponentsByUserIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
export type GetComponentsByUserIdQueryHookResult = ReturnType<typeof useGetComponentsByUserIdQuery>;
export type GetComponentsByUserIdLazyQueryHookResult = ReturnType<typeof useGetComponentsByUserIdLazyQuery>;
export type GetComponentsByUserIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentsByUserIdSuspenseQuery>;
export type GetComponentsByUserIdQueryResult = Apollo.QueryResult<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>;
export const CreateInvitationDocument = gql`
    mutation CreateInvitation($data: CreateInvitationInputDto!) {
  createInvitation(data: $data) {
    ...Invitation
  }
}
    ${InvitationFragmentDoc}`;
export type CreateInvitationMutationFn = Apollo.MutationFunction<CreateInvitationMutation, CreateInvitationMutationVariables>;

/**
 * __useCreateInvitationMutation__
 *
 * To run a mutation, you first call `useCreateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createInvitationMutation, { data, loading, error }] = useCreateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useCreateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<CreateInvitationMutation, CreateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateInvitationMutation, CreateInvitationMutationVariables>(CreateInvitationDocument, options);
      }
export type CreateInvitationMutationHookResult = ReturnType<typeof useCreateInvitationMutation>;
export type CreateInvitationMutationResult = Apollo.MutationResult<CreateInvitationMutation>;
export type CreateInvitationMutationOptions = Apollo.BaseMutationOptions<CreateInvitationMutation, CreateInvitationMutationVariables>;
export const GetInvitationsDocument = gql`
    query GetInvitations {
  getInvitations {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsQuery__
 *
 * To run a query within a React component, call `useGetInvitationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetInvitationsQuery(baseOptions?: Apollo.QueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
      }
export function useGetInvitationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
export type GetInvitationsQueryHookResult = ReturnType<typeof useGetInvitationsQuery>;
export type GetInvitationsLazyQueryHookResult = ReturnType<typeof useGetInvitationsLazyQuery>;
export type GetInvitationsSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsSuspenseQuery>;
export type GetInvitationsQueryResult = Apollo.QueryResult<GetInvitationsQuery, GetInvitationsQueryVariables>;
export const DeleteInvitationDocument = gql`
    mutation DeleteInvitation($deleteInvitationId: String!) {
  deleteInvitation(id: $deleteInvitationId) {
    id
  }
}
    `;
export type DeleteInvitationMutationFn = Apollo.MutationFunction<DeleteInvitationMutation, DeleteInvitationMutationVariables>;

/**
 * __useDeleteInvitationMutation__
 *
 * To run a mutation, you first call `useDeleteInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteInvitationMutation, { data, loading, error }] = useDeleteInvitationMutation({
 *   variables: {
 *      deleteInvitationId: // value for 'deleteInvitationId'
 *   },
 * });
 */
export function useDeleteInvitationMutation(baseOptions?: Apollo.MutationHookOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteInvitationMutation, DeleteInvitationMutationVariables>(DeleteInvitationDocument, options);
      }
export type DeleteInvitationMutationHookResult = ReturnType<typeof useDeleteInvitationMutation>;
export type DeleteInvitationMutationResult = Apollo.MutationResult<DeleteInvitationMutation>;
export type DeleteInvitationMutationOptions = Apollo.BaseMutationOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>;
export const UpdateInvitationDocument = gql`
    mutation UpdateInvitation($data: UpateInvitationInputDto!) {
  updateInvitation(data: $data) {
    id
  }
}
    `;
export type UpdateInvitationMutationFn = Apollo.MutationFunction<UpdateInvitationMutation, UpdateInvitationMutationVariables>;

/**
 * __useUpdateInvitationMutation__
 *
 * To run a mutation, you first call `useUpdateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateInvitationMutation, { data, loading, error }] = useUpdateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateInvitationMutation, UpdateInvitationMutationVariables>(UpdateInvitationDocument, options);
      }
export type UpdateInvitationMutationHookResult = ReturnType<typeof useUpdateInvitationMutation>;
export type UpdateInvitationMutationResult = Apollo.MutationResult<UpdateInvitationMutation>;
export type UpdateInvitationMutationOptions = Apollo.BaseMutationOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>;
export const GetInvitationsByEmailDocument = gql`
    query GetInvitationsByEmail($email: String!) {
  getInvitationsByEmail(email: $email) {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsByEmailQuery__
 *
 * To run a query within a React component, call `useGetInvitationsByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetInvitationsByEmailQuery(baseOptions: Apollo.QueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables> & ({ variables: GetInvitationsByEmailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
      }
export function useGetInvitationsByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
export type GetInvitationsByEmailQueryHookResult = ReturnType<typeof useGetInvitationsByEmailQuery>;
export type GetInvitationsByEmailLazyQueryHookResult = ReturnType<typeof useGetInvitationsByEmailLazyQuery>;
export type GetInvitationsByEmailSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsByEmailSuspenseQuery>;
export type GetInvitationsByEmailQueryResult = Apollo.QueryResult<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>;
export const GetCountriesDocument = gql`
    query GetCountries {
  getCountries {
    countries {
      name
      code
    }
  }
}
    `;

/**
 * __useGetCountriesQuery__
 *
 * To run a query within a React component, call `useGetCountriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCountriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCountriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetCountriesQuery(baseOptions?: Apollo.QueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
      }
export function useGetCountriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
export type GetCountriesQueryHookResult = ReturnType<typeof useGetCountriesQuery>;
export type GetCountriesLazyQueryHookResult = ReturnType<typeof useGetCountriesLazyQuery>;
export type GetCountriesSuspenseQueryHookResult = ReturnType<typeof useGetCountriesSuspenseQuery>;
export type GetCountriesQueryResult = Apollo.QueryResult<GetCountriesQuery, GetCountriesQueryVariables>;
export const GetCurrencySymbolDocument = gql`
    query GetCurrencySymbol($code: String!) {
  getCurrencySymbol(code: $code) {
    code
    symbol
    image
  }
}
    `;

/**
 * __useGetCurrencySymbolQuery__
 *
 * To run a query within a React component, call `useGetCurrencySymbolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCurrencySymbolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCurrencySymbolQuery({
 *   variables: {
 *      code: // value for 'code'
 *   },
 * });
 */
export function useGetCurrencySymbolQuery(baseOptions: Apollo.QueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables> & ({ variables: GetCurrencySymbolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
      }
export function useGetCurrencySymbolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
export type GetCurrencySymbolQueryHookResult = ReturnType<typeof useGetCurrencySymbolQuery>;
export type GetCurrencySymbolLazyQueryHookResult = ReturnType<typeof useGetCurrencySymbolLazyQuery>;
export type GetCurrencySymbolSuspenseQueryHookResult = ReturnType<typeof useGetCurrencySymbolSuspenseQuery>;
export type GetCurrencySymbolQueryResult = Apollo.QueryResult<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>;
export const GetTimezonesDocument = gql`
    query GetTimezones {
  getTimezones {
    timezones {
      value
      abbr
      offset
      isdst
      text
      utc
    }
  }
}
    `;

/**
 * __useGetTimezonesQuery__
 *
 * To run a query within a React component, call `useGetTimezonesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTimezonesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTimezonesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetTimezonesQuery(baseOptions?: Apollo.QueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
      }
export function useGetTimezonesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
export type GetTimezonesQueryHookResult = ReturnType<typeof useGetTimezonesQuery>;
export type GetTimezonesLazyQueryHookResult = ReturnType<typeof useGetTimezonesLazyQuery>;
export type GetTimezonesSuspenseQueryHookResult = ReturnType<typeof useGetTimezonesSuspenseQuery>;
export type GetTimezonesQueryResult = Apollo.QueryResult<GetTimezonesQuery, GetTimezonesQueryVariables>;
export const UpdateNotificationDocument = gql`
    mutation UpdateNotification($isRead: Boolean!, $updateNotificationId: String!) {
  updateNotification(isRead: $isRead, id: $updateNotificationId) {
    id
  }
}
    `;
export type UpdateNotificationMutationFn = Apollo.MutationFunction<UpdateNotificationMutation, UpdateNotificationMutationVariables>;

/**
 * __useUpdateNotificationMutation__
 *
 * To run a mutation, you first call `useUpdateNotificationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNotificationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNotificationMutation, { data, loading, error }] = useUpdateNotificationMutation({
 *   variables: {
 *      isRead: // value for 'isRead'
 *      updateNotificationId: // value for 'updateNotificationId'
 *   },
 * });
 */
export function useUpdateNotificationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNotificationMutation, UpdateNotificationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNotificationMutation, UpdateNotificationMutationVariables>(UpdateNotificationDocument, options);
      }
export type UpdateNotificationMutationHookResult = ReturnType<typeof useUpdateNotificationMutation>;
export type UpdateNotificationMutationResult = Apollo.MutationResult<UpdateNotificationMutation>;
export type UpdateNotificationMutationOptions = Apollo.BaseMutationOptions<UpdateNotificationMutation, UpdateNotificationMutationVariables>;
export const GetDataUsagesDocument = gql`
    query GetDataUsages($data: SimUsagesInputDto!) {
  getDataUsages(data: $data) {
    usages {
      usage
      simId
    }
  }
}
    `;

/**
 * __useGetDataUsagesQuery__
 *
 * To run a query within a React component, call `useGetDataUsagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDataUsagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDataUsagesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetDataUsagesQuery(baseOptions: Apollo.QueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables> & ({ variables: GetDataUsagesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
      }
export function useGetDataUsagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
export type GetDataUsagesQueryHookResult = ReturnType<typeof useGetDataUsagesQuery>;
export type GetDataUsagesLazyQueryHookResult = ReturnType<typeof useGetDataUsagesLazyQuery>;
export type GetDataUsagesSuspenseQueryHookResult = ReturnType<typeof useGetDataUsagesSuspenseQuery>;
export type GetDataUsagesQueryResult = Apollo.QueryResult<GetDataUsagesQuery, GetDataUsagesQueryVariables>;