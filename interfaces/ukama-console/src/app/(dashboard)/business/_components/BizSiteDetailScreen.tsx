/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Business site detail — business summary above admin tabs (biz-customers.jsx). */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import ArrowBackRounded from '@mui/icons-material/ArrowBackRounded';
import ArrowUpwardRounded from '@mui/icons-material/ArrowUpwardRounded';
import DashboardCustomizeRounded from '@mui/icons-material/DashboardCustomizeRounded';
import SiteMap from '@/components/Map/SiteMap';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { BIZ_SITES, BIZ_SITE_DETAIL } from '@/data';

function SiteKpiTile({
  label,
  value,
  delta,
  dir,
  sub,
}: {
  label: string;
  value: string;
  delta?: string;
  dir?: 'up' | 'down';
  sub?: string;
}) {
  return (
    <div style={{ borderLeft: '1px solid var(--uk-line)', padding: '2px 0 2px 20px', minWidth: 0 }}>
      <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}>{label}</div>
      <div
        className="tnum"
        style={{
          fontFamily: 'var(--font-display)',
          fontSize: 24,
          fontWeight: 500,
          margin: '3px 0 2px',
        }}
      >
        {value}
      </div>
      {delta != null ? (
        <span className={`stat-delta ${dir === 'down' ? 'down' : 'up'}`}>
          <ArrowUpwardRounded sx={{ fontSize: 12 }} />
          {delta}
        </span>
      ) : (
        <span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{sub}</span>
      )}
    </div>
  );
}

export default function BizSiteDetailScreen({ siteId }: { siteId: string }) {
  const router = useRouter();
  const d = BIZ_SITE_DETAIL;
  const s = BIZ_SITES.find((x) => x.id === siteId) ?? BIZ_SITES[0];
  const [tab, setTab] = useState('Overview');
  if (!s) return null;

  return (
    <div className="page">
      <PageHeader
        crumb={['Sites', s.name]}
        title="Site detail"
        sub="Business summary above the existing admin / detail tabs."
        actions={
          <Button
            variant="outlined"
            startIcon={<ArrowBackRounded />}
            onClick={() => router.push('/business/sites')}
          >
            Back to sites
          </Button>
        }
      />

      <div
        className="card card-pad"
        style={{
          marginBottom: 'var(--uk-gap)',
          display: 'flex',
          alignItems: 'center',
          gap: 24,
          flexWrap: 'wrap',
        }}
      >
        <div style={{ minWidth: 200, flex: '0 0 auto' }}>
          <div style={{ fontFamily: 'var(--font-display)', fontSize: 30, fontWeight: 500 }}>
            {s.name}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: 8 }}>
            <StatusBadge status={s.status} variant="pill" />
            <span style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>{d.meta}</span>
          </div>
        </div>
        <div
          style={{ display: 'flex', gap: 24, flex: 1, justifyContent: 'flex-end', flexWrap: 'wrap' }}
        >
          {d.kpis.map((k, i) => (
            <SiteKpiTile key={i} {...k} />
          ))}
        </div>
      </div>

      <div className="tabs">
        {d.tabs.map((t) => (
          <button
            key={t}
            type="button"
            className={`tab${tab === t ? ' on' : ''}`}
            onClick={() => setTab(t)}
          >
            {t}
          </button>
        ))}
      </div>

      {tab === 'Overview' ? (
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1.1fr', alignItems: 'stretch' }}>
          <SiteMap sites={BIZ_SITES} title="Coverage / site map" height={330} selected={s.id} />
          <SectionCard title="Resources">
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Resource</th>
                  <th>Status</th>
                  <th>Last seen</th>
                  <th>Issue</th>
                </tr>
              </thead>
              <tbody>
                {d.resources.map((r) => (
                  <tr key={r.res} className="static">
                    <td style={{ fontWeight: 600 }}>{r.res}</td>
                    <td>
                      <StatusBadge status={r.status} variant="pill" />
                    </td>
                    <td className="muted">{r.seen}</td>
                    <td
                      style={{
                        color:
                          r.issue !== '—' ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-ink-3)',
                      }}
                    >
                      {r.issue}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </SectionCard>
        </div>
      ) : (
        <div className="card" style={{ padding: '56px 24px', textAlign: 'center', color: 'var(--uk-ink-3)' }}>
          <DashboardCustomizeRounded sx={{ fontSize: 42 }} />
          <div
            style={{
              fontFamily: 'var(--font-display)',
              fontSize: 18,
              fontWeight: 500,
              marginTop: 12,
              color: 'var(--uk-ink)',
            }}
          >
            {tab}
          </div>
          <div style={{ fontSize: 13.5, marginTop: 6 }}>
            Detailed {tab.toLowerCase()} view for {s.name}.
          </div>
        </div>
      )}
    </div>
  );
}
