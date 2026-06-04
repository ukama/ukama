/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';

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
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Resource</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Site</TableCell>
              <TableCell align="right">Customers affected</TableCell>
              <TableCell>Revenue context</TableCell>
              <TableCell>Last updated</TableCell>
              <TableCell>Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {b.rows.map((r) => (
              <TableRow key={r.res}>
                <TableCell style={{ fontWeight: 600 }}>{r.res}</TableCell>
                <TableCell className="muted">{r.type}</TableCell>
                <TableCell>
                  <StatusBadge status={r.status} variant="pill" />
                </TableCell>
                <TableCell>{r.site}</TableCell>
                <TableCell align="right" className="tnum">{r.affected}</TableCell>
                <TableCell className="muted">{r.context}</TableCell>
                <TableCell className="muted">{r.updated}</TableCell>
                <TableCell>
                  <button type="button" className="link">
                    Open
                  </button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
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
