/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Node detail — six tabs, left health rail, node-board imagery and the
 * "Turn node off" power menu (node-site-detail.jsx NodeDetail).
 */
import { useState } from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ArrowBackRounded from '@mui/icons-material/ArrowBackRounded';
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';
import PowerSettingsNewRounded from '@mui/icons-material/PowerSettingsNewRounded';
import RestartAltRounded from '@mui/icons-material/RestartAltRounded';
import SyncRounded from '@mui/icons-material/SyncRounded';
import WifiOffRounded from '@mui/icons-material/WifiOffRounded';
import { LineChart } from '@/components/charts';
import DetailPicker from '@/components/DetailPicker';
import KV from '@/components/KV';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { NODES } from '@/data';
import { series } from '@/lib/series';

const NODE_TEMP = series(46, 22, 0.18, 0.12);
const NODE_LOAD = series(52, 22, 0.14, 0.18);
const TABS = ['Overview', 'Network', 'Resources', 'Radio', 'Software', 'Schematic'];

const TEMP_LEGEND = [
  { color: 'var(--uk-ac)', label: 'Below 50° normal' },
  { color: 'var(--uk-orange)', label: '51–60° high' },
  { color: 'var(--uk-error)', label: 'Above 60° critical' },
];
const LOAD_LEGEND = [
  { color: 'var(--uk-ac)', label: 'Below 70% normal' },
  { color: 'var(--uk-orange)', label: '71–90% high' },
  { color: 'var(--uk-error)', label: 'Above 90% critical' },
];

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontSize: 12, color: 'var(--uk-ink-2)' }}>
      <span style={{ width: 9, height: 9, borderRadius: 3, background: color }} />
      {label}
    </span>
  );
}

function HealthChartCard({
  title,
  unit,
  data,
  legend,
}: {
  title: string;
  unit: string;
  data: number[];
  legend: { color: string; label: string }[];
}) {
  return (
    <SectionCard title={title} right={<span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{unit}</span>}>
      <LineChart data={data} height={188} />
      <div style={{ display: 'flex', gap: 18, justifyContent: 'center', marginTop: 10 }}>
        {legend.map((l) => (
          <LegendDot key={l.label} {...l} />
        ))}
      </div>
    </SectionCard>
  );
}

function nodeStatusText(status: string, up: string) {
  if (status === 'online') return `is online and well for ${up}`;
  if (status === 'configuring') return 'is configuring · ~4 min remaining';
  if (status === 'degraded') return 'is online with warnings';
  return 'has been offline for 2h 14m';
}

