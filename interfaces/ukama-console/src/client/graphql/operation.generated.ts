import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type OperationFieldsFragment = { __typename?: 'OperationDto', id: string, type: string, status: string, requestedBy?: string | null, startedAt?: string | null, leaseExpiresAt?: string | null };

export type GetNodeOperationStatusQueryVariables = Types.Exact<{
  nodeId: Types.Scalars['String']['input'];
}>;


export type GetNodeOperationStatusQuery = { __typename?: 'Query', getNodeOperationStatus: { __typename?: 'NodeOperationStatusDto', nodeId: string, busy: boolean, operation?: { __typename?: 'OperationDto', id: string, type: string, status: string, requestedBy?: string | null, startedAt?: string | null, leaseExpiresAt?: string | null } | null } };

export type GetSiteOperationStatusQueryVariables = Types.Exact<{
  siteId: Types.Scalars['String']['input'];
}>;


export type GetSiteOperationStatusQuery = { __typename?: 'Query', getSiteOperationStatus: { __typename?: 'SiteOperationStatusDto', siteId: string, busy: boolean, degraded: boolean, nodes: Array<{ __typename?: 'NodeOperationStatusDto', nodeId: string, type?: Types.NodeTypeEnum | null, busy: boolean, operation?: { __typename?: 'OperationDto', id: string, type: string, status: string, requestedBy?: string | null, startedAt?: string | null, leaseExpiresAt?: string | null } | null }>, actions: { __typename?: 'SiteActionsDto', restartSite: { __typename?: 'ActionAvailabilityDto', available: boolean, reason?: string | null }, rf: { __typename?: 'ActionAvailabilityDto', available: boolean, reason?: string | null }, service: { __typename?: 'ActionAvailabilityDto', available: boolean, reason?: string | null } } } };

export const OperationFieldsFragmentDoc = gql`
    fragment OperationFields on OperationDto {
  id
  type
  status
  requestedBy
  startedAt
  leaseExpiresAt
}
    `;
export const GetNodeOperationStatusDocument = gql`
    query GetNodeOperationStatus($nodeId: String!) {
  getNodeOperationStatus(data: {nodeId: $nodeId}) {
    nodeId
    busy
    operation {
      ...OperationFields
    }
  }
}
    ${OperationFieldsFragmentDoc}`;

/**
 * __useGetNodeOperationStatusQuery__
 *
 * To run a query within a React component, call `useGetNodeOperationStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeOperationStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeOperationStatusQuery({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *   },
 * });
 */
export function useGetNodeOperationStatusQuery(baseOptions: Apollo.QueryHookOptions<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables> & ({ variables: GetNodeOperationStatusQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>(GetNodeOperationStatusDocument, options);
      }
export function useGetNodeOperationStatusLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>(GetNodeOperationStatusDocument, options);
        }
// @ts-ignore
export function useGetNodeOperationStatusSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>;
export function useGetNodeOperationStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeOperationStatusQuery | undefined, GetNodeOperationStatusQueryVariables>;
export function useGetNodeOperationStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>(GetNodeOperationStatusDocument, options);
        }
export type GetNodeOperationStatusQueryHookResult = ReturnType<typeof useGetNodeOperationStatusQuery>;
export type GetNodeOperationStatusLazyQueryHookResult = ReturnType<typeof useGetNodeOperationStatusLazyQuery>;
export type GetNodeOperationStatusSuspenseQueryHookResult = ReturnType<typeof useGetNodeOperationStatusSuspenseQuery>;
export type GetNodeOperationStatusQueryResult = Apollo.QueryResult<GetNodeOperationStatusQuery, GetNodeOperationStatusQueryVariables>;
export const GetSiteOperationStatusDocument = gql`
    query GetSiteOperationStatus($siteId: String!) {
  getSiteOperationStatus(data: {siteId: $siteId}) {
    siteId
    busy
    degraded
    nodes {
      nodeId
      type
      busy
      operation {
        ...OperationFields
      }
    }
    actions {
      restartSite {
        available
        reason
      }
      rf {
        available
        reason
      }
      service {
        available
        reason
      }
    }
  }
}
    ${OperationFieldsFragmentDoc}`;

/**
 * __useGetSiteOperationStatusQuery__
 *
 * To run a query within a React component, call `useGetSiteOperationStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteOperationStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteOperationStatusQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetSiteOperationStatusQuery(baseOptions: Apollo.QueryHookOptions<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables> & ({ variables: GetSiteOperationStatusQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>(GetSiteOperationStatusDocument, options);
      }
export function useGetSiteOperationStatusLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>(GetSiteOperationStatusDocument, options);
        }
// @ts-ignore
export function useGetSiteOperationStatusSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>;
export function useGetSiteOperationStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteOperationStatusQuery | undefined, GetSiteOperationStatusQueryVariables>;
export function useGetSiteOperationStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>(GetSiteOperationStatusDocument, options);
        }
export type GetSiteOperationStatusQueryHookResult = ReturnType<typeof useGetSiteOperationStatusQuery>;
export type GetSiteOperationStatusLazyQueryHookResult = ReturnType<typeof useGetSiteOperationStatusLazyQuery>;
export type GetSiteOperationStatusSuspenseQueryHookResult = ReturnType<typeof useGetSiteOperationStatusSuspenseQuery>;
export type GetSiteOperationStatusQueryResult = Apollo.QueryResult<GetSiteOperationStatusQuery, GetSiteOperationStatusQueryVariables>;