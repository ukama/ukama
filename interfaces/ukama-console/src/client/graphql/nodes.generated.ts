import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NodeFragment = { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } };

export type GetNodeQueryVariables = Types.Exact<{
  data: Types.NodeInput;
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesQueryVariables = Types.Exact<{
  data: Types.NodesFilterInput;
}>;


export type GetNodesQuery = { __typename?: 'Query', getNodes: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type DeleteNodeMutationVariables = Types.Exact<{
  data: Types.NodeInput;
}>;


export type DeleteNodeMutation = { __typename?: 'Mutation', deleteNodeFromOrg: { __typename?: 'DeleteNode', id: string } };

export type AttachNodeMutationVariables = Types.Exact<{
  data: Types.AttachNodeInput;
}>;


export type AttachNodeMutation = { __typename?: 'Mutation', attachNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type DetachhNodeMutationVariables = Types.Exact<{
  data: Types.NodeInput;
}>;


export type DetachhNodeMutation = { __typename?: 'Mutation', detachhNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type AddNodeMutationVariables = Types.Exact<{
  data: Types.AddNodeInput;
}>;


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type ReleaseNodeFromSiteMutationVariables = Types.Exact<{
  data: Types.NodeInput;
}>;


export type ReleaseNodeFromSiteMutation = { __typename?: 'Mutation', releaseNodeFromSite: { __typename?: 'CBooleanResponse', success: boolean } };

export type AddNodeToSiteMutationVariables = Types.Exact<{
  data: Types.AddNodeToSiteInput;
}>;


export type AddNodeToSiteMutation = { __typename?: 'Mutation', addNodeToSite: { __typename?: 'CBooleanResponse', success: boolean } };

export type UpdateNodeStateMutationVariables = Types.Exact<{
  data: Types.UpdateNodeStateInput;
}>;


export type UpdateNodeStateMutation = { __typename?: 'Mutation', updateNodeState: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodesForSiteQueryVariables = Types.Exact<{
  siteId: Types.Scalars['String']['input'];
}>;


export type GetNodesForSiteQuery = { __typename?: 'Query', getNodesForSite: { __typename?: 'Nodes', nodes: Array<{ __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> } };

export type UpdateNodeMutationVariables = Types.Exact<{
  data: Types.UpdateNodeInput;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'Node', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ __typename?: 'AttachedNodes', id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }>, site: { __typename?: 'NodeSite', nodeId?: string | null, siteId?: string | null, networkId?: string | null, addedAt?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } } };

export type GetNodeAppsQueryVariables = Types.Exact<{
  data: Types.NodeAppsChangeLogInput;
}>;


export type GetNodeAppsQuery = { __typename?: 'Query', getNodeApps: { __typename?: 'NodeApps', type: Types.NodeTypeEnum, apps: Array<{ __typename?: 'NodeApp', name: string, date: number, version: string, cpu: string, memory: string, notes: string }> } };

export type GetNodeStateQueryVariables = Types.Exact<{
  getNodeStateId: Types.Scalars['String']['input'];
}>;


export type GetNodeStateQuery = { __typename?: 'Query', getNodeState: { __typename?: 'NodeStateRes', id: string, nodeId: string, previousStateId?: string | null, previousState?: Types.NodeStateEnum | null, currentState: Types.NodeStateEnum, createdAt: string } };

export type RestartNodeMutationVariables = Types.Exact<{
  data: Types.RestartNodeInputDto;
}>;


export type RestartNodeMutation = { __typename?: 'Mutation', restartNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleInternetSwitchMutationVariables = Types.Exact<{
  data: Types.ToggleInternetSwitchInputDto;
}>;


export type ToggleInternetSwitchMutation = { __typename?: 'Mutation', toggleInternetSwitch: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleRfStatusMutationVariables = Types.Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleRfStatusMutation = { __typename?: 'Mutation', toggleRFStatus: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleServiceMutationVariables = Types.Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleServiceMutation = { __typename?: 'Mutation', toggleService: { __typename?: 'CBooleanResponse', success: boolean } };

export type GetHealthReportQueryVariables = Types.Exact<{
  data: Types.GetHealthReportInputDto;
}>;


export type GetHealthReportQuery = { __typename?: 'Query', getHealthReport: { __typename?: 'HealthInfo', id: string, nodeId: string, timestamp: string, system: Array<{ __typename?: 'HealthSystemInfo', id: string, healthId: string, name: string, value: string }>, capps: Array<{ __typename?: 'HealthCappInfo', id: string, space: string, name: string, tag: string, status: string, resources: Array<{ __typename?: 'HealthResourceInfo', id: string, cappId: string, name: string, value: string }> }> } };

export const NodeFragmentDoc = gql`
    fragment node on Node {
  id
  name
  latitude
  longitude
  type
  attached {
    id
    name
    latitude
    longitude
    type
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
  site {
    nodeId
    siteId
    networkId
    addedAt
  }
  status {
    connectivity
    state
  }
}
    `;
export const GetNodeDocument = gql`
    query GetNode($data: NodeInput!) {
  getNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodeQuery__
 *
 * To run a query within a React component, call `useGetNodeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeQuery(baseOptions: Apollo.QueryHookOptions<GetNodeQuery, GetNodeQueryVariables> & ({ variables: GetNodeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
      }
export function useGetNodeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
// @ts-ignore
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeQuery, GetNodeQueryVariables>;
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeQuery | undefined, GetNodeQueryVariables>;
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export type GetNodeQueryHookResult = ReturnType<typeof useGetNodeQuery>;
export type GetNodeLazyQueryHookResult = ReturnType<typeof useGetNodeLazyQuery>;
export type GetNodeSuspenseQueryHookResult = ReturnType<typeof useGetNodeSuspenseQuery>;
export type GetNodeQueryResult = Apollo.QueryResult<GetNodeQuery, GetNodeQueryVariables>;
export const GetNodesDocument = gql`
    query GetNodes($data: NodesFilterInput!) {
  getNodes(data: $data) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesQuery__
 *
 * To run a query within a React component, call `useGetNodesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodesQuery(baseOptions: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables> & ({ variables: GetNodesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
      }
export function useGetNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
// @ts-ignore
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesQuery, GetNodesQueryVariables>;
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesQuery | undefined, GetNodesQueryVariables>;
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<typeof useGetNodesLazyQuery>;
export type GetNodesSuspenseQueryHookResult = ReturnType<typeof useGetNodesSuspenseQuery>;
export type GetNodesQueryResult = Apollo.QueryResult<GetNodesQuery, GetNodesQueryVariables>;
export const DeleteNodeDocument = gql`
    mutation deleteNode($data: NodeInput!) {
  deleteNodeFromOrg(data: $data) {
    id
  }
}
    `;
export type DeleteNodeMutationFn = Apollo.MutationFunction<DeleteNodeMutation, DeleteNodeMutationVariables>;

/**
 * __useDeleteNodeMutation__
 *
 * To run a mutation, you first call `useDeleteNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteNodeMutation, { data, loading, error }] = useDeleteNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteNodeMutation(baseOptions?: Apollo.MutationHookOptions<DeleteNodeMutation, DeleteNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteNodeMutation, DeleteNodeMutationVariables>(DeleteNodeDocument, options);
      }
export type DeleteNodeMutationHookResult = ReturnType<typeof useDeleteNodeMutation>;
export type DeleteNodeMutationResult = Apollo.MutationResult<DeleteNodeMutation>;
export type DeleteNodeMutationOptions = Apollo.BaseMutationOptions<DeleteNodeMutation, DeleteNodeMutationVariables>;
export const AttachNodeDocument = gql`
    mutation attachNode($data: AttachNodeInput!) {
  attachNode(data: $data) {
    success
  }
}
    `;
export type AttachNodeMutationFn = Apollo.MutationFunction<AttachNodeMutation, AttachNodeMutationVariables>;

/**
 * __useAttachNodeMutation__
 *
 * To run a mutation, you first call `useAttachNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAttachNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [attachNodeMutation, { data, loading, error }] = useAttachNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAttachNodeMutation(baseOptions?: Apollo.MutationHookOptions<AttachNodeMutation, AttachNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AttachNodeMutation, AttachNodeMutationVariables>(AttachNodeDocument, options);
      }
export type AttachNodeMutationHookResult = ReturnType<typeof useAttachNodeMutation>;
export type AttachNodeMutationResult = Apollo.MutationResult<AttachNodeMutation>;
export type AttachNodeMutationOptions = Apollo.BaseMutationOptions<AttachNodeMutation, AttachNodeMutationVariables>;
export const DetachhNodeDocument = gql`
    mutation detachhNode($data: NodeInput!) {
  detachhNode(data: $data) {
    success
  }
}
    `;
export type DetachhNodeMutationFn = Apollo.MutationFunction<DetachhNodeMutation, DetachhNodeMutationVariables>;

/**
 * __useDetachhNodeMutation__
 *
 * To run a mutation, you first call `useDetachhNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDetachhNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [detachhNodeMutation, { data, loading, error }] = useDetachhNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDetachhNodeMutation(baseOptions?: Apollo.MutationHookOptions<DetachhNodeMutation, DetachhNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DetachhNodeMutation, DetachhNodeMutationVariables>(DetachhNodeDocument, options);
      }
export type DetachhNodeMutationHookResult = ReturnType<typeof useDetachhNodeMutation>;
export type DetachhNodeMutationResult = Apollo.MutationResult<DetachhNodeMutation>;
export type DetachhNodeMutationOptions = Apollo.BaseMutationOptions<DetachhNodeMutation, DetachhNodeMutationVariables>;
export const AddNodeDocument = gql`
    mutation addNode($data: AddNodeInput!) {
  addNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type AddNodeMutationFn = Apollo.MutationFunction<AddNodeMutation, AddNodeMutationVariables>;

/**
 * __useAddNodeMutation__
 *
 * To run a mutation, you first call `useAddNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeMutation, { data, loading, error }] = useAddNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeMutation, AddNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeMutation, AddNodeMutationVariables>(AddNodeDocument, options);
      }
export type AddNodeMutationHookResult = ReturnType<typeof useAddNodeMutation>;
export type AddNodeMutationResult = Apollo.MutationResult<AddNodeMutation>;
export type AddNodeMutationOptions = Apollo.BaseMutationOptions<AddNodeMutation, AddNodeMutationVariables>;
export const ReleaseNodeFromSiteDocument = gql`
    mutation releaseNodeFromSite($data: NodeInput!) {
  releaseNodeFromSite(data: $data) {
    success
  }
}
    `;
export type ReleaseNodeFromSiteMutationFn = Apollo.MutationFunction<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;

/**
 * __useReleaseNodeFromSiteMutation__
 *
 * To run a mutation, you first call `useReleaseNodeFromSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useReleaseNodeFromSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [releaseNodeFromSiteMutation, { data, loading, error }] = useReleaseNodeFromSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useReleaseNodeFromSiteMutation(baseOptions?: Apollo.MutationHookOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>(ReleaseNodeFromSiteDocument, options);
      }
export type ReleaseNodeFromSiteMutationHookResult = ReturnType<typeof useReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationResult = Apollo.MutationResult<ReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationOptions = Apollo.BaseMutationOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;
export const AddNodeToSiteDocument = gql`
    mutation addNodeToSite($data: AddNodeToSiteInput!) {
  addNodeToSite(data: $data) {
    success
  }
}
    `;
export type AddNodeToSiteMutationFn = Apollo.MutationFunction<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;

/**
 * __useAddNodeToSiteMutation__
 *
 * To run a mutation, you first call `useAddNodeToSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeToSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeToSiteMutation, { data, loading, error }] = useAddNodeToSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeToSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>(AddNodeToSiteDocument, options);
      }
export type AddNodeToSiteMutationHookResult = ReturnType<typeof useAddNodeToSiteMutation>;
export type AddNodeToSiteMutationResult = Apollo.MutationResult<AddNodeToSiteMutation>;
export type AddNodeToSiteMutationOptions = Apollo.BaseMutationOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;
export const UpdateNodeStateDocument = gql`
    mutation updateNodeState($data: UpdateNodeStateInput!) {
  updateNodeState(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeStateMutationFn = Apollo.MutationFunction<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;

/**
 * __useUpdateNodeStateMutation__
 *
 * To run a mutation, you first call `useUpdateNodeStateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeStateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeStateMutation, { data, loading, error }] = useUpdateNodeStateMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeStateMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>(UpdateNodeStateDocument, options);
      }
export type UpdateNodeStateMutationHookResult = ReturnType<typeof useUpdateNodeStateMutation>;
export type UpdateNodeStateMutationResult = Apollo.MutationResult<UpdateNodeStateMutation>;
export type UpdateNodeStateMutationOptions = Apollo.BaseMutationOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;
export const GetNodesForSiteDocument = gql`
    query getNodesForSite($siteId: String!) {
  getNodesForSite(siteId: $siteId) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesForSiteQuery__
 *
 * To run a query within a React component, call `useGetNodesForSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesForSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesForSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetNodesForSiteQuery(baseOptions: Apollo.QueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables> & ({ variables: GetNodesForSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
      }
export function useGetNodesForSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
// @ts-ignore
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesForSiteQuery | undefined, GetNodesForSiteQueryVariables>;
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
export type GetNodesForSiteQueryHookResult = ReturnType<typeof useGetNodesForSiteQuery>;
export type GetNodesForSiteLazyQueryHookResult = ReturnType<typeof useGetNodesForSiteLazyQuery>;
export type GetNodesForSiteSuspenseQueryHookResult = ReturnType<typeof useGetNodesForSiteSuspenseQuery>;
export type GetNodesForSiteQueryResult = Apollo.QueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
export const UpdateNodeDocument = gql`
    mutation UpdateNode($data: UpdateNodeInput!) {
  updateNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeMutationFn = Apollo.MutationFunction<UpdateNodeMutation, UpdateNodeMutationVariables>;

/**
 * __useUpdateNodeMutation__
 *
 * To run a mutation, you first call `useUpdateNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeMutation, { data, loading, error }] = useUpdateNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeMutation, UpdateNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeMutation, UpdateNodeMutationVariables>(UpdateNodeDocument, options);
      }
export type UpdateNodeMutationHookResult = ReturnType<typeof useUpdateNodeMutation>;
export type UpdateNodeMutationResult = Apollo.MutationResult<UpdateNodeMutation>;
export type UpdateNodeMutationOptions = Apollo.BaseMutationOptions<UpdateNodeMutation, UpdateNodeMutationVariables>;
export const GetNodeAppsDocument = gql`
    query getNodeApps($data: NodeAppsChangeLogInput!) {
  getNodeApps(data: $data) {
    apps {
      name
      date
      version
      cpu
      memory
      notes
    }
    type
  }
}
    `;

/**
 * __useGetNodeAppsQuery__
 *
 * To run a query within a React component, call `useGetNodeAppsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeAppsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeAppsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeAppsQuery(baseOptions: Apollo.QueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables> & ({ variables: GetNodeAppsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
      }
export function useGetNodeAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
// @ts-ignore
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeAppsQuery | undefined, GetNodeAppsQueryVariables>;
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
export type GetNodeAppsQueryHookResult = ReturnType<typeof useGetNodeAppsQuery>;
export type GetNodeAppsLazyQueryHookResult = ReturnType<typeof useGetNodeAppsLazyQuery>;
export type GetNodeAppsSuspenseQueryHookResult = ReturnType<typeof useGetNodeAppsSuspenseQuery>;
export type GetNodeAppsQueryResult = Apollo.QueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export const GetNodeStateDocument = gql`
    query GetNodeState($getNodeStateId: String!) {
  getNodeState(id: $getNodeStateId) {
    id
    nodeId
    previousStateId
    previousState
    currentState
    createdAt
  }
}
    `;

/**
 * __useGetNodeStateQuery__
 *
 * To run a query within a React component, call `useGetNodeStateQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeStateQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeStateQuery({
 *   variables: {
 *      getNodeStateId: // value for 'getNodeStateId'
 *   },
 * });
 */
export function useGetNodeStateQuery(baseOptions: Apollo.QueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables> & ({ variables: GetNodeStateQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
      }
export function useGetNodeStateLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
// @ts-ignore
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeStateQuery, GetNodeStateQueryVariables>;
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeStateQuery | undefined, GetNodeStateQueryVariables>;
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
export type GetNodeStateQueryHookResult = ReturnType<typeof useGetNodeStateQuery>;
export type GetNodeStateLazyQueryHookResult = ReturnType<typeof useGetNodeStateLazyQuery>;
export type GetNodeStateSuspenseQueryHookResult = ReturnType<typeof useGetNodeStateSuspenseQuery>;
export type GetNodeStateQueryResult = Apollo.QueryResult<GetNodeStateQuery, GetNodeStateQueryVariables>;
export const RestartNodeDocument = gql`
    mutation RestartNode($data: RestartNodeInputDto!) {
  restartNode(data: $data) {
    success
  }
}
    `;
export type RestartNodeMutationFn = Apollo.MutationFunction<RestartNodeMutation, RestartNodeMutationVariables>;

/**
 * __useRestartNodeMutation__
 *
 * To run a mutation, you first call `useRestartNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestartNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restartNodeMutation, { data, loading, error }] = useRestartNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRestartNodeMutation(baseOptions?: Apollo.MutationHookOptions<RestartNodeMutation, RestartNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RestartNodeMutation, RestartNodeMutationVariables>(RestartNodeDocument, options);
      }
export type RestartNodeMutationHookResult = ReturnType<typeof useRestartNodeMutation>;
export type RestartNodeMutationResult = Apollo.MutationResult<RestartNodeMutation>;
export type RestartNodeMutationOptions = Apollo.BaseMutationOptions<RestartNodeMutation, RestartNodeMutationVariables>;
export const ToggleInternetSwitchDocument = gql`
    mutation ToggleInternetSwitch($data: ToggleInternetSwitchInputDto!) {
  toggleInternetSwitch(data: $data) {
    success
  }
}
    `;
export type ToggleInternetSwitchMutationFn = Apollo.MutationFunction<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;

/**
 * __useToggleInternetSwitchMutation__
 *
 * To run a mutation, you first call `useToggleInternetSwitchMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleInternetSwitchMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleInternetSwitchMutation, { data, loading, error }] = useToggleInternetSwitchMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleInternetSwitchMutation(baseOptions?: Apollo.MutationHookOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>(ToggleInternetSwitchDocument, options);
      }
export type ToggleInternetSwitchMutationHookResult = ReturnType<typeof useToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationResult = Apollo.MutationResult<ToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationOptions = Apollo.BaseMutationOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;
export const ToggleRfStatusDocument = gql`
    mutation ToggleRFStatus($data: ToggleRFStatusInputDto!) {
  toggleRFStatus(data: $data) {
    success
  }
}
    `;
export type ToggleRfStatusMutationFn = Apollo.MutationFunction<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;

/**
 * __useToggleRfStatusMutation__
 *
 * To run a mutation, you first call `useToggleRfStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleRfStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleRfStatusMutation, { data, loading, error }] = useToggleRfStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleRfStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>(ToggleRfStatusDocument, options);
      }
export type ToggleRfStatusMutationHookResult = ReturnType<typeof useToggleRfStatusMutation>;
export type ToggleRfStatusMutationResult = Apollo.MutationResult<ToggleRfStatusMutation>;
export type ToggleRfStatusMutationOptions = Apollo.BaseMutationOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;
export const ToggleServiceDocument = gql`
    mutation ToggleService($data: ToggleRFStatusInputDto!) {
  toggleService(data: $data) {
    success
  }
}
    `;
export type ToggleServiceMutationFn = Apollo.MutationFunction<ToggleServiceMutation, ToggleServiceMutationVariables>;

/**
 * __useToggleServiceMutation__
 *
 * To run a mutation, you first call `useToggleServiceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleServiceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleServiceMutation, { data, loading, error }] = useToggleServiceMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleServiceMutation(baseOptions?: Apollo.MutationHookOptions<ToggleServiceMutation, ToggleServiceMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleServiceMutation, ToggleServiceMutationVariables>(ToggleServiceDocument, options);
      }
export type ToggleServiceMutationHookResult = ReturnType<typeof useToggleServiceMutation>;
export type ToggleServiceMutationResult = Apollo.MutationResult<ToggleServiceMutation>;
export type ToggleServiceMutationOptions = Apollo.BaseMutationOptions<ToggleServiceMutation, ToggleServiceMutationVariables>;
export const GetHealthReportDocument = gql`
    query GetHealthReport($data: GetHealthReportInputDto!) {
  getHealthReport(data: $data) {
    id
    nodeId
    timestamp
    system {
      id
      healthId
      name
      value
    }
    capps {
      id
      space
      name
      tag
      status
      resources {
        id
        cappId
        name
        value
      }
    }
  }
}
    `;

/**
 * __useGetHealthReportQuery__
 *
 * To run a query within a React component, call `useGetHealthReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetHealthReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetHealthReportQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetHealthReportQuery(baseOptions: Apollo.QueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables> & ({ variables: GetHealthReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
      }
export function useGetHealthReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
        }
// @ts-ignore
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetHealthReportQuery, GetHealthReportQueryVariables>;
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetHealthReportQuery | undefined, GetHealthReportQueryVariables>;
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
        }
export type GetHealthReportQueryHookResult = ReturnType<typeof useGetHealthReportQuery>;
export type GetHealthReportLazyQueryHookResult = ReturnType<typeof useGetHealthReportLazyQuery>;
export type GetHealthReportSuspenseQueryHookResult = ReturnType<typeof useGetHealthReportSuspenseQuery>;
export type GetHealthReportQueryResult = Apollo.QueryResult<GetHealthReportQuery, GetHealthReportQueryVariables>;