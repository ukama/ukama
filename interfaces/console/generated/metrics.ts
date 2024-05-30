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
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  scopes: Array<Scalars['String']['input']>;
  siteId: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
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
  Network = 'NETWORK',
  Node = 'NODE',
  Org = 'ORG',
  Site = 'SITE',
  Subscriber = 'SUBSCRIBER',
  User = 'USER'
}

export enum Notification_Type {
  Error = 'ERROR',
  Info = 'INFO',
  Unknown = 'UNKNOWN',
  Warning = 'WARNING'
}

export type NotificationRes = {
  __typename?: 'NotificationRes';
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  role: Role_Type;
  scope: Notification_Scope;
  subscriberId: Scalars['String']['output'];
  timeStamp: Scalars['String']['output'];
  title: Scalars['String']['output'];
  type: Notification_Type;
  userId: Scalars['String']['output'];
};

export type NotificationsRes = {
  __typename?: 'NotificationsRes';
  notifications: Array<NotificationRes>;
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
  Admin = 'ADMIN',
  Owner = 'OWNER',
  Users = 'USERS',
  Vendor = 'VENDOR'
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
  getNotificationsSub: NotificationRes;
};


export type SubscriptionGetMetricByTabSubArgs = {
  from: Scalars['Float']['input'];
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
};


export type SubscriptionGetNotificationsSubArgs = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  scopes: Array<Scalars['String']['input']>;
  siteId: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
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


export type GetNotificationsQuery = { __typename?: 'Query', getNotifications: { __typename?: 'NotificationsRes', notifications: Array<{ __typename?: 'NotificationRes', description: string, id: string, isRead: boolean, networkId: string, orgId: string, role: Role_Type, scope: Notification_Scope, subscriberId: string, title: string, type: Notification_Type, userId: string, timeStamp: string }> } };

export type GetNotificationsSubSubscriptionVariables = Exact<{
  orgId: Scalars['String']['input'];
  networkId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  scopes: Array<Scalars['String']['input']> | Scalars['String']['input'];
}>;


export type GetNotificationsSubSubscription = { __typename?: 'Subscription', getNotificationsSub: { __typename?: 'NotificationRes', description: string, id: string, isRead: boolean, networkId: string, orgId: string, role: Role_Type, scope: Notification_Scope, subscriberId: string, timeStamp: string, title: string, type: Notification_Type, userId: string } };


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
    query getNotifications($data: GetNotificationsInput!) {
  getNotifications(data: $data) {
    notifications {
      description
      id
      isRead
      networkId
      orgId
      role
      scope
      subscriberId
      title
      type
      userId
      timeStamp
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
export const GetNotificationsSubDocument = gql`
    subscription GetNotificationsSub($orgId: String!, $networkId: String!, $siteId: String!, $userId: String!, $subscriberId: String!, $scopes: [String!]!) {
  getNotificationsSub(
    orgId: $orgId
    networkId: $networkId
    siteId: $siteId
    userId: $userId
    subscriberId: $subscriberId
    scopes: $scopes
  ) {
    description
    id
    isRead
    networkId
    orgId
    role
    scope
    subscriberId
    timeStamp
    title
    type
    userId
  }
}
    `;

/**
 * __useGetNotificationsSubSubscription__
 *
 * To run a query within a React component, call `useGetNotificationsSubSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetNotificationsSubSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNotificationsSubSubscription({
 *   variables: {
 *      orgId: // value for 'orgId'
 *      networkId: // value for 'networkId'
 *      siteId: // value for 'siteId'
 *      userId: // value for 'userId'
 *      subscriberId: // value for 'subscriberId'
 *      scopes: // value for 'scopes'
 *   },
 * });
 */
export function useGetNotificationsSubSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetNotificationsSubSubscription, GetNotificationsSubSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetNotificationsSubSubscription, GetNotificationsSubSubscriptionVariables>(GetNotificationsSubDocument, options);
      }
export type GetNotificationsSubSubscriptionHookResult = ReturnType<typeof useGetNotificationsSubSubscription>;
export type GetNotificationsSubSubscriptionResult = Apollo.SubscriptionResult<GetNotificationsSubSubscription>;