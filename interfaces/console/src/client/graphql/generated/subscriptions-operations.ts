/** Internal type. DO NOT USE DIRECTLY. */
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
/** Internal type. DO NOT USE DIRECTLY. */
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
import type * as Types from './subscription-schema-types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type GetNotificationsQueryVariables = Exact<{
  networkId: string;
  orgId: string;
  orgName: string;
  role: string;
  startTimestamp: string;
  subscriberId: string;
  userId: string;
}>;


export type GetNotificationsQuery = { getNotifications: { notifications: Array<{ id: string, type: Types.Notification_Type, scope: Types.Notification_Scope, title: string, isRead: boolean, eventKey: string, createdAt: string, resourceId: string, description: string, redirect: { action: string, title: string } }> } };

export type NotificationSubscriptionSubscriptionVariables = Exact<{
  networkId: string;
  orgId: string;
  orgName: string;
  role: string;
  startTimestamp: string;
  subscriberId: string;
  userId: string;
}>;


export type NotificationSubscriptionSubscription = { notificationSubscription: { id: string, type: Types.Notification_Type, scope: Types.Notification_Scope, title: string, isRead: boolean, eventKey: string, createdAt: string, resourceId: string, description: string, redirect: { action: string, title: string } } };

export type GetMetricByTabQueryVariables = Exact<{
  data: Types.GetMetricByTabInput;
}>;


export type GetMetricByTabQuery = { getMetricByTab: { metrics: Array<{ dataPlanId: string | null, format: string | null, msg: string, networkId: string | null, nodeId: string | null, siteId: string | null, packageId: string | null, success: boolean, tickInterval: number | null, tickPositions: Array<number> | null, type: string, unit: string | null, values: Array<Array<number>>, threshold: { max: number, min: number, normal: number } | null }> } };

export type GetMetricsStatQueryVariables = Exact<{
  data: Types.GetMetricsStatInput;
}>;


export type GetMetricsStatQuery = { getMetricsStat: { metrics: Array<{ dataPlanId: string | null, format: string | null, msg: string, networkId: string | null, nodeId: string, packageId: string | null, siteId: string | null, success: boolean, tickInterval: number | null, tickPositions: Array<number> | null, type: string, unit: string | null, value: number, threshold: { max: number, min: number, normal: number } | null }> } };

export type GetSiteStatQueryVariables = Exact<{
  data: Types.GetMetricsSiteStatInput;
}>;


export type GetSiteStatQuery = { getSiteStat: { metrics: Array<{ dataPlanId: string | null, format: string | null, msg: string, networkId: string | null, nodeId: string, siteId: string | null, packageId: string | null, success: boolean, tickInterval: number | null, tickPositions: Array<number> | null, type: string, unit: string | null, value: number, threshold: { max: number, min: number, normal: number } | null }> } };

export type GetMetricBySiteQueryVariables = Exact<{
  data: Types.GetMetricBySiteInput;
}>;


export type GetMetricBySiteQuery = { getMetricBySite: { metrics: Array<{ dataPlanId: string | null, format: string | null, msg: string, networkId: string | null, nodeId: string | null, siteId: string | null, packageId: string | null, success: boolean, tickInterval: number | null, tickPositions: Array<number> | null, type: string, unit: string | null, values: Array<Array<number>>, threshold: { max: number, min: number, normal: number } | null }> } };

export type GetMetricsStatSubSubscriptionVariables = Exact<{
  data: Types.SubMetricsStatInput;
}>;


export type GetMetricsStatSubSubscription = { getMetricStatSub: { msg: string, nodeId: string, success: boolean, type: string, value: Array<number>, networkId: string | null, siteId: string, packageId: string | null, dataPlanId: string | null } };

export type GetSiteMetricStatSubSubscriptionVariables = Exact<{
  data: Types.SubSiteMetricsStatInput;
}>;


export type GetSiteMetricStatSubSubscription = { getSiteMetricStatSub: { msg: string, siteId: string, nodeId: string, success: boolean, type: string, value: Array<number> } };


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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
      }
export function useGetNotificationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
        }
