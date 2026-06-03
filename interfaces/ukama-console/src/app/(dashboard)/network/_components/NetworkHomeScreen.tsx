/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Ops Home — KPIs + live network map with inline inspector (home.jsx). */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import { KpiRow } from '@/components/Kpi';
import MapPanel from '@/components/Map/MapPanel';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { SITES } from '@/data';

export default function NetworkHomeScreen() {
  const router = useRouter();
  const [sel, setSel] = useState<string | null>(null);
  const site = SITES.find((s) => s.id === sel);
  const online = SITES.filter((s) => s.status === 'online').length;
  const avgUptime = (
    SITES.reduce((a, s) => a + s.uptime, 0) / SITES.length
  ).toFixed(1);

  return (
    <div className="page">
      <PageHeader title="Home" />
      <KpiRow
        cols={4}
        items={[
          {
            icon: 'network_check',
            color: 'var(--uk-success-bright)',
            label: 'Network uptime',
            value: `${avgUptime}%`,
            sub: 'last 30 days',
          },
          {
            icon: 'group',
            color: 'var(--uk-secondary)',
            label: 'Active customers',
            value: '1,284',
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-beige)',
            label: 'Data volume',
            value: '312',
            unit: 'GB',
          },
          {
            icon: 'cell_tower',
            color: 'var(--uk-ac)',
            label: 'Sites online',
            value: `${online}/${SITES.length}`,
            sub:
              online < SITES.length
                ? `${SITES.length - online} need attention`
                : 'all healthy',
            danger: online < SITES.length,
          },
        ]}
      />

      <div
        className="card"
        style={{
          padding: 0,
          overflow: 'hidden',
          flex: 1,
          minHeight: 380,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div
          className="sec-head"
          style={{ padding: '16px 20px 12px', margin: 0, borderBottom: '1px solid var(--uk-line-soft)' }}
        >
          <div className="sec-title">Network map</div>
        </div>
        <div style={{ flex: 1, position: 'relative', minHeight: 300, padding: '16px 20px' }}>
          <MapPanel sites={SITES} selected={sel} onSelect={(s) => setSel(s.id)} />
        </div>
        {site && (
          <div style={{ padding: '0 20px 16px' }}>
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 14,
                background: 'var(--uk-page)',
                borderRadius: 10,
                padding: '12px 14px',
              }}
            >
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 9 }}>
                  <span style={{ fontFamily: 'var(--font-display)', fontSize: 15, fontWeight: 500 }}>
                    {site.name}
                  </span>
                  <StatusBadge status={site.status} />
                </div>
                <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 3 }}>
                  {site.subs} customers · {site.uptime}% uptime
                </div>
              </div>
              <Button
                size="small"
                variant="contained"
                onClick={() => router.push(`/network/sites/${site.id}`)}
              >
                Open site
              </Button>
              <Button size="small" color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={() => setSel(null)}>
                Clear
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
