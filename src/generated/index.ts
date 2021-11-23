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
const defaultOptions = {};
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
    eSimNumber: Scalars["String"];
    email?: InputMaybe<Scalars["String"]>;
    firstName: Scalars["String"];
    lastName: Scalars["String"];
    phone?: InputMaybe<Scalars["String"]>;
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

export enum Connected_User_Type {
    Guests = "GUESTS",
    Residents = "RESIDENTS",
}

export type ConnectedUserDto = {
    __typename?: "ConnectedUserDto";
    guestUsers: Scalars["Float"];
    residentUsers: Scalars["Float"];
    totalUser: Scalars["Float"];
};

export type ConnectedUserResponse = {
    __typename?: "ConnectedUserResponse";
    data: ConnectedUserDto;
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

export enum Data_Plan_Type {
    Na = "NA",
    Paid = "PAID",
    Unpaid = "UNPAID",
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
    dataPlan: Data_Plan_Type;
    dataUsage: Scalars["Float"];
    dlActivity: Scalars["String"];
    id: Scalars["String"];
    name: Scalars["String"];
    node: Scalars["String"];
    status: Get_User_Status_Type;
    ulActivity: Scalars["String"];
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

export type Meta = {
    __typename?: "Meta";
    count: Scalars["Float"];
    page: Scalars["Float"];
    pages: Scalars["Float"];
    size: Scalars["Float"];
};

export type Mutation = {
    __typename?: "Mutation";
    activateUser: ActivateUserResponse;
    addNode: AddNodeResponse;
};

export type MutationActivateUserArgs = {
    data: ActivateUserDto;
};

export type MutationAddNodeArgs = {
    data: AddNodeDto;
};

export type NodeDto = {
    __typename?: "NodeDto";
    description: Scalars["String"];
    id: Scalars["String"];
    status: Get_User_Status_Type;
    title: Scalars["String"];
    totalUser: Scalars["Float"];
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
    getConnectedUsers: ConnectedUserDto;
    getDataBill: DataBillDto;
    getDataUsage: DataUsageDto;
    getEsims: Array<EsimDto>;
    getNodes: NodesResponse;
    getResidents: ResidentsResponse;
    getUsers: GetUserResponse;
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

export type QueryGetNodesArgs = {
    data: PaginationDto;
};

export type QueryGetResidentsArgs = {
    data: PaginationDto;
};

export type QueryGetUsersArgs = {
    data: GetUserPaginationDto;
};

export type ResidentDto = {
    __typename?: "ResidentDto";
    dataUsage: Scalars["Float"];
    id: Scalars["String"];
    name: Scalars["String"];
};

export type ResidentResponse = {
    __typename?: "ResidentResponse";
    activeResidents: Scalars["Float"];
    residents: Array<ResidentDto>;
    totalResidents: Scalars["Float"];
};

export type ResidentsResponse = {
    __typename?: "ResidentsResponse";
    meta: Meta;
    residents: ResidentResponse;
};

export enum Time_Filter {
    Month = "MONTH",
    Today = "TODAY",
    Total = "TOTAL",
    Week = "WEEK",
}

export type UserDto = {
    __typename?: "UserDto";
    email: Scalars["String"];
    id: Scalars["String"];
    name: Scalars["String"];
    type: Connected_User_Type;
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

export type GetConnectedUsersQueryVariables = Exact<{
    filter: Time_Filter;
}>;

export type GetConnectedUsersQuery = {
    __typename?: "Query";
    getConnectedUsers: {
        __typename?: "ConnectedUserDto";
        totalUser: number;
        residentUsers: number;
        guestUsers: number;
    };
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
            id?: string | null | undefined;
            type: Alert_Type;
            title?: string | null | undefined;
            description?: string | null | undefined;
            alertDate?: any | null | undefined;
        }>;
    };
};

export type GetNodesQueryVariables = Exact<{
    data: PaginationDto;
}>;

export type GetNodesQuery = {
    __typename?: "Query";
    getNodes: {
        __typename?: "NodesResponse";
        meta: {
            __typename?: "Meta";
            count: number;
            page: number;
            size: number;
            pages: number;
        };
        nodes: {
            __typename?: "NodeResponseDto";
            activeNodes: number;
            totalNodes: number;
            nodes: Array<{
                __typename?: "NodeDto";
                id: string;
                title: string;
                description: string;
                totalUser: number;
            }>;
        };
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
                __typename?: "ResidentDto";
                id: string;
                name: string;
                dataUsage: number;
            }>;
        };
    };
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
export const GetConnectedUsersDocument = gql`
    query getConnectedUsers($filter: TIME_FILTER!) {
        getConnectedUsers(filter: $filter) {
            totalUser
            residentUsers
            guestUsers
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
export const GetNodesDocument = gql`
    query getNodes($data: PaginationDto!) {
        getNodes(data: $data) {
            meta {
                count
                page
                size
                pages
            }
            nodes {
                nodes {
                    id
                    title
                    description
                    totalUser
                }
                activeNodes
                totalNodes
            }
        }
    }
`;

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
export function useGetNodesQuery(
    baseOptions: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables>
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(
        GetNodesDocument,
        options
    );
}
export function useGetNodesLazyQuery(
    baseOptions?: Apollo.LazyQueryHookOptions<
        GetNodesQuery,
        GetNodesQueryVariables
    >
) {
    const options = { ...defaultOptions, ...baseOptions };
    return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(
        GetNodesDocument,
        options
    );
}
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<
    typeof useGetNodesLazyQuery
>;
export type GetNodesQueryResult = Apollo.QueryResult<
    GetNodesQuery,
    GetNodesQueryVariables
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
