import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type MetricsRangeQueryVariables = Types.Exact<{
  data: Types.MetricsRangeInput;
}>;


export type MetricsRangeQuery = { __typename?: 'Query', metricsRange: { __typename?: 'MetricsRes', metrics: Array<{ __typename?: 'MetricRes', type: string, success: boolean, nodeId?: string | null, siteId?: string | null, values: Array<Array<number>>, unit?: string | null, format?: string | null }> } };


export const MetricsRangeDocument = gql`
    query MetricsRange($data: MetricsRangeInput!) {
  metricsRange(data: $data) {
    metrics {
      type
      success
      nodeId
      siteId
      values
      unit
      format
    }
  }
}
    `;

/**
 * __useMetricsRangeQuery__
 *
 * To run a query within a React component, call `useMetricsRangeQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetricsRangeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetricsRangeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useMetricsRangeQuery(baseOptions: Apollo.QueryHookOptions<MetricsRangeQuery, MetricsRangeQueryVariables> & ({ variables: MetricsRangeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MetricsRangeQuery, MetricsRangeQueryVariables>(MetricsRangeDocument, options);
      }
export function useMetricsRangeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MetricsRangeQuery, MetricsRangeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MetricsRangeQuery, MetricsRangeQueryVariables>(MetricsRangeDocument, options);
        }
// @ts-ignore
export function useMetricsRangeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<MetricsRangeQuery, MetricsRangeQueryVariables>): Apollo.UseSuspenseQueryResult<MetricsRangeQuery, MetricsRangeQueryVariables>;
export function useMetricsRangeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MetricsRangeQuery, MetricsRangeQueryVariables>): Apollo.UseSuspenseQueryResult<MetricsRangeQuery | undefined, MetricsRangeQueryVariables>;
export function useMetricsRangeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MetricsRangeQuery, MetricsRangeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MetricsRangeQuery, MetricsRangeQueryVariables>(MetricsRangeDocument, options);
        }
export type MetricsRangeQueryHookResult = ReturnType<typeof useMetricsRangeQuery>;
export type MetricsRangeLazyQueryHookResult = ReturnType<typeof useMetricsRangeLazyQuery>;
export type MetricsRangeSuspenseQueryHookResult = ReturnType<typeof useMetricsRangeSuspenseQuery>;
export type MetricsRangeQueryResult = Apollo.QueryResult<MetricsRangeQuery, MetricsRangeQueryVariables>;