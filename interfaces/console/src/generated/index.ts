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
  /** The javascript `Date` as string. Type represents date and time as the ISO Date string. */
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

export type ActivateUserResponse = {
  __typename?: 'ActivateUserResponse';
  success: Scalars['Boolean'];
};

export type AddNodeDto = {
  associate: Scalars['Boolean'];
  attached: Array<NodeObj>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
};

export type AddNodeResponse = {
  __typename?: 'AddNodeResponse';
  success: Scalars['Boolean'];
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

export type ConnectedUserDto = {
  __typename?: 'ConnectedUserDto';
  totalUser: Scalars['String'];
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

export type DeleteNodeRes = {
  __typename?: 'DeleteNodeRes';
  nodeId: Scalars['String'];
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

export type GetUserDto = {
  __typename?: 'GetUserDto';
  dataPlan: Scalars['String'];
  dataUsage: Scalars['String'];
  eSimNumber: Scalars['String'];
  email: Scalars['String'];
  iccid: Scalars['String'];
  id: Scalars['String'];
  name: Scalars['String'];
  roaming: Scalars['Boolean'];
  status: Scalars['Boolean'];
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

export type HeaderType = {
  __typename?: 'HeaderType';
  Authorization: Scalars['String'];
  Cookie: Scalars['String'];
};

export type LinkNodes = {
  __typename?: 'LinkNodes';
  attached: Array<AttachedNodes>;
  nodeId: Scalars['String'];
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
  addNode: AddNodeResponse;
  addUser: UserResDto;
  deactivateUser: DeactivateResponse;
  deleteNode: DeleteNodeRes;
  deleteUser: ActivateUserResponse;
  updateNode: OrgNodeDto;
  updateUser: UserResDto;
  updateUserRoaming: OrgUserSimDto;
  updateUserStatus: OrgUserSimDto;
};


export type MutationAddNodeArgs = {
  data: AddNodeDto;
};


export type MutationAddUserArgs = {
  data: UserInputDto;
};


export type MutationDeactivateUserArgs = {
  id: Scalars['String'];
};


export type MutationDeleteNodeArgs = {
  id: Scalars['String'];
};


export type MutationDeleteUserArgs = {
  userId: Scalars['String'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeDto;
};


export type MutationUpdateUserArgs = {
  data: UserInputDto;
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

export type NetworkDto = {
  __typename?: 'NetworkDto';
  liveNode: Scalars['Float'];
  status: Network_Status;
  totalNodes: Scalars['Float'];
};

export type NetworkResponse = {
  __typename?: 'NetworkResponse';
  data: NetworkDto;
  status: Scalars['String'];
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
  name: Scalars['String'];
  nodeId: Scalars['String'];
};

export type NodeResponse = {
  __typename?: 'NodeResponse';
  attached: Array<OrgNodeDto>;
  name: Scalars['String'];
  nodeId: Scalars['String'];
  state: Org_Node_State;
  type: Node_Type;
};

export enum Org_Node_State {
  Error = 'ERROR',
  Onboarded = 'ONBOARDED',
  Pending = 'PENDING',
  Undefined = 'UNDEFINED'
}

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

export type PaginationDto = {
  pageNo: Scalars['Float'];
  pageSize: Scalars['Float'];
};

export type PaginationResponse = {
  __typename?: 'PaginationResponse';
  meta: Meta;
};

export type ParsedCookie = {
  __typename?: 'ParsedCookie';
  header: HeaderType;
  orgId: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  getAlerts: AlertsResponse;
  getBillHistory: Array<BillHistoryDto>;
  getConnectedUsers: ConnectedUserDto;
  getCurrentBill: BillResponse;
  getDataBill: DataBillDto;
  getDataUsage: DataUsageDto;
  getEsimQR: ESimQrCodeRes;
  getEsims: Array<EsimDto>;
  getMetricsByTab: GetMetricsRes;
  getNetworkStatus: NetworkDto;
  getNode: NodeResponse;
  getNodeApps: Array<NodeAppResponse>;
  getNodeAppsVersionLogs: Array<NodeAppsVersionLogsResponse>;
  getNodeStatus: GetNodeStatusRes;
  getNodesByOrg: OrgNodeResponseDto;
  getUser: GetUserDto;
  getUsersByOrg: Array<GetUsersDto>;
  getUsersDataUsage: Array<GetUserDto>;
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
  filter: Time_Filter;
};


export type QueryGetEsimQrArgs = {
  data: GetESimQrCodeInput;
};


export type QueryGetMetricsByTabArgs = {
  data: MetricsByTabInputDto;
};


export type QueryGetNodeArgs = {
  nodeId: Scalars['String'];
};


export type QueryGetNodeStatusArgs = {
  data: GetNodeStatusInput;
};


export type QueryGetUserArgs = {
  userId: Scalars['String'];
};


export type QueryGetUsersDataUsageArgs = {
  data: DataUsageInputDto;
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

export enum Time_Filter {
  Month = 'MONTH',
  Today = 'TODAY',
  Total = 'TOTAL',
  Week = 'WEEK'
}

export type UpdateNodeDto = {
  name: Scalars['String'];
  nodeId: Scalars['String'];
};

export type UpdateUserServiceInput = {
  simId: Scalars['String'];
  status: Scalars['Boolean'];
  userId: Scalars['String'];
};

export type UserDataUsageDto = {
  __typename?: 'UserDataUsageDto';
  dataAllowanceBytes?: Maybe<Scalars['String']>;
  dataUsedBytes?: Maybe<Scalars['String']>;
};

export type UserInputDto = {
  email: Scalars['String'];
  name: Scalars['String'];
  status?: InputMaybe<Scalars['Boolean']>;
};

export type UserResDto = {
  __typename?: 'UserResDto';
  email: Scalars['String'];
  iccid?: Maybe<Scalars['String']>;
  id: Scalars['String'];
  name: Scalars['String'];
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

export type GetDataUsageQueryVariables = Exact<{
  filter: Time_Filter;
}>;


export type GetDataUsageQuery = { __typename?: 'Query', getDataUsage: { __typename?: 'DataUsageDto', id: string, dataConsumed: number, dataPackage: string } };

export type GetLatestDataUsageSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetLatestDataUsageSubscription = { __typename?: 'Subscription', getDataUsage: { __typename?: 'DataUsageDto', id: string, dataConsumed: number, dataPackage: string } };

export type GetConnectedUsersQueryVariables = Exact<{
  filter: Time_Filter;
}>;


export type GetConnectedUsersQuery = { __typename?: 'Query', getConnectedUsers: { __typename?: 'ConnectedUserDto', totalUser: string } };

export type GetLatestConnectedUsersSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetLatestConnectedUsersSubscription = { __typename?: 'Subscription', getConnectedUsers: { __typename?: 'ConnectedUserDto', totalUser: string } };

export type GetLatestDataBillSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetLatestDataBillSubscription = { __typename?: 'Subscription', getDataBill: { __typename?: 'DataBillDto', id: string, dataBill: number, billDue: number } };

export type GetBillHistoryQueryVariables = Exact<{ [key: string]: never; }>;


export type GetBillHistoryQuery = { __typename?: 'Query', getBillHistory: Array<{ __typename?: 'BillHistoryDto', id: string, date: string, description: string, totalUsage: number, subtotal: number }> };

export type GetCurrentBillQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCurrentBillQuery = { __typename?: 'Query', getCurrentBill: { __typename?: 'BillResponse', total: number, billMonth: string, dueDate: string, bill: Array<{ __typename?: 'CurrentBillDto', id: string, name: string, dataUsed: number, rate: number, subtotal: number }> } };

export type GetAlertsQueryVariables = Exact<{
  data: PaginationDto;
}>;


export type GetAlertsQuery = { __typename?: 'Query', getAlerts: { __typename?: 'AlertsResponse', meta: { __typename?: 'Meta', count: number, page: number, size: number, pages: number }, alerts: Array<{ __typename?: 'AlertDto', id?: string | null, type: Alert_Type, title?: string | null, description?: string | null, alertDate?: any | null }> } };

export type GetEsimQrQueryVariables = Exact<{
  data: GetESimQrCodeInput;
}>;


export type GetEsimQrQuery = { __typename?: 'Query', getEsimQR: { __typename?: 'ESimQRCodeRes', qrCode: string } };

export type GetLatestAlertsSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetLatestAlertsSubscription = { __typename?: 'Subscription', getAlerts: { __typename?: 'AlertDto', id?: string | null, type: Alert_Type, title?: string | null, description?: string | null, alertDate?: any | null } };

export type GetNodesByOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNodesByOrgQuery = { __typename?: 'Query', getNodesByOrg: { __typename?: 'OrgNodeResponseDto', orgId: string, activeNodes: number, totalNodes: number, nodes: Array<{ __typename?: 'NodeDto', id: string, status: Org_Node_State, name: string, type: string, description: string, totalUser: number, isUpdateAvailable: boolean, updateVersion: string, updateShortNote: string, updateDescription: string }> } };

export type GetNodeAppsVersionLogsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNodeAppsVersionLogsQuery = { __typename?: 'Query', getNodeAppsVersionLogs: Array<{ __typename?: 'NodeAppsVersionLogsResponse', version: string, date: number, notes: string }> };

export type GetNodeAppsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNodeAppsQuery = { __typename?: 'Query', getNodeApps: Array<{ __typename?: 'NodeAppResponse', id: string, title: string, version: string, cpu: string, memory: string }> };

export type GetUsersByOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetUsersByOrgQuery = { __typename?: 'Query', getUsersByOrg: Array<{ __typename?: 'GetUsersDto', id: string, name: string, email: string, dataPlan?: string | null, dataUsage?: string | null }> };

export type GetUserQueryVariables = Exact<{
  userId: Scalars['String'];
}>;


export type GetUserQuery = { __typename?: 'Query', getUser: { __typename?: 'GetUserDto', id: string, status: boolean, name: string, eSimNumber: string, iccid: string, email: string, roaming: boolean, dataPlan: string, dataUsage: string } };

export type GetNetworkStatusQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworkStatusQuery = { __typename?: 'Query', getNetworkStatus: { __typename?: 'NetworkDto', totalNodes: number, liveNode: number, status: Network_Status } };

export type GetNetworkStatusSSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetNetworkStatusSSubscription = { __typename?: 'Subscription', getNetworkStatus: { __typename?: 'NetworkDto', totalNodes: number, liveNode: number, status: Network_Status } };

export type DeactivateUserMutationVariables = Exact<{
  id: Scalars['String'];
}>;


export type DeactivateUserMutation = { __typename?: 'Mutation', deactivateUser: { __typename?: 'DeactivateResponse', uuid: string, name: string, email: string, isDeactivated: boolean } };

export type AddUserMutationVariables = Exact<{
  data: UserInputDto;
}>;


export type AddUserMutation = { __typename?: 'Mutation', addUser: { __typename?: 'UserResDto', name: string, email: string, iccid?: string | null, id: string } };

export type DeleteNodeMutationVariables = Exact<{
  id: Scalars['String'];
}>;


export type DeleteNodeMutation = { __typename?: 'Mutation', deleteNode: { __typename?: 'DeleteNodeRes', nodeId: string } };

export type UpdateUserMutationVariables = Exact<{
  userId: Scalars['String'];
  data: UserInputDto;
}>;


export type UpdateUserMutation = { __typename?: 'Mutation', updateUser: { __typename?: 'UserResDto', name: string, email: string, iccid?: string | null, id: string } };

export type AddNodeMutationVariables = Exact<{
  data: AddNodeDto;
}>;


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'AddNodeResponse', success: boolean } };

export type UpdateNodeMutationVariables = Exact<{
  data: UpdateNodeDto;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'OrgNodeDto', nodeId: string, state: Org_Node_State, type: Node_Type, name: string } };

export type GetDataBillQueryVariables = Exact<{
  data: Data_Bill_Filter;
}>;


export type GetDataBillQuery = { __typename?: 'Query', getDataBill: { __typename?: 'DataBillDto', id: string, dataBill: number, billDue: number } };

export type GetMetricsByTabQueryVariables = Exact<{
  data: MetricsByTabInputDto;
}>;


export type GetMetricsByTabQuery = { __typename?: 'Query', getMetricsByTab: { __typename?: 'GetMetricsRes', to: number, next: boolean, metrics: Array<{ __typename?: 'MetricRes', type: string, name: string, next: boolean, data: Array<{ __typename?: 'MetricDto', y: number, x: number }> }> } };

export type GetMetricsByTabSSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetMetricsByTabSSubscription = { __typename?: 'Subscription', getMetricsByTab: Array<{ __typename?: 'MetricRes', type: string, name: string, next: boolean, data: Array<{ __typename?: 'MetricDto', x: number, y: number }> }> };

export type UpdateUserStatusMutationVariables = Exact<{
  data: UpdateUserServiceInput;
}>;


export type UpdateUserStatusMutation = { __typename?: 'Mutation', updateUserStatus: { __typename?: 'OrgUserSimDto', iccid: string, isPhysical: boolean, ukama: { __typename?: 'UserServicesDto', status: Get_User_Status_Type, services: { __typename?: 'UserSimServices', voice: boolean, data: boolean, sms: boolean } }, carrier: { __typename?: 'UserServicesDto', status: Get_User_Status_Type, services: { __typename?: 'UserSimServices', voice: boolean, data: boolean, sms: boolean } } } };

export type GetNodeQueryVariables = Exact<{
  nodeId: Scalars['String'];
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'NodeResponse', nodeId: string, type: Node_Type, state: Org_Node_State, name: string, attached: Array<{ __typename?: 'OrgNodeDto', nodeId: string, type: Node_Type, state: Org_Node_State, name: string }> } };

export type GetUsersDataUsageQueryVariables = Exact<{
  data: DataUsageInputDto;
}>;


export type GetUsersDataUsageQuery = { __typename?: 'Query', getUsersDataUsage: Array<{ __typename?: 'GetUserDto', id: string, name: string, email: string, dataPlan: string, dataUsage: string }> };

export type GetUsersDataUsageSSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type GetUsersDataUsageSSubscription = { __typename?: 'Subscription', getUsersDataUsage: { __typename?: 'GetUserDto', id: string, name: string, email: string, dataPlan: string, dataUsage: string } };

export type GetNodeStatusQueryVariables = Exact<{
  data: GetNodeStatusInput;
}>;


export type GetNodeStatusQuery = { __typename?: 'Query', getNodeStatus: { __typename?: 'GetNodeStatusRes', uptime: number, status: Org_Node_State } };

export type UpdateUserRoamingMutationVariables = Exact<{
  data: UpdateUserServiceInput;
}>;


export type UpdateUserRoamingMutation = { __typename?: 'Mutation', updateUserRoaming: { __typename?: 'OrgUserSimDto', iccid: string, isPhysical: boolean, ukama: { __typename?: 'UserServicesDto', status: Get_User_Status_Type, services: { __typename?: 'UserSimServices', data: boolean } }, carrier: { __typename?: 'UserServicesDto', services: { __typename?: 'UserSimServices', data: boolean } } } };


export const GetDataUsageDocument = gql`
    query getDataUsage($filter: TIME_FILTER!) {
  getDataUsage(filter: $filter) {
    id
    dataConsumed
    dataPackage
  }
}
    `;

/**
 * __useGetDataUsageQuery__
 *
 * To run a query within a React component, call `useGetDataUsageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDataUsageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDataUsageQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetDataUsageQuery(baseOptions: Apollo.QueryHookOptions<GetDataUsageQuery, GetDataUsageQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDataUsageQuery, GetDataUsageQueryVariables>(GetDataUsageDocument, options);
      }
export function useGetDataUsageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDataUsageQuery, GetDataUsageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDataUsageQuery, GetDataUsageQueryVariables>(GetDataUsageDocument, options);
        }
export type GetDataUsageQueryHookResult = ReturnType<typeof useGetDataUsageQuery>;
export type GetDataUsageLazyQueryHookResult = ReturnType<typeof useGetDataUsageLazyQuery>;
export type GetDataUsageQueryResult = Apollo.QueryResult<GetDataUsageQuery, GetDataUsageQueryVariables>;
export const GetLatestDataUsageDocument = gql`
    subscription getLatestDataUsage {
  getDataUsage {
    id
    dataConsumed
    dataPackage
  }
}
    `;

/**
 * __useGetLatestDataUsageSubscription__
 *
 * To run a query within a React component, call `useGetLatestDataUsageSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestDataUsageSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestDataUsageSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetLatestDataUsageSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetLatestDataUsageSubscription, GetLatestDataUsageSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetLatestDataUsageSubscription, GetLatestDataUsageSubscriptionVariables>(GetLatestDataUsageDocument, options);
      }
export type GetLatestDataUsageSubscriptionHookResult = ReturnType<typeof useGetLatestDataUsageSubscription>;
export type GetLatestDataUsageSubscriptionResult = Apollo.SubscriptionResult<GetLatestDataUsageSubscription>;
export const GetConnectedUsersDocument = gql`
    query getConnectedUsers($filter: TIME_FILTER!) {
  getConnectedUsers(filter: $filter) {
    totalUser
  }
}
    `;

/**
 * __useGetConnectedUsersQuery__
 *
 * To run a query within a React component, call `useGetConnectedUsersQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetConnectedUsersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetConnectedUsersQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetConnectedUsersQuery(baseOptions: Apollo.QueryHookOptions<GetConnectedUsersQuery, GetConnectedUsersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetConnectedUsersQuery, GetConnectedUsersQueryVariables>(GetConnectedUsersDocument, options);
      }
export function useGetConnectedUsersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetConnectedUsersQuery, GetConnectedUsersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetConnectedUsersQuery, GetConnectedUsersQueryVariables>(GetConnectedUsersDocument, options);
        }
export type GetConnectedUsersQueryHookResult = ReturnType<typeof useGetConnectedUsersQuery>;
export type GetConnectedUsersLazyQueryHookResult = ReturnType<typeof useGetConnectedUsersLazyQuery>;
export type GetConnectedUsersQueryResult = Apollo.QueryResult<GetConnectedUsersQuery, GetConnectedUsersQueryVariables>;
export const GetLatestConnectedUsersDocument = gql`
    subscription getLatestConnectedUsers {
  getConnectedUsers {
    totalUser
  }
}
    `;

/**
 * __useGetLatestConnectedUsersSubscription__
 *
 * To run a query within a React component, call `useGetLatestConnectedUsersSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestConnectedUsersSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestConnectedUsersSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetLatestConnectedUsersSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetLatestConnectedUsersSubscription, GetLatestConnectedUsersSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetLatestConnectedUsersSubscription, GetLatestConnectedUsersSubscriptionVariables>(GetLatestConnectedUsersDocument, options);
      }
export type GetLatestConnectedUsersSubscriptionHookResult = ReturnType<typeof useGetLatestConnectedUsersSubscription>;
export type GetLatestConnectedUsersSubscriptionResult = Apollo.SubscriptionResult<GetLatestConnectedUsersSubscription>;
export const GetLatestDataBillDocument = gql`
    subscription getLatestDataBill {
  getDataBill {
    id
    dataBill
    billDue
  }
}
    `;

/**
 * __useGetLatestDataBillSubscription__
 *
 * To run a query within a React component, call `useGetLatestDataBillSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestDataBillSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestDataBillSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetLatestDataBillSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetLatestDataBillSubscription, GetLatestDataBillSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetLatestDataBillSubscription, GetLatestDataBillSubscriptionVariables>(GetLatestDataBillDocument, options);
      }
export type GetLatestDataBillSubscriptionHookResult = ReturnType<typeof useGetLatestDataBillSubscription>;
export type GetLatestDataBillSubscriptionResult = Apollo.SubscriptionResult<GetLatestDataBillSubscription>;
export const GetBillHistoryDocument = gql`
    query getBillHistory {
  getBillHistory {
    id
    date
    description
    totalUsage
    subtotal
  }
}
    `;

/**
 * __useGetBillHistoryQuery__
 *
 * To run a query within a React component, call `useGetBillHistoryQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetBillHistoryQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetBillHistoryQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetBillHistoryQuery(baseOptions?: Apollo.QueryHookOptions<GetBillHistoryQuery, GetBillHistoryQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetBillHistoryQuery, GetBillHistoryQueryVariables>(GetBillHistoryDocument, options);
      }
export function useGetBillHistoryLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetBillHistoryQuery, GetBillHistoryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetBillHistoryQuery, GetBillHistoryQueryVariables>(GetBillHistoryDocument, options);
        }
export type GetBillHistoryQueryHookResult = ReturnType<typeof useGetBillHistoryQuery>;
export type GetBillHistoryLazyQueryHookResult = ReturnType<typeof useGetBillHistoryLazyQuery>;
export type GetBillHistoryQueryResult = Apollo.QueryResult<GetBillHistoryQuery, GetBillHistoryQueryVariables>;
export const GetCurrentBillDocument = gql`
    query getCurrentBill {
  getCurrentBill {
    bill {
      id
      name
      dataUsed
      rate
      subtotal
    }
    total
    billMonth
    dueDate
  }
}
    `;

/**
 * __useGetCurrentBillQuery__
 *
 * To run a query within a React component, call `useGetCurrentBillQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCurrentBillQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCurrentBillQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetCurrentBillQuery(baseOptions?: Apollo.QueryHookOptions<GetCurrentBillQuery, GetCurrentBillQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCurrentBillQuery, GetCurrentBillQueryVariables>(GetCurrentBillDocument, options);
      }
export function useGetCurrentBillLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCurrentBillQuery, GetCurrentBillQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCurrentBillQuery, GetCurrentBillQueryVariables>(GetCurrentBillDocument, options);
        }
export type GetCurrentBillQueryHookResult = ReturnType<typeof useGetCurrentBillQuery>;
export type GetCurrentBillLazyQueryHookResult = ReturnType<typeof useGetCurrentBillLazyQuery>;
export type GetCurrentBillQueryResult = Apollo.QueryResult<GetCurrentBillQuery, GetCurrentBillQueryVariables>;
export const GetAlertsDocument = gql`
    query getAlerts($data: PaginationDto!) {
  getAlerts(data: $data) {
    meta {
      count
      page
      size
      pages
    }
    alerts {
      id
      type
      title
      description
      alertDate
    }
  }
}
    `;

/**
 * __useGetAlertsQuery__
 *
 * To run a query within a React component, call `useGetAlertsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetAlertsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetAlertsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetAlertsQuery(baseOptions: Apollo.QueryHookOptions<GetAlertsQuery, GetAlertsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetAlertsQuery, GetAlertsQueryVariables>(GetAlertsDocument, options);
      }
export function useGetAlertsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetAlertsQuery, GetAlertsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetAlertsQuery, GetAlertsQueryVariables>(GetAlertsDocument, options);
        }
export type GetAlertsQueryHookResult = ReturnType<typeof useGetAlertsQuery>;
export type GetAlertsLazyQueryHookResult = ReturnType<typeof useGetAlertsLazyQuery>;
export type GetAlertsQueryResult = Apollo.QueryResult<GetAlertsQuery, GetAlertsQueryVariables>;
export const GetEsimQrDocument = gql`
    query getEsimQR($data: GetESimQRCodeInput!) {
  getEsimQR(data: $data) {
    qrCode
  }
}
    `;

/**
 * __useGetEsimQrQuery__
 *
 * To run a query within a React component, call `useGetEsimQrQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetEsimQrQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetEsimQrQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetEsimQrQuery(baseOptions: Apollo.QueryHookOptions<GetEsimQrQuery, GetEsimQrQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetEsimQrQuery, GetEsimQrQueryVariables>(GetEsimQrDocument, options);
      }
export function useGetEsimQrLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetEsimQrQuery, GetEsimQrQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetEsimQrQuery, GetEsimQrQueryVariables>(GetEsimQrDocument, options);
        }
export type GetEsimQrQueryHookResult = ReturnType<typeof useGetEsimQrQuery>;
export type GetEsimQrLazyQueryHookResult = ReturnType<typeof useGetEsimQrLazyQuery>;
export type GetEsimQrQueryResult = Apollo.QueryResult<GetEsimQrQuery, GetEsimQrQueryVariables>;
export const GetLatestAlertsDocument = gql`
    subscription getLatestAlerts {
  getAlerts {
    id
    type
    title
    description
    alertDate
  }
}
    `;

/**
 * __useGetLatestAlertsSubscription__
 *
 * To run a query within a React component, call `useGetLatestAlertsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestAlertsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestAlertsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetLatestAlertsSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetLatestAlertsSubscription, GetLatestAlertsSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetLatestAlertsSubscription, GetLatestAlertsSubscriptionVariables>(GetLatestAlertsDocument, options);
      }
export type GetLatestAlertsSubscriptionHookResult = ReturnType<typeof useGetLatestAlertsSubscription>;
export type GetLatestAlertsSubscriptionResult = Apollo.SubscriptionResult<GetLatestAlertsSubscription>;
export const GetNodesByOrgDocument = gql`
    query getNodesByOrg {
  getNodesByOrg {
    orgId
    nodes {
      id
      status
      name
      type
      description
      totalUser
      isUpdateAvailable
      updateVersion
      updateShortNote
      updateDescription
    }
    activeNodes
    totalNodes
  }
}
    `;

/**
 * __useGetNodesByOrgQuery__
 *
 * To run a query within a React component, call `useGetNodesByOrgQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesByOrgQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesByOrgQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNodesByOrgQuery(baseOptions?: Apollo.QueryHookOptions<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>(GetNodesByOrgDocument, options);
      }
export function useGetNodesByOrgLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>(GetNodesByOrgDocument, options);
        }
export type GetNodesByOrgQueryHookResult = ReturnType<typeof useGetNodesByOrgQuery>;
export type GetNodesByOrgLazyQueryHookResult = ReturnType<typeof useGetNodesByOrgLazyQuery>;
export type GetNodesByOrgQueryResult = Apollo.QueryResult<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>;
export const GetNodeAppsVersionLogsDocument = gql`
    query getNodeAppsVersionLogs {
  getNodeAppsVersionLogs {
    version
    date
    notes
  }
}
    `;

/**
 * __useGetNodeAppsVersionLogsQuery__
 *
 * To run a query within a React component, call `useGetNodeAppsVersionLogsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeAppsVersionLogsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeAppsVersionLogsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNodeAppsVersionLogsQuery(baseOptions?: Apollo.QueryHookOptions<GetNodeAppsVersionLogsQuery, GetNodeAppsVersionLogsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeAppsVersionLogsQuery, GetNodeAppsVersionLogsQueryVariables>(GetNodeAppsVersionLogsDocument, options);
      }
export function useGetNodeAppsVersionLogsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeAppsVersionLogsQuery, GetNodeAppsVersionLogsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeAppsVersionLogsQuery, GetNodeAppsVersionLogsQueryVariables>(GetNodeAppsVersionLogsDocument, options);
        }
export type GetNodeAppsVersionLogsQueryHookResult = ReturnType<typeof useGetNodeAppsVersionLogsQuery>;
export type GetNodeAppsVersionLogsLazyQueryHookResult = ReturnType<typeof useGetNodeAppsVersionLogsLazyQuery>;
export type GetNodeAppsVersionLogsQueryResult = Apollo.QueryResult<GetNodeAppsVersionLogsQuery, GetNodeAppsVersionLogsQueryVariables>;
export const GetNodeAppsDocument = gql`
    query getNodeApps {
  getNodeApps {
    id
    title
    version
    cpu
    memory
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
 *   },
 * });
 */
export function useGetNodeAppsQuery(baseOptions?: Apollo.QueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
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
export const GetUsersByOrgDocument = gql`
    query getUsersByOrg {
  getUsersByOrg {
    id
    name
    email
    dataPlan
    dataUsage
  }
}
    `;

/**
 * __useGetUsersByOrgQuery__
 *
 * To run a query within a React component, call `useGetUsersByOrgQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersByOrgQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersByOrgQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetUsersByOrgQuery(baseOptions?: Apollo.QueryHookOptions<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>(GetUsersByOrgDocument, options);
      }
export function useGetUsersByOrgLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>(GetUsersByOrgDocument, options);
        }
export type GetUsersByOrgQueryHookResult = ReturnType<typeof useGetUsersByOrgQuery>;
export type GetUsersByOrgLazyQueryHookResult = ReturnType<typeof useGetUsersByOrgLazyQuery>;
export type GetUsersByOrgQueryResult = Apollo.QueryResult<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>;
export const GetUserDocument = gql`
    query getUser($userId: String!) {
  getUser(userId: $userId) {
    id
    status
    name
    eSimNumber
    iccid
    email
    roaming
    dataPlan
    dataUsage
  }
}
    `;

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
export const GetNetworkStatusDocument = gql`
    query getNetworkStatus {
  getNetworkStatus {
    totalNodes
    liveNode
    status
  }
}
    `;

/**
 * __useGetNetworkStatusQuery__
 *
 * To run a query within a React component, call `useGetNetworkStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkStatusQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNetworkStatusQuery(baseOptions?: Apollo.QueryHookOptions<GetNetworkStatusQuery, GetNetworkStatusQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkStatusQuery, GetNetworkStatusQueryVariables>(GetNetworkStatusDocument, options);
      }
export function useGetNetworkStatusLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkStatusQuery, GetNetworkStatusQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkStatusQuery, GetNetworkStatusQueryVariables>(GetNetworkStatusDocument, options);
        }
export type GetNetworkStatusQueryHookResult = ReturnType<typeof useGetNetworkStatusQuery>;
export type GetNetworkStatusLazyQueryHookResult = ReturnType<typeof useGetNetworkStatusLazyQuery>;
export type GetNetworkStatusQueryResult = Apollo.QueryResult<GetNetworkStatusQuery, GetNetworkStatusQueryVariables>;
export const GetNetworkStatusSDocument = gql`
    subscription getNetworkStatusS {
  getNetworkStatus {
    totalNodes
    liveNode
    status
  }
}
    `;

/**
 * __useGetNetworkStatusSSubscription__
 *
 * To run a query within a React component, call `useGetNetworkStatusSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkStatusSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkStatusSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetNetworkStatusSSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetNetworkStatusSSubscription, GetNetworkStatusSSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetNetworkStatusSSubscription, GetNetworkStatusSSubscriptionVariables>(GetNetworkStatusSDocument, options);
      }
export type GetNetworkStatusSSubscriptionHookResult = ReturnType<typeof useGetNetworkStatusSSubscription>;
export type GetNetworkStatusSSubscriptionResult = Apollo.SubscriptionResult<GetNetworkStatusSSubscription>;
export const DeactivateUserDocument = gql`
    mutation deactivateUser($id: String!) {
  deactivateUser(id: $id) {
    uuid
    name
    email
    isDeactivated
  }
}
    `;
export type DeactivateUserMutationFn = Apollo.MutationFunction<DeactivateUserMutation, DeactivateUserMutationVariables>;

/**
 * __useDeactivateUserMutation__
 *
 * To run a mutation, you first call `useDeactivateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeactivateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deactivateUserMutation, { data, loading, error }] = useDeactivateUserMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeactivateUserMutation(baseOptions?: Apollo.MutationHookOptions<DeactivateUserMutation, DeactivateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeactivateUserMutation, DeactivateUserMutationVariables>(DeactivateUserDocument, options);
      }
export type DeactivateUserMutationHookResult = ReturnType<typeof useDeactivateUserMutation>;
export type DeactivateUserMutationResult = Apollo.MutationResult<DeactivateUserMutation>;
export type DeactivateUserMutationOptions = Apollo.BaseMutationOptions<DeactivateUserMutation, DeactivateUserMutationVariables>;
export const AddUserDocument = gql`
    mutation addUser($data: UserInputDto!) {
  addUser(data: $data) {
    name
    email
    iccid
    id
  }
}
    `;
export type AddUserMutationFn = Apollo.MutationFunction<AddUserMutation, AddUserMutationVariables>;

/**
 * __useAddUserMutation__
 *
 * To run a mutation, you first call `useAddUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addUserMutation, { data, loading, error }] = useAddUserMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddUserMutation(baseOptions?: Apollo.MutationHookOptions<AddUserMutation, AddUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddUserMutation, AddUserMutationVariables>(AddUserDocument, options);
      }
export type AddUserMutationHookResult = ReturnType<typeof useAddUserMutation>;
export type AddUserMutationResult = Apollo.MutationResult<AddUserMutation>;
export type AddUserMutationOptions = Apollo.BaseMutationOptions<AddUserMutation, AddUserMutationVariables>;
export const DeleteNodeDocument = gql`
    mutation deleteNode($id: String!) {
  deleteNode(id: $id) {
    nodeId
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
 *      id: // value for 'id'
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
export const UpdateUserDocument = gql`
    mutation updateUser($userId: String!, $data: UserInputDto!) {
  updateUser(data: $data, userId: $userId) {
    name
    email
    iccid
    id
  }
}
    `;
export type UpdateUserMutationFn = Apollo.MutationFunction<UpdateUserMutation, UpdateUserMutationVariables>;

/**
 * __useUpdateUserMutation__
 *
 * To run a mutation, you first call `useUpdateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateUserMutation, { data, loading, error }] = useUpdateUserMutation({
 *   variables: {
 *      userId: // value for 'userId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateUserMutation(baseOptions?: Apollo.MutationHookOptions<UpdateUserMutation, UpdateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateUserMutation, UpdateUserMutationVariables>(UpdateUserDocument, options);
      }
export type UpdateUserMutationHookResult = ReturnType<typeof useUpdateUserMutation>;
export type UpdateUserMutationResult = Apollo.MutationResult<UpdateUserMutation>;
export type UpdateUserMutationOptions = Apollo.BaseMutationOptions<UpdateUserMutation, UpdateUserMutationVariables>;
export const AddNodeDocument = gql`
    mutation addNode($data: AddNodeDto!) {
  addNode(data: $data) {
    success
  }
}
    `;
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
export const UpdateNodeDocument = gql`
    mutation updateNode($data: UpdateNodeDto!) {
  updateNode(data: $data) {
    nodeId
    state
    type
    name
  }
}
    `;
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
export const GetDataBillDocument = gql`
    query getDataBill($data: DATA_BILL_FILTER!) {
  getDataBill(filter: $data) {
    id
    dataBill
    billDue
  }
}
    `;

/**
 * __useGetDataBillQuery__
 *
 * To run a query within a React component, call `useGetDataBillQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDataBillQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDataBillQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetDataBillQuery(baseOptions: Apollo.QueryHookOptions<GetDataBillQuery, GetDataBillQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDataBillQuery, GetDataBillQueryVariables>(GetDataBillDocument, options);
      }
export function useGetDataBillLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDataBillQuery, GetDataBillQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDataBillQuery, GetDataBillQueryVariables>(GetDataBillDocument, options);
        }
export type GetDataBillQueryHookResult = ReturnType<typeof useGetDataBillQuery>;
export type GetDataBillLazyQueryHookResult = ReturnType<typeof useGetDataBillLazyQuery>;
export type GetDataBillQueryResult = Apollo.QueryResult<GetDataBillQuery, GetDataBillQueryVariables>;
export const GetMetricsByTabDocument = gql`
    query getMetricsByTab($data: MetricsByTabInputDTO!) {
  getMetricsByTab(data: $data) {
    to
    next
    metrics {
      type
      name
      next
      data {
        y
        x
      }
    }
  }
}
    `;

/**
 * __useGetMetricsByTabQuery__
 *
 * To run a query within a React component, call `useGetMetricsByTabQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsByTabQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsByTabQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsByTabQuery(baseOptions: Apollo.QueryHookOptions<GetMetricsByTabQuery, GetMetricsByTabQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricsByTabQuery, GetMetricsByTabQueryVariables>(GetMetricsByTabDocument, options);
      }
export function useGetMetricsByTabLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricsByTabQuery, GetMetricsByTabQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricsByTabQuery, GetMetricsByTabQueryVariables>(GetMetricsByTabDocument, options);
        }
export type GetMetricsByTabQueryHookResult = ReturnType<typeof useGetMetricsByTabQuery>;
export type GetMetricsByTabLazyQueryHookResult = ReturnType<typeof useGetMetricsByTabLazyQuery>;
export type GetMetricsByTabQueryResult = Apollo.QueryResult<GetMetricsByTabQuery, GetMetricsByTabQueryVariables>;
export const GetMetricsByTabSDocument = gql`
    subscription getMetricsByTabS {
  getMetricsByTab {
    type
    name
    next
    data {
      x
      y
    }
  }
}
    `;

/**
 * __useGetMetricsByTabSSubscription__
 *
 * To run a query within a React component, call `useGetMetricsByTabSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsByTabSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsByTabSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsByTabSSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetMetricsByTabSSubscription, GetMetricsByTabSSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetMetricsByTabSSubscription, GetMetricsByTabSSubscriptionVariables>(GetMetricsByTabSDocument, options);
      }
export type GetMetricsByTabSSubscriptionHookResult = ReturnType<typeof useGetMetricsByTabSSubscription>;
export type GetMetricsByTabSSubscriptionResult = Apollo.SubscriptionResult<GetMetricsByTabSSubscription>;
export const UpdateUserStatusDocument = gql`
    mutation updateUserStatus($data: UpdateUserServiceInput!) {
  updateUserStatus(data: $data) {
    iccid
    isPhysical
    ukama {
      status
      services {
        voice
        data
        sms
      }
    }
    carrier {
      status
      services {
        voice
        data
        sms
      }
    }
  }
}
    `;
export type UpdateUserStatusMutationFn = Apollo.MutationFunction<UpdateUserStatusMutation, UpdateUserStatusMutationVariables>;

/**
 * __useUpdateUserStatusMutation__
 *
 * To run a mutation, you first call `useUpdateUserStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateUserStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateUserStatusMutation, { data, loading, error }] = useUpdateUserStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateUserStatusMutation(baseOptions?: Apollo.MutationHookOptions<UpdateUserStatusMutation, UpdateUserStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateUserStatusMutation, UpdateUserStatusMutationVariables>(UpdateUserStatusDocument, options);
      }
export type UpdateUserStatusMutationHookResult = ReturnType<typeof useUpdateUserStatusMutation>;
export type UpdateUserStatusMutationResult = Apollo.MutationResult<UpdateUserStatusMutation>;
export type UpdateUserStatusMutationOptions = Apollo.BaseMutationOptions<UpdateUserStatusMutation, UpdateUserStatusMutationVariables>;
export const GetNodeDocument = gql`
    query getNode($nodeId: String!) {
  getNode(nodeId: $nodeId) {
    nodeId
    type
    state
    name
    attached {
      nodeId
      type
      state
      name
    }
  }
}
    `;

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
 *      nodeId: // value for 'nodeId'
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
export const GetUsersDataUsageDocument = gql`
    query getUsersDataUsage($data: DataUsageInputDto!) {
  getUsersDataUsage(data: $data) {
    id
    name
    email
    dataPlan
    dataUsage
  }
}
    `;

/**
 * __useGetUsersDataUsageQuery__
 *
 * To run a query within a React component, call `useGetUsersDataUsageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersDataUsageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersDataUsageQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetUsersDataUsageQuery(baseOptions: Apollo.QueryHookOptions<GetUsersDataUsageQuery, GetUsersDataUsageQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUsersDataUsageQuery, GetUsersDataUsageQueryVariables>(GetUsersDataUsageDocument, options);
      }
export function useGetUsersDataUsageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUsersDataUsageQuery, GetUsersDataUsageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUsersDataUsageQuery, GetUsersDataUsageQueryVariables>(GetUsersDataUsageDocument, options);
        }
export type GetUsersDataUsageQueryHookResult = ReturnType<typeof useGetUsersDataUsageQuery>;
export type GetUsersDataUsageLazyQueryHookResult = ReturnType<typeof useGetUsersDataUsageLazyQuery>;
export type GetUsersDataUsageQueryResult = Apollo.QueryResult<GetUsersDataUsageQuery, GetUsersDataUsageQueryVariables>;
export const GetUsersDataUsageSDocument = gql`
    subscription getUsersDataUsageS {
  getUsersDataUsage {
    id
    name
    email
    dataPlan
    dataUsage
  }
}
    `;

/**
 * __useGetUsersDataUsageSSubscription__
 *
 * To run a query within a React component, call `useGetUsersDataUsageSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersDataUsageSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersDataUsageSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetUsersDataUsageSSubscription(baseOptions?: Apollo.SubscriptionHookOptions<GetUsersDataUsageSSubscription, GetUsersDataUsageSSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetUsersDataUsageSSubscription, GetUsersDataUsageSSubscriptionVariables>(GetUsersDataUsageSDocument, options);
      }
export type GetUsersDataUsageSSubscriptionHookResult = ReturnType<typeof useGetUsersDataUsageSSubscription>;
export type GetUsersDataUsageSSubscriptionResult = Apollo.SubscriptionResult<GetUsersDataUsageSSubscription>;
export const GetNodeStatusDocument = gql`
    query getNodeStatus($data: GetNodeStatusInput!) {
  getNodeStatus(data: $data) {
    uptime
    status
  }
}
    `;

/**
 * __useGetNodeStatusQuery__
 *
 * To run a query within a React component, call `useGetNodeStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeStatusQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeStatusQuery(baseOptions: Apollo.QueryHookOptions<GetNodeStatusQuery, GetNodeStatusQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeStatusQuery, GetNodeStatusQueryVariables>(GetNodeStatusDocument, options);
      }
export function useGetNodeStatusLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeStatusQuery, GetNodeStatusQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeStatusQuery, GetNodeStatusQueryVariables>(GetNodeStatusDocument, options);
        }
export type GetNodeStatusQueryHookResult = ReturnType<typeof useGetNodeStatusQuery>;
export type GetNodeStatusLazyQueryHookResult = ReturnType<typeof useGetNodeStatusLazyQuery>;
export type GetNodeStatusQueryResult = Apollo.QueryResult<GetNodeStatusQuery, GetNodeStatusQueryVariables>;
export const UpdateUserRoamingDocument = gql`
    mutation updateUserRoaming($data: UpdateUserServiceInput!) {
  updateUserRoaming(data: $data) {
    iccid
    isPhysical
    ukama {
      status
      services {
        data
      }
    }
    carrier {
      services {
        data
      }
    }
  }
}
    `;
export type UpdateUserRoamingMutationFn = Apollo.MutationFunction<UpdateUserRoamingMutation, UpdateUserRoamingMutationVariables>;

/**
 * __useUpdateUserRoamingMutation__
 *
 * To run a mutation, you first call `useUpdateUserRoamingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateUserRoamingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateUserRoamingMutation, { data, loading, error }] = useUpdateUserRoamingMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateUserRoamingMutation(baseOptions?: Apollo.MutationHookOptions<UpdateUserRoamingMutation, UpdateUserRoamingMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateUserRoamingMutation, UpdateUserRoamingMutationVariables>(UpdateUserRoamingDocument, options);
      }
export type UpdateUserRoamingMutationHookResult = ReturnType<typeof useUpdateUserRoamingMutation>;
export type UpdateUserRoamingMutationResult = Apollo.MutationResult<UpdateUserRoamingMutation>;
export type UpdateUserRoamingMutationOptions = Apollo.BaseMutationOptions<UpdateUserRoamingMutation, UpdateUserRoamingMutationVariables>;