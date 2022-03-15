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
    regPolling: Scalars["Boolean"];
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
    getCurrentBill: BillResponse;
    getDataBill: DataBillDto;
    getDataUsage: DataUsageDto;
    getEsims: Array<EsimDto>;
    getMetricsCpuCOM: Array<MetricDto>;
    getMetricsCpuTRX: Array<MetricDto>;
    getMetricsDiskCOM: Array<MetricDto>;
    getMetricsDiskTRX: Array<MetricDto>;
    getMetricsERAB: Array<MetricDto>;
    getMetricsMemoryCOM: Array<MetricDto>;
    getMetricsMemoryTRX: Array<MetricDto>;
    getMetricsPaPower: Array<MetricDto>;
    getMetricsPower: Array<MetricDto>;
    getMetricsRLC: Array<MetricDto>;
    getMetricsRRC: Array<MetricDto>;
    getMetricsRxPower: Array<MetricDto>;
    getMetricsSubActive: Array<MetricDto>;
    getMetricsSubAttached: Array<MetricDto>;
    getMetricsTempCOM: Array<MetricDto>;
    getMetricsTempTRX: Array<MetricDto>;
    getMetricsThroughputDL: Array<MetricDto>;
    getMetricsThroughputUL: Array<MetricDto>;
    getMetricsTxPower: Array<MetricDto>;
    getMetricsUptime: Array<MetricDto>;
    getNetwork: NetworkDto;
    getNodeDetails: NodeDetailDto;
    getNodeNetwork: NetworkDto;
    getNodes: NodesResponse;
    getNodesByOrg: OrgNodeResponseDto;
    getResidents: ResidentsResponse;
    getUser: GetUserDto;
    getUsers: GetUserResponse;
    getUsersByOrg: OrgUserResponseDto;
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

export type QueryGetMetricsCpuComArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsCpuTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsDiskComArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsDiskTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsErabArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsMemoryComArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsMemoryTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsPaPowerArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsPowerArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsRlcArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsRrcArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsRxPowerArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsSubActiveArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsSubAttachedArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsTempComArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsTempTrxArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsThroughputDlArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsThroughputUlArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsTxPowerArgs = {
    data: MetricsInputDto;
};

export type QueryGetMetricsUptimeArgs = {
    data: MetricsInputDto;
};

