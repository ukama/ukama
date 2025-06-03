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

export enum Graphs_Type {
  Battery = 'BATTERY',
  Controller = 'CONTROLLER',
  DataUsage = 'DATA_USAGE',
  Home = 'HOME',
  MainBackhaul = 'MAIN_BACKHAUL',
  NetworkBackhaul = 'NETWORK_BACKHAUL',
  NetworkCellular = 'NETWORK_CELLULAR',
  NodeHealth = 'NODE_HEALTH',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Site = 'SITE',
  Solar = 'SOLAR',
  Subscribers = 'SUBSCRIBERS',
  Switch = 'SWITCH'
}

export type GetMetricBySiteInput = {
  from: Scalars['Float']['input'];
  orgName: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricByTabInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  siteId?: InputMaybe<Scalars['String']['input']>;
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricsSiteStatInput = {
  from: Scalars['Float']['input'];
  nodeIds?: InputMaybe<Array<Scalars['String']['input']>>;
  orgName: Scalars['String']['input'];
  siteIds?: InputMaybe<Array<Scalars['String']['input']>>;
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Stats_Type;
  userId?: InputMaybe<Scalars['String']['input']>;
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricsStatInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  operation?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  siteId?: InputMaybe<Scalars['String']['input']>;
  step?: Scalars['Float']['input'];
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Stats_Type;
  userId?: InputMaybe<Scalars['String']['input']>;
  withSubscription?: Scalars['Boolean']['input'];
};

export type LatestMetricSubRes = {
  __typename?: 'LatestMetricSubRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId: Scalars['String']['output'];
  packageId?: Maybe<Scalars['String']['output']>;
  siteId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  values: Array<Array<Scalars['Float']['output']>>;
};

export type MetricStateRes = {
  __typename?: 'MetricStateRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId: Scalars['String']['output'];
  packageId?: Maybe<Scalars['String']['output']>;
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Scalars['Float']['output'];
};

export type MetricsRes = {
  __typename?: 'MetricsRes';
  metrics: Array<MetricRes>;
};

export type MetricsStateRes = {
  __typename?: 'MetricsStateRes';
  metrics: Array<MetricStateRes>;
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

export type NotificationRedirect = {
  __typename?: 'NotificationRedirect';
  action: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type NotificationsRes = {
  __typename?: 'NotificationsRes';
  notifications: Array<NotificationsResDto>;
};

export type NotificationsResDto = {
  __typename?: 'NotificationsResDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  eventKey: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  redirect: NotificationRedirect;
  resourceId: Scalars['String']['output'];
  scope: Notification_Scope;
  title: Scalars['String']['output'];
  type: Notification_Type;
};

export type Query = {
  __typename?: 'Query';
  getMetricBySite: MetricsRes;
  getMetricByTab: MetricsRes;
  getMetricsStat: MetricsStateRes;
  getNotifications: NotificationsRes;
  getSiteStat: SiteMetricsStateRes;
};


export type QueryGetMetricBySiteArgs = {
  data: GetMetricBySiteInput;
};


export type QueryGetMetricByTabArgs = {
  data: GetMetricByTabInput;
};


export type QueryGetMetricsStatArgs = {
  data: GetMetricsStatInput;
};


export type QueryGetNotificationsArgs = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};


export type QueryGetSiteStatArgs = {
  data: GetMetricsSiteStatInput;
};

export enum Stats_Type {
  AllNode = 'ALL_NODE',
  Battery = 'BATTERY',
  DataUsage = 'DATA_USAGE',
  Home = 'HOME',
  MainBackhaul = 'MAIN_BACKHAUL',
  Network = 'NETWORK',
  Overview = 'OVERVIEW',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Site = 'SITE'
}

export type SiteMetricStateRes = {
  __typename?: 'SiteMetricStateRes';
  msg: Scalars['String']['output'];
  nodeId?: Maybe<Scalars['String']['output']>;
  siteId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Scalars['Float']['output'];
};

export type SiteMetricsStateRes = {
  __typename?: 'SiteMetricsStateRes';
  metrics: Array<SiteMetricStateRes>;
};

export type SubMetricsStatInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  type: Stats_Type;
  userId: Scalars['String']['input'];
};

export type SubSiteMetricByTabInput = {
  from: Scalars['Float']['input'];
  orgName: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
};

export type SubSiteMetricsStatInput = {
  from: Scalars['Float']['input'];
  nodeIds?: InputMaybe<Array<Scalars['String']['input']>>;
  orgName: Scalars['String']['input'];
  siteIds?: InputMaybe<Array<Scalars['String']['input']>>;
  type: Stats_Type;
  userId: Scalars['String']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetricStatSub: LatestMetricSubRes;
  getSiteMetricByTabSub: LatestMetricSubRes;
  getSiteMetricStatSub: LatestMetricSubRes;
  notificationSubscription: NotificationsResDto;
};


export type SubscriptionGetMetricStatSubArgs = {
  data: SubMetricsStatInput;
};


export type SubscriptionGetSiteMetricByTabSubArgs = {
  data: SubSiteMetricByTabInput;
};


export type SubscriptionGetSiteMetricStatSubArgs = {
  data: SubSiteMetricsStatInput;
};


export type SubscriptionNotificationSubscriptionArgs = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type GetNotificationsQueryVariables = Exact<{
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
}>;


export type GetNotificationsQuery = { __typename?: 'Query', getNotifications: { __typename?: 'NotificationsRes', notifications: Array<{ __typename?: 'NotificationsResDto', id: string, type: Notification_Type, scope: Notification_Scope, title: string, isRead: boolean, eventKey: string, createdAt: string, resourceId: string, description: string, redirect: { __typename?: 'NotificationRedirect', action: string, title: string } }> } };

export type NotificationSubscriptionSubscriptionVariables = Exact<{
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
}>;


export type NotificationSubscriptionSubscription = { __typename?: 'Subscription', notificationSubscription: { __typename?: 'NotificationsResDto', id: string, type: Notification_Type, scope: Notification_Scope, title: string, isRead: boolean, eventKey: string, createdAt: string, resourceId: string, description: string, redirect: { __typename?: 'NotificationRedirect', action: string, title: string } } };

export type GetMetricByTabQueryVariables = Exact<{
  data: GetMetricByTabInput;
}>;


export type GetMetricByTabQuery = { __typename?: 'Query', getMetricByTab: { __typename?: 'MetricsRes', metrics: Array<{ __typename?: 'MetricRes', msg: string, nodeId?: string | null, success: boolean, type: string, values: Array<Array<number>>, networkId?: string | null, siteId?: string | null, packageId?: string | null, dataPlanId?: string | null }> } };

export type GetMetricsStatQueryVariables = Exact<{
  data: GetMetricsStatInput;
}>;


export type GetMetricsStatQuery = { __typename?: 'Query', getMetricsStat: { __typename?: 'MetricsStateRes', metrics: Array<{ __typename?: 'MetricStateRes', success: boolean, msg: string, nodeId: string, type: string, value: number, networkId?: string | null, packageId?: string | null, dataPlanId?: string | null }> } };

export type GetSiteStatQueryVariables = Exact<{
  data: GetMetricsSiteStatInput;
}>;


export type GetSiteStatQuery = { __typename?: 'Query', getSiteStat: { __typename?: 'SiteMetricsStateRes', metrics: Array<{ __typename?: 'SiteMetricStateRes', success: boolean, msg: string, siteId: string, nodeId?: string | null, type: string, value: number }> } };

export type GetMetricBySiteQueryVariables = Exact<{
  data: GetMetricBySiteInput;
}>;


export type GetMetricBySiteQuery = { __typename?: 'Query', getMetricBySite: { __typename?: 'MetricsRes', metrics: Array<{ __typename?: 'MetricRes', msg: string, success: boolean, type: string, values: Array<Array<number>> }> } };

export type GetMetricsStatSubSubscriptionVariables = Exact<{
  data: SubMetricsStatInput;
}>;


export type GetMetricsStatSubSubscription = { __typename?: 'Subscription', getMetricStatSub: { __typename?: 'LatestMetricSubRes', msg: string, nodeId: string, success: boolean, type: string, value: Array<number>, networkId?: string | null, siteId: string, packageId?: string | null, dataPlanId?: string | null } };

export type GetSiteMetricStatSubSubscriptionVariables = Exact<{
  data: SubSiteMetricsStatInput;
}>;


export type GetSiteMetricStatSubSubscription = { __typename?: 'Subscription', getSiteMetricStatSub: { __typename?: 'LatestMetricSubRes', msg: string, siteId: string, nodeId: string, success: boolean, type: string, value: Array<number> } };


export const GetNotificationsDocument = gql`
    query GetNotifications($networkId: String!, $orgId: String!, $orgName: String!, $role: String!, $startTimestamp: String!, $subscriberId: String!, $userId: String!) {
  getNotifications(
    networkId: $networkId
    orgId: $orgId
    orgName: $orgName
    startTimestamp: $startTimestamp
    subscriberId: $subscriberId
    userId: $userId
    role: $role
  ) {
    notifications {
      id
      type
      scope
      title
      isRead
      eventKey
      createdAt
      resourceId
      description
      redirect {
        action
        title
      }
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
 *      networkId: // value for 'networkId'
 *      orgId: // value for 'orgId'
 *      orgName: // value for 'orgName'
 *      role: // value for 'role'
 *      startTimestamp: // value for 'startTimestamp'
 *      subscriberId: // value for 'subscriberId'
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetNotificationsQuery(baseOptions: Apollo.QueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables> & ({ variables: GetNotificationsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
}
export function useGetNotificationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useLazyQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
}
export function useGetNotificationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
  const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions};
  return Apollo.useSuspenseQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
}
export type GetNotificationsQueryHookResult = ReturnType<typeof useGetNotificationsQuery>;
export type GetNotificationsLazyQueryHookResult = ReturnType<typeof useGetNotificationsLazyQuery>;
export type GetNotificationsSuspenseQueryHookResult = ReturnType<typeof useGetNotificationsSuspenseQuery>;
export type GetNotificationsQueryResult = Apollo.QueryResult<GetNotificationsQuery, GetNotificationsQueryVariables>;
export const NotificationSubscriptionDocument = gql`
    subscription NotificationSubscription($networkId: String!, $orgId: String!, $orgName: String!, $role: String!, $startTimestamp: String!, $subscriberId: String!, $userId: String!) {
  notificationSubscription(
    networkId: $networkId
    orgId: $orgId
    orgName: $orgName
    startTimestamp: $startTimestamp
    subscriberId: $subscriberId
    userId: $userId
    role: $role
  ) {
    id
    type
    scope
    title
    isRead
    eventKey
    createdAt
    resourceId
    description
    redirect {
      action
      title
    }
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
 *      networkId: // value for 'networkId'
 *      orgId: // value for 'orgId'
 *      orgName: // value for 'orgName'
 *      role: // value for 'role'
 *      startTimestamp: // value for 'startTimestamp'
 *      subscriberId: // value for 'subscriberId'
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useNotificationSubscriptionSubscription(baseOptions: Apollo.SubscriptionHookOptions<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables> & ({ variables: NotificationSubscriptionSubscriptionVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useSubscription<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables>(NotificationSubscriptionDocument, options);
}
export type NotificationSubscriptionSubscriptionHookResult = ReturnType<typeof useNotificationSubscriptionSubscription>;
export type NotificationSubscriptionSubscriptionResult = Apollo.SubscriptionResult<NotificationSubscriptionSubscription>;
export const GetMetricByTabDocument = gql`
    query GetMetricByTab($data: GetMetricByTabInput!) {
  getMetricByTab(data: $data) {
    metrics {
      msg
      nodeId
      success
      type
      values
      networkId
      siteId
      packageId
      dataPlanId
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
export function useGetMetricByTabQuery(baseOptions: Apollo.QueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables> & ({ variables: GetMetricByTabQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
}
export function useGetMetricByTabLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useLazyQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
}
export function useGetMetricByTabSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
  const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions};
  return Apollo.useSuspenseQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
}
export type GetMetricByTabQueryHookResult = ReturnType<typeof useGetMetricByTabQuery>;
export type GetMetricByTabLazyQueryHookResult = ReturnType<typeof useGetMetricByTabLazyQuery>;
export type GetMetricByTabSuspenseQueryHookResult = ReturnType<typeof useGetMetricByTabSuspenseQuery>;
export type GetMetricByTabQueryResult = Apollo.QueryResult<GetMetricByTabQuery, GetMetricByTabQueryVariables>;
export const GetMetricsStatDocument = gql`
    query GetMetricsStat($data: GetMetricsStatInput!) {
  getMetricsStat(data: $data) {
    metrics {
      success
      msg
      nodeId
      type
      value
      networkId
      packageId
      dataPlanId
    }
  }
}
    `;

/**
 * __useGetMetricsStatQuery__
 *
 * To run a query within a React component, call `useGetMetricsStatQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsStatQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsStatQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsStatQuery(baseOptions: Apollo.QueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables> & ({ variables: GetMetricsStatQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useQuery<GetMetricsStatQuery, GetMetricsStatQueryVariables>(GetMetricsStatDocument, options);
}
export function useGetMetricsStatLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useLazyQuery<GetMetricsStatQuery, GetMetricsStatQueryVariables>(GetMetricsStatDocument, options);
}
export function useGetMetricsStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>) {
  const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions};
  return Apollo.useSuspenseQuery<GetMetricsStatQuery, GetMetricsStatQueryVariables>(GetMetricsStatDocument, options);
}
export type GetMetricsStatQueryHookResult = ReturnType<typeof useGetMetricsStatQuery>;
export type GetMetricsStatLazyQueryHookResult = ReturnType<typeof useGetMetricsStatLazyQuery>;
export type GetMetricsStatSuspenseQueryHookResult = ReturnType<typeof useGetMetricsStatSuspenseQuery>;
export type GetMetricsStatQueryResult = Apollo.QueryResult<GetMetricsStatQuery, GetMetricsStatQueryVariables>;
export const GetSiteStatDocument = gql`
    query GetSiteStat($data: GetMetricsSiteStatInput!) {
  getSiteStat(data: $data) {
    metrics {
      success
      msg
      siteId
      nodeId
      type
      value
    }
  }
}
    `;

/**
 * __useGetSiteStatQuery__
 *
 * To run a query within a React component, call `useGetSiteStatQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteStatQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteStatQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSiteStatQuery(baseOptions: Apollo.QueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables> & ({ variables: GetSiteStatQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useQuery<GetSiteStatQuery, GetSiteStatQueryVariables>(GetSiteStatDocument, options);
}
export function useGetSiteStatLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useLazyQuery<GetSiteStatQuery, GetSiteStatQueryVariables>(GetSiteStatDocument, options);
}
export function useGetSiteStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>) {
  const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions};
  return Apollo.useSuspenseQuery<GetSiteStatQuery, GetSiteStatQueryVariables>(GetSiteStatDocument, options);
}
export type GetSiteStatQueryHookResult = ReturnType<typeof useGetSiteStatQuery>;
export type GetSiteStatLazyQueryHookResult = ReturnType<typeof useGetSiteStatLazyQuery>;
export type GetSiteStatSuspenseQueryHookResult = ReturnType<typeof useGetSiteStatSuspenseQuery>;
export type GetSiteStatQueryResult = Apollo.QueryResult<GetSiteStatQuery, GetSiteStatQueryVariables>;
export const GetMetricBySiteDocument = gql`
    query GetMetricBySite($data: GetMetricBySiteInput!) {
  getMetricBySite(data: $data) {
    metrics {
      msg
      success
      type
      values
    }
  }
}
    `;

/**
 * __useGetMetricBySiteQuery__
 *
 * To run a query within a React component, call `useGetMetricBySiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricBySiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricBySiteQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricBySiteQuery(baseOptions: Apollo.QueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables> & ({ variables: GetMetricBySiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useQuery<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>(GetMetricBySiteDocument, options);
}
export function useGetMetricBySiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useLazyQuery<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>(GetMetricBySiteDocument, options);
}
export function useGetMetricBySiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>) {
  const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions};
  return Apollo.useSuspenseQuery<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>(GetMetricBySiteDocument, options);
}
export type GetMetricBySiteQueryHookResult = ReturnType<typeof useGetMetricBySiteQuery>;
export type GetMetricBySiteLazyQueryHookResult = ReturnType<typeof useGetMetricBySiteLazyQuery>;
export type GetMetricBySiteSuspenseQueryHookResult = ReturnType<typeof useGetMetricBySiteSuspenseQuery>;
export type GetMetricBySiteQueryResult = Apollo.QueryResult<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>;
export const GetMetricsStatSubDocument = gql`
    subscription GetMetricsStatSub($data: SubMetricsStatInput!) {
  getMetricStatSub(data: $data) {
    msg
    nodeId
    success
    type
    value
    networkId
    siteId
    packageId
    dataPlanId
  }
}
    `;

/**
 * __useGetMetricsStatSubSubscription__
 *
 * To run a query within a React component, call `useGetMetricsStatSubSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsStatSubSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsStatSubSubscription({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsStatSubSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetMetricsStatSubSubscription, GetMetricsStatSubSubscriptionVariables> & ({ variables: GetMetricsStatSubSubscriptionVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useSubscription<GetMetricsStatSubSubscription, GetMetricsStatSubSubscriptionVariables>(GetMetricsStatSubDocument, options);
}
export type GetMetricsStatSubSubscriptionHookResult = ReturnType<typeof useGetMetricsStatSubSubscription>;
export type GetMetricsStatSubSubscriptionResult = Apollo.SubscriptionResult<GetMetricsStatSubSubscription>;
export const GetSiteMetricStatSubDocument = gql`
    subscription GetSiteMetricStatSub($data: SubSiteMetricsStatInput!) {
  getSiteMetricStatSub(data: $data) {
    msg
    siteId
    nodeId
    success
    type
    value
  }
}
    `;

/**
 * __useGetSiteMetricStatSubSubscription__
 *
 * To run a query within a React component, call `useGetSiteMetricStatSubSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteMetricStatSubSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteMetricStatSubSubscription({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSiteMetricStatSubSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetSiteMetricStatSubSubscription, GetSiteMetricStatSubSubscriptionVariables> & ({ variables: GetSiteMetricStatSubSubscriptionVariables; skip?: boolean; } | { skip: boolean; }) ) {
  const options = {...defaultOptions, ...baseOptions};
  return Apollo.useSubscription<GetSiteMetricStatSubSubscription, GetSiteMetricStatSubSubscriptionVariables>(GetSiteMetricStatSubDocument, options);
}
export type GetSiteMetricStatSubSubscriptionHookResult = ReturnType<typeof useGetSiteMetricStatSubSubscription>;
export type GetSiteMetricStatSubSubscriptionResult = Apollo.SubscriptionResult<GetSiteMetricStatSubSubscription>;