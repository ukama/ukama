export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
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

export type App = {
  __typename?: 'App';
  metricsKeys: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  notes: Scalars['String']['output'];
  space: Scalars['String']['output'];
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

export type Apps = {
  __typename?: 'Apps';
  apps: Array<App>;
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

export type GetHealthReportInputDto = {
  id: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  timeframe: Timeframe_Filter;
  timestamp: Scalars['String']['input'];
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

export type GetSoftwaresInput = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  status: SoftwareStatusEnum;
};

export type HealthCappInfo = {
  __typename?: 'HealthCappInfo';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  resources: Array<HealthResourceInfo>;
  space: Scalars['String']['output'];
  status: Scalars['String']['output'];
  tag: Scalars['String']['output'];
};

export type HealthInfo = {
  __typename?: 'HealthInfo';
  capps: Array<HealthCappInfo>;
  id: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  system: Array<HealthSystemInfo>;
  timestamp: Scalars['String']['output'];
};

export type HealthResourceInfo = {
  __typename?: 'HealthResourceInfo';
  cappId: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type HealthSystemInfo = {
  __typename?: 'HealthSystemInfo';
  healthId: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
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
  toggleService: CBooleanResponse;
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
  updateSoftware: StringResponse;
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


export type MutationToggleServiceArgs = {
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


export type MutationUpdateSoftwareArgs = {
  data: UpdateSoftwareInputDto;
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
  Cnode = 'cnode',
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
  getApps?: Maybe<Apps>;
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
  getHealthReport: HealthInfo;
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
  getSoftwares: Softwares;
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


export type QueryGetHealthReportArgs = {
  data: GetHealthReportInputDto;
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


export type QueryGetSoftwaresArgs = {
  data: GetSoftwaresInput;
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

export type Software = {
  __typename?: 'Software';
  changeLog: Array<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  currentVersion: Scalars['String']['output'];
  desiredVersion: Scalars['String']['output'];
  id: Scalars['String']['output'];
  metricsKeys: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  notes: Scalars['String']['output'];
  releaseDate: Scalars['String']['output'];
  space: Scalars['String']['output'];
  status: SoftwareStatusEnum;
  updatedAt: Scalars['String']['output'];
};

/** Software status enums */
export enum SoftwareStatusEnum {
  Unknown = 'unknown',
  UpToDate = 'up_to_date',
  UpdateAvailable = 'update_available',
  UpdateFailed = 'update_failed',
  UpdateInProgress = 'update_in_progress'
}

export type Softwares = {
  __typename?: 'Softwares';
  software: Array<Software>;
};

export type StringResponse = {
  __typename?: 'StringResponse';
  message: Scalars['String']['output'];
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

export enum Timeframe_Filter {
  All = 'ALL',
  Latest = 'LATEST',
  Unknown = 'UNKNOWN'
}

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

export type UpdateSoftwareInputDto = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  tag: Scalars['String']['input'];
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
