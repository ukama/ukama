import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string | number; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type AddDraftInput = {
  lastSaved: Scalars['Float']['input'];
  name: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type AddNodeInput = {
  id: Scalars['String']['input'];
  state: Scalars['String']['input'];
};

export type AddNodeToNetworkInput = {
  networkId: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
};

export type AttachNodeInput = {
  anodel: Scalars['String']['input'];
  anoder: Scalars['String']['input'];
  parentNode: Scalars['String']['input'];
};

export type CBooleanResponse = {
  __typename?: 'CBooleanResponse';
  success: Scalars['Boolean']['output'];
};

export type CoverageInput = {
  height: Scalars['Float']['input'];
  lat: Scalars['Float']['input'];
  lng: Scalars['Float']['input'];
  mode: Scalars['String']['input'];
};

export type DeleteDraftRes = {
  __typename?: 'DeleteDraftRes';
  id: Scalars['String']['output'];
};

export type DeleteLinkRes = {
  __typename?: 'DeleteLinkRes';
  id: Scalars['String']['output'];
};

export type DeleteNode = {
  __typename?: 'DeleteNode';
  id: Scalars['String']['output'];
};

export type DeleteSiteRes = {
  __typename?: 'DeleteSiteRes';
  id: Scalars['String']['output'];
};

export type Draft = {
  __typename?: 'Draft';
  events: Array<Event>;
  id: Scalars['ID']['output'];
  lastSaved: Scalars['Float']['output'];
  links: Array<Link>;
  name: Scalars['String']['output'];
  sites: Array<Site>;
  userId: Scalars['String']['output'];
};

export type Event = {
  __typename?: 'Event';
  createdAt: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  operation: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type GetNode = {
  __typename?: 'GetNode';
  node: Node;
};

export type GetNodes = {
  __typename?: 'GetNodes';
  nodes: Array<Node>;
};

export type Link = {
  __typename?: 'Link';
  id: Scalars['String']['output'];
  siteA: Scalars['String']['output'];
  siteB: Scalars['String']['output'];
};

export type LinkInput = {
  lastSaved: Scalars['Float']['input'];
  siteA: Scalars['String']['input'];
  siteB: Scalars['String']['input'];
};

export type Location = {
  __typename?: 'Location';
  address: Scalars['String']['output'];
  id: Scalars['String']['output'];
  lat: Scalars['String']['output'];
  lng: Scalars['String']['output'];
};

export type LocationInput = {
  address: Scalars['String']['input'];
  lastSaved: Scalars['Float']['input'];
  lat: Scalars['String']['input'];
  lng: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addDraft: Draft;
  addLink: Draft;
  addNode: GetNode;
  addNodeToNetwork: CBooleanResponse;
  addSite: Draft;
  attachNode: CBooleanResponse;
  coverage: Site;
  deleteDraft: DeleteDraftRes;
  deleteLink: DeleteLinkRes;
  deleteNodeFromOrg: DeleteNode;
  deleteSite: DeleteSiteRes;
  detachhNode: CBooleanResponse;
  releaseNodeFromNetwork: CBooleanResponse;
  updateDraftName: Draft;
  updateEvent: Event;
  updateLocation: Location;
  updateNode: GetNode;
  updateNodeState: NodeState;
  updateSite: Draft;
};


export type MutationAddDraftArgs = {
  data: AddDraftInput;
};


export type MutationAddLinkArgs = {
  data: LinkInput;
  draftId: Scalars['String']['input'];
};


export type MutationAddNodeArgs = {
  data: AddNodeInput;
};


export type MutationAddNodeToNetworkArgs = {
  data: AddNodeToNetworkInput;
};


export type MutationAddSiteArgs = {
  data: SiteInput;
  draftId: Scalars['String']['input'];
};


export type MutationAttachNodeArgs = {
  data: AttachNodeInput;
};


export type MutationCoverageArgs = {
  data: CoverageInput;
  siteId: Scalars['String']['input'];
};


export type MutationDeleteDraftArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteLinkArgs = {
  draftId: Scalars['String']['input'];
  lastSaved: Scalars['Float']['input'];
  linkId: Scalars['String']['input'];
};


export type MutationDeleteNodeFromOrgArgs = {
  data: NodeInput;
};


export type MutationDeleteSiteArgs = {
  id: Scalars['String']['input'];
};


export type MutationDetachhNodeArgs = {
  data: NodeInput;
};


export type MutationReleaseNodeFromNetworkArgs = {
  data: NodeInput;
};


export type MutationUpdateDraftNameArgs = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};


export type MutationUpdateEventArgs = {
  data: UpdateEventInput;
  eventId: Scalars['String']['input'];
};


export type MutationUpdateLocationArgs = {
  data: LocationInput;
  draftId: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeInput;
};


export type MutationUpdateNodeStateArgs = {
  data: UpdateNodeStateInput;
};


export type MutationUpdateSiteArgs = {
  data: SiteInput;
  draftId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export type Node = {
  __typename?: 'Node';
  allocated: Scalars['Boolean']['output'];
  attached: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  network: Scalars['String']['output'];
  node: Scalars['String']['output'];
  state: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type NodeInput = {
  id: Scalars['String']['input'];
};

export type NodeState = {
  __typename?: 'NodeState';
  id: Scalars['String']['output'];
  state: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  getDraft: Draft;
  getDrafts: Array<Draft>;
  getFreeNodes: GetNodes;
  getNode: GetNode;
  getNodes: GetNodes;
};


export type QueryGetDraftArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetDraftsArgs = {
  userId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  data: NodeInput;
};

export type Site = {
  __typename?: 'Site';
  apOption: Scalars['String']['output'];
  draftId: Scalars['String']['output'];
  east: Scalars['Float']['output'];
  height: Scalars['Float']['output'];
  id: Scalars['String']['output'];
  isSetlite: Scalars['Boolean']['output'];
  location: Location;
  name: Scalars['String']['output'];
  north: Scalars['Float']['output'];
  populationCovered: Scalars['Float']['output'];
  populationUrl: Scalars['String']['output'];
  solarUptime: Scalars['Float']['output'];
  south: Scalars['Float']['output'];
  status: Scalars['String']['output'];
  totalBoxesCovered: Scalars['Float']['output'];
  url: Scalars['String']['output'];
  west: Scalars['Float']['output'];
};

export type SiteInput = {
  address: Scalars['String']['input'];
  apOption: Scalars['String']['input'];
  height: Scalars['Float']['input'];
  isSetlite: Scalars['Boolean']['input'];
  lastSaved: Scalars['Float']['input'];
  lat: Scalars['String']['input'];
  lng: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
  siteName: Scalars['String']['input'];
  solarUptime: Scalars['Float']['input'];
};

export type UpdateEventInput = {
  operation: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type UpdateNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type UpdateNodeStateInput = {
  id: Scalars['String']['input'];
  state: Scalars['String']['input'];
};

export type NodeFragment = { __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string };

export type GetNodeQueryVariables = Exact<{
  data: NodeInput;
}>;


export type GetNodeQuery = { __typename?: 'Query', getNode: { __typename?: 'GetNode', node: { __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string } } };

export type GetNodesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNodesQuery = { __typename?: 'Query', getNodes: { __typename?: 'GetNodes', nodes: Array<{ __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string }> } };

export type GetFreeNodesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetFreeNodesQuery = { __typename?: 'Query', getFreeNodes: { __typename?: 'GetNodes', nodes: Array<{ __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string }> } };

export type DeleteNodeMutationVariables = Exact<{
  data: NodeInput;
}>;


export type DeleteNodeMutation = { __typename?: 'Mutation', deleteNodeFromOrg: { __typename?: 'DeleteNode', id: string } };

export type AttachNodeMutationVariables = Exact<{
  data: AttachNodeInput;
}>;


export type AttachNodeMutation = { __typename?: 'Mutation', attachNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type DetachhNodeMutationVariables = Exact<{
  data: NodeInput;
}>;


export type DetachhNodeMutation = { __typename?: 'Mutation', detachhNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type AddNodeMutationVariables = Exact<{
  data: AddNodeInput;
}>;


export type AddNodeMutation = { __typename?: 'Mutation', addNode: { __typename?: 'GetNode', node: { __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string } } };

export type ReleaseNodeFromNetworkMutationVariables = Exact<{
  data: NodeInput;
}>;


export type ReleaseNodeFromNetworkMutation = { __typename?: 'Mutation', releaseNodeFromNetwork: { __typename?: 'CBooleanResponse', success: boolean } };

export type AddNodeToNetworkMutationVariables = Exact<{
  data: AddNodeToNetworkInput;
}>;


export type AddNodeToNetworkMutation = { __typename?: 'Mutation', addNodeToNetwork: { __typename?: 'CBooleanResponse', success: boolean } };

export type UpdateNodeStateMutationVariables = Exact<{
  data: UpdateNodeStateInput;
}>;


export type UpdateNodeStateMutation = { __typename?: 'Mutation', updateNodeState: { __typename?: 'NodeState', id: string, state: string } };

export type UpdateNodeMutationVariables = Exact<{
  data: UpdateNodeInput;
}>;


export type UpdateNodeMutation = { __typename?: 'Mutation', updateNode: { __typename?: 'GetNode', node: { __typename?: 'Node', allocated: boolean, attached: Array<string>, name: string, network: string, node: string, state: string, type: string } } };

export type LocationFragment = { __typename?: 'Location', id: string, lat: string, lng: string, address: string };

export type LinkFragment = { __typename?: 'Link', id: string, siteA: string, siteB: string };

export type SiteFragment = { __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } };

export type EventFragment = { __typename?: 'Event', id: string, value: string, operation: string, createdAt: string };

export type DraftFragment = { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> };

export type AddDraftMutationVariables = Exact<{
  data: AddDraftInput;
}>;


export type AddDraftMutation = { __typename?: 'Mutation', addDraft: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateDraftNameMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  name: Scalars['String']['input'];
}>;


export type UpdateDraftNameMutation = { __typename?: 'Mutation', updateDraftName: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type GetDraftsQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetDraftsQuery = { __typename?: 'Query', getDrafts: Array<{ __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> }> };

export type GetDraftQueryVariables = Exact<{
  draftId: Scalars['String']['input'];
}>;


export type GetDraftQuery = { __typename?: 'Query', getDraft: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type AddSiteMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  data: SiteInput;
}>;


export type AddSiteMutation = { __typename?: 'Mutation', addSite: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateSiteMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  data: SiteInput;
}>;


export type UpdateSiteMutation = { __typename?: 'Mutation', updateSite: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateLocationMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  locationId: Scalars['String']['input'];
  data: LocationInput;
}>;


export type UpdateLocationMutation = { __typename?: 'Mutation', updateLocation: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } };

export type DeleteDraftMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
}>;


export type DeleteDraftMutation = { __typename?: 'Mutation', deleteDraft: { __typename?: 'DeleteDraftRes', id: string } };

export type DeleteSiteMutationVariables = Exact<{
  siteId: Scalars['String']['input'];
}>;


export type DeleteSiteMutation = { __typename?: 'Mutation', deleteSite: { __typename?: 'DeleteSiteRes', id: string } };

export type DeleteLinkMutationVariables = Exact<{
  lastSaved: Scalars['Float']['input'];
  draftId: Scalars['String']['input'];
  linkId: Scalars['String']['input'];
}>;


export type DeleteLinkMutation = { __typename?: 'Mutation', deleteLink: { __typename?: 'DeleteLinkRes', id: string } };

export type AddLinkMutationVariables = Exact<{
  data: LinkInput;
  draftId: Scalars['String']['input'];
}>;


export type AddLinkMutation = { __typename?: 'Mutation', addLink: { __typename?: 'Draft', id: string, name: string, userId: string, lastSaved: number, links: Array<{ __typename?: 'Link', id: string, siteA: string, siteB: string }>, sites: Array<{ __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } }>, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type CoverageMutationVariables = Exact<{
  siteId: Scalars['String']['input'];
  data: CoverageInput;
}>;


export type CoverageMutation = { __typename?: 'Mutation', coverage: { __typename?: 'Site', id: string, url: string, east: number, name: string, west: number, north: number, south: number, status: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, populationUrl: string, populationCovered: number, totalBoxesCovered: number, location: { __typename?: 'Location', id: string, lat: string, lng: string, address: string } } };

export const NodeFragmentDoc = gql`
    fragment node on Node {
  allocated
  attached
  name
  network
  node
  state
  type
}
    `;
export const LinkFragmentDoc = gql`
    fragment link on Link {
  id
  siteA
  siteB
}
    `;
export const LocationFragmentDoc = gql`
    fragment location on Location {
  id
  lat
  lng
  address
}
    `;
export const SiteFragmentDoc = gql`
    fragment site on Site {
  id
  url
  east
  name
  west
  north
  south
  status
  height
  apOption
  isSetlite
  solarUptime
  populationUrl
  populationCovered
  totalBoxesCovered
  location {
    ...location
  }
}
    ${LocationFragmentDoc}`;
export const EventFragmentDoc = gql`
    fragment event on Event {
  id
  value
  operation
  createdAt
}
    `;
export const DraftFragmentDoc = gql`
    fragment draft on Draft {
  id
  name
  userId
  lastSaved
  links {
    ...link
  }
  sites {
    ...site
  }
  events {
    ...event
  }
}
    ${LinkFragmentDoc}
${SiteFragmentDoc}
${EventFragmentDoc}`;
export const GetNodeDocument = gql`
    query getNode($data: NodeInput!) {
  getNode(data: $data) {
    node {
      ...node
    }
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
export function useGetNodeQuery(baseOptions: Apollo.QueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
      }
export function useGetNodeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export type GetNodeQueryHookResult = ReturnType<typeof useGetNodeQuery>;
export type GetNodeLazyQueryHookResult = ReturnType<typeof useGetNodeLazyQuery>;
export type GetNodeQueryResult = Apollo.QueryResult<GetNodeQuery, GetNodeQueryVariables>;
export const GetNodesDocument = gql`
    query getNodes {
  getNodes {
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
 *   },
 * });
 */
export function useGetNodesQuery(baseOptions?: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
      }
export function useGetNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<typeof useGetNodesLazyQuery>;
export type GetNodesQueryResult = Apollo.QueryResult<GetNodesQuery, GetNodesQueryVariables>;
export const GetFreeNodesDocument = gql`
    query getFreeNodes {
  getFreeNodes {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetFreeNodesQuery__
 *
 * To run a query within a React component, call `useGetFreeNodesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetFreeNodesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetFreeNodesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetFreeNodesQuery(baseOptions?: Apollo.QueryHookOptions<GetFreeNodesQuery, GetFreeNodesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetFreeNodesQuery, GetFreeNodesQueryVariables>(GetFreeNodesDocument, options);
      }
export function useGetFreeNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetFreeNodesQuery, GetFreeNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetFreeNodesQuery, GetFreeNodesQueryVariables>(GetFreeNodesDocument, options);
        }
export type GetFreeNodesQueryHookResult = ReturnType<typeof useGetFreeNodesQuery>;
export type GetFreeNodesLazyQueryHookResult = ReturnType<typeof useGetFreeNodesLazyQuery>;
export type GetFreeNodesQueryResult = Apollo.QueryResult<GetFreeNodesQuery, GetFreeNodesQueryVariables>;
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
    node {
      ...node
    }
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
export const ReleaseNodeFromNetworkDocument = gql`
    mutation releaseNodeFromNetwork($data: NodeInput!) {
  releaseNodeFromNetwork(data: $data) {
    success
  }
}
    `;
export type ReleaseNodeFromNetworkMutationFn = Apollo.MutationFunction<ReleaseNodeFromNetworkMutation, ReleaseNodeFromNetworkMutationVariables>;

/**
 * __useReleaseNodeFromNetworkMutation__
 *
 * To run a mutation, you first call `useReleaseNodeFromNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useReleaseNodeFromNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [releaseNodeFromNetworkMutation, { data, loading, error }] = useReleaseNodeFromNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useReleaseNodeFromNetworkMutation(baseOptions?: Apollo.MutationHookOptions<ReleaseNodeFromNetworkMutation, ReleaseNodeFromNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ReleaseNodeFromNetworkMutation, ReleaseNodeFromNetworkMutationVariables>(ReleaseNodeFromNetworkDocument, options);
      }
export type ReleaseNodeFromNetworkMutationHookResult = ReturnType<typeof useReleaseNodeFromNetworkMutation>;
export type ReleaseNodeFromNetworkMutationResult = Apollo.MutationResult<ReleaseNodeFromNetworkMutation>;
export type ReleaseNodeFromNetworkMutationOptions = Apollo.BaseMutationOptions<ReleaseNodeFromNetworkMutation, ReleaseNodeFromNetworkMutationVariables>;
export const AddNodeToNetworkDocument = gql`
    mutation addNodeToNetwork($data: AddNodeToNetworkInput!) {
  addNodeToNetwork(data: $data) {
    success
  }
}
    `;
export type AddNodeToNetworkMutationFn = Apollo.MutationFunction<AddNodeToNetworkMutation, AddNodeToNetworkMutationVariables>;

/**
 * __useAddNodeToNetworkMutation__
 *
 * To run a mutation, you first call `useAddNodeToNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeToNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeToNetworkMutation, { data, loading, error }] = useAddNodeToNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeToNetworkMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeToNetworkMutation, AddNodeToNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeToNetworkMutation, AddNodeToNetworkMutationVariables>(AddNodeToNetworkDocument, options);
      }
export type AddNodeToNetworkMutationHookResult = ReturnType<typeof useAddNodeToNetworkMutation>;
export type AddNodeToNetworkMutationResult = Apollo.MutationResult<AddNodeToNetworkMutation>;
export type AddNodeToNetworkMutationOptions = Apollo.BaseMutationOptions<AddNodeToNetworkMutation, AddNodeToNetworkMutationVariables>;
export const UpdateNodeStateDocument = gql`
    mutation updateNodeState($data: UpdateNodeStateInput!) {
  updateNodeState(data: $data) {
    id
    state
  }
}
    `;
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
export const UpdateNodeDocument = gql`
    mutation UpdateNode($data: UpdateNodeInput!) {
  updateNode(data: $data) {
    node {
      ...node
    }
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
export const AddDraftDocument = gql`
    mutation AddDraft($data: AddDraftInput!) {
  addDraft(data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddDraftMutationFn = Apollo.MutationFunction<AddDraftMutation, AddDraftMutationVariables>;

/**
 * __useAddDraftMutation__
 *
 * To run a mutation, you first call `useAddDraftMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddDraftMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addDraftMutation, { data, loading, error }] = useAddDraftMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddDraftMutation(baseOptions?: Apollo.MutationHookOptions<AddDraftMutation, AddDraftMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddDraftMutation, AddDraftMutationVariables>(AddDraftDocument, options);
      }
export type AddDraftMutationHookResult = ReturnType<typeof useAddDraftMutation>;
export type AddDraftMutationResult = Apollo.MutationResult<AddDraftMutation>;
export type AddDraftMutationOptions = Apollo.BaseMutationOptions<AddDraftMutation, AddDraftMutationVariables>;
export const UpdateDraftNameDocument = gql`
    mutation UpdateDraftName($draftId: String!, $name: String!) {
  updateDraftName(id: $draftId, name: $name) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type UpdateDraftNameMutationFn = Apollo.MutationFunction<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>;

/**
 * __useUpdateDraftNameMutation__
 *
 * To run a mutation, you first call `useUpdateDraftNameMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateDraftNameMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateDraftNameMutation, { data, loading, error }] = useUpdateDraftNameMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      name: // value for 'name'
 *   },
 * });
 */
export function useUpdateDraftNameMutation(baseOptions?: Apollo.MutationHookOptions<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>(UpdateDraftNameDocument, options);
      }
export type UpdateDraftNameMutationHookResult = ReturnType<typeof useUpdateDraftNameMutation>;
export type UpdateDraftNameMutationResult = Apollo.MutationResult<UpdateDraftNameMutation>;
export type UpdateDraftNameMutationOptions = Apollo.BaseMutationOptions<UpdateDraftNameMutation, UpdateDraftNameMutationVariables>;
export const GetDraftsDocument = gql`
    query GetDrafts($userId: String!) {
  getDrafts(userId: $userId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;

/**
 * __useGetDraftsQuery__
 *
 * To run a query within a React component, call `useGetDraftsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDraftsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDraftsQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetDraftsQuery(baseOptions: Apollo.QueryHookOptions<GetDraftsQuery, GetDraftsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDraftsQuery, GetDraftsQueryVariables>(GetDraftsDocument, options);
      }
export function useGetDraftsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDraftsQuery, GetDraftsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDraftsQuery, GetDraftsQueryVariables>(GetDraftsDocument, options);
        }
export type GetDraftsQueryHookResult = ReturnType<typeof useGetDraftsQuery>;
export type GetDraftsLazyQueryHookResult = ReturnType<typeof useGetDraftsLazyQuery>;
export type GetDraftsQueryResult = Apollo.QueryResult<GetDraftsQuery, GetDraftsQueryVariables>;
export const GetDraftDocument = gql`
    query GetDraft($draftId: String!) {
  getDraft(id: $draftId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;

/**
 * __useGetDraftQuery__
 *
 * To run a query within a React component, call `useGetDraftQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDraftQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDraftQuery({
 *   variables: {
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useGetDraftQuery(baseOptions: Apollo.QueryHookOptions<GetDraftQuery, GetDraftQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDraftQuery, GetDraftQueryVariables>(GetDraftDocument, options);
      }
export function useGetDraftLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDraftQuery, GetDraftQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDraftQuery, GetDraftQueryVariables>(GetDraftDocument, options);
        }
export type GetDraftQueryHookResult = ReturnType<typeof useGetDraftQuery>;
export type GetDraftLazyQueryHookResult = ReturnType<typeof useGetDraftLazyQuery>;
export type GetDraftQueryResult = Apollo.QueryResult<GetDraftQuery, GetDraftQueryVariables>;
export const AddSiteDocument = gql`
    mutation addSite($draftId: String!, $data: SiteInput!) {
  addSite(draftId: $draftId, data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddSiteMutationFn = Apollo.MutationFunction<AddSiteMutation, AddSiteMutationVariables>;

/**
 * __useAddSiteMutation__
 *
 * To run a mutation, you first call `useAddSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSiteMutation, { data, loading, error }] = useAddSiteMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddSiteMutation, AddSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSiteMutation, AddSiteMutationVariables>(AddSiteDocument, options);
      }
export type AddSiteMutationHookResult = ReturnType<typeof useAddSiteMutation>;
export type AddSiteMutationResult = Apollo.MutationResult<AddSiteMutation>;
export type AddSiteMutationOptions = Apollo.BaseMutationOptions<AddSiteMutation, AddSiteMutationVariables>;
export const UpdateSiteDocument = gql`
    mutation UpdateSite($draftId: String!, $siteId: String!, $data: SiteInput!) {
  updateSite(draftId: $draftId, siteId: $siteId, data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type UpdateSiteMutationFn = Apollo.MutationFunction<UpdateSiteMutation, UpdateSiteMutationVariables>;

/**
 * __useUpdateSiteMutation__
 *
 * To run a mutation, you first call `useUpdateSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSiteMutation, { data, loading, error }] = useUpdateSiteMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      siteId: // value for 'siteId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSiteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSiteMutation, UpdateSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSiteMutation, UpdateSiteMutationVariables>(UpdateSiteDocument, options);
      }
