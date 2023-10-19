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
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  values: Array<Array<Scalars['Float']['output']>>;
};

export type Query = {
  __typename?: 'Query';
  _service: _Service;
  getLatestMetric: LatestMetricRes;
  getMetricRange: MetricRes;
  getNodeRangeMetric: MetricRes;
};


export type QueryGetLatestMetricArgs = {
  data: GetLatestMetricInput;
};


export type QueryGetMetricRangeArgs = {
  data: GetMetricRangeInput;
};


export type QueryGetNodeRangeMetricArgs = {
  data: GetMetricRangeInput;
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetricRangeSub: LatestMetricRes;
};


export type SubscriptionGetMetricRangeSubArgs = {
  from: Scalars['Float']['input'];
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type _Service = {
  __typename?: '_Service';
  sdl?: Maybe<Scalars['String']['output']>;
};

export type GetLatestMetricQueryVariables = Exact<{
  data: GetLatestMetricInput;
}>;


export type GetLatestMetricQuery = { __typename?: 'Query', getLatestMetric: { __typename?: 'LatestMetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, value: Array<number> } };

export type GetMetricRangeQueryVariables = Exact<{
  data: GetMetricRangeInput;
}>;


export type GetMetricRangeQuery = { __typename?: 'Query', getMetricRange: { __typename?: 'MetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, values: Array<Array<number>> } };

export type GetNodeRangeMetricQueryVariables = Exact<{
  data: GetMetricRangeInput;
}>;


export type GetNodeRangeMetricQuery = { __typename?: 'Query', getNodeRangeMetric: { __typename?: 'MetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, values: Array<Array<number>> } };

export type MetricRangeSubscriptionVariables = Exact<{
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  from: Scalars['Float']['input'];
}>;


export type MetricRangeSubscription = { __typename?: 'Subscription', getMetricRangeSub: { __typename?: 'LatestMetricRes', success: boolean, msg: string, orgId: string, nodeId: string, type: string, value: Array<number> } };


export const GetLatestMetricDocument = gql`
    query GetLatestMetric($data: GetLatestMetricInput!) {
  getLatestMetric(data: $data) {
    success
    msg
    orgId
    nodeId
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
    success
    msg
    orgId
    nodeId
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
export const GetNodeRangeMetricDocument = gql`
    query GetNodeRangeMetric($data: GetMetricRangeInput!) {
  getNodeRangeMetric(data: $data) {
    success
    msg
    orgId
    nodeId
    type
    values
  }
}
    `;

/**
 * __useGetNodeRangeMetricQuery__
 *
 * To run a query within a React component, call `useGetNodeRangeMetricQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeRangeMetricQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeRangeMetricQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeRangeMetricQuery(baseOptions: Apollo.QueryHookOptions<GetNodeRangeMetricQuery, GetNodeRangeMetricQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeRangeMetricQuery, GetNodeRangeMetricQueryVariables>(GetNodeRangeMetricDocument, options);
      }
export function useGetNodeRangeMetricLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeRangeMetricQuery, GetNodeRangeMetricQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeRangeMetricQuery, GetNodeRangeMetricQueryVariables>(GetNodeRangeMetricDocument, options);
        }
export type GetNodeRangeMetricQueryHookResult = ReturnType<typeof useGetNodeRangeMetricQuery>;
export type GetNodeRangeMetricLazyQueryHookResult = ReturnType<typeof useGetNodeRangeMetricLazyQuery>;
export type GetNodeRangeMetricQueryResult = Apollo.QueryResult<GetNodeRangeMetricQuery, GetNodeRangeMetricQueryVariables>;
export const MetricRangeDocument = gql`
    subscription MetricRange($nodeId: String!, $orgId: String!, $type: String!, $userId: String!, $from: Float!) {
  getMetricRangeSub(
    nodeId: $nodeId
    orgId: $orgId
    type: $type
    userId: $userId
    from: $from
  ) {
    success
    msg
    orgId
    nodeId
    type
    value
  }
}
    `;

/**
 * __useMetricRangeSubscription__
 *
 * To run a query within a React component, call `useMetricRangeSubscription` and pass it any options that fit your needs.
 * When your component renders, `useMetricRangeSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetricRangeSubscription({
 *   variables: {
 *      nodeId: // value for 'nodeId'
 *      orgId: // value for 'orgId'
 *      type: // value for 'type'
 *      userId: // value for 'userId'
 *      from: // value for 'from'
 *   },
 * });
 */
export function useMetricRangeSubscription(baseOptions: Apollo.SubscriptionHookOptions<MetricRangeSubscription, MetricRangeSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<MetricRangeSubscription, MetricRangeSubscriptionVariables>(MetricRangeDocument, options);
      }
export type MetricRangeSubscriptionHookResult = ReturnType<typeof useMetricRangeSubscription>;
export type MetricRangeSubscriptionResult = Apollo.SubscriptionResult<MetricRangeSubscription>;