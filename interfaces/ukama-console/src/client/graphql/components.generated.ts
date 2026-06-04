import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type UComponentFragment = { __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string };

export type GetComponentByIdQueryVariables = Types.Exact<{
  componentId: Types.Scalars['String']['input'];
}>;


export type GetComponentByIdQuery = { __typename?: 'Query', getComponentById: { __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string } };

export type GetComponentsByUserIdQueryVariables = Types.Exact<{
  data: Types.ComponentTypeInputDto;
}>;


export type GetComponentsByUserIdQuery = { __typename?: 'Query', getComponentsByUserId: { __typename?: 'ComponentsResDto', components: Array<{ __typename?: 'ComponentDto', id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string }> } };

export const UComponentFragmentDoc = gql`
    fragment UComponent on ComponentDto {
  id
  inventoryId
  type
  userId
  description
  category
  datasheetUrl
  imageUrl
  partNumber
  manufacturer
  managed
  warranty
  specification
}
    `;
export const GetComponentByIdDocument = gql`
    query getComponentById($componentId: String!) {
  getComponentById(componentId: $componentId) {
    ...UComponent
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentByIdQuery__
 *
 * To run a query within a React component, call `useGetComponentByIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentByIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentByIdQuery({
 *   variables: {
 *      componentId: // value for 'componentId'
 *   },
 * });
 */
export function useGetComponentByIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables> & ({ variables: GetComponentByIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
      }
export function useGetComponentByIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
// @ts-ignore
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentByIdQuery, GetComponentByIdQueryVariables>;
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentByIdQuery | undefined, GetComponentByIdQueryVariables>;
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
export type GetComponentByIdQueryHookResult = ReturnType<typeof useGetComponentByIdQuery>;
export type GetComponentByIdLazyQueryHookResult = ReturnType<typeof useGetComponentByIdLazyQuery>;
export type GetComponentByIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentByIdSuspenseQuery>;
export type GetComponentByIdQueryResult = Apollo.QueryResult<GetComponentByIdQuery, GetComponentByIdQueryVariables>;
export const GetComponentsByUserIdDocument = gql`
    query GetComponentsByUserId($data: ComponentTypeInputDto!) {
  getComponentsByUserId(data: $data) {
    components {
      ...UComponent
    }
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentsByUserIdQuery__
 *
 * To run a query within a React component, call `useGetComponentsByUserIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentsByUserIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentsByUserIdQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetComponentsByUserIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables> & ({ variables: GetComponentsByUserIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
      }
export function useGetComponentsByUserIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
// @ts-ignore
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>;
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentsByUserIdQuery | undefined, GetComponentsByUserIdQueryVariables>;
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
export type GetComponentsByUserIdQueryHookResult = ReturnType<typeof useGetComponentsByUserIdQuery>;
export type GetComponentsByUserIdLazyQueryHookResult = ReturnType<typeof useGetComponentsByUserIdLazyQuery>;
export type GetComponentsByUserIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentsByUserIdSuspenseQuery>;
export type GetComponentsByUserIdQueryResult = Apollo.QueryResult<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>;