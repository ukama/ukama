import { gql } from "@apollo/client";
import * as Apollo from "@apollo/client";
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
    [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
    [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
    [SubKey in K]: Maybe<T[SubKey]>;
};
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
    Error = "ERROR",
    Info = "INFO",
    Warning = "WARNING",
}

export enum Api_Method_Type {
    Delete = "DELETE",
    Get = "GET",
    Post = "POST",
    Put = "PUT",
}

export type ActivateUserDto = {
    dataPlan: Scalars["Float"];
    dataUsage: Scalars["Float"];
    eSimNumber: Scalars["String"];
    email?: InputMaybe<Scalars["String"]>;
    iccid: Scalars["String"];
    name: Scalars["String"];
    phone?: InputMaybe<Scalars["String"]>;
    roaming: Scalars["Boolean"];
};

export type ActivateUserResponse = {
    __typename?: "ActivateUserResponse";
    success: Scalars["Boolean"];
};

export type ActiveUserResponseDto = {
    __typename?: "ActiveUserResponseDto";
    data: ActivateUserResponse;
    status: Scalars["String"];
};

export type AddNodeDto = {
    name: Scalars["String"];
    securityCode: Scalars["String"];
    serialNo: Scalars["String"];
};

export type AddNodeResponse = {
    __typename?: "AddNodeResponse";
    success: Scalars["Boolean"];
};

export type AddNodeResponseDto = {
    __typename?: "AddNodeResponseDto";
    data: AddNodeResponse;
    status: Scalars["String"];
};

export type AddUserDto = {
    email: Scalars["String"];
    firstName: Scalars["String"];
    imsi: Scalars["String"];
    lastName: Scalars["String"];
};

export type AddUserResponse = {
    __typename?: "AddUserResponse";
    email: Scalars["String"];
    firstName: Scalars["String"];
    imsi: Scalars["String"];
    lastName: Scalars["String"];
    uuid: Scalars["String"];
};

export type AlertDto = {
    __typename?: "AlertDto";
    alertDate?: Maybe<Scalars["DateTime"]>;
    description?: Maybe<Scalars["String"]>;
    id?: Maybe<Scalars["String"]>;
    title?: Maybe<Scalars["String"]>;
    type: Alert_Type;
};

export type AlertResponse = {
    __typename?: "AlertResponse";
    data: Array<AlertDto>;
    length: Scalars["Float"];
    status: Scalars["String"];
};

export type AlertsResponse = {
    __typename?: "AlertsResponse";
    alerts: Array<AlertDto>;
    meta: Meta;
};

export type ApiMethodDataDto = {
    __typename?: "ApiMethodDataDto";
    body?: Maybe<Scalars["String"]>;
    headers?: Maybe<Scalars["String"]>;
    params?: Maybe<Scalars["String"]>;
    path: Scalars["String"];
    type: Api_Method_Type;
};

export type BillHistoryDto = {
    __typename?: "BillHistoryDto";
    date: Scalars["String"];
    description: Scalars["String"];
    id: Scalars["String"];
    subtotal: Scalars["Float"];
    totalUsage: Scalars["Float"];
};

export type BillHistoryResponse = {
    __typename?: "BillHistoryResponse";
    data: Array<BillHistoryDto>;
    status: Scalars["String"];
};

export type BillResponse = {
    __typename?: "BillResponse";
    bill: Array<CurrentBillDto>;
    billMonth: Scalars["String"];
    dueDate: Scalars["String"];
    total: Scalars["Float"];
};

export enum Connected_User_Type {
    Guests = "GUESTS",
    Residents = "RESIDENTS",
}

export type ConnectedUserDto = {
    __typename?: "ConnectedUserDto";
    totalUser: Scalars["Float"];
};

export type ConnectedUserResponse = {
    __typename?: "ConnectedUserResponse";
    data: ConnectedUserDto;
    status: Scalars["String"];
};

export type CpuUsageMetricsDto = {
    __typename?: "CpuUsageMetricsDto";
    id?: Maybe<Scalars["String"]>;
    timestamp: Scalars["Float"];
    usage: Scalars["Float"];
};

export type CurrentBillDto = {
    __typename?: "CurrentBillDto";
    dataUsed: Scalars["Float"];
    id: Scalars["String"];
    name: Scalars["String"];
    rate: Scalars["Float"];
    subtotal: Scalars["Float"];
};

export type CurrentBillResponse = {
    __typename?: "CurrentBillResponse";
    data: Array<CurrentBillDto>;
    status: Scalars["String"];
};

export enum Data_Bill_Filter {
    April = "APRIL",
    August = "AUGUST",
    Current = "CURRENT",
    December = "DECEMBER",
    Februray = "FEBRURAY",
    January = "JANUARY",
    July = "JULY",
    June = "JUNE",
    March = "MARCH",
    May = "MAY",
    Novermber = "NOVERMBER",
    October = "OCTOBER",
    September = "SEPTEMBER",
}

export type DataBillDto = {
    __typename?: "DataBillDto";
    billDue: Scalars["Float"];
    dataBill: Scalars["Float"];
    id: Scalars["String"];
};

export type DataBillResponse = {
    __typename?: "DataBillResponse";
    data: DataBillDto;
    status: Scalars["String"];
};

export type DataUsageDto = {
    __typename?: "DataUsageDto";
    dataConsumed: Scalars["Float"];
    dataPackage: Scalars["String"];
    id: Scalars["String"];
};

export type DataUsageResponse = {
    __typename?: "DataUsageResponse";
    data: DataUsageDto;
    status: Scalars["String"];
};

export type DeactivateResponse = {
    __typename?: "DeactivateResponse";
    id: Scalars["String"];
    success: Scalars["Boolean"];
};

export type ErrorType = {
    __typename?: "ErrorType";
    code: Scalars["Float"];
    description?: Maybe<Scalars["String"]>;
    message: Scalars["String"];
};

export type EsimDto = {
    __typename?: "EsimDto";
    active: Scalars["Boolean"];
    esim: Scalars["String"];
};

export type EsimResponse = {
    __typename?: "EsimResponse";
    data: Array<EsimDto>;
    status: Scalars["String"];
};

export enum Get_User_Status_Type {
    Active = "ACTIVE",
    Inactive = "INACTIVE",
}

export enum Get_User_Type {
    All = "ALL",
    Guest = "GUEST",
    Home = "HOME",
    Resident = "RESIDENT",
    Visitor = "VISITOR",
}

export enum Graph_Filter {
    Day = "DAY",
    Month = "MONTH",
    Week = "WEEK",
}

export type GetUserDto = {
    __typename?: "GetUserDto";
    dataPlan: Scalars["Float"];
    dataUsage: Scalars["Float"];
    eSimNumber: Scalars["String"];
    email?: Maybe<Scalars["String"]>;
    iccid: Scalars["String"];
    id: Scalars["String"];
    name: Scalars["String"];
    phone?: Maybe<Scalars["String"]>;
    roaming: Scalars["Boolean"];
    status: Get_User_Status_Type;
};

export type GetUserPaginationDto = {
    pageNo: Scalars["Float"];
    pageSize: Scalars["Float"];
    type: Get_User_Type;
};

export type GetUserResponse = {
    __typename?: "GetUserResponse";
    meta: Meta;
    users: Array<GetUserDto>;
};

export type GetUserResponseDto = {
    __typename?: "GetUserResponseDto";
    data: Array<GetUserDto>;
    length: Scalars["Float"];
    status: Scalars["String"];
};

export type HeaderType = {
    __typename?: "HeaderType";
    Authorization: Scalars["String"];
    Cookie: Scalars["String"];
};

export type IoMetricsDto = {
    __typename?: "IOMetricsDto";
    id?: Maybe<Scalars["String"]>;
    input: Scalars["Float"];
    output: Scalars["Float"];
    timestamp: Scalars["Float"];
};

export type MemoryUsageMetricsDto = {
    __typename?: "MemoryUsageMetricsDto";
    id?: Maybe<Scalars["String"]>;
    timestamp: Scalars["Float"];
    usage: Scalars["Float"];
};

export type Meta = {
    __typename?: "Meta";
    count: Scalars["Float"];
    page: Scalars["Float"];
    pages: Scalars["Float"];
    size: Scalars["Float"];
};

export type MetricDto = {
    __typename?: "MetricDto";
    x: Scalars["Float"];
    y: Scalars["Float"];
};

export type MetricsInputDto = {
    from: Scalars["Float"];
    nodeId: Scalars["String"];
    orgId: Scalars["String"];
    step: Scalars["Float"];
    to: Scalars["Float"];
};

