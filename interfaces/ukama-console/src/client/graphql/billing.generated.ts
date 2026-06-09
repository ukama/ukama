import * as Types from './types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type PaymentFragment = { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string };

export type UpdatePaymentMutationVariables = Types.Exact<{
  data: Types.UpdatePaymentInputDto;
}>;


export type UpdatePaymentMutation = { __typename?: 'Mutation', updatePayment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } };

export type ProcessPaymentMutationVariables = Types.Exact<{
  data: Types.ProcessPaymentInputDto;
}>;


export type ProcessPaymentMutation = { __typename?: 'Mutation', processPayment: { __typename?: 'ProcessPaymentDto', payment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } } };

export type GetPaymentQueryVariables = Types.Exact<{
  paymentId: Types.Scalars['String']['input'];
}>;


export type GetPaymentQuery = { __typename?: 'Query', getPayment: { __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string } };

export type GetPaymentsQueryVariables = Types.Exact<{
  data: Types.GetPaymentsInputDto;
}>;


export type GetPaymentsQuery = { __typename?: 'Query', getPayments: { __typename?: 'PaymentsDto', payments: Array<{ __typename?: 'PaymentDto', id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string }> } };

export type CustomerFragment = { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string };

export type SubscriptionFragment = { __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null };

export type FeeFragment = { __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } };

export type RawReportFragment = { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> };

export type GetReportsQueryVariables = Types.Exact<{
  data: Types.GetReportsInputDto;
}>;


export type GetReportsQuery = { __typename?: 'Query', getReports: { __typename?: 'GetReportsDto', reports: Array<{ __typename?: 'ReportDto', id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> } }> } };

export type GetReportQueryVariables = Types.Exact<{
  id: Types.Scalars['String']['input'];
}>;


export type GetReportQuery = { __typename?: 'Query', getReport: { __typename?: 'GetReportDto', report: { __typename?: 'ReportDto', id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { __typename?: 'RawReportDto', issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency?: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { __typename?: 'CustomerDto', externalId: string, name: string, email?: string | null, addressLine1?: string | null, legalName?: string | null, legalNumber?: string | null, phone?: string | null, currency: string, timezone?: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ __typename?: 'SubscriptionDto', externalCustomerId: string, externalId: string, planCode: string, name?: string | null, status: string, createdAt: string, startedAt: string, canceledAt?: string | null, terminatedAt?: string | null }>, fees: Array<{ __typename?: 'FeeDto', taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { __typename?: 'ItemResDto', type: string, code: string, name: string } }> } } } };

export type GetGeneratedPdfReportQueryVariables = Types.Exact<{
  Id: Types.Scalars['String']['input'];
}>;


export type GetGeneratedPdfReportQuery = { __typename?: 'Query', getGeneratedPdfReport: { __typename?: 'GetPdfReportUrlDto', contentType: string, filename: string, downloadUrl: string } };

export const PaymentFragmentDoc = gql`
    fragment payment on PaymentDto {
  id
  itemId
  itemType
  amount
  currency
  paymentMethod
  depositedAmount
  paidAt
  payerName
  payerEmail
  payerPhone
  correspondent
  country
  description
  status
  failureReason
  extra
  createdAt
}
    `;
export const CustomerFragmentDoc = gql`
    fragment customer on CustomerDto {
  externalId
  name
  email
  addressLine1
  legalName
  legalNumber
  phone
  currency
  timezone
  vatRate
  createdAt
}
    `;
export const SubscriptionFragmentDoc = gql`
    fragment subscription on SubscriptionDto {
  externalCustomerId
  externalId
  planCode
  name
  status
  createdAt
  startedAt
  canceledAt
  terminatedAt
}
    `;
export const FeeFragmentDoc = gql`
    fragment fee on FeeDto {
  taxesAmountCents
  taxesPreciseAmount
  totalAmountCents
  totalAmountCurrency
  eventsCount
  units
  item {
    type
    code
    name
  }
}
    `;
