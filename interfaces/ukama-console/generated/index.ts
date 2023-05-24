import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  DateTime: any;
};

export enum Alert_Type {
  Error = 'ERROR',
  Info = 'INFO',
  Warning = 'WARNING'
}

export enum Api_Method_Type {
  Delete = 'DELETE',
  Get = 'GET',
  Patch = 'PATCH',
  Post = 'POST',
  Put = 'PUT'
}

export type AddNetworkInputDto = {
  network_name: Scalars['String'];
  org: Scalars['String'];
};

export type AddNodeDto = {
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Scalars['String'];
};

export type AddNodeResponse = {
  __typename?: 'AddNodeResponse';
  attached: Array<Scalars['String']>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Scalars['String'];
  type: Scalars['String'];
};

export type AddOrgInputDto = {
  certificate: Scalars['String'];
  name: Scalars['String'];
  owner_uuid: Scalars['String'];
};

export type AddPackageInputDto = {
  active: Scalars['Boolean'];
  data_volume: Scalars['Int'];
  duration: Scalars['Int'];
  name: Scalars['String'];
  org_id: Scalars['String'];
  org_rates_id: Scalars['Int'];
  sim_type: Scalars['String'];
  sms_volume: Scalars['Int'];
  voice_volume: Scalars['Int'];
};

export type AddPackageSimResDto = {
  __typename?: 'AddPackageSimResDto';
  packageId?: Maybe<Scalars['String']>;
};

export type AddPackageToSimInputDto = {
  packageId: Scalars['String'];
  simId: Scalars['String'];
  startDate: Scalars['String'];
};

export type AddSiteInputDto = {
  site: Scalars['String'];
};

export type AddUserServiceRes = {
  __typename?: 'AddUserServiceRes';
  iccid: Scalars['String'];
  user: OrgUserDto;
};

export type AlertDto = {
  __typename?: 'AlertDto';
  alertDate?: Maybe<Scalars['DateTime']>;
  description?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  title?: Maybe<Scalars['String']>;
  type: Alert_Type;
};

export type AlertResponse = {
  __typename?: 'AlertResponse';
  data: Array<AlertDto>;
  length: Scalars['Float'];
  status: Scalars['String'];
};

export type AlertsResponse = {
  __typename?: 'AlertsResponse';
  alerts: Array<AlertDto>;
  meta: Meta;
};

export type AllocateSimInputDto = {
  iccid: Scalars['String'];
  networkId: Scalars['String'];
  packageId: Scalars['String'];
  simType: Scalars['String'];
  subscriberId: Scalars['String'];
};

export type ApiMethodDataDto = {
  __typename?: 'ApiMethodDataDto';
  body?: Maybe<Scalars['String']>;
  headers?: Maybe<Scalars['String']>;
  params?: Maybe<Scalars['String']>;
  path: Scalars['String'];
  type: Api_Method_Type;
};

export type AttachedNodes = {
  __typename?: 'AttachedNodes';
  nodeId: Scalars['String'];
};

export type AuthType = {
  __typename?: 'AuthType';
  Authorization: Scalars['String'];
  Cookie: Scalars['String'];
};

export type BillHistoryDto = {
  __typename?: 'BillHistoryDto';
  date: Scalars['String'];
  description: Scalars['String'];
  id: Scalars['String'];
  subtotal: Scalars['Float'];
  totalUsage: Scalars['Float'];
};

export type BillHistoryResponse = {
  __typename?: 'BillHistoryResponse';
  data: Array<BillHistoryDto>;
  status: Scalars['String'];
};

export type BillResponse = {
  __typename?: 'BillResponse';
  bill: Array<CurrentBillDto>;
  billMonth: Scalars['String'];
  dueDate: Scalars['String'];
  total: Scalars['Float'];
};

export type BoolResponse = {
  __typename?: 'BoolResponse';
  success: Scalars['Boolean'];
};

export type ConnectedUserDto = {
  __typename?: 'ConnectedUserDto';
  totalUser: Scalars['String'];
};

export type CreateCustomerDto = {
  email: Scalars['String'];
  name: Scalars['String'];
};

export type CurrentBillDto = {
  __typename?: 'CurrentBillDto';
  dataUsed: Scalars['Float'];
  id: Scalars['String'];
  name: Scalars['String'];
  rate: Scalars['Float'];
  subtotal: Scalars['Float'];
};

export type CurrentBillResponse = {
  __typename?: 'CurrentBillResponse';
  data: Array<CurrentBillDto>;
  status: Scalars['String'];
};

export enum Data_Bill_Filter {
  April = 'APRIL',
  August = 'AUGUST',
  Current = 'CURRENT',
  December = 'DECEMBER',
  Februray = 'FEBRURAY',
  January = 'JANUARY',
  July = 'JULY',
  June = 'JUNE',
  March = 'MARCH',
  May = 'MAY',
  Novermber = 'NOVERMBER',
  October = 'OCTOBER',
  September = 'SEPTEMBER'
}

export type DataBillDto = {
  __typename?: 'DataBillDto';
  billDue: Scalars['Float'];
  dataBill: Scalars['Float'];
  id: Scalars['String'];
};

export type DataBillResponse = {
  __typename?: 'DataBillResponse';
  data: DataBillDto;
  status: Scalars['String'];
};

export type DataUsageDto = {
  __typename?: 'DataUsageDto';
  dataConsumed: Scalars['Float'];
  dataPackage: Scalars['String'];
  id: Scalars['String'];
};

export type DataUsageInputDto = {
  ids: Array<Scalars['String']>;
};

export type DataUsageNetworkResponse = {
  __typename?: 'DataUsageNetworkResponse';
  usage: Scalars['Float'];
};

export type DataUsageResponse = {
  __typename?: 'DataUsageResponse';
  data: DataUsageDto;
  status: Scalars['String'];
};

export type DeactivateResponse = {
  __typename?: 'DeactivateResponse';
  email: Scalars['String'];
  isDeactivated: Scalars['Boolean'];
  name: Scalars['String'];
  uuid: Scalars['String'];
};

export type DefaultMarkupApiResDto = {
  __typename?: 'DefaultMarkupAPIResDto';
  markup: Scalars['Float'];
};

