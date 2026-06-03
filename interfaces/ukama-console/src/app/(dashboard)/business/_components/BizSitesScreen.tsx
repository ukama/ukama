/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

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
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Site</th>
                  <th>Status</th>
                  <th className="num">Revenue</th>
                  <th className="num">Customers</th>
                  <th className="num">Data sold</th>
                  <th className="num">Uptime</th>
                  <th>Top package</th>
                  <th>Issue</th>
                </tr>
              </thead>
              <tbody>
                {BIZ_SITES.map((s) => (
                  <tr
                    key={s.id}
                    role="button"
                    tabIndex={0}
                    onClick={() => go(s.id)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') go(s.id);
                    }}
                  >
                    <td style={{ fontWeight: 600 }}>{s.name}</td>
                    <td>
                      <StatusBadge status={s.status} variant="pill" />
                    </td>
                    <td className="num tnum" style={{ fontWeight: 600 }}>
                      ${s.revenue.toLocaleString()}
                    </td>
                    <td className="num tnum">{s.customers}</td>
                    <td className="num tnum muted">{s.data}</td>
                    <td className="num tnum">{s.uptime}%</td>
                    <td className="muted">{s.top}</td>
                    <td
                      style={{
                        color: s.issue ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-ink-3)',
                      }}
                    >
                      {s.issue ?? '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
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
