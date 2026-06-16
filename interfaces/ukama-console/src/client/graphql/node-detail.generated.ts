import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewNodeFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NodeDetailQueryVariables = Types.Exact<{
  nodeId: Types.Scalars['String']['input'];
}>;


export type NodeDetailQuery = { __typename?: 'Query', nodeView: { __typename?: 'NodeView', nodeId: string, node: { __typename?: 'NodeSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, node?: { __typename?: 'Node', latitude: string, longitude: string, id: string, name: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string }>, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } | null }, site: { __typename?: 'SiteSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, site?: { __typename?: 'SiteDto', id: string, name: string } | null }, siblings: { __typename?: 'NodesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, nodes?: Array<{ __typename?: 'Node', id: string, name: string }> | null }, software: { __typename?: 'SoftwareSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, softwares?: { __typename?: 'Softwares', software: Array<{ __typename?: 'Software', name: string, status: Types.SoftwareStatusEnum, currentVersion: string, desiredVersion: string, releaseDate: string }> } | null } } };


export const NodeDetailDocument = gql`
    query NodeDetail($nodeId: String!) {
  nodeView(nodeId: $nodeId) {
    nodeId
    node {
      error {
        ...SectionErrorFields
      }
      node {
        ...ViewNode
        latitude
        longitude
        attached {
          id
          name
        }
      }
    }
    site {
      error {
        ...SectionErrorFields
      }
      site {
        id
        name
      }
    }
    siblings {
      error {
        ...SectionErrorFields
      }
      nodes {
        id
        name
      }
    }
    software {
      error {
        ...SectionErrorFields
      }
      softwares {
        software {
          name
          status
          currentVersion
          desiredVersion
          releaseDate
        }
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewNodeFragmentDoc}`;

/**
 * __useNodeDetailQuery__
 *
 * To run a query within a React component, call `useNodeDetailQuery` and pass it any options that fit your needs.
 * When your component renders, `useNodeDetailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNodeDetailQuery({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *   },
 * });
 */
export function useNodeDetailQuery(baseOptions: Apollo.QueryHookOptions<NodeDetailQuery, NodeDetailQueryVariables> & ({ variables: NodeDetailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NodeDetailQuery, NodeDetailQueryVariables>(NodeDetailDocument, options);
      }
export function useNodeDetailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NodeDetailQuery, NodeDetailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NodeDetailQuery, NodeDetailQueryVariables>(NodeDetailDocument, options);
        }
// @ts-ignore
export function useNodeDetailSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NodeDetailQuery, NodeDetailQueryVariables>): Apollo.UseSuspenseQueryResult<NodeDetailQuery, NodeDetailQueryVariables>;
export function useNodeDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodeDetailQuery, NodeDetailQueryVariables>): Apollo.UseSuspenseQueryResult<NodeDetailQuery | undefined, NodeDetailQueryVariables>;
export function useNodeDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NodeDetailQuery, NodeDetailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NodeDetailQuery, NodeDetailQueryVariables>(NodeDetailDocument, options);
        }
export type NodeDetailQueryHookResult = ReturnType<typeof useNodeDetailQuery>;
export type NodeDetailLazyQueryHookResult = ReturnType<typeof useNodeDetailLazyQuery>;
export type NodeDetailSuspenseQueryHookResult = ReturnType<typeof useNodeDetailSuspenseQuery>;
export type NodeDetailQueryResult = Apollo.QueryResult<NodeDetailQuery, NodeDetailQueryVariables>;