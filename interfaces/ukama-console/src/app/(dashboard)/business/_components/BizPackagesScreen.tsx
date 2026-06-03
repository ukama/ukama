/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Packages — how the data packages sell and perform (biz-home.jsx). */
import BarList from '@/components/BarList';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { BIZ_HOME, BIZ_PACKAGES, PLANS } from '@/data';

export default function BizPackagesScreen() {
  const b = BIZ_PACKAGES;
  const total = PLANS.reduce((s, p) => s + p.subs, 0);
  const mrr = PLANS.reduce((s, p) => s + p.subs * p.price, 0);
  const arpu = mrr / total;
  const byRev = [...PLANS]
    .map((p) => ({ name: p.name, rev: p.subs * p.price }))
    .sort((a, z) => z.rev - a.rev);
  const top = byRev[0];
  const topShare = top ? Math.round((top.rev / mrr) * 100) : 0;
  const topPkgs = BIZ_HOME.topPackages;
  const maxSold = Math.max(...topPkgs.map((p) => p.sold));

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
            value: `$${mrr.toLocaleString()}`,
            sub: 'estimated · all plans',
          },
          {
            icon: 'payments',
            color: 'var(--uk-ac)',
            label: 'ARPU',
            value: `$${arpu.toFixed(2)}`,
            sub: 'avg revenue / customer',
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-success-bright)',
            label: 'Top plan by revenue',
            value: top ? top.name : '—',
            sub: `${topShare}% of revenue`,
          },
          {
            icon: 'group',
            color: 'var(--uk-secondary)',
            label: 'Customers on a plan',
            value: total.toLocaleString(),
            sub: `across ${PLANS.length} plans`,
          },
        ]}
      />

      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div className="sec-head">
          <div className="sec-title">Package performance</div>
        </div>
        <div className="tbl-wrap">
          <table className="tbl">
            <thead>
              <tr className="static">
                <th>Package</th>
                <th className="num">Price</th>
                <th>Validity</th>
                <th className="num">Sold</th>
                <th className="num">Revenue</th>
                <th className="num">Data used</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {b.rows.map((r) => (
                <tr key={r.pkg} className="static">
                  <td style={{ fontWeight: 600 }}>{r.pkg}</td>
                  <td className="num tnum">{r.price}</td>
                  <td className="muted">{r.validity}</td>
                  <td className="num tnum">{r.sold}</td>
                  <td className="num tnum" style={{ fontWeight: 600 }}>
                    {r.revenue}
                  </td>
                  <td className="num tnum muted">{r.data}</td>
                  <td>
                    <StatusBadge status={r.status} variant="pill" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        <SectionCard
          title="Top packages"
          right={
            <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>
              Best sellers this month
            </span>
          }
        >
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            {topPkgs.map((p, i) => (
              <div
                key={p.name}
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
                      <b style={{ color: 'var(--uk-ink)' }}>${p.revenue.toLocaleString()}</b> ·{' '}
                      {p.sold} sold
                    </span>
                  </div>
                  <div className="meter">
                    <span
                      style={{
                        width: Math.round((p.sold / maxSold) * 100) + '%',
                        background: p.color,
                      }}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </SectionCard>
        <SectionCard title="Package revenue mix">
          <BarList rows={b.mix} />
        </SectionCard>
      </div>
    </div>
  );
}
