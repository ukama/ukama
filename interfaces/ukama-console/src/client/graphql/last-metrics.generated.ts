import * as Types from './types';

import { gql } from '@apollo/client';
import { SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type MetricsLastQueryVariables = Types.Exact<{
  data: Types.MetricsLastInput;
}>;


export type MetricsLastQuery = { __typename?: 'Query', metricsLast: { __typename?: 'KpisSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, metrics?: Array<{ __typename?: 'KpiEntryDto', key: string, value: number, timestamp: number, success: boolean, label?: string | null, unit?: string | null, format?: string | null }> | null } };


export const MetricsLastDocument = gql`
    query MetricsLast($data: MetricsLastInput!) {
  metricsLast(data: $data) {
    error {
      ...SectionErrorFields
    }
    metrics {
      key
      value
      timestamp
      success
      label
      unit
      format
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useMetricsLastQuery__
 *
 * To run a query within a React component, call `useMetricsLastQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetricsLastQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetricsLastQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useMetricsLastQuery(baseOptions: Apollo.QueryHookOptions<MetricsLastQuery, MetricsLastQueryVariables> & ({ variables: MetricsLastQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MetricsLastQuery, MetricsLastQueryVariables>(MetricsLastDocument, options);
      }
export function useMetricsLastLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MetricsLastQuery, MetricsLastQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MetricsLastQuery, MetricsLastQueryVariables>(MetricsLastDocument, options);
        }
// @ts-ignore
export function useMetricsLastSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<MetricsLastQuery, MetricsLastQueryVariables>): Apollo.UseSuspenseQueryResult<MetricsLastQuery, MetricsLastQueryVariables>;
export function useMetricsLastSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MetricsLastQuery, MetricsLastQueryVariables>): Apollo.UseSuspenseQueryResult<MetricsLastQuery | undefined, MetricsLastQueryVariables>;
export function useMetricsLastSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MetricsLastQuery, MetricsLastQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MetricsLastQuery, MetricsLastQueryVariables>(MetricsLastDocument, options);
        }
export type MetricsLastQueryHookResult = ReturnType<typeof useMetricsLastQuery>;
export type MetricsLastLazyQueryHookResult = ReturnType<typeof useMetricsLastLazyQuery>;
export type MetricsLastSuspenseQueryHookResult = ReturnType<typeof useMetricsLastSuspenseQuery>;
export type MetricsLastQueryResult = Apollo.QueryResult<MetricsLastQuery, MetricsLastQueryVariables>;