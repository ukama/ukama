/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Sites — card grid with status count chips (screens-console.jsx). */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import CellTowerRounded from '@mui/icons-material/CellTowerRounded';
import ErrorOutlineRounded from '@mui/icons-material/ErrorOutlineRounded';
import { EmptyState } from '@/components/EmptyState';
import FilterChips from '@/components/FilterChips';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import StatusBadge from '@/components/StatusBadge';
import { SITES } from '@/data';
import type { Site } from '@/data';

function SiteCard({ s, onOpen }: { s: Site; onOpen: (s: Site) => void }) {
  const issueColor = s.status === 'offline' ? 'var(--uk-error-deep, #cf121b)' : '#b5591b';
  return (
    <div
      className="card ecard"
      role="button"
      tabIndex={0}
      onClick={() => onOpen(s)}
      onKeyDown={(e) => {
        if (e.key === 'Enter') onOpen(s);
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 10 }}>
        <div style={{ display: 'flex', gap: 12, minWidth: 0 }}>
          <div
            style={{
              width: 42,
              height: 42,
              borderRadius: 10,
              background: 'var(--uk-ac-soft)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flex: 'none',
            }}
          >
            <CellTowerRounded sx={{ fontSize: 22, color: 'var(--uk-ac)' }} />
          </div>
          <div style={{ minWidth: 0 }}>
            <div style={{ fontFamily: 'var(--font-display)', fontSize: 16, fontWeight: 500 }}>
              {s.name}
            </div>
            <div
              style={{
                fontSize: 12.5,
                color: 'var(--uk-ink-3)',
                marginTop: 1,
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
              }}
            >
              {s.area}
            </div>
          </div>
        </div>
        <StatusBadge status={s.status} />
      </div>
      {s.issue && (
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 7,
            marginTop: 13,
            fontSize: 12.5,
            fontWeight: 500,
            color: issueColor,
          }}
        >
          <ErrorOutlineRounded sx={{ fontSize: 16 }} />
          {s.issue}
        </div>
      )}
      <hr className="divider" style={{ margin: '14px 0' }} />
      <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 12.5, color: 'var(--uk-ink-2)' }}>
        <span>
          {s.subs} customers · {s.nodes} nodes
        </span>
        <span className="tnum">{s.uptime}% uptime</span>
      </div>
    </div>
  );
}

export default function SitesScreen() {
  const router = useRouter();
  const [filter, setFilter] = useState('all');
  const [q, setQ] = useState('');
  const counts = {
    all: SITES.length,
    online: SITES.filter((s) => s.status === 'online').length,
    degraded: SITES.filter((s) => s.status === 'degraded').length,
    offline: SITES.filter((s) => s.status === 'offline').length,
  };
  const list = SITES.filter(
    (s) =>
      (filter === 'all' || s.status === filter) &&
      s.name.toLowerCase().includes(q.toLowerCase()),
  );
  const open = (s: Site) => router.push(`/network/sites/${s.id}`);

  return (
    <div className="page">
      <PageHeader
        title="Sites"
        count={SITES.length}
        sub="Physical locations where your network is installed."
      />
      <div style={{ display: 'flex', gap: 10, marginBottom: 18, flexWrap: 'wrap', alignItems: 'center' }}>
        <SearchField value={q} onChange={setQ} placeholder="Search sites" width={260} />
        <FilterChips
          value={filter}
          onChange={setFilter}
          options={[
            { value: 'all', label: 'All', count: counts.all },
            { value: 'online', label: 'Online', count: counts.online },
            { value: 'degraded', label: 'Degraded', count: counts.degraded },
            { value: 'offline', label: 'Offline', count: counts.offline },
          ]}
        />
      </div>
      <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(310px, 1fr))' }}>
        {list.map((s) => (
          <SiteCard key={s.id} s={s} onOpen={open} />
        ))}
      </div>
      {list.length === 0 && (
        <div className="card">
          <EmptyState art="search" title="No matching sites" sub="Try a different filter or search term." />
        </div>
      )}
    </div>
  );
}
