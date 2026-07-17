import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type KpiValueFieldsFragment = { __typename?: 'KpiValueDto', kpi: string, value: number, span?: string | null, unit?: string | null, symbol?: string | null, isPartial?: boolean | null, trend?: { __typename?: 'TrendDto', direction?: string | null, changePct?: number | null, changeAbs?: number | null, prevValue?: number | null, hasPrevious?: boolean | null } | null };

export type ReportCellFieldsFragment = { __typename?: 'ReportCellDto', column: string, value: number, unit?: string | null, symbol?: string | null, format?: string | null };

export type GetKpiValuesQueryVariables = Types.Exact<{
  data: Types.KpiValuesInput;
}>;


export type GetKpiValuesQuery = { __typename?: 'Query', getKpiValues: { __typename?: 'GetKpiValuesDto', values: Array<{ __typename?: 'KpiValueDto', kpi: string, value: number, span?: string | null, unit?: string | null, symbol?: string | null, isPartial?: boolean | null, trend?: { __typename?: 'TrendDto', direction?: string | null, changePct?: number | null, changeAbs?: number | null, prevValue?: number | null, hasPrevious?: boolean | null } | null }> } };

export type GetPerformanceReportQueryVariables = Types.Exact<{
  data: Types.PerformanceReportInput;
}>;


export type GetPerformanceReportQuery = { __typename?: 'Query', getPerformanceReport: { __typename?: 'GetPerformanceReportDto', report: string, title?: string | null, span?: string | null, rows: Array<{ __typename?: 'ReportRowDto', entityId: string, status?: string | null, attributes: Array<{ __typename?: 'ScopeEntryDto', key: string, value: string }>, cells: Array<{ __typename?: 'ReportCellDto', column: string, value: number, unit?: string | null, symbol?: string | null, format?: string | null }> }> } };

export const KpiValueFieldsFragmentDoc = gql`
    fragment KpiValueFields on KpiValueDto {
  kpi
  value
  span
  unit
  symbol
  isPartial
  trend {
    direction
    changePct
    changeAbs
    prevValue
    hasPrevious
  }
}
    `;
export const ReportCellFieldsFragmentDoc = gql`
    fragment ReportCellFields on ReportCellDto {
  column
  value
  unit
  symbol
  format
}
    `;
export const GetKpiValuesDocument = gql`
    query GetKpiValues($data: KpiValuesInput!) {
  getKpiValues(data: $data) {
    values {
      ...KpiValueFields
    }
  }
}
    ${KpiValueFieldsFragmentDoc}`;

