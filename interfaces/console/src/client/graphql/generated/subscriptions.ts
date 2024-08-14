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
  Network = 'NETWORK',
  NodeHealth = 'NODE_HEALTH',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Subscribers = 'SUBSCRIBERS'
}

export type GetMetricByTabInput = {
  from?: InputMaybe<Scalars['Float']['input']>;
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  step?: InputMaybe<Scalars['Float']['input']>;
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription: Scalars['Boolean']['input'];
};

export type GetNotificationsInput = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Role_Type;
  startTimestamp: Scalars['String']['input'];
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
  getMetricByTab: MetricsRes;
  getNotifications: NotificationsRes;
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
  data: GetNotificationsInput;
};

export type GetNotificationsQueryVariables = Exact<{
  data: GetNotificationsInput;
}>;


export type GetNotificationsQuery = { __typename?: 'Query', getNotifications: { __typename?: 'NotificationsRes', notifications: Array<{ __typename?: 'NotificationsResDto', id: string, type: Notification_Type, scope: Notification_Scope, title: string, isRead: boolean, createdAt: string, description: string }> } };

export type NotificationSubscriptionSubscriptionVariables = Exact<{
  data: GetNotificationsInput;
}>;


export type NotificationSubscriptionSubscription = { __typename?: 'Subscription', notificationSubscription: { __typename?: 'NotificationsResDto', id: string, type: Notification_Type, scope: Notification_Scope, title: string, isRead: boolean, createdAt: string, description: string } };


export const GetNotificationsDocument = gql`
    query GetNotifications($data: GetNotificationsInput!) {
  getNotifications(data: $data) {
    notifications {
      id
      type
      scope
      title
      isRead
      createdAt
      description
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
export function useGetNotificationsQuery(baseOptions: Apollo.QueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables> & ({ variables: GetNotificationsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
      }
export function useGetNotificationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
        }
export function useGetNotificationsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNotificationsQuery, GetNotificationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNotificationsQuery, GetNotificationsQueryVariables>(GetNotificationsDocument, options);
        }
export type GetNotificationsQueryHookResult = ReturnType<typeof useGetNotificationsQuery>;
export type GetNotificationsLazyQueryHookResult = ReturnType<typeof useGetNotificationsLazyQuery>;
export type GetNotificationsSuspenseQueryHookResult = ReturnType<typeof useGetNotificationsSuspenseQuery>;
export type GetNotificationsQueryResult = Apollo.QueryResult<GetNotificationsQuery, GetNotificationsQueryVariables>;
export const NotificationSubscriptionDocument = gql`
    subscription NotificationSubscription($data: GetNotificationsInput!) {
  notificationSubscription(data: $data) {
    id
    type
    scope
    title
    isRead
    createdAt
    description
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
 *      data: // value for 'data'
 *   },
 * });
 */
export function useNotificationSubscriptionSubscription(baseOptions: Apollo.SubscriptionHookOptions<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables> & ({ variables: NotificationSubscriptionSubscriptionVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<NotificationSubscriptionSubscription, NotificationSubscriptionSubscriptionVariables>(NotificationSubscriptionDocument, options);
      }
export type NotificationSubscriptionSubscriptionHookResult = ReturnType<typeof useNotificationSubscriptionSubscription>;
export type NotificationSubscriptionSubscriptionResult = Apollo.SubscriptionResult<NotificationSubscriptionSubscription>;