export type DefaultMarkupHistoryApiResDto = {
  __typename?: 'DefaultMarkupHistoryAPIResDto';
  markupRates?: Maybe<Array<DefaultMarkupHistoryDto>>;
};

export type DefaultMarkupHistoryDto = {
  __typename?: 'DefaultMarkupHistoryDto';
  Markup: Scalars['Float'];
  createdAt: Scalars['String'];
  deletedAt: Scalars['String'];
};

export type DefaultMarkupHistoryResDto = {
  __typename?: 'DefaultMarkupHistoryResDto';
  markupRates?: Maybe<Array<DefaultMarkupHistoryDto>>;
};

export type DefaultMarkupInputDto = {
  markup: Scalars['Float'];
};

export type DefaultMarkupResDto = {
  __typename?: 'DefaultMarkupResDto';
  markup: Scalars['Float'];
};

export type DeleteNodeRes = {
  __typename?: 'DeleteNodeRes';
  nodeId: Scalars['String'];
};

export type DeleteSimInputDto = {
  simId: Scalars['String'];
};

export type DeleteSimResDto = {
  __typename?: 'DeleteSimResDto';
  simId?: Maybe<Scalars['String']>;
};

export type ESimQrCodeRes = {
  __typename?: 'ESimQRCodeRes';
  qrCode: Scalars['String'];
};

export type ErrorType = {
  __typename?: 'ErrorType';
  code: Scalars['Float'];
  description?: Maybe<Scalars['String']>;
  message: Scalars['String'];
};

export type EsimDto = {
  __typename?: 'EsimDto';
  active: Scalars['Boolean'];
  esim: Scalars['String'];
};

export type EsimResponse = {
  __typename?: 'EsimResponse';
  data: Array<EsimDto>;
  status: Scalars['String'];
};

export enum Get_User_Status_Type {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
  Unknown = 'UNKNOWN'
}

export enum Graphs_Tab {
  Home = 'HOME',
  Network = 'NETWORK',
  NodeStatus = 'NODE_STATUS',
  Overview = 'OVERVIEW',
  Radio = 'RADIO',
  Resources = 'RESOURCES'
}

export type GetAccountDetailsDto = {
  __typename?: 'GetAccountDetailsDto';
  email: Scalars['String'];
  isFirstVisit: Scalars['Boolean'];
};

export type GetESimQrCodeInput = {
  simId: Scalars['String'];
  userId: Scalars['String'];
};

export type GetMetricsRes = {
  __typename?: 'GetMetricsRes';
  metrics: Array<MetricRes>;
  next: Scalars['Boolean'];
  to: Scalars['Float'];
};

export type GetNodeStatusInput = {
  nodeId: Scalars['String'];
  nodeType: Node_Type;
};

export type GetNodeStatusRes = {
  __typename?: 'GetNodeStatusRes';
  status: Org_Node_State;
  uptime: Scalars['Float'];
};

export type GetPackagesForSimInputDto = {
  simId: Scalars['String'];
};

export type GetPackagesForSimResDto = {
  __typename?: 'GetPackagesForSimResDto';
  Packages?: Maybe<Array<SimPackageDto>>;
};

export type GetSimApiResDto = {
  __typename?: 'GetSimAPIResDto';
  sim: SimDetailsDto;
};

export type GetSimByNetworkInputDto = {
  networkId: Scalars['String'];
};

export type GetSimBySubscriberIdInputDto = {
  subscriberId: Scalars['String'];
};

export type GetSimInputDto = {
  simId: Scalars['String'];
};

export type GetUserDto = {
  __typename?: 'GetUserDto';
  email: Scalars['String'];
  id: Scalars['String'];
  name: Scalars['String'];
  phone: Scalars['String'];
};

export type GetUserResponseDto = {
  __typename?: 'GetUserResponseDto';
  data: Array<GetUserDto>;
  length: Scalars['Float'];
  status: Scalars['String'];
};

export type GetUsersDto = {
  __typename?: 'GetUsersDto';
  dataPlan?: Maybe<Scalars['String']>;
  dataUsage?: Maybe<Scalars['String']>;
  email: Scalars['String'];
  id: Scalars['String'];
  name: Scalars['String'];
};

export type IdResponse = {
  __typename?: 'IdResponse';
  uuid: Scalars['String'];
};

export type LinkNodes = {
  __typename?: 'LinkNodes';
  attachedNodeIds: Array<Scalars['String']>;
  nodeId: Scalars['String'];
};

export type MemberApiObj = {
  __typename?: 'MemberAPIObj';
  is_deactivated: Scalars['Boolean'];
  member_since: Scalars['String'];
  org_id: Scalars['String'];
  user_id: Scalars['String'];
  uuid: Scalars['String'];
};

export type MemberObj = {
  __typename?: 'MemberObj';
  isDeactivated: Scalars['Boolean'];
  memberSince: Scalars['String'];
  orgId: Scalars['String'];
  userId: Scalars['String'];
  uuid: Scalars['String'];
};

export type Meta = {
  __typename?: 'Meta';
  count: Scalars['Float'];
  page: Scalars['Float'];
  pages: Scalars['Float'];
  size: Scalars['Float'];
};

export type MetricDto = {
  __typename?: 'MetricDto';
  x: Scalars['Float'];
  y: Scalars['Float'];
};

export type MetricInfo = {
  __typename?: 'MetricInfo';
  org: Scalars['String'];
};

export type MetricLatestValueRes = {
  __typename?: 'MetricLatestValueRes';
  time: Scalars['String'];
  value: Scalars['String'];
};

export type MetricRes = {
  __typename?: 'MetricRes';
  data: Array<MetricDto>;
  name: Scalars['String'];
  next: Scalars['Boolean'];
  type: Scalars['String'];
};

export type MetricServiceValueRes = {
  __typename?: 'MetricServiceValueRes';
  metric: MetricInfo;
  value: Array<MetricValues>;
};

export type MetricValues = {
  __typename?: 'MetricValues';
  x: Scalars['Float'];
  y: Scalars['String'];
};

