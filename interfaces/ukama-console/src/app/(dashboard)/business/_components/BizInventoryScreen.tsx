/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';

/** Inventory — SIMs / Nodes / Hardware pill tabs (biz-ops.jsx BizInventory). */
import { useState } from 'react';
import MemoryRounded from '@mui/icons-material/MemoryRounded';
import DateChip from '@/components/DateChip';
import FilterChips from '@/components/FilterChips';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { BIZ_INVENTORY } from '@/data';

export default function BizInventoryScreen() {
  const b = BIZ_INVENTORY;
  const [tab, setTab] = useState('SIMs');

  return (
    <div className="page">
      <PageHeader
        title="Inventory"
        sub="Do I have enough SIMs and nodes to operate and grow?"
        actions={<DateChip />}
      />
      <KpiRow items={b.kpis} />
      <div style={{ marginBottom: 18 }}>
        <FilterChips
          options={[
            { value: 'SIMs', label: 'SIMs' },
            { value: 'Nodes', label: 'Nodes' },
            { value: 'Hardware', label: 'Hardware' },
          ]}
          value={tab}
          onChange={setTab}
        />
      </div>

      <div className="card card-pad">
        {tab === 'SIMs' && (
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>SIM / ICCID</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Assigned customer</TableCell>
                <TableCell>Site / network</TableCell>
                <TableCell>Activation date</TableCell>
                <TableCell>Issue</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {b.sims.map((r) => (
                <TableRow key={r.iccid}>
                  <TableCell className="tnum" style={{ fontWeight: 600 }}>
                    {r.iccid}
                  </TableCell>
                  <TableCell>
                    <StatusBadge status={r.status} variant="pill" />
                  </TableCell>
                  <TableCell className="tnum">{r.cust}</TableCell>
                  <TableCell className="muted">{r.site}</TableCell>
                  <TableCell className="muted">{r.date}</TableCell>
                  <TableCell
                    style={{
                      color: r.issue !== '—' ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-ink-3)',
                    }}
                  >
                    {r.issue}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
        {tab === 'Nodes' && (
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Serial</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Site</TableCell>
                <TableCell>Date</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {b.nodes.map((r) => (
                <TableRow key={r.serial}>
                  <TableCell className="tnum" style={{ fontWeight: 600 }}>
                    {r.serial}
                  </TableCell>
                  <TableCell className="muted">{r.type}</TableCell>
                  <TableCell>
                    <StatusBadge status={r.status} variant="pill" />
                  </TableCell>
                  <TableCell>{r.site}</TableCell>
                  <TableCell className="muted">{r.date}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
        {tab === 'Hardware' && (
          <div style={{ textAlign: 'center', padding: 48, color: 'var(--uk-ink-3)' }}>
            <MemoryRounded sx={{ fontSize: 42 }} />
            <div
              style={{
                fontFamily: 'var(--font-display)',
                fontSize: 18,
                fontWeight: 500,
                marginTop: 12,
                color: 'var(--uk-ink)',
              }}
            >
              Hardware
            </div>
            <div style={{ fontSize: 13.5, marginTop: 6 }}>
              Routers, amplifiers and accessories in inventory.
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