export const RawReportFragmentDoc = gql`
    fragment rawReport on RawReportDto {
  issuingDate
  paymentDueDate
  paymentOverdue
  invoiceType
  status
  paymentStatus
  feesAmountCents
  taxesAmountCents
  subTotalExcludingTaxesAmountCents
  subTotalIncludingTaxesAmountCents
  vatAmountCents
  vatAmountCurrency
  totalAmountCents
  currency
  fileUrl
  customer {
    ...customer
  }
  subscriptions {
    ...subscription
  }
  fees {
    ...fee
  }
}
    ${CustomerFragmentDoc}
${SubscriptionFragmentDoc}
${FeeFragmentDoc}`;
export const UpdatePaymentDocument = gql`
    mutation UpdatePayment($data: UpdatePaymentInputDto!) {
  updatePayment(data: $data) {
    id
    itemId
    itemType
    amount
    currency
    paymentMethod
    depositedAmount
    paidAt
    payerName
    payerEmail
    payerPhone
    correspondent
    country
    description
    status
    failureReason
    createdAt
  }
}
    `;
export type UpdatePaymentMutationFn = Apollo.MutationFunction<UpdatePaymentMutation, UpdatePaymentMutationVariables>;

/**
 * __useUpdatePaymentMutation__
 *
 * To run a mutation, you first call `useUpdatePaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePaymentMutation, { data, loading, error }] = useUpdatePaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePaymentMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePaymentMutation, UpdatePaymentMutationVariables>(UpdatePaymentDocument, options);
      }
export type UpdatePaymentMutationHookResult = ReturnType<typeof useUpdatePaymentMutation>;
export type UpdatePaymentMutationResult = Apollo.MutationResult<UpdatePaymentMutation>;
export type UpdatePaymentMutationOptions = Apollo.BaseMutationOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>;
export const ProcessPaymentDocument = gql`
    mutation ProcessPayment($data: ProcessPaymentInputDto!) {
  processPayment(data: $data) {
    payment {
      id
      itemId
      itemType
      amount
      currency
      paymentMethod
      depositedAmount
      paidAt
      payerName
      payerEmail
      payerPhone
      correspondent
      country
      description
      status
      failureReason
      createdAt
    }
  }
}
    `;
export type ProcessPaymentMutationFn = Apollo.MutationFunction<ProcessPaymentMutation, ProcessPaymentMutationVariables>;

/**
 * __useProcessPaymentMutation__
 *
 * To run a mutation, you first call `useProcessPaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useProcessPaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [processPaymentMutation, { data, loading, error }] = useProcessPaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useProcessPaymentMutation(baseOptions?: Apollo.MutationHookOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ProcessPaymentMutation, ProcessPaymentMutationVariables>(ProcessPaymentDocument, options);
      }
export type ProcessPaymentMutationHookResult = ReturnType<typeof useProcessPaymentMutation>;
export type ProcessPaymentMutationResult = Apollo.MutationResult<ProcessPaymentMutation>;
export type ProcessPaymentMutationOptions = Apollo.BaseMutationOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>;
export const GetPaymentDocument = gql`
    query GetPayment($paymentId: String!) {
  getPayment(paymentId: $paymentId) {
    ...payment
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentQuery__
 *
 * To run a query within a React component, call `useGetPaymentQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentQuery({
 *   variables: {
 *      paymentId: // value for 'paymentId'
 *   },
 * });
 */
