import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type GetCountriesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetCountriesQuery = { __typename?: 'Query', getCountries: { __typename?: 'CountriesRes', countries: Array<{ __typename?: 'CountryDto', name: string, code: string }> } };

export type GetCurrencySymbolQueryVariables = Types.Exact<{
  code: Types.Scalars['String']['input'];
}>;


export type GetCurrencySymbolQuery = { __typename?: 'Query', getCurrencySymbol: { __typename?: 'CurrencyRes', code: string, symbol: string, image: string } };

export type GetTimezonesQueryVariables = Types.Exact<{ [key: string]: never; }>;


export type GetTimezonesQuery = { __typename?: 'Query', getTimezones: { __typename?: 'TimezoneRes', timezones: Array<{ __typename?: 'TimezoneDto', value: string, abbr: string, offset: number, isdst: boolean, text: string, utc: Array<string> }> } };


export const GetCountriesDocument = gql`
    query GetCountries {
  getCountries {
    countries {
      name
      code
    }
  }
}
    `;

/**
 * __useGetCountriesQuery__
 *
 * To run a query within a React component, call `useGetCountriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCountriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCountriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetCountriesQuery(baseOptions?: Apollo.QueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
      }
export function useGetCountriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
// @ts-ignore
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetCountriesQuery, GetCountriesQueryVariables>;
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetCountriesQuery | undefined, GetCountriesQueryVariables>;
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
export type GetCountriesQueryHookResult = ReturnType<typeof useGetCountriesQuery>;
export type GetCountriesLazyQueryHookResult = ReturnType<typeof useGetCountriesLazyQuery>;
export type GetCountriesSuspenseQueryHookResult = ReturnType<typeof useGetCountriesSuspenseQuery>;
export type GetCountriesQueryResult = Apollo.QueryResult<GetCountriesQuery, GetCountriesQueryVariables>;
export const GetCurrencySymbolDocument = gql`
    query GetCurrencySymbol($code: String!) {
  getCurrencySymbol(code: $code) {
    code
    symbol
    image
  }
}
    `;

/**
 * __useGetCurrencySymbolQuery__
 *
 * To run a query within a React component, call `useGetCurrencySymbolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCurrencySymbolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCurrencySymbolQuery({
 *   variables: {
 *      code: // value for 'code'
 *   },
 * });
 */
export function useGetCurrencySymbolQuery(baseOptions: Apollo.QueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables> & ({ variables: GetCurrencySymbolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
      }
export function useGetCurrencySymbolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
// @ts-ignore
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>): Apollo.UseSuspenseQueryResult<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>;
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>): Apollo.UseSuspenseQueryResult<GetCurrencySymbolQuery | undefined, GetCurrencySymbolQueryVariables>;
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
export type GetCurrencySymbolQueryHookResult = ReturnType<typeof useGetCurrencySymbolQuery>;
export type GetCurrencySymbolLazyQueryHookResult = ReturnType<typeof useGetCurrencySymbolLazyQuery>;
export type GetCurrencySymbolSuspenseQueryHookResult = ReturnType<typeof useGetCurrencySymbolSuspenseQuery>;
export type GetCurrencySymbolQueryResult = Apollo.QueryResult<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>;
export const GetTimezonesDocument = gql`
    query GetTimezones {
  getTimezones {
    timezones {
      value
      abbr
      offset
      isdst
      text
      utc
    }
  }
}
    `;

/**
 * __useGetTimezonesQuery__
 *
 * To run a query within a React component, call `useGetTimezonesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTimezonesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTimezonesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetTimezonesQuery(baseOptions?: Apollo.QueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
      }
export function useGetTimezonesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
// @ts-ignore
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>): Apollo.UseSuspenseQueryResult<GetTimezonesQuery, GetTimezonesQueryVariables>;
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>): Apollo.UseSuspenseQueryResult<GetTimezonesQuery | undefined, GetTimezonesQueryVariables>;
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
export type GetTimezonesQueryHookResult = ReturnType<typeof useGetTimezonesQuery>;
export type GetTimezonesLazyQueryHookResult = ReturnType<typeof useGetTimezonesLazyQuery>;
export type GetTimezonesSuspenseQueryHookResult = ReturnType<typeof useGetTimezonesSuspenseQuery>;
export type GetTimezonesQueryResult = Apollo.QueryResult<GetTimezonesQuery, GetTimezonesQueryVariables>;