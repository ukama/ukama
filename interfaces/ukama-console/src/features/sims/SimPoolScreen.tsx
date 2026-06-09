/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * SIM pool — inventory cockpit wired to the `simPoolView` composite. Stats
 * (incl. pctAssigned / lowStock) are derived server-side; the table lists
 * pool SIMs. Business can act; Network is view-only.
 */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import UploadFileRounded from '@mui/icons-material/UploadFileRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';

import { useSimPoolOverviewQuery } from '@/client/graphql/sim-pool.generated';
import { EmptyState } from '@/components/EmptyState';
import FilterChips from '@/components/FilterChips';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { useToast } from '@/components/ToastProvider';
import { formatDate, parseTimestamp } from '@/lib/parsers';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import { publicEnv } from '@/lib/runtime-env';
import UploadSimsDialog from './UploadSimsDialog';

type SimRow = { isAllocated: boolean; isFailed: boolean; isPhysical: boolean };
type TypeFilter = 'all' | 'physical' | 'esim';
type StatusFilter = 'all' | 'available' | 'assigned' | 'faulty';
type SortKey = 'type' | 'status' | 'added';
type SortDir = 'asc' | 'desc';

const statusKey = (s: SimRow): StatusFilter =>
  s.isFailed ? 'faulty' : s.isAllocated ? 'assigned' : 'available';

// Matches the BFF cap (MAX_POOL_SIMS) so the table can show the full pool.
const LIST_LIMIT = 100;

function SimMenu({ iccid }: { iccid: string }) {
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
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Viewing SIM ${iccid}`);
          }}
        >
          <VisibilityRounded sx={{ fontSize: 18 }} /> View SIM
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Exported ${iccid}`);
          }}
        >
          <DownloadRounded sx={{ fontSize: 18 }} /> Export
        </MenuItem>
      </Menu>
    </>
  );
}

const simStatus = (sim: { isAllocated: boolean; isFailed: boolean }) =>
  sim.isFailed ? 'offline' : sim.isAllocated ? 'online' : 'configuring';

const simStatusLabel = (sim: { isAllocated: boolean; isFailed: boolean }) =>
  sim.isFailed ? 'Faulty' : sim.isAllocated ? 'Assigned' : 'Available';