export type Mutation = {
    __typename?: "Mutation";
    activateUser: ActivateUserResponse;
    addNode: AddNodeResponse;
    addUser: AddUserResponse;
    deactivateUser: DeactivateResponse;
    deleteNode: DeactivateResponse;
    deleteUser: ActivateUserResponse;
    updateNode: UpdateNodeResponse;
    updateUser: UserResponse;
};

export type MutationActivateUserArgs = {
    data: ActivateUserDto;
};

export type MutationAddNodeArgs = {
    data: AddNodeDto;
};

export type MutationAddUserArgs = {
    data: AddUserDto;
    orgId: Scalars["String"];
};

export type MutationDeactivateUserArgs = {
    id: Scalars["String"];
};

export type MutationDeleteNodeArgs = {
    id: Scalars["String"];
};

export type MutationDeleteUserArgs = {
    orgId: Scalars["String"];
    userId: Scalars["String"];
};

export type MutationUpdateNodeArgs = {
    data: UpdateNodeDto;
};

export type MutationUpdateUserArgs = {
    data: UpdateUserDto;
};

export enum Network_Status {
    BeingConfigured = "BEING_CONFIGURED",
    Online = "ONLINE",
}

export enum Network_Type {
    Private = "PRIVATE",
    Public = "PUBLIC",
}

export type NetworkDto = {
    __typename?: "NetworkDto";
    description?: Maybe<Scalars["String"]>;
    id: Scalars["String"];
    status: Network_Status;
};

export type NetworkResponse = {
    __typename?: "NetworkResponse";
    data: NetworkDto;
    status: Scalars["String"];
};

export type NodeDetailDto = {
    __typename?: "NodeDetailDto";
    description: Scalars["String"];
    hardware: Scalars["Float"];
    id: Scalars["String"];
    macAddress: Scalars["Float"];
    manufacturing: Scalars["Float"];
    modelType: Scalars["String"];
    osVersion: Scalars["Float"];
    serial: Scalars["Float"];
    ukamaOS: Scalars["Float"];
};

export type NodeDto = {
    __typename?: "NodeDto";
    description: Scalars["String"];
    id: Scalars["String"];
    status: Org_Node_State;
    title: Scalars["String"];
    totalUser: Scalars["Float"];
    type: Scalars["String"];
};

export type NodeMetaDataDto = {
    __typename?: "NodeMetaDataDto";
    throughput: Scalars["Float"];
    usersAttached: Scalars["Float"];
};

export type NodePhysicalHealthDto = {
    __typename?: "NodePhysicalHealthDto";
    Memory: Scalars["Float"];
    cpu: Scalars["Float"];
    io: Scalars["Float"];
    temperature: Scalars["Float"];
};

export type NodeRfDto = {
    __typename?: "NodeRFDto";
    qam: Scalars["Float"];
    rfOutput: Scalars["Float"];
    rssi: Scalars["Float"];
    timestamp: Scalars["Float"];
};

export type NodeResponse = {
    __typename?: "NodeResponse";
    data: Array<NodeDto>;
    length: Scalars["Float"];
    status: Scalars["String"];
};

export type NodeResponseDto = {
    __typename?: "NodeResponseDto";
    activeNodes: Scalars["Float"];
    nodes: Array<NodeDto>;
    totalNodes: Scalars["Float"];
};

export type NodesResponse = {
    __typename?: "NodesResponse";
    meta: Meta;
    nodes: NodeResponseDto;
};

export enum Org_Node_State {
    Onboarded = "ONBOARDED",
    Pending = "PENDING",
    Undefined = "UNDEFINED",
}

export type OrgMetricDto = {
    __typename?: "OrgMetricDto";
    nodeId: Scalars["String"];
    receive: Scalars["String"];
    tenant_id: Scalars["String"];
};

export type OrgMetricResponse = {
    __typename?: "OrgMetricResponse";
    metric: OrgMetricDto;
    values: Array<OrgMetricValueDto>;
};

export type OrgMetricValueDto = {
    __typename?: "OrgMetricValueDto";
    x: Scalars["Float"];
    y: Scalars["String"];
};

export type OrgNodeDto = {
    __typename?: "OrgNodeDto";
    nodeId: Scalars["String"];
    state: Org_Node_State;
    type: Scalars["String"];
};

export type OrgNodeResponse = {
    __typename?: "OrgNodeResponse";
    nodes: Array<OrgNodeDto>;
    orgName: Scalars["String"];
};

export type OrgNodeResponseDto = {
    __typename?: "OrgNodeResponseDto";
    activeNodes: Scalars["Float"];
    nodes: Array<NodeDto>;
    orgName: Scalars["String"];
    totalNodes: Scalars["Float"];
};

export type OrgUserDto = {
    __typename?: "OrgUserDto";
    email: Scalars["String"];
    firstName: Scalars["String"];
    lastName: Scalars["String"];
    uuid: Scalars["String"];
};

export type OrgUserResponse = {
    __typename?: "OrgUserResponse";
    org: Scalars["String"];
    users: Array<OrgUserDto>;
};

export type OrgUserResponseDto = {
    __typename?: "OrgUserResponseDto";
    orgName: Scalars["String"];
    users: Array<GetUserDto>;
};

export type PaginationDto = {
    pageNo: Scalars["Float"];
    pageSize: Scalars["Float"];
};

export type PaginationResponse = {
    __typename?: "PaginationResponse";
    meta: Meta;
};

export type Query = {
    __typename?: "Query";
    getAlerts: AlertsResponse;
    getBillHistory: Array<BillHistoryDto>;
    getConnectedUsers: ConnectedUserDto;
    getCpuUsageMetrics: Array<CpuUsageMetricsDto>;
    getCurrentBill: BillResponse;
    getDataBill: DataBillDto;
    getDataUsage: DataUsageDto;
    getEsims: Array<EsimDto>;
    getIOMetrics: Array<IoMetricsDto>;
    getMemoryUsageMetrics: Array<MemoryUsageMetricsDto>;
    getMetricsCpuTRX: Array<MetricDto>;
    getMetricsMemoryTRX: Array<MetricDto>;
    getMetricsUptime: Array<MetricDto>;
    getNetwork: NetworkDto;
    getNodeDetails: NodeDetailDto;
    getNodeMetaData: NodeMetaDataDto;
    getNodeNetwork: NetworkDto;
    getNodePhysicalHealth: NodePhysicalHealthDto;
    getNodeRFKPI: Array<NodeRfDto>;
    getNodes: NodesResponse;
    getNodesByOrg: OrgNodeResponseDto;
    getResidents: ResidentsResponse;
    getTemperatureMetrics: Array<TemperatureMetricsDto>;
    getThroughputMetrics: Array<ThroughputMetricsDto>;
    getUser: GetUserDto;
    getUsers: GetUserResponse;
    getUsersAttachedMetrics: Array<UsersAttachedMetricsDto>;
    myUsers: OrgUserResponseDto;
};

export type QueryGetAlertsArgs = {
    data: PaginationDto;
};

export type QueryGetConnectedUsersArgs = {
    filter: Time_Filter;
};

export type QueryGetCpuUsageMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryGetDataBillArgs = {
    filter: Data_Bill_Filter;
};

export type QueryGetDataUsageArgs = {
    filter: Time_Filter;
};

export type QueryGetIoMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryGetMemoryUsageMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryGetMetricsCpuTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsMemoryTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsUptimeArgs = {
    data: MetricsInputDto;
};

export type QueryGetNetworkArgs = {
    filter: Network_Type;
};

export type QueryGetNodeRfkpiArgs = {
    filter: Graph_Filter;
};

export type QueryGetNodesArgs = {
    data: PaginationDto;
};

export type QueryGetNodesByOrgArgs = {
    orgId: Scalars["String"];
};

export type QueryGetResidentsArgs = {
    data: PaginationDto;
};

export type QueryGetTemperatureMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryGetThroughputMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryGetUserArgs = {
    id: Scalars["String"];
};

export type QueryGetUsersArgs = {
    data: GetUserPaginationDto;
};

export type QueryGetUsersAttachedMetricsArgs = {
    filter: Graph_Filter;
};

export type QueryMyUsersArgs = {
    orgId: Scalars["String"];
};

export type ResidentResponse = {
    __typename?: "ResidentResponse";
    activeResidents: Scalars["Float"];
    residents: Array<GetUserDto>;
    totalResidents: Scalars["Float"];
};

