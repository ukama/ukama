/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Network health framed by business impact (biz-ops.jsx BizNetwork). */
import DateChip from '@/components/DateChip';
import FeedRow from '@/components/FeedRow';
import { KpiRow } from '@/components/Kpi';
import SiteMap from '@/components/Map/SiteMap';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { BIZ_NETWORK, BIZ_SITES } from '@/data';

export default function BizNetworkScreen() {
  const b = BIZ_NETWORK;
  return (
    <div className="page">
      <PageHeader
        title="Network health"
        sub="Is the network healthy, and what is the business impact?"
        actions={<DateChip />}
      />
      <KpiRow items={b.kpis} />
      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <table className="tbl">
          <thead>
            <tr className="static">
              <th>Resource</th>
              <th>Type</th>
              <th>Status</th>
              <th>Site</th>
              <th className="num">Customers affected</th>
              <th>Revenue context</th>
              <th>Last updated</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {b.rows.map((r) => (
              <tr key={r.res} className="static">
                <td style={{ fontWeight: 600 }}>{r.res}</td>
                <td className="muted">{r.type}</td>
                <td>
                  <StatusBadge status={r.status} variant="pill" />
                </td>
                <td>{r.site}</td>
                <td className="num tnum">{r.affected}</td>
                <td className="muted">{r.context}</td>
                <td className="muted">{r.updated}</td>
                <td>
                  <button type="button" className="link">
                    Open
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="tile-grid" style={{ gridTemplateColumns: '1.4fr 1fr', alignItems: 'stretch' }}>
        <SiteMap sites={BIZ_SITES} title="Network health map" height={300} />
        <SectionCard title="Health summary">
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            {b.summary.map((s, i) => (
              <div
                key={i}
                style={{
                  borderBottom:
                    i < b.summary.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
                }}
              >
                <FeedRow tone={s.tone} title={s.title} detail={s.detail} />
              </div>
            ))}
          </div>
        </SectionCard>
      </div>
    </div>
  );
}
