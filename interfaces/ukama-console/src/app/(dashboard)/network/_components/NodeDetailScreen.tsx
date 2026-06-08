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
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';
import WifiOffRounded from '@mui/icons-material/WifiOffRounded';
import Skeleton from '@mui/material/Skeleton';

import { useNodeDetailQuery } from '@/client/graphql/node-detail.generated';
import { useNodeKpisQuery } from '@/client/graphql/node-kpis.generated';
import AppTabs from '@/components/AppTabs';
import { LineChart } from '@/components/charts';
import DetailPicker from '@/components/DetailPicker';
import { EmptyState } from '@/components/EmptyState';
import KV from '@/components/KV';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { useToast } from '@/components/ToastProvider';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { series } from '@/lib/series';
import { StateChip } from './nodeStatus';

const TABS = ['Overview', 'Network', 'Resources', 'Radio', 'Software'];

/** Product imagery keyed by raw node type (tnode/anode/cnode/hnode). */
const NODE_IMAGE_BASE = 'https://ukama-site-assets.s3.amazonaws.com/images';
const NODE_IMAGES: Record<string, string> = {
  tnode: `${NODE_IMAGE_BASE}/ukama_tower_node.png`,
  anode: `${NODE_IMAGE_BASE}/ukama_amplifier_node.png`,
  cnode: `${NODE_IMAGE_BASE}/ukama_home_node.png`,
  hnode: `${NODE_IMAGE_BASE}/ukama_home_node.png`,
};
const NODE_IMAGE_FALLBACK = `${NODE_IMAGE_BASE}/ukama_tower_node.png`;
const nodeImage = (type: string): string =>
  NODE_IMAGES[type.toLowerCase()] ?? NODE_IMAGE_FALLBACK;

const GENERIC_LEGEND = [{ color: 'var(--uk-ac)', label: 'Trend' }];

type Range = 'Day' | 'Week' | 'Month';
const RANGES: Range[] = ['Day', 'Week', 'Month'];
const WINDOW_LABEL: Record<Range, string> = {
  Day: 'last 24h',
  Week: 'last 7d',
  Month: 'last 30d',
};

interface ChartDef {
  key: string;
  title: string;
  unit: string;
  legend: { color: string; label: string }[];
  data: Record<Range, number[]>;
}

/** Three placeholder series (one per range) for a metric — until the metric
 *  service is wired (backend gap #6). */
function ranged(base: number, a: number, b: number): Record<Range, number[]> {
  return {
    Day: series(base, 24, a, b),
    Week: series(base, 28, a, b * 1.3),
    Month: series(base, 30, a, b * 1.7),
  };
}

/* -------------------------------------------------------------------------- *
 * Node KPI catalog — labels, units and per-node-type key lists mirrored from
 * the legacy console (constants NODE KPI config) and console-bff
 * getGraphsKeyByType(GRAPHS_TYPE.*, nodeType). Values come from the NodeKpis
 * query once the metric service is wired (backend gap #6).
 * -------------------------------------------------------------------------- */
type NodeKind = 'tnode' | 'anode' | 'cnode' | 'hnode';
const nodeKind = (raw: string): NodeKind => {
  const t = raw.toLowerCase();
  if (t.includes('anode')) return 'anode';
  if (t.includes('cnode')) return 'cnode';
  if (t.includes('hnode')) return 'hnode';
  return 'tnode';
};

const METRIC_META: Record<string, { label: string; unit: string; base: number }> = {
  uptime: { label: 'Uptime', unit: 's', base: 86 },
  cpu_temperature: { label: 'Temp. (CPU)', unit: '°C', base: 46 },
  fem1_temperature: { label: 'FEM 1 temp.', unit: '°C', base: 44 },
  fem2_temperature: { label: 'FEM 2 temp.', unit: '°C', base: 48 },
  memory: { label: 'Memory', unit: '%', base: 52 },
  cpu: { label: 'CPU', unit: '%', base: 38 },
  disk: { label: 'Disk', unit: 'MB', base: 60 },
  subscribers_active: { label: 'Active subscribers', unit: '', base: 30 },
  cellular_uplink: { label: 'Cellular uplink', unit: 'Mbps', base: 18 },
  cellular_downlink: { label: 'Cellular downlink', unit: 'Mbps', base: 64 },
  backhaul_uplink: { label: 'Backhaul uplink', unit: 'Mbps', base: 22 },
  backhaul_downlink: { label: 'Backhaul downlink', unit: 'Mbps', base: 70 },
  backhaul_latency: { label: 'Backhaul latency', unit: 'ms', base: 35 },
  power: { label: 'TX power', unit: 'dBm', base: 31 },
  pa_power: { label: 'PA power', unit: 'dBm', base: 30 },
  rx_power: { label: 'RX power', unit: 'dBm', base: 28 },
  tx_power: { label: 'TX power', unit: 'dBm', base: 31 },
};

