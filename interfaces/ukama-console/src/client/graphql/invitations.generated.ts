import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type InvitationFragment = { __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status };

export type CreateInvitationMutationVariables = Types.Exact<{
  data: Types.CreateInvitationInputDto;
}>;


export type CreateInvitationMutation = { __typename?: 'Mutation', createInvitation: { __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status } };

export type GetInvitationsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetInvitationsQuery = { __typename?: 'Query', getInvitations: { __typename?: 'InvitationsResDto', invitations: Array<{ __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status }> } };

export type DeleteInvitationMutationVariables = Types.Exact<{
  deleteInvitationId: Types.Scalars['String']['input'];
}>;


export type DeleteInvitationMutation = { __typename?: 'Mutation', deleteInvitation: { __typename?: 'DeleteInvitationResDto', id: string } };

export type UpdateInvitationMutationVariables = Types.Exact<{
  data: Types.UpateInvitationInputDto;
}>;


export type UpdateInvitationMutation = { __typename?: 'Mutation', updateInvitation: { __typename?: 'UpdateInvitationResDto', id: string } };

export type GetInvitationsByEmailQueryVariables = Types.Exact<{
  email: Types.Scalars['String']['input'];
}>;


export type GetInvitationsByEmailQuery = { __typename?: 'Query', getInvitationsByEmail: { __typename?: 'InvitationsResDto', invitations: Array<{ __typename?: 'InvitationDto', email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status }> } };

export const InvitationFragmentDoc = gql`
    fragment Invitation on InvitationDto {
  email
  expireAt
  id
  name
  role
  link
  userId
  status
}
    `;
export const CreateInvitationDocument = gql`
    mutation CreateInvitation($data: CreateInvitationInputDto!) {
  createInvitation(data: $data) {
    ...Invitation
  }
}
    ${InvitationFragmentDoc}`;
export type CreateInvitationMutationFn = Apollo.MutationFunction<CreateInvitationMutation, CreateInvitationMutationVariables>;

/**
 * __useCreateInvitationMutation__
 *
 * To run a mutation, you first call `useCreateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createInvitationMutation, { data, loading, error }] = useCreateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useCreateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<CreateInvitationMutation, CreateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateInvitationMutation, CreateInvitationMutationVariables>(CreateInvitationDocument, options);
      }
export type CreateInvitationMutationHookResult = ReturnType<typeof useCreateInvitationMutation>;
export type CreateInvitationMutationResult = Apollo.MutationResult<CreateInvitationMutation>;
export type CreateInvitationMutationOptions = Apollo.BaseMutationOptions<CreateInvitationMutation, CreateInvitationMutationVariables>;
export const GetInvitationsDocument = gql`
    query GetInvitations {
  getInvitations {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsQuery__
 *
 * To run a query within a React component, call `useGetInvitationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetInvitationsQuery(baseOptions?: Apollo.QueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
      }
export function useGetInvitationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
// @ts-ignore
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsQuery, GetInvitationsQueryVariables>;
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsQuery | undefined, GetInvitationsQueryVariables>;
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
export type GetInvitationsQueryHookResult = ReturnType<typeof useGetInvitationsQuery>;
export type GetInvitationsLazyQueryHookResult = ReturnType<typeof useGetInvitationsLazyQuery>;
export type GetInvitationsSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsSuspenseQuery>;
export type GetInvitationsQueryResult = Apollo.QueryResult<GetInvitationsQuery, GetInvitationsQueryVariables>;
export const DeleteInvitationDocument = gql`
    mutation DeleteInvitation($deleteInvitationId: String!) {
  deleteInvitation(id: $deleteInvitationId) {
    id
  }
}
    `;
export type DeleteInvitationMutationFn = Apollo.MutationFunction<DeleteInvitationMutation, DeleteInvitationMutationVariables>;

/**
 * __useDeleteInvitationMutation__
 *
 * To run a mutation, you first call `useDeleteInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteInvitationMutation, { data, loading, error }] = useDeleteInvitationMutation({
 *   variables: {
 *      deleteInvitationId: // value for 'deleteInvitationId'
 *   },
 * });
 */
export function useDeleteInvitationMutation(baseOptions?: Apollo.MutationHookOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteInvitationMutation, DeleteInvitationMutationVariables>(DeleteInvitationDocument, options);
      }
export type DeleteInvitationMutationHookResult = ReturnType<typeof useDeleteInvitationMutation>;
export type DeleteInvitationMutationResult = Apollo.MutationResult<DeleteInvitationMutation>;
export type DeleteInvitationMutationOptions = Apollo.BaseMutationOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>;
export const UpdateInvitationDocument = gql`
    mutation UpdateInvitation($data: UpateInvitationInputDto!) {
  updateInvitation(data: $data) {
    id
  }
}
    `;
export type UpdateInvitationMutationFn = Apollo.MutationFunction<UpdateInvitationMutation, UpdateInvitationMutationVariables>;

/**
 * __useUpdateInvitationMutation__
 *
 * To run a mutation, you first call `useUpdateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateInvitationMutation, { data, loading, error }] = useUpdateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateInvitationMutation, UpdateInvitationMutationVariables>(UpdateInvitationDocument, options);
      }
export type UpdateInvitationMutationHookResult = ReturnType<typeof useUpdateInvitationMutation>;
export type UpdateInvitationMutationResult = Apollo.MutationResult<UpdateInvitationMutation>;
export type UpdateInvitationMutationOptions = Apollo.BaseMutationOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>;
export const GetInvitationsByEmailDocument = gql`
    query GetInvitationsByEmail($email: String!) {
  getInvitationsByEmail(email: $email) {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsByEmailQuery__
 *
 * To run a query within a React component, call `useGetInvitationsByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetInvitationsByEmailQuery(baseOptions: Apollo.QueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables> & ({ variables: GetInvitationsByEmailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
      }
export function useGetInvitationsByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
// @ts-ignore
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>;
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsByEmailQuery | undefined, GetInvitationsByEmailQueryVariables>;
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
export type GetInvitationsByEmailQueryHookResult = ReturnType<typeof useGetInvitationsByEmailQuery>;
export type GetInvitationsByEmailLazyQueryHookResult = ReturnType<typeof useGetInvitationsByEmailLazyQuery>;
export type GetInvitationsByEmailSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsByEmailSuspenseQuery>;
export type GetInvitationsByEmailQueryResult = Apollo.QueryResult<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>;