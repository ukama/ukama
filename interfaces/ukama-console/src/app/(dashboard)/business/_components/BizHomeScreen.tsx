/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Business Home — KPIs + full-height sites map + site summary modal. */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import ListAltRounded from '@mui/icons-material/ListAltRounded';
import AppModal from '@/components/AppModal';
import DateChip from '@/components/DateChip';
import { KpiRow } from '@/components/Kpi';
import SiteMap, { StatusDot } from '@/components/Map/SiteMap';
import PageHeader from '@/components/PageHeader';
import { BIZ_HOME, BIZ_SITES } from '@/data';
import type { BizSite } from '@/data';

function SiteSummaryList({ onSite }: { onSite: (s: BizSite) => void }) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      {BIZ_SITES.map((s, i) => (
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
            borderBottom: i < BIZ_SITES.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
          }}
        >
          <span style={{ marginTop: 4, display: 'inline-flex' }}>
            <StatusDot status={s.status} />
          </span>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 13.5, fontWeight: 600 }}>{s.name}</div>
            <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 1 }}>
              {s.status === 'offline'
                ? `$0 today · ${s.customers} affected · Offline`
                : `$${s.revToday} today · ${s.custToday} customers · ${s.status === 'warning' ? 'Warning' : 'Online'}`}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default function BizHomeScreen() {
  const router = useRouter();
  const [showSummary, setShowSummary] = useState(false);
  const online = BIZ_SITES.filter((s) => s.status === 'online').length;
  const goSite = (id: string) => router.push(`/business/sites/${id}`);

  return (
    <div className="page">
      <PageHeader
        title="Home"
        sub="Your network served 142 customers and sold 87 GB today."
        actions={<DateChip />}
      />
      <KpiRow items={BIZ_HOME.kpis} />

      <div style={{ flex: 1, minHeight: 420, display: 'flex', flexDirection: 'column' }}>
        <SiteMap
          sites={BIZ_SITES}
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
            {online} of {BIZ_SITES.length} sites online
          </div>
          <SiteSummaryList
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