export type MetricsByTabInputDto = {
  from: Scalars['Float'];
  nodeId: Scalars['String'];
  nodeType: Node_Type;
  regPolling: Scalars['Boolean'];
  step: Scalars['Float'];
  tab: Graphs_Tab;
  to: Scalars['Float'];
};

export type MetricsInputDto = {
  from: Scalars['Float'];
  nodeId: Scalars['String'];
  orgId: Scalars['String'];
  regPolling: Scalars['Boolean'];
  step: Scalars['Float'];
  to: Scalars['Float'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addMember: MemberObj;
  addNetwork: NetworkDto;
  addNode: AddNodeResponse;
  addOrg: OrgDto;
  addPackage: PackageDto;
  addSubscriber: SubscriberDto;
  addUser: UserResDto;
  allocateSim: SimResDto;
  attachPaymentWithCustomer: Scalars['Boolean'];
  createCustomer: StripeCustomer;
  deactivateUser: UserResDto;
  defaultMarkup: BoolResponse;
  deleteNode: DeleteNodeRes;
  deletePackage: IdResponse;
  deleteSubscriber: BoolResponse;
  deleteUser: BoolResponse;
  getSim: SetActivePackageForSimResDto;
  removeMember: BoolResponse;
  toggleSimStatus: SimStatusResDto;
  updateFirstVisit: UserFistVisitResDto;
  updateMember: BoolResponse;
  updateNode: UpdateNodeResponse;
  updatePackage: PackageDto;
  updateSubscriber: BoolResponse;
  updateUser: UserResDto;
  updateUserRoaming: OrgUserSimDto;
  updateUserStatus: OrgUserSimDto;
};


export type MutationAddMemberArgs = {
  userId: Scalars['String'];
};


export type MutationAddNetworkArgs = {
  data: AddNetworkInputDto;
};


export type MutationAddNodeArgs = {
  data: AddNodeDto;
};


export type MutationAddOrgArgs = {
  data: AddOrgInputDto;
};


export type MutationAddPackageArgs = {
  data: AddPackageInputDto;
};


export type MutationAddSubscriberArgs = {
  data: SubscriberInputDto;
};


export type MutationAddUserArgs = {
  data: UserInputDto;
};


export type MutationAllocateSimArgs = {
  data: AllocateSimInputDto;
};


export type MutationAttachPaymentWithCustomerArgs = {
  paymentId: Scalars['String'];
};


export type MutationCreateCustomerArgs = {
  data: CreateCustomerDto;
};


export type MutationDeactivateUserArgs = {
  uuid: Scalars['String'];
};


export type MutationDefaultMarkupArgs = {
  data: DefaultMarkupInputDto;
};


export type MutationDeleteNodeArgs = {
  id: Scalars['String'];
};


export type MutationDeletePackageArgs = {
  packageId: Scalars['String'];
};


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String'];
};


export type MutationDeleteUserArgs = {
  userId: Scalars['String'];
};


export type MutationGetSimArgs = {
  data: SetActivePackageForSimInputDto;
};


export type MutationToggleSimStatusArgs = {
  data: ToggleSimStatusInputDto;
};


export type MutationUpdateFirstVisitArgs = {
  data: UserFistVisitInputDto;
};


export type MutationUpdateMemberArgs = {
  data: UpdateMemberInputDto;
  memberId: Scalars['String'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeDto;
};


export type MutationUpdatePackageArgs = {
  data: UpdatePackageInputDto;
  packageId: Scalars['String'];
};


export type MutationUpdateSubscriberArgs = {
  data: UpdateSubscriberInputDto;
  subscriberId: Scalars['String'];
};


export type MutationUpdateUserArgs = {
  data: UpdateUserInputDto;
  userId: Scalars['String'];
};


export type MutationUpdateUserRoamingArgs = {
  data: UpdateUserServiceInput;
};


export type MutationUpdateUserStatusArgs = {
  data: UpdateUserServiceInput;
};

export enum Network_Status {
  Down = 'DOWN',
  Online = 'ONLINE',
  Undefined = 'UNDEFINED'
}

export enum Node_Type {
  Amplifier = 'AMPLIFIER',
  Home = 'HOME',
  Tower = 'TOWER'
}

export type NetworkApiDto = {
  __typename?: 'NetworkAPIDto';
  created_at: Scalars['String'];
  id: Scalars['String'];
  is_deactivated: Scalars['String'];
  name: Scalars['String'];
  org_id: Scalars['String'];
};

export type NetworkApiResDto = {
  __typename?: 'NetworkAPIResDto';
  network: NetworkApiDto;
};

export type NetworkDto = {
  __typename?: 'NetworkDto';
  createdAt: Scalars['String'];
  id: Scalars['String'];
  isDeactivated: Scalars['String'];
  name: Scalars['String'];
  orgId: Scalars['String'];
};

export type NetworkStatusDto = {
  __typename?: 'NetworkStatusDto';
  liveNode: Scalars['Float'];
  status: Network_Status;
  totalNodes: Scalars['Float'];
};

export type NetworkStatusResponse = {
  __typename?: 'NetworkStatusResponse';
  data: NetworkStatusDto;
  status: Scalars['String'];
};

export type NetworksApiResDto = {
  __typename?: 'NetworksAPIResDto';
  networks: Array<NetworkApiDto>;
  org_id: Scalars['String'];
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
  orgId: Scalars['String'];
};

export type NodeAppResponse = {
  __typename?: 'NodeAppResponse';
  cpu: Scalars['String'];
  id: Scalars['String'];
  memory: Scalars['String'];
  title: Scalars['String'];
  version: Scalars['String'];
};

export type NodeAppsVersionLogsResponse = {
  __typename?: 'NodeAppsVersionLogsResponse';
  date: Scalars['Float'];
  notes: Scalars['String'];
  version: Scalars['String'];
};

export type NodeDto = {
  __typename?: 'NodeDto';
  description: Scalars['String'];
  id: Scalars['String'];
  isUpdateAvailable: Scalars['Boolean'];
  name: Scalars['String'];
  status: Org_Node_State;
  totalUser: Scalars['Float'];
  type: Scalars['String'];
  updateDescription: Scalars['String'];
  updateShortNote: Scalars['String'];
  updateVersion: Scalars['String'];
};

export type NodeObj = {
  attached?: InputMaybe<Array<NodeObj>>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Scalars['String'];
};

export type NodeResponse = {
  __typename?: 'NodeResponse';
  attached: Array<OrgNodeDto>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Org_Node_State;
  type: Node_Type;
};

export type NodeStatsResponse = {
  __typename?: 'NodeStatsResponse';
  claimCount: Scalars['Float'];
  totalCount: Scalars['Float'];
  upCount: Scalars['Float'];
};

export enum Org_Node_State {
  Error = 'ERROR',
  Onboarded = 'ONBOARDED',
  Pending = 'PENDING',
  Undefined = 'UNDEFINED'
}

export type OrgApiDto = {
  __typename?: 'OrgAPIDto';
  certificate: Scalars['Boolean'];
  created_at: Scalars['String'];
  id: Scalars['String'];
  is_deactivated: Scalars['Boolean'];
  name: Scalars['String'];
  owner: Scalars['String'];
};

export type OrgApiResDto = {
  __typename?: 'OrgAPIResDto';
  org: OrgApiDto;
};

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['Boolean'];
  createdAt: Scalars['String'];
  id: Scalars['String'];
  isDeactivated: Scalars['Boolean'];
  name: Scalars['String'];
  owner: Scalars['String'];
};

