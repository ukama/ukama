import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type KpiFieldsFragment = { __typename?: 'KpiDto', key: string, value: number, formatted?: string | null, delta?: number | null, deltaPeriod?: string | null, stale?: boolean | null, asOf?: string | null };

export type NamedValueFieldsFragment = { __typename?: 'NamedValueDto', id?: string | null, name?: string | null, value: number };

export type GetBusinessHomeQueryVariables = Types.Exact<{
  data: Types.AnalyticsWindowInput;
}>;


export type GetBusinessHomeQuery = { __typename?: 'Query', getBusinessHome: { __typename?: 'BusinessHomeDto', kpis: Array<{ __typename?: 'KpiDto', key: string, value: number, formatted?: string | null, delta?: number | null, deltaPeriod?: string | null, stale?: boolean | null, asOf?: string | null }>, sites: Array<{ __typename?: 'SiteSummaryDto', siteId: string, name?: string | null, status?: string | null, revenue: number, customers: number }>, topPackages: Array<{ __typename?: 'NamedValueDto', id?: string | null, name?: string | null, value: number }>, recentActivity: Array<{ __typename?: 'ActivityItemDto', routingKey?: string | null, description?: string | null, occurredAt?: string | null }> } };

export type GetBusinessSitesQueryVariables = Types.Exact<{
  data: Types.AnalyticsWindowInput;
}>;


export type GetBusinessSitesQuery = { __typename?: 'Query', getBusinessSites: { __typename?: 'BusinessSitesDto', sites: Array<{ __typename?: 'BusinessSiteRowDto', siteId: string, name?: string | null, status?: string | null, revenue: number, revenueToday: number, customers: number, dataUsed: number, uptime: number, topPackage?: string | null, issue?: string | null, latitude: number, longitude: number }>, meta?: { __typename?: 'MetaDto', count: number } | null } };

export type GetSalesOverviewQueryVariables = Types.Exact<{
  data: Types.AnalyticsWindowInput;
}>;


export type GetSalesOverviewQuery = { __typename?: 'Query', getSalesOverview: { __typename?: 'SalesOverviewDto', kpis: Array<{ __typename?: 'KpiDto', key: string, value: number, formatted?: string | null, delta?: number | null, deltaPeriod?: string | null, stale?: boolean | null, asOf?: string | null }>, revenueTrend?: { __typename?: 'TimeSeriesDto', key: string, points: Array<{ __typename?: 'PointDto', time?: string | null, value: number }> } | null, revenueBySite: Array<{ __typename?: 'NamedValueDto', id?: string | null, name?: string | null, value: number }>, revenueByPackage: Array<{ __typename?: 'NamedValueDto', id?: string | null, name?: string | null, value: number }> } };

export type GetPackagePerformanceQueryVariables = Types.Exact<{
  data: Types.AnalyticsWindowInput;
}>;


export type GetPackagePerformanceQuery = { __typename?: 'Query', getPackagePerformance: { __typename?: 'PackagePerformanceDto', kpis: Array<{ __typename?: 'KpiDto', key: string, value: number, formatted?: string | null, delta?: number | null, deltaPeriod?: string | null, stale?: boolean | null, asOf?: string | null }>, packages: Array<{ __typename?: 'PackageRowDto', packageId: string, name?: string | null, price: number, validity?: string | null, dataQuota?: string | null, status?: string | null, soldCount: number, revenue: number, revenueSharePct?: number | null, dataUsed: number, activeSubscribers: number }>, revenueMix: Array<{ __typename?: 'NamedValueDto', id?: string | null, name?: string | null, value: number }>, meta?: { __typename?: 'MetaDto', count: number } | null } };

export const KpiFieldsFragmentDoc = gql`
    fragment KpiFields on KpiDto {
  key
  value
  formatted
  delta
  deltaPeriod
  stale
  asOf
}
    `;
export const NamedValueFieldsFragmentDoc = gql`
    fragment NamedValueFields on NamedValueDto {
  id
  name
  value
}
    `;
export const GetBusinessHomeDocument = gql`
    query GetBusinessHome($data: AnalyticsWindowInput!) {
  getBusinessHome(data: $data) {
    kpis {
      ...KpiFields
    }
    sites {
      siteId
      name
      status
      revenue
      customers
    }
    topPackages {
      ...NamedValueFields
    }
    recentActivity {
      routingKey
      description
      occurredAt
    }
  }
}
    ${KpiFieldsFragmentDoc}
${NamedValueFieldsFragmentDoc}`;

