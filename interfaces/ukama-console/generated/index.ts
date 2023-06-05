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
  DateTime: { input: any; output: any; }
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

export type AddMemberInputDto = {
  email: Scalars['String']['input'];
  role: Scalars['String']['input'];
};

export type AddNetworkInputDto = {
  network_name: Scalars['String']['input'];
  org: Scalars['String']['input'];
};

export type AddNodeDto = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  state: Scalars['String']['input'];
};

export type AddNodeResponse = {
  __typename?: 'AddNodeResponse';
  attached: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  state: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type AddOrgInputDto = {
  certificate: Scalars['String']['input'];
  name: Scalars['String']['input'];
  owner_uuid: Scalars['String']['input'];
};

export type AddPackageInputDto = {
  amount: Scalars['Int']['input'];
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
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
  startDate: Scalars['String']['input'];
};

export type AddSiteInputDto = {
  site: Scalars['String']['input'];
};

export type AddUserServiceRes = {
  __typename?: 'AddUserServiceRes';
  iccid: Scalars['String']['output'];
  user: OrgUserDto;
};

export type AlertDto = {
  __typename?: 'AlertDto';
  alertDate?: Maybe<Scalars['DateTime']['output']>;
  description?: Maybe<Scalars['String']['output']>;
  id?: Maybe<Scalars['String']['output']>;
  title?: Maybe<Scalars['String']['output']>;
  type: Alert_Type;
};

export type AlertResponse = {
  __typename?: 'AlertResponse';
  data: Array<AlertDto>;
  length: Scalars['Float']['output'];
  status: Scalars['String']['output'];
};

export type AlertsResponse = {
  __typename?: 'AlertsResponse';
  alerts: Array<AlertDto>;
  meta: Meta;
};

export type AllocateSimInputDto = {
  iccid: Scalars['String']['input'];
  networkId: Scalars['String']['input'];
  packageId: Scalars['String']['input'];
  simType: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
};

export type ApiMethodDataDto = {
  __typename?: 'ApiMethodDataDto';
  body?: Maybe<Scalars['String']['output']>;
  headers?: Maybe<Scalars['String']['output']>;
  params?: Maybe<Scalars['String']['output']>;
  path: Scalars['String']['output'];
  type: Api_Method_Type;
};

export type AttachedNodes = {
  __typename?: 'AttachedNodes';
  nodeId: Scalars['String']['output'];
};

export type AuthType = {
  __typename?: 'AuthType';
  Authorization: Scalars['String']['output'];
  Cookie: Scalars['String']['output'];
};

export type BillHistoryDto = {
  __typename?: 'BillHistoryDto';
  date: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  subtotal: Scalars['Float']['output'];
  totalUsage: Scalars['Float']['output'];
};

export type BillHistoryResponse = {
  __typename?: 'BillHistoryResponse';
  data: Array<BillHistoryDto>;
  status: Scalars['String']['output'];
};

export type BillResponse = {
  __typename?: 'BillResponse';
  bill: Array<CurrentBillDto>;
  billMonth: Scalars['String']['output'];
  dueDate: Scalars['String']['output'];
  total: Scalars['Float']['output'];
};

export type BoolResponse = {
  __typename?: 'BoolResponse';
  success: Scalars['Boolean']['output'];
};

export type ConnectedUserDto = {
  __typename?: 'ConnectedUserDto';
  totalUser: Scalars['String']['output'];
};

export type CreateCustomerDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type CurrentBillDto = {
  __typename?: 'CurrentBillDto';
  dataUsed: Scalars['Float']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  rate: Scalars['Float']['output'];
  subtotal: Scalars['Float']['output'];
};

export type CurrentBillResponse = {
  __typename?: 'CurrentBillResponse';
  data: Array<CurrentBillDto>;
  status: Scalars['String']['output'];
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
  billDue: Scalars['Float']['output'];
  dataBill: Scalars['Float']['output'];
  id: Scalars['String']['output'];
};

export type DataBillResponse = {
  __typename?: 'DataBillResponse';
  data: DataBillDto;
  status: Scalars['String']['output'];
};

export type DataUsageDto = {
  __typename?: 'DataUsageDto';
  dataConsumed: Scalars['Float']['output'];
  dataPackage: Scalars['String']['output'];
  id: Scalars['String']['output'];
};

export type DataUsageInputDto = {
  ids: Array<Scalars['String']['input']>;
};

export type DataUsageNetworkResponse = {
  __typename?: 'DataUsageNetworkResponse';
  usage: Scalars['Float']['output'];
};

export type DataUsageResponse = {
  __typename?: 'DataUsageResponse';
  data: DataUsageDto;
  status: Scalars['String']['output'];
};

export type DeactivateResponse = {
  __typename?: 'DeactivateResponse';
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type DefaultMarkupApiResDto = {
  __typename?: 'DefaultMarkupAPIResDto';
  markup: Scalars['Float']['output'];
};

export type DefaultMarkupHistoryApiResDto = {
  __typename?: 'DefaultMarkupHistoryAPIResDto';
  markupRates?: Maybe<Array<DefaultMarkupHistoryDto>>;
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

export type DeleteNodeRes = {
  __typename?: 'DeleteNodeRes';
  nodeId: Scalars['String']['output'];
};

export type DeleteSimInputDto = {
  simId: Scalars['String']['input'];
};

export type DeleteSimResDto = {
  __typename?: 'DeleteSimResDto';
  simId?: Maybe<Scalars['String']['output']>;
};

export type ESimQrCodeRes = {
  __typename?: 'ESimQRCodeRes';
  qrCode: Scalars['String']['output'];
};

export type ErrorType = {
  __typename?: 'ErrorType';
  code: Scalars['Float']['output'];
  description?: Maybe<Scalars['String']['output']>;
  message: Scalars['String']['output'];
};

export type EsimDto = {
  __typename?: 'EsimDto';
  active: Scalars['Boolean']['output'];
  esim: Scalars['String']['output'];
};

export type EsimResponse = {
  __typename?: 'EsimResponse';
  data: Array<EsimDto>;
  status: Scalars['String']['output'];
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
  email: Scalars['String']['output'];
  isFirstVisit: Scalars['Boolean']['output'];
};

export type GetESimQrCodeInput = {
  simId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type GetMetricsRes = {
  __typename?: 'GetMetricsRes';
  metrics: Array<MetricRes>;
  next: Scalars['Boolean']['output'];
  to: Scalars['Float']['output'];
};

export type GetNodeStatusInput = {
  nodeId: Scalars['String']['input'];
  nodeType: Node_Type;
};

export type GetNodeStatusRes = {
  __typename?: 'GetNodeStatusRes';
  status: Org_Node_State;
  uptime: Scalars['Float']['output'];
};

export type GetPackagesForSimInputDto = {
  simId: Scalars['String']['input'];
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
  networkId: Scalars['String']['input'];
};

export type GetSimBySubscriberIdInputDto = {
  subscriberId: Scalars['String']['input'];
};

export type GetSimInputDto = {
  simId: Scalars['String']['input'];
};

export type GetUserDto = {
  __typename?: 'GetUserDto';
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  phone: Scalars['String']['output'];
};

export type GetUserResponseDto = {
  __typename?: 'GetUserResponseDto';
  data: Array<GetUserDto>;
  length: Scalars['Float']['output'];
  status: Scalars['String']['output'];
};

export type GetUsersDto = {
  __typename?: 'GetUsersDto';
  dataPlan?: Maybe<Scalars['String']['output']>;
  dataUsage?: Maybe<Scalars['String']['output']>;
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type IdResponse = {
  __typename?: 'IdResponse';
  uuid: Scalars['String']['output'];
};

export type LinkNodes = {
  __typename?: 'LinkNodes';
  attachedNodeIds: Array<Scalars['String']['output']>;
  nodeId: Scalars['String']['output'];
};

export type MemberApiObj = {
  __typename?: 'MemberAPIObj';
  is_deactivated: Scalars['Boolean']['output'];
  member_since: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
  role: Scalars['String']['output'];
  user_id: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
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

export type Meta = {
  __typename?: 'Meta';
  count: Scalars['Float']['output'];
  page: Scalars['Float']['output'];
  pages: Scalars['Float']['output'];
  size: Scalars['Float']['output'];
};

export type MetricDto = {
  __typename?: 'MetricDto';
  x: Scalars['Float']['output'];
  y: Scalars['Float']['output'];
};

export type MetricInfo = {
  __typename?: 'MetricInfo';
  org: Scalars['String']['output'];
};

export type MetricLatestValueRes = {
  __typename?: 'MetricLatestValueRes';
  time: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type MetricRes = {
  __typename?: 'MetricRes';
  data: Array<MetricDto>;
  name: Scalars['String']['output'];
  next: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
};

export type MetricServiceValueRes = {
  __typename?: 'MetricServiceValueRes';
  metric: MetricInfo;
  value: Array<MetricValues>;
};

export type MetricValues = {
  __typename?: 'MetricValues';
  x: Scalars['Float']['output'];
  y: Scalars['String']['output'];
};

export type MetricsByTabInputDto = {
  from: Scalars['Float']['input'];
  nodeId: Scalars['String']['input'];
  nodeType: Node_Type;
  regPolling: Scalars['Boolean']['input'];
  step: Scalars['Float']['input'];
  tab: Graphs_Tab;
  to: Scalars['Float']['input'];
};

export type MetricsInputDto = {
  from: Scalars['Float']['input'];
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  regPolling: Scalars['Boolean']['input'];
  step: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
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
  allocateSim: SimDto;
  attachPaymentWithCustomer: Scalars['Boolean']['output'];
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
  uploadSims: UploadSimsResDto;
};


export type MutationAddMemberArgs = {
  data: AddMemberInputDto;
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
  paymentId: Scalars['String']['input'];
};


export type MutationCreateCustomerArgs = {
  data: CreateCustomerDto;
};


export type MutationDeactivateUserArgs = {
  uuid: Scalars['String']['input'];
};


export type MutationDefaultMarkupArgs = {
  data: DefaultMarkupInputDto;
};


export type MutationDeleteNodeArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeletePackageArgs = {
  packageId: Scalars['String']['input'];
};


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type MutationDeleteUserArgs = {
  userId: Scalars['String']['input'];
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
  memberId: Scalars['String']['input'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeDto;
};


export type MutationUpdatePackageArgs = {
  data: UpdatePackageInputDto;
  packageId: Scalars['String']['input'];
};


export type MutationUpdateSubscriberArgs = {
  data: UpdateSubscriberInputDto;
  subscriberId: Scalars['String']['input'];
};


export type MutationUpdateUserArgs = {
  data: UpdateUserInputDto;
  userId: Scalars['String']['input'];
};


export type MutationUpdateUserRoamingArgs = {
  data: UpdateUserServiceInput;
};


export type MutationUpdateUserStatusArgs = {
  data: UpdateUserServiceInput;
};


export type MutationUploadSimsArgs = {
  data: UploadSimsInputDto;
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
  created_at: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_deactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
};

export type NetworkApiResDto = {
  __typename?: 'NetworkAPIResDto';
  network: NetworkApiDto;
};

export type NetworkDto = {
  __typename?: 'NetworkDto';
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
};

export type NetworkStatusDto = {
  __typename?: 'NetworkStatusDto';
  liveNode: Scalars['Float']['output'];
  status: Network_Status;
  totalNodes: Scalars['Float']['output'];
};

export type NetworkStatusResponse = {
  __typename?: 'NetworkStatusResponse';
  data: NetworkStatusDto;
  status: Scalars['String']['output'];
};

export type NetworksApiResDto = {
  __typename?: 'NetworksAPIResDto';
  networks: Array<NetworkApiDto>;
  org_id: Scalars['String']['output'];
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
  orgId: Scalars['String']['output'];
};

export type NodeAppResponse = {
  __typename?: 'NodeAppResponse';
  cpu: Scalars['String']['output'];
  id: Scalars['String']['output'];
  memory: Scalars['String']['output'];
  title: Scalars['String']['output'];
  version: Scalars['String']['output'];
};

export type NodeAppsVersionLogsResponse = {
  __typename?: 'NodeAppsVersionLogsResponse';
  date: Scalars['Float']['output'];
  notes: Scalars['String']['output'];
  version: Scalars['String']['output'];
};

export type NodeDto = {
  __typename?: 'NodeDto';
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isUpdateAvailable: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  status: Org_Node_State;
  totalUser: Scalars['Float']['output'];
  type: Scalars['String']['output'];
  updateDescription: Scalars['String']['output'];
  updateShortNote: Scalars['String']['output'];
  updateVersion: Scalars['String']['output'];
};

export type NodeObj = {
  attached?: InputMaybe<Array<NodeObj>>;
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  state: Scalars['String']['input'];
};

export type NodeResponse = {
  __typename?: 'NodeResponse';
  attached: Array<OrgNodeDto>;
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  state: Org_Node_State;
  type: Node_Type;
};

export type NodeStatsResponse = {
  __typename?: 'NodeStatsResponse';
  claimCount: Scalars['Float']['output'];
  totalCount: Scalars['Float']['output'];
  upCount: Scalars['Float']['output'];
};

export enum Org_Node_State {
  Error = 'ERROR',
  Onboarded = 'ONBOARDED',
  Pending = 'PENDING',
  Undefined = 'UNDEFINED'
}

export type OrgApiDto = {
  __typename?: 'OrgAPIDto';
  certificate: Scalars['Boolean']['output'];
  created_at: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_deactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgApiResDto = {
  __typename?: 'OrgAPIResDto';
  org: OrgApiDto;
};

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['Boolean']['output'];
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgMemberApiResDto = {
  __typename?: 'OrgMemberAPIResDto';
  member: MemberApiObj;
};

export type OrgMembersApiResDto = {
  __typename?: 'OrgMembersAPIResDto';
  members: Array<MemberApiObj>;
  org: Scalars['String']['output'];
};

export type OrgMembersResDto = {
  __typename?: 'OrgMembersResDto';
  members: Array<MemberObj>;
  org: Scalars['String']['output'];
};

export type OrgMetricValueDto = {
  __typename?: 'OrgMetricValueDto';
  x: Scalars['Float']['output'];
  y: Scalars['String']['output'];
};

export type OrgNodeDto = {
  __typename?: 'OrgNodeDto';
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  state: Org_Node_State;
  type: Node_Type;
};

export type OrgNodeResponse = {
  __typename?: 'OrgNodeResponse';
  nodes: Array<OrgNodeDto>;
  orgName: Scalars['String']['output'];
};

export type OrgNodeResponseDto = {
  __typename?: 'OrgNodeResponseDto';
  activeNodes: Scalars['Float']['output'];
  nodes: Array<NodeDto>;
  orgId: Scalars['String']['output'];
  totalNodes: Scalars['Float']['output'];
};

export type OrgResDto = {
  __typename?: 'OrgResDto';
  org: OrgDto;
};

export type OrgUserDto = {
  __typename?: 'OrgUserDto';
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type OrgUserResponse = {
  __typename?: 'OrgUserResponse';
  sim: OrgUserSimDto;
  user: OrgUserDto;
};

export type OrgUserSimDto = {
  __typename?: 'OrgUserSimDto';
  carrier: UserServicesDto;
  iccid: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  ukama: UserServicesDto;
};

export type OrgUsersResponse = {
  __typename?: 'OrgUsersResponse';
  org: Scalars['String']['output'];
  users: Array<OrgUserDto>;
};

export type OrgsApiResDto = {
  __typename?: 'OrgsAPIResDto';
  orgs: Array<OrgApiDto>;
  owner: Scalars['String']['output'];
};

export type OrgsResDto = {
  __typename?: 'OrgsResDto';
  orgs: Array<OrgDto>;
  owner: Scalars['String']['output'];
};

export type PackageApiDto = {
  __typename?: 'PackageAPIDto';
  active: Scalars['Boolean']['output'];
  amount: Scalars['Float']['output'];
  apn: Scalars['String']['output'];
  country: Scalars['String']['output'];
  created_at: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  data_unit: Scalars['String']['output'];
  data_volume: Scalars['String']['output'];
  deleted_at: Scalars['String']['output'];
  dlbr: Scalars['String']['output'];
  duration: Scalars['String']['output'];
  flatrate: Scalars['Boolean']['output'];
  from: Scalars['String']['output'];
  markup: PackageMarkupApiDto;
  message_unit: Scalars['String']['output'];
  name: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
  owner_id: Scalars['String']['output'];
  provider: Scalars['String']['output'];
  rate: PackageRateApiDto;
  sim_type: Scalars['String']['output'];
  sms_volume: Scalars['String']['output'];
  to: Scalars['String']['output'];
  type: Scalars['String']['output'];
  ulbr: Scalars['String']['output'];
  updated_at: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
  voice_unit: Scalars['String']['output'];
  voice_volume: Scalars['String']['output'];
};

export type PackageApiResDto = {
  __typename?: 'PackageAPIResDto';
  package: PackageApiDto;
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

export type PackageMarkupDto = {
  __typename?: 'PackageMarkupDto';
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

export type PackageRateDto = {
  __typename?: 'PackageRateDto';
  amount: Scalars['Float']['output'];
  data: Scalars['Float']['output'];
  smsMo: Scalars['String']['output'];
  smsMt: Scalars['Float']['output'];
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
  pageNo: Scalars['Float']['input'];
  pageSize: Scalars['Float']['input'];
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
  getSim: SimDto;
  getSimPoolStats: SimPoolStatsDto;
  getSims: SimsResDto;
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
  networkId: Scalars['String']['input'];
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
  simId: Scalars['String']['input'];
};


export type QueryGetEsimQrArgs = {
  data: GetESimQrCodeInput;
};


export type QueryGetMetricsByTabArgs = {
  data: MetricsByTabInputDto;
};


export type QueryGetNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNetworkDataUsageArgs = {
  filter: Time_Filter;
  networkId: Scalars['String']['input'];
};


export type QueryGetNetworkNodesStatArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  nodeId: Scalars['String']['input'];
};


export type QueryGetNodeStatusArgs = {
  data: GetNodeStatusInput;
};


export type QueryGetOrgArgs = {
  orgName: Scalars['String']['input'];
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


export type QueryGetUsersDataUsageArgs = {
  data: DataUsageInputDto;
};

export type RemovePackageFormSimInputDto = {
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type RemovePackageFromSimResDto = {
  __typename?: 'RemovePackageFromSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type SetActivePackageForSimInputDto = {
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type SetActivePackageForSimResDto = {
  __typename?: 'SetActivePackageForSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type SimApiDto = {
  __typename?: 'SimAPIDto';
  activation_code: Scalars['String']['output'];
  created_at: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_allocated: Scalars['String']['output'];
  is_physical: Scalars['String']['output'];
  msisdn: Scalars['String']['output'];
  qr_code: Scalars['String']['output'];
  sim_type: Scalars['String']['output'];
  sm_ap_address: Scalars['String']['output'];
};

export type SimApiResDto = {
  __typename?: 'SimAPIResDto';
  sim: SimApiDto;
};

export type SimDataUsage = {
  __typename?: 'SimDataUsage';
  usage: Scalars['String']['output'];
};

export type SimDetailsDto = {
  __typename?: 'SimDetailsDto';
  Package: SimPackageDto;
  activationsCount: Scalars['Float']['output'];
  allocatedAt: Scalars['String']['output'];
  deactivationsCount: Scalars['Float']['output'];
  firstActivatedOn: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  lastActivatedOn: Scalars['String']['output'];
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  status: Scalars['String']['output'];
  subscriberId?: Maybe<Scalars['String']['output']>;
  type: Scalars['String']['output'];
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

export type SimPackageDto = {
  __typename?: 'SimPackageDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
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

export type SimsApiResDto = {
  __typename?: 'SimsAPIResDto';
  sims: Array<SimApiDto>;
};

export type SimsResDto = {
  __typename?: 'SimsResDto';
  sim: Array<SimDto>;
};

export type SiteApiDto = {
  __typename?: 'SiteAPIDto';
  created_at: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_deactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
};

export type SiteApiResDto = {
  __typename?: 'SiteAPIResDto';
  site: SiteApiDto;
};

export type SiteDto = {
  __typename?: 'SiteDto';
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
};

export type SitesApiResDto = {
  __typename?: 'SitesAPIResDto';
  network_id: Scalars['String']['output'];
  sites: Array<SiteApiDto>;
};

export type SitesResDto = {
  __typename?: 'SitesResDto';
  networkId: Scalars['String']['output'];
  sites: Array<SiteDto>;
};

export type StripeCustomer = {
  __typename?: 'StripeCustomer';
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type StripePaymentMethods = {
  __typename?: 'StripePaymentMethods';
  brand: Scalars['String']['output'];
  country?: Maybe<Scalars['String']['output']>;
  created: Scalars['Float']['output'];
  cvc_check?: Maybe<Scalars['String']['output']>;
  exp_month: Scalars['Float']['output'];
  exp_year: Scalars['Float']['output'];
  funding: Scalars['String']['output'];
  id: Scalars['String']['output'];
  last4: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type SubSimApiDto = {
  __typename?: 'SubSimAPIDto';
  activations_count: Scalars['String']['output'];
  allocated_at: Scalars['String']['output'];
  deactivations_count: Scalars['String']['output'];
  first_activated_on?: Maybe<Scalars['String']['output']>;
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  is_physical?: Maybe<Scalars['Boolean']['output']>;
  last_activated_on?: Maybe<Scalars['String']['output']>;
  msisdn: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
  package?: Maybe<Scalars['String']['output']>;
  status: Scalars['String']['output'];
  subscriber_id: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type SubscriberApiDto = {
  __typename?: 'SubscriberAPIDto';
  address: Scalars['String']['output'];
  date_of_birth: Scalars['String']['output'];
  email: Scalars['String']['output'];
  first_name: Scalars['String']['output'];
  gender: Scalars['String']['output'];
  id_serial: Scalars['String']['output'];
  last_name: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
  org_id: Scalars['String']['output'];
  phone_number: Scalars['String']['output'];
  proof_of_identification: Scalars['String']['output'];
  sim: Array<SubSimApiDto>;
  subscriber_id: Scalars['String']['output'];
};

export type SubscriberApiResDto = {
  __typename?: 'SubscriberAPIResDto';
  subscriber: SubscriberApiDto;
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
  orgId: Scalars['String']['output'];
  orgName: Scalars['String']['output'];
  userId: Scalars['String']['output'];
};

export enum Time_Filter {
  Month = 'MONTH',
  Today = 'TODAY',
  Total = 'TOTAL',
  Week = 'WEEK'
}

export type ToggleSimStatusInputDto = {
  simId: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type UpdateMemberInputDto = {
  isDeactivated: Scalars['Boolean']['input'];
};

export type UpdateNodeDto = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
};

export type UpdateNodeResponse = {
  __typename?: 'UpdateNodeResponse';
  attached: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  state: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type UpdatePackageInputDto = {
  active?: InputMaybe<Scalars['Boolean']['input']>;
  data_volume?: InputMaybe<Scalars['Int']['input']>;
  duration?: InputMaybe<Scalars['Int']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  org_rates_id?: InputMaybe<Scalars['Int']['input']>;
  sim_type?: InputMaybe<Scalars['String']['input']>;
  sms_volume?: InputMaybe<Scalars['Int']['input']>;
  voice_volume?: InputMaybe<Scalars['Int']['input']>;
};

export type UpdateSubscriberInputDto = {
  address?: InputMaybe<Scalars['String']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  id_serial?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  proof_of_identification?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateUserInputDto = {
  email?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  phone: Scalars['String']['input'];
};

export type UpdateUserServiceInput = {
  simId: Scalars['String']['input'];
  status: Scalars['Boolean']['input'];
  userId: Scalars['String']['input'];
};

export type UploadSimsInputDto = {
  data: Scalars['String']['input'];
  simType: Scalars['String']['input'];
};

export type UploadSimsResDto = {
  __typename?: 'UploadSimsResDto';
  iccid: Array<Scalars['String']['output']>;
};

export type UserApiObj = {
  __typename?: 'UserAPIObj';
  email: Scalars['String']['output'];
  is_deactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  registered_since: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type UserApiResDto = {
  __typename?: 'UserAPIResDto';
  user: Array<UserApiObj>;
};

export type UserDataUsageDto = {
  __typename?: 'UserDataUsageDto';
  dataAllowanceBytes?: Maybe<Scalars['String']['output']>;
  dataUsedBytes?: Maybe<Scalars['String']['output']>;
};

export type UserFistVisitInputDto = {
  firstVisit: Scalars['Boolean']['input'];
};

export type UserFistVisitResDto = {
  __typename?: 'UserFistVisitResDto';
  firstVisit: Scalars['Boolean']['output'];
};

export type UserInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  phone: Scalars['String']['input'];
};

export type UserResDto = {
  __typename?: 'UserResDto';
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  registeredSince: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type UserServicesDto = {
  __typename?: 'UserServicesDto';
  services: UserSimServices;
  status: Get_User_Status_Type;
  usage?: Maybe<UserDataUsageDto>;
};

export type UserSimServices = {
  __typename?: 'UserSimServices';
  data: Scalars['Boolean']['output'];
  sms: Scalars['Boolean']['output'];
  voice: Scalars['Boolean']['output'];
};

export type WhoamiApiDto = {
  __typename?: 'WhoamiAPIDto';
  email: Scalars['String']['output'];
  first_visit: Scalars['Boolean']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
};

export type WhoamiDto = {
  __typename?: 'WhoamiDto';
  email: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isFirstVisit: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
};

export type WhoamiQueryVariables = Exact<{ [key: string]: never; }>;


export type WhoamiQuery = { __typename?: 'Query', whoami: { __typename?: 'WhoamiDto', id: string, name: string, email: string, role: string, isFirstVisit: boolean } };

export type GetSubscriberMetricByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscriberMetricByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };

export type GetNetworkNodesStatQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetNetworkNodesStatQuery = { __typename?: 'Query', getNetworkNodesStat: { __typename?: 'NodeStatsResponse', totalCount: number, upCount: number, claimCount: number } };

export type GetNetworkDataUsageQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
  filter: Time_Filter;
}>;


export type GetNetworkDataUsageQuery = { __typename?: 'Query', getNetworkDataUsage: { __typename?: 'DataUsageNetworkResponse', usage: number } };

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


export type UpdateSubscriberMutation = { __typename?: 'Mutation', updateSubscriber: { __typename?: 'BoolResponse', success: boolean } };

export type DeleteSubscriberMutationVariables = Exact<{
  subscriberId: Scalars['String']['input'];
}>;


export type DeleteSubscriberMutation = { __typename?: 'Mutation', deleteSubscriber: { __typename?: 'BoolResponse', success: boolean } };

export type GetSubscribersByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscribersByNetworkQuery = { __typename?: 'Query', getSubscribersByNetwork: { __typename?: 'SubscribersResDto', subscribers: Array<{ __typename?: 'SubscriberDto', uuid: string, address: string, dob: string, email: string, firstName: string, lastName: string, gender: string, idSerial: string, networkId: string, orgId: string, phone: string, proofOfIdentification: string, sim: Array<{ __typename?: 'SubscriberSimDto', id: string, subscriberId: string, networkId: string, orgId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, firstActivatedOn?: string | null, lastActivatedOn?: string | null, activationsCount: string, deactivationsCount: string, allocatedAt: string, isPhysical?: boolean | null, package?: string | null }> }> } };

export type GetSubscriberMetricsByNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetSubscriberMetricsByNetworkQuery = { __typename?: 'Query', getSubscriberMetricsByNetwork: { __typename?: 'SubscriberMetricsByNetworkDto', total: number, active: number, inactive: number, terminated: number } };

export type SimPoolFragment = { __typename?: 'SimDto', activationCode: string, createdAt: string, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string };

export type GetSimsQueryVariables = Exact<{
  type: Scalars['String']['input'];
}>;


export type GetSimsQuery = { __typename?: 'Query', getSims: { __typename?: 'SimsResDto', sim: Array<{ __typename?: 'SimDto', activationCode: string, createdAt: string, iccid: string, id: string, isAllocated: string, isPhysical: string, msisdn: string, qrCode: string, simType: string, smapAddress: string }> } };

export type NetworkFragment = { __typename?: 'NetworkDto', id: string, name: string, orgId: string, isDeactivated: string, createdAt: string };

export type GetNetworkQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
}>;


export type GetNetworkQuery = { __typename?: 'Query', getNetwork: { __typename?: 'NetworkDto', id: string, name: string, orgId: string, isDeactivated: string, createdAt: string } };

export type GetNetworksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { __typename?: 'Query', getNetworks: { __typename?: 'NetworksResDto', orgId: string, networks: Array<{ __typename?: 'NetworkDto', id: string, name: string, orgId: string, isDeactivated: string, createdAt: string }> } };

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

export type GetSimpoolStatsQueryVariables = Exact<{
  type: Scalars['String']['input'];
}>;


export type GetSimpoolStatsQuery = { __typename?: 'Query', getSimPoolStats: { __typename?: 'SimPoolStatsDto', total: number, available: number, consumed: number, failed: number, physical: number, esim: number } };

export type UploadSimsMutationVariables = Exact<{
  data: UploadSimsInputDto;
}>;


export type UploadSimsMutation = { __typename?: 'Mutation', uploadSims: { __typename?: 'UploadSimsResDto', iccid: Array<string> } };

export type OrgUserFragment = { __typename?: 'UserResDto', name: string, email: string, uuid: string, phone: string, isDeactivated: boolean, registeredSince: string };

export type MemberFragment = { __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null };

export type GetOrgMemberQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgMemberQuery = { __typename?: 'Query', getOrgMembers: { __typename?: 'OrgMembersResDto', org: string, members: Array<{ __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null, user: { __typename?: 'UserResDto', name: string, email: string, uuid: string, phone: string, isDeactivated: boolean, registeredSince: string } }> } };

export type AddMemberMutationVariables = Exact<{
  data: AddMemberInputDto;
}>;


export type AddMemberMutation = { __typename?: 'Mutation', addMember: { __typename?: 'MemberObj', uuid: string, userId: string, orgId: string, role: string, isDeactivated: boolean, memberSince?: string | null } };

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
export const NetworkFragmentDoc = gql`
    fragment Network on NetworkDto {
  id
  name
  orgId
  isDeactivated
  createdAt
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
export const GetNetworkDocument = gql`
    query getNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    ...Network
  }
}
    ${NetworkFragmentDoc}`;

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
export function useGetNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
      }
export function useGetNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
export type GetNetworkQueryHookResult = ReturnType<typeof useGetNetworkQuery>;
export type GetNetworkLazyQueryHookResult = ReturnType<typeof useGetNetworkLazyQuery>;
export type GetNetworkQueryResult = Apollo.QueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
export const GetNetworksDocument = gql`
    query getNetworks {
  getNetworks {
    orgId
    networks {
      ...Network
    }
  }
}
    ${NetworkFragmentDoc}`;

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