export type ResidentsResponse = {
    __typename?: "ResidentsResponse";
    meta: Meta;
    residents: ResidentResponse;
};

export type Subscription = {
    __typename?: "Subscription";
    getAlerts: AlertDto;
    getConnectedUsers: ConnectedUserDto;
    getCpuUsageMetrics: CpuUsageMetricsDto;
    getDataBill: DataBillDto;
    getDataUsage: DataUsageDto;
    getIOMetrics: IoMetricsDto;
    getMemoryUsageMetrics: MemoryUsageMetricsDto;
    getNetwork: NetworkDto;
    getNodeMetaData: NodeMetaDataDto;
    getNodePhysicalHealth: NodePhysicalHealthDto;
    getNodeRFKPI: NodeRfDto;
    getTemperatureMetrics: TemperatureMetricsDto;
    getThroughputMetrics: ThroughputMetricsDto;
    getUsersAttachedMetrics: UsersAttachedMetricsDto;
};

export enum Time_Filter {
    Month = "MONTH",
    Today = "TODAY",
    Total = "TOTAL",
    Week = "WEEK",
}

export type TemperatureMetricsDto = {
    __typename?: "TemperatureMetricsDto";
    id?: Maybe<Scalars["String"]>;
    temperature: Scalars["Float"];
    timestamp: Scalars["Float"];
};

export type ThroughputMetricsDto = {
    __typename?: "ThroughputMetricsDto";
    amount: Scalars["Float"];
    id?: Maybe<Scalars["String"]>;
    timestamp: Scalars["Float"];
};

export type UpdateNodeDto = {
    id: Scalars["String"];
    name?: InputMaybe<Scalars["String"]>;
    securityCode?: InputMaybe<Scalars["String"]>;
    serialNo?: InputMaybe<Scalars["String"]>;
};

export type UpdateNodeResponse = {
    __typename?: "UpdateNodeResponse";
    id: Scalars["String"];
    name: Scalars["String"];
    serialNo: Scalars["String"];
};

export type UpdateUserDto = {
    eSimNumber?: InputMaybe<Scalars["String"]>;
    email?: InputMaybe<Scalars["String"]>;
    firstName?: InputMaybe<Scalars["String"]>;
    id: Scalars["String"];
    lastName?: InputMaybe<Scalars["String"]>;
    phone?: InputMaybe<Scalars["String"]>;
};

export type UserDto = {
    __typename?: "UserDto";
    email: Scalars["String"];
    id: Scalars["String"];
    name: Scalars["String"];
    type: Connected_User_Type;
};

export type UserResponse = {
    __typename?: "UserResponse";
    email: Scalars["String"];
    id: Scalars["String"];
    name: Scalars["String"];
    phone: Scalars["String"];
    sim: Scalars["String"];
};

export type UsersAttachedMetricsDto = {
    __typename?: "UsersAttachedMetricsDto";
    id?: Maybe<Scalars["String"]>;
    timestamp: Scalars["Float"];
    users: Scalars["Float"];
};

export type GetDataUsageQueryVariables = Exact<{
    filter: Time_Filter;
}>;

export type GetDataUsageQuery = {
    __typename?: "Query";
    getDataUsage: {
        __typename?: "DataUsageDto";
        id: string;
        dataConsumed: number;
        dataPackage: string;
    };
};

export type GetLatestDataUsageSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetLatestDataUsageSubscription = {
    __typename?: "Subscription";
    getDataUsage: {
        __typename?: "DataUsageDto";
        id: string;
        dataConsumed: number;
        dataPackage: string;
    };
};

export type GetConnectedUsersQueryVariables = Exact<{
    filter: Time_Filter;
}>;

export type GetConnectedUsersQuery = {
    __typename?: "Query";
    getConnectedUsers: { __typename?: "ConnectedUserDto"; totalUser: number };
};

export type GetLatestConnectedUsersSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetLatestConnectedUsersSubscription = {
    __typename?: "Subscription";
    getConnectedUsers: { __typename?: "ConnectedUserDto"; totalUser: number };
};

export type GetDataBillQueryVariables = Exact<{
    filter: Data_Bill_Filter;
}>;

export type GetDataBillQuery = {
    __typename?: "Query";
    getDataBill: {
        __typename?: "DataBillDto";
        id: string;
        dataBill: number;
        billDue: number;
    };
};

export type GetLatestDataBillSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetLatestDataBillSubscription = {
    __typename?: "Subscription";
    getDataBill: {
        __typename?: "DataBillDto";
        id: string;
        dataBill: number;
        billDue: number;
    };
};

export type GetAlertsQueryVariables = Exact<{
    data: PaginationDto;
}>;

export type GetAlertsQuery = {
    __typename?: "Query";
    getAlerts: {
        __typename?: "AlertsResponse";
        meta: {
            __typename?: "Meta";
            count: number;
            page: number;
            size: number;
            pages: number;
        };
        alerts: Array<{
            __typename?: "AlertDto";
            id?: string | null;
            type: Alert_Type;
            title?: string | null;
            description?: string | null;
            alertDate?: any | null;
        }>;
    };
};

export type GetLatestAlertsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetLatestAlertsSubscription = {
    __typename?: "Subscription";
    getAlerts: {
        __typename?: "AlertDto";
        id?: string | null;
        type: Alert_Type;
        title?: string | null;
        description?: string | null;
        alertDate?: any | null;
    };
};

export type GetNodesByOrgQueryVariables = Exact<{
    orgId: Scalars["String"];
}>;

export type GetNodesByOrgQuery = {
    __typename?: "Query";
    getNodesByOrg: {
        __typename?: "OrgNodeResponseDto";
        orgName: string;
        activeNodes: number;
        totalNodes: number;
        nodes: Array<{
            __typename?: "NodeDto";
            id: string;
            status: Org_Node_State;
            title: string;
            type: string;
            description: string;
            totalUser: number;
        }>;
    };
};

export type GetNodeDetailsQueryVariables = Exact<{ [key: string]: never }>;

export type GetNodeDetailsQuery = {
    __typename?: "Query";
    getNodeDetails: {
        __typename?: "NodeDetailDto";
        id: string;
        modelType: string;
        serial: number;
        macAddress: number;
        osVersion: number;
        manufacturing: number;
        ukamaOS: number;
        hardware: number;
        description: string;
    };
};

export type MyUsersQueryVariables = Exact<{
    orgId: Scalars["String"];
}>;

export type MyUsersQuery = {
    __typename?: "Query";
    myUsers: {
        __typename?: "OrgUserResponseDto";
        orgName: string;
        users: Array<{
            __typename?: "GetUserDto";
            id: string;
            name: string;
            email?: string | null;
            eSimNumber: string;
            dataPlan: number;
            dataUsage: number;
            phone?: string | null;
            roaming: boolean;
            iccid: string;
            status: Get_User_Status_Type;
        }>;
    };
};

export type GetUserQueryVariables = Exact<{
    id: Scalars["String"];
}>;

export type GetUserQuery = {
    __typename?: "Query";
    getUser: {
        __typename?: "GetUserDto";
        id: string;
        status: Get_User_Status_Type;
        name: string;
        eSimNumber: string;
        iccid: string;
        email?: string | null;
        phone?: string | null;
        roaming: boolean;
        dataPlan: number;
        dataUsage: number;
    };
};

export type GetResidentsQueryVariables = Exact<{
    data: PaginationDto;
}>;

export type GetResidentsQuery = {
    __typename?: "Query";
    getResidents: {
        __typename?: "ResidentsResponse";
        meta: {
            __typename?: "Meta";
            count: number;
            page: number;
            size: number;
            pages: number;
        };
        residents: {
            __typename?: "ResidentResponse";
            activeResidents: number;
            totalResidents: number;
            residents: Array<{
                __typename?: "GetUserDto";
                id: string;
                name: string;
                dataUsage: number;
            }>;
        };
    };
};

export type GetNetworkQueryVariables = Exact<{
    filter: Network_Type;
}>;

export type GetNetworkQuery = {
    __typename?: "Query";
    getNetwork: {
        __typename?: "NetworkDto";
        id: string;
        status: Network_Status;
        description?: string | null;
    };
};

export type GetLatestNetworkSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetLatestNetworkSubscription = {
    __typename?: "Subscription";
    getNetwork: {
        __typename?: "NetworkDto";
        id: string;
        status: Network_Status;
        description?: string | null;
    };
};

export type DeactivateUserMutationVariables = Exact<{
    id: Scalars["String"];
}>;