export type OrgMemberApiResDto = {
  __typename?: 'OrgMemberAPIResDto';
  member: MemberApiObj;
};

export type OrgMembersApiResDto = {
  __typename?: 'OrgMembersAPIResDto';
  members: Array<MemberApiObj>;
  org: Scalars['String'];
};

export type OrgMembersResDto = {
  __typename?: 'OrgMembersResDto';
  members: Array<MemberObj>;
  org: Scalars['String'];
};

export type OrgMetricValueDto = {
  __typename?: 'OrgMetricValueDto';
  x: Scalars['Float'];
  y: Scalars['String'];
};

export type OrgNodeDto = {
  __typename?: 'OrgNodeDto';
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Org_Node_State;
  type: Node_Type;
};

export type OrgNodeResponse = {
  __typename?: 'OrgNodeResponse';
  nodes: Array<OrgNodeDto>;
  orgName: Scalars['String'];
};

export type OrgNodeResponseDto = {
  __typename?: 'OrgNodeResponseDto';
  activeNodes: Scalars['Float'];
  nodes: Array<NodeDto>;
  orgId: Scalars['String'];
  totalNodes: Scalars['Float'];
};

export type OrgResDto = {
  __typename?: 'OrgResDto';
  org: OrgDto;
};

export type OrgUserDto = {
  __typename?: 'OrgUserDto';
  email: Scalars['String'];
  isDeactivated: Scalars['Boolean'];
  name: Scalars['String'];
  uuid: Scalars['String'];
};

export type OrgUserResponse = {
  __typename?: 'OrgUserResponse';
  sim: OrgUserSimDto;
  user: OrgUserDto;
};

export type OrgUserSimDto = {
  __typename?: 'OrgUserSimDto';
  carrier: UserServicesDto;
  iccid: Scalars['String'];
  isPhysical: Scalars['Boolean'];
  ukama: UserServicesDto;
};

export type OrgUsersResponse = {
  __typename?: 'OrgUsersResponse';
  org: Scalars['String'];
  users: Array<OrgUserDto>;
};

export type OrgsApiResDto = {
  __typename?: 'OrgsAPIResDto';
  orgs: Array<OrgApiDto>;
  owner: Scalars['String'];
};

export type OrgsResDto = {
  __typename?: 'OrgsResDto';
  orgs: Array<OrgDto>;
  owner: Scalars['String'];
};

export type PackageApiDto = {
  __typename?: 'PackageAPIDto';
  active: Scalars['Boolean'];
  amount: Scalars['Float'];
  apn: Scalars['String'];
  country: Scalars['String'];
  created_at: Scalars['String'];
  currency: Scalars['String'];
  data_unit: Scalars['String'];
  data_volume: Scalars['String'];
  deleted_at: Scalars['String'];
  dlbr: Scalars['String'];
  duration: Scalars['String'];
  flatrate: Scalars['Boolean'];
  from: Scalars['String'];
  markup: PackageMarkupApiDto;
  message_unit: Scalars['String'];
  name: Scalars['String'];
  org_id: Scalars['String'];
  owner_id: Scalars['String'];
  provider: Scalars['String'];
  rate: PackageRateApiDto;
  sim_type: Scalars['String'];
  sms_volume: Scalars['String'];
  to: Scalars['String'];
  type: Scalars['String'];
  ulbr: Scalars['String'];
  updated_at: Scalars['String'];
  uuid: Scalars['String'];
  voice_unit: Scalars['String'];
  voice_volume: Scalars['String'];
};

export type PackageApiResDto = {
  __typename?: 'PackageAPIResDto';
  package: PackageApiDto;
};

export type PackageDto = {
  __typename?: 'PackageDto';
  active: Scalars['Boolean'];
  amount: Scalars['Float'];
  apn: Scalars['String'];
  country: Scalars['String'];
  createdAt: Scalars['String'];
  currency: Scalars['String'];
  dataUnit: Scalars['String'];
  dataVolume: Scalars['String'];
  deletedAt: Scalars['String'];
  dlbr: Scalars['String'];
  duration: Scalars['String'];
  flatrate: Scalars['Boolean'];
  from: Scalars['String'];
  markup: PackageMarkupApiDto;
  messageUnit: Scalars['String'];
  name: Scalars['String'];
  orgId: Scalars['String'];
  ownerId: Scalars['String'];
  provider: Scalars['String'];
  rate: PackageRateApiDto;
  simType: Scalars['String'];
  smsVolume: Scalars['String'];
  to: Scalars['String'];
  type: Scalars['String'];
  ulbr: Scalars['String'];
  updatedAt: Scalars['String'];
  uuid: Scalars['String'];
  voiceUnit: Scalars['String'];
  voiceVolume: Scalars['String'];
};

export type PackageMarkupApiDto = {
  __typename?: 'PackageMarkupAPIDto';
  baserate: Scalars['String'];
  markup: Scalars['Float'];
};