export type UpdateSiteMutationHookResult = ReturnType<typeof useUpdateSiteMutation>;
export type UpdateSiteMutationResult = Apollo.MutationResult<UpdateSiteMutation>;
export type UpdateSiteMutationOptions = Apollo.BaseMutationOptions<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const UpdateLocationDocument = gql`
    mutation UpdateLocation($draftId: String!, $locationId: String!, $data: LocationInput!) {
  updateLocation(draftId: $draftId, locationId: $locationId, data: $data) {
    ...location
  }
}
    ${LocationFragmentDoc}`;
export type UpdateLocationMutationFn = Apollo.MutationFunction<UpdateLocationMutation, UpdateLocationMutationVariables>;

/**
 * __useUpdateLocationMutation__
 *
 * To run a mutation, you first call `useUpdateLocationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateLocationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateLocationMutation, { data, loading, error }] = useUpdateLocationMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      locationId: // value for 'locationId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateLocationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateLocationMutation, UpdateLocationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateLocationMutation, UpdateLocationMutationVariables>(UpdateLocationDocument, options);
      }
export type UpdateLocationMutationHookResult = ReturnType<typeof useUpdateLocationMutation>;
export type UpdateLocationMutationResult = Apollo.MutationResult<UpdateLocationMutation>;
export type UpdateLocationMutationOptions = Apollo.BaseMutationOptions<UpdateLocationMutation, UpdateLocationMutationVariables>;
export const DeleteDraftDocument = gql`
    mutation DeleteDraft($draftId: String!) {
  deleteDraft(id: $draftId) {
    id
  }
}
    `;
export type DeleteDraftMutationFn = Apollo.MutationFunction<DeleteDraftMutation, DeleteDraftMutationVariables>;

/**
 * __useDeleteDraftMutation__
 *
 * To run a mutation, you first call `useDeleteDraftMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteDraftMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteDraftMutation, { data, loading, error }] = useDeleteDraftMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useDeleteDraftMutation(baseOptions?: Apollo.MutationHookOptions<DeleteDraftMutation, DeleteDraftMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteDraftMutation, DeleteDraftMutationVariables>(DeleteDraftDocument, options);
      }
export type DeleteDraftMutationHookResult = ReturnType<typeof useDeleteDraftMutation>;
export type DeleteDraftMutationResult = Apollo.MutationResult<DeleteDraftMutation>;
export type DeleteDraftMutationOptions = Apollo.BaseMutationOptions<DeleteDraftMutation, DeleteDraftMutationVariables>;
export const DeleteSiteDocument = gql`
    mutation DeleteSite($siteId: String!) {
  deleteSite(id: $siteId) {
    id
  }
}
    `;
export type DeleteSiteMutationFn = Apollo.MutationFunction<DeleteSiteMutation, DeleteSiteMutationVariables>;

/**
 * __useDeleteSiteMutation__
 *
 * To run a mutation, you first call `useDeleteSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSiteMutation, { data, loading, error }] = useDeleteSiteMutation({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useDeleteSiteMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSiteMutation, DeleteSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSiteMutation, DeleteSiteMutationVariables>(DeleteSiteDocument, options);
      }
export type DeleteSiteMutationHookResult = ReturnType<typeof useDeleteSiteMutation>;
export type DeleteSiteMutationResult = Apollo.MutationResult<DeleteSiteMutation>;
export type DeleteSiteMutationOptions = Apollo.BaseMutationOptions<DeleteSiteMutation, DeleteSiteMutationVariables>;
export const DeleteLinkDocument = gql`
    mutation DeleteLink($lastSaved: Float!, $draftId: String!, $linkId: String!) {
  deleteLink(lastSaved: $lastSaved, draftId: $draftId, linkId: $linkId) {
    id
  }
}
    `;
export type DeleteLinkMutationFn = Apollo.MutationFunction<DeleteLinkMutation, DeleteLinkMutationVariables>;

/**
 * __useDeleteLinkMutation__
 *
 * To run a mutation, you first call `useDeleteLinkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteLinkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteLinkMutation, { data, loading, error }] = useDeleteLinkMutation({
 *   variables: {
 *      lastSaved: // value for 'lastSaved'
 *      draftId: // value for 'draftId'
 *      linkId: // value for 'linkId'
 *   },
 * });
 */
