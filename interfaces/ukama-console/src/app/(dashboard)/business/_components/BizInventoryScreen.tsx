/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

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
          <table className="tbl">
            <thead>
              <tr className="static">
                <th>SIM / ICCID</th>
                <th>Status</th>
                <th>Assigned customer</th>
                <th>Site / network</th>
                <th>Activation date</th>
                <th>Issue</th>
              </tr>
            </thead>
            <tbody>
              {b.sims.map((r) => (
                <tr key={r.iccid} className="static">
                  <td className="tnum" style={{ fontWeight: 600 }}>
                    {r.iccid}
                  </td>
                  <td>
                    <StatusBadge status={r.status} variant="pill" />
                  </td>
                  <td className="tnum">{r.cust}</td>
                  <td className="muted">{r.site}</td>
                  <td className="muted">{r.date}</td>
                  <td
                    style={{
                      color: r.issue !== '—' ? 'var(--uk-error-deep, #cf121b)' : 'var(--uk-ink-3)',
                    }}
                  >
                    {r.issue}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {tab === 'Nodes' && (
          <table className="tbl">
            <thead>
              <tr className="static">
                <th>Serial</th>
                <th>Type</th>
                <th>Status</th>
                <th>Site</th>
                <th>Date</th>
              </tr>
            </thead>
            <tbody>
              {b.nodes.map((r) => (
                <tr key={r.serial} className="static">
                  <td className="tnum" style={{ fontWeight: 600 }}>
                    {r.serial}
                  </td>
                  <td className="muted">{r.type}</td>
                  <td>
                    <StatusBadge status={r.status} variant="pill" />
                  </td>
                  <td>{r.site}</td>
                  <td className="muted">{r.date}</td>
                </tr>
              ))}
            </tbody>
          </table>
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
