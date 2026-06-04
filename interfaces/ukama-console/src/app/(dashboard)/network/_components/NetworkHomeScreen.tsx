/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Ops Home — KPIs + live network map, wired to the `networkOverview`
 * composite (NetworkHome query). Section data follows the §4.5 contract:
 * loading → skeleton, failed/not-implemented → "—", empty → empty state.
 */
import { useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';

import { useNetworkHomeQuery } from '@/client/graphql/network-home.generated';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import MapPanel from '@/components/Map/MapPanel';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import type { Site } from '@/data';
import { useUiPrefs } from '@/lib/store';

export default function NetworkHomeScreen() {
  const router = useRouter();
  const networkId = useUiPrefs((s) => s.networkId);
  const [sel, setSel] = useState<string | null>(null);

  const { data, loading, refetch } = useNetworkHomeQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const overview = data?.networkOverview;

  // Map SiteDto → the map's Site view-model. Per-site uptime/subscriber
  // figures are metrics-phase data (kpis backend gap) — placeholders until
  // then; status derives from isDeactivated (honest, if coarse).
  const mapSites: Site[] = useMemo(
    () =>
      (overview?.siteStats.sites ?? []).map((s) => ({
        id: s.id,
        name: s.name,
        area: s.location,
        status: s.isDeactivated ? 'offline' : 'online',
        subs: 0,
        nodes: 0,
        uptime: 0,
        battery: 0,
        signal: null,
        data: '',
        lat: parseFloat(s.latitude) || 0,
        lng: parseFloat(s.longitude) || 0,
        plan: '',
      })),
    [overview?.siteStats.sites]
  );

  const site = mapSites.find((s) => s.id === sel);
  const sitesTotal = mapSites.length;
  const sitesOnline = mapSites.filter((s) => s.status === 'online').length;
  const subStats = overview?.subscriberStats;
  const kpisGap = overview?.kpis.error ?? null;

  return (
    <div className="page">
      <PageHeader title="Home" actions={<DateChip />} />
      {loading ? (
        <Skeleton variant="rounded" sx={{ width: '100%', height: 96 }} />
      ) : (
        <KpiRow
          cols={4}
          items={[
            {
              icon: 'network_check',
              color: 'var(--uk-success-bright)',
              // TODO(metrics-phase): networkOverview.kpis (backend gap #5)
              label: 'Network uptime',
              value: sectionValue(null, kpisGap),
              sub: 'last 30 days',
            },
            {
              icon: 'group',
              color: 'var(--uk-secondary)',
              label: 'Active customers',
              value: sectionValue(subStats?.active, subStats?.error),
            },
            {
              icon: 'donut_small',
              color: 'var(--uk-beige)',
              // TODO(metrics-phase): networkOverview.kpis (backend gap #5)
              label: 'Data volume',
              value: sectionValue(null, kpisGap),
              unit: 'GB',
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-ac)',
              label: 'Sites online',
              value: overview?.siteStats.error
                ? '—'
                : `${sitesOnline}/${sitesTotal}`,
              sub:
                sitesOnline < sitesTotal
                  ? `${sitesTotal - sitesOnline} need attention`
                  : 'all healthy',
              danger: sitesOnline < sitesTotal,
            },
          ]}
        />
      )}

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
          <div className="sec-title">Network</div>
        </div>
        <div style={{ flex: 1, position: 'relative', minHeight: 300, padding: '16px 20px' }}>
          {loading ? (
            <Skeleton variant="rounded" sx={{ width: '100%', height: '100%', minHeight: 280 }} />
          ) : overview?.siteStats.error ? (
            <div style={{ padding: 24, color: 'var(--uk-ink-2)', fontSize: 13 }}>
              Couldn&apos;t load sites.{' '}
              <Button size="small" onClick={() => refetch()}>
                Retry
              </Button>
            </div>
          ) : (
            <MapPanel sites={mapSites} selected={sel} onSelect={(s) => setSel(s.id)} />
          )}
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
                  {site.area || '—'}
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