export function useDeleteLinkMutation(baseOptions?: Apollo.MutationHookOptions<DeleteLinkMutation, DeleteLinkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteLinkMutation, DeleteLinkMutationVariables>(DeleteLinkDocument, options);
      }
export type DeleteLinkMutationHookResult = ReturnType<typeof useDeleteLinkMutation>;
export type DeleteLinkMutationResult = Apollo.MutationResult<DeleteLinkMutation>;
export type DeleteLinkMutationOptions = Apollo.BaseMutationOptions<DeleteLinkMutation, DeleteLinkMutationVariables>;
export const AddLinkDocument = gql`
    mutation AddLink($data: LinkInput!, $draftId: String!) {
  addLink(data: $data, draftId: $draftId) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type AddLinkMutationFn = Apollo.MutationFunction<AddLinkMutation, AddLinkMutationVariables>;

/**
 * __useAddLinkMutation__
 *
 * To run a mutation, you first call `useAddLinkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddLinkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addLinkMutation, { data, loading, error }] = useAddLinkMutation({
 *   variables: {
 *      data: // value for 'data'
 *      draftId: // value for 'draftId'
 *   },
 * });
 */
export function useAddLinkMutation(baseOptions?: Apollo.MutationHookOptions<AddLinkMutation, AddLinkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddLinkMutation, AddLinkMutationVariables>(AddLinkDocument, options);
      }
