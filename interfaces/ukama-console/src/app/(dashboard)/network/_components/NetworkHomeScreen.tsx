/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Ops Home — KPIs + live network map, wired to the shared analytics home
 * queries (`getHomeKpis` / `getHomeSites`, lens = NETWORK). KPI keys live in
 * docs/analytics-backend-gaps.md and degrade to "—".
 */
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

import { useGetHomeKpisQuery } from '@/client/graphql/analytics.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { HomeLens } from '@/client/graphql/types';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import UkamaMap, { HOME_MAP_ZOOM } from '@/components/Map/UkamaMap';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { KPI_KEYS, kpiText } from '@/lib/kpis';
import { toMapSites } from '@/lib/mappers/sites';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import { pinColor } from '@/lib/status';
import { useUiPrefs } from '@/lib/store';

export default function NetworkHomeScreen() {
  const router = useRouter();
  const networkId = useUiPrefs((s) => s.networkId);
  const [sel, setSel] = useState<string | null>(null);

  // KPIs come from the analytics rollup; sites come live from the registry
  // (sitesView) so the map doesn't depend on the analytics collector.
  const { data: kpiData, loading: kpiLoading } = useGetHomeKpisQuery({
    variables: { data: { lens: HomeLens.Network, networkId } },
    skip: !networkId,
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const {
    data: sitesData,
    loading: sitesLoading,
    error: sitesError,
    refetch,
  } = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const kpis = kpiData?.getHomeKpis.kpis;
  const loading = kpiLoading || sitesLoading;

  // The home map only needs each site's name, status and coordinates.
  const mapSites = useMemo(
    () => toMapSites(sitesData?.sitesView.sites.sites ?? []),
    [sitesData?.sitesView.sites.sites],
  );

  const mapMarkers = mapSites
    .filter((s) => s.lat !== 0 || s.lng !== 0)
    .map((s) => ({
      id: s.id,
      lat: s.lat,
      lng: s.lng,
      color: pinColor(s.status),
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
              value: kpiText(kpis, KPI_KEYS.networkUptime, (v) => `${v}%`),
              sub: 'latest reading',
            },
            {
              icon: 'group',
              color: 'var(--uk-secondary)',
              label: 'Active customers',
              value: kpiText(kpis, KPI_KEYS.activeCustomers),
            },
            {
              icon: 'donut_small',
              color: 'var(--uk-beige)',
              label: 'Data volume',
              value: kpiText(kpis, KPI_KEYS.dataUsage, (v) => `${v} GB`),
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-ac)',
              label: 'Sites online',
              value: sitesError ? '—' : `${sitesOnline}/${sitesTotal}`,
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
              {sitesError ? (
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