export type QueryGetNetworkArgs = {
    filter: Network_Type;
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

export type QueryGetUserArgs = {
    id: Scalars["String"];
};

export type QueryGetUsersArgs = {
    data: GetUserPaginationDto;
};

export type QueryGetUsersByOrgArgs = {
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
    getDataBill: DataBillDto;
    getDataUsage: DataUsageDto;
    getMetricsCpuCOM: Array<MetricDto>;
    getMetricsCpuTRX: Array<MetricDto>;
    getMetricsDiskCOM: Array<MetricDto>;
    getMetricsDiskTRX: Array<MetricDto>;
    getMetricsERAB: Array<MetricDto>;
    getMetricsMemoryCOM: Array<MetricDto>;
    getMetricsMemoryTRX: Array<MetricDto>;
    getMetricsPaPower: Array<MetricDto>;
    getMetricsPower: Array<MetricDto>;
    getMetricsRLC: Array<MetricDto>;
    getMetricsRRC: Array<MetricDto>;
    getMetricsRxPower: Array<MetricDto>;
    getMetricsSubActive: Array<MetricDto>;
    getMetricsSubAttached: Array<MetricDto>;
    getMetricsTempCOM: Array<MetricDto>;
    getMetricsTempTRX: Array<MetricDto>;
    getMetricsThroughputDL: Array<MetricDto>;
    getMetricsThroughputUL: Array<MetricDto>;
    getMetricsTxPower: Array<MetricDto>;
    getMetricsUptime: Array<MetricDto>;
    getNetwork: NetworkDto;
};

export enum Time_Filter {
    Month = "MONTH",
    Today = "TODAY",
    Total = "TOTAL",
    Week = "WEEK",
}

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

export type GetUsersByOrgQueryVariables = Exact<{
    orgId: Scalars["String"];
}>;

export type GetUsersByOrgQuery = {
    __typename?: "Query";
    getUsersByOrg: {
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

export type GetMetricsUptimeQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsUptimeQuery = {
    __typename?: "Query";
    getMetricsUptime: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsUptimeSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsUptimeSSubscription = {
    __typename?: "Subscription";
    getMetricsUptime: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsThroughputUlQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsThroughputUlQuery = {
    __typename?: "Query";
    getMetricsThroughputUL: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsThroughputUlsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsThroughputUlsSubscription = {
    __typename?: "Subscription";
    getMetricsThroughputUL: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsRlcQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsRlcQuery = {
    __typename?: "Query";
    getMetricsRLC: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsRlCsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsRlCsSubscription = {
    __typename?: "Subscription";
    getMetricsRLC: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsErabQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsErabQuery = {
    __typename?: "Query";
    getMetricsERAB: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsEraBsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsEraBsSubscription = {
    __typename?: "Subscription";
    getMetricsERAB: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsRrcQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsRrcQuery = {
    __typename?: "Query";
    getMetricsRRC: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsRrCsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsRrCsSubscription = {
    __typename?: "Subscription";
    getMetricsRRC: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsThroughputDlQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsThroughputDlQuery = {
    __typename?: "Query";
    getMetricsThroughputDL: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsMemoryComsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsMemoryComsSubscription = {
    __typename?: "Subscription";
    getMetricsMemoryCOM: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsMemoryComQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsMemoryComQuery = {
    __typename?: "Query";
    getMetricsMemoryCOM: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsThroughputDlsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsThroughputDlsSubscription = {
    __typename?: "Subscription";
    getMetricsThroughputDL: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
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

export type GetMetricsMemoryTrxsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsMemoryTrxsSubscription = {
    __typename?: "Subscription";
    getMetricsMemoryTRX: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsCpuTrxQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsCpuTrxQuery = {
    __typename?: "Query";
    getMetricsCpuTRX: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsCpuTrxsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsCpuTrxsSubscription = {
    __typename?: "Subscription";
    getMetricsCpuTRX: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsPowerQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsPowerQuery = {
    __typename?: "Query";
    getMetricsPower: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsPowerSSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsPowerSSubscription = {
    __typename?: "Subscription";
    getMetricsPower: Array<{ __typename?: "MetricDto"; y: number; x: number }>;
};

export type GetMetricsTempTrxQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsTempTrxQuery = {
    __typename?: "Query";
    getMetricsTempTRX: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsTempTrxsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsTempTrxsSubscription = {
    __typename?: "Subscription";
    getMetricsTempTRX: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsTempComQueryVariables = Exact<{
    data: MetricsInputDto;
}>;

export type GetMetricsTempComQuery = {
    __typename?: "Query";
    getMetricsTempCOM: Array<{
        __typename?: "MetricDto";
        y: number;
        x: number;
    }>;
};

export type GetMetricsTempComsSubscriptionVariables = Exact<{
    [key: string]: never;
}>;

export type GetMetricsTempComsSubscription = {
    __typename?: "Subscription";
    getMetricsTempCOM: Array<{
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
export const GetUsersByOrgDocument = gql`
    query getUsersByOrg($orgId: String!) {
        getUsersByOrg(orgId: $orgId) {
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
 *      orgId: // value for 'orgId'
 *   },
 * });
 */
export function useGetUsersByOrgQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetUsersByOrgQuery,
        GetUsersByOrgQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>(
        GetUsersByOrgDocument,
        options
    );
}
export function useGetUsersByOrgLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetUsersByOrgQuery,
        GetUsersByOrgQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetUsersByOrgQuery, GetUsersByOrgQueryVariables>(
        GetUsersByOrgDocument,
        options
    );
}
export type GetUsersByOrgQueryHookResult = ReturnType<
    typeof useGetUsersByOrgQuery
>;
export type GetUsersByOrgLazyQueryHookResult = ReturnType<
    typeof useGetUsersByOrgLazyQuery
>;
export type GetUsersByOrgQueryResult = Apollo.QueryResult<
    GetUsersByOrgQuery,
    GetUsersByOrgQueryVariables
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
export const GetMetricsUptimeSDocument = gql`
    subscription getMetricsUptimeS {
        getMetricsUptime {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsUptimeSSubscription__
 *
 * To run a query within a React component, call `useGetMetricsUptimeSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsUptimeSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsUptimeSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsUptimeSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsUptimeSSubscription,
        GetMetricsUptimeSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsUptimeSSubscription,
        GetMetricsUptimeSSubscriptionVariables
    >(GetMetricsUptimeSDocument, options);
}
export type GetMetricsUptimeSSubscriptionHookResult = ReturnType<
    typeof useGetMetricsUptimeSSubscription
>;
export type GetMetricsUptimeSSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsUptimeSSubscription>;
export const GetMetricsThroughputUlDocument = gql`
    query getMetricsThroughputUL($data: MetricsInputDTO!) {
        getMetricsThroughputUL(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsThroughputUlQuery__
 *
 * To run a query within a React component, call `useGetMetricsThroughputUlQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsThroughputUlQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsThroughputUlQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsThroughputUlQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsThroughputUlQuery,
        GetMetricsThroughputUlQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsThroughputUlQuery,
        GetMetricsThroughputUlQueryVariables
    >(GetMetricsThroughputUlDocument, options);
}
export function useGetMetricsThroughputUlLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsThroughputUlQuery,
        GetMetricsThroughputUlQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsThroughputUlQuery,
        GetMetricsThroughputUlQueryVariables
    >(GetMetricsThroughputUlDocument, options);
}
export type GetMetricsThroughputUlQueryHookResult = ReturnType<
    typeof useGetMetricsThroughputUlQuery
>;
export type GetMetricsThroughputUlLazyQueryHookResult = ReturnType<
    typeof useGetMetricsThroughputUlLazyQuery
>;
export type GetMetricsThroughputUlQueryResult = Apollo.QueryResult<
    GetMetricsThroughputUlQuery,
    GetMetricsThroughputUlQueryVariables
>;
export const GetMetricsThroughputUlsDocument = gql`
    subscription getMetricsThroughputULS {
        getMetricsThroughputUL {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsThroughputUlsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsThroughputUlsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsThroughputUlsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsThroughputUlsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsThroughputUlsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsThroughputUlsSubscription,
        GetMetricsThroughputUlsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsThroughputUlsSubscription,
        GetMetricsThroughputUlsSubscriptionVariables
    >(GetMetricsThroughputUlsDocument, options);
}
export type GetMetricsThroughputUlsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsThroughputUlsSubscription
>;
export type GetMetricsThroughputUlsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsThroughputUlsSubscription>;
export const GetMetricsRlcDocument = gql`
    query getMetricsRLC($data: MetricsInputDTO!) {
        getMetricsRLC(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsRlcQuery__
 *
 * To run a query within a React component, call `useGetMetricsRlcQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsRlcQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsRlcQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsRlcQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsRlcQuery,
        GetMetricsRlcQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetMetricsRlcQuery, GetMetricsRlcQueryVariables>(
        GetMetricsRlcDocument,
        options
    );
}
export function useGetMetricsRlcLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsRlcQuery,
        GetMetricsRlcQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetMetricsRlcQuery, GetMetricsRlcQueryVariables>(
        GetMetricsRlcDocument,
        options
    );
}
export type GetMetricsRlcQueryHookResult = ReturnType<
    typeof useGetMetricsRlcQuery
>;
export type GetMetricsRlcLazyQueryHookResult = ReturnType<
    typeof useGetMetricsRlcLazyQuery
>;
export type GetMetricsRlcQueryResult = Apollo.QueryResult<
    GetMetricsRlcQuery,
    GetMetricsRlcQueryVariables
>;
export const GetMetricsRlCsDocument = gql`
    subscription getMetricsRLCs {
        getMetricsRLC {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsRlCsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsRlCsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsRlCsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsRlCsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsRlCsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsRlCsSubscription,
        GetMetricsRlCsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsRlCsSubscription,
        GetMetricsRlCsSubscriptionVariables
    >(GetMetricsRlCsDocument, options);
}
export type GetMetricsRlCsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsRlCsSubscription
>;
export type GetMetricsRlCsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsRlCsSubscription>;
export const GetMetricsErabDocument = gql`
    query getMetricsERAB($data: MetricsInputDTO!) {
        getMetricsERAB(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsErabQuery__
 *
 * To run a query within a React component, call `useGetMetricsErabQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsErabQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsErabQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsErabQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsErabQuery,
        GetMetricsErabQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetMetricsErabQuery, GetMetricsErabQueryVariables>(
        GetMetricsErabDocument,
        options
    );
}
export function useGetMetricsErabLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsErabQuery,
        GetMetricsErabQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsErabQuery,
        GetMetricsErabQueryVariables
    >(GetMetricsErabDocument, options);
}
export type GetMetricsErabQueryHookResult = ReturnType<
    typeof useGetMetricsErabQuery
>;
export type GetMetricsErabLazyQueryHookResult = ReturnType<
    typeof useGetMetricsErabLazyQuery
>;
export type GetMetricsErabQueryResult = Apollo.QueryResult<
    GetMetricsErabQuery,
    GetMetricsErabQueryVariables
>;
export const GetMetricsEraBsDocument = gql`
    subscription getMetricsERABs {
        getMetricsERAB {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsEraBsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsEraBsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsEraBsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsEraBsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsEraBsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsEraBsSubscription,
        GetMetricsEraBsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsEraBsSubscription,
        GetMetricsEraBsSubscriptionVariables
    >(GetMetricsEraBsDocument, options);
}
export type GetMetricsEraBsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsEraBsSubscription
>;
export type GetMetricsEraBsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsEraBsSubscription>;
export const GetMetricsRrcDocument = gql`
    query getMetricsRRC($data: MetricsInputDTO!) {
        getMetricsRRC(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsRrcQuery__
 *
 * To run a query within a React component, call `useGetMetricsRrcQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsRrcQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsRrcQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsRrcQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsRrcQuery,
        GetMetricsRrcQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetMetricsRrcQuery, GetMetricsRrcQueryVariables>(
        GetMetricsRrcDocument,
        options
    );
}
export function useGetMetricsRrcLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsRrcQuery,
        GetMetricsRrcQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetMetricsRrcQuery, GetMetricsRrcQueryVariables>(
        GetMetricsRrcDocument,
        options
    );
}
export type GetMetricsRrcQueryHookResult = ReturnType<
    typeof useGetMetricsRrcQuery
>;
export type GetMetricsRrcLazyQueryHookResult = ReturnType<
    typeof useGetMetricsRrcLazyQuery
>;
export type GetMetricsRrcQueryResult = Apollo.QueryResult<
    GetMetricsRrcQuery,
    GetMetricsRrcQueryVariables
>;
export const GetMetricsRrCsDocument = gql`
    subscription getMetricsRRCs {
        getMetricsRRC {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsRrCsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsRrCsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsRrCsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsRrCsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsRrCsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsRrCsSubscription,
        GetMetricsRrCsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsRrCsSubscription,
        GetMetricsRrCsSubscriptionVariables
    >(GetMetricsRrCsDocument, options);
}
export type GetMetricsRrCsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsRrCsSubscription
>;
export type GetMetricsRrCsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsRrCsSubscription>;
export const GetMetricsThroughputDlDocument = gql`
    query getMetricsThroughputDL($data: MetricsInputDTO!) {
        getMetricsThroughputDL(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsThroughputDlQuery__
 *
 * To run a query within a React component, call `useGetMetricsThroughputDlQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsThroughputDlQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsThroughputDlQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsThroughputDlQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsThroughputDlQuery,
        GetMetricsThroughputDlQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsThroughputDlQuery,
        GetMetricsThroughputDlQueryVariables
    >(GetMetricsThroughputDlDocument, options);
}
export function useGetMetricsThroughputDlLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsThroughputDlQuery,
        GetMetricsThroughputDlQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsThroughputDlQuery,
        GetMetricsThroughputDlQueryVariables
    >(GetMetricsThroughputDlDocument, options);
}
export type GetMetricsThroughputDlQueryHookResult = ReturnType<
    typeof useGetMetricsThroughputDlQuery
>;
export type GetMetricsThroughputDlLazyQueryHookResult = ReturnType<
    typeof useGetMetricsThroughputDlLazyQuery
>;
export type GetMetricsThroughputDlQueryResult = Apollo.QueryResult<
    GetMetricsThroughputDlQuery,
    GetMetricsThroughputDlQueryVariables
>;
export const GetMetricsMemoryComsDocument = gql`
    subscription getMetricsMemoryCOMS {
        getMetricsMemoryCOM {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsMemoryComsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsMemoryComsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsMemoryComsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsMemoryComsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsMemoryComsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsMemoryComsSubscription,
        GetMetricsMemoryComsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsMemoryComsSubscription,
        GetMetricsMemoryComsSubscriptionVariables
    >(GetMetricsMemoryComsDocument, options);
}
export type GetMetricsMemoryComsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsMemoryComsSubscription
>;
export type GetMetricsMemoryComsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsMemoryComsSubscription>;
export const GetMetricsMemoryComDocument = gql`
    query getMetricsMemoryCOM($data: MetricsInputDTO!) {
        getMetricsMemoryCOM(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsMemoryComQuery__
 *
 * To run a query within a React component, call `useGetMetricsMemoryComQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsMemoryComQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsMemoryComQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsMemoryComQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsMemoryComQuery,
        GetMetricsMemoryComQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsMemoryComQuery,
        GetMetricsMemoryComQueryVariables
    >(GetMetricsMemoryComDocument, options);
}
export function useGetMetricsMemoryComLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsMemoryComQuery,
        GetMetricsMemoryComQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsMemoryComQuery,
        GetMetricsMemoryComQueryVariables
    >(GetMetricsMemoryComDocument, options);
}
export type GetMetricsMemoryComQueryHookResult = ReturnType<
    typeof useGetMetricsMemoryComQuery
>;
export type GetMetricsMemoryComLazyQueryHookResult = ReturnType<
    typeof useGetMetricsMemoryComLazyQuery
>;
export type GetMetricsMemoryComQueryResult = Apollo.QueryResult<
    GetMetricsMemoryComQuery,
    GetMetricsMemoryComQueryVariables
>;
export const GetMetricsThroughputDlsDocument = gql`
    subscription getMetricsThroughputDLS {
        getMetricsThroughputDL {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsThroughputDlsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsThroughputDlsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsThroughputDlsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsThroughputDlsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsThroughputDlsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsThroughputDlsSubscription,
        GetMetricsThroughputDlsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsThroughputDlsSubscription,
        GetMetricsThroughputDlsSubscriptionVariables
    >(GetMetricsThroughputDlsDocument, options);
}
export type GetMetricsThroughputDlsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsThroughputDlsSubscription
>;
export type GetMetricsThroughputDlsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsThroughputDlsSubscription>;
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
export const GetMetricsMemoryTrxsDocument = gql`
    subscription getMetricsMemoryTRXS {
        getMetricsMemoryTRX {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsMemoryTrxsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsMemoryTrxsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsMemoryTrxsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsMemoryTrxsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsMemoryTrxsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsMemoryTrxsSubscription,
        GetMetricsMemoryTrxsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsMemoryTrxsSubscription,
        GetMetricsMemoryTrxsSubscriptionVariables
    >(GetMetricsMemoryTrxsDocument, options);
}
export type GetMetricsMemoryTrxsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsMemoryTrxsSubscription
>;
export type GetMetricsMemoryTrxsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsMemoryTrxsSubscription>;
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
export const GetMetricsCpuTrxsDocument = gql`
    subscription getMetricsCpuTRXS {
        getMetricsCpuTRX {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsCpuTrxsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsCpuTrxsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsCpuTrxsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsCpuTrxsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsCpuTrxsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsCpuTrxsSubscription,
        GetMetricsCpuTrxsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsCpuTrxsSubscription,
        GetMetricsCpuTrxsSubscriptionVariables
    >(GetMetricsCpuTrxsDocument, options);
}
export type GetMetricsCpuTrxsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsCpuTrxsSubscription
>;
export type GetMetricsCpuTrxsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsCpuTrxsSubscription>;
export const GetMetricsPowerDocument = gql`
    query getMetricsPower($data: MetricsInputDTO!) {
        getMetricsPower(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsPowerQuery__
 *
 * To run a query within a React component, call `useGetMetricsPowerQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsPowerQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsPowerQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsPowerQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsPowerQuery,
        GetMetricsPowerQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetMetricsPowerQuery, GetMetricsPowerQueryVariables>(
        GetMetricsPowerDocument,
        options
    );
}
export function useGetMetricsPowerLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsPowerQuery,
        GetMetricsPowerQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsPowerQuery,
        GetMetricsPowerQueryVariables
    >(GetMetricsPowerDocument, options);
}
export type GetMetricsPowerQueryHookResult = ReturnType<
    typeof useGetMetricsPowerQuery
>;
export type GetMetricsPowerLazyQueryHookResult = ReturnType<
    typeof useGetMetricsPowerLazyQuery
>;
export type GetMetricsPowerQueryResult = Apollo.QueryResult<
    GetMetricsPowerQuery,
    GetMetricsPowerQueryVariables
>;
export const GetMetricsPowerSDocument = gql`
    subscription getMetricsPowerS {
        getMetricsPower {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsPowerSSubscription__
 *
 * To run a query within a React component, call `useGetMetricsPowerSSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsPowerSSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsPowerSSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsPowerSSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsPowerSSubscription,
        GetMetricsPowerSSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsPowerSSubscription,
        GetMetricsPowerSSubscriptionVariables
    >(GetMetricsPowerSDocument, options);
}
export type GetMetricsPowerSSubscriptionHookResult = ReturnType<
    typeof useGetMetricsPowerSSubscription
>;
export type GetMetricsPowerSSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsPowerSSubscription>;
export const GetMetricsTempTrxDocument = gql`
    query getMetricsTempTRX($data: MetricsInputDTO!) {
        getMetricsTempTRX(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsTempTrxQuery__
 *
 * To run a query within a React component, call `useGetMetricsTempTrxQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsTempTrxQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsTempTrxQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsTempTrxQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsTempTrxQuery,
        GetMetricsTempTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsTempTrxQuery,
        GetMetricsTempTrxQueryVariables
    >(GetMetricsTempTrxDocument, options);
}
export function useGetMetricsTempTrxLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsTempTrxQuery,
        GetMetricsTempTrxQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsTempTrxQuery,
        GetMetricsTempTrxQueryVariables
    >(GetMetricsTempTrxDocument, options);
}
export type GetMetricsTempTrxQueryHookResult = ReturnType<
    typeof useGetMetricsTempTrxQuery
>;
export type GetMetricsTempTrxLazyQueryHookResult = ReturnType<
    typeof useGetMetricsTempTrxLazyQuery
>;
export type GetMetricsTempTrxQueryResult = Apollo.QueryResult<
    GetMetricsTempTrxQuery,
    GetMetricsTempTrxQueryVariables
>;
export const GetMetricsTempTrxsDocument = gql`
    subscription getMetricsTempTRXS {
        getMetricsTempTRX {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsTempTrxsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsTempTrxsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsTempTrxsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsTempTrxsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsTempTrxsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsTempTrxsSubscription,
        GetMetricsTempTrxsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsTempTrxsSubscription,
        GetMetricsTempTrxsSubscriptionVariables
    >(GetMetricsTempTrxsDocument, options);
}
export type GetMetricsTempTrxsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsTempTrxsSubscription
>;
export type GetMetricsTempTrxsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsTempTrxsSubscription>;
export const GetMetricsTempComDocument = gql`
    query getMetricsTempCOM($data: MetricsInputDTO!) {
        getMetricsTempCOM(data: $data) {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsTempComQuery__
 *
 * To run a query within a React component, call `useGetMetricsTempComQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsTempComQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsTempComQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsTempComQuery(
    baseOptions: Apollo.QueryHookOptions<
        GetMetricsTempComQuery,
        GetMetricsTempComQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<
        GetMetricsTempComQuery,
        GetMetricsTempComQueryVariables
    >(GetMetricsTempComDocument, options);
}
export function useGetMetricsTempComLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetMetricsTempComQuery,
        GetMetricsTempComQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<
        GetMetricsTempComQuery,
        GetMetricsTempComQueryVariables
    >(GetMetricsTempComDocument, options);
}
export type GetMetricsTempComQueryHookResult = ReturnType<
    typeof useGetMetricsTempComQuery
>;
export type GetMetricsTempComLazyQueryHookResult = ReturnType<
    typeof useGetMetricsTempComLazyQuery
>;
export type GetMetricsTempComQueryResult = Apollo.QueryResult<
    GetMetricsTempComQuery,
    GetMetricsTempComQueryVariables
>;
export const GetMetricsTempComsDocument = gql`
    subscription getMetricsTempCOMS {
        getMetricsTempCOM {
            y
            x
        }
    }
`;

/**
 * __useGetMetricsTempComsSubscription__
 *
 * To run a query within a React component, call `useGetMetricsTempComsSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsTempComsSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsTempComsSubscription({
 *   variables: {
 *   },
 * });
 */
export function useGetMetricsTempComsSubscription(
    baseOptions?: Apollo.SubscriptionHookOptions<
        GetMetricsTempComsSubscription,
        GetMetricsTempComsSubscriptionVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useSubscription<
        GetMetricsTempComsSubscription,
        GetMetricsTempComsSubscriptionVariables
    >(GetMetricsTempComsDocument, options);
}
export type GetMetricsTempComsSubscriptionHookResult = ReturnType<
    typeof useGetMetricsTempComsSubscription
>;
export type GetMetricsTempComsSubscriptionResult =
    Apollo.SubscriptionResult<GetMetricsTempComsSubscription>;