export type PackageMarkupDto = {
  __typename?: 'PackageMarkupDto';
  baserate: Scalars['String'];
  markup: Scalars['Float'];
};

export type PackageRateApiDto = {
  __typename?: 'PackageRateAPIDto';
  amount: Scalars['Float'];
  data: Scalars['Float'];
  sms_mo: Scalars['String'];
  sms_mt: Scalars['Float'];
};

export type PackageRateDto = {
  __typename?: 'PackageRateDto';
  amount: Scalars['Float'];
  data: Scalars['Float'];
  smsMo: Scalars['String'];
  smsMt: Scalars['Float'];
};

export type PackagesApiResDto = {
  __typename?: 'PackagesAPIResDto';
  packages: Array<PackageApiDto>;
};

export type PackagesResDto = {
  __typename?: 'PackagesResDto';
  packages: Array<PackageDto>;
};

export type PaginationDto = {
  pageNo: Scalars['Float'];
  pageSize: Scalars['Float'];
};

export type PaginationResponse = {
  __typename?: 'PaginationResponse';
  meta: Meta;
};

export type Query = {
  __typename?: 'Query';
  addSite: SiteDto;
  getAccountDetails: GetAccountDetailsDto;
  getAlerts: AlertsResponse;
  getBillHistory: Array<BillHistoryDto>;
  getConnectedUsers: ConnectedUserDto;
  getCurrentBill: BillResponse;
  getDataBill: DataBillDto;
  getDataUsage: SimDataUsage;
  getDefaultMarkup: DefaultMarkupResDto;
  getDefaultMarkupHistory: DefaultMarkupHistoryResDto;
  getEsimQR: ESimQrCodeRes;
  getEsims: Array<EsimDto>;
  getMetricsByTab: GetMetricsRes;
  getNetwork: NetworkDto;
  getNetworkDataUsage: DataUsageNetworkResponse;
  getNetworkNodesStat: NodeStatsResponse;
  getNetworkStatus: NetworkStatusDto;
  getNetworks: NetworksResDto;
  getNode: NodeResponse;
  getNodeApps: Array<NodeAppResponse>;
  getNodeAppsVersionLogs: Array<NodeAppsVersionLogsResponse>;
  getNodeStatus: GetNodeStatusRes;
  getNodesByOrg: OrgNodeResponseDto;
  getOrg: OrgDto;
  getOrgMember: MemberObj;
  getOrgMembers: OrgMembersResDto;
  getOrgs: OrgsResDto;
  getPackage: PackageDto;
  getPackages: PackagesResDto;
  getSim: GetPackagesForSimResDto;
  getSite: SiteDto;
  getSites: SitesResDto;
  getStripeCustomer: StripeCustomer;
  getSubscriber: SubscriberDto;
  getSubscriberMetricsByNetwork: SubscriberMetricsByNetworkDto;
  getSubscribersByNetwork: SubscribersResDto;
  getUser: UserResDto;
  getUsersDataUsage: Array<GetUserDto>;
  retrivePaymentMethods: Array<StripePaymentMethods>;
  whoami: WhoamiDto;
};


export type QueryAddSiteArgs = {
  data: AddSiteInputDto;
  networkId: Scalars['String'];
};


export type QueryGetAlertsArgs = {
  data: PaginationDto;
};


export type QueryGetConnectedUsersArgs = {
  filter: Time_Filter;
};


export type QueryGetDataBillArgs = {
  filter: Data_Bill_Filter;
};


export type QueryGetDataUsageArgs = {
  simId: Scalars['String'];
};


export type QueryGetEsimQrArgs = {
  data: GetESimQrCodeInput;
};


export type QueryGetMetricsByTabArgs = {
  data: MetricsByTabInputDto;
};


export type QueryGetNetworkArgs = {
  networkId: Scalars['String'];
};


export type QueryGetNetworkDataUsageArgs = {
  filter: Time_Filter;
  networkId: Scalars['String'];
};


export type QueryGetNetworkNodesStatArgs = {
  networkId: Scalars['String'];
};


export type QueryGetNodeArgs = {
  nodeId: Scalars['String'];
};


export type QueryGetNodeStatusArgs = {
  data: GetNodeStatusInput;
};


export type QueryGetOrgArgs = {
  orgName: Scalars['String'];
};


export type QueryGetPackageArgs = {
  packageId: Scalars['String'];
};


export type QueryGetSimArgs = {
  data: GetPackagesForSimInputDto;
};


export type QueryGetSiteArgs = {
  networkId: Scalars['String'];
  siteId: Scalars['String'];
};


export type QueryGetSitesArgs = {
  networkId: Scalars['String'];
};


export type QueryGetSubscriberArgs = {
  subscriberId: Scalars['String'];
};


export type QueryGetSubscriberMetricsByNetworkArgs = {
  networkId: Scalars['String'];
};


export type QueryGetSubscribersByNetworkArgs = {
  networkId: Scalars['String'];
};


export type QueryGetUserArgs = {
  userId: Scalars['String'];
};


export type QueryGetUsersDataUsageArgs = {
  data: DataUsageInputDto;
};

export type RemovePackageFormSimInputDto = {
  packageId: Scalars['String'];
  simId: Scalars['String'];
};

export type RemovePackageFromSimResDto = {
  __typename?: 'RemovePackageFromSimResDto';
  packageId?: Maybe<Scalars['String']>;
};

export type SetActivePackageForSimInputDto = {
  packageId: Scalars['String'];
  simId: Scalars['String'];
};

export type SetActivePackageForSimResDto = {
  __typename?: 'SetActivePackageForSimResDto';
  packageId?: Maybe<Scalars['String']>;
};

export type SimApiDto = {
  __typename?: 'SimAPIDto';
  activationCode: Scalars['String'];
  createdAt: Scalars['String'];
  iccid: Scalars['String'];
  id: Scalars['String'];
  isAllocated: Scalars['String'];
  isPhysical: Scalars['String'];
  msisdn: Scalars['String'];
  qrCode: Scalars['String'];
  simType: Scalars['String'];
  smDpAddress: Scalars['String'];
};

export type SimApiResDto = {
  __typename?: 'SimAPIResDto';
  sim: SimApiDto;
};

