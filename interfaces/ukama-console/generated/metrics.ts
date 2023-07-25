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

export type GetMetricInput = {
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type MetricRes = {
  __typename?: 'MetricRes';
  env: Scalars['String']['output'];
  nodeid: Scalars['String']['output'];
  type: Scalars['String']['output'];
  value: Array<Array<Scalars['Float']['output']>>;
};

export type Query = {
  __typename?: 'Query';
  _service: _Service;
  getMetrics: MetricRes;
};


export type QueryGetMetricsArgs = {
  data: GetMetricInput;
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetric: MetricRes;
};


export type SubscriptionGetMetricArgs = {
  nodeId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type _Service = {
  __typename?: '_Service';
  sdl?: Maybe<Scalars['String']['output']>;
};

export type GetMetricsQueryVariables = Exact<{
  data: GetMetricInput;
}>;


export type GetMetricsQuery = { __typename?: 'Query', getMetrics: { __typename?: 'MetricRes', env: string, type: string, nodeid: string, value: Array<Array<number>> } };

export type GetMetricSubscriptionVariables = Exact<{
  orgId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
  type: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
}>;


export type GetMetricSubscription = { __typename?: 'Subscription', getMetric: { __typename?: 'MetricRes', env: string, type: string, nodeid: string, value: Array<Array<number>> } };


export const GetMetricsDocument = gql`
    query GetMetrics($data: GetMetricInput!) {
  getMetrics(data: $data) {
    env
    type
    nodeid
    value
  }
}
    `;

/**
 * __useGetMetricsQuery__
 *
 * To run a query within a React component, call `useGetMetricsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetMetricsQuery(baseOptions: Apollo.QueryHookOptions<GetMetricsQuery, GetMetricsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMetricsQuery, GetMetricsQueryVariables>(GetMetricsDocument, options);
      }
export function useGetMetricsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMetricsQuery, GetMetricsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMetricsQuery, GetMetricsQueryVariables>(GetMetricsDocument, options);
        }
export type GetMetricsQueryHookResult = ReturnType<typeof useGetMetricsQuery>;
export type GetMetricsLazyQueryHookResult = ReturnType<typeof useGetMetricsLazyQuery>;
export type GetMetricsQueryResult = Apollo.QueryResult<GetMetricsQuery, GetMetricsQueryVariables>;
export const GetMetricDocument = gql`
    subscription getMetric($orgId: String!, $userId: String!, $type: String!, $nodeId: String!) {
  getMetric(orgId: $orgId, userId: $userId, type: $type, nodeId: $nodeId) {
    env
    type
    nodeid
    value
  }
}
    `;

/**
 * __useGetMetricSubscription__
 *
 * To run a query within a React component, call `useGetMetricSubscription` and pass it any options that fit your needs.
 * When your component renders, `useGetMetricSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMetricSubscription({
 *   variables: {
 *      orgId: // value for 'orgId'
 *      userId: // value for 'userId'
 *      type: // value for 'type'
 *      nodeId: // value for 'nodeId'
 *   },
 * });
 */
export function useGetMetricSubscription(baseOptions: Apollo.SubscriptionHookOptions<GetMetricSubscription, GetMetricSubscriptionVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useSubscription<GetMetricSubscription, GetMetricSubscriptionVariables>(GetMetricDocument, options);
      }
export type GetMetricSubscriptionHookResult = ReturnType<typeof useGetMetricSubscription>;
export type GetMetricSubscriptionResult = Apollo.SubscriptionResult<GetMetricSubscription>;