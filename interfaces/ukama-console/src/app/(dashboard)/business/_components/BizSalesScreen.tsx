/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
'use client';

/** Revenue — the single most important number, wired to the analytics gateway's
 *  generic KPI API. Headline figures come from `getKpiValues` (revenue @ monthly
 *  with a month-over-month trend, plus mrr); the by-package breakdown from the
 *  `package_performance` report (`getPerformanceReport`). The gateway has no
 *  cumulative "collected to date" or "pending" KPI, so those are not shown.
 *  Absent keys degrade to "—". */
import Skeleton from '@mui/material/Skeleton';
import { useState } from 'react';

import {
  useGetKpiValuesQuery,
  useGetPerformanceReportQuery,
} from '@/client/graphql/analytics.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { BAR_COLORS } from '@/lib/charts';
import { useCurrency } from '@/lib/currency';
import { DEFAULT_RANGE, rangeToSpan } from '@/lib/dateRange';
import {
  attrValue,
  cellMoney,
  KPI_KEYS,
  kpiAmount,
  kpiAmountPrev,
} from '@/lib/kpis';
import { useNetworkId } from '@/lib/useNetworkId';

export default function BizSalesScreen() {
  const networkId = useNetworkId();
  // Org currency symbol from getCurrencySymbol (shared via CurrencyProvider).
  const { money } = useCurrency();
  const [range, setRange] = useState<string>(DEFAULT_RANGE);
  const span = rangeToSpan(range);

  const { data: kpiData, loading: kpiLoading } = useGetKpiValuesQuery({
    variables: {
      data: {
        keys: [KPI_KEYS.revenue, KPI_KEYS.mrr],
        span,
        networkId,
      },
    },
    skip: !networkId,
  });

  const { data: reportData, loading: reportLoading, error: reportError } =
    useGetPerformanceReportQuery({
      variables: {
        data: { report: 'package_performance', span, networkId },
      },
      skip: !networkId,
    });

  const kpis = kpiData?.getKpiValues.values;
  // The by-package breakdown depends only on the report; KPI tiles degrade to
  // "—" on their own, so a KPI failure must not blank the breakdown.
  const error = reportError;
  const loading = kpiLoading || reportLoading;

  // Revenue-by-package bars from the package_performance report rows (money
  // columns are reported in cents).
  const byPackage = (reportData?.getPerformanceReport.rows ?? [])
    .map((r) => ({
      name: attrValue(r.attributes, 'name') ?? '—',
      value: cellMoney(r.cells, 'revenue'),
    }))
    .filter((r) => r.value > 0)
    .sort((a, z) => z.value - a.value)
    .map((r, i) => ({
      ...r,
      color: BAR_COLORS[i % BAR_COLORS.length] ?? 'var(--uk-ac)',
    }));

  return (
    <div className="page">
      <PageHeader
        title="Revenue"
        sub="Revenue across your network — your single most important number."
        actions={<DateChip value={range} onChange={setRange} />}
      />

      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div className="sec-title" style={{ marginBottom: 12 }}>
          Recurring revenue
        </div>
        {loading ? (
          <Skeleton variant="rounded" sx={{ height: 72 }} />
        ) : (
          <div
            style={{
              display: 'flex',
              alignItems: 'flex-end',
              gap: 20,
              flexWrap: 'wrap',
            }}
          >
            <div style={{ display: 'flex', alignItems: 'baseline', gap: 12 }}>
              <span
                className="tnum"
                style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 48,
                  fontWeight: 500,
                  lineHeight: 1,
                }}
              >
                {kpiAmount(kpis, KPI_KEYS.mrr, money)}
              </span>
            </div>
            <div style={{ flex: 1, minWidth: 20 }} />
            <div style={{ display: 'flex', gap: 30, flexWrap: 'wrap' }}>
              {(
                [
                  ['This period', kpiAmount(kpis, KPI_KEYS.revenue, money)],
                  ['Previous', kpiAmountPrev(kpis, KPI_KEYS.revenue, money)],
                ] as const
              ).map(([k, v]) => (
                <div key={k}>
                  <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                    {k}
                  </div>
                  <div
                    className="tnum"
                    style={{
                      fontFamily: 'var(--font-display)',
                      fontSize: 19,
                      fontWeight: 500,
                      marginTop: 2,
                    }}
                  >
                    {v}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
        <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)', marginTop: 8 }}>
          Recurring revenue and revenue for the selected period
        </div>
      </div>

      <KpiRow
        items={[
          {
            icon: 'monetization_on',
            color: 'var(--uk-beige)',
            label: 'Revenue',
            value: kpiAmount(kpis, KPI_KEYS.revenue, money),
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'Recurring revenue',
            value: kpiAmount(kpis, KPI_KEYS.mrr, money),
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Plans earning revenue',
            value: error ? '—' : String(byPackage.length),
          },
        ]}
      />

      <SectionCard title="Revenue by package">
        {error || byPackage.length === 0 ? (
          <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
            {error ? '—' : 'No package revenue yet.'}
          </div>
        ) : (
          <BarList rows={byPackage} />
        )}
      </SectionCard>
    </div>
  );
}