export type SimDataUsage = {
  __typename?: 'SimDataUsage';
  usage: Scalars['String'];
};

export type SimDetailsDto = {
  __typename?: 'SimDetailsDto';
  Package: SimPackageDto;
  activationsCount: Scalars['Float'];
  allocatedAt: Scalars['String'];
  deactivationsCount: Scalars['Float'];
  firstActivatedOn: Scalars['String'];
  iccid: Scalars['String'];
  id: Scalars['String'];
  imsi: Scalars['String'];
  isPhysical: Scalars['Boolean'];
  lastActivatedOn: Scalars['String'];
  msisdn: Scalars['String'];
  networkId: Scalars['String'];
  orgId: Scalars['String'];
  status: Scalars['String'];
  subscriberId: Scalars['String'];
  type: Scalars['String'];
};

export type SimDto = {
  __typename?: 'SimDto';
  activationsCount: Scalars['String'];
  allocatedAt: Scalars['String'];
  deactivationsCount: Scalars['String'];
  firstActivatedOn?: Maybe<Scalars['String']>;
  iccid: Scalars['String'];
  id: Scalars['String'];
  imsi: Scalars['String'];
  isPhysical?: Maybe<Scalars['Boolean']>;
  lastActivatedOn?: Maybe<Scalars['String']>;
  msisdn: Scalars['String'];
  networkId: Scalars['String'];
  orgId: Scalars['String'];
  package?: Maybe<Scalars['String']>;
  status: Scalars['String'];
  subscriberId: Scalars['String'];
  type: Scalars['String'];
};

