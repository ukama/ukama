import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type GetCurrencySymbolQueryVariables = Types.Exact<{
  code: Types.Scalars['String']['input'];
}>;


export type GetCurrencySymbolQuery = { __typename?: 'Query', getCurrencySymbol: { __typename?: 'CurrencyRes', code: string, symbol: string, image: string } };


export const GetCurrencySymbolDocument = gql`
    query getCurrencySymbol($code: String!) {
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