type MetricGroup =
  | 'health'
  | 'customers'
  | 'cellular'
  | 'backhaul'
  | 'resources'
  | 'radio';
const GROUP_KEYS: Record<MetricGroup, Record<NodeKind, string[]>> = {
  health: {
    tnode: ['uptime', 'cpu_temperature', 'memory'],
    anode: ['uptime', 'fem1_temperature', 'fem2_temperature'],
    cnode: ['uptime', 'memory'],
    hnode: [],
  },
  customers: {
    tnode: ['subscribers_active'],
    anode: [],
    cnode: [],
    hnode: [],
  },
  cellular: {
    tnode: ['cellular_uplink', 'cellular_downlink'],
    anode: [],
    cnode: [],
    hnode: [],
  },
  backhaul: {
    tnode: ['backhaul_uplink', 'backhaul_downlink', 'backhaul_latency'],
    anode: [],
    cnode: [],
    hnode: [],
  },
  resources: {
    tnode: ['cpu', 'memory', 'disk'],
    anode: ['cpu', 'memory', 'disk'],
    cnode: ['cpu', 'memory', 'disk'],
    hnode: [],
  },
  radio: {
    tnode: ['power'],
    anode: ['pa_power', 'rx_power', 'tx_power'],
    cnode: [],
    hnode: [],
  },
};

const groupKeys = (group: MetricGroup, kind: NodeKind): string[] =>
  GROUP_KEYS[group][kind] ?? [];

const chartFor = (key: string): ChartDef => {
  const m = METRIC_META[key] ?? { label: key, unit: '', base: 40 };
  return {
    key,
    title: m.label,
    unit: m.unit,
    legend: GENERIC_LEGEND,
    data: ranged(m.base, 0.12, 0.14),
  };
};
const chartsFor = (group: MetricGroup, kind: NodeKind): ChartDef[] =>
  groupKeys(group, kind).map(chartFor);

/* Left-rail sections per tab. The 'info' card shows the node board; every
 * other card lists its group's KPIs and drives the right-side charts. Mirrors
 * the legacy console, where the left rail changes with the active tab. */
interface SectionDef {
  key: string;
  title: string;
  group?: MetricGroup;
}
const TAB_SECTIONS: Record<string, SectionDef[]> = {
  Overview: [
    { key: 'info', title: 'Node information' },
    { key: 'health', title: 'Node health', group: 'health' },
    { key: 'customers', title: 'Customers', group: 'customers' },
  ],
  Network: [
    { key: 'cellular', title: 'Cellular', group: 'cellular' },
    { key: 'backhaul', title: 'Backhaul', group: 'backhaul' },
  ],
  Resources: [{ key: 'resources', title: 'Resources', group: 'resources' }],
  Radio: [{ key: 'radio', title: 'Radio', group: 'radio' }],
};

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontSize: 12, color: 'var(--uk-ink-2)' }}>
      <span style={{ width: 9, height: 9, borderRadius: 3, background: color }} />
      {label}
    </span>
  );
}

function RangeToggle({ value, onChange }: { value: Range; onChange: (r: Range) => void }) {
  return (
    <div className="range-toggle" role="group" aria-label="Time range">
      {RANGES.map((r) => (
        <button
          key={r}
          type="button"
          className={r === value ? 'is-active' : ''}
          aria-pressed={r === value}
          onClick={() => onChange(r)}
        >
          {r}
        </button>
      ))}
    </div>
  );
}

