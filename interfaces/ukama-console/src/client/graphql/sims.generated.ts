import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type GetSimPoolStatsQueryVariables = Types.Exact<{
  data: Types.GetSimsInput;
}>;


export type GetSimPoolStatsQuery = { __typename?: 'Query', getSimPoolStats: { __typename?: 'SimPoolStatsDto', total: number, available: number, consumed: number, failed: number, esim: number, physical: number } };

export type GetSimsFromPoolQueryVariables = Types.Exact<{
  data: Types.GetSimsInput;
}>;


export type GetSimsFromPoolQuery = { __typename?: 'Query', getSimsFromPool: { __typename?: 'SimsPoolResDto', sims: Array<{ __typename?: 'SimPoolResDto', id: string, qrCode: string, iccid: string, msisdn: string, isAllocated: boolean, isFailed: boolean, simType: string, smApAddress: string, activationCode: string, createdAt: string, deletedAt: string, updatedAt: string, isPhysical: boolean }> } };

export type UploadSimsMutationVariables = Types.Exact<{
  data: Types.UploadSimsInputDto;
}>;


export type UploadSimsMutation = { __typename?: 'Mutation', uploadSims: { __typename?: 'UploadSimsResDto', iccid: Array<string> } };

export type SimPackageFragment = { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean };

export type SimFragment = { __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null };

export type SimAllocationPackageFragment = { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null };

export type SimAllocationFragment = { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } };

export type AllocateSimMutationVariables = Types.Exact<{
  data: Types.AllocateSimInputDto;
}>;


export type AllocateSimMutation = { __typename?: 'Mutation', allocateSim: { __typename?: 'AllocateSimAPIDto', id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi?: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { __typename?: 'SimAllocatePackageDto', id?: string | null, packageId?: string | null, startDate?: string | null, endDate?: string | null, isActive?: boolean | null } } };

export type ToggleSimStatusMutationVariables = Types.Exact<{
  data: Types.ToggleSimStatusInputDto;
}>;


export type ToggleSimStatusMutation = { __typename?: 'Mutation', toggleSimStatus: { __typename?: 'SimStatusResDto', simId?: string | null } };

export type GetSimQueryVariables = Types.Exact<{
  data: Types.GetSimInputDto;
}>;


export type GetSimQuery = { __typename?: 'Query', getSim: { __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null } };

export type GetSimsQueryVariables = Types.Exact<{
  data: Types.ListSimsInput;
}>;


