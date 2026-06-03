/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Site drawer — reveal-on-demand site summary (detail.jsx SiteDrawer). */
import Button from '@mui/material/Button';
import BuildRounded from '@mui/icons-material/BuildRounded';
import EditRounded from '@mui/icons-material/EditRounded';
import ErrorOutlineRounded from '@mui/icons-material/ErrorOutlineRounded';
import PlaceRounded from '@mui/icons-material/PlaceRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import AppDrawer, { DetailRow, DrawerHead } from '@/components/AppDrawer';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { ALERTS, NODES } from '@/data';
import type { Site, UkamaNode } from '@/data';
import { Ic } from '../../_components/icons';

export default function SiteDrawer({
  site,
  onClose,
  onManage,
  onOpenNode,
}: {
  site: Site;
  onClose: () => void;
  onManage: (site: Site) => void;
  onOpenNode: (node: UkamaNode) => void;
}) {
  const toast = useToast();
  const nodes = NODES.filter((n) => n.site === site.name);
  const issues = ALERTS.filter((a) => a.site === site.name);
  const battColor =
    site.battery < 20
      ? 'var(--uk-error)'
      : site.battery < 50
        ? 'var(--uk-orange)'
        : 'var(--uk-success-bright)';

  return (
    <AppDrawer onClose={onClose} width={460}>
      <DrawerHead
        title={site.name}
        sub={
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: 3 }}>
            <PlaceRounded sx={{ fontSize: 14 }} /> {site.area}
          </span>
        }
        badge={<StatusBadge status={site.status} />}
        onClose={onClose}
      />
      <div style={{ flex: 1, overflow: 'auto', padding: '18px 24px' }}>
        {site.issue && (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 9,
              marginBottom: 16,
              padding: '11px 13px',
              borderRadius: 10,
              background:
                site.status === 'offline' ? 'var(--uk-error-fill)' : 'rgba(226,116,41,.13)',
              color: site.status === 'offline' ? 'var(--uk-error-deep, #cf121b)' : '#b5591b',
              fontSize: 13,
              fontWeight: 500,
            }}
          >
            <ErrorOutlineRounded sx={{ fontSize: 18 }} />
            {site.issue}
          </div>
        )}
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', marginBottom: 18 }}>
          <div className="card card-pad" style={{ padding: 14 }}>
            <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>Subscribers</div>
            <div
              className="tnum"
              style={{ fontFamily: 'var(--font-display)', fontSize: 22, fontWeight: 500 }}
            >
              {site.subs}
            </div>
          </div>
          <div className="card card-pad" style={{ padding: 14 }}>
            <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>Uptime (30d)</div>
            <div
              className="tnum"
              style={{ fontFamily: 'var(--font-display)', fontSize: 22, fontWeight: 500 }}
            >
              {site.uptime}%
            </div>
          </div>
        </div>
        <div style={{ marginBottom: 6 }}>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              fontSize: 12.5,
              color: 'var(--uk-ink-2)',
              marginBottom: 6,
            }}
          >
            <span>Battery · {site.plan}</span>
            <span className="tnum">{site.battery}%</span>
          </div>
          <div className="meter">
            <span style={{ width: site.battery + '%', background: battColor }} />
          </div>
        </div>
        <div style={{ marginTop: 14 }}>
          <DetailRow k="Signal strength" v={site.signal ? site.signal + ' dBm' : '—'} />
          <DetailRow k="Data (30d)" v={site.data} />
          <DetailRow k="Power plan" v={site.plan} />
        </div>
        <div className="sec-head" style={{ margin: '22px 0 10px' }}>
          <div className="sec-title" style={{ fontSize: 14 }}>
            Nodes <span className="cnt tnum">{nodes.length}</span>
          </div>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {nodes.map((n) => {
            const NIcon = n.type.startsWith('Amp')
              ? SettingsInputAntennaRounded
              : RouterRounded;
            return (
              <button
                key={n.id}
                type="button"
                className="card"
                onClick={() => onOpenNode(n)}
                style={{
                  padding: '11px 13px',
                  display: 'flex',
                  alignItems: 'center',
                  gap: 11,
                  cursor: 'pointer',
                  textAlign: 'left',
                  fontFamily: 'inherit',
                  border: '1px solid var(--uk-line)',
                }}
              >
                <NIcon sx={{ fontSize: 20, color: 'var(--uk-ac)' }} />
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 13, fontWeight: 600 }}>{n.type}</div>
                  <div className="tnum" style={{ fontSize: 11.5, color: 'var(--uk-ink-3)' }}>
                    {n.serial}
                  </div>
                </div>
                <StatusBadge status={n.status} />
              </button>
            );
          })}
        </div>
        {issues.length > 0 && (
          <>
            <div className="sec-head" style={{ margin: '22px 0 10px' }}>
              <div className="sec-title" style={{ fontSize: 14 }}>
                Recent alerts <span className="cnt tnum">{issues.length}</span>
              </div>
            </div>
            {issues.map((a) => (
              <div
                key={a.id}
                style={{ display: 'flex', gap: 10, padding: '9px 0', fontSize: 12.5, color: 'var(--uk-ink-2)' }}
              >
                <Ic name={a.icon} sx={{ fontSize: 17, color: 'var(--uk-orange)' }} />
                <span>
                  {a.title} · <span style={{ color: 'var(--uk-ink-3)' }}>{a.age} ago</span>
                </span>
              </div>
            ))}
          </>
        )}
      </div>
      <div style={{ padding: '14px 24px', borderTop: '1px solid var(--uk-line)', display: 'flex', gap: 10 }}>
        <Button
          variant="contained"
          startIcon={<BuildRounded />}
          sx={{ flex: 1 }}
          onClick={() => onManage(site)}
        >
          {site.status === 'offline' ? 'Diagnose' : 'Manage site'}
        </Button>
        <Button
          variant="outlined"
          startIcon={<EditRounded />}
          onClick={() => toast(`Rename ${site.name}`)}
        >
          Rename
        </Button>
      </div>
    </AppDrawer>
  );
}