export type AddLinkMutationHookResult = ReturnType<typeof useAddLinkMutation>;
export type AddLinkMutationResult = Apollo.MutationResult<AddLinkMutation>;
export type AddLinkMutationOptions = Apollo.BaseMutationOptions<AddLinkMutation, AddLinkMutationVariables>;
export const CoverageDocument = gql`
    mutation Coverage($siteId: String!, $data: CoverageInput!) {
  coverage(data: $data, siteId: $siteId) {
    ...site
  }
}
    ${SiteFragmentDoc}`;
export type CoverageMutationFn = Apollo.MutationFunction<CoverageMutation, CoverageMutationVariables>;

/**
 * __useCoverageMutation__
 *
 * To run a mutation, you first call `useCoverageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCoverageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [coverageMutation, { data, loading, error }] = useCoverageMutation({
 *   variables: {
 *      siteId: // value for 'siteId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useCoverageMutation(baseOptions?: Apollo.MutationHookOptions<CoverageMutation, CoverageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CoverageMutation, CoverageMutationVariables>(CoverageDocument, options);
      }
export type CoverageMutationHookResult = ReturnType<typeof useCoverageMutation>;
export type CoverageMutationResult = Apollo.MutationResult<CoverageMutation>;
export type CoverageMutationOptions = Apollo.BaseMutationOptions<CoverageMutation, CoverageMutationVariables>;