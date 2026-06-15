/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Node pool — hardware inventory, wired to the `nodesView` composite
 *  (NodePool operation: nodes without a site are available to install). */
import { useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import Button from '@mui/material/Button';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import { useNodePoolQuery } from '@/client/graphql/nodes-list.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useUiPrefs } from '@/lib/store';
import { toUkamaNode } from '@/lib/mappers/nodes';

type PoolStatus = 'available' | 'assigned';

interface PoolRow {
  id: string;
  serial: string;
  type: string;
  status: PoolStatus;
  site: string;
  connectivity: string;
  state: string;
}

const NP_LABEL: Record<PoolStatus, string> = {
  available: 'Available',
  assigned: 'Assigned',
};

/** Sortable columns and the value each row sorts by. */
type SortKey = 'type' | 'connectivity' | 'status';
const sortValue = (n: PoolRow, by: SortKey): string => {
  if (by === 'type') return n.type;
  if (by === 'connectivity') return connectivity(n.connectivity).label;
  return NP_LABEL[n.status];
};

/** Maps a node's raw connectivity to a status badge + label. */
function connectivity(raw: string): { kind: string; label: string } {
  const c = raw.toLowerCase();
  if (c === 'online') return { kind: 'online', label: 'Online' };
  if (c === 'offline') return { kind: 'offline', label: 'Offline' };
  return { kind: 'configuring', label: 'Unknown' };
}

/**
 * Row action, driven by the node's lifecycle state — the label — and its
 * connectivity — whether the action is live:
 *  - Unknown state → not yet set up → "Configure" (routes to the flow).
 *  - Configured / Faulty state → "View detail" (routes to the node page).
 * Either action requires the node to be reachable, so the button is disabled
 * unless the node is Online. (An unconfigured node reports Unknown connectivity
 * until it first checks in, so it shows a disabled Configure until it's live.)
 */
function RowAction({ item }: { item: PoolRow }) {
  const router = useRouter();
  const isOnline = item.connectivity.toLowerCase() === 'online';
  const needsConfigure = item.state.toLowerCase() === 'unknown';

  const label = needsConfigure ? 'Configure' : 'View detail';
  const onClick = () =>
    needsConfigure
      ? router.push('/configure/select-network')
      : router.push(`/network/nodes/${item.id}`);

  return (
    <Button
      variant="text"
      size="small"
      disabled={!isOnline}
      endIcon={<ChevronRightRounded />}
      onClick={onClick}
      sx={{
        fontSize: 13.5,
        fontWeight: 600,
        textTransform: 'none',
        whiteSpace: 'nowrap',
        color: needsConfigure ? 'var(--uk-ac)' : 'var(--uk-ink-2)',
        '& .MuiButton-endIcon': { ml: 0.25 },
      }}
    >
      {label}
    </Button>
  );
}

export default function NodePoolScreen() {
  const { data, loading, refetch } = useNodePoolQuery();
  const nodesSection = data?.nodesView.nodes;

  // Resolve siteId → site name for the Site column. NodePool isn't scoped to
  // a network, so use the currently-selected network for the site lookup.
  const networkId = useUiPrefs((s) => s.networkId);
  const { data: sitesData } = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const siteNameById = useMemo(() => {
    const map = new Map<string, string>();
    for (const s of sitesData?.sitesView.sites.sites ?? [])
      map.set(s.id, s.name);
    return map;
  }, [sitesData]);

  // Pool view-model: a node without a site is available to install.
  const pool: PoolRow[] = useMemo(
    () =>
      (nodesSection?.nodes ?? []).map((n) => {
        const mapped = toUkamaNode(n);
        const siteId = n.site?.siteId ?? '';
        return {
          id: n.id,
          serial: mapped.serial,
          type: mapped.type,
          status: siteId ? ('assigned' as const) : ('available' as const),
          site: siteId ? (siteNameById.get(siteId) ?? siteId) : '—',
          connectivity: n.status.connectivity,
          state: n.status.state,
        };
      }),
    [nodesSection?.nodes, siteNameById],
  );

  const avail = pool.filter((n) => n.status === 'available').length;
  const deployed = pool.length - avail;

  // Click-to-sort on Type / Connectivity / Status (toggles asc → desc → off).
  const [sort, setSort] = useState<{ by: SortKey; dir: 'asc' | 'desc' } | null>(
    null,
  );
  const toggleSort = (by: SortKey) =>
    setSort((cur) =>
      cur?.by !== by
        ? { by, dir: 'asc' }
        : cur.dir === 'asc'
          ? { by, dir: 'desc' }
          : null,
    );
  const sortedPool = useMemo(() => {
    if (!sort) return pool;
    return [...pool].sort((a, b) => {
      const r = sortValue(a, sort.by).localeCompare(sortValue(b, sort.by));
      return sort.dir === 'asc' ? r : -r;
    });
  }, [pool, sort]);

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'Node pool']}
        title="Node pool"
        count={pool.length}
        sub="Hardware in inventory, ready to install at a site."
      />
      <KpiRow
        items={[
          {
            icon: 'info',
            label: 'Available to install',
            value: avail,
            color: 'var(--uk-success-bright)',
          },
          {
            icon: 'cell_tower',
            label: 'Deployed (live)',
            value: deployed,
            color: 'var(--uk-ac)',
          },
          { icon: 'account_tree', label: 'In inventory', value: pool.length },
        ]}
      />
      <div className="card card-pad">
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={6} rows={5} />
          ) : nodesSection?.error ? (
            <EmptyState
              art="error"
              title="Couldn't load node pool"
              sub={nodesSection.error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : pool.length === 0 ? (
            <EmptyState
              art="node"
              title="No nodes in inventory"
              sub="Registered nodes appear here."
            />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Node ID</TableCell>
                  {(
                    [
                      ['type', 'Type'],
                      ['connectivity', 'Connectivity'],
                      ['status', 'Status'],
                    ] as [SortKey, string][]
                  ).map(([key, label]) => (
                    <TableCell
                      key={key}
                      sortDirection={sort?.by === key ? sort.dir : false}
                    >
                      <TableSortLabel
                        active={sort?.by === key}
                        direction={sort?.by === key ? sort.dir : 'asc'}
                        onClick={() => toggleSort(key)}
                      >
                        {label}
                      </TableSortLabel>
                    </TableCell>
                  ))}
                  <TableCell>Site</TableCell>
                  <TableCell align="right" sx={{ width: 130 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {sortedPool.map((n) => {
                  const conn = connectivity(n.connectivity);
                  return (
                    <TableRow key={n.id}>
                      <TableCell className="tnum" style={{ fontWeight: 600 }}>
                        {n.serial}
                      </TableCell>
                      <TableCell>{n.type}</TableCell>
                      <TableCell>
                        <StatusBadge status={conn.kind}>
                          {conn.label}
                        </StatusBadge>
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={n.status}>
                          {NP_LABEL[n.status]}
                        </StatusBadge>
                      </TableCell>
                      <TableCell className="muted">{n.site}</TableCell>
                      <TableCell align="right">
                        <RowAction item={n} />
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && !nodesSection?.error && (
          <TableFooter count={pool.length} noun="nodes" />
        )}
      </div>
    </div>
  );
}
