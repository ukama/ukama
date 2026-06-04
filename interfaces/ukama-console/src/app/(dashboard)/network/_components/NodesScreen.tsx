/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Nodes — radio hardware card grid, wired to the `nodesView` composite. */
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';
import PlaceRounded from '@mui/icons-material/PlaceRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import Skeleton from '@mui/material/Skeleton';

import { useNodesListQuery } from '@/client/graphql/nodes-list.generated';
import { EmptyState } from '@/components/EmptyState';
import FilterChips from '@/components/FilterChips';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import type { UkamaNode } from '@/data';
import { useUiPrefs } from '@/lib/store';
import { toUkamaNode } from '@/lib/mappers/nodes';
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
  const networkId = useUiPrefs((s) => s.networkId);
  const [filter, setFilter] = useState('all');
  const [drawerNode, setDrawerNode] = useState<UkamaNode | null>(null);

  const { data, loading, refetch } = useNodesListQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const nodesSection = data?.nodesView.nodes;
  const nodes: UkamaNode[] = useMemo(
    () => (nodesSection?.nodes ?? []).map((n) => toUkamaNode(n)),
    [nodesSection?.nodes]
  );

  const counts = {
    all: nodes.length,
    online: nodes.filter((n) => n.status === 'online').length,
    degraded: nodes.filter((n) => n.status === 'degraded' || n.status === 'configuring').length,
    offline: nodes.filter((n) => n.status === 'offline').length,
  };
  const list = nodes.filter(
    (n) =>
      filter === 'all' ||
      n.status === filter ||
      (filter === 'degraded' && n.status === 'configuring'),
  );

  return (
    <div className="page">
      <PageHeader title="Nodes" count={nodes.length} sub="Radio hardware deployed across your sites." />
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
      {loading ? (
        <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
          {[0, 1, 2].map((i) => (
            <Skeleton key={i} variant="rounded" sx={{ height: 132 }} />
          ))}
        </div>
      ) : nodesSection?.error ? (
        <EmptyState
          art="error"
          title="Couldn't load nodes"
          sub={nodesSection.error.message}
          cta="Try again"
          onCta={() => refetch()}
        />
      ) : list.length === 0 ? (
        <EmptyState art="node" title="No nodes" sub="No nodes match this filter." />
      ) : (
        <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
          {list.map((n) => (
            <NodeCard key={n.id} n={n} onOpen={(node) => setDrawerNode(node)} />
          ))}
        </div>
      )}
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
