import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type RestartNodeMutationVariables = Types.Exact<{
  data: Types.RestartNodeInputDto;
}>;


export type RestartNodeMutation = { __typename?: 'Mutation', restartNode: { __typename?: 'CBooleanResponse', success: boolean } };

export type RestartSiteMutationVariables = Types.Exact<{
  data: Types.RestartSiteInputDto;
}>;


export type RestartSiteMutation = { __typename?: 'Mutation', restartSite: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleRfStatusMutationVariables = Types.Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleRfStatusMutation = { __typename?: 'Mutation', toggleRFStatus: { __typename?: 'CBooleanResponse', success: boolean } };

export type ToggleServiceMutationVariables = Types.Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleServiceMutation = { __typename?: 'Mutation', toggleService: { __typename?: 'CBooleanResponse', success: boolean } };


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
export const RestartSiteDocument = gql`
    mutation RestartSite($data: RestartSiteInputDto!) {
  restartSite(data: $data) {
    success
  }
}
    `;
export type RestartSiteMutationFn = Apollo.MutationFunction<RestartSiteMutation, RestartSiteMutationVariables>;

/**
 * __useRestartSiteMutation__
 *
 * To run a mutation, you first call `useRestartSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestartSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restartSiteMutation, { data, loading, error }] = useRestartSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRestartSiteMutation(baseOptions?: Apollo.MutationHookOptions<RestartSiteMutation, RestartSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RestartSiteMutation, RestartSiteMutationVariables>(RestartSiteDocument, options);
      }
export type RestartSiteMutationHookResult = ReturnType<typeof useRestartSiteMutation>;
export type RestartSiteMutationResult = Apollo.MutationResult<RestartSiteMutation>;
export type RestartSiteMutationOptions = Apollo.BaseMutationOptions<RestartSiteMutation, RestartSiteMutationVariables>;
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