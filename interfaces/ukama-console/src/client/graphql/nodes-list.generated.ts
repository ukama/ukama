import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewNodeFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NodesListQueryVariables = Types.Exact<{
  networkId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type NodesListQuery = { __typename?: 'Query', nodesView: { __typename?: 'NodesView', networkId?: string | null, nodes: { __typename?: 'NodesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, nodes?: Array<{ __typename?: 'Node', id: string, name: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> | null }, health: { __typename?: 'GapSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };

export type NodePoolQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type NodePoolQuery = { __typename?: 'Query', nodesView: { __typename?: 'NodesView', networkId?: string | null, nodes: { __typename?: 'NodesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, nodes?: Array<{ __typename?: 'Node', id: string, name: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> | null } } };


export const NodesListDocument = gql`
    query NodesList($networkId: String) {
  nodesView(networkId: $networkId) {
    networkId
    nodes {
      error {
        ...SectionErrorFields
      }
      nodes {
        ...ViewNode
      }
    }
    health {
      error {
        ...SectionErrorFields
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewNodeFragmentDoc}`;

/**
 * __useNodesListQuery__
 *
 * To run a query within a React component, call `useNodesListQuery` and pass it any options that fit your needs.
 * When your component renders, `useNodesListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNodesListQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useNodesListQuery(baseOptions?: Apollo.QueryHookOptions<NodesListQuery, NodesListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NodesListQuery, NodesListQueryVariables>(NodesListDocument, options);
      }
export function useNodesListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NodesListQuery, NodesListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NodesListQuery, NodesListQueryVariables>(NodesListDocument, options);
        }
// @ts-ignore
export function useNodesListSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NodesListQuery, NodesListQueryVariables>): Apollo.UseSuspenseQueryResult<NodesListQuery, NodesListQueryVariables>;
export function useNodesListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodesListQuery, NodesListQueryVariables>): Apollo.UseSuspenseQueryResult<NodesListQuery | undefined, NodesListQueryVariables>;
export function useNodesListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodesListQuery, NodesListQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NodesListQuery, NodesListQueryVariables>(NodesListDocument, options);
        }
export type NodesListQueryHookResult = ReturnType<typeof useNodesListQuery>;
export type NodesListLazyQueryHookResult = ReturnType<typeof useNodesListLazyQuery>;
export type NodesListSuspenseQueryHookResult = ReturnType<typeof useNodesListSuspenseQuery>;
export type NodesListQueryResult = Apollo.QueryResult<NodesListQuery, NodesListQueryVariables>;
export const NodePoolDocument = gql`
    query NodePool {
  nodesView {
    networkId
    nodes {
      error {
        ...SectionErrorFields
      }
      nodes {
        ...ViewNode
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewNodeFragmentDoc}`;

/**
 * __useNodePoolQuery__
 *
 * To run a query within a React component, call `useNodePoolQuery` and pass it any options that fit your needs.
 * When your component renders, `useNodePoolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNodePoolQuery({
 *   variables: {
 *   },
 * });
 */
export function useNodePoolQuery(baseOptions?: Apollo.QueryHookOptions<NodePoolQuery, NodePoolQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NodePoolQuery, NodePoolQueryVariables>(NodePoolDocument, options);
      }
export function useNodePoolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NodePoolQuery, NodePoolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NodePoolQuery, NodePoolQueryVariables>(NodePoolDocument, options);
        }
// @ts-ignore
export function useNodePoolSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NodePoolQuery, NodePoolQueryVariables>): Apollo.UseSuspenseQueryResult<NodePoolQuery, NodePoolQueryVariables>;
export function useNodePoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodePoolQuery, NodePoolQueryVariables>): Apollo.UseSuspenseQueryResult<NodePoolQuery | undefined, NodePoolQueryVariables>;
export function useNodePoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodePoolQuery, NodePoolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NodePoolQuery, NodePoolQueryVariables>(NodePoolDocument, options);
        }
export type NodePoolQueryHookResult = ReturnType<typeof useNodePoolQuery>;
export type NodePoolLazyQueryHookResult = ReturnType<typeof useNodePoolLazyQuery>;
export type NodePoolSuspenseQueryHookResult = ReturnType<typeof useNodePoolSuspenseQuery>;
export type NodePoolQueryResult = Apollo.QueryResult<NodePoolQuery, NodePoolQueryVariables>;