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
import { useState } from 'react';
import Button from '@mui/material/Button';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import ShoppingCartRounded from '@mui/icons-material/ShoppingCartRounded';
import UploadFileRounded from '@mui/icons-material/UploadFileRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';

import { useSimPoolOverviewQuery } from '@/client/graphql/sim-pool.generated';
import { EmptyState } from '@/components/EmptyState';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { useToast } from '@/components/ToastProvider';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import UploadSimsDialog from './UploadSimsDialog';

const SIM_TYPE = 'ukama_data';
const LIST_LIMIT = 50;

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

  const { data, loading, refetch } = useSimPoolOverviewQuery({
    variables: { simType: SIM_TYPE, limit: LIST_LIMIT },
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const stats = data?.simPoolView.stats;
  const simsSection = data?.simPoolView.sims;
  const sims = simsSection?.sims ?? [];

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
      {stats?.lowStock && (
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 11,
            background: 'rgba(226,116,41,.09)',
            border: '1px solid rgba(226,116,41,.22)',
            borderRadius: 10,
            padding: '12px 16px',
            marginBottom: 'var(--uk-gap)',
          }}
        >
          <InfoRounded sx={{ color: '#b5591b', fontSize: 20 }} />
          <span style={{ fontSize: 13, color: 'var(--uk-ink-2)', flex: 1 }}>
            <b style={{ color: 'var(--uk-ink)' }}>Stock is getting low.</b>{' '}
            {stats.available} SIMs available — below the reorder threshold.
          </span>
          {canAct && (
            <Button
              size="small"
              variant="outlined"
              startIcon={<ShoppingCartRounded />}
              onClick={() => toast('Order placed with Ukama supply')}
            >
              Order SIMs
            </Button>
          )}
        </div>
      )}

      <div className="card card-pad">
        <div className="sec-head">
          <div className="sec-title">
            SIMs <span className="cnt tnum">{sims.length}</span>
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
          ) : sims.length === 0 ? (
            <EmptyState art="sim" title="No SIMs" sub="Upload a SIM batch to get started." />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ICCID</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Added</TableCell>
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
                    <TableCell className="muted">{sim.createdAt}</TableCell>
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
          <TableFooter count={sims.length} noun="SIMs" />
        )}
      </div>
      {showUpload && <UploadSimsDialog onClose={() => setShowUpload(false)} />}
    </div>
  );
}
