/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * SIM pool — inventory cockpit with stock levels and proactive low-stock
 * nudge (screens-manage.jsx). Business can act; Network is view-only.
 */
import { useState } from 'react';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import ShoppingCartRounded from '@mui/icons-material/ShoppingCartRounded';
import UploadFileRounded from '@mui/icons-material/UploadFileRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import { KpiRow } from '@/components/Kpi';
import PageHeader from '@/components/PageHeader';
import { useToast } from '@/components/ToastProvider';
import { SIM_BATCHES, SIMS_SUMMARY } from '@/data';
import type { SimBatch } from '@/data';
import { useFirstLoad } from '@/lib/useFirstLoad';
import UploadSimsDialog from './UploadSimsDialog';

function BatchMenu({ batch }: { batch: SimBatch }) {
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
            toast(`Viewing SIMs in ${batch.batch}`);
          }}
        >
          <VisibilityRounded sx={{ fontSize: 18 }} /> View SIMs
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Exported ${batch.batch}`);
          }}
        >
          <DownloadRounded sx={{ fontSize: 18 }} /> Export batch
        </MenuItem>
      </Menu>
    </>
  );
}

export default function SimPoolScreen({ canAct }: { canAct: boolean }) {
  const s = SIMS_SUMMARY;
  const loading = useFirstLoad('simpool');
  const toast = useToast();
  const [showUpload, setShowUpload] = useState(false);

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'SIM pool']}
        title="SIM pool"
        count={s.total.toLocaleString()}
        sub={
          canAct
            ? 'Inventory of SIMs available to assign to customers.'
            : 'SIM batches uploaded for this network (view-only).'
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
          { icon: 'sim_card', label: 'Assigned', value: s.assigned.toLocaleString(), color: 'var(--uk-ac)' },
          { icon: 'info', label: 'Available', value: s.available, color: 'var(--uk-success-bright)' },
          { icon: 'warning', label: 'Suspended', value: s.suspended, color: 'var(--uk-orange)' },
          { icon: 'error', label: 'Faulty', value: s.faulty, color: 'var(--uk-error)' },
        ]}
      />
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
          <b style={{ color: 'var(--uk-ink)' }}>Stock is getting low.</b> {s.available} SIMs
          available — below your 700 reorder threshold.
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

      <div className="card card-pad">
        <div className="sec-head">
          <div className="sec-title">
            Batches <span className="cnt tnum">{SIM_BATCHES.length}</span>
          </div>
        </div>
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={6} rows={4} />
          ) : (
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Batch</th>
                  <th>Type</th>
                  <th className="num">Quantity</th>
                  <th>Assigned</th>
                  <th>Uploaded</th>
                  {canAct && <th style={{ width: 40 }} />}
                </tr>
              </thead>
              <tbody>
                {SIM_BATCHES.map((b) => {
                  const pct = Math.round((b.assigned / b.qty) * 100);
                  return (
                    <tr key={b.id} className="static">
                      <td className="tnum" style={{ fontWeight: 600 }}>
                        {b.batch}
                      </td>
                      <td>{b.type}</td>
                      <td className="num tnum">{b.qty.toLocaleString()}</td>
                      <td>
                        <div className="usebar" style={{ width: 184 }}>
                          <div className="meter">
                            <span style={{ width: pct + '%' }} />
                          </div>
                          <span
                            className="tnum"
                            style={{ fontSize: 12, color: 'var(--uk-ink-2)', whiteSpace: 'nowrap' }}
                          >
                            {b.assigned} · {pct}%
                          </span>
                        </div>
                      </td>
                      <td className="muted">{b.uploaded}</td>
                      {canAct && (
                        <td>
                          <BatchMenu batch={b} />
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>
        {!loading && <TableFooter count={SIM_BATCHES.length} noun="batches" />}
      </div>
      {showUpload && <UploadSimsDialog onClose={() => setShowUpload(false)} />}
    </div>
  );
}
