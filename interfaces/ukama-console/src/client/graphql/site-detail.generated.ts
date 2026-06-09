import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewNodeFragmentDoc, ViewSiteFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NetworkSiteDetailQueryVariables = Types.Exact<{
  siteId: Types.Scalars['String']['input'];
}>;


export type NetworkSiteDetailQuery = { __typename?: 'Query', siteView: { __typename?: 'SiteView', siteId: string, site: { __typename?: 'SiteSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, site?: { __typename?: 'SiteDto', id: string, name: string, networkId: string, latitude: string, longitude: string, location: string, isDeactivated: boolean, installDate: string, createdAt: string } | null }, nodes: { __typename?: 'NodesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, nodes?: Array<{ __typename?: 'Node', id: string, name: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> | null }, components: { __typename?: 'SiteComponentsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, components?: Array<{ __typename?: 'SiteComponentDto', elementType: string, componentId?: string | null, componentName?: string | null }> | null }, power: { __typename?: 'KpisSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, metrics?: Array<{ __typename?: 'KpiEntryDto', key: string, value: number, timestamp: number, success: boolean }> | null }, kpis: { __typename?: 'KpisSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, metrics?: Array<{ __typename?: 'KpiEntryDto', key: string, value: number, timestamp: number, success: boolean }> | null } } };


export const NetworkSiteDetailDocument = gql`
    query NetworkSiteDetail($siteId: String!) {
  siteView(siteId: $siteId) {
    siteId
    site {
      error {
        ...SectionErrorFields
      }
      site {
        ...ViewSite
      }
    }
    nodes {
      error {
        ...SectionErrorFields
      }
      nodes {
        ...ViewNode
      }
    }
    components {
      error {
        ...SectionErrorFields
      }
      components {
        elementType
        componentId
        componentName
      }
    }
    power {
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
${ViewSiteFragmentDoc}
${ViewNodeFragmentDoc}`;

/**
 * __useNetworkSiteDetailQuery__
 *
 * To run a query within a React component, call `useNetworkSiteDetailQuery` and pass it any options that fit your needs.
 * When your component renders, `useNetworkSiteDetailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNetworkSiteDetailQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useNetworkSiteDetailQuery(baseOptions: Apollo.QueryHookOptions<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables> & ({ variables: NetworkSiteDetailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>(NetworkSiteDetailDocument, options);
      }
export function useNetworkSiteDetailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>(NetworkSiteDetailDocument, options);
        }
// @ts-ignore
export function useNetworkSiteDetailSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>;
export function useNetworkSiteDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkSiteDetailQuery | undefined, NetworkSiteDetailQueryVariables>;
export function useNetworkSiteDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>(NetworkSiteDetailDocument, options);
        }
export type NetworkSiteDetailQueryHookResult = ReturnType<typeof useNetworkSiteDetailQuery>;
export type NetworkSiteDetailLazyQueryHookResult = ReturnType<typeof useNetworkSiteDetailLazyQuery>;
export type NetworkSiteDetailSuspenseQueryHookResult = ReturnType<typeof useNetworkSiteDetailSuspenseQuery>;
export type NetworkSiteDetailQueryResult = Apollo.QueryResult<NetworkSiteDetailQuery, NetworkSiteDetailQueryVariables>;