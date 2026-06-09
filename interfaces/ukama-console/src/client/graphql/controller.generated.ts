import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type RestartNodeMutationVariables = Types.Exact<{
  data: Types.RestartNodeInputDto;
}>;


export type RestartNodeMutation = { __typename?: 'Mutation', restartNode: { __typename?: 'CBooleanResponse', success: boolean } };


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