import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type PackageRateFragment = { __typename?: 'PackageDto', rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number } };

export type PackageMarkupFragment = { __typename?: 'PackageDto', markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type SimPackagesFragment = { __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean };

export type SubscriberSimsFragment = { __typename?: 'SubscriberToSimsDto', subscriberId: string, sims: Array<{ __typename?: 'SubscriberSimsDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> };

export type PackageFragment = { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } };

export type GetPackagesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { __typename?: 'Query', getPackages: { __typename?: 'PackagesResDto', packages: Array<{ __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } }> } };

export type GetPackageQueryVariables = Types.Exact<{
  packageId: Types.Scalars['String']['input'];
}>;


export type GetPackageQuery = { __typename?: 'Query', getPackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type GetSimsBySubscriberQueryVariables = Types.Exact<{
  data: Types.GetSimBySubscriberInputDto;
}>;


export type GetSimsBySubscriberQuery = { __typename?: 'Query', getSimsBySubscriber: { __typename?: 'SubscriberToSimsDto', subscriberId: string, sims: Array<{ __typename?: 'SubscriberSimsDto', id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> } };

export type AddPackageMutationVariables = Types.Exact<{
  data: Types.AddPackageInputDto;
}>;


export type AddPackageMutation = { __typename?: 'Mutation', addPackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export type RemovePackageForSimMutationVariables = Types.Exact<{
  data: Types.RemovePackageFormSimInputDto;
}>;


export type RemovePackageForSimMutation = { __typename?: 'Mutation', removePackageForSim: { __typename?: 'RemovePackageFromSimResDto', packageId?: string | null } };

export type DeletePackageMutationVariables = Types.Exact<{
  packageId: Types.Scalars['String']['input'];
}>;


export type DeletePackageMutation = { __typename?: 'Mutation', deletePackage: { __typename?: 'IdResponse', uuid: string } };

export type GetPackagesForSimQueryVariables = Types.Exact<{
  data: Types.GetPackagesForSimInputDto;
}>;


export type GetPackagesForSimQuery = { __typename?: 'Query', getPackagesForSim: { __typename?: 'GetSimPackagesDtoAPI', sim_id: string, packages: Array<{ __typename?: 'SimToPackagesDto', id: string, package_id: string, start_date: string, end_date: string, is_active: boolean }> } };

export type AddPackagesToSimMutationVariables = Types.Exact<{
  data: Types.AddPackagesToSimInputDto;
}>;


export type AddPackagesToSimMutation = { __typename?: 'Mutation', addPackagesToSim: { __typename?: 'AddPackagesSimResDto', packages: Array<{ __typename?: 'AddPackagSimResDto', packageId?: string | null }> } };

export type DeleteSimMutationVariables = Types.Exact<{
  data: Types.DeleteSimInputDto;
}>;


export type DeleteSimMutation = { __typename?: 'Mutation', deleteSim: { __typename?: 'DeleteSimResDto', simId?: string | null } };

export type UpdatePacakgeMutationVariables = Types.Exact<{
  packageId: Types.Scalars['String']['input'];
  data: Types.UpdatePackageInputDto;
}>;


export type UpdatePacakgeMutation = { __typename?: 'Mutation', updatePackage: { __typename?: 'PackageDto', uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { __typename?: 'PackageRateAPIDto', sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { __typename?: 'PackageMarkupAPIDto', baserate: string, markup: number } } };

export const SimPackagesFragmentDoc = gql`
    fragment SimPackages on SimToPackagesDto {
  id
  package_id
  start_date
  end_date
  is_active
}
    `;
export const SubscriberSimsFragmentDoc = gql`
    fragment SubscriberSims on SubscriberToSimsDto {
  subscriberId
  sims {
    id
    subscriberId
    networkId
    iccid
    msisdn
    imsi
    type
    status
    allocatedAt
    isPhysical
  }
}
    `;
export const PackageRateFragmentDoc = gql`
    fragment PackageRate on PackageDto {
  rate {
    sms_mo
    sms_mt
    data
    amount
  }
}
    `;
export const PackageMarkupFragmentDoc = gql`
    fragment PackageMarkup on PackageDto {
  markup {
    baserate
    markup
  }
}
    `;
export const PackageFragmentDoc = gql`
    fragment Package on PackageDto {
  uuid
  name
  active
  duration
  simType
  createdAt
  deletedAt
  updatedAt
  smsVolume
  dataVolume
  voiceVolume
  ulbr
  dlbr
  type
  dataUnit
  voiceUnit
  messageUnit
  flatrate
  currency
  from
  to
  country
  provider
  apn
  ownerId
  amount
  ...PackageRate
  ...PackageMarkup
}
    ${PackageRateFragmentDoc}
${PackageMarkupFragmentDoc}`;
export const GetPackagesDocument = gql`
    query getPackages {
  getPackages {
    packages {
      ...Package
    }
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackagesQuery__
 *
 * To run a query within a React component, call `useGetPackagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetPackagesQuery(baseOptions?: Apollo.QueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
      }
export function useGetPackagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
// @ts-ignore
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesQuery, GetPackagesQueryVariables>;
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesQuery | undefined, GetPackagesQueryVariables>;
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
export type GetPackagesQueryHookResult = ReturnType<typeof useGetPackagesQuery>;
export type GetPackagesLazyQueryHookResult = ReturnType<typeof useGetPackagesLazyQuery>;
export type GetPackagesSuspenseQueryHookResult = ReturnType<typeof useGetPackagesSuspenseQuery>;
export type GetPackagesQueryResult = Apollo.QueryResult<GetPackagesQuery, GetPackagesQueryVariables>;
export const GetPackageDocument = gql`
    query getPackage($packageId: String!) {
  getPackage(packageId: $packageId) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackageQuery__
 *
 * To run a query within a React component, call `useGetPackageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackageQuery({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useGetPackageQuery(baseOptions: Apollo.QueryHookOptions<GetPackageQuery, GetPackageQueryVariables> & ({ variables: GetPackageQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
      }
export function useGetPackageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
// @ts-ignore
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackageQuery, GetPackageQueryVariables>;
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackageQuery | undefined, GetPackageQueryVariables>;
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
export type GetPackageQueryHookResult = ReturnType<typeof useGetPackageQuery>;
export type GetPackageLazyQueryHookResult = ReturnType<typeof useGetPackageLazyQuery>;
export type GetPackageSuspenseQueryHookResult = ReturnType<typeof useGetPackageSuspenseQuery>;
export type GetPackageQueryResult = Apollo.QueryResult<GetPackageQuery, GetPackageQueryVariables>;
export const GetSimsBySubscriberDocument = gql`
    query getSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    ...SubscriberSims
  }
}
    ${SubscriberSimsFragmentDoc}`;

/**
 * __useGetSimsBySubscriberQuery__
 *
 * To run a query within a React component, call `useGetSimsBySubscriberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsBySubscriberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsBySubscriberQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsBySubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables> & ({ variables: GetSimsBySubscriberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
      }
export function useGetSimsBySubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
// @ts-ignore
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>;
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsBySubscriberQuery | undefined, GetSimsBySubscriberQueryVariables>;
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
export type GetSimsBySubscriberQueryHookResult = ReturnType<typeof useGetSimsBySubscriberQuery>;
export type GetSimsBySubscriberLazyQueryHookResult = ReturnType<typeof useGetSimsBySubscriberLazyQuery>;
export type GetSimsBySubscriberSuspenseQueryHookResult = ReturnType<typeof useGetSimsBySubscriberSuspenseQuery>;
export type GetSimsBySubscriberQueryResult = Apollo.QueryResult<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>;
export const AddPackageDocument = gql`
    mutation addPackage($data: AddPackageInputDto!) {
  addPackage(data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type AddPackageMutationFn = Apollo.MutationFunction<AddPackageMutation, AddPackageMutationVariables>;

/**
 * __useAddPackageMutation__
 *
 * To run a mutation, you first call `useAddPackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackageMutation, { data, loading, error }] = useAddPackageMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackageMutation(baseOptions?: Apollo.MutationHookOptions<AddPackageMutation, AddPackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackageMutation, AddPackageMutationVariables>(AddPackageDocument, options);
      }
export type AddPackageMutationHookResult = ReturnType<typeof useAddPackageMutation>;
export type AddPackageMutationResult = Apollo.MutationResult<AddPackageMutation>;
export type AddPackageMutationOptions = Apollo.BaseMutationOptions<AddPackageMutation, AddPackageMutationVariables>;
export const RemovePackageForSimDocument = gql`
    mutation removePackageForSim($data: RemovePackageFormSimInputDto!) {
  removePackageForSim(data: $data) {
    packageId
  }
}
    `;
export type RemovePackageForSimMutationFn = Apollo.MutationFunction<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;

/**
 * __useRemovePackageForSimMutation__
 *
 * To run a mutation, you first call `useRemovePackageForSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemovePackageForSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removePackageForSimMutation, { data, loading, error }] = useRemovePackageForSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRemovePackageForSimMutation(baseOptions?: Apollo.MutationHookOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>(RemovePackageForSimDocument, options);
      }
export type RemovePackageForSimMutationHookResult = ReturnType<typeof useRemovePackageForSimMutation>;
export type RemovePackageForSimMutationResult = Apollo.MutationResult<RemovePackageForSimMutation>;
export type RemovePackageForSimMutationOptions = Apollo.BaseMutationOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;
export const DeletePackageDocument = gql`
    mutation deletePackage($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}
    `;
export type DeletePackageMutationFn = Apollo.MutationFunction<DeletePackageMutation, DeletePackageMutationVariables>;

/**
 * __useDeletePackageMutation__
 *
 * To run a mutation, you first call `useDeletePackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeletePackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deletePackageMutation, { data, loading, error }] = useDeletePackageMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useDeletePackageMutation(baseOptions?: Apollo.MutationHookOptions<DeletePackageMutation, DeletePackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeletePackageMutation, DeletePackageMutationVariables>(DeletePackageDocument, options);
      }
export type DeletePackageMutationHookResult = ReturnType<typeof useDeletePackageMutation>;
export type DeletePackageMutationResult = Apollo.MutationResult<DeletePackageMutation>;
export type DeletePackageMutationOptions = Apollo.BaseMutationOptions<DeletePackageMutation, DeletePackageMutationVariables>;
export const GetPackagesForSimDocument = gql`
    query getPackagesForSim($data: GetPackagesForSimInputDto!) {
  getPackagesForSim(data: $data) {
    sim_id
    packages {
      ...SimPackages
    }
  }
}
    ${SimPackagesFragmentDoc}`;

/**
 * __useGetPackagesForSimQuery__
 *
 * To run a query within a React component, call `useGetPackagesForSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesForSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesForSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPackagesForSimQuery(baseOptions: Apollo.QueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables> & ({ variables: GetPackagesForSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
      }
export function useGetPackagesForSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
// @ts-ignore
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesForSimQuery | undefined, GetPackagesForSimQueryVariables>;
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
export type GetPackagesForSimQueryHookResult = ReturnType<typeof useGetPackagesForSimQuery>;
export type GetPackagesForSimLazyQueryHookResult = ReturnType<typeof useGetPackagesForSimLazyQuery>;
export type GetPackagesForSimSuspenseQueryHookResult = ReturnType<typeof useGetPackagesForSimSuspenseQuery>;
export type GetPackagesForSimQueryResult = Apollo.QueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export const AddPackagesToSimDocument = gql`
    mutation addPackagesToSim($data: AddPackagesToSimInputDto!) {
  addPackagesToSim(data: $data) {
    packages {
      packageId
    }
  }
}
    `;
export type AddPackagesToSimMutationFn = Apollo.MutationFunction<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;

/**
 * __useAddPackagesToSimMutation__
 *
 * To run a mutation, you first call `useAddPackagesToSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackagesToSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackagesToSimMutation, { data, loading, error }] = useAddPackagesToSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackagesToSimMutation(baseOptions?: Apollo.MutationHookOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>(AddPackagesToSimDocument, options);
      }
export type AddPackagesToSimMutationHookResult = ReturnType<typeof useAddPackagesToSimMutation>;
export type AddPackagesToSimMutationResult = Apollo.MutationResult<AddPackagesToSimMutation>;
export type AddPackagesToSimMutationOptions = Apollo.BaseMutationOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;
export const DeleteSimDocument = gql`
    mutation deleteSim($data: DeleteSimInputDto!) {
  deleteSim(data: $data) {
    simId
  }
}
    `;
export type DeleteSimMutationFn = Apollo.MutationFunction<DeleteSimMutation, DeleteSimMutationVariables>;

/**
 * __useDeleteSimMutation__
 *
 * To run a mutation, you first call `useDeleteSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSimMutation, { data, loading, error }] = useDeleteSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteSimMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSimMutation, DeleteSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSimMutation, DeleteSimMutationVariables>(DeleteSimDocument, options);
      }
export type DeleteSimMutationHookResult = ReturnType<typeof useDeleteSimMutation>;
export type DeleteSimMutationResult = Apollo.MutationResult<DeleteSimMutation>;
export type DeleteSimMutationOptions = Apollo.BaseMutationOptions<DeleteSimMutation, DeleteSimMutationVariables>;
export const UpdatePacakgeDocument = gql`
    mutation updatePacakge($packageId: String!, $data: UpdatePackageInputDto!) {
  updatePackage(packageId: $packageId, data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type UpdatePacakgeMutationFn = Apollo.MutationFunction<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;

/**
 * __useUpdatePacakgeMutation__
 *
 * To run a mutation, you first call `useUpdatePacakgeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePacakgeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePacakgeMutation, { data, loading, error }] = useUpdatePacakgeMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePacakgeMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>(UpdatePacakgeDocument, options);
      }
export type UpdatePacakgeMutationHookResult = ReturnType<typeof useUpdatePacakgeMutation>;
export type UpdatePacakgeMutationResult = Apollo.MutationResult<UpdatePacakgeMutation>;
export type UpdatePacakgeMutationOptions = Apollo.BaseMutationOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;