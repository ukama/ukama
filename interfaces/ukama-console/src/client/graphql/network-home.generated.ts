import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewSiteFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type TopBarAlertsQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
}>;


export type TopBarAlertsQuery = { __typename?: 'Query', networkOverview: { __typename?: 'NetworkOverview', networkId: string, latestAlerts: { __typename?: 'AlertsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, notifications?: Array<{ __typename?: 'NotificationsDto', id: string, title: string, description: string, type: Types.Notification_Type, isRead: boolean, createdAt: string }> | null } } };

export type NetworkHomeQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
  alertLimit?: Types.Scalars['Int']['input'];
}>;


export type NetworkHomeQuery = { __typename?: 'Query', networkOverview: { __typename?: 'NetworkOverview', networkId: string, network: { __typename?: 'NetworkSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, network?: { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean } | null }, nodeStats: { __typename?: 'NodeStatsSection', total?: number | null, online?: number | null, offline?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, siteStats: { __typename?: 'SitesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, sites?: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, latitude: string, longitude: string, location: string, isDeactivated: boolean, installDate: string, createdAt: string }> | null }, subscriberStats: { __typename?: 'SubscriberStatsSection', total?: number | null, active?: number | null, inactive?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, latestAlerts: { __typename?: 'AlertsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, notifications?: Array<{ __typename?: 'NotificationsDto', id: string, title: string, description: string, type: Types.Notification_Type, isRead: boolean, createdAt: string }> | null }, kpis: { __typename?: 'KpisSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, metrics?: Array<{ __typename?: 'KpiEntryDto', key: string, value: number, timestamp: number, success: boolean }> | null } } };


export const TopBarAlertsDocument = gql`
    query TopBarAlerts($networkId: String!) {
  networkOverview(networkId: $networkId) {
    networkId
    latestAlerts(limit: 10) {
      error {
        ...SectionErrorFields
      }
      notifications {
        id
        title
        description
        type
        isRead
        createdAt
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useTopBarAlertsQuery__
 *
 * To run a query within a React component, call `useTopBarAlertsQuery` and pass it any options that fit your needs.
 * When your component renders, `useTopBarAlertsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useTopBarAlertsQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useTopBarAlertsQuery(baseOptions: Apollo.QueryHookOptions<TopBarAlertsQuery, TopBarAlertsQueryVariables> & ({ variables: TopBarAlertsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<TopBarAlertsQuery, TopBarAlertsQueryVariables>(TopBarAlertsDocument, options);
      }
export function useTopBarAlertsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<TopBarAlertsQuery, TopBarAlertsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<TopBarAlertsQuery, TopBarAlertsQueryVariables>(TopBarAlertsDocument, options);
        }
// @ts-ignore
export function useTopBarAlertsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<TopBarAlertsQuery, TopBarAlertsQueryVariables>): Apollo.UseSuspenseQueryResult<TopBarAlertsQuery, TopBarAlertsQueryVariables>;
export function useTopBarAlertsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<TopBarAlertsQuery, TopBarAlertsQueryVariables>): Apollo.UseSuspenseQueryResult<TopBarAlertsQuery | undefined, TopBarAlertsQueryVariables>;
export function useTopBarAlertsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<TopBarAlertsQuery, TopBarAlertsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<TopBarAlertsQuery, TopBarAlertsQueryVariables>(TopBarAlertsDocument, options);
        }
export type TopBarAlertsQueryHookResult = ReturnType<typeof useTopBarAlertsQuery>;
export type TopBarAlertsLazyQueryHookResult = ReturnType<typeof useTopBarAlertsLazyQuery>;
export type TopBarAlertsSuspenseQueryHookResult = ReturnType<typeof useTopBarAlertsSuspenseQuery>;
export type TopBarAlertsQueryResult = Apollo.QueryResult<TopBarAlertsQuery, TopBarAlertsQueryVariables>;
export const NetworkHomeDocument = gql`
    query NetworkHome($networkId: String!, $alertLimit: Int! = 5) {
  networkOverview(networkId: $networkId) {
    networkId
    network {
      error {
        ...SectionErrorFields
      }
      network {
        id
        name
        isDefault
      }
    }
    nodeStats {
      error {
        ...SectionErrorFields
      }
      total
      online
      offline
    }
    siteStats {
      error {
        ...SectionErrorFields
      }
      sites {
        ...ViewSite
      }
    }
    subscriberStats {
      error {
        ...SectionErrorFields
      }
      total
      active
      inactive
    }
    latestAlerts(limit: $alertLimit) {
      error {
        ...SectionErrorFields
      }
      notifications {
        id
        title
        description
        type
        isRead
        createdAt
      }
    }
    kpis {
      error {
        ...SectionErrorFields
      }
      metrics {
        key
        value
        timestamp
        success
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewSiteFragmentDoc}`;

/**
 * __useNetworkHomeQuery__
 *
 * To run a query within a React component, call `useNetworkHomeQuery` and pass it any options that fit your needs.
 * When your component renders, `useNetworkHomeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNetworkHomeQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *      alertLimit: // value for 'alertLimit'
 *   },
 * });
 */
export function useNetworkHomeQuery(baseOptions: Apollo.QueryHookOptions<NetworkHomeQuery, NetworkHomeQueryVariables> & ({ variables: NetworkHomeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NetworkHomeQuery, NetworkHomeQueryVariables>(NetworkHomeDocument, options);
      }
export function useNetworkHomeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NetworkHomeQuery, NetworkHomeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NetworkHomeQuery, NetworkHomeQueryVariables>(NetworkHomeDocument, options);
        }
// @ts-ignore
export function useNetworkHomeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NetworkHomeQuery, NetworkHomeQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkHomeQuery, NetworkHomeQueryVariables>;
export function useNetworkHomeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkHomeQuery, NetworkHomeQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkHomeQuery | undefined, NetworkHomeQueryVariables>;
export function useNetworkHomeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkHomeQuery, NetworkHomeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NetworkHomeQuery, NetworkHomeQueryVariables>(NetworkHomeDocument, options);
        }
export type NetworkHomeQueryHookResult = ReturnType<typeof useNetworkHomeQuery>;
export type NetworkHomeLazyQueryHookResult = ReturnType<typeof useNetworkHomeLazyQuery>;
export type NetworkHomeSuspenseQueryHookResult = ReturnType<typeof useNetworkHomeSuspenseQuery>;
export type NetworkHomeQueryResult = Apollo.QueryResult<NetworkHomeQuery, NetworkHomeQueryVariables>;