/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Business Home — KPIs + full-height sites map, wired to `commerceView`
 * (revenue) and `networkOverview` (customers + sites). Per-site revenue is
 * backend gap #10 and renders as "—" until it lands.
 */
import { useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import ListAltRounded from '@mui/icons-material/ListAltRounded';

import {
  useBizHomeNetworkQuery,
  useBizHomeRevenueQuery,
} from '@/client/graphql/commerce.generated';
import AppModal from '@/components/AppModal';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import SiteMap, { StatusDot } from '@/components/Map/SiteMap';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import type { BizSite } from '@/data';
import { useUiPrefs } from '@/lib/store';

const money = (value?: number | null): string =>
  value == null ? '—' : `$${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`;

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
            borderBottom: i < sites.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
          }}
        >
          <span style={{ marginTop: 4, display: 'inline-flex' }}>
            <StatusDot status={s.status} />
          </span>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 13.5, fontWeight: 600 }}>{s.name}</div>
            <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 1 }}>
              {/* per-site revenue/customers: backend gap #10 */}
              {s.status === 'offline' ? 'Offline' : 'Online'} · revenue —
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

  const { data: revenueData, loading: revenueLoading } = useBizHomeRevenueQuery({
    variables: { networkId },
  });
  const { data: networkData, loading: networkLoading } = useBizHomeNetworkQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const revenue = revenueData?.commerceView.revenue;
  const subStats = networkData?.networkOverview.subscriberStats;
  const siteStats = networkData?.networkOverview.siteStats;
  const loading = revenueLoading || networkLoading;

  const sites: BizSite[] = useMemo(
    () =>
      (siteStats?.sites ?? []).map((s) => ({
        id: s.id,
        name: s.name,
        status: s.isDeactivated ? 'offline' : 'online',
        // TODO(backend-gap #10): per-site revenue/customer rollup
        revenue: 0,
        revToday: 0,
        customers: 0,
        custToday: 0,
        data: '—',
        uptime: 0,
        top: '—',
        issue: null,
        lat: parseFloat(s.latitude) || 0,
        lng: parseFloat(s.longitude) || 0,
      })),
    [siteStats?.sites]
  );
  const online = sites.filter((s) => s.status === 'online').length;
  const goSite = (id: string) => router.push(`/business/sites/${id}`);

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
              value: revenue?.error ? '—' : money(revenue?.monthPaid),
              sub:
                revenue?.momPct != null
                  ? `${revenue.momPct >= 0 ? '+' : ''}${revenue.momPct}% vs last month`
                  : undefined,
            },
            {
              icon: 'group',
              color: 'var(--uk-secondary)',
              label: 'Active customers',
              value: sectionValue(subStats?.active, subStats?.error),
              sub:
                subStats?.total != null ? `${subStats.total} total` : undefined,
            },
            {
              icon: 'donut_small',
              color: 'var(--uk-ac)',
              label: 'Collected to date',
              value: revenue?.error ? '—' : money(revenue?.totalPaid),
            },
            {
              icon: 'cell_tower',
              color: 'var(--uk-success-bright)',
              label: 'Sites online',
              value: siteStats?.error ? '—' : `${online}/${sites.length}`,
              danger: online < sites.length,
            },
          ]}
        />
      )}

      <div style={{ flex: 1, minHeight: 420, display: 'flex', flexDirection: 'column' }}>
        {networkLoading ? (
          <Skeleton variant="rounded" sx={{ flex: 1, minHeight: 380 }} />
        ) : (
          <SiteMap
            sites={sites}
            title="Sites"
            fill
            action={
              <Button
                variant="text"
                startIcon={<ListAltRounded />}
                onClick={() => setShowSummary(true)}
              >
                View summary
              </Button>
            }
            onSelect={(s) => goSite(s.id)}
          />
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
          <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)', marginBottom: 4 }}>
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
