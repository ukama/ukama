import * as Types from './types';

import { gql } from '@apollo/client';
import { SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type SimPoolOverviewQueryVariables = Types.Exact<{
  simType: Types.Scalars['String']['input'];
  limit?: Types.Scalars['Int']['input'];
}>;


export type SimPoolOverviewQuery = { __typename?: 'Query', simPoolView: { __typename?: 'SimPoolView', simType: string, stats: { __typename?: 'SimPoolStatsSection', total?: number | null, available?: number | null, consumed?: number | null, failed?: number | null, esim?: number | null, physical?: number | null, pctAssigned?: number | null, lowStock?: boolean | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, sims: { __typename?: 'PoolSimsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, sims?: Array<{ __typename?: 'SimPoolResDto', id: string, iccid: string, msisdn: string, isAllocated: boolean, isFailed: boolean, simType: string, isPhysical: boolean, createdAt: string }> | null } } };


export const SimPoolOverviewDocument = gql`
    query SimPoolOverview($simType: String!, $limit: Int! = 20) {
  simPoolView(simType: $simType) {
    simType
    stats {
      error {
        ...SectionErrorFields
      }
      total
      available
      consumed
      failed
      esim
      physical
      pctAssigned
      lowStock
    }
    sims(limit: $limit) {
      error {
        ...SectionErrorFields
      }
      sims {
        id
        iccid
        msisdn
        isAllocated
        isFailed
        simType
        isPhysical
        createdAt
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useSimPoolOverviewQuery__
 *
 * To run a query within a React component, call `useSimPoolOverviewQuery` and pass it any options that fit your needs.
 * When your component renders, `useSimPoolOverviewQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSimPoolOverviewQuery({
 *   variables: {
 *      simType: // value for 'simType'
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useSimPoolOverviewQuery(baseOptions: Apollo.QueryHookOptions<SimPoolOverviewQuery, SimPoolOverviewQueryVariables> & ({ variables: SimPoolOverviewQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>(SimPoolOverviewDocument, options);
      }
export function useSimPoolOverviewLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>(SimPoolOverviewDocument, options);
        }
// @ts-ignore
export function useSimPoolOverviewSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>;
export function useSimPoolOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<SimPoolOverviewQuery | undefined, SimPoolOverviewQueryVariables>;
export function useSimPoolOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>(SimPoolOverviewDocument, options);
        }
export type SimPoolOverviewQueryHookResult = ReturnType<typeof useSimPoolOverviewQuery>;
export type SimPoolOverviewLazyQueryHookResult = ReturnType<typeof useSimPoolOverviewLazyQuery>;
export type SimPoolOverviewSuspenseQueryHookResult = ReturnType<typeof useSimPoolOverviewSuspenseQuery>;
export type SimPoolOverviewQueryResult = Apollo.QueryResult<SimPoolOverviewQuery, SimPoolOverviewQueryVariables>;