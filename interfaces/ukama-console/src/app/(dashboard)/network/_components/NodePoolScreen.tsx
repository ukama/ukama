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
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import AddLocationAltRounded from '@mui/icons-material/AddLocationAltRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import { useNodePoolQuery } from '@/client/graphql/nodes-list.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
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
}

const NP_LABEL: Record<PoolStatus, string> = {
  available: 'Available',
  assigned: 'Assigned',
};

/** Maps a node's raw connectivity to a status badge + label. */
function connectivity(raw: string): { kind: string; label: string } {
  const c = raw.toLowerCase();
  if (c === 'online') return { kind: 'online', label: 'Online' };
  if (c === 'offline') return { kind: 'offline', label: 'Offline' };
  return { kind: 'configuring', label: 'Unknown' };
}

function PoolMenu({ item }: { item: PoolRow }) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const router = useRouter();
  const toast = useToast();
  return (
    <>
      <IconButton
        size="small"
        aria-label="More actions"
        sx={{ color: 'var(--uk-ink-3)' }}
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        <MoreVertRounded sx={{ fontSize: 20 }} />
      </IconButton>
      <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
        {item.status === 'available' && (
          <MenuItem
            sx={{ fontSize: 13.5, gap: 1.25 }}
            onClick={() => {
              setAnchor(null);
              toast(`Assign ${item.serial} to a site`);
            }}
          >
            <AddLocationAltRounded sx={{ fontSize: 18 }} /> Assign to site
          </MenuItem>
        )}
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            router.push(`/network/nodes/${item.id}`);
          }}
        >
          <InfoRounded sx={{ fontSize: 18 }} /> Details
        </MenuItem>
      </Menu>
    </>
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
    for (const s of sitesData?.sitesView.sites.sites ?? []) map.set(s.id, s.name);
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
        };
      }),
    [nodesSection?.nodes, siteNameById]
  );

  const avail = pool.filter((n) => n.status === 'available').length;
  const deployed = pool.length - avail;

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
          { icon: 'cell_tower', label: 'Deployed (live)', value: deployed, color: 'var(--uk-ac)' },
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
            <EmptyState art="node" title="No nodes in inventory" sub="Registered nodes appear here." />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Node ID</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Connectivity</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Site</TableCell>
                  <TableCell sx={{ width: 44 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {pool.map((n) => {
                  const conn = connectivity(n.connectivity);
                  return (
                    <TableRow key={n.id}>
                      <TableCell className="tnum" style={{ fontWeight: 600 }}>
                        {n.serial}
                      </TableCell>
                      <TableCell>{n.type}</TableCell>
                      <TableCell>
                        <StatusBadge status={conn.kind}>{conn.label}</StatusBadge>
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={n.status}>{NP_LABEL[n.status]}</StatusBadge>
                      </TableCell>
                      <TableCell className="muted">{n.site}</TableCell>
                      <TableCell>
                        <PoolMenu item={n} />
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && !nodesSection?.error && <TableFooter count={pool.length} noun="nodes" />}
      </div>
    </div>
  );
}