/**
 * __useGetBusinessHomeQuery__
 *
 * To run a query within a React component, call `useGetBusinessHomeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetBusinessHomeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetBusinessHomeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetBusinessHomeQuery(baseOptions: Apollo.QueryHookOptions<GetBusinessHomeQuery, GetBusinessHomeQueryVariables> & ({ variables: GetBusinessHomeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>(GetBusinessHomeDocument, options);
      }
export function useGetBusinessHomeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>(GetBusinessHomeDocument, options);
        }
// @ts-ignore
export function useGetBusinessHomeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>): Apollo.UseSuspenseQueryResult<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>;
export function useGetBusinessHomeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>): Apollo.UseSuspenseQueryResult<GetBusinessHomeQuery | undefined, GetBusinessHomeQueryVariables>;
export function useGetBusinessHomeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>(GetBusinessHomeDocument, options);
        }
export type GetBusinessHomeQueryHookResult = ReturnType<typeof useGetBusinessHomeQuery>;
export type GetBusinessHomeLazyQueryHookResult = ReturnType<typeof useGetBusinessHomeLazyQuery>;
export type GetBusinessHomeSuspenseQueryHookResult = ReturnType<typeof useGetBusinessHomeSuspenseQuery>;
export type GetBusinessHomeQueryResult = Apollo.QueryResult<GetBusinessHomeQuery, GetBusinessHomeQueryVariables>;
export const GetBusinessSitesDocument = gql`
    query GetBusinessSites($data: AnalyticsWindowInput!) {
  getBusinessSites(data: $data) {
    sites {
      siteId
      name
      status
      revenue
      revenueToday
      customers
      dataUsed
      uptime
      topPackage
      issue
      latitude
      longitude
    }
    meta {
      count
    }
  }
}
    `;

/**
 * __useGetBusinessSitesQuery__
 *
 * To run a query within a React component, call `useGetBusinessSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetBusinessSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetBusinessSitesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetBusinessSitesQuery(baseOptions: Apollo.QueryHookOptions<GetBusinessSitesQuery, GetBusinessSitesQueryVariables> & ({ variables: GetBusinessSitesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>(GetBusinessSitesDocument, options);
      }
export function useGetBusinessSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>(GetBusinessSitesDocument, options);
        }
// @ts-ignore
export function useGetBusinessSitesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>;
export function useGetBusinessSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetBusinessSitesQuery | undefined, GetBusinessSitesQueryVariables>;
export function useGetBusinessSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>(GetBusinessSitesDocument, options);
        }
export type GetBusinessSitesQueryHookResult = ReturnType<typeof useGetBusinessSitesQuery>;
export type GetBusinessSitesLazyQueryHookResult = ReturnType<typeof useGetBusinessSitesLazyQuery>;
export type GetBusinessSitesSuspenseQueryHookResult = ReturnType<typeof useGetBusinessSitesSuspenseQuery>;
export type GetBusinessSitesQueryResult = Apollo.QueryResult<GetBusinessSitesQuery, GetBusinessSitesQueryVariables>;
export const GetSalesOverviewDocument = gql`
    query GetSalesOverview($data: AnalyticsWindowInput!) {
  getSalesOverview(data: $data) {
    kpis {
      ...KpiFields
    }
    revenueTrend {
      key
      points {
        time
        value
      }
    }
    revenueBySite {
      ...NamedValueFields
    }
    revenueByPackage {
      ...NamedValueFields
    }
  }
}
    ${KpiFieldsFragmentDoc}
${NamedValueFieldsFragmentDoc}`;

/**
 * __useGetSalesOverviewQuery__
 *
 * To run a query within a React component, call `useGetSalesOverviewQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSalesOverviewQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSalesOverviewQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSalesOverviewQuery(baseOptions: Apollo.QueryHookOptions<GetSalesOverviewQuery, GetSalesOverviewQueryVariables> & ({ variables: GetSalesOverviewQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>(GetSalesOverviewDocument, options);
      }
export function useGetSalesOverviewLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>(GetSalesOverviewDocument, options);
        }
// @ts-ignore
export function useGetSalesOverviewSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>;
export function useGetSalesOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>): Apollo.UseSuspenseQueryResult<GetSalesOverviewQuery | undefined, GetSalesOverviewQueryVariables>;
export function useGetSalesOverviewSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>(GetSalesOverviewDocument, options);
        }
export type GetSalesOverviewQueryHookResult = ReturnType<typeof useGetSalesOverviewQuery>;
export type GetSalesOverviewLazyQueryHookResult = ReturnType<typeof useGetSalesOverviewLazyQuery>;
export type GetSalesOverviewSuspenseQueryHookResult = ReturnType<typeof useGetSalesOverviewSuspenseQuery>;
export type GetSalesOverviewQueryResult = Apollo.QueryResult<GetSalesOverviewQuery, GetSalesOverviewQueryVariables>;
export const GetPackagePerformanceDocument = gql`
    query GetPackagePerformance($data: AnalyticsWindowInput!) {
  getPackagePerformance(data: $data) {
    kpis {
      ...KpiFields
    }
    packages {
      packageId
      name
      price
      validity
      dataQuota
      status
      soldCount
      revenue
      revenueSharePct
      dataUsed
      activeSubscribers
    }
    revenueMix {
      ...NamedValueFields
    }
    meta {
      count
    }
  }
}
    ${KpiFieldsFragmentDoc}
${NamedValueFieldsFragmentDoc}`;

/**
 * __useGetPackagePerformanceQuery__
 *
 * To run a query within a React component, call `useGetPackagePerformanceQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagePerformanceQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagePerformanceQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPackagePerformanceQuery(baseOptions: Apollo.QueryHookOptions<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables> & ({ variables: GetPackagePerformanceQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>(GetPackagePerformanceDocument, options);
      }
export function useGetPackagePerformanceLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>(GetPackagePerformanceDocument, options);
        }
// @ts-ignore
export function useGetPackagePerformanceSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>;
export function useGetPackagePerformanceSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagePerformanceQuery | undefined, GetPackagePerformanceQueryVariables>;
export function useGetPackagePerformanceSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>(GetPackagePerformanceDocument, options);
        }
export type GetPackagePerformanceQueryHookResult = ReturnType<typeof useGetPackagePerformanceQuery>;
export type GetPackagePerformanceLazyQueryHookResult = ReturnType<typeof useGetPackagePerformanceLazyQuery>;
export type GetPackagePerformanceSuspenseQueryHookResult = ReturnType<typeof useGetPackagePerformanceSuspenseQuery>;
export type GetPackagePerformanceQueryResult = Apollo.QueryResult<GetPackagePerformanceQuery, GetPackagePerformanceQueryVariables>;