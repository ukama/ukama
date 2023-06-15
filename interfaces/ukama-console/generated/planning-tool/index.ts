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
  _Any: { input: any; output: any; }
  _FieldSet: { input: any; output: any; }
};

export type AddDraftInput = {
  lastSaved: Scalars['Float']['input'];
  name: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type Draft = {
  __typename?: 'Draft';
  events: Array<Event>;
  id: Scalars['ID']['output'];
  lastSaved: Scalars['Float']['output'];
  name: Scalars['String']['output'];
  site: Site;
  userId: Scalars['String']['output'];
};

export type Event = {
  __typename?: 'Event';
  createdAt: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  operation: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type Location = {
  __typename?: 'Location';
  address: Scalars['String']['output'];
  lat: Scalars['String']['output'];
  lng: Scalars['String']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addDraft: Draft;
  updateDraftName: Draft;
  updateEvent: Draft;
  updateSite: Draft;
};


export type MutationAddDraftArgs = {
  data: AddDraftInput;
};


export type MutationUpdateDraftNameArgs = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};


export type MutationUpdateEventArgs = {
  data: UpdateEventInput;
  draftId: Scalars['String']['input'];
};


export type MutationUpdateSiteArgs = {
  data: UpdateSiteInput;
  id: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  _service: _Service;
  getDraft: Draft;
  getDrafts: Array<Draft>;
};


export type QueryGetDraftArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetDraftsArgs = {
  userId: Scalars['String']['input'];
};

export type Site = {
  __typename?: 'Site';
  apOption: Scalars['String']['output'];
  height: Scalars['Float']['output'];
  isSetlite: Scalars['Boolean']['output'];
  location: Location;
  name: Scalars['String']['output'];
  solarUptime: Scalars['Float']['output'];
};

export type UpdateEventInput = {
  operation: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type UpdateSiteInput = {
  address: Scalars['String']['input'];
  apOption: Scalars['String']['input'];
  height: Scalars['Float']['input'];
  isSetlite: Scalars['Boolean']['input'];
  lastSaved: Scalars['Float']['input'];
  lat: Scalars['String']['input'];
  lng: Scalars['String']['input'];
  siteName: Scalars['String']['input'];
  solarUptime: Scalars['Float']['input'];
};

export type _Service = {
  __typename?: '_Service';
  sdl?: Maybe<Scalars['String']['output']>;
};

export type LocationFragment = { __typename?: 'Location', lat: string, lng: string, address: string };

export type SiteFragment = { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } };

export type EventFragment = { __typename?: 'Event', id: string, value: string, operation: string, createdAt: string };

export type DraftFragment = { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> };

export type AddDraftMutationVariables = Exact<{
  data: AddDraftInput;
}>;


export type AddDraftMutation = { __typename?: 'Mutation', addDraft: { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type GetDraftsQueryVariables = Exact<{
  userId: Scalars['String']['input'];
}>;


export type GetDraftsQuery = { __typename?: 'Query', getDrafts: Array<{ __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> }> };

export type GetDraftQueryVariables = Exact<{
  draftId: Scalars['String']['input'];
}>;


export type GetDraftQuery = { __typename?: 'Query', getDraft: { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateEventMutationVariables = Exact<{
  draftId: Scalars['String']['input'];
  data: UpdateEventInput;
}>;


export type UpdateEventMutation = { __typename?: 'Mutation', updateEvent: { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateDraftNameMutationVariables = Exact<{
  updateDraftNameId: Scalars['String']['input'];
  name: Scalars['String']['input'];
}>;


export type UpdateDraftNameMutation = { __typename?: 'Mutation', updateDraftName: { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export type UpdateSiteMutationVariables = Exact<{
  updateSiteId: Scalars['String']['input'];
  data: UpdateSiteInput;
}>;


export type UpdateSiteMutation = { __typename?: 'Mutation', updateSite: { __typename?: 'Draft', id: string, name: string, lastSaved: number, userId: string, site: { __typename?: 'Site', name: string, height: number, apOption: string, isSetlite: boolean, solarUptime: number, location: { __typename?: 'Location', lat: string, lng: string, address: string } }, events: Array<{ __typename?: 'Event', id: string, value: string, operation: string, createdAt: string }> } };

export const LocationFragmentDoc = gql`
    fragment location on Location {
  lat
  lng
  address
}
    `;
export const SiteFragmentDoc = gql`
    fragment site on Site {
  name
  height
  apOption
  isSetlite
  solarUptime
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
  lastSaved
  userId
  site {
    ...site
  }
  events {
    ...event
  }
}
    ${SiteFragmentDoc}
${EventFragmentDoc}`;
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
export const UpdateEventDocument = gql`
    mutation UpdateEvent($draftId: String!, $data: UpdateEventInput!) {
  updateEvent(draftId: $draftId, data: $data) {
    ...draft
  }
}
    ${DraftFragmentDoc}`;
export type UpdateEventMutationFn = Apollo.MutationFunction<UpdateEventMutation, UpdateEventMutationVariables>;

/**
 * __useUpdateEventMutation__
 *
 * To run a mutation, you first call `useUpdateEventMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateEventMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateEventMutation, { data, loading, error }] = useUpdateEventMutation({
 *   variables: {
 *      draftId: // value for 'draftId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateEventMutation(baseOptions?: Apollo.MutationHookOptions<UpdateEventMutation, UpdateEventMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateEventMutation, UpdateEventMutationVariables>(UpdateEventDocument, options);
      }
export type UpdateEventMutationHookResult = ReturnType<typeof useUpdateEventMutation>;
export type UpdateEventMutationResult = Apollo.MutationResult<UpdateEventMutation>;
export type UpdateEventMutationOptions = Apollo.BaseMutationOptions<UpdateEventMutation, UpdateEventMutationVariables>;
export const UpdateDraftNameDocument = gql`
    mutation UpdateDraftName($updateDraftNameId: String!, $name: String!) {
  updateDraftName(id: $updateDraftNameId, name: $name) {
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
 *      updateDraftNameId: // value for 'updateDraftNameId'
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
export const UpdateSiteDocument = gql`
    mutation UpdateSite($updateSiteId: String!, $data: UpdateSiteInput!) {
  updateSite(id: $updateSiteId, data: $data) {
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
 *      updateSiteId: // value for 'updateSiteId'
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