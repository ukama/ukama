/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
'use client';

/** Revenue — wired to the analytics gateway's generic KPI API. The headline
 *  tiles are all REVENUE at different ops (SUM = revenue, COUNT = purchases,
 *  AVG = avg purchase) plus PAID_CUSTOMERS, each at the selected rolling span;
 *  the weekly trend comes from `getKpiTimeSeries`; the by-package breakdown
 *  from the `package_performance` report. Absent keys degrade to "—". */
import Skeleton from '@mui/material/Skeleton';
import { useState } from 'react';

import {
  useGetKpiTimeSeriesQuery,
  useGetKpiValuesQuery,
  useGetPerformanceReportQuery,
} from '@/client/graphql/analytics.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import RevenueTrendChart from '@/components/RevenueTrendChart';
import SectionCard from '@/components/SectionCard';
import { BAR_COLORS } from '@/lib/charts';
import { useCurrency } from '@/lib/currency';
import { DEFAULT_RANGE, type KpiSpan, rangeToSpan } from '@/lib/dateRange';
import {
  attrValue,
  cellMoney,
  KPI_KEYS,
  kpiAmount,
  kpiChangeAbs,
  kpiDelta,
  kpiText,
} from '@/lib/kpis';
import { POLL_LIVE_MS, visiblePoll } from '@/lib/polling';
import { useNetworkId } from '@/lib/useNetworkId';

// Trend/sub phrasing keyed by the selected rolling span, matching the DateChip.
const TREND_SUFFIX: Record<KpiSpan, string> = {
  last_24h: 'vs prev 24h',
  last_7d: 'vs prev 7 days',
  last_30d: 'vs prev 30 days',
};
const PERIOD_LABEL: Record<KpiSpan, string> = {
  last_24h: 'last 24h',
  last_7d: 'last 7 days',
  last_30d: 'last 30 days',
};

const intFmt = (v: number) => String(Math.round(v));