export type DeactivateUserMutation = {
    __typename?: "Mutation";
    deactivateUser: {
        __typename?: "DeactivateResponse";
        id: string;
        success: boolean;
    };
};

export type UpdateUserMutationVariables = Exact<{
    data: UpdateUserDto;
}>;

export type UpdateUserMutation = {
    __typename?: "Mutation";
    updateUser: {
        __typename?: "UserResponse";
        id: string;
        name: string;
        sim: string;
        email: string;
        phone: string;
    };
};

export type ActivateUserMutationVariables = Exact<{
    data: ActivateUserDto;
}>;

export type ActivateUserMutation = {
    __typename?: "Mutation";
    activateUser: { __typename?: "ActivateUserResponse"; success: boolean };
};

export type UpdateNodeMutationVariables = Exact<{
    data: UpdateNodeDto;
}>;

export type UpdateNodeMutation = {
    __typename?: "Mutation";
    updateNode: {
        __typename: "UpdateNodeResponse";
        id: string;
        name: string;
        serialNo: string;
    };
};

export type AddNodeMutationVariables = Exact<{
    data: AddNodeDto;
}>;

export type AddNodeMutation = {
    __typename?: "Mutation";
    addNode: { __typename: "AddNodeResponse"; success: boolean };
};

export type DeleteNodeMutationVariables = Exact<{
    id: Scalars["String"];
}>;

export type DeleteNodeMutation = {
    __typename?: "Mutation";
    deleteNode: {
        __typename?: "DeactivateResponse";
        id: string;
        success: boolean;
    };
};

export type GetNodeRfkpisSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetNodeRfkpisSubscription = {
    __typename?: "Subscription";
    getNodeRFKPI: {
        __typename?: "NodeRFDto";
        qam: number;
        rfOutput: number;
        rssi: number;
        timestamp: number;
    };
};

export type GetNodeRfkpiqQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetNodeRfkpiqQuery = {
    __typename?: "Query";
    getNodeRFKPI: Array<{
        __typename?: "NodeRFDto";
        qam: number;
        rfOutput: number;
        rssi: number;
        timestamp: number;
    }>;
};

export type GetUsersAttachedMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetUsersAttachedMetricsSSubscription = {
    __typename?: "Subscription";
    getUsersAttachedMetrics: {
        __typename?: "UsersAttachedMetricsDto";
        id?: string | null;
        users: number;
        timestamp: number;
    };
};

export type GetUsersAttachedMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetUsersAttachedMetricsQQuery = {
    __typename?: "Query";
    getUsersAttachedMetrics: Array<{
        __typename?: "UsersAttachedMetricsDto";
        id?: string | null;
        users: number;
        timestamp: number;
    }>;
};

export type GetThroughputMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetThroughputMetricsSSubscription = {
    __typename?: "Subscription";
    getThroughputMetrics: {
        __typename?: "ThroughputMetricsDto";
        id?: string | null;
        amount: number;
        timestamp: number;
    };
};

export type GetThroughputMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetThroughputMetricsQQuery = {
    __typename?: "Query";
    getThroughputMetrics: Array<{
        __typename?: "ThroughputMetricsDto";
        id?: string | null;
        amount: number;
        timestamp: number;
    }>;
};

export type GetTemperatureMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetTemperatureMetricsSSubscription = {
    __typename?: "Subscription";
    getTemperatureMetrics: {
        __typename?: "TemperatureMetricsDto";
        id?: string | null;
        temperature: number;
        timestamp: number;
    };
};

export type GetTemperatureMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetTemperatureMetricsQQuery = {
    __typename?: "Query";
    getTemperatureMetrics: Array<{
        __typename?: "TemperatureMetricsDto";
        id?: string | null;
        temperature: number;
        timestamp: number;
    }>;
};

export type GetCpuUsageMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetCpuUsageMetricsSSubscription = {
    __typename?: "Subscription";
    getCpuUsageMetrics: {
        __typename?: "CpuUsageMetricsDto";
        id?: string | null;
        usage: number;
        timestamp: number;
    };
};

export type GetCpuUsageMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetCpuUsageMetricsQQuery = {
    __typename?: "Query";
    getCpuUsageMetrics: Array<{
        __typename?: "CpuUsageMetricsDto";
        id?: string | null;
        usage: number;
        timestamp: number;
    }>;
};

export type GetMemoryUsageMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMemoryUsageMetricsSSubscription = {
    __typename?: "Subscription";
    getMemoryUsageMetrics: {
        __typename?: "MemoryUsageMetricsDto";
        id?: string | null;
        usage: number;
        timestamp: number;
    };
};

export type GetMemoryUsageMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetMemoryUsageMetricsQQuery = {
    __typename?: "Query";
    getMemoryUsageMetrics: Array<{
        __typename?: "MemoryUsageMetricsDto";
        id?: string | null;
        usage: number;
        timestamp: number;
    }>;
};

export type GetIoMetricsSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetIoMetricsSSubscription = {
    __typename?: "Subscription";
    getIOMetrics: {
        __typename?: "IOMetricsDto";
        id?: string | null;
        input: number;
        output: number;
        timestamp: number;
    };
};

export type GetIoMetricsQQueryVariables = Exact<{
    filter: Graph_Filter;
}>;

export type GetIoMetricsQQuery = {
    __typename?: "Query";
    getIOMetrics: Array<{
        __typename?: "IOMetricsDto";
        id?: string | null;
        input: number;
        output: number;
        timestamp: number;
    }>;
};

export type GetMetricsCpuTrxQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsCpuTrxQuery = {
    __typename?: "Query";
    getMetricsCpuTRX: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsUptimeQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsUptimeQuery = {
    __typename?: "Query";
    getMetricsUptime: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsMemoryTrxQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsMemoryTrxQuery = {
    __typename?: "Query";
    getMetricsMemoryTRX: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

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
export function useGetDataUsageQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetDataUsageQuery,
        GetDataUsageQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetDataUsageQuery, GetDataUsageQueryVariables>(
        GetDataUsageDocument,
        options
    );
}
export function useGetDataUsageLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetDataUsageQuery,
        GetDataUsageQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetDataUsageQuery, GetDataUsageQueryVariables>(
        GetDataUsageDocument,
        options
    );
}
export type GetDataUsageQueryHookResult = ReturnType<
    typeof useGetDataUsageQuery
>;
export type GetDataUsageLazyQueryHookResult = ReturnType<
    typeof useGetDataUsageLazyQuery
>;
export type GetDataUsageQueryResult = Apollo.QueryResult<
    GetDataUsageQuery,
    GetDataUsageQueryVariables
>;
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
export function useGetLatestDataUsageSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetLatestDataUsageSubscription,
        GetLatestDataUsageSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetLatestDataUsageSubscription,
        GetLatestDataUsageSubscriptionVariables
    >(GetLatestDataUsageDocument, options);
}
export type GetLatestDataUsageSubscriptionHookResult = ReturnType<
    typeof useGetLatestDataUsageSubscription
>;
export type GetLatestDataUsageSubscriptionResult =
    Apollo.SubscriptionResult<GetLatestDataUsageSubscription>;
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
export function useGetConnectedUsersQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetConnectedUsersQuery,
        GetConnectedUsersQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetConnectedUsersQuery,
        GetConnectedUsersQueryVariables
    >(GetConnectedUsersDocument, options);
}
export function useGetConnectedUsersLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetConnectedUsersQuery,
        GetConnectedUsersQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetConnectedUsersQuery,
        GetConnectedUsersQueryVariables
    >(GetConnectedUsersDocument, options);
}
export type GetConnectedUsersQueryHookResult = ReturnType<
    typeof useGetConnectedUsersQuery
>;
export type GetConnectedUsersLazyQueryHookResult = ReturnType<
    typeof useGetConnectedUsersLazyQuery
>;
export type GetConnectedUsersQueryResult = Apollo.QueryResult<
    GetConnectedUsersQuery,
    GetConnectedUsersQueryVariables
>;
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
export function useGetLatestConnectedUsersSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetLatestConnectedUsersSubscription,
        GetLatestConnectedUsersSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetLatestConnectedUsersSubscription,
        GetLatestConnectedUsersSubscriptionVariables
    >(GetLatestConnectedUsersDocument, options);
}
export type GetLatestConnectedUsersSubscriptionHookResult = ReturnType<
    typeof useGetLatestConnectedUsersSubscription
>;
export type GetLatestConnectedUsersSubscriptionResult =
    Apollo.SubscriptionResult<GetLatestConnectedUsersSubscription>;