function PowerMenu({ serial }: { serial: string }) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const toast = useToast();
  return (
    <>
      <Button
        variant="contained"
        startIcon={<PowerSettingsNewRounded />}
        endIcon={<ExpandMoreRounded />}
        sx={{
          bgcolor: '#1C1E22',
          '&:hover': { bgcolor: '#2c2f36' },
        }}
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        Turn node off
      </Button>
      <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`${serial} powered off`);
          }}
        >
          <PowerSettingsNewRounded sx={{ fontSize: 18 }} /> Turn node off
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Restarting ${serial}…`);
          }}
        >
          <RestartAltRounded sx={{ fontSize: 18 }} /> Restart node
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast('RF disabled');
          }}
        >
          <WifiOffRounded sx={{ fontSize: 18 }} /> Turn RF off
        </MenuItem>
      </Menu>
    </>
  );
}

export default function NodeDetailScreen({ nodeId }: { nodeId: string }) {
  const router = useRouter();
  const toast = useToast();
  const n = NODES.find((x) => x.id === nodeId) ?? NODES[0];
  const [tab, setTab] = useState('Overview');
  if (!n) return null;

  const off = n.status === 'offline';
  const hot = (n.temp ?? 0) > 55;
  const tempTrx = off ? null : n.temp;
  const tempCom = off || n.temp == null ? null : n.temp + 3;

  return (
    <div className="page">
      <PageHeader
        crumb={['Nodes', n.serial]}
        title={`Node ${n.serial}`}
        actions={
          <>
            <Button
              variant="outlined"
              startIcon={<ArrowBackRounded />}
              onClick={() => router.push('/network/nodes')}
            >
              Back
            </Button>
            <PowerMenu serial={n.serial} />
          </>
        }
      />

      <div className="detail-subrow">
        <DetailPicker
          value={{ id: n.id, label: n.serial, status: n.status }}
          items={NODES.map((x) => ({ id: x.id, label: x.serial, status: x.status }))}
          onPick={(it) => router.push(`/network/nodes/${it.id}`)}
        />
        <StatusBadge status={n.status} />
        <span style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>
          {nodeStatusText(n.status, n.up)}
        </span>
      </div>

      <div className="tabs scroll-x">
        {TABS.map((t) => (
          <button
            key={t}
            type="button"
            className={`tab${tab === t ? ' on' : ''}`}
            onClick={() => setTab(t)}
          >
            {t}
          </button>
        ))}
      </div>

      <div className="detail-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
          <SectionCard title="Node information">
            <KV k="Model type" v={n.type} />
            <KV k="Serial #" v={n.serial} />
            <KV k="Firmware" v={n.fw} />
            <KV k="Node group" v="Default group" />
          </SectionCard>
          <SectionCard title="Node health">
            <KV
              k="Temp. (TRX)"
              v={tempTrx != null ? tempTrx + ' °C' : '—'}
              warn={hot}
              vColor={hot ? 'var(--uk-orange)' : null}
            />
            <KV
              k="Temp. (COM)"
              v={tempCom != null ? tempCom + ' °C' : '—'}
              warn={hot}
              vColor={hot ? 'var(--uk-orange)' : null}
            />
            <KV k="CPU load" v={off ? '—' : n.cpu + '%'} />
            <KV k="Memory" v={off ? '—' : n.mem + '%'} />
          </SectionCard>
          <SectionCard title="Customers">
            <KV k="Attached" v={off ? '0' : '100'} />
            <KV k="Active" v={off ? '0' : '86'} />
          </SectionCard>
        </div>

        <div style={{ minWidth: 0 }}>
          {tab === 'Overview' && (
            <SectionCard
              title="Node hardware"
              right={<span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{n.type}</span>}
            >
              <div style={{ display: 'flex', justifyContent: 'center', padding: '12px 0' }}>
                <Image src="/node-board.png" alt="Node board" width={300} height={420} style={{ height: 'auto' }} />
              </div>
            </SectionCard>
          )}
          {(tab === 'Network' || tab === 'Resources' || tab === 'Radio') && (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
              <HealthChartCard
                title="Temperature"
                unit="°C · last 6h"
                data={off ? NODE_TEMP.map(() => 0) : NODE_TEMP}
                legend={TEMP_LEGEND}
              />
              <HealthChartCard
                title="Load index"
                unit="% · last 6h"
                data={off ? NODE_LOAD.map(() => 0) : NODE_LOAD}
                legend={LOAD_LEGEND}
              />
            </div>
          )}
          {tab === 'Software' && (
            <SectionCard title="Software & firmware">
              <KV k="Firmware version" v={n.fw} />
              <KV k="Channel" v="Stable" />
              <KV k="Last update" v="12 Oct 2025" />
              <KV
                k="Update available"
                v={n.fw === '13.2.1' ? 'Up to date' : '13.2.1'}
                vColor={n.fw === '13.2.1' ? 'var(--uk-success)' : 'var(--uk-ac-dark)'}
              />
              <div style={{ marginTop: 16 }}>
                <Button
                  variant="outlined"
                  startIcon={<SyncRounded />}
                  onClick={() => toast('Checking for updates…')}
                >
                  Check for updates
                </Button>
              </div>
            </SectionCard>
          )}
          {tab === 'Schematic' && (
            <SectionCard
              title="Schematic"
              right={<span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>Port layout</span>}
            >
              <div style={{ display: 'flex', justifyContent: 'center', padding: '12px 0' }}>
                <Image
                  src="/node-board.png"
                  alt="Node schematic"
                  width={300}
                  height={420}
                  style={{ height: 'auto' }}
                />
              </div>
            </SectionCard>
          )}
        </div>
      </div>
    </div>
  );
}
