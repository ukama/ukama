/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Revenue — the single most important number, wired to the analytics service
 *  (`getSalesOverview`). Headline figures come from the keyed `kpis` array;
 *  the by-package / by-site breakdowns from `revenueByPackage` / `revenueBySite`.
 *  KPI keys are listed in docs/analytics-backend-gaps.md and degrade to "—". */
import Skeleton from '@mui/material/Skeleton';

import { useGetSalesOverviewQuery } from '@/client/graphql/analytics.generated';
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { useCurrency } from '@/lib/currency';
import { KPI_KEYS, kpiAmount } from '@/lib/kpis';
import { useUiPrefs } from '@/lib/store';

// KPI keys this screen reads. Centralised so a backend rename is a one-line
const BAR_COLORS = [
  'var(--uk-ac)',
  'var(--uk-secondary)',
  'var(--uk-success-bright)',
  'var(--uk-beige)',
  'var(--uk-orange)',
];

const toBars = (rows: { name?: string | null; value: number }[]) =>
  rows
    .filter((r) => r.value > 0)
    .sort((a, z) => z.value - a.value)
    .map((r, i) => ({
      name: r.name ?? '—',
      value: r.value,
      color: BAR_COLORS[i % BAR_COLORS.length] ?? 'var(--uk-ac)',
    }));

export default function BizSalesScreen() {
  const networkId = useUiPrefs((s) => s.networkId);
  // Org currency symbol from getCurrencySymbol (shared via CurrencyProvider).
  const { symbol } = useCurrency();
  const money = (value?: number | null): string =>
    value == null
      ? '—'
      : `${symbol}${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;
  const { data, loading, error } = useGetSalesOverviewQuery({
    variables: { data: { networkId } },
  });
  const overview = data?.getSalesOverview;
  const kpis = overview?.kpis;

  const byPackage = toBars(overview?.revenueByPackage ?? []);
  const bySite = toBars(overview?.revenueBySite ?? []);

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
                {error ? '—' : kpiAmount(kpis, KPI_KEYS.revenueCollected, money)}
              </span>
              {/* MoM delta hidden until the backend provides a reliable
                  previous-period revenue (current value is unreliable). */}
            </div>
            <div style={{ flex: 1, minWidth: 20 }} />
            <div style={{ display: 'flex', gap: 30, flexWrap: 'wrap' }}>
              {(
                [
                  [
                    'This month',
                    error ? '—' : kpiAmount(kpis, KPI_KEYS.revenueMonth, money),
                  ],
                  [
                    'Last month',
                    error ? '—' : kpiAmount(kpis, KPI_KEYS.revenuePrevMonth, money),
                  ],
                  [
                    'Pending',
                    error ? '—' : kpiAmount(kpis, KPI_KEYS.revenuePending, money),
                  ],
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
          All payments collected to date
        </div>
      </div>

      <KpiRow
        items={[
          {
            icon: 'monetization_on',
            color: 'var(--uk-beige)',
            label: 'Revenue this month',
            value: error ? '—' : kpiAmount(kpis, KPI_KEYS.revenueMonth, money),
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'Pending payments',
            value: error ? '—' : kpiAmount(kpis, KPI_KEYS.revenuePending, money),
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Collected to date',
            value: error ? '—' : kpiAmount(kpis, KPI_KEYS.revenueCollected, money),
          },
          {
            icon: 'sell',
            color: 'var(--uk-secondary)',
            label: 'Plans earning revenue',
            value: error ? '—' : String(byPackage.length),
          },
        ]}
      />

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        <SectionCard title="Revenue by package">
          {error || byPackage.length === 0 ? (
            <div
              style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}
            >
              {error ? '—' : 'No package revenue yet.'}
            </div>
          ) : (
            <BarList rows={byPackage} />
          )}
        </SectionCard>
        <SectionCard title="Revenue by site">
          {error || bySite.length === 0 ? (
            <div
              style={{ padding: 24, fontSize: 13, color: 'var(--uk-ink-3)' }}
            >
              {error ? '—' : 'No site revenue yet.'}
            </div>
          ) : (
            <BarList rows={bySite} />
          )}
        </SectionCard>
      </div>
    </div>
  );
}
