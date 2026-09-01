/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Business Home — KPIs + full-height sites map, wired to the analytics service:
 * `getKpiValues` (headline KPIs) and `getBusinessSites` (per-site coordinates /
 * status). The KPI strip shows the four "at a glance" numbers: revenue,
 * customers, data sold and network uptime. KPI keys live in
 * src/lib/kpis.ts; any not-yet-emitted key degrades to "—".
 */
import ListAltRounded from '@mui/icons-material/ListAltRounded';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

import { useGetKpiValuesQuery } from '@/client/graphql/analytics.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import AppModal from '@/components/AppModal';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import { DEFAULT_RANGE, type KpiSpan, rangeToSpan } from '@/lib/dateRange';
import { StatusDot } from '@/components/Map/SiteMap';
import UkamaMap, { HOME_MAP_ZOOM } from '@/components/Map/UkamaMap';
import PageHeader from '@/components/PageHeader';
import { useCurrency } from '@/lib/currency';
import {
  KPI_KEYS,
  kpiAmount,
  kpiChangeAbs,
  kpiDelta,
  kpiText,
  kpiValue,
} from '@/lib/kpis';
import { type MapSite, toMapSites } from '@/lib/mappers/sites';
import { sitesOnlineTile } from '@/lib/sitesOnline';
import { pinColor } from '@/lib/status';
import { formatBytes } from '@/lib/usage';
import { useNetworkId } from '@/lib/useNetworkId';

// Period phrasing for the KPI trend/sub lines, keyed by the selected rolling
// span so the wording matches the DateChip ("Last 24h" -> "vs prev 24h").
const TREND_SUFFIX: Record<KpiSpan, string> = {
  last_24h: 'vs prev 24h',
  last_7d: 'vs prev 7 days',
  last_30d: 'vs prev 30 days',
};
const PERIOD_LABEL: Record<KpiSpan, string> = {
  last_24h: 'last 24h',
  last_7d: 'last 7 days',
  last_30d: 'last 30 days',
};