export type GetSimsQuery = { __typename?: 'Query', getSims: { __typename?: 'SimsResDto', sims: Array<{ __typename?: 'SimDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package?: { __typename?: 'SimPackage', id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null }> } };

export type GetDataUsagesQueryVariables = Types.Exact<{
  data: Types.SimUsagesInputDto;
}>;


export type GetDataUsagesQuery = { __typename?: 'Query', getDataUsages: { __typename?: 'SimDataUsages', usages: Array<{ __typename?: 'SimDataUsage', usage: string, simId: string }> } };

export const SimPackageFragmentDoc = gql`
    fragment SimPackage on SimPackage {
  id
  packageId
  startDate
  endDate
  defaultDuration
  isActive
  asExpired
}
    `;
export const SimFragmentDoc = gql`
    fragment Sim on SimDto {
  id
  subscriberId
  networkId
  iccid
  msisdn
  imsi
  type
  status
  isPhysical
  trafficPolicy
  firstActivatedOn
  lastActivatedOn
  activationsCount
  deactivationsCount
  allocatedAt
  syncStatus
  package {
    ...SimPackage
  }
}
    ${SimPackageFragmentDoc}`;
export const SimAllocationPackageFragmentDoc = gql`
    fragment SimAllocationPackage on SimAllocatePackageDto {
  id
  packageId
  startDate
  endDate
  isActive
}
    `;
export const SimAllocationFragmentDoc = gql`
    fragment SimAllocation on AllocateSimAPIDto {
  id
  subscriber_id
  network_id
  package {
    ...SimAllocationPackage
  }
  iccid
  msisdn
  imsi
  type
  status
  is_physical
  traffic_policy
  allocated_at
  sync_status
}
    ${SimAllocationPackageFragmentDoc}`;
export const GetSimPoolStatsDocument = gql`
    query GetSimPoolStats($data: GetSimsInput!) {
  getSimPoolStats(data: $data) {
    total
    available
    consumed
    failed
    esim
    physical
  }
}
    `;

/**
 * __useGetSimPoolStatsQuery__
 *
 * To run a query within a React component, call `useGetSimPoolStatsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimPoolStatsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimPoolStatsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimPoolStatsQuery(baseOptions: Apollo.QueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables> & ({ variables: GetSimPoolStatsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
      }
export function useGetSimPoolStatsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
// @ts-ignore
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>;
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimPoolStatsQuery | undefined, GetSimPoolStatsQueryVariables>;
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
export type GetSimPoolStatsQueryHookResult = ReturnType<typeof useGetSimPoolStatsQuery>;
export type GetSimPoolStatsLazyQueryHookResult = ReturnType<typeof useGetSimPoolStatsLazyQuery>;
export type GetSimPoolStatsSuspenseQueryHookResult = ReturnType<typeof useGetSimPoolStatsSuspenseQuery>;
export type GetSimPoolStatsQueryResult = Apollo.QueryResult<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>;
export const GetSimsFromPoolDocument = gql`
    query GetSimsFromPool($data: GetSimsInput!) {
  getSimsFromPool(data: $data) {
    sims {
      id
      qrCode
      iccid
      msisdn
      isAllocated
      isFailed
      simType
      smApAddress
      activationCode
      createdAt
      deletedAt
      updatedAt
      isPhysical
    }
  }
}
    `;

/**
 * __useGetSimsFromPoolQuery__
 *
 * To run a query within a React component, call `useGetSimsFromPoolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsFromPoolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsFromPoolQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsFromPoolQuery(baseOptions: Apollo.QueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables> & ({ variables: GetSimsFromPoolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
      }
export function useGetSimsFromPoolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
// @ts-ignore
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>;
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsFromPoolQuery | undefined, GetSimsFromPoolQueryVariables>;
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
export type GetSimsFromPoolQueryHookResult = ReturnType<typeof useGetSimsFromPoolQuery>;
export type GetSimsFromPoolLazyQueryHookResult = ReturnType<typeof useGetSimsFromPoolLazyQuery>;
export type GetSimsFromPoolSuspenseQueryHookResult = ReturnType<typeof useGetSimsFromPoolSuspenseQuery>;
export type GetSimsFromPoolQueryResult = Apollo.QueryResult<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>;
export const UploadSimsDocument = gql`
    mutation uploadSims($data: UploadSimsInputDto!) {
  uploadSims(data: $data) {
    iccid
  }
}
    `;
export type UploadSimsMutationFn = Apollo.MutationFunction<UploadSimsMutation, UploadSimsMutationVariables>;

/**
 * __useUploadSimsMutation__
 *
 * To run a mutation, you first call `useUploadSimsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUploadSimsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [uploadSimsMutation, { data, loading, error }] = useUploadSimsMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUploadSimsMutation(baseOptions?: Apollo.MutationHookOptions<UploadSimsMutation, UploadSimsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UploadSimsMutation, UploadSimsMutationVariables>(UploadSimsDocument, options);
      }
export type UploadSimsMutationHookResult = ReturnType<typeof useUploadSimsMutation>;
export type UploadSimsMutationResult = Apollo.MutationResult<UploadSimsMutation>;
export type UploadSimsMutationOptions = Apollo.BaseMutationOptions<UploadSimsMutation, UploadSimsMutationVariables>;
export const AllocateSimDocument = gql`
    mutation allocateSim($data: AllocateSimInputDto!) {
  allocateSim(data: $data) {
    ...SimAllocation
  }
}
    ${SimAllocationFragmentDoc}`;
export type AllocateSimMutationFn = Apollo.MutationFunction<AllocateSimMutation, AllocateSimMutationVariables>;

/**
 * __useAllocateSimMutation__
 *
 * To run a mutation, you first call `useAllocateSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAllocateSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [allocateSimMutation, { data, loading, error }] = useAllocateSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAllocateSimMutation(baseOptions?: Apollo.MutationHookOptions<AllocateSimMutation, AllocateSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AllocateSimMutation, AllocateSimMutationVariables>(AllocateSimDocument, options);
      }
export type AllocateSimMutationHookResult = ReturnType<typeof useAllocateSimMutation>;
export type AllocateSimMutationResult = Apollo.MutationResult<AllocateSimMutation>;
export type AllocateSimMutationOptions = Apollo.BaseMutationOptions<AllocateSimMutation, AllocateSimMutationVariables>;
export const ToggleSimStatusDocument = gql`
    mutation toggleSimStatus($data: ToggleSimStatusInputDto!) {
  toggleSimStatus(data: $data) {
    simId
  }
}
    `;
export type ToggleSimStatusMutationFn = Apollo.MutationFunction<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;

/**
 * __useToggleSimStatusMutation__
 *
 * To run a mutation, you first call `useToggleSimStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleSimStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleSimStatusMutation, { data, loading, error }] = useToggleSimStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleSimStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>(ToggleSimStatusDocument, options);
      }
export type ToggleSimStatusMutationHookResult = ReturnType<typeof useToggleSimStatusMutation>;
export type ToggleSimStatusMutationResult = Apollo.MutationResult<ToggleSimStatusMutation>;
export type ToggleSimStatusMutationOptions = Apollo.BaseMutationOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;
export const GetSimDocument = gql`
    query getSim($data: GetSimInputDto!) {
  getSim(data: $data) {
    ...Sim
  }
}
    ${SimFragmentDoc}`;

/**
 * __useGetSimQuery__
 *
 * To run a query within a React component, call `useGetSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimQuery(baseOptions: Apollo.QueryHookOptions<GetSimQuery, GetSimQueryVariables> & ({ variables: GetSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
      }
export function useGetSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
// @ts-ignore
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimQuery, GetSimQueryVariables>;
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimQuery | undefined, GetSimQueryVariables>;
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
export type GetSimQueryHookResult = ReturnType<typeof useGetSimQuery>;
export type GetSimLazyQueryHookResult = ReturnType<typeof useGetSimLazyQuery>;
export type GetSimSuspenseQueryHookResult = ReturnType<typeof useGetSimSuspenseQuery>;
export type GetSimQueryResult = Apollo.QueryResult<GetSimQuery, GetSimQueryVariables>;
export const GetSimsDocument = gql`
    query GetSims($data: ListSimsInput!) {
  getSims(data: $data) {
    sims {
      ...Sim
    }
  }
}
    ${SimFragmentDoc}`;

/**
 * __useGetSimsQuery__
 *
 * To run a query within a React component, call `useGetSimsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsQuery(baseOptions: Apollo.QueryHookOptions<GetSimsQuery, GetSimsQueryVariables> & ({ variables: GetSimsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
      }
export function useGetSimsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
// @ts-ignore
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsQuery, GetSimsQueryVariables>;
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsQuery | undefined, GetSimsQueryVariables>;
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
export type GetSimsQueryHookResult = ReturnType<typeof useGetSimsQuery>;
export type GetSimsLazyQueryHookResult = ReturnType<typeof useGetSimsLazyQuery>;
export type GetSimsSuspenseQueryHookResult = ReturnType<typeof useGetSimsSuspenseQuery>;
export type GetSimsQueryResult = Apollo.QueryResult<GetSimsQuery, GetSimsQueryVariables>;
export const GetDataUsagesDocument = gql`
    query GetDataUsages($data: SimUsagesInputDto!) {
  getDataUsages(data: $data) {
    usages {
      usage
      simId
    }
  }
}
    `;

/**
 * __useGetDataUsagesQuery__
 *
 * To run a query within a React component, call `useGetDataUsagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDataUsagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDataUsagesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetDataUsagesQuery(baseOptions: Apollo.QueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables> & ({ variables: GetDataUsagesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
      }
export function useGetDataUsagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
// @ts-ignore
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetDataUsagesQuery, GetDataUsagesQueryVariables>;
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetDataUsagesQuery | undefined, GetDataUsagesQueryVariables>;
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
export type GetDataUsagesQueryHookResult = ReturnType<typeof useGetDataUsagesQuery>;
export type GetDataUsagesLazyQueryHookResult = ReturnType<typeof useGetDataUsagesLazyQuery>;
export type GetDataUsagesSuspenseQueryHookResult = ReturnType<typeof useGetDataUsagesSuspenseQuery>;
export type GetDataUsagesQueryResult = Apollo.QueryResult<GetDataUsagesQuery, GetDataUsagesQueryVariables>;