export function useGetPaymentQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables> & ({ variables: GetPaymentQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
      }
export function useGetPaymentLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
// @ts-ignore
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentQuery, GetPaymentQueryVariables>;
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentQuery | undefined, GetPaymentQueryVariables>;
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
export type GetPaymentQueryHookResult = ReturnType<typeof useGetPaymentQuery>;
export type GetPaymentLazyQueryHookResult = ReturnType<typeof useGetPaymentLazyQuery>;
export type GetPaymentSuspenseQueryHookResult = ReturnType<typeof useGetPaymentSuspenseQuery>;
export type GetPaymentQueryResult = Apollo.QueryResult<GetPaymentQuery, GetPaymentQueryVariables>;
export const GetPaymentsDocument = gql`
    query GetPayments($data: GetPaymentsInputDto!) {
  getPayments(data: $data) {
    payments {
      ...payment
    }
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentsQuery__
 *
 * To run a query within a React component, call `useGetPaymentsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPaymentsQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables> & ({ variables: GetPaymentsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
      }
export function useGetPaymentsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
// @ts-ignore
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentsQuery, GetPaymentsQueryVariables>;
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentsQuery | undefined, GetPaymentsQueryVariables>;
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
export type GetPaymentsQueryHookResult = ReturnType<typeof useGetPaymentsQuery>;
export type GetPaymentsLazyQueryHookResult = ReturnType<typeof useGetPaymentsLazyQuery>;
export type GetPaymentsSuspenseQueryHookResult = ReturnType<typeof useGetPaymentsSuspenseQuery>;
export type GetPaymentsQueryResult = Apollo.QueryResult<GetPaymentsQuery, GetPaymentsQueryVariables>;
export const GetReportsDocument = gql`
    query GetReports($data: GetReportsInputDto!) {
  getReports(data: $data) {
    reports {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportsQuery__
 *
 * To run a query within a React component, call `useGetReportsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetReportsQuery(baseOptions: Apollo.QueryHookOptions<GetReportsQuery, GetReportsQueryVariables> & ({ variables: GetReportsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
      }
export function useGetReportsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
// @ts-ignore
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportsQuery, GetReportsQueryVariables>;
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportsQuery | undefined, GetReportsQueryVariables>;
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
export type GetReportsQueryHookResult = ReturnType<typeof useGetReportsQuery>;
export type GetReportsLazyQueryHookResult = ReturnType<typeof useGetReportsLazyQuery>;
export type GetReportsSuspenseQueryHookResult = ReturnType<typeof useGetReportsSuspenseQuery>;
export type GetReportsQueryResult = Apollo.QueryResult<GetReportsQuery, GetReportsQueryVariables>;
export const GetReportDocument = gql`
    query GetReport($id: String!) {
  getReport(id: $id) {
    report {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportQuery__
 *
 * To run a query within a React component, call `useGetReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetReportQuery(baseOptions: Apollo.QueryHookOptions<GetReportQuery, GetReportQueryVariables> & ({ variables: GetReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
      }
export function useGetReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
// @ts-ignore
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportQuery, GetReportQueryVariables>;
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportQuery | undefined, GetReportQueryVariables>;
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
export type GetReportQueryHookResult = ReturnType<typeof useGetReportQuery>;
export type GetReportLazyQueryHookResult = ReturnType<typeof useGetReportLazyQuery>;
export type GetReportSuspenseQueryHookResult = ReturnType<typeof useGetReportSuspenseQuery>;
export type GetReportQueryResult = Apollo.QueryResult<GetReportQuery, GetReportQueryVariables>;
export const GetGeneratedPdfReportDocument = gql`
    query getGeneratedPdfReport($Id: String!) {
  getGeneratedPdfReport(id: $Id) {
    contentType
    filename
    downloadUrl
  }
}
    `;

/**
 * __useGetGeneratedPdfReportQuery__
 *
 * To run a query within a React component, call `useGetGeneratedPdfReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetGeneratedPdfReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetGeneratedPdfReportQuery({
 *   variables: {
 *      Id: // value for 'Id'
 *   },
 * });
 */
export function useGetGeneratedPdfReportQuery(baseOptions: Apollo.QueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables> & ({ variables: GetGeneratedPdfReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
      }
export function useGetGeneratedPdfReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
// @ts-ignore
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>;
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetGeneratedPdfReportQuery | undefined, GetGeneratedPdfReportQueryVariables>;
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
export type GetGeneratedPdfReportQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportQuery>;
export type GetGeneratedPdfReportLazyQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportLazyQuery>;
export type GetGeneratedPdfReportSuspenseQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportSuspenseQuery>;
export type GetGeneratedPdfReportQueryResult = Apollo.QueryResult<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>;