export default function BizSalesScreen() {
  const networkId = useNetworkId();
  // Org currency symbol from getCurrencySymbol (shared via CurrencyProvider).
  const { money } = useCurrency();
  const [range, setRange] = useState<string>(DEFAULT_RANGE);
  const span = rangeToSpan(range);
  // Fixed 30-day window for the revenue trend (daily buckets). Computed once
  // so the query variables stay stable across renders.
  const [trendRange] = useState(() => ({
    to: new Date().toISOString(),
    from: new Date(Date.now() - 30 * 864e5).toISOString(),
  }));

  // getKpiValues resolves ONE op per call, so REVENUE at SUM/COUNT/AVG needs
  // three calls. The base call uses each KPI's default op: REVENUE -> SUM,
  // PAID_CUSTOMERS -> LAST.
  const [fetchedAt, setFetchedAt] = useState(() => new Date());
  const { data: baseData, loading: baseLoading } = useGetKpiValuesQuery({
    variables: {
      data: { keys: [KPI_KEYS.revenue, KPI_KEYS.paidCustomers], span, networkId },
    },
    skip: !networkId,
    onCompleted: () => setFetchedAt(new Date()),
    ...visiblePoll(POLL_LIVE_MS, true),
  });
  const { data: countData, loading: countLoading } = useGetKpiValuesQuery({
    variables: {
      data: { keys: [KPI_KEYS.revenue], span, op: 'COUNT', networkId },
    },
    skip: !networkId,
    ...visiblePoll(POLL_LIVE_MS, true),
  });
  const { data: avgData, loading: avgLoading } = useGetKpiValuesQuery({
    variables: {
      data: { keys: [KPI_KEYS.revenue], span, op: 'AVG', networkId },
    },
    skip: !networkId,
    ...visiblePoll(POLL_LIVE_MS, true),
  });

  const { data: reportData, error: reportError } = useGetPerformanceReportQuery({
    variables: {
      data: { report: 'package_performance', span: 'last_30d', networkId },
    },
    skip: !networkId,
    ...visiblePoll(POLL_LIVE_MS, true),
  });

  // Revenue trend: daily REVENUE (SUM) over the last 30 days.
  const { data: trendData, loading: trendLoading } = useGetKpiTimeSeriesQuery({
    variables: {
      data: {
        keys: [KPI_KEYS.revenue],
        span: 'daily',
        op: 'SUM',
        from: trendRange.from,
        to: trendRange.to,
        networkId,
      },
    },
    skip: !networkId,
    ...visiblePoll(POLL_LIVE_MS, true),
  });

  const base = baseData?.getKpiValues.values;
  const countVals = countData?.getKpiValues.values;
  const avgVals = avgData?.getKpiValues.values;
  const kpiLoading = baseLoading || countLoading || avgLoading;
  const error = reportError;

  // Per-tile trends: revenue/avg are % (changePct); purchases/paid are counts
  // (changeAbs).
  const revDelta = kpiDelta(base, KPI_KEYS.revenue);
  const purchasesDelta = kpiChangeAbs(countVals, KPI_KEYS.revenue);
  const avgDelta = kpiDelta(avgVals, KPI_KEYS.revenue);
  const paidDelta = kpiChangeAbs(base, KPI_KEYS.paidCustomers);

  // Daily revenue points for the trend chart (REVENUE is reported in cents).
  const trendPoints = (trendData?.getKpiTimeSeries.values ?? [])
    .filter((v) => v.from)
    .map((v) => ({
      t: new Date(v.from as string).getTime(),
      value: (v.unit ?? '').toLowerCase() === 'cents' ? v.value / 100 : v.value,
    }))
    .sort((a, z) => a.t - z.t);

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
        fetchedAt={fetchedAt}
      />

      {kpiLoading ? (
        <Skeleton variant="rounded" sx={{ height: 96 }} />
      ) : (
        <KpiRow
          items={[
            {
              icon: 'monetization_on',
              color: 'var(--uk-beige)',
              label: 'Revenue',
              value: kpiAmount(base, KPI_KEYS.revenue, money),
              delta:
                revDelta != null
                  ? `${revDelta >= 0 ? '+' : ''}${Math.round(revDelta)}% ${TREND_SUFFIX[span]}`
                  : undefined,
              dir: revDelta != null && revDelta < 0 ? 'down' : 'up',
            },
            {
              icon: 'sell',
              color: 'var(--uk-secondary)',
              label: 'Purchases',
              value: kpiText(countVals, KPI_KEYS.revenue, intFmt),
              delta:
                purchasesDelta != null
                  ? `${purchasesDelta >= 0 ? '+' : ''}${Math.round(purchasesDelta)} ${PERIOD_LABEL[span]}`
                  : undefined,
              dir: purchasesDelta != null && purchasesDelta < 0 ? 'down' : 'up',
            },
            {
              icon: 'payments',
              color: 'var(--uk-ac)',
              label: 'Avg purchase',
              value: kpiAmount(avgVals, KPI_KEYS.revenue, money),
              // A ~0 change reads as "stable"; otherwise show the % move.
              sub:
                avgDelta == null || Math.round(avgDelta) === 0
                  ? 'stable'
                  : undefined,
              delta:
                avgDelta != null && Math.round(avgDelta) !== 0
                  ? `${avgDelta >= 0 ? '+' : ''}${Math.round(avgDelta)}% ${TREND_SUFFIX[span]}`
                  : undefined,
              dir: avgDelta != null && avgDelta < 0 ? 'down' : 'up',
            },
            {
              icon: 'group',
              color: 'var(--uk-success-bright)',
              label: 'Paid customers',
              value: kpiText(base, KPI_KEYS.paidCustomers, intFmt),
              delta:
                paidDelta != null
                  ? `${paidDelta >= 0 ? '+' : ''}${Math.round(paidDelta)} ${PERIOD_LABEL[span]}`
                  : undefined,
              dir: paidDelta != null && paidDelta < 0 ? 'down' : 'up',
            },
          ]}
        />
      )}

      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div
          className="sec-head"
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            marginBottom: 8,
          }}
        >
          <div className="sec-title">Revenue trend</div>
          <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>
            Last 30 days
          </span>
        </div>
        {trendLoading ? (
          <Skeleton variant="rounded" sx={{ height: 300 }} />
        ) : trendPoints.length === 0 ? (
          <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
            No revenue in this period.
          </div>
        ) : (
          <RevenueTrendChart
            points={trendPoints}
            formatValue={money}
            height={300}
          />
        )}
      </div>

      <SectionCard title="Revenue by package · Last 30 days">
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