export const GetDataBillDocument = gql`
    query getDataBill($filter: DATA_BILL_FILTER!) {
        getDataBill(filter: $filter) {
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
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetDataBillQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetDataBillQuery,
        GetDataBillQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetDataBillQuery, GetDataBillQueryVariables>(
        GetDataBillDocument,
        options
    );
}
export function useGetDataBillLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetDataBillQuery,
        GetDataBillQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetDataBillQuery, GetDataBillQueryVariables>(
        GetDataBillDocument,
        options
    );
}
export type GetDataBillQueryHookResult = ReturnType<typeof useGetDataBillQuery>;
export type GetDataBillLazyQueryHookResult = ReturnType<
    typeof useGetDataBillLazyQuery
>;
export type GetDataBillQueryResult = Apollo.QueryResult<
    GetDataBillQuery,
    GetDataBillQueryVariables
>;
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
export function useGetLatestDataBillSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetLatestDataBillSubscription,
        GetLatestDataBillSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetLatestDataBillSubscription,
        GetLatestDataBillSubscriptionVariables
    >(GetLatestDataBillDocument, options);
}
export type GetLatestDataBillSubscriptionHookResult = ReturnType<
    typeof useGetLatestDataBillSubscription
>;
export type GetLatestDataBillSubscriptionResult =
    Apollo.SubscriptionResult<GetLatestDataBillSubscription>;
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
export function useGetAlertsQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetAlertsQuery,
        GetAlertsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetAlertsQuery, GetAlertsQueryVariables>(
        GetAlertsDocument,
        options
    );
}
export function useGetAlertsLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetAlertsQuery,
        GetAlertsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetAlertsQuery, GetAlertsQueryVariables>(
        GetAlertsDocument,
        options
    );
}
export type GetAlertsQueryHookResult = ReturnType<typeof useGetAlertsQuery>;
export type GetAlertsLazyQueryHookResult = ReturnType<
    typeof useGetAlertsLazyQuery
>;
export type GetAlertsQueryResult = Apollo.QueryResult<
    GetAlertsQuery,
    GetAlertsQueryVariables
>;
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
export function useGetLatestAlertsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetLatestAlertsSubscription,
        GetLatestAlertsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetLatestAlertsSubscription,
        GetLatestAlertsSubscriptionVariables
    >(GetLatestAlertsDocument, options);
}
export type GetLatestAlertsSubscriptionHookResult = ReturnType<
    typeof useGetLatestAlertsSubscription
>;
export type GetLatestAlertsSubscriptionResult =
    Apollo.SubscriptionResult<GetLatestAlertsSubscription>;
export const GetNodesByOrgDocument = gql`
    query getNodesByOrg($orgId: String!) {
        getNodesByOrg(orgId: $orgId) {
            orgName
            nodes {
                id
                status
                title
                type
                description
                totalUser
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
 *      orgId: // value for 'orgId'
 *   },
 * });
 */
export function useGetNodesByOrgQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetNodesByOrgQuery,
        GetNodesByOrgQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>(
        GetNodesByOrgDocument,
        options
    );
}
export function useGetNodesByOrgLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetNodesByOrgQuery,
        GetNodesByOrgQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetNodesByOrgQuery, GetNodesByOrgQueryVariables>(
        GetNodesByOrgDocument,
        options
    );
}
export type GetNodesByOrgQueryHookResult = ReturnType<
    typeof useGetNodesByOrgQuery
>;
export type GetNodesByOrgLazyQueryHookResult = ReturnType<
    typeof useGetNodesByOrgLazyQuery
>;
export type GetNodesByOrgQueryResult = Apollo.QueryResult<
    GetNodesByOrgQuery,
    GetNodesByOrgQueryVariables
>;
export const GetNodeDetailsDocument = gql`
    query getNodeDetails {
        getNodeDetails {
            id
            modelType
            serial
            macAddress
            osVersion
            manufacturing
            ukamaOS
            hardware
            description
        }
    }
`;

/**
 * __useGetNodeDetailsQuery__
 *
 * To run a query within a React component, call `useGetNodeDetailsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeDetailsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeDetailsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNodeDetailsQuery(
    baseOptions?: Apollo.QueryHookOptions<
        GetNodeDetailsQuery,
        GetNodeDetailsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetNodeDetailsQuery, GetNodeDetailsQueryVariables>(
        GetNodeDetailsDocument,
        options
    );
}
export function useGetNodeDetailsLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetNodeDetailsQuery,
        GetNodeDetailsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetNodeDetailsQuery,
        GetNodeDetailsQueryVariables
    >(GetNodeDetailsDocument, options);
}
export type GetNodeDetailsQueryHookResult = ReturnType<
    typeof useGetNodeDetailsQuery
>;
export type GetNodeDetailsLazyQueryHookResult = ReturnType<
    typeof useGetNodeDetailsLazyQuery
>;
export type GetNodeDetailsQueryResult = Apollo.QueryResult<
    GetNodeDetailsQuery,
    GetNodeDetailsQueryVariables
>;
export const MyUsersDocument = gql`
    query myUsers($orgId: String!) {
        myUsers(orgId: $orgId) {
            orgName
            users {
                id
                name
                email
                eSimNumber
                dataPlan
                dataUsage
                phone
                roaming
                iccid
                status
            }
        }
    }
`;

/**
 * __useMyUsersQuery__
 *
 * To run a query within a React component, call `useMyUsersQuery` and pass it any options that fit your needs.
 * When your component renders, `useMyUsersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMyUsersQuery({
 *   variables: {
 *      orgId: // value for 'orgId'
 *   },
 * });
 */
export function useMyUsersQuery(
    baseOptions: Apollo.QueryHookOptions<MyUsersQuery, MyUsersQueryVariables>
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<MyUsersQuery, MyUsersQueryVariables>(
        MyUsersDocument,
        options
    );
}
export function useMyUsersLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        MyUsersQuery,
        MyUsersQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<MyUsersQuery, MyUsersQueryVariables>(
        MyUsersDocument,
        options
    );
}
export type MyUsersQueryHookResult = ReturnType<typeof useMyUsersQuery>;
export type MyUsersLazyQueryHookResult = ReturnType<typeof useMyUsersLazyQuery>;
export type MyUsersQueryResult = Apollo.QueryResult<
    MyUsersQuery,
    MyUsersQueryVariables
>;
export const GetUserDocument = gql`
    query getUser($id: String!) {
        getUser(id: $id) {
            id
            status
            name
            eSimNumber
            iccid
            email
            phone
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
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetUserQuery(
    baseOptions: Apollo.QueryHookOptions<GetUserQuery, GetUserQueryVariables>
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetUserQuery, GetUserQueryVariables>(
        GetUserDocument,
        options
    );
}
export function useGetUserLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetUserQuery,
        GetUserQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetUserQuery, GetUserQueryVariables>(
        GetUserDocument,
        options
    );
}
export type GetUserQueryHookResult = ReturnType<typeof useGetUserQuery>;
export type GetUserLazyQueryHookResult = ReturnType<typeof useGetUserLazyQuery>;
export type GetUserQueryResult = Apollo.QueryResult<
    GetUserQuery,
    GetUserQueryVariables
>;
export const GetResidentsDocument = gql`
    query getResidents($data: PaginationDto!) {
        getResidents(data: $data) {
            meta {
                count
                page
                size
                pages
            }
            residents {
                residents {
                    id
                    name
                    dataUsage
                }
                activeResidents
                totalResidents
            }
        }
    }
`;

/**
 * __useGetResidentsQuery__
 *
 * To run a query within a React component, call `useGetResidentsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetResidentsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetResidentsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetResidentsQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetResidentsQuery,
        GetResidentsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetResidentsQuery, GetResidentsQueryVariables>(
        GetResidentsDocument,
        options
    );
}
export function useGetResidentsLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetResidentsQuery,
        GetResidentsQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetResidentsQuery, GetResidentsQueryVariables>(
        GetResidentsDocument,
        options
    );
}
export type GetResidentsQueryHookResult = ReturnType<
    typeof useGetResidentsQuery
>;
export type GetResidentsLazyQueryHookResult = ReturnType<
    typeof useGetResidentsLazyQuery
>;
export type GetResidentsQueryResult = Apollo.QueryResult<
    GetResidentsQuery,
    GetResidentsQueryVariables
