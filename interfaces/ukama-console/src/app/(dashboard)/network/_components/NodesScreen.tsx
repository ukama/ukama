/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Nodes — radio hardware card grid (screens-console.jsx). */
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import PlaceRounded from '@mui/icons-material/PlaceRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import FilterChips from '@/components/FilterChips';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { NODES } from '@/data';
import type { UkamaNode } from '@/data';
import NodeDrawer from './NodeDrawer';

function NodeCard({ n, onOpen }: { n: UkamaNode; onOpen: (n: UkamaNode) => void }) {
  const off = n.status === 'offline';
  const Icon = n.type.startsWith('Amp') ? SettingsInputAntennaRounded : RouterRounded;
  return (
    <div
      className="card ecard"
      role="button"
      tabIndex={0}
      onClick={() => onOpen(n)}
      onKeyDown={(e) => {
        if (e.key === 'Enter') onOpen(n);
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div style={{ display: 'flex', gap: 12 }}>
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
            <Icon sx={{ fontSize: 22, color: 'var(--uk-ac)' }} />
          </div>
          <div>
            <div style={{ fontSize: 14, fontWeight: 600 }}>{n.type}</div>
            <div className="tnum" style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
              {n.serial}
            </div>
          </div>
        </div>
        <StatusBadge status={n.status} />
      </div>
      <hr className="divider" style={{ margin: '14px 0' }} />
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          fontSize: 12.5,
          color: 'var(--uk-ink-2)',
        }}
      >
        <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
          <PlaceRounded sx={{ fontSize: 14 }} /> {n.site}
        </span>
        <span className="tnum" style={{ color: 'var(--uk-ink-3)' }}>
          {off ? 'No telemetry' : 'Up ' + n.up}
        </span>
      </div>
    </div>
  );
}

export default function NodesScreen() {
  const router = useRouter();
  const [filter, setFilter] = useState('all');
  const [drawerNode, setDrawerNode] = useState<UkamaNode | null>(null);
  const counts = {
    all: NODES.length,
    online: NODES.filter((n) => n.status === 'online').length,
    degraded: NODES.filter((n) => n.status === 'degraded' || n.status === 'configuring').length,
    offline: NODES.filter((n) => n.status === 'offline').length,
  };
  const list = NODES.filter(
    (n) =>
      filter === 'all' ||
      n.status === filter ||
      (filter === 'degraded' && n.status === 'configuring'),
  );

  return (
    <div className="page">
      <PageHeader title="Nodes" count={NODES.length} sub="Radio hardware deployed across your sites." />
      <div style={{ display: 'flex', gap: 8, marginBottom: 18, flexWrap: 'wrap' }}>
        <FilterChips
          value={filter}
          onChange={setFilter}
          options={[
            { value: 'all', label: 'All', count: counts.all },
            { value: 'online', label: 'Online', count: counts.online },
            { value: 'degraded', label: 'Needs attention', count: counts.degraded },
            { value: 'offline', label: 'Offline', count: counts.offline },
          ]}
        />
      </div>
      <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
        {list.map((n) => (
          <NodeCard key={n.id} n={n} onOpen={(node) => setDrawerNode(node)} />
        ))}
      </div>
      {drawerNode && (
        <NodeDrawer
          node={drawerNode}
          onClose={() => setDrawerNode(null)}
          onOpenDetail={(n) => {
            setDrawerNode(null);
            router.push(`/network/nodes/${n.id}`);
          }}
        />
      )}
    </div>
  );
}
