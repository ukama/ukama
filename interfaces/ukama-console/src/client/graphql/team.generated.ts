import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewNodeFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type TeamListQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type TeamListQuery = { __typename?: 'Query', membersView: { __typename?: 'MembersView', orgName: string, team: { __typename?: 'TeamSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, rows?: Array<{ __typename?: 'TeamMemberDto', id: string, name?: string | null, email?: string | null, role: string, status: string, memberSince?: string | null, inviteExpiresAt?: string | null }> | null } } };

export type InventoryOverviewQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type InventoryOverviewQuery = { __typename?: 'Query', inventoryView: { __typename?: 'InventoryView', orgName: string, components: { __typename?: 'ComponentStatsSection', total?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, byCategory?: Array<{ __typename?: 'CategoryCountDto', category: string, count: number }> | null }, unassignedNodes: { __typename?: 'NodesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, nodes?: Array<{ __typename?: 'Node', id: string, name: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } }> | null }, simStock: { __typename?: 'SimPoolStatsSection', total?: number | null, available?: number | null, consumed?: number | null, pctAssigned?: number | null, lowStock?: boolean | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };

export type SubscriberDetailQueryVariables = Types.Exact<{
  subscriberId: Types.Scalars['String']['input'];
}>;


export type SubscriberDetailQuery = { __typename?: 'Query', subscriberView: { __typename?: 'SubscriberView', subscriberId: string, subscriber: { __typename?: 'SubscriberSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, subscriber?: { __typename?: 'SubscriberDto', uuid: string, name: string, email: string, phone: string, networkId: string, sim?: Array<{ __typename?: 'SubscriberSimDto', id: string, iccid: string, msisdn: string, status: string, type: string }> | null } | null }, plans: { __typename?: 'SubscriberPlansSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, plans?: Array<{ __typename?: 'PlanNameDto', packageId: string, name: string }> | null }, billing: { __typename?: 'SubscriberBillingSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, payments?: Array<{ __typename?: 'PaymentDto', id: string, amount: string, currency: string, status: string, paidAt: string, paymentMethod: string }> | null }, usage: { __typename?: 'GapSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };


export const TeamListDocument = gql`
    query TeamList {
  membersView {
    orgName
    team {
      error {
        ...SectionErrorFields
      }
      rows {
        id
        name
        email
        role
        status
        memberSince
        inviteExpiresAt
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useTeamListQuery__
 *
 * To run a query within a React component, call `useTeamListQuery` and pass it any options that fit your needs.
 * When your component renders, `useTeamListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useTeamListQuery({
 *   variables: {
 *   },
 * });
 */
export function useTeamListQuery(baseOptions?: Apollo.QueryHookOptions<TeamListQuery, TeamListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<TeamListQuery, TeamListQueryVariables>(TeamListDocument, options);
      }
export function useTeamListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<TeamListQuery, TeamListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<TeamListQuery, TeamListQueryVariables>(TeamListDocument, options);
        }
// @ts-ignore
export function useTeamListSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<TeamListQuery, TeamListQueryVariables>): Apollo.UseSuspenseQueryResult<TeamListQuery, TeamListQueryVariables>;
export function useTeamListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<TeamListQuery, TeamListQueryVariables>): Apollo.UseSuspenseQueryResult<TeamListQuery | undefined, TeamListQueryVariables>;
export function useTeamListSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<TeamListQuery, TeamListQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<TeamListQuery, TeamListQueryVariables>(TeamListDocument, options);
        }
export type TeamListQueryHookResult = ReturnType<typeof useTeamListQuery>;
export type TeamListLazyQueryHookResult = ReturnType<typeof useTeamListLazyQuery>;
export type TeamListSuspenseQueryHookResult = ReturnType<typeof useTeamListSuspenseQuery>;
export type TeamListQueryResult = Apollo.QueryResult<TeamListQuery, TeamListQueryVariables>;
export const InventoryOverviewDocument = gql`
    query InventoryOverview {
  inventoryView {
    orgName
    components {
      error {
        ...SectionErrorFields
      }
      total
      byCategory {
        category
        count
      }
    }
    unassignedNodes {
      error {
        ...SectionErrorFields
      }
      nodes {
        ...ViewNode
      }
    }
    simStock {
      error {
        ...SectionErrorFields
      }
      total
      available
      consumed
      pctAssigned
      lowStock
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewNodeFragmentDoc}`;

/**
 * __useInventoryOverviewQuery__
 *
 * To run a query within a React component, call `useInventoryOverviewQuery` and pass it any options that fit your needs.
 * When your component renders, `useInventoryOverviewQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useInventoryOverviewQuery({
 *   variables: {
 *   },
 * });
 */
export function useInventoryOverviewQuery(baseOptions?: Apollo.QueryHookOptions<InventoryOverviewQuery, InventoryOverviewQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<InventoryOverviewQuery, InventoryOverviewQueryVariables>(InventoryOverviewDocument, options);
      }
export function useInventoryOverviewLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<InventoryOverviewQuery, InventoryOverviewQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<InventoryOverviewQuery, InventoryOverviewQueryVariables>(InventoryOverviewDocument, options);
        }
// @ts-ignore
export function useInventoryOverviewSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<InventoryOverviewQuery, InventoryOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<InventoryOverviewQuery, InventoryOverviewQueryVariables>;
export function useInventoryOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<InventoryOverviewQuery, InventoryOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<InventoryOverviewQuery | undefined, InventoryOverviewQueryVariables>;
export function useInventoryOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<InventoryOverviewQuery, InventoryOverviewQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<InventoryOverviewQuery, InventoryOverviewQueryVariables>(InventoryOverviewDocument, options);
        }
export type InventoryOverviewQueryHookResult = ReturnType<typeof useInventoryOverviewQuery>;
export type InventoryOverviewLazyQueryHookResult = ReturnType<typeof useInventoryOverviewLazyQuery>;
export type InventoryOverviewSuspenseQueryHookResult = ReturnType<typeof useInventoryOverviewSuspenseQuery>;
export type InventoryOverviewQueryResult = Apollo.QueryResult<InventoryOverviewQuery, InventoryOverviewQueryVariables>;
export const SubscriberDetailDocument = gql`
    query SubscriberDetail($subscriberId: String!) {
  subscriberView(subscriberId: $subscriberId) {
    subscriberId
    subscriber {
      error {
        ...SectionErrorFields
      }
      subscriber {
        uuid
        name
        email
        phone
        networkId
        sim {
          id
          iccid
          msisdn
          status
          type
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
    billing {
      error {
        ...SectionErrorFields
      }
      payments {
        id
        amount
        currency
        status
        paidAt
        paymentMethod
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
 * __useSubscriberDetailQuery__
 *
 * To run a query within a React component, call `useSubscriberDetailQuery` and pass it any options that fit your needs.
 * When your component renders, `useSubscriberDetailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSubscriberDetailQuery({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *   },
 * });
 */
export function useSubscriberDetailQuery(baseOptions: Apollo.QueryHookOptions<SubscriberDetailQuery, SubscriberDetailQueryVariables> & ({ variables: SubscriberDetailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SubscriberDetailQuery, SubscriberDetailQueryVariables>(SubscriberDetailDocument, options);
      }
export function useSubscriberDetailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SubscriberDetailQuery, SubscriberDetailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SubscriberDetailQuery, SubscriberDetailQueryVariables>(SubscriberDetailDocument, options);
        }
// @ts-ignore
export function useSubscriberDetailSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SubscriberDetailQuery, SubscriberDetailQueryVariables>): Apollo.UseSuspenseQueryResult<SubscriberDetailQuery, SubscriberDetailQueryVariables>;
export function useSubscriberDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SubscriberDetailQuery, SubscriberDetailQueryVariables>): Apollo.UseSuspenseQueryResult<SubscriberDetailQuery | undefined, SubscriberDetailQueryVariables>;
export function useSubscriberDetailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SubscriberDetailQuery, SubscriberDetailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SubscriberDetailQuery, SubscriberDetailQueryVariables>(SubscriberDetailDocument, options);
        }
export type SubscriberDetailQueryHookResult = ReturnType<typeof useSubscriberDetailQuery>;
export type SubscriberDetailLazyQueryHookResult = ReturnType<typeof useSubscriberDetailLazyQuery>;
export type SubscriberDetailSuspenseQueryHookResult = ReturnType<typeof useSubscriberDetailSuspenseQuery>;
export type SubscriberDetailQueryResult = Apollo.QueryResult<SubscriberDetailQuery, SubscriberDetailQueryVariables>;