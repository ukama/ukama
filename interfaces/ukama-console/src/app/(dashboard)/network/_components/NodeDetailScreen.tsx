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
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';
import PowerSettingsNewRounded from '@mui/icons-material/PowerSettingsNewRounded';
import RestartAltRounded from '@mui/icons-material/RestartAltRounded';
import SyncRounded from '@mui/icons-material/SyncRounded';
import WifiOffRounded from '@mui/icons-material/WifiOffRounded';
import Skeleton from '@mui/material/Skeleton';

import { useNodeDetailQuery } from '@/client/graphql/node-detail.generated';
import AppTabs from '@/components/AppTabs';
import { LineChart } from '@/components/charts';
import DetailPicker from '@/components/DetailPicker';
import { EmptyState } from '@/components/EmptyState';
import KV from '@/components/KV';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { sectionValue } from '@/components/SectionFallback';
import { useToast } from '@/components/ToastProvider';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { POLL_LIVE_MS, visiblePoll } from '@/lib/polling';
import { series } from '@/lib/series';
import { StateChip } from './nodeStatus';

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
  const [tab, setTab] = useState('Overview');

  const { data, loading, refetch } = useNodeDetailQuery({
    variables: { nodeId },
    ...visiblePoll(POLL_LIVE_MS),
  });
  const view = data?.nodeView;
  const nodeSection = view?.node;
  const healthSection = view?.health;
  const softwareSection = view?.software;
  const kpisGap = view?.kpis.error ?? null;
  const kpiByKey = new Map(
    (view?.kpis.metrics ?? [])
      .filter((m) => m.success)
      .map((m) => [m.key, Math.round(m.value * 100) / 100])
  );

  if (loading) {
    return (
      <div className="page">
        <PageHeader crumb={['Nodes', nodeId]} title={`Node ${nodeId}`} />
        <Skeleton variant="rounded" sx={{ height: 42, mb: 2 }} />
        <Skeleton variant="rounded" sx={{ height: 420 }} />
      </div>
    );
  }
  if (!nodeSection?.node) {
    return (
      <div className="page">
        <PageHeader crumb={['Nodes', nodeId]} title={`Node ${nodeId}`} />
        <div className="card">
          <EmptyState
            art="error"
            title="Couldn't load node"
            sub={nodeSection?.error?.message ?? 'Node not found.'}
            cta="Try again"
            onCta={() => refetch()}
          />
        </div>
      </div>
    );
  }

  const n = toUkamaNode(nodeSection.node);
  const nodeName = n.name ?? n.serial;
  const off = n.status === 'offline';
  const healthRows = healthSection?.health?.system?.slice(0, 4) ?? [];
  const firstSoftware = softwareSection?.softwares?.software?.[0];
  const updateAvailable = firstSoftware?.status === 'update_available';

  return (
    <div className="page">
      <PageHeader
        crumb={['Nodes', n.serial]}
        title={nodeName}
        actions={<PowerMenu serial={n.serial} />}
      />

      <div className="detail-subrow">
        <DetailPicker
          value={{ id: n.id, label: `${nodeName} (${n.id})`, status: n.status }}
          items={[{ id: n.id, label: `${nodeName} (${n.id})`, status: n.status }]}
          onPick={(it) => router.push(`/network/nodes/${it.id}`)}
        />
        <StateChip state={n.state} />
      </div>

      <AppTabs tabs={TABS} value={tab} onChange={setTab} scrollable />

      <div className="detail-grid">
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
          <SectionCard title="Node information">
            <KV k="Model type" v={n.type} />
            <KV k="Serial #" v={n.serial} />
            <KV
              k="Firmware"
              v={firstSoftware?.currentVersion ?? '—'}
            />
            <KV k="Site" v={n.site} />
          </SectionCard>
          <SectionCard title="Node health">
            {healthRows.length > 0 && !healthSection?.error ? (
              healthRows.map((row) => <KV key={row.name} k={row.name} v={row.value} />)
            ) : kpiByKey.size > 0 ? (
              // Polled node KPIs from the metric service (Phase 4)
              <>
                <KV k="Uptime" v={kpiByKey.has('uptime') ? `${kpiByKey.get('uptime')}` : '—'} />
                <KV
                  k="Temp. (CPU)"
                  v={kpiByKey.has('cpu_temperature') ? `${kpiByKey.get('cpu_temperature')} °C` : '—'}
                />
                <KV
                  k="Memory"
                  v={kpiByKey.has('memory') ? `${kpiByKey.get('memory')}%` : '—'}
                />
              </>
            ) : (
              <>
                <KV k="Uptime" v="—" />
                <KV k="Temp. (CPU)" v="—" />
                <KV k="Memory" v="—" />
              </>
            )}
          </SectionCard>
          <SectionCard title="Customers">
            {/* attach counts not in metric keys yet — renders "—" */}
            <KV k="Attached" v={sectionValue(null, kpisGap)} />
            <KV k="Active" v={sectionValue(null, kpisGap)} />
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
              {softwareSection?.error ? (
                <>
                  <KV k="Firmware version" v="—" />
                  <KV k="Update available" v="—" />
                </>
              ) : (
                <>
                  <KV k="Firmware version" v={firstSoftware?.currentVersion ?? '—'} />
                  <KV k="Release date" v={firstSoftware?.releaseDate ?? '—'} />
                  <KV
                    k="Update available"
                    v={updateAvailable ? (firstSoftware?.desiredVersion ?? 'Yes') : 'Up to date'}
                    vColor={updateAvailable ? 'var(--uk-ac-dark)' : 'var(--uk-success)'}
                  />
                </>
              )}
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