function MetricChartCard({ chart, off }: { chart: ChartDef; off: boolean }) {
  const [range, setRange] = useState<Range>('Day');
  const data = off ? chart.data[range].map(() => 0) : chart.data[range];
  return (
    <SectionCard title={chart.title} right={<RangeToggle value={range} onChange={setRange} />}>
      <div style={{ fontSize: 12, color: 'var(--uk-ink-3)', marginTop: -2, marginBottom: 10 }}>
        {chart.unit} · {WINDOW_LABEL[range]}
      </div>
      <LineChart data={data} height={188} />
      <div style={{ display: 'flex', gap: 18, justifyContent: 'center', marginTop: 10 }}>
        {chart.legend.map((l) => (
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
  // Which left-rail section is selected (key from TAB_SECTIONS for the tab).
  const [section, setSection] = useState<string>('info');

  // Structural data — node core, site name and sibling nodes — as one
  // composite query (one-shot; no polling).
  const { data, loading, refetch } = useNodeDetailQuery({
    variables: { nodeId },
  });
  // Node health KPIs live in their own query so they can poll independently
  // once the metric service is wired (backend gap #6); plain fetch for now.
  const { data: kpisData } = useNodeKpisQuery({ variables: { nodeId } });

  const view = data?.nodeView;
  const nodeSection = view?.node;
  const softwareSection = view?.software;
  // Sibling nodes power the switcher dropdown.
  const pickerItems = (view?.siblings.nodes ?? []).map((nd) => ({
    id: nd.id,
    label: `${nd.name || nd.id} (${nd.id})`,
    status: '',
  }));
  const kpiByKey = new Map(
    (kpisData?.nodeView.kpis.metrics ?? [])
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

  const n = toUkamaNode(nodeSection.node, view?.site.site?.name ?? undefined);
  const nodeName = n.name ?? n.serial;
  const off = n.status === 'offline';
  const apps = softwareSection?.softwares?.software ?? [];

  // KPIs are node-type specific (legacy console parity, backend gap #6).
  const kind = nodeKind(nodeSection.node.type);
  const fmtKpi = (key: string): string => {
    if (!kpiByKey.has(key)) return '—';
    const v = kpiByKey.get(key);
    const unit = METRIC_META[key]?.unit ?? '';
    if (!unit) return `${v}`;
    return unit === '%' ? `${v}%` : `${v} ${unit}`;
  };
  // KV rows for a metric group (one per node-type key, value or "—").
  const groupRows = (group: MetricGroup) => {
    const keys = groupKeys(group, kind);
    if (keys.length === 0) return <KV k="Metrics" v="—" />;
    return keys.map((k) => (
      <KV key={k} k={METRIC_META[k]?.label ?? k} v={fmtKpi(k)} />
    ));
  };

  // Hide tabs that have no KPIs for this node type. Overview always shows
  // (node info), Software always shows (apps aren't KPI-gated).
  const tabHasKpis = (t: string) =>
    (TAB_SECTIONS[t] ?? []).some(
      (s) => s.group && groupKeys(s.group, kind).length > 0,
    );
  const visibleTabs = TABS.filter(
    (t) => t === 'Overview' || t === 'Software' || tabHasKpis(t),
  );
  const activeTab = visibleTabs.includes(tab) ? tab : 'Overview';

  // The left rail (and its selection) changes with the active tab.
  const sections = TAB_SECTIONS[activeTab] ?? [];
  const activeKey = sections.some((s) => s.key === section)
    ? section
    : (sections[0]?.key ?? 'info');
  const activeGroup = sections.find((s) => s.key === activeKey)?.group;
  const renderCharts = (charts: ChartDef[]) =>
    charts.length > 0 ? (
      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
        {charts.map((c) => (
          <MetricChartCard key={c.key} chart={c} off={off} />
        ))}
      </div>
    ) : (
      <SectionCard title="Metrics">
        <EmptyState
          art="search"
          title="No metrics for this node type"
          sub="This node type doesn't report metrics in this category."
        />
      </SectionCard>
    );

  return (
    <div className="page">
      <PageHeader
        crumb={['Nodes', n.serial]}
        title={nodeName}
        onBack={() => router.push('/network/nodes')}
        actions={<PowerMenu serial={n.serial} />}
      />

      <div className="detail-subrow">
        <DetailPicker
          value={{ id: n.id, label: `${nodeName} (${n.id})`, status: n.status }}
          items={
            pickerItems.length > 0
              ? pickerItems
              : [{ id: n.id, label: `${nodeName} (${n.id})`, status: n.status }]
          }
          onPick={(it) => router.push(`/network/nodes/${it.id}`)}
        />
        <StateChip state={n.state} />
      </div>

      <AppTabs tabs={visibleTabs} value={activeTab} onChange={setTab} scrollable />

      {activeTab === 'Software' ? (
        <NodeApps apps={apps} error={!!softwareSection?.error} onCheck={() => toast('Checking for updates…')} />
      ) : (
        <div className="detail-grid">
          <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
            {sections.map((s) => (
              <SectionCard
                key={s.key}
                title={s.title}
                selectable
                active={activeKey === s.key}
                onClick={() => setSection(s.key)}
              >
                {s.key === 'info' ? (
                  <>
                    <KV k="Model type" v={n.type} />
                    <KV k="Serial #" v={n.serial} />
                    <KV k="Site" v={n.site} />
                  </>
                ) : (
                  groupRows(s.group as MetricGroup)
                )}
              </SectionCard>
            ))}
          </div>

          <div style={{ minWidth: 0 }}>
            {activeKey === 'info' ? (
              <SectionCard
                title="Node hardware"
                right={<span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{n.type}</span>}
              >
                <div
                  style={{
                    position: 'relative',
                    height: 440,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '12px 0',
                  }}
                >
                  <Image
                    src={nodeImage(nodeSection.node.type)}
                    alt={n.type}
                    fill
                    priority
                    sizes="(max-width: 900px) 100vw, 420px"
                    style={{ objectFit: 'contain' }}
                  />
                </div>
              </SectionCard>
            ) : (
              renderCharts(activeGroup ? chartsFor(activeGroup, kind) : [])
            )}
          </div>
        </div>
      )}
    </div>
  );
}

/** Software tab — node apps as a card grid (legacy NodeSoftwareTab). */
function NodeApps({
  apps,
  error,
  onCheck,
}: {
  apps: { name: string; status: string; currentVersion: string; desiredVersion: string; releaseDate: string }[];
  error: boolean;
  onCheck: () => void;
}) {
  if (error) {
    return (
      <div className="card">
        <EmptyState art="error" title="Couldn't load apps" sub="The software service didn't respond." />
      </div>
    );
  }
  if (apps.length === 0) {
    return (
      <div className="card">
        <EmptyState art="search" title="No apps" sub="This node isn't reporting any installed apps." />
      </div>
    );
  }
  return (
    <SectionCard
      title="Node apps"
      count={apps.length}
      right={
        <Button variant="outlined" size="small" startIcon={<SyncRounded />} onClick={onCheck}>
          Check for updates
        </Button>
      }
    >
      <div className="apps-grid">
        {apps.map((app) => {
          const update = app.status === 'update_available';
          return (
            <div key={app.name} className="app-card">
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 6 }}>
                <CheckCircleRounded sx={{ fontSize: 18, color: update ? 'var(--uk-ac)' : 'var(--uk-success)' }} />
                <span style={{ fontWeight: 600, textTransform: 'capitalize' }}>{app.name}</span>
              </div>
              <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}>
                Version: <span className="tnum">{app.currentVersion || '—'}</span>
              </div>
              {app.releaseDate && (
                <div style={{ fontSize: 12, color: 'var(--uk-ink-3)', marginTop: 2 }}>
                  Released {app.releaseDate}
                </div>
              )}
              <div style={{ marginTop: 10 }}>
                {update ? (
                  <span style={{ fontSize: 12.5, fontWeight: 600, color: 'var(--uk-ac-dark)' }}>
                    Update available{app.desiredVersion ? ` → ${app.desiredVersion}` : ''}
                  </span>
                ) : (
                  <span style={{ fontSize: 12.5, color: 'var(--uk-success)' }}>Up to date</span>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </SectionCard>
  );
}
