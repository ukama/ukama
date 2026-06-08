import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewSiteFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type SitesListQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
}>;


export type SitesListQuery = { __typename?: 'Query', sitesView: { __typename?: 'SitesView', networkId: string, sites: { __typename?: 'SitesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, sites?: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, latitude: string, longitude: string, location: string, isDeactivated: boolean, installDate: string, createdAt: string }> | null }, nodeCounts: { __typename?: 'SiteNodeCountsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, counts?: Array<{ __typename?: 'SiteNodeCountDto', siteId: string, total: number, online: number, offline: number }> | null }, customers: { __typename?: 'SiteCustomersSection', count?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, kpis: { __typename?: 'GapSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };


export const SitesListDocument = gql`
    query SitesList($networkId: String!) {
  sitesView(networkId: $networkId) {
    networkId
    sites {
      error {
        ...SectionErrorFields
      }
      sites {
        ...ViewSite
      }
    }
    nodeCounts {
      error {
        ...SectionErrorFields
      }
      counts {
        siteId
        total
        online
        offline
      }
    }
    customers {
      error {
        ...SectionErrorFields
      }
      count
    }
    kpis {
      error {
        ...SectionErrorFields
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewSiteFragmentDoc}`;

/**
 * __useSitesListQuery__
 *
 * To run a query within a React component, call `useSitesListQuery` and pass it any options that fit your needs.
 * When your component renders, `useSitesListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSitesListQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useSitesListQuery(baseOptions: Apollo.QueryHookOptions<SitesListQuery, SitesListQueryVariables> & ({ variables: SitesListQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SitesListQuery, SitesListQueryVariables>(SitesListDocument, options);
      }
export function useSitesListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SitesListQuery, SitesListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SitesListQuery, SitesListQueryVariables>(SitesListDocument, options);
        }
// @ts-ignore
export function useSitesListSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SitesListQuery, SitesListQueryVariables>): Apollo.UseSuspenseQueryResult<SitesListQuery, SitesListQueryVariables>;
export function useSitesListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SitesListQuery, SitesListQueryVariables>): Apollo.UseSuspenseQueryResult<SitesListQuery | undefined, SitesListQueryVariables>;
export function useSitesListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SitesListQuery, SitesListQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SitesListQuery, SitesListQueryVariables>(SitesListDocument, options);
        }
export type SitesListQueryHookResult = ReturnType<typeof useSitesListQuery>;
export type SitesListLazyQueryHookResult = ReturnType<typeof useSitesListLazyQuery>;
export type SitesListSuspenseQueryHookResult = ReturnType<typeof useSitesListSuspenseQuery>;
export type SitesListQueryResult = Apollo.QueryResult<SitesListQuery, SitesListQueryVariables>;