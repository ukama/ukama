import * as Types from './types';

import { gql } from '@apollo/client';
import { SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NodeKpisQueryVariables = Types.Exact<{
  nodeId: Types.Scalars['String']['input'];
}>;


export type NodeKpisQuery = { __typename?: 'Query', nodeView: { __typename?: 'NodeView', nodeId: string, kpis: { __typename?: 'KpisSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, metrics?: Array<{ __typename?: 'KpiEntryDto', key: string, value: number, timestamp: number, success: boolean, label?: string | null, unit?: string | null, format?: string | null, threshold?: { __typename?: 'MetricThreshold', min: number, normal: number, max: number } | null }> | null } } };


export const NodeKpisDocument = gql`
    query NodeKpis($nodeId: String!) {
  nodeView(nodeId: $nodeId) {
    nodeId
    kpis {
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
        threshold {
          min
          normal
          max
        }
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useNodeKpisQuery__
 *
 * To run a query within a React component, call `useNodeKpisQuery` and pass it any options that fit your needs.
 * When your component renders, `useNodeKpisQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNodeKpisQuery({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *   },
 * });
 */
export function useNodeKpisQuery(baseOptions: Apollo.QueryHookOptions<NodeKpisQuery, NodeKpisQueryVariables> & ({ variables: NodeKpisQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NodeKpisQuery, NodeKpisQueryVariables>(NodeKpisDocument, options);
      }
export function useNodeKpisLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NodeKpisQuery, NodeKpisQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NodeKpisQuery, NodeKpisQueryVariables>(NodeKpisDocument, options);
        }
// @ts-ignore
export function useNodeKpisSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NodeKpisQuery, NodeKpisQueryVariables>): Apollo.UseSuspenseQueryResult<NodeKpisQuery, NodeKpisQueryVariables>;
export function useNodeKpisSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodeKpisQuery, NodeKpisQueryVariables>): Apollo.UseSuspenseQueryResult<NodeKpisQuery | undefined, NodeKpisQueryVariables>;
export function useNodeKpisSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodeKpisQuery, NodeKpisQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NodeKpisQuery, NodeKpisQueryVariables>(NodeKpisDocument, options);
        }
export type NodeKpisQueryHookResult = ReturnType<typeof useNodeKpisQuery>;
export type NodeKpisLazyQueryHookResult = ReturnType<typeof useNodeKpisLazyQuery>;
export type NodeKpisSuspenseQueryHookResult = ReturnType<typeof useNodeKpisSuspenseQuery>;
export type NodeKpisQueryResult = Apollo.QueryResult<NodeKpisQuery, NodeKpisQueryVariables>;