>;
export const GetNetworkDocument = gql`
    query getNetwork($filter: NETWORK_TYPE!) {
        getNetwork(filter: $filter) {
            id
            status
            description
        }
    }
`;

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
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetNetworkQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetNetworkQuery,
        GetNetworkQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetNetworkQuery, GetNetworkQueryVariables>(
        GetNetworkDocument,
        options
    );
}
export function useGetNetworkLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetNetworkQuery,
        GetNetworkQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetNetworkQuery, GetNetworkQueryVariables>(
        GetNetworkDocument,
        options
    );
}
export type GetNetworkQueryHookResult = ReturnType<typeof useGetNetworkQuery>;
export type GetNetworkLazyQueryHookResult = ReturnType<
    typeof useGetNetworkLazyQuery
>;
export type GetNetworkQueryResult = Apollo.QueryResult<
    GetNetworkQuery,
    GetNetworkQueryVariables
>;
export const GetLatestNetworkDocument = gql`
    subscription getLatestNetwork {
        getNetwork {
            id
            status
            description
        }
    }
`;

/**
 * __useGetLatestNetworkSubscription__
 *
 * To run a query within a React component, call `useGetLatestNetworkSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestNetworkSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestNetworkSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetLatestNetworkSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetLatestNetworkSubscription,
        GetLatestNetworkSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetLatestNetworkSubscription,
        GetLatestNetworkSubscriptionVariables
    >(GetLatestNetworkDocument, options);
}
export type GetLatestNetworkSubscriptionHookResult = ReturnType<
    typeof useGetLatestNetworkSubscription
>;
export type GetLatestNetworkSubscriptionResult =
    Apollo.SubscriptionResult<GetLatestNetworkSubscription>;
export const DeactivateUserDocument = gql`
    mutation deactivateUser($id: String!) {
        deactivateUser(id: $id) {
            id
            success
        }
    }
`;
export type DeactivateUserMutationFn = Apollo.MutationFunction<
    DeactivateUserMutation,
    DeactivateUserMutationVariables
>;

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
export function useDeactivateUserMutation(
    baseOptions?: Apollo.MutationHookOptions<
        DeactivateUserMutation,
        DeactivateUserMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<
        DeactivateUserMutation,
        DeactivateUserMutationVariables
    >(DeactivateUserDocument, options);
}
export type DeactivateUserMutationHookResult = ReturnType<
    typeof useDeactivateUserMutation
>;
export type DeactivateUserMutationResult =
    Apollo.MutationResult<DeactivateUserMutation>;
export type DeactivateUserMutationOptions = Apollo.BaseMutationOptions<
    DeactivateUserMutation,
    DeactivateUserMutationVariables
>;
export const UpdateUserDocument = gql`
    mutation updateUser($data: UpdateUserDto!) {
        updateUser(data: $data) {
            id
            name
            sim
            email
            phone
        }
    }
`;
export type UpdateUserMutationFn = Apollo.MutationFunction<
    UpdateUserMutation,
    UpdateUserMutationVariables
>;

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
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateUserMutation(
    baseOptions?: Apollo.MutationHookOptions<
        UpdateUserMutation,
        UpdateUserMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<UpdateUserMutation, UpdateUserMutationVariables>(
        UpdateUserDocument,
        options
    );
}
export type UpdateUserMutationHookResult = ReturnType<
    typeof useUpdateUserMutation
>;
export type UpdateUserMutationResult =
    Apollo.MutationResult<UpdateUserMutation>;
export type UpdateUserMutationOptions = Apollo.BaseMutationOptions<
    UpdateUserMutation,
    UpdateUserMutationVariables
>;
export const ActivateUserDocument = gql`
    mutation activateUser($data: ActivateUserDto!) {
        activateUser(data: $data) {
            success
        }
    }
`;
export type ActivateUserMutationFn = Apollo.MutationFunction<
    ActivateUserMutation,
    ActivateUserMutationVariables
>;

/**
 * __useActivateUserMutation__
 *
 * To run a mutation, you first call `useActivateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useActivateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [activateUserMutation, { data, loading, error }] = useActivateUserMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useActivateUserMutation(
    baseOptions?: Apollo.MutationHookOptions<
        ActivateUserMutation,
        ActivateUserMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<
        ActivateUserMutation,
        ActivateUserMutationVariables
    >(ActivateUserDocument, options);
}
export type ActivateUserMutationHookResult = ReturnType<
    typeof useActivateUserMutation
>;
export type ActivateUserMutationResult =
    Apollo.MutationResult<ActivateUserMutation>;
export type ActivateUserMutationOptions = Apollo.BaseMutationOptions<
    ActivateUserMutation,
    ActivateUserMutationVariables
>;
export const UpdateNodeDocument = gql`
    mutation updateNode($data: UpdateNodeDto!) {
        updateNode(data: $data) {
            id
            name
            serialNo
            __typename
        }
    }
`;
export type UpdateNodeMutationFn = Apollo.MutationFunction<
    UpdateNodeMutation,
    UpdateNodeMutationVariables
>;

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
export function useUpdateNodeMutation(
    baseOptions?: Apollo.MutationHookOptions<
        UpdateNodeMutation,
        UpdateNodeMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<UpdateNodeMutation, UpdateNodeMutationVariables>(
        UpdateNodeDocument,
        options
    );
}
export type UpdateNodeMutationHookResult = ReturnType<
    typeof useUpdateNodeMutation
>;
export type UpdateNodeMutationResult =
    Apollo.MutationResult<UpdateNodeMutation>;
export type UpdateNodeMutationOptions = Apollo.BaseMutationOptions<
    UpdateNodeMutation,
    UpdateNodeMutationVariables
>;
export const AddNodeDocument = gql`
    mutation addNode($data: AddNodeDto!) {
        addNode(data: $data) {
            success
            __typename
        }
    }
`;
export type AddNodeMutationFn = Apollo.MutationFunction<
    AddNodeMutation,
    AddNodeMutationVariables
>;

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
export function useAddNodeMutation(
    baseOptions?: Apollo.MutationHookOptions<
        AddNodeMutation,
        AddNodeMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<AddNodeMutation, AddNodeMutationVariables>(
        AddNodeDocument,
        options
    );
}
export type AddNodeMutationHookResult = ReturnType<typeof useAddNodeMutation>;
export type AddNodeMutationResult = Apollo.MutationResult<AddNodeMutation>;
export type AddNodeMutationOptions = Apollo.BaseMutationOptions<
    AddNodeMutation,
    AddNodeMutationVariables
>;
export const DeleteNodeDocument = gql`
    mutation deleteNode($id: String!) {
        deleteNode(id: $id) {
            id
            success
        }
    }
`;
export type DeleteNodeMutationFn = Apollo.MutationFunction<
    DeleteNodeMutation,
    DeleteNodeMutationVariables
>;

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
export function useDeleteNodeMutation(
    baseOptions?: Apollo.MutationHookOptions<
        DeleteNodeMutation,
        DeleteNodeMutationVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useMutation<DeleteNodeMutation, DeleteNodeMutationVariables>(
        DeleteNodeDocument,
        options
    );
}
export type DeleteNodeMutationHookResult = ReturnType<
    typeof useDeleteNodeMutation
>;
export type DeleteNodeMutationResult =
    Apollo.MutationResult<DeleteNodeMutation>;
export type DeleteNodeMutationOptions = Apollo.BaseMutationOptions<
    DeleteNodeMutation,
    DeleteNodeMutationVariables
>;
export const GetNodeRfkpisDocument = gql`
    subscription getNodeRFKPIS {
        getNodeRFKPI {
            qam
            rfOutput
            rssi
            timestamp
        }
    }
`;

/**
 * __useGetNodeRfkpisSubscription__
 *
 * To run a query within a React component, call `useGetNodeRfkpisSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeRfkpisSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeRfkpisSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetNodeRfkpisSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetNodeRfkpisSubscription,
        GetNodeRfkpisSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetNodeRfkpisSubscription,
        GetNodeRfkpisSubscriptionVariables
    >(GetNodeRfkpisDocument, options);
}
export type GetNodeRfkpisSubscriptionHookResult = ReturnType<
    typeof useGetNodeRfkpisSubscription
>;
export type GetNodeRfkpisSubscriptionResult =
    Apollo.SubscriptionResult<GetNodeRfkpisSubscription>;
export const GetNodeRfkpiqDocument = gql`
    query getNodeRFKPIQ($filter: GRAPH_FILTER!) {
        getNodeRFKPI(filter: $filter) {
            qam
            rfOutput
            rssi
            timestamp
        }
    }
`;

/**
 * __useGetNodeRfkpiqQuery__
 *
 * To run a query within a React component, call `useGetNodeRfkpiqQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeRfkpiqQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeRfkpiqQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetNodeRfkpiqQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetNodeRfkpiqQuery,
        GetNodeRfkpiqQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetNodeRfkpiqQuery, GetNodeRfkpiqQueryVariables>(
        GetNodeRfkpiqDocument,
        options
    );
}
export function useGetNodeRfkpiqLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetNodeRfkpiqQuery,
        GetNodeRfkpiqQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetNodeRfkpiqQuery, GetNodeRfkpiqQueryVariables>(
        GetNodeRfkpiqDocument,
        options
    );
}
export type GetNodeRfkpiqQueryHookResult = ReturnType<
    typeof useGetNodeRfkpiqQuery
>;
export type GetNodeRfkpiqLazyQueryHookResult = ReturnType<
    typeof useGetNodeRfkpiqLazyQuery
>;
export type GetNodeRfkpiqQueryResult = Apollo.QueryResult<
    GetNodeRfkpiqQuery,
    GetNodeRfkpiqQueryVariables
>;
export const GetUsersAttachedMetricsSDocument = gql`
    subscription getUsersAttachedMetricsS {
        getUsersAttachedMetrics {
            id
            users
            timestamp
        }
    }
`;

/**
 * __useGetUsersAttachedMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetUsersAttachedMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersAttachedMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersAttachedMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetUsersAttachedMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetUsersAttachedMetricsSSubscription,
        GetUsersAttachedMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetUsersAttachedMetricsSSubscription,
        GetUsersAttachedMetricsSSubscriptionVariables
    >(GetUsersAttachedMetricsSDocument, options);
}
export type GetUsersAttachedMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetUsersAttachedMetricsSSubscription
>;
export type GetUsersAttachedMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetUsersAttachedMetricsSSubscription>;
export const GetUsersAttachedMetricsQDocument = gql`
    query getUsersAttachedMetricsQ($filter: GRAPH_FILTER!) {
        getUsersAttachedMetrics(filter: $filter) {
            id
            users
            timestamp
        }
    }
`;

/**
 * __useGetUsersAttachedMetricsQQuery__
 *
 * To run a query within a React component, call `useGetUsersAttachedMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersAttachedMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersAttachedMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetUsersAttachedMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetUsersAttachedMetricsQQuery,
        GetUsersAttachedMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetUsersAttachedMetricsQQuery,
        GetUsersAttachedMetricsQQueryVariables
    >(GetUsersAttachedMetricsQDocument, options);
}
export function useGetUsersAttachedMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetUsersAttachedMetricsQQuery,
        GetUsersAttachedMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetUsersAttachedMetricsQQuery,
        GetUsersAttachedMetricsQQueryVariables
    >(GetUsersAttachedMetricsQDocument, options);
}
export type GetUsersAttachedMetricsQQueryHookResult = ReturnType<
    typeof useGetUsersAttachedMetricsQQuery
>;
export type GetUsersAttachedMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetUsersAttachedMetricsQLazyQuery
>;
export type GetUsersAttachedMetricsQQueryResult = Apollo.QueryResult<
    GetUsersAttachedMetricsQQuery,
    GetUsersAttachedMetricsQQueryVariables
>;
export const GetThroughputMetricsSDocument = gql`
    subscription getThroughputMetricsS {
        getThroughputMetrics {
            id
            amount
            timestamp
        }
    }
`;

/**
 * __useGetThroughputMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetThroughputMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetThroughputMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetThroughputMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetThroughputMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetThroughputMetricsSSubscription,
        GetThroughputMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetThroughputMetricsSSubscription,
        GetThroughputMetricsSSubscriptionVariables
    >(GetThroughputMetricsSDocument, options);
}
export type GetThroughputMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetThroughputMetricsSSubscription
>;
export type GetThroughputMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetThroughputMetricsSSubscription>;
export const GetThroughputMetricsQDocument = gql`
    query getThroughputMetricsQ($filter: GRAPH_FILTER!) {
        getThroughputMetrics(filter: $filter) {
            id
            amount
            timestamp
        }
    }
`;

/**
 * __useGetThroughputMetricsQQuery__
 *
 * To run a query within a React component, call `useGetThroughputMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetThroughputMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetThroughputMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetThroughputMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetThroughputMetricsQQuery,
        GetThroughputMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetThroughputMetricsQQuery,
        GetThroughputMetricsQQueryVariables
    >(GetThroughputMetricsQDocument, options);
}
export function useGetThroughputMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetThroughputMetricsQQuery,
        GetThroughputMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetThroughputMetricsQQuery,
        GetThroughputMetricsQQueryVariables
    >(GetThroughputMetricsQDocument, options);
}
export type GetThroughputMetricsQQueryHookResult = ReturnType<
    typeof useGetThroughputMetricsQQuery
>;
export type GetThroughputMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetThroughputMetricsQLazyQuery
>;
export type GetThroughputMetricsQQueryResult = Apollo.QueryResult<
    GetThroughputMetricsQQuery,
    GetThroughputMetricsQQueryVariables
>;
export const GetTemperatureMetricsSDocument = gql`
    subscription getTemperatureMetricsS {
        getTemperatureMetrics {
            id
            temperature
            timestamp
        }
    }
`;

/**
 * __useGetTemperatureMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetTemperatureMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetTemperatureMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTemperatureMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetTemperatureMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetTemperatureMetricsSSubscription,
        GetTemperatureMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetTemperatureMetricsSSubscription,
        GetTemperatureMetricsSSubscriptionVariables
    >(GetTemperatureMetricsSDocument, options);
}
export type GetTemperatureMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetTemperatureMetricsSSubscription
>;
export type GetTemperatureMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetTemperatureMetricsSSubscription>;
export const GetTemperatureMetricsQDocument = gql`
    query getTemperatureMetricsQ($filter: GRAPH_FILTER!) {
        getTemperatureMetrics(filter: $filter) {
            id
            temperature
            timestamp
        }
    }
`;

/**
 * __useGetTemperatureMetricsQQuery__
 *
 * To run a query within a React component, call `useGetTemperatureMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTemperatureMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTemperatureMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetTemperatureMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetTemperatureMetricsQQuery,
        GetTemperatureMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetTemperatureMetricsQQuery,
        GetTemperatureMetricsQQueryVariables
    >(GetTemperatureMetricsQDocument, options);
}
export function useGetTemperatureMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetTemperatureMetricsQQuery,
        GetTemperatureMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetTemperatureMetricsQQuery,
        GetTemperatureMetricsQQueryVariables
    >(GetTemperatureMetricsQDocument, options);
}
export type GetTemperatureMetricsQQueryHookResult = ReturnType<
    typeof useGetTemperatureMetricsQQuery
>;
export type GetTemperatureMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetTemperatureMetricsQLazyQuery
>;
export type GetTemperatureMetricsQQueryResult = Apollo.QueryResult<
    GetTemperatureMetricsQQuery,
    GetTemperatureMetricsQQueryVariables
>;
export const GetCpuUsageMetricsSDocument = gql`
    subscription getCpuUsageMetricsS {
        getCpuUsageMetrics {
            id
            usage
            timestamp
        }
    }
`;

/**
 * __useGetCpuUsageMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetCpuUsageMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetCpuUsageMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCpuUsageMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetCpuUsageMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetCpuUsageMetricsSSubscription,
        GetCpuUsageMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetCpuUsageMetricsSSubscription,
        GetCpuUsageMetricsSSubscriptionVariables
    >(GetCpuUsageMetricsSDocument, options);
}
export type GetCpuUsageMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetCpuUsageMetricsSSubscription
>;
export type GetCpuUsageMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetCpuUsageMetricsSSubscription>;
export const GetCpuUsageMetricsQDocument = gql`
    query getCpuUsageMetricsQ($filter: GRAPH_FILTER!) {
        getCpuUsageMetrics(filter: $filter) {
            id
            usage
            timestamp
        }
    }
`;

/**
 * __useGetCpuUsageMetricsQQuery__
 *
 * To run a query within a React component, call `useGetCpuUsageMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCpuUsageMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCpuUsageMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetCpuUsageMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetCpuUsageMetricsQQuery,
        GetCpuUsageMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetCpuUsageMetricsQQuery,
        GetCpuUsageMetricsQQueryVariables
    >(GetCpuUsageMetricsQDocument, options);
}
export function useGetCpuUsageMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetCpuUsageMetricsQQuery,
        GetCpuUsageMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetCpuUsageMetricsQQuery,
        GetCpuUsageMetricsQQueryVariables
    >(GetCpuUsageMetricsQDocument, options);
}
export type GetCpuUsageMetricsQQueryHookResult = ReturnType<
    typeof useGetCpuUsageMetricsQQuery
>;
export type GetCpuUsageMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetCpuUsageMetricsQLazyQuery
>;
export type GetCpuUsageMetricsQQueryResult = Apollo.QueryResult<
    GetCpuUsageMetricsQQuery,
    GetCpuUsageMetricsQQueryVariables
>;
export const GetMemoryUsageMetricsSDocument = gql`
    subscription getMemoryUsageMetricsS {
        getMemoryUsageMetrics {
            id
            usage
            timestamp
        }
    }
`;

/**
 * __useGetMemoryUsageMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetMemoryUsageMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMemoryUsageMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemoryUsageMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMemoryUsageMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMemoryUsageMetricsSSubscription,
        GetMemoryUsageMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMemoryUsageMetricsSSubscription,
        GetMemoryUsageMetricsSSubscriptionVariables
    >(GetMemoryUsageMetricsSDocument, options);
}
export type GetMemoryUsageMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetMemoryUsageMetricsSSubscription
>;
export type GetMemoryUsageMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetMemoryUsageMetricsSSubscription>;
export const GetMemoryUsageMetricsQDocument = gql`
    query getMemoryUsageMetricsQ($filter: GRAPH_FILTER!) {
        getMemoryUsageMetrics(filter: $filter) {
            id
            usage
            timestamp
        }
    }
`;

/**
 * __useGetMemoryUsageMetricsQQuery__
 *
 * To run a query within a React component, call `useGetMemoryUsageMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMemoryUsageMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemoryUsageMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetMemoryUsageMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMemoryUsageMetricsQQuery,
        GetMemoryUsageMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMemoryUsageMetricsQQuery,
        GetMemoryUsageMetricsQQueryVariables
    >(GetMemoryUsageMetricsQDocument, options);
}
export function useGetMemoryUsageMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMemoryUsageMetricsQQuery,
        GetMemoryUsageMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMemoryUsageMetricsQQuery,
        GetMemoryUsageMetricsQQueryVariables
    >(GetMemoryUsageMetricsQDocument, options);
}
export type GetMemoryUsageMetricsQQueryHookResult = ReturnType<
    typeof useGetMemoryUsageMetricsQQuery
>;
export type GetMemoryUsageMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetMemoryUsageMetricsQLazyQuery
>;
export type GetMemoryUsageMetricsQQueryResult = Apollo.QueryResult<
    GetMemoryUsageMetricsQQuery,
    GetMemoryUsageMetricsQQueryVariables
>;
export const GetIoMetricsSDocument = gql`
    subscription getIOMetricsS {
        getIOMetrics {
            id
            input
            output
            timestamp
        }
    }
`;

/**
 * __useGetIoMetricsSSubscription__
 *
 * To run a query within a React component, call `useGetIoMetricsSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetIoMetricsSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetIoMetricsSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetIoMetricsSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetIoMetricsSSubscription,
        GetIoMetricsSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetIoMetricsSSubscription,
        GetIoMetricsSSubscriptionVariables
    >(GetIoMetricsSDocument, options);
}
export type GetIoMetricsSSubscriptionHookResult = ReturnType<
    typeof useGetIoMetricsSSubscription
>;
export type GetIoMetricsSSubscriptionResult =
    Apollo.SubscriptionResult<GetIoMetricsSSubscription>;
export const GetIoMetricsQDocument = gql`
    query getIOMetricsQ($filter: GRAPH_FILTER!) {
        getIOMetrics(filter: $filter) {
            id
            input
            output
            timestamp
        }
    }
`;

/**
 * __useGetIoMetricsQQuery__
 *
 * To run a query within a React component, call `useGetIoMetricsQQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetIoMetricsQQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetIoMetricsQQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useGetIoMetricsQQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetIoMetricsQQuery,
        GetIoMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetIoMetricsQQuery, GetIoMetricsQQueryVariables>(
        GetIoMetricsQDocument,
        options
    );
}
export function useGetIoMetricsQLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetIoMetricsQQuery,
        GetIoMetricsQQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetIoMetricsQQuery, GetIoMetricsQQueryVariables>(
        GetIoMetricsQDocument,
        options
    );
}
export type GetIoMetricsQQueryHookResult = ReturnType<
    typeof useGetIoMetricsQQuery
>;
export type GetIoMetricsQLazyQueryHookResult = ReturnType<
    typeof useGetIoMetricsQLazyQuery
>;
export type GetIoMetricsQQueryResult = Apollo.QueryResult<
    GetIoMetricsQQuery,
    GetIoMetricsQQueryVariables
>;
export const GetMetricsCpuTrxDocument = gql`
    query getMetricsCpuTRX($data: MetricsInputDTO!) {
        getMetricsCpuTRX(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsCpuTrxQuery__
 *
 * To run a query within a React component, call `useGetMetricsCpuTrxQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsCpuTrxQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsCpuTrxQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsCpuTrxQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsCpuTrxQuery,
        GetMetricsCpuTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsCpuTrxQuery,
        GetMetricsCpuTrxQueryVariables
    >(GetMetricsCpuTrxDocument, options);
}
export function useGetMetricsCpuTrxLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsCpuTrxQuery,
        GetMetricsCpuTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsCpuTrxQuery,
        GetMetricsCpuTrxQueryVariables
    >(GetMetricsCpuTrxDocument, options);
}
export type GetMetricsCpuTrxQueryHookResult = ReturnType<
    typeof useGetMetricsCpuTrxQuery
>;
export type GetMetricsCpuTrxLazyQueryHookResult = ReturnType<
    typeof useGetMetricsCpuTrxLazyQuery
>;
export type GetMetricsCpuTrxQueryResult = Apollo.QueryResult<
    GetMetricsCpuTrxQuery,
    GetMetricsCpuTrxQueryVariables
>;
export const GetMetricsUptimeDocument = gql`
    query getMetricsUptime($data: MetricsInputDTO!) {
        getMetricsUptime(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsUptimeQuery__
 *
 * To run a query within a React component, call `useGetMetricsUptimeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsUptimeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsUptimeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsUptimeQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsUptimeQuery,
        GetMetricsUptimeQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsUptimeQuery,
        GetMetricsUptimeQueryVariables
    >(GetMetricsUptimeDocument, options);
}
export function useGetMetricsUptimeLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsUptimeQuery,
        GetMetricsUptimeQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsUptimeQuery,
        GetMetricsUptimeQueryVariables
    >(GetMetricsUptimeDocument, options);
}
export type GetMetricsUptimeQueryHookResult = ReturnType<
    typeof useGetMetricsUptimeQuery
>;
export type GetMetricsUptimeLazyQueryHookResult = ReturnType<
    typeof useGetMetricsUptimeLazyQuery
>;
export type GetMetricsUptimeQueryResult = Apollo.QueryResult<
    GetMetricsUptimeQuery,
    GetMetricsUptimeQueryVariables
>;
export const GetMetricsMemoryTrxDocument = gql`
    query getMetricsMemoryTRX($data: MetricsInputDTO!) {
        getMetricsMemoryTRX(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsMemoryTrxQuery__
 *
 * To run a query within a React component, call `useGetMetricsMemoryTrxQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsMemoryTrxQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsMemoryTrxQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsMemoryTrxQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsMemoryTrxQuery,
        GetMetricsMemoryTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsMemoryTrxQuery,
        GetMetricsMemoryTrxQueryVariables
    >(GetMetricsMemoryTrxDocument, options);
}
export function useGetMetricsMemoryTrxLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsMemoryTrxQuery,
        GetMetricsMemoryTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsMemoryTrxQuery,
        GetMetricsMemoryTrxQueryVariables
    >(GetMetricsMemoryTrxDocument, options);
}
export type GetMetricsMemoryTrxQueryHookResult = ReturnType<
    typeof useGetMetricsMemoryTrxQuery
>;
export type GetMetricsMemoryTrxLazyQueryHookResult = ReturnType<
    typeof useGetMetricsMemoryTrxLazyQuery
>;
export type GetMetricsMemoryTrxQueryResult = Apollo.QueryResult<
    GetMetricsMemoryTrxQuery,
    GetMetricsMemoryTrxQueryVariables
>;
