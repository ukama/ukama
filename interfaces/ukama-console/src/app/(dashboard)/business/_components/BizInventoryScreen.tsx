/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Inventory — SIMs / Nodes / Hardware pill tabs, wired to `inventoryView`. */
import { useState } from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';

import { useInventoryOverviewQuery } from '@/client/graphql/team.generated';
import DateChip from '@/components/DateChip';
import { EmptyState } from '@/components/EmptyState';
import FilterChips from '@/components/FilterChips';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import StatusBadge from '@/components/StatusBadge';
import { toUkamaNode } from '@/lib/mappers/nodes';

export default function BizInventoryScreen() {
  const [tab, setTab] = useState('SIMs');
  const { data, loading, refetch } = useInventoryOverviewQuery();
  const view = data?.inventoryView;
  const simStock = view?.simStock;
  const componentsSection = view?.components;
  const nodesSection = view?.unassignedNodes;
  const nodes = (nodesSection?.nodes ?? []).map((n) => toUkamaNode(n));
  const categories = componentsSection?.byCategory ?? [];

  return (
    <div className="page">
      <PageHeader
        title="Inventory"
        sub="Do I have enough SIMs and nodes to operate and grow?"
        actions={<DateChip />}
      />
      <KpiRow
        items={[
          {
            icon: 'sim_card',
            color: 'var(--uk-ac)',
            label: 'SIMs available',
            value: sectionValue(simStock?.available, simStock?.error),
            sub: simStock?.lowStock ? 'low stock' : undefined,
            danger: !!simStock?.lowStock,
          },
          {
            icon: 'donut_small',
            color: 'var(--uk-secondary)',
            label: 'SIMs assigned',
            value: sectionValue(simStock?.consumed, simStock?.error),
            sub:
              simStock?.pctAssigned != null
                ? `${simStock.pctAssigned}% of pool`
                : undefined,
          },
          {
            icon: 'cell_tower',
            color: 'var(--uk-success-bright)',
            label: 'Nodes ready to install',
            value: sectionValue(nodes.length || null, nodesSection?.error),
          },
          {
            icon: 'memory',
            color: 'var(--uk-beige)',
            label: 'Components in stock',
            value: sectionValue(componentsSection?.total, componentsSection?.error),
          },
        ]}
      />
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
        {loading ? (
          <SkeletonTable cols={4} rows={4} />
        ) : (
          <>
            {tab === 'SIMs' &&
              (simStock?.error ? (
                <EmptyState
                  art="error"
                  title="Couldn't load SIM stock"
                  sub={simStock.error.message}
                  cta="Try again"
                  onCta={() => refetch()}
                />
              ) : (
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Metric</TableCell>
                      <TableCell align="right">Count</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(
                      [
                        ['Total in pool', simStock?.total],
                        ['Available', simStock?.available],
                        ['Assigned', simStock?.consumed],
                      ] as const
                    ).map(([label, value]) => (
                      <TableRow key={label}>
                        <TableCell style={{ fontWeight: 600 }}>{label}</TableCell>
                        <TableCell align="right" className="tnum">
                          {value ?? '—'}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ))}
            {tab === 'Nodes' &&
              (nodesSection?.error ? (
                <EmptyState
                  art="error"
                  title="Couldn't load nodes"
                  sub={nodesSection.error.message}
                  cta="Try again"
                  onCta={() => refetch()}
                />
              ) : nodes.length === 0 ? (
                <EmptyState
                  art="node"
                  title="No unassigned nodes"
                  sub="All registered nodes are deployed to sites."
                />
              ) : (
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Serial</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Status</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {nodes.map((n) => (
                      <TableRow key={n.id}>
                        <TableCell className="tnum" style={{ fontWeight: 600 }}>
                          {n.serial}
                        </TableCell>
                        <TableCell className="muted">{n.type}</TableCell>
                        <TableCell>
                          <StatusBadge status="available" variant="pill">
                            Available
                          </StatusBadge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ))}
            {tab === 'Hardware' &&
              (componentsSection?.error ? (
                <EmptyState
                  art="error"
                  title="Couldn't load components"
                  sub={componentsSection.error.message}
                  cta="Try again"
                  onCta={() => refetch()}
                />
              ) : categories.length === 0 ? (
                <EmptyState
                  art="search"
                  title="No components"
                  sub="Components registered to your org appear here."
                />
              ) : (
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Category</TableCell>
                      <TableCell align="right">Count</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {categories.map((c) => (
                      <TableRow key={c.category}>
                        <TableCell style={{ fontWeight: 600 }}>{c.category}</TableCell>
                        <TableCell align="right" className="tnum">
                          {c.count}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ))}
          </>
        )}
      </div>
    </div>
  );
}
