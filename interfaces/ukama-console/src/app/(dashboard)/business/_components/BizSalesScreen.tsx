/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Revenue — the single most important number, wired to `commerceView`. */
import Skeleton from '@mui/material/Skeleton';

import { useRevenueOverviewQuery } from '@/client/graphql/commerce.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { Delta, KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { sectionValue } from '@/components/SectionFallback';
import { useUiPrefs } from '@/lib/store';

const money = (value?: number | null): string =>
  value == null ? '—' : `$${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;

export default function BizSalesScreen() {
  const networkId = useUiPrefs((s) => s.networkId);
  const { data, loading } = useRevenueOverviewQuery({
    variables: { networkId },
  });
  const revenue = data?.commerceView.revenue;
  const plansSection = data?.commerceView.plans;
  const BAR_COLORS = [
    'var(--uk-ac)',
    'var(--uk-secondary)',
    'var(--uk-success-bright)',
    'var(--uk-beige)',
    'var(--uk-orange)',
  ];
  const byPackage = (plansSection?.plans ?? [])
    .filter((p) => p.revenue > 0)
    .sort((a, z) => z.revenue - a.revenue)
    .map((p, i) => ({
      name: p.name,
      value: p.revenue,
      color: BAR_COLORS[i % BAR_COLORS.length] ?? 'var(--uk-ac)',
    }));

  return (
    <div className="page">
      <PageHeader
        title="Revenue"
        sub="Revenue collected across your network — your single most important number."
        actions={<DateChip />}
      />

      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div className="sec-title" style={{ marginBottom: 12 }}>
          Revenue collected
        </div>
        {loading ? (
          <Skeleton variant="rounded" sx={{ height: 72 }} />
        ) : (
          <div style={{ display: 'flex', alignItems: 'flex-end', gap: 20, flexWrap: 'wrap' }}>
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
                {revenue?.error ? '—' : money(revenue?.totalPaid)}
              </span>
              {revenue?.momPct != null && (
                <Delta fontSize={14}>
                  {revenue.momPct >= 0 ? '+' : ''}
                  {revenue.momPct}% MoM
                </Delta>
              )}
            </div>
            <div style={{ flex: 1, minWidth: 20 }} />
            <div style={{ display: 'flex', gap: 30, flexWrap: 'wrap' }}>
              {(
                [
                  ['This month', revenue?.error ? '—' : money(revenue?.monthPaid)],
                  ['Last month', revenue?.error ? '—' : money(revenue?.prevMonthPaid)],
                  ['Pending', revenue?.error ? '—' : money(revenue?.totalPending)],
                ] as const
              ).map(([k, v]) => (
                <div key={k}>
                  <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{k}</div>
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
          All payments collected to date
        </div>
      </div>

      <KpiRow
        items={[
          {
            icon: 'monetization_on',
            color: 'var(--uk-beige)',
            label: 'Revenue this month',
            value: revenue?.error ? '—' : money(revenue?.monthPaid),
            sub:
              revenue?.momPct != null
                ? `${revenue.momPct >= 0 ? '+' : ''}${revenue.momPct}% vs last month`
                : undefined,
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'Pending payments',
            value: revenue?.error ? '—' : money(revenue?.totalPending),
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Collected to date',
            value: revenue?.error ? '—' : money(revenue?.totalPaid),
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Plans earning revenue',
            value: sectionValue(byPackage.length || null, plansSection?.error),
          },
        ]}
      />

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        <SectionCard title="Revenue by package">
          {plansSection?.error || byPackage.length === 0 ? (
            <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
              {plansSection?.error ? '—' : 'No package revenue yet.'}
            </div>
          ) : (
            <BarList rows={byPackage} />
          )}
        </SectionCard>
        <SectionCard title="Revenue by site">
          {/* TODO(backend-gap #10): per-site revenue rollup */}
          <div style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}>
            Per-site revenue is not available yet.
          </div>
        </SectionCard>
      </div>
    </div>
  );
}
