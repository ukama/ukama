/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Node drawer — quick health summary (detail.jsx NodeDrawer). */
import Button from '@mui/material/Button';
import SyncRounded from '@mui/icons-material/SyncRounded';
import AppDrawer, { DetailRow, DrawerHead } from '@/components/AppDrawer';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { UkamaNode } from '@/data';

function MeterCard({ label, value }: { label: string; value: number | null }) {
  const v = value ?? 0;
  const color =
    v > 75 ? 'var(--uk-error)' : v > 60 ? 'var(--uk-orange)' : 'var(--uk-success-bright)';
  return (
    <div className="card card-pad" style={{ padding: 14 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          fontSize: 12,
          color: 'var(--uk-ink-2)',
          marginBottom: 6,
        }}
      >
        <span>{label}</span>
        <span className="tnum">{value == null ? '—' : v + '%'}</span>
      </div>
      <div className="meter">
        <span style={{ width: v + '%', background: color }} />
      </div>
    </div>
  );
}

export default function NodeDrawer({
  node,
  onClose,
  onOpenDetail,
}: {
  node: UkamaNode;
  onClose: () => void;
  onOpenDetail: (node: UkamaNode) => void;
}) {
  const toast = useToast();
  const off = node.status === 'offline';

  return (
    <AppDrawer onClose={onClose} width={420}>
      <DrawerHead
        title={node.type}
        sub={<span className="tnum">{node.serial}</span>}
        badge={<StatusBadge status={node.status} />}
        onClose={onClose}
      />
      <div style={{ flex: 1, overflow: 'auto', padding: '18px 24px' }}>
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', marginBottom: 18 }}>
          <MeterCard label="CPU" value={off ? null : node.cpu} />
          <MeterCard label="Memory" value={off ? null : node.mem} />
        </div>
        <DetailRow k="Site" v={node.site} />
        <DetailRow k="Temperature" v={node.temp ? node.temp + ' °C' : '—'} />
        <DetailRow k="Firmware" v={node.fw} />
        <DetailRow k="Uptime" v={node.up} />
        {node.note && <DetailRow k="Note" v={node.note} vColor="var(--uk-orange)" />}
        <div style={{ marginTop: 16 }}>
          <Button variant="text" onClick={() => onOpenDetail(node)}>
            Open full details →
          </Button>
        </div>
      </div>
      <div style={{ padding: '14px 24px', borderTop: '1px solid var(--uk-line)', display: 'flex', gap: 10 }}>
        <Button
          variant="contained"
          sx={{ flex: 1 }}
          onClick={() => toast(off ? `Diagnosing ${node.serial}…` : `Restarting ${node.serial}…`)}
        >
          {off ? 'Diagnose' : 'Restart node'}
        </Button>
        <Button
          variant="outlined"
          startIcon={<SyncRounded />}
          onClick={() => toast('Firmware update queued')}
        >
          Update FW
        </Button>
      </div>
    </AppDrawer>
  );
}