export type SimPackageDto = {
  __typename?: 'SimPackageDto';
  createdAt: Scalars['String'];
  description: Scalars['String'];
  id: Scalars['String'];
  name: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type SimResDto = {
  __typename?: 'SimResDto';
  activationCode: Scalars['String'];
  createdAt: Scalars['String'];
  iccid: Scalars['String'];
  id: Scalars['String'];
  isAllocated: Scalars['String'];
  isPhysical: Scalars['String'];
  msisdn: Scalars['String'];
  qrCode: Scalars['String'];
  simType: Scalars['String'];
  smDpAddress: Scalars['String'];
};

export type SimStatusResDto = {
  __typename?: 'SimStatusResDto';
  simId?: Maybe<Scalars['String']>;
};

export type SiteApiDto = {
  __typename?: 'SiteAPIDto';
  created_at: Scalars['String'];
  id: Scalars['String'];
  is_deactivated: Scalars['String'];
  name: Scalars['String'];
  network_id: Scalars['String'];
};

export type SiteApiResDto = {
  __typename?: 'SiteAPIResDto';
  site: SiteApiDto;
};

export type SiteDto = {
  __typename?: 'SiteDto';
  createdAt: Scalars['String'];
  id: Scalars['String'];
  isDeactivated: Scalars['String'];
  name: Scalars['String'];
  networkId: Scalars['String'];
};

export type SitesApiResDto = {
  __typename?: 'SitesAPIResDto';
  network_id: Scalars['String'];
  sites: Array<SiteApiDto>;
};

export type SitesResDto = {
  __typename?: 'SitesResDto';
  networkId: Scalars['String'];
  sites: Array<SiteDto>;
};

export type StripeCustomer = {
  __typename?: 'StripeCustomer';
  email: Scalars['String'];
  id: Scalars['String'];
  name: Scalars['String'];
};

export type StripePaymentMethods = {
  __typename?: 'StripePaymentMethods';
  brand: Scalars['String'];
  country?: Maybe<Scalars['String']>;
  created: Scalars['Float'];
  cvc_check?: Maybe<Scalars['String']>;
  exp_month: Scalars['Float'];
  exp_year: Scalars['Float'];
  funding: Scalars['String'];
  id: Scalars['String'];
  last4: Scalars['String'];
  type: Scalars['String'];
};

export type SubSimApiDto = {
  __typename?: 'SubSimAPIDto';
  activations_count: Scalars['String'];
  allocated_at: Scalars['String'];
  deactivations_count: Scalars['String'];
  first_activated_on?: Maybe<Scalars['String']>;
  iccid: Scalars['String'];
  id: Scalars['String'];
  imsi: Scalars['String'];
  is_physical?: Maybe<Scalars['Boolean']>;
  last_activated_on?: Maybe<Scalars['String']>;
  msisdn: Scalars['String'];
  network_id: Scalars['String'];
  org_id: Scalars['String'];
  package?: Maybe<Scalars['String']>;
  status: Scalars['String'];
  subscriber_id: Scalars['String'];
  type: Scalars['String'];
};

export type SubscriberApiDto = {
  __typename?: 'SubscriberAPIDto';
  address: Scalars['String'];
  date_of_birth: Scalars['String'];
  email: Scalars['String'];
  first_name: Scalars['String'];
  gender: Scalars['String'];
  id_serial: Scalars['String'];
  last_name: Scalars['String'];
  network_id: Scalars['String'];
  org_id: Scalars['String'];
  phone_number: Scalars['String'];
  proof_of_identification: Scalars['String'];
  sim: Array<SubSimApiDto>;
  subscriber_id: Scalars['String'];
};

export type SubscriberApiResDto = {
  __typename?: 'SubscriberAPIResDto';
  subscriber: SubscriberApiDto;
};

export type SubscriberDto = {
  __typename?: 'SubscriberDto';
  address: Scalars['String'];
  dob: Scalars['String'];
  email: Scalars['String'];
  firstName: Scalars['String'];
  gender: Scalars['String'];
  idSerial: Scalars['String'];
  lastName: Scalars['String'];
  networkId: Scalars['String'];
  orgId: Scalars['String'];
  phone: Scalars['String'];
  proofOfIdentification: Scalars['String'];
  sim: Array<SimDto>;
  uuid: Scalars['String'];
};

export type SubscriberInputDto = {
  address: Scalars['String'];
  dob: Scalars['String'];
  email: Scalars['String'];
  first_name: Scalars['String'];
  gender: Scalars['String'];
  id_serial: Scalars['String'];
  last_name: Scalars['String'];
  network_id: Scalars['String'];
  org_id: Scalars['String'];
  phone: Scalars['String'];
  proof_of_identification: Scalars['String'];
};

export type SubscriberMetricsByNetworkDto = {
  __typename?: 'SubscriberMetricsByNetworkDto';
  active: Scalars['Float'];
  inactive: Scalars['Float'];
  terminated: Scalars['Float'];
  total: Scalars['Float'];
};

export type SubscribersResDto = {
  __typename?: 'SubscribersResDto';
  subscribers: Array<SubscriberDto>;
};

export type Subscription = {
  __typename?: 'Subscription';
  getAlerts: AlertDto;
  getConnectedUsers: ConnectedUserDto;
  getDataBill: DataBillDto;
  getDataUsage: DataUsageDto;
  getMetricsByTab: Array<MetricRes>;
  getNetworkStatus: NetworkDto;
  getUsersDataUsage: GetUserDto;
};

export type THeaders = {
  __typename?: 'THeaders';
  auth: AuthType;
  orgId: Scalars['String'];
  orgName: Scalars['String'];
  userId: Scalars['String'];
};

export enum Time_Filter {
  Month = 'MONTH',
  Today = 'TODAY',
  Total = 'TOTAL',
  Week = 'WEEK'
}

export type ToggleSimStatusInputDto = {
  simId: Scalars['String'];
  status: Scalars['String'];
};

export type UpdateMemberInputDto = {
  isDeactivated: Scalars['Boolean'];
};

export type UpdateNodeDto = {
  name: Scalars['String'];
  nodeId: Scalars['String'];
};

export type UpdateNodeResponse = {
  __typename?: 'UpdateNodeResponse';
  attached: Array<Scalars['String']>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Scalars['String'];
  type: Scalars['String'];
};

export type UpdatePackageInputDto = {
  active?: InputMaybe<Scalars['Boolean']>;
  data_volume?: InputMaybe<Scalars['Int']>;
  duration?: InputMaybe<Scalars['Int']>;
  name?: InputMaybe<Scalars['String']>;
  org_rates_id?: InputMaybe<Scalars['Int']>;
  sim_type?: InputMaybe<Scalars['String']>;
  sms_volume?: InputMaybe<Scalars['Int']>;
  voice_volume?: InputMaybe<Scalars['Int']>;
};

export type UpdateSubscriberInputDto = {
  address?: InputMaybe<Scalars['String']>;
  email?: InputMaybe<Scalars['String']>;
  id_serial?: InputMaybe<Scalars['String']>;
  phone?: InputMaybe<Scalars['String']>;
  proof_of_identification?: InputMaybe<Scalars['String']>;
};

export type UpdateUserInputDto = {
  email?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  phone?: InputMaybe<Scalars['String']>;
};

export type UpdateUserServiceInput = {
  simId: Scalars['String'];
  status: Scalars['Boolean'];
  userId: Scalars['String'];
};

export type UserApiObj = {
  __typename?: 'UserAPIObj';
  email: Scalars['String'];
  is_deactivated: Scalars['Boolean'];
  name: Scalars['String'];
  phone: Scalars['String'];
  registered_since: Scalars['String'];
  uuid: Scalars['String'];
};

export type UserApiResDto = {
  __typename?: 'UserAPIResDto';
  user: Array<UserApiObj>;
};

export type UserDataUsageDto = {
  __typename?: 'UserDataUsageDto';
  dataAllowanceBytes?: Maybe<Scalars['String']>;
  dataUsedBytes?: Maybe<Scalars['String']>;
};

export type UserFistVisitInputDto = {
  firstVisit: Scalars['Boolean'];
};

export type UserFistVisitResDto = {
  __typename?: 'UserFistVisitResDto';
  firstVisit: Scalars['Boolean'];
};

export type UserInputDto = {
  email: Scalars['String'];
  name: Scalars['String'];
  phone: Scalars['String'];
};

export type UserResDto = {
  __typename?: 'UserResDto';
  email: Scalars['String'];
  isDeactivated: Scalars['Boolean'];
  name: Scalars['String'];
  phone: Scalars['String'];
  registeredSince: Scalars['String'];
  uuid: Scalars['String'];
};

export type UserServicesDto = {
  __typename?: 'UserServicesDto';
  services: UserSimServices;
  status: Get_User_Status_Type;
  usage?: Maybe<UserDataUsageDto>;
};

export type UserSimServices = {
  __typename?: 'UserSimServices';
  data: Scalars['Boolean'];
  sms: Scalars['Boolean'];
  voice: Scalars['Boolean'];
};

export type WhoamiApiDto = {
  __typename?: 'WhoamiAPIDto';
  email: Scalars['String'];
  first_visit: Scalars['Boolean'];
  id: Scalars['String'];
  name: Scalars['String'];
  role: Scalars['String'];
};

export type WhoamiDto = {
  __typename?: 'WhoamiDto';
  email: Scalars['String'];
  id: Scalars['String'];
  isFirstVisit: Scalars['Boolean'];
  name: Scalars['String'];
  role: Scalars['String'];
};

export type WhoamiQueryVariables = Exact<{ [key: string]: never; }>;


export type WhoamiQuery = { __typename?: 'Query', whoami: { __typename?: 'WhoamiDto', id: string, name: string, email: string, role: string, isFirstVisit: boolean } };

export type GetSubscriberMetricByNetworkQueryVariables = Exact<{
  networkId: Scalars['String'];
}>;


export type GetSubscriberMetricByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };

export type GetNetworkNodesStatQueryVariables = Exact<{
  networkId: Scalars['String'];
}>;


export type GetNetworkNodesStatQuery = { __typename?: 'Query', getNetworkNodesStat: { __typename?: 'NodeStatsResponse', totalCount: number, upCount: number, claimCount: number } };

export type GetNetworkDataUsageQueryVariables = Exact<{
  networkId: Scalars['String'];
  filter: Time_Filter;
}>;


export type GetNetworkDataUsageQuery = { __typename?: 'Query', getNetworkDataUsage: { __typename?: 'DataUsageNetworkResponse', usage: number } };

