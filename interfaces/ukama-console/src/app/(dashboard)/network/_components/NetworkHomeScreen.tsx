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
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

import { useNetworkHomeQuery } from '@/client/graphql/network-home.generated';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import UkamaMap, { HOME_MAP_ZOOM } from '@/components/Map/UkamaMap';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import type { Site } from '@/data';
import { normalizeCoords } from '@/lib/geo';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import { useUiPrefs } from '@/lib/store';

export default function NetworkHomeScreen() {
  const router = useRouter();
  const networkId = useUiPrefs((s) => s.networkId);
  const [sel, setSel] = useState<string | null>(null);

  const { data, loading, refetch } = useNetworkHomeQuery({
    variables: { networkId },
    skip: !networkId,
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const overview = data?.networkOverview;
  const kpiByKey = new Map(
    (overview?.kpis.metrics ?? []).map((m) => [m.key, m]),
  );
  const kpiValue = (key: string, unit = ''): string => {
    const entry = kpiByKey.get(key);
    if (!entry || !entry.success) return '—';
    return `${Math.round(entry.value * 100) / 100}${unit}`;
  };

  // Map SiteDto → the map's Site view-model. Per-site uptime/subscriber
  // figures are metrics-phase data (kpis backend gap) — placeholders until
  // then; status derives from isDeactivated (honest, if coarse).
  const mapSites: Site[] = useMemo(
    () =>
      (overview?.siteStats.sites ?? []).map((s) => {
        const geo = normalizeCoords(s.latitude, s.longitude);
        return {
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
          lat: geo?.lat ?? 0,
          lng: geo?.lng ?? 0,
          plan: '',
        };
      }),
    [overview?.siteStats.sites],
  );

  const SITE_PIN: Record<string, string> = {
    online: 'var(--uk-success-bright)',
    degraded: 'var(--uk-warning)',
    offline: 'var(--uk-error)',
  };
  const mapMarkers = mapSites
    .filter((s) => s.lat !== 0 || s.lng !== 0)
    .map((s) => ({
      id: s.id,
      lat: s.lat,
      lng: s.lng,
      color: SITE_PIN[s.status] ?? 'var(--uk-ac)',
      popup: (
        <div style={{ minWidth: 120 }}>
          <div style={{ fontWeight: 600, marginBottom: 2 }}>{s.name}</div>
          <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{s.area}</div>
        </div>
      ),
    }));

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
              label: 'Network uptime',
              value: kpisGap ? '—' : kpiValue('network_uptime', '%'),
              sub: 'latest reading',
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
              label: 'Data volume',
              value: kpisGap ? '—' : kpiValue('data_usage'),
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
        style={{
          flex: 1,
          minHeight: 420,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        {loading ? (
          <Skeleton
            variant="rounded"
            sx={{ width: '100%', height: '100%', minHeight: 280, mt: 1 }}
          />
        ) : (
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
              style={{
                padding: '16px 20px 12px',
                margin: 0,
                borderBottom: '1px solid var(--uk-line-soft)',
              }}
            >
              <div className="sec-title">Network</div>
            </div>
            <div
              style={{
                flex: 1,
                minHeight: 300,
              }}
            >
              {overview?.siteStats.error ? (
                <div
                  style={{
                    padding: 24,
                    color: 'var(--uk-ink-2)',
                    fontSize: 13,
                  }}
                >
                  Couldn&apos;t load sites.{' '}
                  <Button size="small" onClick={() => refetch()}>
                    Retry
                  </Button>
                </div>
              ) : (
                <UkamaMap
                  markers={mapMarkers}
                  onSelect={setSel}
                  zoom={HOME_MAP_ZOOM}
                  height="100%"
                />
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
                    <div
                      style={{ display: 'flex', alignItems: 'center', gap: 9 }}
                    >
                      <span
                        style={{
                          fontFamily: 'var(--font-display)',
                          fontSize: 15,
                          fontWeight: 500,
                        }}
                      >
                        {site.name}
                      </span>
                      <StatusBadge status={site.status} />
                    </div>
                    <div
                      style={{
                        fontSize: 12.5,
                        color: 'var(--uk-ink-2)',
                        marginTop: 3,
                      }}
                    >
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
                  <Button
                    size="small"
                    color="inherit"
                    sx={{ color: 'var(--uk-ink-3)' }}
                    onClick={() => setSel(null)}
                  >
                    Clear
                  </Button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
