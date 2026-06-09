import * as Types from './types';

import { gql } from '@apollo/client';
import { SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NetworkCustomersQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
}>;


export type NetworkCustomersQuery = { __typename?: 'Query', subscribersView: { __typename?: 'SubscribersView', networkId: string, subscribers: { __typename?: 'SubscribersSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, subscribers?: Array<{ __typename?: 'SubscriberDto', uuid: string, name: string, phone: string, email: string, networkId: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, iccid: string, msisdn: string, status: string, type: string, package?: { __typename?: 'SimPackageDto', package_id: string, is_active: boolean } | null }> | null }> | null }, plans: { __typename?: 'PlansSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, plans?: Array<{ __typename?: 'PlanNameDto', packageId: string, name: string }> | null }, usage: { __typename?: 'GapSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };


export const NetworkCustomersDocument = gql`
    query NetworkCustomers($networkId: String!) {
  subscribersView(networkId: $networkId) {
    networkId
    subscribers {
      error {
        ...SectionErrorFields
      }
      subscribers {
        uuid
        name
        phone
        email
        networkId
        sim {
          id
          iccid
          msisdn
          status
          type
          package {
            package_id
            is_active
          }
        }
      }
    }
    plans {
      error {
        ...SectionErrorFields
      }
      plans {
        packageId
        name
      }
    }
    usage {
      error {
        ...SectionErrorFields
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useNetworkCustomersQuery__
 *
 * To run a query within a React component, call `useNetworkCustomersQuery` and pass it any options that fit your needs.
 * When your component renders, `useNetworkCustomersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useNetworkCustomersQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useNetworkCustomersQuery(baseOptions: Apollo.QueryHookOptions<NetworkCustomersQuery, NetworkCustomersQueryVariables> & ({ variables: NetworkCustomersQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<NetworkCustomersQuery, NetworkCustomersQueryVariables>(NetworkCustomersDocument, options);
      }
export function useNetworkCustomersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<NetworkCustomersQuery, NetworkCustomersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<NetworkCustomersQuery, NetworkCustomersQueryVariables>(NetworkCustomersDocument, options);
        }
// @ts-ignore
export function useNetworkCustomersSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<NetworkCustomersQuery, NetworkCustomersQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkCustomersQuery, NetworkCustomersQueryVariables>;
export function useNetworkCustomersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkCustomersQuery, NetworkCustomersQueryVariables>): Apollo.UseSuspenseQueryResult<NetworkCustomersQuery | undefined, NetworkCustomersQueryVariables>;
export function useNetworkCustomersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<NetworkCustomersQuery, NetworkCustomersQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<NetworkCustomersQuery, NetworkCustomersQueryVariables>(NetworkCustomersDocument, options);
        }
export type NetworkCustomersQueryHookResult = ReturnType<typeof useNetworkCustomersQuery>;
export type NetworkCustomersLazyQueryHookResult = ReturnType<typeof useNetworkCustomersLazyQuery>;
export type NetworkCustomersSuspenseQueryHookResult = ReturnType<typeof useNetworkCustomersSuspenseQuery>;
export type NetworkCustomersQueryResult = Apollo.QueryResult<NetworkCustomersQuery, NetworkCustomersQueryVariables>;