// @ts-ignore
export function useGetNotificationsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNotificationsQuery, GetNotificationsQueryVariables>;
export function useGetNotificationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNotificationsQuery | undefined, GetNotificationsQueryVariables>;
export function useGetNotificationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables>(NotificationSubscriptionDocument, options);
      }
export type NotificationSubscriptionSubscriptionHookResult = ReturnType<typeof useNotificationSubscriptionSubscription>;
export type NotificationSubscriptionSubscriptionResult = Apollo.SubscriptionResult<NotificationSubscriptionSubscription>;
export const GetMetricByTabDocument = gql`
    query GetMetricByTab($data: GetMetricByTabInput!) {
  getMetricByTab(data: $data) {
    metrics {
      dataPlanId
      format
      msg
      networkId
      nodeId
      siteId
      packageId
      success
      tickInterval
      threshold {
        max
        min
        normal
      }
      tickPositions
      type
      unit
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
export function useGetMetricByTabQuery(baseOptions: Apollo.QueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables> & ({ variables: GetMetricByTabQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
      }
export function useGetMetricByTabLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricByTabQuery, GetMetricByTabQueryVariables>(GetMetricByTabDocument, options);
        }
// @ts-ignore
export function useGetMetricByTabSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricByTabQuery, GetMetricByTabQueryVariables>;
export function useGetMetricByTabSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricByTabQuery | undefined, GetMetricByTabQueryVariables>;
export function useGetMetricByTabSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricByTabQuery, GetMetricByTabQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
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
      dataPlanId
      format
      msg
      networkId
      nodeId
      packageId
      siteId
      success
      threshold {
        max
        min
        normal
      }
      tickInterval
      tickPositions
      type
      unit
      value
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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricsStatQuery, GetMetricsStatQueryVariables>(GetMetricsStatDocument, options);
      }
export function useGetMetricsStatLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricsStatQuery, GetMetricsStatQueryVariables>(GetMetricsStatDocument, options);
        }
// @ts-ignore
export function useGetMetricsStatSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricsStatQuery, GetMetricsStatQueryVariables>;
export function useGetMetricsStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricsStatQuery | undefined, GetMetricsStatQueryVariables>;
export function useGetMetricsStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricsStatQuery, GetMetricsStatQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
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
      dataPlanId
      format
      msg
      networkId
      nodeId
      siteId
      packageId
      success
      tickInterval
      threshold {
        max
        min
        normal
      }
      tickPositions
      type
      unit
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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSiteStatQuery, GetSiteStatQueryVariables>(GetSiteStatDocument, options);
      }
export function useGetSiteStatLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSiteStatQuery, GetSiteStatQueryVariables>(GetSiteStatDocument, options);
        }
// @ts-ignore
export function useGetSiteStatSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteStatQuery, GetSiteStatQueryVariables>;
export function useGetSiteStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteStatQuery | undefined, GetSiteStatQueryVariables>;
export function useGetSiteStatSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteStatQuery, GetSiteStatQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
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
      dataPlanId
      format
      msg
      networkId
      nodeId
      siteId
      packageId
      success
      tickInterval
      threshold {
        max
        min
        normal
      }
      tickPositions
      type
      unit
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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>(GetMetricBySiteDocument, options);
      }
export function useGetMetricBySiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>(GetMetricBySiteDocument, options);
        }
// @ts-ignore
export function useGetMetricBySiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>;
export function useGetMetricBySiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetMetricBySiteQuery | undefined, GetMetricBySiteQueryVariables>;
export function useGetMetricBySiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMetricBySiteQuery, GetMetricBySiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
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
        const options = {...defaultOptions, ...baseOptions}
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
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetSiteMetricStatSubSubscription, GetSiteMetricStatSubSubscriptionVariables>(GetSiteMetricStatSubDocument, options);
      }
export type GetSiteMetricStatSubSubscriptionHookResult = ReturnType<typeof useGetSiteMetricStatSubSubscription>;
export type GetSiteMetricStatSubSubscriptionResult = Apollo.SubscriptionResult<GetSiteMetricStatSubSubscription>;