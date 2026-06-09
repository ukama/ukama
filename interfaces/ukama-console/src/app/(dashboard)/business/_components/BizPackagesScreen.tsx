/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Packages — selling & performance, wired to `commerceView.plans`
 *  (MRR/ARPU/revenue share computed server-side). */
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Meter from '@/components/Meter';

import { usePackagesDashboardQuery } from '@/client/graphql/commerce.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { EmptyState } from '@/components/EmptyState';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { sectionValue } from '@/components/SectionFallback';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import StatusBadge from '@/components/StatusBadge';
import { useUiPrefs } from '@/lib/store';

const BAR_COLORS = [
  'var(--uk-ac)',
  'var(--uk-secondary)',
  'var(--uk-success-bright)',
  'var(--uk-beige)',
  'var(--uk-orange)',
];

const money = (value?: number | null): string =>
  value == null ? '—' : `$${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;

export default function BizPackagesScreen() {
  const networkId = useUiPrefs((s) => s.networkId);
  const { data, loading, refetch } = usePackagesDashboardQuery({
    variables: { networkId },
  });
  const section = data?.commerceView.plans;
  const plans = section?.plans ?? [];
  const byRevenue = [...plans].sort((a, z) => z.revenue - a.revenue);
  const top = byRevenue[0];
  const topPkgs = byRevenue.slice(0, 3);
  const maxRevenue = Math.max(...topPkgs.map((p) => p.revenue), 1);
  const mix = byRevenue
    .filter((p) => p.revenue > 0)
    .map((p, i) => ({
      name: p.name,
      value: p.revenue,
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
            value: section?.error ? '—' : money(section?.mrr),
            sub: 'paid this month',
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'ARPU',
            value: section?.error ? '—' : money(section?.arpu),
            sub: 'avg revenue / active SIM',
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Top plan by revenue',
            value: top && top.revenue > 0 ? top.name : '—',
            sub: top && top.revenue > 0 ? `${top.revenueSharePct}% of revenue` : undefined,
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Active plans',
            value: sectionValue(
              plans.filter((p) => p.active).length || null,
              section?.error
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
          ) : section?.error ? (
            <EmptyState
              art="error"
              title="Couldn't load packages"
              sub={section.error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : plans.length === 0 ? (
            <EmptyState art="invoice" title="No packages" sub="Create a data plan to get started." />
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
                    <TableCell style={{ fontWeight: 600 }}>{p.name}</TableCell>
                    <TableCell align="right" className="tnum">
                      {p.amount} {p.currency}
                    </TableCell>
                    <TableCell align="right" className="tnum">
                      {p.attachCount ?? '—'}
                    </TableCell>
                    <TableCell align="right" className="tnum" style={{ fontWeight: 600 }}>
                      {money(p.revenue)}
                    </TableCell>
                    <TableCell align="right" className="tnum muted">
                      {p.revenueSharePct}%
                    </TableCell>
                    <TableCell>
                      <StatusBadge status={p.active ? 'active' : 'inactive'} variant="pill" />
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
          {topPkgs.length === 0 || section?.error ? (
            <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
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
                      i < topPkgs.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
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
                      <span style={{ fontSize: 13.5, fontWeight: 600 }}>{p.name}</span>
                      <span
                        className="tnum"
                        style={{ fontSize: 13, color: 'var(--uk-ink-2)', whiteSpace: 'nowrap' }}
                      >
                        <b style={{ color: 'var(--uk-ink)' }}>{money(p.revenue)}</b>
                        {p.attachCount != null ? ` · ${p.attachCount} active` : ''}
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
          {mix.length === 0 || section?.error ? (
            <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
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
