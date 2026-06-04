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
import Divider from '@mui/material/Divider';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import AddLocationAltRounded from '@mui/icons-material/AddLocationAltRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import { useNodePoolQuery } from '@/client/graphql/nodes-list.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { NodePoolItem } from '@/data';
import { toUkamaNode } from '@/lib/mappers/nodes';

const NP_LABEL: Record<NodePoolItem['status'], string> = {
  available: 'Available',
  assigned: 'Assigned',
  rma: 'RMA',
};

function PoolMenu({ item }: { item: NodePoolItem }) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
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
            toast(`${item.serial} · ${item.type}`);
          }}
        >
          <InfoRounded sx={{ fontSize: 18 }} /> Details
        </MenuItem>
        <Divider />
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25, color: 'var(--uk-error)' }}
          onClick={() => {
            setAnchor(null);
            toast(`${item.serial} removed`, {
              action: { label: 'Undo', fn: () => toast(`${item.serial} restored`) },
            });
          }}
        >
          <DeleteOutlineRounded sx={{ fontSize: 18 }} /> Remove
        </MenuItem>
      </Menu>
    </>
  );
}

export default function NodePoolScreen() {
  const { data, loading, refetch } = useNodePoolQuery();
  const nodesSection = data?.nodesView.nodes;

  // Pool view-model: a node without a site is available to install.
  const pool: NodePoolItem[] = useMemo(
    () =>
      (nodesSection?.nodes ?? []).map((n) => {
        const mapped = toUkamaNode(n);
        const assigned = !!n.site?.siteId;
        return {
          id: n.id,
          serial: mapped.serial,
          type: mapped.type,
          status: assigned ? ('assigned' as const) : ('available' as const),
          site: n.site?.siteId ?? undefined,
          added: '—',
        };
      }),
    [nodesSection?.nodes]
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
                  <TableCell>Serial</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Assigned to</TableCell>
                  <TableCell>Added</TableCell>
                  <TableCell sx={{ width: 44 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {pool.map((n) => (
                  <TableRow key={n.id}>
                    <TableCell className="tnum" style={{ fontWeight: 600 }}>
                      {n.serial}
                    </TableCell>
                    <TableCell>{n.type}</TableCell>
                    <TableCell>
                      <StatusBadge status={n.status}>{NP_LABEL[n.status]}</StatusBadge>
                    </TableCell>
                    <TableCell className="muted">{n.site ?? '—'}</TableCell>
                    <TableCell className="muted">{n.added}</TableCell>
                    <TableCell>
                      <PoolMenu item={n} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && !nodesSection?.error && <TableFooter count={pool.length} noun="nodes" />}
      </div>
    </div>
  );
}
