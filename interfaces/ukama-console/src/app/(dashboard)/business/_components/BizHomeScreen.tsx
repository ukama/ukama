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
 * `getBusinessHome` (headline KPIs) and `getBusinessSites` (per-site revenue /
 * customers / coordinates). Per-site revenue (was backend gap #10) is now
 * served by the analytics rollup. KPI keys and any not-yet-emitted fields are
 * listed in docs/analytics-backend-gaps.md and degrade to "—".
 */
import ListAltRounded from '@mui/icons-material/ListAltRounded';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';

import {
  useGetBusinessHomeQuery,
  useGetBusinessSitesQuery,
} from '@/client/graphql/analytics.generated';
import AppModal from '@/components/AppModal';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import { StatusDot } from '@/components/Map/SiteMap';
import UkamaMap, { HOME_MAP_ZOOM } from '@/components/Map/UkamaMap';
import PageHeader from '@/components/PageHeader';
import type { BizSite } from '@/data';
import { kpiByKey, kpiText, kpiValue } from '@/lib/kpis';
import { useUiPrefs } from '@/lib/store';

// KPI keys this screen reads — see docs/analytics-backend-gaps.md.
const KEY = {
  revenueMonth: 'revenue_month',
  revenueCollected: 'revenue_collected',
  activeCustomers: 'active_customers',
  customersTotal: 'customers_total',
} as const;

const money = (value?: number | null): string =>
  value == null
    ? '—'
    : `$${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;

/** Map a backend status string onto the BizSite status union. */
const normStatus = (status?: string | null): BizSite['status'] => {
  const s = (status ?? '').toLowerCase();
  if (s === 'offline' || s === 'down' || s === 'deactivated') return 'offline';
  if (s === 'warning' || s === 'degraded') return 'warning';
  return 'online';
};

function SiteSummaryList({
  sites,
  onSite,
}: {
  sites: BizSite[];
  onSite: (s: BizSite) => void;
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
              {money(s.revenue)} · {s.customers} customers
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default function BizHomeScreen() {
  const router = useRouter();
  const networkId = useUiPrefs((s) => s.networkId);
  const [showSummary, setShowSummary] = useState(false);

  const { data: homeData, loading: homeLoading } = useGetBusinessHomeQuery({
    variables: { data: { networkId } },
  });
  const { data: sitesData, loading: sitesLoading } = useGetBusinessSitesQuery({
    variables: { data: { networkId } },
    skip: !networkId,
  });
  const kpis = homeData?.getBusinessHome.kpis;
  const monthDelta = kpiByKey(kpis, KEY.revenueMonth)?.delta;
  const totalCustomers = kpiValue(kpis, KEY.customersTotal);
  const loading = homeLoading || sitesLoading;

  const sites: BizSite[] = useMemo(
    () =>
      (sitesData?.getBusinessSites.sites ?? []).map((s) => ({
        id: s.siteId,
        name: s.name ?? s.siteId,
        status: normStatus(s.status),
        revenue: s.revenue,
        revToday: s.revenueToday,
        customers: s.customers,
        custToday: 0,
        data: '—',
        uptime: s.uptime,
        top: s.topPackage ?? '—',
        issue: s.issue ?? null,
        lat: s.latitude,
        lng: s.longitude,
      })),
    [sitesData?.getBusinessSites.sites],
  );
  const online = sites.filter((s) => s.status !== 'offline').length;
  const goSite = (id: string) => router.push(`/business/sites/${id}`);

  const BIZ_PIN: Record<string, string> = {
    online: 'var(--uk-success-bright)',
    warning: 'var(--uk-warning)',
    degraded: 'var(--uk-warning)',
    offline: 'var(--uk-error)',
  };
  const bizMarkers = sites
    .filter((s) => s.lat !== 0 || s.lng !== 0)
    .map((s) => ({
      id: s.id,
      lat: s.lat,
      lng: s.lng,
      color: BIZ_PIN[s.status] ?? 'var(--uk-ac)',
      popup: <div style={{ fontWeight: 600 }}>{s.name}</div>,
    }));

  return (
    <div className="page">
      <PageHeader
        title="Home"
        sub="Revenue, customers and sites at a glance."
        actions={<DateChip />}
      />
      {loading ? (
        <Skeleton variant="rounded" sx={{ height: 96 }} />
      ) : (
        <KpiRow
          items={[
            {
              icon: 'monetization_on',
              color: 'var(--uk-beige)',
              label: 'Revenue this month',
              value: kpiText(kpis, KEY.revenueMonth, money),
              sub:
                monthDelta != null
                  ? `${monthDelta >= 0 ? '+' : ''}${monthDelta}% vs last month`
                  : undefined,
            },
            {
              icon: 'group',
              color: 'var(--uk-secondary)',
              label: 'Active customers',
              value: kpiText(kpis, KEY.activeCustomers),
              sub:
                totalCustomers != null ? `${totalCustomers} total` : undefined,
            },
            {
              icon: 'donut_small',
              color: 'var(--uk-ac)',
              label: 'Collected to date',
              value: kpiText(kpis, KEY.revenueCollected, money),
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-success-bright)',
              label: 'Sites online',
              value: sites.length === 0 ? '—' : `${online}/${sites.length}`,
              danger: sites.length > 0 && online < sites.length,
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
            {online} of {sites.length} sites online
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
