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
  _Any: { input: any; output: any; }
  _FieldSet: { input: any; output: any; }
};

export enum Graphs_Type {
  Network = 'NETWORK',
  NodeHealth = 'NODE_HEALTH',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Subscribers = 'SUBSCRIBERS'
}

export type GetLatestMetricInput = {
  nodeId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type GetMetricByTabInput = {
  from?: InputMaybe<Scalars['Float']['input']>;
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  step?: InputMaybe<Scalars['Float']['input']>;
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription: Scalars['Boolean']['input'];
};

export type GetNotificationsInput = {
  forRole: Role_Type;
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  scopes: Array<Notification_Scope>;
  userId: Scalars['String']['input'];
};

export type LatestMetricRes = {
  __typename?: 'LatestMetricRes';
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  values: Array<Array<Scalars['Float']['output']>>;
};

export type MetricsRes = {
  __typename?: 'MetricsRes';
  metrics: Array<MetricRes>;
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
  NotifCritical = 'NOTIF_CRITICAL',
  NotifError = 'NOTIF_ERROR',
  NotifInfo = 'NOTIF_INFO',
  NotifInvalid = 'NOTIF_INVALID',
  NotifWarning = 'NOTIF_WARNING'
}

export type NotificationsRes = {
  __typename?: 'NotificationsRes';
  notifications: Array<NotificationsResDto>;
};

export type NotificationsResDto = {
  __typename?: 'NotificationsResDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  scope: Notification_Scope;
  title: Scalars['String']['output'];
  type: Notification_Type;
};

export type Query = {
  __typename?: 'Query';
  _service: _Service;
  getLatestMetric: LatestMetricRes;
  getMetricByTab: MetricsRes;
  getNotifications: NotificationsRes;
  getStatsMetric: StatsMetric;
};


export type QueryGetLatestMetricArgs = {
  data: GetLatestMetricInput;
};


export type QueryGetMetricByTabArgs = {
  data: GetMetricByTabInput;
};


export type QueryGetNotificationsArgs = {
  data: GetNotificationsInput;
};

export enum Role_Type {
  RoleAdmin = 'ROLE_ADMIN',
  RoleInvalid = 'ROLE_INVALID',
  RoleNetworkOwner = 'ROLE_NETWORK_OWNER',
  RoleOwner = 'ROLE_OWNER',
  RoleUser = 'ROLE_USER',
  RoleVendor = 'ROLE_VENDOR'
}

export type StatsMetric = {
  __typename?: 'StatsMetric';
  activeSubscriber: Scalars['Float']['output'];
  averageSignalStrength: Scalars['Float']['output'];
  averageThroughput: Scalars['Float']['output'];
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetricByTabSub: LatestMetricRes;
  notificationSubscription: NotificationsResDto;
};


export type SubscriptionGetMetricByTabSubArgs = {
  from: Scalars['Float']['input'];
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
};


export type SubscriptionNotificationSubscriptionArgs = {
  forRole: Role_Type;
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  scopes: Array<Notification_Scope>;
  userId: Scalars['String']['input'];
};

export type _Service = {
  __typename?: '_Service';
  sdl?: Maybe<Scalars['String']['output']>;
};

export type GetLatestMetricQueryVariables = Exact<{
  data: GetLatestMetricInput;
}>;


export type GetLatestMetricQuery = { __typename?: 'Query', getLatestMetric: { __typename?: 'LatestMetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, value: Array<number> } };

export type GetMetricByTabQueryVariables = Exact<{
  data: GetMetricByTabInput;
}>;


export type GetMetricByTabQuery = { __typename?: 'Query', getMetricByTab: { __typename?: 'MetricsRes', metrics: Array<{ __typename?: 'MetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, values: Array<Array<number>> }> } };

export type GetMetricByTabSubSubscriptionVariables = Exact<{
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  from: Scalars['Float']['input'];
}>;


export type GetMetricByTabSubSubscription = { __typename?: 'Subscription', getMetricByTabSub: { __typename?: 'LatestMetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, value: Array<number> } };

export type GetStatsMetricQueryVariables = Exact<{ [key: string]: never; }>;


export type GetStatsMetricQuery = { __typename?: 'Query', getStatsMetric: { __typename?: 'StatsMetric', activeSubscriber: number, averageSignalStrength: number, averageThroughput: number } };

export type GetNotificationsQueryVariables = Exact<{
  data: GetNotificationsInput;
}>;


export type GetNotificationsQuery = { __typename?: 'Query', getNotifications: { __typename?: 'NotificationsRes', notifications: Array<{ __typename?: 'NotificationsResDto', id: string, title: string, description: string, createdAt: string, type: Notification_Type, scope: Notification_Scope, isRead: boolean }> } };

export type NotificationSubscriptionSubscriptionVariables = Exact<{
  orgId: Scalars['String']['input'];
  networkId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  forRole: Role_Type;
  scopes: Array<Notification_Scope> | Notification_Scope;
}>;


export type NotificationSubscriptionSubscription = { __typename?: 'Subscription', notificationSubscription: { __typename?: 'NotificationsResDto', id: string, title: string, description: string, createdAt: string, type: Notification_Type, scope: Notification_Scope, isRead: boolean } };


export const GetLatestMetricDocument = gql`
    query GetLatestMetric($data: GetLatestMetricInput!) {
  getLatestMetric(data: $data) {
    success
    msg
    orgId
    nodeId
    type
    value
  }
}
    `;

/**
 * __useGetLatestMetricQuery__
 *
 * To run a query within a React component, call `useGetLatestMetricQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestMetricQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestMetricQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetLatestMetricQuery(baseOptions: Apollo.QueryHookOptions<GetLatestMetricQuery, GetLatestMetricQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetLatestMetricQuery, GetLatestMetricQueryVariables>(GetLatestMetricDocument, options);
      }
export function useGetLatestMetricLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetLatestMetricQuery, GetLatestMetricQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetLatestMetricQuery, GetLatestMetricQueryVariables>(GetLatestMetricDocument, options);
        }
export type GetLatestMetricQueryHookResult = ReturnType<typeof useGetLatestMetricQuery>;
export type GetLatestMetricLazyQueryHookResult = ReturnType<typeof useGetLatestMetricLazyQuery>;
export type GetLatestMetricQueryResult = Apollo.QueryResult<GetLatestMetricQuery, GetLatestMetricQueryVariables>;
export const GetMetricByTabDocument = gql`
    query GetMetricByTab($data: GetMetricByTabInput!) {
  getMetricByTab(data: $data) {
    metrics {
      success
      msg
      orgId
      nodeId
      type
      values
    }
  }
}
    `;

/**
 * __useGetMetricByTabQuery__
 *
 * To run a query within a React component, call `useGetMetricByTabQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricByTabQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricByTabQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricByTabQuery(baseOptions: Apollo.QueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
      }
export function useGetMetricByTabLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
        }
export type GetMetricByTabQueryHookResult = ReturnType<typeof useGetMetricByTabQuery>;
export type GetMetricByTabLazyQueryHookResult = ReturnType<typeof useGetMetricByTabLazyQuery>;
export type GetMetricByTabQueryResult = Apollo.QueryResult<GetMetricByTabQuery, GetMetricByTabQueryVariables>;
export const GetMetricByTabSubDocument = gql`
    subscription GetMetricByTabSub($nodeId: String!, $orgId: String!, $type: GRAPHS_TYPE!, $userId: String!, $from: Float!) {
  getMetricByTabSub(
    nodeId: $nodeId
    orgId: $orgId
    type: $type
    userId: $userId
    from: $from
  ) {
    success
    msg
    orgId
    nodeId
    type
    value
  }
}
    `;

/**
 * __useGetMetricByTabSubSubscription__
 *
 * To run a query within a React component, call `useGetMetricByTabSubSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricByTabSubSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricByTabSubSubscription({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *      orgId: // value for 'orgId'
 *      type: // value for 'type'
 *      userId: // value for 'userId'
 *      from: // value for 'from'
 *   },
 * });
 */
export function useGetMetricByTabSubSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetMetricByTabSubSubscription, GetMetricByTabSubSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetMetricByTabSubSubscription, GetMetricByTabSubSubscriptionVariables>(GetMetricByTabSubDocument, options);
      }
export type GetMetricByTabSubSubscriptionHookResult = ReturnType<typeof useGetMetricByTabSubSubscription>;
export type GetMetricByTabSubSubscriptionResult = Apollo.SubscriptionResult<GetMetricByTabSubSubscription>;
export const GetStatsMetricDocument = gql`
    query GetStatsMetric {
  getStatsMetric {
    activeSubscriber
    averageSignalStrength
    averageThroughput
  }
}
    `;

/**
 * __useGetStatsMetricQuery__
 *
 * To run a query within a React component, call `useGetStatsMetricQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetStatsMetricQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetStatsMetricQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetStatsMetricQuery(baseOptions?: Apollo.QueryHookOptions<GetStatsMetricQuery, GetStatsMetricQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetStatsMetricQuery, GetStatsMetricQueryVariables>(GetStatsMetricDocument, options);
      }
export function useGetStatsMetricLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetStatsMetricQuery, GetStatsMetricQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetStatsMetricQuery, GetStatsMetricQueryVariables>(GetStatsMetricDocument, options);
        }
export type GetStatsMetricQueryHookResult = ReturnType<typeof useGetStatsMetricQuery>;
export type GetStatsMetricLazyQueryHookResult = ReturnType<typeof useGetStatsMetricLazyQuery>;
export type GetStatsMetricQueryResult = Apollo.QueryResult<GetStatsMetricQuery, GetStatsMetricQueryVariables>;
export const GetNotificationsDocument = gql`
    query GetNotifications($data: GetNotificationsInput!) {
  getNotifications(data: $data) {
    notifications {
      id
      title
      description
      createdAt
      type
      scope
      isRead
    }
  }
}
    `;

/**
 * __useGetNotificationsQuery__
 *
 * To run a query within a React component, call `useGetNotificationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNotificationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNotificationsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNotificationsQuery(baseOptions: Apollo.QueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
      }
export function useGetNotificationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
        }
export type GetNotificationsQueryHookResult = ReturnType<typeof useGetNotificationsQuery>;
export type GetNotificationsLazyQueryHookResult = ReturnType<typeof useGetNotificationsLazyQuery>;
export type GetNotificationsQueryResult = Apollo.QueryResult<GetNotificationsQuery, GetNotificationsQueryVariables>;
export const NotificationSubscriptionDocument = gql`
    subscription NotificationSubscription($orgId: String!, $networkId: String!, $userId: String!, $forRole: ROLE_TYPE!, $scopes: [NOTIFICATION_SCOPE!]!) {
  notificationSubscription(
    orgId: $orgId
    networkId: $networkId
    userId: $userId
    forRole: $forRole
    scopes: $scopes
  ) {
    id
    title
    description
    createdAt
    type
    scope
    isRead
  }
}
    `;

/**
 * __useNotificationSubscriptionSubscription__
 *
 * To run a query within a React component, call `useNotificationSubscriptionSubscription` and pass it any options that fit your needs.
 * When your component renders, `useNotificationSubscriptionSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNotificationSubscriptionSubscription({
 *   variables: {
 *      orgId: // value for 'orgId'
 *      networkId: // value for 'networkId'
 *      userId: // value for 'userId'
 *      forRole: // value for 'forRole'
 *      scopes: // value for 'scopes'
 *   },
 * });
 */
export function useNotificationSubscriptionSubscription(baseOptions: Apollo.SubscriptionHookOptions<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables>(NotificationSubscriptionDocument, options);
      }
export type NotificationSubscriptionSubscriptionHookResult = ReturnType<typeof useNotificationSubscriptionSubscription>;
export type NotificationSubscriptionSubscriptionResult = Apollo.SubscriptionResult<NotificationSubscriptionSubscription>;