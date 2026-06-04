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

/** Business Sites — which sites perform well as businesses (biz-customers.jsx). */
import { useRouter } from 'next/navigation';
import TableFooter from '@/components/data-table/TableFooter';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import SiteMap from '@/components/Map/SiteMap';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { BIZ_SITES } from '@/data';
import { useFirstLoad } from '@/lib/useFirstLoad';

export default function BizSitesScreen() {
  const router = useRouter();
  const loading = useFirstLoad('biz-sites');
  const go = (id: string) => router.push(`/business/sites/${id}`);

  return (
    <div className="page">
      <PageHeader
        title="Sites"
        sub="Which sites are performing well as businesses?"
        actions={<DateChip />}
      />
      <KpiRow
        items={[
          { label: 'Total sites', value: '7', sub: '6 online' },
          { label: 'Sites online', value: '6/7', sub: 'one offline', danger: true },
          { label: 'Site revenue', value: '$3,420', sub: 'this month' },
          { label: 'Active customers', value: '376', sub: 'across sites' },
        ]}
      />
      <div className="card card-pad" style={{ marginBottom: 'var(--uk-gap)' }}>
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={8} rows={6} />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Site</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell align="right">Revenue</TableCell>
                  <TableCell align="right">Customers</TableCell>
                  <TableCell align="right">Data sold</TableCell>
                  <TableCell align="right">Uptime</TableCell>
                  <TableCell>Top package</TableCell>
                  <TableCell>Issue</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {BIZ_SITES.map((s) => (
                  <TableRow
                    hover
                    sx={{ cursor: 'pointer' }}
                    key={s.id}
                    role="button"
                    tabIndex={0}
                    onClick={() => go(s.id)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') go(s.id);
                    }}
                  >
                    <TableCell style={{ fontWeight: 600 }}>{s.name}</TableCell>
                    <TableCell>
                      <StatusBadge status={s.status} variant="pill" />
                    </TableCell>
                    <TableCell align="right" className="tnum" style={{ fontWeight: 600 }}>
                      ${s.revenue.toLocaleString()}
                    </TableCell>
                    <TableCell align="right" className="tnum">{s.customers}</TableCell>
                    <TableCell align="right" className="tnum muted">{s.data}</TableCell>
                    <TableCell align="right" className="tnum">{s.uptime}%</TableCell>
                    <TableCell className="muted">{s.top}</TableCell>
                    <TableCell
                      style={{
                        color: s.issue ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-ink-3)',
                      }}
                    >
                      {s.issue ?? '—'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && <TableFooter count={BIZ_SITES.length} noun="sites" />}
      </div>
      <SiteMap
        sites={BIZ_SITES}
        title="Site mini-map"
        height={230}
        onSelect={(s) => go(s.id)}
      />
    </div>
  );
}
