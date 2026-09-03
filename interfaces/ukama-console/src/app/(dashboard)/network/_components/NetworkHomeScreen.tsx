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

import { useGetKpiValuesQuery } from '@/client/graphql/analytics.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import DateChip from '@/components/DateChip';
import { DEFAULT_RANGE, rangeToSpan } from '@/lib/dateRange';
import { KpiRow } from '@/components/Kpi';
import UkamaMap, { HOME_MAP_ZOOM } from '@/components/Map/UkamaMap';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { KPI_KEYS, kpiText, kpiValue } from '@/lib/kpis';
import { formatBytes } from '@/lib/usage';
import { toMapSites } from '@/lib/mappers/sites';
import { heldQuery } from '@/lib/heldQuery';
import { POLL_LIVE_MS, visiblePoll } from '@/lib/polling';
import { sitesOnlineTile } from '@/lib/sitesOnline';
import { pinColor } from '@/lib/status';
import { useNetworkId, useNetworkQueryPending } from '@/lib/useNetworkId';

export default function NetworkHomeScreen() {
  const router = useRouter();
  const networkId = useNetworkId();
  const networkPending = useNetworkQueryPending();
  const [sel, setSel] = useState<string | null>(null);
  const [range, setRange] = useState<string>(DEFAULT_RANGE);

  // KPIs come from the analytics rollup; sites come live from the registry
  // (sitesView) so the map doesn't depend on the analytics collector.
  const [fetchedAt, setFetchedAt] = useState(() => new Date());
  const kpiResult = useGetKpiValuesQuery({
    variables: {
      data: {
        keys: [
          KPI_KEYS.networkUptime,
          KPI_KEYS.activeCustomers,
          KPI_KEYS.dataUsage,
          KPI_KEYS.sitesOnline,
        ],
        span: rangeToSpan(range),
        networkId,
      },
    },
    skip: !networkId,
    onCompleted: () => setFetchedAt(new Date()),
    ...visiblePoll(POLL_LIVE_MS, true),
  });
  const sitesResult = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
    ...visiblePoll(POLL_LIVE_MS, true),
  });
  const { error: sitesError, refetch } = sitesResult;
  const { data: kpiData, loading: kpiLoading } = heldQuery(
    kpiResult,
    networkPending,
  );
  const { data: sitesData, loading: sitesLoading } = heldQuery(
    sitesResult,
    networkPending,
  );
  const kpis = kpiData?.getKpiValues.values;
  const loading = kpiLoading || sitesLoading;

  // `sitesView` resolves per section, so a failed sites read arrives as an
  // empty list with no top-level error — that is unknown, not zero.
  const sitesSection = sitesData?.sitesView.sites;
  const siteRows =
    sitesError || sitesSection?.error ? undefined : sitesSection?.sites;

  // The home map only needs each site's name, status and coordinates.
  const mapSites = useMemo(() => toMapSites(siteRows ?? []), [siteRows]);

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
  // Sites online comes from the analytics SITES_ONLINE KPI only (no registry
  // fallback); undefined when the KPI hasn't been emitted. The KPI counts a
  // site only when every one of its tnode/anode/cnode is online, so a site
  // with an offline node — or with no nodes registered yet — is excluded and
  // the count reads below the registry site total. The total comes from a
  // different response, so `sitesOnlineTile` renders the pair only when it is
  // coherent (`siteRows` is undefined, not [], when unknown).
  const sitesTile = sitesOnlineTile(
    kpiValue(kpis, KPI_KEYS.sitesOnline),
    siteRows?.length,
  );

  return (
    <div className="page">
      <PageHeader
        title="Home"
        actions={<DateChip value={range} onChange={setRange} />}
        fetchedAt={fetchedAt}
      />
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
              // Computed KPI only, one decimal (e.g. 15.1%); "—" when absent.
              value: kpiText(
                kpis,
                KPI_KEYS.networkUptime,
                (v) => `${v.toFixed(1)}%`,
              ),
              // NETWORK_UPTIME allows AVG/MIN/MAX (no LAST), so the gateway
              // defaults to AVG: this is the share of site-time up across the
              // selected window, not a spot reading.
              sub: `average over ${range.toLowerCase()}`,
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
              value: kpiText(kpis, KPI_KEYS.dataUsage, (v) =>
                formatBytes(v),
              ),
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-ac)',
              label: 'Sites online',
              value: sitesTile.text,
              // A site counts as online only when all of its nodes are
              // connected, so the shortfall is "not fully online" rather than
              // "down". No sub-line at all while the KPI is absent — the old
              // copy claimed "all healthy" for a missing reading.
              sub:
                sitesTile.offline == null
                  ? undefined
                  : sitesTile.offline > 0
                    ? `${sitesTile.offline} not fully online`
                    : 'all sites online',
              danger: sitesTile.offline != null && sitesTile.offline > 0,
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