export default function SimPoolScreen({ canAct }: { canAct: boolean }) {
  const toast = useToast();
  const [showUpload, setShowUpload] = useState(false);
  const [typeFilter, setTypeFilter] = useState<TypeFilter>('all');
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');
  const [sortKey, setSortKey] = useState<SortKey>('added');
  const [sortDir, setSortDir] = useState<SortDir>('desc');

  const { data, loading, refetch } = useSimPoolOverviewQuery({
    variables: { simType: publicEnv().simType, limit: LIST_LIMIT },
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const stats = data?.simPoolView.stats;
  const simsSection = data?.simPoolView.sims;
  const allSims = useMemo(() => simsSection?.sims ?? [], [simsSection]);

  const toggleSort = (key: SortKey) => {
    if (sortKey === key) setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'));
    else {
      setSortKey(key);
      setSortDir('asc');
    }
  };

  const sims = useMemo(() => {
    const filtered = allSims.filter((s) => {
      if (typeFilter !== 'all') {
        if (typeFilter === 'physical' && !s.isPhysical) return false;
        if (typeFilter === 'esim' && s.isPhysical) return false;
      }
      if (statusFilter !== 'all' && statusKey(s) !== statusFilter) return false;
      return true;
    });
    const dir = sortDir === 'asc' ? 1 : -1;
    return [...filtered].sort((a, b) => {
      let cmp = 0;
      if (sortKey === 'type') {
        cmp = Number(a.isPhysical) - Number(b.isPhysical);
      } else if (sortKey === 'status') {
        cmp = statusKey(a).localeCompare(statusKey(b));
      } else {
        cmp = (parseTimestamp(a.createdAt) || 0) - (parseTimestamp(b.createdAt) || 0);
      }
      return cmp * dir;
    });
  }, [allSims, typeFilter, statusFilter, sortKey, sortDir]);

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'SIM pool']}
        title="SIM pool"
        count={sectionValue(stats?.total, stats?.error)}
        sub={
          canAct
            ? 'Inventory of SIMs available to assign to customers.'
            : 'SIMs uploaded for this network (view-only).'
        }
        actions={
          canAct ? (
            <>
              <Button
                variant="outlined"
                startIcon={<DownloadRounded />}
                onClick={() => toast('Exported SIM pool')}
              >
                Export
              </Button>
              <Button
                variant="contained"
                startIcon={<UploadFileRounded />}
                onClick={() => setShowUpload(true)}
              >
                Upload SIMs
              </Button>
            </>
          ) : undefined
        }
      />
      <KpiRow
        items={[
          {
            icon: 'sim_card',
            label: 'Assigned',
            value: sectionValue(stats?.consumed, stats?.error),
            sub: stats?.pctAssigned != null ? `${stats.pctAssigned}% of pool` : undefined,
            color: 'var(--uk-ac)',
          },
          {
            icon: 'info',
            label: 'Available',
            value: sectionValue(stats?.available, stats?.error),
            color: 'var(--uk-success-bright)',
          },
          {
            icon: 'error',
            label: 'Faulty',
            value: sectionValue(stats?.failed, stats?.error),
            color: 'var(--uk-error)',
          },
          {
            icon: 'sim_card',
            label: 'eSIM / physical',
            value: stats?.error
              ? '—'
              : `${stats?.esim ?? 0} / ${stats?.physical ?? 0}`,
            color: 'var(--uk-secondary)',
          },
        ]}
      />

      <div className="card card-pad">
        <div
          className="sec-head"
          style={{ flexWrap: 'wrap', gap: 12, rowGap: 12 }}
        >
          <div className="sec-title">
            SIMs <span className="cnt tnum">{sims.length}</span>
          </div>
          <div style={{ display: 'flex', gap: 16, flexWrap: 'wrap' }}>
            <FilterChips
              value={typeFilter}
              onChange={(v) => setTypeFilter(v as TypeFilter)}
              options={[
                { value: 'all', label: 'All types' },
                { value: 'physical', label: 'Physical' },
                { value: 'esim', label: 'eSIM' },
              ]}
            />
            <FilterChips
              value={statusFilter}
              onChange={(v) => setStatusFilter(v as StatusFilter)}
              options={[
                { value: 'all', label: 'All statuses' },
                { value: 'available', label: 'Available' },
                { value: 'assigned', label: 'Assigned' },
                { value: 'faulty', label: 'Faulty' },
              ]}
            />
          </div>
        </div>
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={5} rows={4} />
          ) : simsSection?.error ? (
            <EmptyState
              art="error"
              title="Couldn't load SIMs"
              sub={simsSection.error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : allSims.length === 0 ? (
            <EmptyState art="sim" title="No SIMs" sub="Upload a SIM batch to get started." />
          ) : sims.length === 0 ? (
            <EmptyState
              art="search"
              title="No SIMs match"
              sub="Try a different type or status filter."
            />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ICCID</TableCell>
                  <TableCell sortDirection={sortKey === 'type' ? sortDir : false}>
                    <TableSortLabel
                      active={sortKey === 'type'}
                      direction={sortKey === 'type' ? sortDir : 'asc'}
                      onClick={() => toggleSort('type')}
                    >
                      Type
                    </TableSortLabel>
                  </TableCell>
                  <TableCell sortDirection={sortKey === 'status' ? sortDir : false}>
                    <TableSortLabel
                      active={sortKey === 'status'}
                      direction={sortKey === 'status' ? sortDir : 'asc'}
                      onClick={() => toggleSort('status')}
                    >
                      Status
                    </TableSortLabel>
                  </TableCell>
                  <TableCell sortDirection={sortKey === 'added' ? sortDir : false}>
                    <TableSortLabel
                      active={sortKey === 'added'}
                      direction={sortKey === 'added' ? sortDir : 'asc'}
                      onClick={() => toggleSort('added')}
                    >
                      Added
                    </TableSortLabel>
                  </TableCell>
                  {canAct && <TableCell sx={{ width: 44 }} />}
                </TableRow>
              </TableHead>
              <TableBody>
                {sims.map((sim) => (
                  <TableRow key={sim.id}>
                    <TableCell className="tnum" style={{ fontWeight: 600 }}>
                      {sim.iccid}
                    </TableCell>
                    <TableCell>{sim.isPhysical ? 'Physical' : 'eSIM'}</TableCell>
                    <TableCell>
                      <StatusBadge status={simStatus(sim)}>{simStatusLabel(sim)}</StatusBadge>
                    </TableCell>
                    <TableCell className="muted">{formatDate(sim.createdAt)}</TableCell>
                    {canAct && (
                      <TableCell>
                        <SimMenu iccid={sim.iccid} />
                      </TableCell>
                    )}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && !simsSection?.error && (
          (() => {
            // When the pool is larger than what the list returns, show
            // "Showing N of M" so the cap is clear vs. the stats total.
            const total = stats?.error ? undefined : stats?.total;
            const noFilter = typeFilter === 'all' && statusFilter === 'all';
            return noFilter && total != null && total > allSims.length ? (
              <TableFooter showing={allSims.length} total={total} />
            ) : (
              <TableFooter count={sims.length} noun="SIMs" />
            );
          })()
        )}
      </div>
      {showUpload && (
        <UploadSimsDialog
          onClose={() => setShowUpload(false)}
          onUploaded={() => void refetch()}
        />
      )}
    </div>
  );
}