/**
 * __useGetKpiValuesQuery__
 *
 * To run a query within a React component, call `useGetKpiValuesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetKpiValuesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetKpiValuesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetKpiValuesQuery(baseOptions: Apollo.QueryHookOptions<GetKpiValuesQuery, GetKpiValuesQueryVariables> & ({ variables: GetKpiValuesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetKpiValuesQuery, GetKpiValuesQueryVariables>(GetKpiValuesDocument, options);
      }
export function useGetKpiValuesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetKpiValuesQuery, GetKpiValuesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetKpiValuesQuery, GetKpiValuesQueryVariables>(GetKpiValuesDocument, options);
        }
// @ts-ignore
export function useGetKpiValuesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetKpiValuesQuery, GetKpiValuesQueryVariables>): Apollo.UseSuspenseQueryResult<GetKpiValuesQuery, GetKpiValuesQueryVariables>;
export function useGetKpiValuesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetKpiValuesQuery, GetKpiValuesQueryVariables>): Apollo.UseSuspenseQueryResult<GetKpiValuesQuery | undefined, GetKpiValuesQueryVariables>;
export function useGetKpiValuesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetKpiValuesQuery, GetKpiValuesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetKpiValuesQuery, GetKpiValuesQueryVariables>(GetKpiValuesDocument, options);
        }
export type GetKpiValuesQueryHookResult = ReturnType<typeof useGetKpiValuesQuery>;
export type GetKpiValuesLazyQueryHookResult = ReturnType<typeof useGetKpiValuesLazyQuery>;
export type GetKpiValuesSuspenseQueryHookResult = ReturnType<typeof useGetKpiValuesSuspenseQuery>;
export type GetKpiValuesQueryResult = Apollo.QueryResult<GetKpiValuesQuery, GetKpiValuesQueryVariables>;
export const GetPerformanceReportDocument = gql`
    query GetPerformanceReport($data: PerformanceReportInput!) {
  getPerformanceReport(data: $data) {
    report
    title
    span
    rows {
      entityId
      status
      attributes {
        key
        value
      }
      cells {
        ...ReportCellFields
      }
    }
  }
}
    ${ReportCellFieldsFragmentDoc}`;

/**
 * __useGetPerformanceReportQuery__
 *
 * To run a query within a React component, call `useGetPerformanceReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPerformanceReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPerformanceReportQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPerformanceReportQuery(baseOptions: Apollo.QueryHookOptions<GetPerformanceReportQuery, GetPerformanceReportQueryVariables> & ({ variables: GetPerformanceReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>(GetPerformanceReportDocument, options);
      }
export function useGetPerformanceReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>(GetPerformanceReportDocument, options);
        }
// @ts-ignore
export function useGetPerformanceReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>;
export function useGetPerformanceReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetPerformanceReportQuery | undefined, GetPerformanceReportQueryVariables>;
export function useGetPerformanceReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>(GetPerformanceReportDocument, options);
        }
export type GetPerformanceReportQueryHookResult = ReturnType<typeof useGetPerformanceReportQuery>;
export type GetPerformanceReportLazyQueryHookResult = ReturnType<typeof useGetPerformanceReportLazyQuery>;
export type GetPerformanceReportSuspenseQueryHookResult = ReturnType<typeof useGetPerformanceReportSuspenseQuery>;
export type GetPerformanceReportQueryResult = Apollo.QueryResult<GetPerformanceReportQuery, GetPerformanceReportQueryVariables>;export type GetKpiTimeSeriesQueryVariables = Types.Exact<{
  data: Types.KpiTimeSeriesInput;
}>;


export type GetKpiTimeSeriesQuery = { __typename?: 'Query', getKpiTimeSeries: { __typename?: 'GetKpiTimeSeriesDto', values: Array<{ __typename?: 'KpiValueDto', kpi: string, value: number, span?: string | null, op?: string | null, from?: string | null, to?: string | null, unit?: string | null, symbol?: string | null, isPartial?: boolean | null }> } };

export const GetKpiTimeSeriesDocument = gql`
    query GetKpiTimeSeries($data: KpiTimeSeriesInput!) {
  getKpiTimeSeries(data: $data) {
    values {
      kpi
      value
      span
      op
      from
      to
      unit
      symbol
      isPartial
    }
  }
}
    `;

/**
 * __useGetKpiTimeSeriesQuery__
 *
 * To run a query within a React component, call `useGetKpiTimeSeriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetKpiTimeSeriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetKpiTimeSeriesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetKpiTimeSeriesQuery(baseOptions: Apollo.QueryHookOptions<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables> & ({ variables: GetKpiTimeSeriesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>(GetKpiTimeSeriesDocument, options);
      }
export function useGetKpiTimeSeriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>(GetKpiTimeSeriesDocument, options);
        }
// @ts-ignore
export function useGetKpiTimeSeriesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>;
export function useGetKpiTimeSeriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetKpiTimeSeriesQuery | undefined, GetKpiTimeSeriesQueryVariables>;
export function useGetKpiTimeSeriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>(GetKpiTimeSeriesDocument, options);
        }
export type GetKpiTimeSeriesQueryHookResult = ReturnType<typeof useGetKpiTimeSeriesQuery>;
export type GetKpiTimeSeriesLazyQueryHookResult = ReturnType<typeof useGetKpiTimeSeriesLazyQuery>;
export type GetKpiTimeSeriesSuspenseQueryHookResult = ReturnType<typeof useGetKpiTimeSeriesSuspenseQuery>;
export type GetKpiTimeSeriesQueryResult = Apollo.QueryResult<GetKpiTimeSeriesQuery, GetKpiTimeSeriesQueryVariables>;
