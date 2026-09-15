/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Node pool — every registered node, wired to the `nodesView` composite
 *  (NodePool operation). Status uses the same connectivity dot and state chip
 *  as the Nodes page, so a node reads the same in both places. */
import { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useApolloClient } from '@apollo/client';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import Tooltip from '@mui/material/Tooltip';
import Button from '@mui/material/Button';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import { useNodePoolQuery } from '@/client/graphql/nodes-list.generated';
import {
  SitesListDocument,
  type SitesListQuery,
  type SitesListQueryVariables,
} from '@/client/graphql/sites-list.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import { heldQuery } from '@/lib/heldQuery';
import { toUkamaNode } from '@/lib/mappers/nodes';
import {
  ConnectivityDot,
  StateChip,
  connLabel,
  describeNode,
} from './nodeStatus';

interface PoolRow {
  id: string;
  serial: string;
  type: string;
  assigned: boolean;
  site: string;
  connectivity: string;
  state: string;
}

const isOnline = (n: PoolRow) => n.connectivity.toLowerCase() === 'online';
const isState = (n: PoolRow, s: string) => n.state.toLowerCase() === s;

/** Sortable columns and the value each row sorts by. */
type SortKey = 'type' | 'status' | 'site';
const sortValue = (n: PoolRow, by: SortKey): string => {
  if (by === 'type') return n.type;
  if (by === 'site') return n.site;
  return `${connLabel(n.connectivity)} ${n.state}`;
};

/**
 * Row action, driven by site membership: a node without a site is installed
 * through the configure flow, a node with one opens its detail page.
 * Lifecycle state does not decide the label: a node released back into the
 * pool reports Ready, not Unknown, and must offer Configure again.
 * Configuring needs the node reachable, so it is disabled while offline.
 */
function RowAction({ item }: { item: PoolRow }) {
  const router = useRouter();
  const needsConfigure = !item.assigned;
  const blocked = needsConfigure && !isOnline(item);

  const button = (
    <Button
      variant="text"
      size="small"
      disabled={blocked}
      endIcon={<ChevronRightRounded />}
      onClick={() =>
        needsConfigure
          ? router.push('/configure/select-network')
          : router.push(`/network/nodes/${item.id}`)
      }
      sx={{
        fontSize: 13.5,
        fontWeight: 600,
        textTransform: 'none',
        whiteSpace: 'nowrap',
        color: needsConfigure ? 'var(--uk-ac)' : 'var(--uk-ink-2)',
        '& .MuiButton-endIcon': { ml: 0.25 },
      }}
    >
      {needsConfigure ? 'Configure' : 'View detail'}
    </Button>
  );

  if (!blocked) return button;
  return (
    <Tooltip title="Power the node on to configure it">
      <span>{button}</span>
    </Tooltip>
  );
}

export default function NodePoolScreen() {
  const poolResult = useNodePoolQuery();
  const refetch = poolResult.refetch;
  const { data, loading } = heldQuery(poolResult);
  const nodesSection = data?.nodesView.nodes;

  // NodePool spans every network, so resolve site names per network the
  // nodes actually belong to rather than only the selected one.
  const client = useApolloClient();
  const [siteNameById, setSiteNameById] = useState<Map<string, string>>(
    () => new Map(),
  );
  const networkIds = useMemo(() => {
    const ids = new Set<string>();
    for (const n of nodesSection?.nodes ?? []) {
      if (n.site?.networkId) ids.add(n.site.networkId);
    }
    return [...ids].sort();
  }, [nodesSection?.nodes]);
  const networkKey = networkIds.join(',');
  useEffect(() => {
    if (networkIds.length === 0) return;
    let cancelled = false;
    Promise.all(
      networkIds.map((networkId) =>
        client
          .query<SitesListQuery, SitesListQueryVariables>({
            query: SitesListDocument,
            variables: { networkId },
          })
          .then((r) => r.data?.sitesView.sites.sites ?? [])
          .catch(() => []),
      ),
    ).then((lists) => {
      if (cancelled) return;
      const map = new Map<string, string>();
      for (const s of lists.flat()) map.set(s.id, s.name);
      setSiteNameById(map);
    });
    return () => {
      cancelled = true;
    };
    // networkKey stands in for networkIds so the effect runs once per set.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [client, networkKey]);

  const pool: PoolRow[] = useMemo(
    () =>
      (nodesSection?.nodes ?? []).map((n) => {
        const mapped = toUkamaNode(n);
        const siteId = n.site?.siteId ?? '';
        return {
          id: n.id,
          serial: mapped.serial,
          type: mapped.type,
          assigned: Boolean(siteId),
          site: siteId ? (siteNameById.get(siteId) ?? siteId) : '—',
          connectivity: n.status.connectivity,
          state: n.status.state,
        };
      }),
    [nodesSection?.nodes, siteNameById],
  );

  const readyToInstall = pool.filter(
    (n) => !n.assigned && isOnline(n) && isState(n, 'ready'),
  ).length;
  const operational = pool.filter((n) => isState(n, 'operational')).length;
  const offline = pool.filter((n) => !isOnline(n)).length;

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
        sub="Every registered node, installed or not. A node that is online and Ready can be installed at a site."
      />
      <KpiRow
        items={[
          {
            icon: 'info',
            label: 'Ready to install',
            value: readyToInstall,
            color: 'var(--uk-ac)',
          },
          {
            icon: 'cell_tower',
            label: 'Operational',
            value: operational,
            color: 'var(--uk-success-bright)',
          },
          {
            icon: 'warning',
            label: 'Offline',
            value: offline,
            color: 'var(--uk-error)',
          },
          { icon: 'account_tree', label: 'In inventory', value: pool.length },
        ]}
      />
      <div className="card card-pad">
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={5} rows={5} />
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
                      ['status', 'Status'],
                      ['site', 'Site'],
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
                  <TableCell align="right" sx={{ width: 130 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {sortedPool.map((n) => (
                  <TableRow key={n.id}>
                    <TableCell className="tnum" style={{ fontWeight: 600 }}>
                      {n.serial}
                    </TableCell>
                    <TableCell>{n.type}</TableCell>
                    <TableCell>
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 8,
                        }}
                      >
                        <ConnectivityDot connectivity={n.connectivity} />
                        <span>{connLabel(n.connectivity)}</span>
                        <StateChip state={n.state} />
                      </div>
                      <div className="muted" style={{ fontSize: 12, marginTop: 2 }}>
                        {describeNode(n.connectivity, n.state, n.assigned)}
                      </div>
                    </TableCell>
                    <TableCell className="muted">{n.site}</TableCell>
                    <TableCell align="right">
                      <RowAction item={n} />
                    </TableCell>
                  </TableRow>
                ))}
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
