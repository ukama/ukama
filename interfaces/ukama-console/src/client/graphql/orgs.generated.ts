import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type OrgFragment = { __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean };

export type GetOrgsQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetOrgsQuery = { __typename?: 'Query', getOrgs: { __typename?: 'OrgsResDto', user: string, ownerOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }>, memberOf: Array<{ __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }> } };

export type GetOrgQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetOrgQuery = { __typename?: 'Query', getOrg: { __typename?: 'OrgDto', id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean } };

export const OrgFragmentDoc = gql`
    fragment Org on OrgDto {
  id
  name
  owner
  country
  currency
  createdAt
  certificate
  isDeactivated
}
    `;
export const GetOrgsDocument = gql`
    query getOrgs {
  getOrgs {
    user
    ownerOf {
      ...Org
    }
    memberOf {
      ...Org
    }
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgsQuery__
 *
 * To run a query within a React component, call `useGetOrgsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgsQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
      }
export function useGetOrgsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
// @ts-ignore
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgsQuery, GetOrgsQueryVariables>;
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgsQuery | undefined, GetOrgsQueryVariables>;
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
export type GetOrgsQueryHookResult = ReturnType<typeof useGetOrgsQuery>;
export type GetOrgsLazyQueryHookResult = ReturnType<typeof useGetOrgsLazyQuery>;
export type GetOrgsSuspenseQueryHookResult = ReturnType<typeof useGetOrgsSuspenseQuery>;
export type GetOrgsQueryResult = Apollo.QueryResult<GetOrgsQuery, GetOrgsQueryVariables>;
export const GetOrgDocument = gql`
    query getOrg {
  getOrg {
    ...Org
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgQuery__
 *
 * To run a query within a React component, call `useGetOrgQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
      }
export function useGetOrgLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
// @ts-ignore
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgQuery, GetOrgQueryVariables>;
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgQuery | undefined, GetOrgQueryVariables>;
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
export type GetOrgQueryHookResult = ReturnType<typeof useGetOrgQuery>;
export type GetOrgLazyQueryHookResult = ReturnType<typeof useGetOrgLazyQuery>;
export type GetOrgSuspenseQueryHookResult = ReturnType<typeof useGetOrgSuspenseQuery>;
export type GetOrgQueryResult = Apollo.QueryResult<GetOrgQuery, GetOrgQueryVariables>;