/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useMemo } from 'react';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import { BILLING_HISTORY_TABLE_MENU, BILLING_TABLE_COLUMNS } from '@/constants';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import { GetReportResDto } from '@/client/graphql/generated';

interface BillingHistoryProps {
  bills: GetReportResDto[];
  loading?: boolean;
  onViewDetails?: (reportId: string) => void;
}

const BillingHistory: React.FC<BillingHistoryProps> = ({
  bills,
  loading = false,
  onViewDetails,
}) => {
  const billingHistoryDataset = useMemo(() => {
    return bills.map((report: GetReportResDto) => ({
      date: new Date(report.createdAt).toLocaleDateString(),
      amount: `${report.rawReport.totalAmountCurrency} ${(report.rawReport.totalAmountCents / 100).toFixed(2)}`,
      status: report.rawReport.paymentStatus || report.rawReport.status,
      period: report.period,
      id: report.id,
    }));
  }, [bills]);

  const handleMenuItemClick = (id: string, type: string) => {
    if (type === 'view_details' && onViewDetails) {
      onViewDetails(id);
    }
  };

  return (
    <DataTableWithOptions
      columns={BILLING_TABLE_COLUMNS}
      icon={SubscriberIcon}
      dataset={billingHistoryDataset}
      menuOptions={BILLING_HISTORY_TABLE_MENU}
      onMenuItemClick={handleMenuItemClick}
      emptyViewLabel="No billing history found"
      isRowClickable={false}
      // loading={loading}
    />
  );
};

export default BillingHistory;
