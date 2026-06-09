import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type UNetworkFragment = { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> };

export type GetNetworksQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { __typename?: 'Query', getNetworks: { __typename?: 'NetworksResDto', networks: Array<{ __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> }> } };

export type GetNetworkQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
}>;


export type GetNetworkQuery = { __typename?: 'Query', getNetwork: { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type AddNetworkMutationVariables = Types.Exact<{
  data: Types.AddNetworkInputDto;
}>;


export type AddNetworkMutation = { __typename?: 'Mutation', addNetwork: { __typename?: 'NetworkDto', id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type SetDefaultNetworkMutationVariables = Types.Exact<{
  data: Types.SetDefaultNetworkInputDto;
}>;


export type SetDefaultNetworkMutation = { __typename?: 'Mutation', setDefaultNetwork: { __typename?: 'CBooleanResponse', success: boolean } };

export const UNetworkFragmentDoc = gql`
    fragment UNetwork on NetworkDto {
  id
  name
  isDefault
  budget
  overdraft
  trafficPolicy
  isDeactivated
  paymentLinks
  createdAt
  countries
  networks
}
    `;
export const GetNetworksDocument = gql`
    query getNetworks {
  getNetworks {
    networks {
      ...UNetwork
    }
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworksQuery__
 *
 * To run a query within a React component, call `useGetNetworksQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworksQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworksQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNetworksQuery(baseOptions?: Apollo.QueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
      }
export function useGetNetworksLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
// @ts-ignore
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworksQuery | undefined, GetNetworksQueryVariables>;
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
export type GetNetworksQueryHookResult = ReturnType<typeof useGetNetworksQuery>;
export type GetNetworksLazyQueryHookResult = ReturnType<typeof useGetNetworksLazyQuery>;
export type GetNetworksSuspenseQueryHookResult = ReturnType<typeof useGetNetworksSuspenseQuery>;
export type GetNetworksQueryResult = Apollo.QueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export const GetNetworkDocument = gql`
    query getNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworkQuery__
 *
 * To run a query within a React component, call `useGetNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables> & ({ variables: GetNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
      }
export function useGetNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
// @ts-ignore
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworkQuery | undefined, GetNetworkQueryVariables>;
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
export type GetNetworkQueryHookResult = ReturnType<typeof useGetNetworkQuery>;
export type GetNetworkLazyQueryHookResult = ReturnType<typeof useGetNetworkLazyQuery>;
export type GetNetworkSuspenseQueryHookResult = ReturnType<typeof useGetNetworkSuspenseQuery>;
export type GetNetworkQueryResult = Apollo.QueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
export const AddNetworkDocument = gql`
    mutation AddNetwork($data: AddNetworkInputDto!) {
  addNetwork(data: $data) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;
export type AddNetworkMutationFn = Apollo.MutationFunction<AddNetworkMutation, AddNetworkMutationVariables>;

/**
 * __useAddNetworkMutation__
 *
 * To run a mutation, you first call `useAddNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNetworkMutation, { data, loading, error }] = useAddNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNetworkMutation(baseOptions?: Apollo.MutationHookOptions<AddNetworkMutation, AddNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNetworkMutation, AddNetworkMutationVariables>(AddNetworkDocument, options);
      }
export type AddNetworkMutationHookResult = ReturnType<typeof useAddNetworkMutation>;
export type AddNetworkMutationResult = Apollo.MutationResult<AddNetworkMutation>;
export type AddNetworkMutationOptions = Apollo.BaseMutationOptions<AddNetworkMutation, AddNetworkMutationVariables>;
export const SetDefaultNetworkDocument = gql`
    mutation SetDefaultNetwork($data: SetDefaultNetworkInputDto!) {
  setDefaultNetwork(data: $data) {
    success
  }
}
    `;
export type SetDefaultNetworkMutationFn = Apollo.MutationFunction<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;

/**
 * __useSetDefaultNetworkMutation__
 *
 * To run a mutation, you first call `useSetDefaultNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetDefaultNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setDefaultNetworkMutation, { data, loading, error }] = useSetDefaultNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSetDefaultNetworkMutation(baseOptions?: Apollo.MutationHookOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>(SetDefaultNetworkDocument, options);
      }
export type SetDefaultNetworkMutationHookResult = ReturnType<typeof useSetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationResult = Apollo.MutationResult<SetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationOptions = Apollo.BaseMutationOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;