/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
'use client';

/** Packages — selling & performance, wired to the analytics gateway. Everything
 *  on this page derives from the `package_performance` report
 *  (`getPerformanceReport`): per-package attributes (name/price/validity/
 *  data_volume/active) and columns (sold/revenue/data_used), plus a
 *  threshold-derived status. The four headline tiles are report-derived totals:
 *  Package revenue (Σ revenue), Packages sold (Σ sold), Best package (best
 *  seller's allowance / validity) and Data consumed (DATA_USAGE). */
import { useState } from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Meter from '@/components/Meter';

import {
  useGetKpiValuesQuery,
  useGetPerformanceReportQuery,
} from '@/client/graphql/analytics.generated';
import DateChip from '@/components/DateChip';
import { EmptyState } from '@/components/EmptyState';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import StatusBadge from '@/components/StatusBadge';
import { BAR_COLORS } from '@/lib/charts';
import { useCurrency } from '@/lib/currency';
import { formatDuration } from '@/lib/duration';
import { DEFAULT_RANGE, rangeToSpan } from '@/lib/dateRange';
import { attrValue, cellMoney, cellValue, KPI_KEYS, kpiValue } from '@/lib/kpis';
import { dataVolumeToBytes, formatBytes } from '@/lib/usage';
import { heldQuery } from '@/lib/heldQuery';
import { useNetworkId, useNetworkQueryPending } from '@/lib/useNetworkId';

interface Plan {
  packageId: string;
  name: string;
  price: number;
  revenue: number;
  sold: number;
  dataVolume: number;
  dataUnit: string;
  duration: number;
  active: boolean;
  statusLabel?: string | null;
}

export default function BizPackagesScreen() {
  const networkId = useNetworkId();
  const networkPending = useNetworkQueryPending();
  // Org currency symbol comes from getCurrencySymbol (shared via CurrencyProvider).
  const { money } = useCurrency();
  const [range, setRange] = useState<string>(DEFAULT_RANGE);

  // Every Packages KPI + the table/mix derive from the package_performance
  // report, windowed by the DateChip filter (last 24h / 7 days / 30 days).
  const reportResult = useGetPerformanceReportQuery({
    variables: {
      data: {
        report: 'package_performance',
        span: rangeToSpan(range),
        networkId,
      },
    },
    skip: !networkId,
  });
  const { error, refetch } = reportResult;
  const { data, loading } = heldQuery(reportResult, networkPending);

  // Data consumed = total network data usage (DATA_USAGE filtered by network,
  // which folds the per-iccid scope to one network total), read at the selected
  // rolling span — no per-package active-assignment attribution required.
  const { data: usageData } = useGetKpiValuesQuery({
    variables: {
      data: {
        keys: [KPI_KEYS.dataUsage],
        span: rangeToSpan(range),
        networkId,
      },
    },
    skip: !networkId,
  });
  const usageBytes = kpiValue(
    usageData?.getKpiValues.values,
    KPI_KEYS.dataUsage,
  );

  const plans: Plan[] = (data?.getPerformanceReport.rows ?? []).map((r) => ({
    packageId: r.entityId,
    name: attrValue(r.attributes, 'name') ?? '—',
    price: Number(attrValue(r.attributes, 'price') ?? 0),
    revenue: cellMoney(r.cells, 'revenue'),
    sold: cellValue(r.cells, 'sold') ?? 0,
    dataVolume: Number(attrValue(r.attributes, 'data_volume') ?? 0),
    dataUnit: attrValue(r.attributes, 'data_unit') ?? '',
    duration: Number(attrValue(r.attributes, 'validity') ?? 0),
    active: (attrValue(r.attributes, 'active') ?? '').toLowerCase() === 'true',
    statusLabel: r.status,
  }));

  const totalRevenue = plans.reduce((sum, p) => sum + p.revenue, 0);
  const sharePct = (p: Plan): number =>
    totalRevenue > 0 ? Math.round((p.revenue / totalRevenue) * 100) : 0;

  const byRevenue = [...plans].sort((a, z) => z.revenue - a.revenue);
  const topPkgs = byRevenue.slice(0, 3);
  const maxRevenue = Math.max(...topPkgs.map((p) => p.revenue), 1);
  const mix = byRevenue
    .filter((p) => p.revenue > 0)
    .map((p, i) => ({
      name: p.name,
      value: p.revenue,
      color: BAR_COLORS[i % BAR_COLORS.length] ?? 'var(--uk-ac)',
    }));

  // Headline KPIs (value-only) derived from the report.
  const totalSold = plans.reduce((sum, p) => sum + p.sold, 0);
  // Best package = the best seller, shown by its spec (e.g. "1 GB / 7d").
  const bestSeller = [...plans]
    .filter((p) => p.sold > 0)
    .sort((a, z) => z.sold - a.sold)[0];
  const bestLabel = bestSeller
    ? `${formatBytes(dataVolumeToBytes(bestSeller.dataVolume, bestSeller.dataUnit))} / ${formatDuration(bestSeller.duration)}`
    : '—';

  return (
    <div className="page">
      <PageHeader
        title="Packages"
        sub="How your data packages are selling and performing."
        actions={<DateChip value={range} onChange={setRange} />}
      />
      <KpiRow
        cols={4}
        items={[
          {
            icon: 'monetization_on',
            color: 'var(--uk-beige)',
            label: 'Package revenue',
            value: error ? '—' : money(totalRevenue),
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Packages sold',
            value: error ? '—' : String(totalSold),
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Best package',
            value: error ? '—' : bestLabel,
            truncate: true,
          },
          {
            icon: 'data_usage',
            color: 'var(--uk-ac)',
            label: 'Data consumed (all packages)',
            value: usageBytes == null ? '—' : formatBytes(usageBytes),
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
                  <TableCell align="right">Sold</TableCell>
                  <TableCell align="right">Revenue</TableCell>
                  <TableCell align="right">Share</TableCell>
                  <TableCell>Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {byRevenue.map((p) => (
                  <TableRow key={p.packageId}>
                    <TableCell style={{ fontWeight: 600 }}>{p.name}</TableCell>
                    <TableCell align="right" className="tnum">
                      {money(p.price)}
                    </TableCell>
                    <TableCell align="right" className="tnum">
                      {p.sold || '—'}
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
                        status={p.active ? 'active' : 'inactive'}
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
              By revenue
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
                        {p.name}
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
                        {p.sold > 0 ? ` · ${p.sold} sold` : ''}
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
        <SectionCard
          title="Package revenue mix"
          right={
            <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>
              Share of revenue
            </span>
          }
        >
          {mix.length === 0 || error ? (
            <div
              style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}
            >
              No package revenue yet.
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column' }}>
              {mix.map((m, i) => (
                <div
                  key={`${m.name}-${i}`}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 14,
                    padding: '13px 0',
                    borderBottom:
                      i < mix.length - 1
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
                        {m.name}
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
                          {money(m.value)}
                        </b>
                        {totalRevenue > 0
                          ? ` · ${Math.round((m.value / totalRevenue) * 100)}%`
                          : ''}
                      </span>
                    </div>
                    <Meter
                      value={
                        totalRevenue > 0
                          ? Math.round((m.value / totalRevenue) * 100)
                          : 0
                      }
                      color={m.color}
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </SectionCard>
      </div>
    </div>
  );
}
