/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Packages — selling & performance, wired to the analytics service
 *  (`getPackagePerformance`). KPIs come from the keyed `kpis` array; per-plan
 *  rows from `packages`. Some fields/keys are exposed ahead of backend support
 *  (see docs/analytics-backend-gaps.md) and degrade to "—". */
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Meter from '@/components/Meter';

import { useGetPackagePerformanceQuery } from '@/client/graphql/analytics.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { EmptyState } from '@/components/EmptyState';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { sectionValue } from '@/components/SectionFallback';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import StatusBadge from '@/components/StatusBadge';
import { useCurrency } from '@/lib/currency';
import { kpiAmount } from '@/lib/kpis';
import { useUiPrefs } from '@/lib/store';

const BAR_COLORS = [
  'var(--uk-ac)',
  'var(--uk-secondary)',
  'var(--uk-success-bright)',
  'var(--uk-beige)',
  'var(--uk-orange)',
];

const isActive = (status?: string | null): boolean =>
  (status ?? '').toLowerCase() === 'active';

export default function BizPackagesScreen() {
  const networkId = useUiPrefs((s) => s.networkId);
  // Org currency symbol comes from getCurrencySymbol (shared via CurrencyProvider).
  const { symbol } = useCurrency();
  const money = (value?: number | null): string =>
    value == null
      ? '—'
      : `${symbol}${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;
  const { data, loading, error, refetch } = useGetPackagePerformanceQuery({
    variables: { data: { networkId } },
  });
  const perf = data?.getPackagePerformance;
  const kpis = perf?.kpis;
  const plans = perf?.packages ?? [];
  const totalRevenue = plans.reduce((sum, p) => sum + (p.revenue ?? 0), 0);
  // revenueSharePct is a backend gap — derive from total until it lands.
  const sharePct = (p: (typeof plans)[number]): number =>
    p.revenueSharePct ??
    (totalRevenue > 0 ? Math.round((p.revenue / totalRevenue) * 100) : 0);

  const byRevenue = [...plans].sort((a, z) => z.revenue - a.revenue);
  const top = byRevenue[0];
  const topPkgs = byRevenue.slice(0, 3);
  const maxRevenue = Math.max(...topPkgs.map((p) => p.revenue), 1);
  const mix = (perf?.revenueMix ?? [])
    .filter((m) => m.value > 0)
    .map((m, i) => ({
      name: m.name ?? '—',
      value: m.value,
      color: BAR_COLORS[i % BAR_COLORS.length] ?? 'var(--uk-ac)',
    }));

  return (
    <div className="page">
      <PageHeader
        title="Packages"
        sub="How your data packages are selling and performing."
        actions={<DateChip />}
      />
      <KpiRow
        cols={4}
        items={[
          {
            icon: 'monetization_on',
            color: 'var(--uk-beige)',
            label: 'Monthly recurring revenue',
            value: error ? '—' : kpiAmount(kpis, 'mrr', money),
            sub: 'paid this month',
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'ARPU',
            value: error ? '—' : kpiAmount(kpis, 'arpu', money),
            sub: 'avg revenue / active SIM',
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Top plan by revenue',
            value: top && top.revenue > 0 ? (top.name ?? '—') : '—',
            sub:
              top && top.revenue > 0
                ? `${sharePct(top)}% of revenue`
                : undefined,
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Active plans',
            value: error
              ? '—'
              : sectionValue(
                  plans.filter((p) => isActive(p.status)).length || null,
                ),
            sub: `across ${plans.length} packages`,
          },
        ]}
      />

      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div className="sec-head">
          <div className="sec-title">Package performance</div>
        </div>
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={6} rows={4} />
          ) : error ? (
            <EmptyState
              art="error"
              title="Couldn't load packages"
              sub={error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : plans.length === 0 ? (
            <EmptyState
              art="invoice"
              title="No packages"
              sub="Create a data plan to get started."
            />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Package</TableCell>
                  <TableCell align="right">Price</TableCell>
                  <TableCell align="right">Active SIMs</TableCell>
                  <TableCell align="right">Revenue</TableCell>
                  <TableCell align="right">Share</TableCell>
                  <TableCell>Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {byRevenue.map((p) => (
                  <TableRow key={p.packageId}>
                    <TableCell style={{ fontWeight: 600 }}>
                      {p.name ?? '—'}
                    </TableCell>
                    <TableCell align="right" className="tnum">
                      {money(p.price)}
                    </TableCell>
                    <TableCell align="right" className="tnum">
                      {p.activeSubscribers ?? '—'}
                    </TableCell>
                    <TableCell
                      align="right"
                      className="tnum"
                      style={{ fontWeight: 600 }}
                    >
                      {money(p.revenue)}
                    </TableCell>
                    <TableCell align="right" className="tnum muted">
                      {sharePct(p)}%
                    </TableCell>
                    <TableCell>
                      <StatusBadge
                        status={isActive(p.status) ? 'active' : 'inactive'}
                        variant="pill"
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </div>

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        <SectionCard
          title="Top packages"
          right={
            <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>
              By revenue collected
            </span>
          }
        >
          {topPkgs.length === 0 || error ? (
            <div
              style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}
            >
              No package revenue yet.
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column' }}>
              {topPkgs.map((p, i) => (
                <div
                  key={p.packageId}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 14,
                    padding: '13px 0',
                    borderBottom:
                      i < topPkgs.length - 1
                        ? '1px solid var(--uk-line-soft)'
                        : 'none',
                  }}
                >
                  <span
                    className="tnum"
                    style={{
                      fontFamily: 'var(--font-display)',
                      fontSize: 15,
                      fontWeight: 500,
                      color: 'var(--uk-ink-3)',
                      width: 18,
                      flex: 'none',
                    }}
                  >
                    {i + 1}
                  </span>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        gap: 12,
                        marginBottom: 6,
                      }}
                    >
                      <span style={{ fontSize: 13.5, fontWeight: 600 }}>
                        {p.name ?? '—'}
                      </span>
                      <span
                        className="tnum"
                        style={{
                          fontSize: 13,
                          color: 'var(--uk-ink-2)',
                          whiteSpace: 'nowrap',
                        }}
                      >
                        <b style={{ color: 'var(--uk-ink)' }}>
                          {money(p.revenue)}
                        </b>
                        {p.activeSubscribers != null
                          ? ` · ${p.activeSubscribers} active`
                          : ''}
                      </span>
                    </div>
                    <Meter
                      value={Math.round((p.revenue / maxRevenue) * 100)}
                      color={BAR_COLORS[i % BAR_COLORS.length]}
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </SectionCard>
        <SectionCard title="Package revenue mix">
          {mix.length === 0 || error ? (
            <div
              style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}
            >
              No package revenue yet.
            </div>
          ) : (
            <BarList rows={mix} />
          )}
        </SectionCard>
      </div>
    </div>
  );
}
