/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Revenue — the single most important number (biz-home.jsx BizSales). */
import ArrowUpwardRounded from '@mui/icons-material/ArrowUpwardRounded';
import BarList from '@/components/BarList';
import { LineChart } from '@/components/charts';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { BIZ_SALES } from '@/data';

export default function BizSalesScreen() {
  const b = BIZ_SALES;
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
              $48,210
            </span>
            <span className="stat-delta up" style={{ fontSize: 14 }}>
              <ArrowUpwardRounded sx={{ fontSize: 16 }} />
              +22% YoY
            </span>
          </div>
          <div style={{ flex: 1, minWidth: 20 }} />
          <div style={{ display: 'flex', gap: 30, flexWrap: 'wrap' }}>
            {(
              [
                ['This month', '$1,284'],
                ['Last month', '$1,206'],
                ['Monthly average', '$4,017'],
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
        <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)', marginTop: 8 }}>
          Collected over the last 12 months
        </div>
      </div>

      <KpiRow items={b.kpis} />

      <SectionCard
        title="Revenue trend"
        style={{ marginBottom: 'var(--uk-gap)' }}
        right={<span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>Last 9 weeks</span>}
      >
        <LineChart data={b.trend} height={240} />
      </SectionCard>

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        <SectionCard title="Revenue by site">
          <BarList rows={b.bySite} />
        </SectionCard>
        <SectionCard title="Revenue by package">
          <BarList rows={b.byPackage} />
        </SectionCard>
      </div>
    </div>
  );
}