export type GetPackagesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { __typename?: 'Query', getPackages: { __typename?: 'PackagesResDto', packages: Array<{ __typename?: 'PackageDto', uuid: string, name: string, orgId: string, active: boolean, duration: string, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: string, dataVolume: string, voiceVolume: string, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } }> } };

export type AddSubscriberMutationVariables = Exact<{
  data: SubscriberInputDto;
}>;


export type AddSubscriberMutation = { __typename?: 'Mutation', addSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> } };

export type GetSubscriberQueryVariables = Exact<{
  subscriberId: Scalars['String'];
}>;


export type GetSubscriberQuery = { __typename?: 'Query', getSubscriber: { __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> } };

export type UpdateSubscriberMutationVariables = Exact<{
  subscriberId: Scalars['String'];
  data: UpdateSubscriberInputDto;
}>;


export type UpdateSubscriberMutation = { __typename?: 'Mutation', updateSubscriber: { __typename?: 'BoolResponse', success: boolean } };

export type DeleteSubscriberMutationVariables = Exact<{
  subscriberId: Scalars['String'];
}>;


export type DeleteSubscriberMutation = { __typename?: 'Mutation', deleteSubscriber: { __typename?: 'BoolResponse', success: boolean } };

export type GetSubscribersByNetworkQueryVariables = Exact<{
  networkId: Scalars['String'];
}>;


export type GetSubscribersByNetworkQuery = { __typename?: 'Query', getSubscribersByNetwork: { __typename?: 'SubscribersResDto', subscribers: Array<{ __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> }> } };

export type GetSubscriberMetricsByNetworkQueryVariables = Exact<{
  networkId: Scalars['String'];
}>;


export type GetSubscriberMetricsByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };


export const WhoamiDocument = gql`
    query Whoami {
  whoami {
    id
    name
    email
    role
    isFirstVisit
  }
}
    `;

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
export const GetSubscriberMetricByNetworkDocument = gql`
    query GetSubscriberMetricByNetwork($networkId: String!) {
  getSubscriberMetricsByNetwork(networkId: $networkId) {
    total
    active
    inactive
    terminated
  }
}
    `;

/**
 * __useGetSubscriberMetricByNetworkQuery__
 *
 * To run a query within a React component, call `useGetSubscriberMetricByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscriberMetricByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscriberMetricByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSubscriberMetricByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberMetricByNetworkQuery, GetSubscriberMetricByNetworkQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberMetricByNetworkQuery, GetSubscriberMetricByNetworkQueryVariables>(GetSubscriberMetricByNetworkDocument, options);
      }
export function useGetSubscriberMetricByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberMetricByNetworkQuery, GetSubscriberMetricByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberMetricByNetworkQuery, GetSubscriberMetricByNetworkQueryVariables>(GetSubscriberMetricByNetworkDocument, options);
        }
export type GetSubscriberMetricByNetworkQueryHookResult = ReturnType<typeof useGetSubscriberMetricByNetworkQuery>;
export type GetSubscriberMetricByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscriberMetricByNetworkLazyQuery>;
export type GetSubscriberMetricByNetworkQueryResult = Apollo.QueryResult<GetSubscriberMetricByNetworkQuery, GetSubscriberMetricByNetworkQueryVariables>;
export const GetNetworkNodesStatDocument = gql`
    query GetNetworkNodesStat($networkId: String!) {
  getNetworkNodesStat(networkId: $networkId) {
    totalCount
    upCount
    claimCount
  }
}
    `;

/**
 * __useGetNetworkNodesStatQuery__
 *
 * To run a query within a React component, call `useGetNetworkNodesStatQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkNodesStatQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkNodesStatQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetNetworkNodesStatQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkNodesStatQuery, GetNetworkNodesStatQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkNodesStatQuery, GetNetworkNodesStatQueryVariables>(GetNetworkNodesStatDocument, options);
      }
export function useGetNetworkNodesStatLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkNodesStatQuery, GetNetworkNodesStatQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkNodesStatQuery, GetNetworkNodesStatQueryVariables>(GetNetworkNodesStatDocument, options);
        }
export type GetNetworkNodesStatQueryHookResult = ReturnType<typeof useGetNetworkNodesStatQuery>;
export type GetNetworkNodesStatLazyQueryHookResult = ReturnType<typeof useGetNetworkNodesStatLazyQuery>;
export type GetNetworkNodesStatQueryResult = Apollo.QueryResult<GetNetworkNodesStatQuery, GetNetworkNodesStatQueryVariables>;
export const GetNetworkDataUsageDocument = gql`
    query GetNetworkDataUsage($networkId: String!, $filter: TIME_FILTER!) {
  getNetworkDataUsage(networkId: $networkId, filter: $filter) {
    usage
  }
}
    `;

/**
 * __useGetNetworkDataUsageQuery__
 *
 * To run a query within a React component, call `useGetNetworkDataUsageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkDataUsageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkDataUsageQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetNetworkDataUsageQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkDataUsageQuery, GetNetworkDataUsageQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkDataUsageQuery, GetNetworkDataUsageQueryVariables>(GetNetworkDataUsageDocument, options);
      }
export function useGetNetworkDataUsageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkDataUsageQuery, GetNetworkDataUsageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkDataUsageQuery, GetNetworkDataUsageQueryVariables>(GetNetworkDataUsageDocument, options);
        }
export type GetNetworkDataUsageQueryHookResult = ReturnType<typeof useGetNetworkDataUsageQuery>;
export type GetNetworkDataUsageLazyQueryHookResult = ReturnType<typeof useGetNetworkDataUsageLazyQuery>;
export type GetNetworkDataUsageQueryResult = Apollo.QueryResult<GetNetworkDataUsageQuery, GetNetworkDataUsageQueryVariables>;
export const GetPackagesDocument = gql`
    query getPackages {
  getPackages {
    packages {
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
      rate {
        sms_mo
        sms_mt
        data
        amount
      }
      markup {
        baserate
        markup
      }
    }
  }
}
    `;

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
export const AddSubscriberDocument = gql`
    mutation addSubscriber($data: SubscriberInputDto!) {
  addSubscriber(data: $data) {
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
}
    `;
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
}
    `;

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
  }
}
    `;

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