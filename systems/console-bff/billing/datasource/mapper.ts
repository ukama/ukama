/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  CustomerDto,
  CustomerResDto,
  FeeDto,
  FeeResDto,
  GetReportDto,
  GetReportResDto,
  GetReportsDto,
  GetReportsResDto,
  ItemResDto,
  RawReportDto,
  RawReportResDto,
  ReportDto,
  ReportResDto,
  SubscriptionDto,
  SubscriptionResDto,
} from "../resolvers/types";

const mapItem = (item: ItemResDto) => {
  return {
    type: item.type,
    code: item.code,
    name: item.name,
  };
};

const mapCustomer = (customer: CustomerResDto): CustomerDto => {
  return {
    externalId: customer.external_id,
    name: customer.name,
    email: customer.email,
    addressLine1: customer.address_line1,
    legalName: customer.legal_name,
    legalNumber: customer.legal_number,
    phone: customer.phone,
    currency: customer.currency,
    timezone: customer.timezone,
    vatRate: customer.vat_rate,
    createdAt: customer.created_at,
  };
};

const mapFee = (fee: FeeResDto): FeeDto => {
  return {
    taxesAmountCents: fee.taxes_amount_cents,
    taxesPreciseAmount: fee.taxes_precise_amount,
    totalAmountCents: fee.total_amount_cents,
    totalAmountCurrency: fee.total_amount_currency,
    eventsCount: fee.events_count,
    units: fee.units,
    item: mapItem(fee.item),
  };
};

const mapSubscription = (subscription: SubscriptionResDto): SubscriptionDto => {
  return {
    externalCustomerId: subscription.external_customer_id,
    externalId: subscription.external_id,
    planCode: subscription.plan_code,
    name: subscription.name,
    status: subscription.status,
    createdAt: subscription.created_at,
    startedAt: subscription.started_at,
    canceledAt: subscription.canceled_at,
    terminatedAt: subscription.terminated_at,
  };
};

const mapRawReport = (rawReport: RawReportResDto): RawReportDto => {
  return {
    issuingDate: rawReport.issuing_date,
    paymentDueDate: rawReport.payment_due_date,
    paymentOverdue: rawReport.payment_overdue,
    invoiceType: rawReport.invoice_type,
    status: rawReport.status,
    paymentStatus: rawReport.payment_status,
    feesAmountCents: rawReport.fees_amount_cents,
    taxesAmountCents: rawReport.taxes_amount_cents,
    subTotalExcludingTaxesAmountCents:
      rawReport.sub_total_excluding_taxes_amount_cents,
    subTotalIncludingTaxesAmountCents:
      rawReport.sub_total_including_taxes_amount_cents,
    vatAmountCents: rawReport.vat_amount_cents,
    vatAmountCurrency: rawReport.vat_amount_currency,
    totalAmountCents: rawReport.total_amount_cents,
    currency: rawReport.currency,
    fileUrl: rawReport.file_url,
    customer: mapCustomer(rawReport.customer),
    subscriptions: rawReport.subscriptions.map(mapSubscription),
    fees: rawReport.fees.map(mapFee),
  };
};

const mapReport = (report: ReportResDto): ReportDto => {
  return {
    id: report.id,
    ownerId: report.owner_id,
    ownerType: report.owner_type,
    networkId: report.network_Id,
    period: report.period,
    type: report.Type,
    rawReport: mapRawReport(report.raw_report),
    isPaid: report.is_paid,
    createdAt: report.created_at,
  };
};

export const dtoToReportsDto = (res: GetReportsResDto): GetReportsDto => {
  return {
    reports: res.reports.map(mapReport),
  };
};

export const dtoToReportDto = (res: GetReportResDto): GetReportDto => {
  return {
    report: mapReport(res.report),
  };
};
