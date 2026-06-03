/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Node pool — hardware in inventory, ready to install (screens-manage.jsx). */
import { useState } from 'react';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import AddLocationAltRounded from '@mui/icons-material/AddLocationAltRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { NODE_POOL, NODES } from '@/data';
import type { NodePoolItem } from '@/data';
import { useFirstLoad } from '@/lib/useFirstLoad';

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
  const avail = NODE_POOL.filter((n) => n.status === 'available').length;
  const loading = useFirstLoad('nodepool');

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'Node pool']}
        title="Node pool"
        count={NODE_POOL.length}
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
          { icon: 'cell_tower', label: 'Deployed (live)', value: NODES.length, color: 'var(--uk-ac)' },
          { icon: 'account_tree', label: 'In inventory', value: NODE_POOL.length },
        ]}
      />
      <div className="card card-pad">
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={6} rows={5} />
          ) : (
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Serial</th>
                  <th>Type</th>
                  <th>Status</th>
                  <th>Assigned to</th>
                  <th>Added</th>
                  <th style={{ width: 44 }} />
                </tr>
              </thead>
              <tbody>
                {NODE_POOL.map((n) => (
                  <tr key={n.id} className="static">
                    <td className="tnum" style={{ fontWeight: 600 }}>
                      {n.serial}
                    </td>
                    <td>{n.type}</td>
                    <td>
                      <StatusBadge status={n.status}>{NP_LABEL[n.status]}</StatusBadge>
                    </td>
                    <td className="muted">{n.site ?? '—'}</td>
                    <td className="muted">{n.added}</td>
                    <td>
                      <PoolMenu item={n} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {!loading && <TableFooter count={NODE_POOL.length} noun="nodes" />}
      </div>
    </div>
  );
}
