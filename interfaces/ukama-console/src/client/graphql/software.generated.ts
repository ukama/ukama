import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type SoftwareQueryVariables = Types.Exact<{
  data: Types.GetSoftwaresInput;
}>;


export type SoftwareQuery = { __typename?: 'Query', getSoftwares: { __typename?: 'Softwares', software: Array<{ __typename?: 'Software', id: string, releaseDate: string, nodeId: string, status: Types.SoftwareStatusEnum, changeLog: Array<string>, currentVersion: string, desiredVersion: string, name: string, space: string, notes: string, metricsKeys: Array<string>, createdAt: string, updatedAt: string }> } };

export type UpdateSoftwareMutationVariables = Types.Exact<{
  data: Types.UpdateSoftwareInputDto;
}>;


export type UpdateSoftwareMutation = { __typename?: 'Mutation', updateSoftware: { __typename?: 'StringResponse', message: string } };

export type GetAppsQueryVariables = Types.Exact<{
  data: Types.GetAppsInputDto;
}>;


export type GetAppsQuery = { __typename?: 'Query', getApps?: { __typename?: 'Apps', apps: Array<{ __typename?: 'App', name: string, version: string, tag: string, status: string, resource?: { __typename?: 'AppResource', cpuPercent: number, memoryRssKb: number, diskReadBytes: number, diskWriteBytes: number } | null }> } | null };


export const SoftwareDocument = gql`
    query Software($data: GetSoftwaresInput!) {
  getSoftwares(data: $data) {
    software {
      id
      releaseDate
      nodeId
      status
      changeLog
      currentVersion
      desiredVersion
      name
      space
      notes
      metricsKeys
      createdAt
      updatedAt
    }
  }
}
    `;

/**
 * __useSoftwareQuery__
 *
 * To run a query within a React component, call `useSoftwareQuery` and pass it any options that fit your needs.
 * When your component renders, `useSoftwareQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSoftwareQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSoftwareQuery(baseOptions: Apollo.QueryHookOptions<SoftwareQuery, SoftwareQueryVariables> & ({ variables: SoftwareQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
      }
export function useSoftwareLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
        }
// @ts-ignore
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>): Apollo.UseSuspenseQueryResult<SoftwareQuery, SoftwareQueryVariables>;
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>): Apollo.UseSuspenseQueryResult<SoftwareQuery | undefined, SoftwareQueryVariables>;
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
        }
export type SoftwareQueryHookResult = ReturnType<typeof useSoftwareQuery>;
export type SoftwareLazyQueryHookResult = ReturnType<typeof useSoftwareLazyQuery>;
export type SoftwareSuspenseQueryHookResult = ReturnType<typeof useSoftwareSuspenseQuery>;
export type SoftwareQueryResult = Apollo.QueryResult<SoftwareQuery, SoftwareQueryVariables>;
export const UpdateSoftwareDocument = gql`
    mutation UpdateSoftware($data: UpdateSoftwareInputDto!) {
  updateSoftware(data: $data) {
    message
  }
}
    `;
export type UpdateSoftwareMutationFn = Apollo.MutationFunction<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>;

/**
 * __useUpdateSoftwareMutation__
 *
 * To run a mutation, you first call `useUpdateSoftwareMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSoftwareMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSoftwareMutation, { data, loading, error }] = useUpdateSoftwareMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSoftwareMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>(UpdateSoftwareDocument, options);
      }
export type UpdateSoftwareMutationHookResult = ReturnType<typeof useUpdateSoftwareMutation>;
export type UpdateSoftwareMutationResult = Apollo.MutationResult<UpdateSoftwareMutation>;
export type UpdateSoftwareMutationOptions = Apollo.BaseMutationOptions<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>;
export const GetAppsDocument = gql`
    query GetApps($data: GetAppsInputDto!) {
  getApps(data: $data) {
    apps {
      name
      version
      tag
      status
      resource {
        cpuPercent
        memoryRssKb
        diskReadBytes
        diskWriteBytes
      }
    }
  }
}
    `;

/**
 * __useGetAppsQuery__
 *
 * To run a query within a React component, call `useGetAppsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetAppsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetAppsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetAppsQuery(baseOptions: Apollo.QueryHookOptions<GetAppsQuery, GetAppsQueryVariables> & ({ variables: GetAppsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
      }
export function useGetAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
        }
// @ts-ignore
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetAppsQuery, GetAppsQueryVariables>;
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetAppsQuery | undefined, GetAppsQueryVariables>;
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
        }
export type GetAppsQueryHookResult = ReturnType<typeof useGetAppsQuery>;
export type GetAppsLazyQueryHookResult = ReturnType<typeof useGetAppsLazyQuery>;
export type GetAppsSuspenseQueryHookResult = ReturnType<typeof useGetAppsSuspenseQuery>;
export type GetAppsQueryResult = Apollo.QueryResult<GetAppsQuery, GetAppsQueryVariables>;