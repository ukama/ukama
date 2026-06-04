import * as Types from './types';

import { gql } from '@apollo/client';
import { ViewSiteFragmentDoc, SectionErrorFieldsFragmentDoc } from './views-shared.generated';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type BizHomeRevenueQueryVariables = Types.Exact<{
  networkId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type BizHomeRevenueQuery = { __typename?: 'Query', commerceView: { __typename?: 'CommerceView', networkId?: string | null, revenue: { __typename?: 'RevenueSection', totalPaid?: number | null, monthPaid?: number | null, momPct?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };

export type RevenueOverviewQueryVariables = Types.Exact<{
  networkId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type RevenueOverviewQuery = { __typename?: 'Query', commerceView: { __typename?: 'CommerceView', networkId?: string | null, revenue: { __typename?: 'RevenueSection', totalPaid?: number | null, totalPending?: number | null, monthPaid?: number | null, prevMonthPaid?: number | null, momPct?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, plans: { __typename?: 'PlanStatsSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, plans?: Array<{ __typename?: 'PlanStatsDto', packageId: string, name: string, revenue: number, revenueSharePct: number }> | null }, invoices: { __typename?: 'InvoicesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, reports?: Array<{ __typename?: 'ReportDto', id: string, period: string, type: string, isPaid: boolean, networkId: string }> | null } } };

export type BizHomeNetworkQueryVariables = Types.Exact<{
  networkId: Types.Scalars['String']['input'];
}>;


export type BizHomeNetworkQuery = { __typename?: 'Query', networkOverview: { __typename?: 'NetworkOverview', networkId: string, subscriberStats: { __typename?: 'SubscriberStatsSection', total?: number | null, active?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null }, siteStats: { __typename?: 'SitesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, sites?: Array<{ __typename?: 'SiteDto', id: string, name: string, networkId: string, latitude: string, longitude: string, location: string, isDeactivated: boolean, installDate: string, createdAt: string }> | null } } };

export type PackagesDashboardQueryVariables = Types.Exact<{
  networkId?: Types.InputMaybe<Types.Scalars['String']['input']>;
}>;


export type PackagesDashboardQuery = { __typename?: 'Query', commerceView: { __typename?: 'CommerceView', networkId?: string | null, plans: { __typename?: 'PlanStatsSection', mrr?: number | null, arpu?: number | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, plans?: Array<{ __typename?: 'PlanStatsDto', packageId: string, name: string, amount: number, currency: string, active: boolean, attachCount?: number | null, revenue: number, revenueSharePct: number }> | null } } };

export type BillingOverviewQueryVariables = Types.Exact<{
  networkId?: Types.InputMaybe<Types.Scalars['String']['input']>;
  invoiceLimit?: Types.Scalars['Int']['input'];
}>;


export type BillingOverviewQuery = { __typename?: 'Query', commerceView: { __typename?: 'CommerceView', networkId?: string | null, invoices: { __typename?: 'InvoicesSection', error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null, reports?: Array<{ __typename?: 'ReportDto', id: string, period: string, type: string, isPaid: boolean, networkId: string, createdAt: string }> | null }, balance: { __typename?: 'BalanceSection', outstandingCount?: number | null, latestUnpaidPeriod?: string | null, error?: { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string } | null } } };


export const BizHomeRevenueDocument = gql`
    query BizHomeRevenue($networkId: String) {
  commerceView(networkId: $networkId) {
    networkId
    revenue {
      error {
        ...SectionErrorFields
      }
      totalPaid
      monthPaid
      momPct
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useBizHomeRevenueQuery__
 *
 * To run a query within a React component, call `useBizHomeRevenueQuery` and pass it any options that fit your needs.
 * When your component renders, `useBizHomeRevenueQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useBizHomeRevenueQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useBizHomeRevenueQuery(baseOptions?: Apollo.QueryHookOptions<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>(BizHomeRevenueDocument, options);
      }
export function useBizHomeRevenueLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>(BizHomeRevenueDocument, options);
        }
// @ts-ignore
export function useBizHomeRevenueSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>): Apollo.UseSuspenseQueryResult<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>;
export function useBizHomeRevenueSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>): Apollo.UseSuspenseQueryResult<BizHomeRevenueQuery | undefined, BizHomeRevenueQueryVariables>;
export function useBizHomeRevenueSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>(BizHomeRevenueDocument, options);
        }
export type BizHomeRevenueQueryHookResult = ReturnType<typeof useBizHomeRevenueQuery>;
export type BizHomeRevenueLazyQueryHookResult = ReturnType<typeof useBizHomeRevenueLazyQuery>;
export type BizHomeRevenueSuspenseQueryHookResult = ReturnType<typeof useBizHomeRevenueSuspenseQuery>;
export type BizHomeRevenueQueryResult = Apollo.QueryResult<BizHomeRevenueQuery, BizHomeRevenueQueryVariables>;
export const RevenueOverviewDocument = gql`
    query RevenueOverview($networkId: String) {
  commerceView(networkId: $networkId) {
    networkId
    revenue {
      error {
        ...SectionErrorFields
      }
      totalPaid
      totalPending
      monthPaid
      prevMonthPaid
      momPct
    }
    plans {
      error {
        ...SectionErrorFields
      }
      plans {
        packageId
        name
        revenue
        revenueSharePct
      }
    }
    invoices(limit: 10) {
      error {
        ...SectionErrorFields
      }
      reports {
        id
        period
        type
        isPaid
        networkId
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useRevenueOverviewQuery__
 *
 * To run a query within a React component, call `useRevenueOverviewQuery` and pass it any options that fit your needs.
 * When your component renders, `useRevenueOverviewQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRevenueOverviewQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useRevenueOverviewQuery(baseOptions?: Apollo.QueryHookOptions<RevenueOverviewQuery, RevenueOverviewQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RevenueOverviewQuery, RevenueOverviewQueryVariables>(RevenueOverviewDocument, options);
      }
export function useRevenueOverviewLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RevenueOverviewQuery, RevenueOverviewQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RevenueOverviewQuery, RevenueOverviewQueryVariables>(RevenueOverviewDocument, options);
        }
// @ts-ignore
export function useRevenueOverviewSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<RevenueOverviewQuery, RevenueOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<RevenueOverviewQuery, RevenueOverviewQueryVariables>;
export function useRevenueOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RevenueOverviewQuery, RevenueOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<RevenueOverviewQuery | undefined, RevenueOverviewQueryVariables>;
export function useRevenueOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RevenueOverviewQuery, RevenueOverviewQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<RevenueOverviewQuery, RevenueOverviewQueryVariables>(RevenueOverviewDocument, options);
        }
export type RevenueOverviewQueryHookResult = ReturnType<typeof useRevenueOverviewQuery>;
export type RevenueOverviewLazyQueryHookResult = ReturnType<typeof useRevenueOverviewLazyQuery>;
export type RevenueOverviewSuspenseQueryHookResult = ReturnType<typeof useRevenueOverviewSuspenseQuery>;
export type RevenueOverviewQueryResult = Apollo.QueryResult<RevenueOverviewQuery, RevenueOverviewQueryVariables>;
export const BizHomeNetworkDocument = gql`
    query BizHomeNetwork($networkId: String!) {
  networkOverview(networkId: $networkId) {
    networkId
    subscriberStats {
      error {
        ...SectionErrorFields
      }
      total
      active
    }
    siteStats {
      error {
        ...SectionErrorFields
      }
      sites {
        ...ViewSite
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}
${ViewSiteFragmentDoc}`;

/**
 * __useBizHomeNetworkQuery__
 *
 * To run a query within a React component, call `useBizHomeNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useBizHomeNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useBizHomeNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useBizHomeNetworkQuery(baseOptions: Apollo.QueryHookOptions<BizHomeNetworkQuery, BizHomeNetworkQueryVariables> & ({ variables: BizHomeNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>(BizHomeNetworkDocument, options);
      }
export function useBizHomeNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>(BizHomeNetworkDocument, options);
        }
// @ts-ignore
export function useBizHomeNetworkSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>;
export function useBizHomeNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<BizHomeNetworkQuery | undefined, BizHomeNetworkQueryVariables>;
export function useBizHomeNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>(BizHomeNetworkDocument, options);
        }
export type BizHomeNetworkQueryHookResult = ReturnType<typeof useBizHomeNetworkQuery>;
export type BizHomeNetworkLazyQueryHookResult = ReturnType<typeof useBizHomeNetworkLazyQuery>;
export type BizHomeNetworkSuspenseQueryHookResult = ReturnType<typeof useBizHomeNetworkSuspenseQuery>;
export type BizHomeNetworkQueryResult = Apollo.QueryResult<BizHomeNetworkQuery, BizHomeNetworkQueryVariables>;
export const PackagesDashboardDocument = gql`
    query PackagesDashboard($networkId: String) {
  commerceView(networkId: $networkId) {
    networkId
    plans {
      error {
        ...SectionErrorFields
      }
      mrr
      arpu
      plans {
        packageId
        name
        amount
        currency
        active
        attachCount
        revenue
        revenueSharePct
      }
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __usePackagesDashboardQuery__
 *
 * To run a query within a React component, call `usePackagesDashboardQuery` and pass it any options that fit your needs.
 * When your component renders, `usePackagesDashboardQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePackagesDashboardQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function usePackagesDashboardQuery(baseOptions?: Apollo.QueryHookOptions<PackagesDashboardQuery, PackagesDashboardQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PackagesDashboardQuery, PackagesDashboardQueryVariables>(PackagesDashboardDocument, options);
      }
export function usePackagesDashboardLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PackagesDashboardQuery, PackagesDashboardQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PackagesDashboardQuery, PackagesDashboardQueryVariables>(PackagesDashboardDocument, options);
        }
// @ts-ignore
export function usePackagesDashboardSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<PackagesDashboardQuery, PackagesDashboardQueryVariables>): Apollo.UseSuspenseQueryResult<PackagesDashboardQuery, PackagesDashboardQueryVariables>;
export function usePackagesDashboardSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<PackagesDashboardQuery, PackagesDashboardQueryVariables>): Apollo.UseSuspenseQueryResult<PackagesDashboardQuery | undefined, PackagesDashboardQueryVariables>;
export function usePackagesDashboardSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<PackagesDashboardQuery, PackagesDashboardQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<PackagesDashboardQuery, PackagesDashboardQueryVariables>(PackagesDashboardDocument, options);
        }
export type PackagesDashboardQueryHookResult = ReturnType<typeof usePackagesDashboardQuery>;
export type PackagesDashboardLazyQueryHookResult = ReturnType<typeof usePackagesDashboardLazyQuery>;
export type PackagesDashboardSuspenseQueryHookResult = ReturnType<typeof usePackagesDashboardSuspenseQuery>;
export type PackagesDashboardQueryResult = Apollo.QueryResult<PackagesDashboardQuery, PackagesDashboardQueryVariables>;
export const BillingOverviewDocument = gql`
    query BillingOverview($networkId: String, $invoiceLimit: Int! = 20) {
  commerceView(networkId: $networkId) {
    networkId
    invoices(limit: $invoiceLimit) {
      error {
        ...SectionErrorFields
      }
      reports {
        id
        period
        type
        isPaid
        networkId
        createdAt
      }
    }
    balance {
      error {
        ...SectionErrorFields
      }
      outstandingCount
      latestUnpaidPeriod
    }
  }
}
    ${SectionErrorFieldsFragmentDoc}`;

/**
 * __useBillingOverviewQuery__
 *
 * To run a query within a React component, call `useBillingOverviewQuery` and pass it any options that fit your needs.
 * When your component renders, `useBillingOverviewQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useBillingOverviewQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *      invoiceLimit: // value for 'invoiceLimit'
 *   },
 * });
 */
export function useBillingOverviewQuery(baseOptions?: Apollo.QueryHookOptions<BillingOverviewQuery, BillingOverviewQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<BillingOverviewQuery, BillingOverviewQueryVariables>(BillingOverviewDocument, options);
      }
export function useBillingOverviewLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<BillingOverviewQuery, BillingOverviewQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<BillingOverviewQuery, BillingOverviewQueryVariables>(BillingOverviewDocument, options);
        }
// @ts-ignore
export function useBillingOverviewSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<BillingOverviewQuery, BillingOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<BillingOverviewQuery, BillingOverviewQueryVariables>;
export function useBillingOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BillingOverviewQuery, BillingOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<BillingOverviewQuery | undefined, BillingOverviewQueryVariables>;
export function useBillingOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<BillingOverviewQuery, BillingOverviewQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<BillingOverviewQuery, BillingOverviewQueryVariables>(BillingOverviewDocument, options);
        }
export type BillingOverviewQueryHookResult = ReturnType<typeof useBillingOverviewQuery>;
export type BillingOverviewLazyQueryHookResult = ReturnType<typeof useBillingOverviewLazyQuery>;
export type BillingOverviewSuspenseQueryHookResult = ReturnType<typeof useBillingOverviewSuspenseQuery>;
export type BillingOverviewQueryResult = Apollo.QueryResult<BillingOverviewQuery, BillingOverviewQueryVariables>;