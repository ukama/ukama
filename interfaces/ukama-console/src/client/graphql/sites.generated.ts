import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type USiteFragment = { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string };

export type GetSiteQueryVariables = Types.Exact<{
  siteId: Types.Scalars['String']['input'];
}>;


export type GetSiteQuery = { __typename?: 'Query', getSite: { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type AddSiteMutationVariables = Types.Exact<{
  data: Types.AddSiteInputDto;
}>;


export type AddSiteMutation = { __typename?: 'Mutation', addSite: { __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type GetSitesQueryVariables = Types.Exact<{
  data: Types.SitesInputDto;
}>;


export type GetSitesQuery = { __typename?: 'Query', getSites: { __typename?: 'SitesResDto', sites: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string }> } };

export type UpdateSiteMutationVariables = Types.Exact<{
  siteId: Types.Scalars['String']['input'];
  data: Types.UpdateSiteInputDto;
}>;


export type UpdateSiteMutation = { __typename?: 'Mutation', updateSite: { __typename?: 'SiteDto', name: string } };

export const USiteFragmentDoc = gql`
    fragment USite on SiteDto {
  id
  name
  networkId
  backhaulId
  powerId
  accessId
  spectrumId
  switchId
  isDeactivated
  latitude
  longitude
  installDate
  createdAt
  location
}
    `;
export const GetSiteDocument = gql`
    query getSite($siteId: String!) {
  getSite(siteId: $siteId) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;

/**
 * __useGetSiteQuery__
 *
 * To run a query within a React component, call `useGetSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetSiteQuery(baseOptions: Apollo.QueryHookOptions<GetSiteQuery, GetSiteQueryVariables> & ({ variables: GetSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
      }
export function useGetSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
// @ts-ignore
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteQuery, GetSiteQueryVariables>;
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteQuery | undefined, GetSiteQueryVariables>;
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
export type GetSiteQueryHookResult = ReturnType<typeof useGetSiteQuery>;
export type GetSiteLazyQueryHookResult = ReturnType<typeof useGetSiteLazyQuery>;
export type GetSiteSuspenseQueryHookResult = ReturnType<typeof useGetSiteSuspenseQuery>;
export type GetSiteQueryResult = Apollo.QueryResult<GetSiteQuery, GetSiteQueryVariables>;
export const AddSiteDocument = gql`
    mutation addSite($data: AddSiteInputDto!) {
  addSite(data: $data) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;
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
export const GetSitesDocument = gql`
    query GetSites($data: SitesInputDto!) {
  getSites(data: $data) {
    sites {
      ...USite
    }
  }
}
    ${USiteFragmentDoc}`;

/**
 * __useGetSitesQuery__
 *
 * To run a query within a React component, call `useGetSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSitesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSitesQuery(baseOptions: Apollo.QueryHookOptions<GetSitesQuery, GetSitesQueryVariables> & ({ variables: GetSitesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
      }
export function useGetSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
// @ts-ignore
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetSitesQuery, GetSitesQueryVariables>;
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetSitesQuery | undefined, GetSitesQueryVariables>;
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
export type GetSitesQueryHookResult = ReturnType<typeof useGetSitesQuery>;
export type GetSitesLazyQueryHookResult = ReturnType<typeof useGetSitesLazyQuery>;
export type GetSitesSuspenseQueryHookResult = ReturnType<typeof useGetSitesSuspenseQuery>;
export type GetSitesQueryResult = Apollo.QueryResult<GetSitesQuery, GetSitesQueryVariables>;
export const UpdateSiteDocument = gql`
    mutation updateSite($siteId: String!, $data: UpdateSiteInputDto!) {
  updateSite(siteId: $siteId, data: $data) {
    name
  }
}
    `;
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