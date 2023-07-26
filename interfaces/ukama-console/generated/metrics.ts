import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string | number; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  _Any: { input: any; output: any; }
  _FieldSet: { input: any; output: any; }
};

export type GetLatestMetricInput = {
  nodeId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type GetMetricRangeInput = {
  from?: InputMaybe<Scalars['Float']['input']>;
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  step?: InputMaybe<Scalars['Float']['input']>;
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  withSubscription?: InputMaybe<Scalars['Boolean']['input']>;
};

export type LatestMetricRes = {
  __typename?: 'LatestMetricRes';
  env: Scalars['String']['output'];
  nodeid: Scalars['String']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  env: Scalars['String']['output'];
  nodeid: Scalars['String']['output'];
  type: Scalars['String']['output'];
  values: Array<Array<Scalars['Float']['output']>>;
};

export type Query = {
  __typename?: 'Query';
  _service: _Service;
  getLatestMetric: LatestMetricRes;
  getMetricRange: MetricRes;
};


export type QueryGetLatestMetricArgs = {
  data: GetLatestMetricInput;
};


export type QueryGetMetricRangeArgs = {
  data: GetMetricRangeInput;
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetricRangeSub: MetricRes;
};


export type SubscriptionGetMetricRangeSubArgs = {
  from?: InputMaybe<Scalars['Float']['input']>;
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  step?: InputMaybe<Scalars['Float']['input']>;
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  withSubscription?: InputMaybe<Scalars['Boolean']['input']>;
};

export type _Service = {
  __typename?: '_Service';
  sdl?: Maybe<Scalars['String']['output']>;
};

export type GetLatestMetricQueryVariables = Exact<{
  data: GetLatestMetricInput;
}>;


export type GetLatestMetricQuery = { __typename?: 'Query', getLatestMetric: { __typename?: 'LatestMetricRes', env: string, nodeid: string, type: string, value: Array<number> } };

export type GetMetricRangeQueryVariables = Exact<{
  data: GetMetricRangeInput;
}>;


export type GetMetricRangeQuery = { __typename?: 'Query', getMetricRange: { __typename?: 'MetricRes', env: string, nodeid: string, type: string, values: Array<Array<number>> } };

export type GetMetricRangeSubSubscriptionVariables = Exact<{
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
}>;


export type GetMetricRangeSubSubscription = { __typename?: 'Subscription', getMetricRangeSub: { __typename?: 'MetricRes', env: string, nodeid: string, type: string, values: Array<Array<number>> } };


export const GetLatestMetricDocument = gql`
    query GetLatestMetric($data: GetLatestMetricInput!) {
  getLatestMetric(data: $data) {
    env
    nodeid
    type
    value
  }
}
    `;

/**
 * __useGetLatestMetricQuery__
 *
 * To run a query within a React component, call `useGetLatestMetricQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetLatestMetricQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetLatestMetricQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetLatestMetricQuery(baseOptions: Apollo.QueryHookOptions<GetLatestMetricQuery, GetLatestMetricQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetLatestMetricQuery, GetLatestMetricQueryVariables>(GetLatestMetricDocument, options);
      }
export function useGetLatestMetricLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetLatestMetricQuery, GetLatestMetricQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetLatestMetricQuery, GetLatestMetricQueryVariables>(GetLatestMetricDocument, options);
        }
export type GetLatestMetricQueryHookResult = ReturnType<typeof useGetLatestMetricQuery>;
export type GetLatestMetricLazyQueryHookResult = ReturnType<typeof useGetLatestMetricLazyQuery>;
export type GetLatestMetricQueryResult = Apollo.QueryResult<GetLatestMetricQuery, GetLatestMetricQueryVariables>;
export const GetMetricRangeDocument = gql`
    query GetMetricRange($data: GetMetricRangeInput!) {
  getMetricRange(data: $data) {
    env
    nodeid
    type
    values
  }
}
    `;

/**
 * __useGetMetricRangeQuery__
 *
 * To run a query within a React component, call `useGetMetricRangeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricRangeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricRangeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricRangeQuery(baseOptions: Apollo.QueryHookOptions<GetMetricRangeQuery, GetMetricRangeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricRangeQuery, GetMetricRangeQueryVariables>(GetMetricRangeDocument, options);
      }
export function useGetMetricRangeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricRangeQuery, GetMetricRangeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricRangeQuery, GetMetricRangeQueryVariables>(GetMetricRangeDocument, options);
        }
export type GetMetricRangeQueryHookResult = ReturnType<typeof useGetMetricRangeQuery>;
export type GetMetricRangeLazyQueryHookResult = ReturnType<typeof useGetMetricRangeLazyQuery>;
export type GetMetricRangeQueryResult = Apollo.QueryResult<GetMetricRangeQuery, GetMetricRangeQueryVariables>;
export const GetMetricRangeSubDocument = gql`
    subscription GetMetricRangeSub($nodeId: String!, $orgId: String!, $type: String!, $userId: String!) {
  getMetricRangeSub(nodeId: $nodeId, orgId: $orgId, type: $type, userId: $userId) {
    env
    nodeid
    type
    values
  }
}
    `;

/**
 * __useGetMetricRangeSubSubscription__
 *
 * To run a query within a React component, call `useGetMetricRangeSubSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricRangeSubSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricRangeSubSubscription({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *      orgId: // value for 'orgId'
 *      type: // value for 'type'
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetMetricRangeSubSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetMetricRangeSubSubscription, GetMetricRangeSubSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetMetricRangeSubSubscription, GetMetricRangeSubSubscriptionVariables>(GetMetricRangeSubDocument, options);
      }
export type GetMetricRangeSubSubscriptionHookResult = ReturnType<typeof useGetMetricRangeSubSubscription>;
export type GetMetricRangeSubSubscriptionResult = Apollo.SubscriptionResult<GetMetricRangeSubSubscription>;