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
import { useMemo } from 'react';
import PlaceRounded from '@mui/icons-material/PlaceRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import Skeleton from '@mui/material/Skeleton';

import { useNodesListQuery } from '@/client/graphql/nodes-list.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { EmptyState } from '@/components/EmptyState';
import PageHeader from '@/components/PageHeader';
import type { UkamaNode } from '@/data';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import { useUiPrefs } from '@/lib/store';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { ConnectivityDot, StateChip } from './nodeStatus';

function NodeCard({ n, onOpen }: { n: UkamaNode; onOpen: (n: UkamaNode) => void }) {
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
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 10 }}>
        <div style={{ display: 'flex', gap: 12, flex: 1, minWidth: 0 }}>
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
          <div style={{ minWidth: 0 }}>
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 7,
                fontSize: 14,
                fontWeight: 600,
              }}
            >
              <ConnectivityDot connectivity={n.connectivity} />
              <span
                style={{
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                }}
              >
                {n.name ?? n.type}
              </span>
            </div>
            <div
              className="tnum"
              style={{
                fontSize: 12,
                color: 'var(--uk-ink-3)',
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
              }}
            >
              {n.type} · {n.serial}
            </div>
          </div>
        </div>
        <span style={{ flex: 'none' }}>
          <StateChip state={n.state} />
        </span>
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
      </div>
    </div>
  );
}

export default function NodesScreen() {
  const router = useRouter();
  const networkId = useUiPrefs((s) => s.networkId);

  const { data, loading, refetch } = useNodesListQuery({
    variables: { networkId },
    skip: !networkId,
    ...visiblePoll(POLL_OVERVIEW_MS),
  });

  // Resolve each node's siteId → site name for the location line.
  const { data: sitesData } = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const siteNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const s of sitesData?.sitesView.sites.sites ?? []) m.set(s.id, s.name);
    return m;
  }, [sitesData]);

  const nodesSection = data?.nodesView.nodes;
  const nodes: UkamaNode[] = useMemo(
    () =>
      (nodesSection?.nodes ?? []).map((n) =>
        toUkamaNode(
          n,
          n.site?.siteId ? siteNameById.get(n.site.siteId) : undefined,
        ),
      ),
    [nodesSection?.nodes, siteNameById]
  );

  return (
    <div className="page">
      <PageHeader title="Nodes" count={nodes.length} sub="Radio hardware deployed across your sites." />
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
      ) : nodes.length === 0 ? (
        <EmptyState art="node" title="No nodes" sub="Registered nodes appear here." />
      ) : (
        <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
          {nodes.map((n) => (
            <NodeCard
              key={n.id}
              n={n}
              onOpen={(node) => router.push(`/network/nodes/${node.id}`)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