function SiteSummaryList({
  sites,
  onSite,
}: {
  sites: MapSite[];
  onSite: (s: MapSite) => void;
}) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      {sites.map((s, i) => (
        <div
          key={s.id}
          role="button"
          tabIndex={0}
          onClick={() => onSite(s)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') onSite(s);
          }}
          style={{
            display: 'flex',
            alignItems: 'flex-start',
            gap: 11,
            padding: '13px 0',
            cursor: 'pointer',
            borderBottom:
              i < sites.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
          }}
        >
          <span style={{ marginTop: 4, display: 'inline-flex' }}>
            <StatusDot status={s.status} />
          </span>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 13.5, fontWeight: 600 }}>{s.name}</div>
            <div
              style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 1 }}
            >
              {s.status === 'offline' ? 'Offline' : 'Online'}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default function BizHomeScreen() {
  const router = useRouter();
  const networkId = useNetworkId();
  const [showSummary, setShowSummary] = useState(false);
  const [range, setRange] = useState<string>(DEFAULT_RANGE);
  const span = rangeToSpan(range);
  // Org currency symbol from getCurrencySymbol (shared via CurrencyProvider).
  const { money } = useCurrency();

  // KPIs come from the analytics rollup; sites come live from the registry
  // (sitesView) so the map doesn't depend on the analytics collector. No `op`
  // is sent — each key resolves to its spec default (revenue -> SUM,
  // customers -> LAST, data_sold -> SUM, network_uptime -> AVG).
  const { data: homeData, loading: homeLoading } = useGetKpiValuesQuery({
    variables: {
      data: {
        keys: [
          KPI_KEYS.revenue,
          KPI_KEYS.customers,
          KPI_KEYS.dataSold,
          KPI_KEYS.networkUptime,
          KPI_KEYS.sitesOnline,
        ],
        span,
        networkId,
      },
    },
    skip: !networkId,
  });
  // Last-30-days data-sold total for the "… last 30 days" sub-line. Skipped
  // when the strip is already showing the 30-day window (headline == sub).
  const { data: monthData } = useGetKpiValuesQuery({
    variables: {
      data: { keys: [KPI_KEYS.dataSold], span: 'last_30d', networkId },
    },
    skip: !networkId || span === 'last_30d',
  });
  const {
    data: sitesCurrent,
    previousData: sitesPrevious,
    loading: sitesLoading,
    error: sitesError,
  } = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
  });
  // Hold the last delivered response while a cache-TTL expiry refetches, so
  // the site total never momentarily reads as zero.
  const sitesData = sitesCurrent ?? sitesPrevious;
  const kpis = homeData?.getKpiValues.values;
  const monthKpis = monthData?.getKpiValues.values;
  const loading = homeLoading || sitesLoading;

  // `sitesView` resolves per section, so a failed sites read arrives as an
  // empty list with no top-level error — that is unknown, not zero.
  const sitesSection = sitesData?.sitesView.sites;
  const siteRows =
    sitesError || sitesSection?.error ? undefined : sitesSection?.sites;
  const sites = useMemo(() => toMapSites(siteRows ?? []), [siteRows]);
  // Sites online comes from the analytics SITES_ONLINE KPI only — undefined
  // when it hasn't been emitted, same as the Ops home screen. The old registry
  // fallback counted every site that wasn't deactivated, which no longer means
  // the same thing: the KPI counts a site only when every one of its
  // tnode/anode/cnode is online, so the fallback silently over-reported. The
  // KPI and the total come from separate responses, so `sitesOnlineTile`
  // renders the pair only when it is coherent.
  const sitesTile = sitesOnlineTile(
    kpiValue(kpis, KPI_KEYS.sitesOnline),
    siteRows?.length,
  );
  // The business site-detail page was removed; drill into the canonical
  // Network site detail instead.
  const goSite = (id: string) => router.push(`/network/sites/${id}`);

  // --- KPI strip values ---
  // Revenue (SUM) with a period-over-period trend arrow.
  const revDelta = kpiDelta(kpis, KPI_KEYS.revenue);
  // Customers (LAST) with the absolute change over the period.
  const custDelta = kpiChangeAbs(kpis, KPI_KEYS.customers);
  // Data sold (SUM, bytes) for the period; monthly total for the sub-line.
  const soldBytes = kpiValue(kpis, KPI_KEYS.dataSold);
  const soldMonthBytes = kpiValue(monthKpis, KPI_KEYS.dataSold);
  // Network uptime (AVG, percent) — the computed KPI only; "—" when absent.
  const uptime = kpiValue(kpis, KPI_KEYS.networkUptime);

  // Data sold display: auto data unit for real values (e.g. "1 GB", "1.8 TB"),
  // a bare "0" (no unit) when nothing was sold, and "—" when the KPI is absent.
  const fmtDataSold = (bytes: number | undefined): string =>
    bytes == null ? '—' : bytes === 0 ? '0' : formatBytes(bytes);

  const bizMarkers = sites
    .filter((s) => s.lat !== 0 || s.lng !== 0)
    .map((s) => ({
      id: s.id,
      lat: s.lat,
      lng: s.lng,
      color: pinColor(s.status),
      popup: <div style={{ fontWeight: 600 }}>{s.name}</div>,
    }));

  return (
    <div className="page">
      <PageHeader
        title="Home"
        sub="Revenue, customers and sites at a glance."
        actions={<DateChip value={range} onChange={setRange} />}
      />
      {loading ? (
        <Skeleton variant="rounded" sx={{ height: 96 }} />
      ) : (
        <KpiRow
          items={[
            {
              icon: 'monetization_on',
              color: 'var(--uk-beige)',
              label: 'Revenue',
              value: kpiAmount(kpis, KPI_KEYS.revenue, money),
              delta:
                revDelta != null
                  ? `${revDelta >= 0 ? '+' : ''}${Math.round(revDelta)}% ${TREND_SUFFIX[span]}`
                  : undefined,
              dir: revDelta != null && revDelta < 0 ? 'down' : 'up',
            },
            {
              icon: 'group',
              color: 'var(--uk-secondary)',
              label: 'Customers',
              value: kpiText(kpis, KPI_KEYS.customers),
              delta:
                custDelta != null
                  ? `${custDelta >= 0 ? '+' : ''}${Math.round(custDelta)} ${PERIOD_LABEL[span]}`
                  : undefined,
              dir: custDelta != null && custDelta < 0 ? 'down' : 'up',
            },
            {
              icon: 'data_usage',
              color: 'var(--uk-ac)',
              label: 'Data sold',
              value: fmtDataSold(soldBytes),
              sub:
                span !== 'last_30d' && soldMonthBytes != null
                  ? `${fmtDataSold(soldMonthBytes)} last 30 days`
                  : undefined,
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-success-bright)',
              label: 'Network uptime',
              value: uptime != null ? `${uptime.toFixed(1)}%` : '—',
              sub: sitesTile.total
                ? `${sitesTile.text} sites online`
                : undefined,
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
        {sitesLoading ? (
          <Skeleton variant="rounded" sx={{ flex: 1, minHeight: 380, mt: 1 }} />
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
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
              }}
            >
              <div className="sec-title">Sites</div>
              <Button
                variant="text"
                startIcon={<ListAltRounded />}
                onClick={() => setShowSummary(true)}
              >
                View summary
              </Button>
            </div>
            <div style={{ flex: 1, minHeight: 300 }}>
              <UkamaMap
                markers={bizMarkers}
                onSelect={goSite}
                zoom={HOME_MAP_ZOOM}
                height="100%"
              />
            </div>
          </div>
        )}
      </div>

      {showSummary && (
        <AppModal
          title="Site summary"
          width={520}
          onClose={() => setShowSummary(false)}
          footer={
            <Button color="inherit" onClick={() => setShowSummary(false)}>
              Close
            </Button>
          }
        >
          <div
            style={{
              fontSize: 12.5,
              color: 'var(--uk-ink-3)',
              marginBottom: 4,
            }}
          >
            {sitesTile.online == null
              ? `${sites.length} sites`
              : `${sitesTile.online} of ${sitesTile.total} sites online`}
          </div>
          <SiteSummaryList
            sites={sites}
            onSite={(s) => {
              setShowSummary(false);
              goSite(s.id);
            }}
          />
        </AppModal>
      )}
    </div>
  );
}
