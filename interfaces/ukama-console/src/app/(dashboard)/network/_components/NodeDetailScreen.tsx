/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Node detail — per-type tabs, left KPI rail, node product imagery and the
 * "Turn node off" power menu (node-site-detail.jsx NodeDetail).
 */
import { useCallback, useEffect, useState } from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import RestartAltRounded from '@mui/icons-material/RestartAltRounded';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';
import SystemUpdateAltRounded from '@mui/icons-material/SystemUpdateAltRounded';
import ErrorRounded from '@mui/icons-material/ErrorRounded';
import HelpRounded from '@mui/icons-material/HelpRounded';
import CircularProgress from '@mui/material/CircularProgress';
import Skeleton from '@mui/material/Skeleton';

import { useNodeDetailQuery } from '@/client/graphql/node-detail.generated';
import type { MetricsRangeQuery } from '@/client/graphql/range-metrics.generated';
import { useMetricsRangeQuery } from '@/client/graphql/range-metrics.generated';
import { useRestartNodeMutation } from '@/client/graphql/controller.generated';
import {
  useGetAppsQuery,
  useUpdateSoftwareMutation,
} from '@/client/graphql/software.generated';
import AppModal from '@/components/AppModal';
import AppTabs from '@/components/AppTabs';
import MetricLineChart, {
  ChartMessage,
  thresholdLegendRows,
} from '@/components/MetricLineChart';
import DetailPicker from '@/components/DetailPicker';
import { EmptyState } from '@/components/EmptyState';
import KV from '@/components/KV';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { useToast } from '@/components/ToastProvider';
import RangeToggle from '@/components/RangeToggle';
import { metricLabel } from '@/lib/labels';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { formatDate } from '@/lib/parsers';
import { RANGE_SECONDS, type Range } from '@/lib/ranges';
import { ConnectivityDot, StateChip } from './nodeStatus';

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

/* -------------------------------------------------------------------------- *
 * Per-node-type metric key lists (which KPIs/graphs exist for a node type),
 * mirrored from console-bff getGraphsKeyByType. Labels, units, thresholds and
 * values all come from the BFF (nodeView.kpis / metricsRange) — the console
 * owns none of that; it only decides which keys to ask for.
 * -------------------------------------------------------------------------- */
type NodeKind = 'tnode' | 'anode' | 'cnode' | 'hnode';
const nodeKind = (raw: string): NodeKind => {
  const t = raw.toLowerCase();
  if (t.includes('anode')) return 'anode';
  if (t.includes('cnode')) return 'cnode';
  if (t.includes('hnode')) return 'hnode';
  return 'tnode';
};

type MetricGroup =
  | 'health'
  | 'customers'
  | 'cellular'
  | 'backhaul'
  | 'resources'
  | 'radio';
const GROUP_KEYS: Record<MetricGroup, Record<NodeKind, string[]>> = {
  // Uptime is shown as a value in Node information, not as a graph.
  health: {
    tnode: ['cpu_temperature', 'memory'],
    anode: ['fem1_temperature', 'fem2_temperature'],
    cnode: ['memory'],
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
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 6,
        fontSize: 12,
        color: 'var(--uk-ink-2)',
      }}
    >
      <span
        style={{ width: 9, height: 9, borderRadius: 3, background: color }}
      />
      {label}
    </span>
  );
}

/** One metric series straight from the BFF (metricsRange). */
type MetricSeries = MetricsRangeQuery['metricsRange']['metrics'][number];

/** Latest value + presentation metadata for one rail KPI, derived from a
 *  metricsRange series (its last real sample) — no separate latest fetch. */
type LatestEntry = {
  value: number;
  success: boolean;
  label?: string | null;
  unit?: string | null;
  format?: string | null;
};

/** Collapse a series to its latest KPI value: the most recent sample that
 *  isn't a gap-fill placeholder (-1). The chart's last point IS the latest. */
const seriesLatest = (m: MetricSeries): LatestEntry => {
  const vals = m.values ?? [];
  let last: number | null = null;
  for (let i = vals.length - 1; i >= 0; i--) {
    const v = vals[i]?.[1];
    if (v != null && v !== -1) {
      last = v;
      break;
    }
  }
  return {
    value: last ?? 0,
    success: m.success !== false && last != null,
    label: m.label,
    unit: m.unit,
    format: m.format,
  };
};

/** Hidden fetcher: one batched metricsRange (Day window) for rail KPIs whose
 *  chart isn't currently shown, so every left-rail value still has a source
 *  without a separate latest query. Reports each series' latest upward. */
function RailLatest({
  nodeId,
  keys,
  onLatest,
}: {
  nodeId: string;
  keys: string[];
  onLatest: (key: string, entry: LatestEntry) => void;
}) {
  const [nowSec] = useState(() => Math.floor(Date.now() / 1000));
  const { data } = useMetricsRangeQuery({
    variables: {
      data: { keys, nodeId, from: nowSec - RANGE_SECONDS.Day, to: nowSec },
    },
    skip: keys.length === 0,
  });
  useEffect(() => {
    for (const m of data?.metricsRange.metrics ?? []) {
      onLatest(m.type, seriesLatest(m));
    }
  }, [data, onLatest]);
  return null;
}

/** One metric chart with its own range filter — self-fetches its series so
 *  every graph filters independently. Rendered with Recharts. */
function MetricChart({
  nodeId,
  metricKey,
  off,
  onLatest,
}: {
  nodeId: string;
  metricKey: string;
  off: boolean;
  onLatest?: (key: string, entry: LatestEntry) => void;
}) {
  const [range, setRange] = useState<Range>('Day');
  const [nowSec] = useState(() => Math.floor(Date.now() / 1000));
  const to = nowSec;
  const from = nowSec - RANGE_SECONDS[range];
  const { data, loading, error } = useMetricsRangeQuery({
    variables: { data: { keys: [metricKey], nodeId, from, to } },
  });

  const m: MetricSeries | undefined = data?.metricsRange.metrics?.[0];
  // The chart already holds the series — feed its latest sample to the rail so
  // the same metric isn't fetched twice (chart + a separate latest query).
  useEffect(() => {
    if (m && onLatest) onLatest(metricKey, seriesLatest(m));
  }, [m, metricKey, onLatest]);
  const hasData = !!m && m.values.length > 0 && m.success !== false;
  const values: [number, number][] = hasData
    ? off
      ? m!.values.map((v) => [v[0] ?? 0, 0])
      : m!.values.map((v) => [v[0] ?? 0, v[1] ?? 0])
    : [];
  const legend = thresholdLegendRows(m?.threshold ?? null, m?.unit);
  const title = m?.label || metricKey;
  return (
    <SectionCard
      title={title}
      right={<RangeToggle value={range} onChange={setRange} />}
    >
      {error ? (
        <ChartMessage kind="error" message={error.message} height={300} />
      ) : loading && !m ? (
        <Skeleton variant="rounded" sx={{ height: 300 }} />
      ) : !hasData ? (
        <ChartMessage kind="empty" height={300} />
      ) : (
        <>
          <MetricLineChart
            values={values}
            title={title}
            unit={m?.unit}
            format={m?.format}
            threshold={m?.threshold ?? null}
            height={300}
          />
          <div
            style={{
              display: 'flex',
              gap: 18,
              justifyContent: 'center',
              marginTop: 10,
              flexWrap: 'wrap',
            }}
          >
            {legend.map((l) => (
              <LegendDot key={l.label} {...l} />
            ))}
          </div>
        </>
      )}
    </SectionCard>
  );
}

/** Right-panel charts for a metric group — one self-filtering chart per key. */
function GroupCharts({
  nodeId,
  keys,
  off,
  onLatest,
}: {
  nodeId: string;
  keys: string[];
  off: boolean;
  onLatest?: (key: string, entry: LatestEntry) => void;
}) {
  if (keys.length === 0) {
    return (
      <SectionCard title="Metrics">
        <EmptyState
          art="search"
          title="No metrics for this node type"
          sub="This node type doesn't report metrics in this category."
        />
      </SectionCard>
    );
  }
  return (
    <div
      style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}
    >
      {keys.map((k) => (
        <MetricChart
          key={k}
          nodeId={nodeId}
          metricKey={k}
          off={off}
          onLatest={onLatest}
        />
      ))}
    </div>
  );
}

/** Restart node — single action, confirmation dialog, then mutation. */
function RestartAction({ nodeId, name }: { nodeId: string; name: string }) {
  const [open, setOpen] = useState(false);
  const toast = useToast();
  const [restart, { loading }] = useRestartNodeMutation({
    onCompleted: (d) => {
      setOpen(false);
      toast(
        d.restartNode.success
          ? `Restarting ${name}…`
          : `Couldn't restart ${name}`,
      );
    },
    onError: () => {
      setOpen(false);
      toast(`Couldn't restart ${name}`);
    },
  });
  return (
    <>
      <Button
        variant="contained"
        startIcon={<RestartAltRounded />}
        sx={{ bgcolor: '#1C1E22', '&:hover': { bgcolor: '#2c2f36' } }}
        onClick={() => setOpen(true)}
      >
        Restart node
      </Button>
      {open && (
        <AppModal
          title="Restart node"
          width={440}
          onClose={() => {
            if (!loading) setOpen(false);
          }}
          footer={
            <>
              <Button
                color="inherit"
                sx={{ color: 'var(--uk-ink-3)' }}
                disabled={loading}
                onClick={() => setOpen(false)}
              >
                Cancel
              </Button>
              <Button
                variant="contained"
                startIcon={<RestartAltRounded />}
                disabled={loading}
                onClick={() => restart({ variables: { data: { nodeId } } })}
              >
                {loading ? 'Restarting…' : 'Restart node'}
              </Button>
            </>
          }
        >
          <div
            style={{ fontSize: 14, color: 'var(--uk-ink-2)', lineHeight: 1.55 }}
          >
            This will reboot <b style={{ color: 'var(--uk-ink)' }}>{name}</b>.
            The node will be briefly offline while it restarts, and active
            sessions may be interrupted.
          </div>
        </AppModal>
      )}
    </>
  );
}

export default function NodeDetailScreen({ nodeId }: { nodeId: string }) {
  const router = useRouter();
  const [tab, setTab] = useState('Overview');
  // Which left-rail section is selected (key from TAB_SECTIONS for the tab).
  const [section, setSection] = useState<string>('info');

  // Structural data — node core, site name and sibling nodes — as one
  // composite query (one-shot; no polling).
  const { data, loading, refetch } = useNodeDetailQuery({
    variables: { nodeId },
  });
  // Latest rail KPI values keyed by metric key. Single source of truth is the
  // metricsRange series (chart for charted keys, a batched Day fetch for the
  // rest) — no separate latest query, so a metric is never fetched twice.
  const [latestByKey, setLatestByKey] = useState<Record<string, LatestEntry>>(
    {},
  );
  const reportLatest = useCallback((key: string, entry: LatestEntry) => {
    setLatestByKey((prev) => {
      const cur = prev[key];
      if (cur && cur.value === entry.value && cur.success === entry.success) {
        return prev; // unchanged — skip to avoid a render loop
      }
      return { ...prev, [key]: entry };
    });
  }, []);

  const view = data?.nodeView;
  const nodeSection = view?.node;
  const softwareSection = view?.software;
  // Sibling nodes power the switcher dropdown.
  const pickerItems = (view?.siblings.nodes ?? []).map((nd) => ({
    id: nd.id,
    label: `${nd.name || nd.id} (${nd.id})`,
    status: '',
  }));
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
  // Render straight from series-derived metadata: label, unit, value. Never
  // shows a raw snake_case key — falls back to a humanized key.
  const labelFor = (key: string) => metricLabel(latestByKey[key]?.label, key);
  const fmtKpi = (key: string): string => {
    const e = latestByKey[key];
    if (!e || !e.success) return '—';
    const v = e.format === 'decimal' ? e.value.toFixed(2) : Math.round(e.value);
    const unit = e.unit ?? '';
    if (!unit) return `${v}`;
    return unit === '%' ? `${v}%` : `${v} ${unit}`;
  };
  // Uptime is reported in seconds — show a human-readable "Nd Nh" / "Nh Nm" /
  // "Nm" / "Ns". Sub-minute uptime renders as seconds (not a misleading "0m").
  const fmtUptime = (): string => {
    const e = latestByKey['uptime'];
    if (!e || !e.success) return '—';
    let s = Math.max(0, Math.round(e.value));
    const d = Math.floor(s / 86400);
    s -= d * 86400;
    const h = Math.floor(s / 3600);
    s -= h * 3600;
    const m = Math.floor(s / 60);
    s -= m * 60;
    if (d) return `${d}d ${h}h`;
    if (h) return `${h}h ${m}m`;
    if (m) return `${m}m`;
    return `${s}s`;
  };
  // KV rows for a metric group (one per node-type key, value or "—").
  const groupRows = (group: MetricGroup) => {
    const keys = groupKeys(group, kind);
    if (keys.length === 0) return <KV k="Metrics" v="—" />;
    return keys.map((k) => <KV key={k} k={labelFor(k)} v={fmtKpi(k)} />);
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

  // The left rail (and its selection) changes with the active tab. Drop any
  // section whose group has no KPIs for this node type (e.g. Customers on
  // controller/amplifier nodes); 'info' has no group and always shows.
  const sections = (TAB_SECTIONS[activeTab] ?? []).filter(
    (s) => !s.group || groupKeys(s.group, kind).length > 0,
  );
  const activeKey = sections.some((s) => s.key === section)
    ? section
    : (sections[0]?.key ?? 'info');
  const activeGroup = sections.find((s) => s.key === activeKey)?.group;

  // Charted keys (active group) get their latest from the chart series; every
  // other rail key (other cards + uptime) is fetched once by the batched
  // RailLatest below. Union by key so nothing is fetched twice.
  const chartedKeys = activeGroup ? groupKeys(activeGroup, kind) : [];
  const showsUptime = sections.some((s) => s.key === 'info');
  const railKeys = Array.from(
    new Set([
      ...sections.flatMap((s) => (s.group ? groupKeys(s.group, kind) : [])),
      ...(showsUptime ? ['uptime'] : []),
    ]),
  );
  const nonChartedKeys = railKeys.filter((k) => !chartedKeys.includes(k));

  return (
    <div className="page">
      <RailLatest
        nodeId={n.id}
        keys={nonChartedKeys}
        onLatest={reportLatest}
      />
      <PageHeader
        crumb={['Nodes', n.serial]}
        title={
          <span
            style={{ display: 'inline-flex', alignItems: 'center', gap: 10 }}
          >
            <ConnectivityDot connectivity={n.connectivity} />
            {nodeName}
          </span>
        }
        onBack={() => router.push('/network/nodes')}
        actions={<RestartAction nodeId={n.id} name={nodeName} />}
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

      <AppTabs
        tabs={visibleTabs}
        value={activeTab}
        onChange={setTab}
        scrollable
      />

      {activeTab === 'Software' ? (
        <NodeApps
          apps={apps}
          error={!!softwareSection?.error}
          nodeId={nodeId}
        />
      ) : (
        <div className="detail-grid">
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: 'var(--uk-gap)',
            }}
          >
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
                    <KV k={labelFor('uptime')} v={fmtUptime()} />
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
                right={
                  <span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                    {n.type}
                  </span>
                }
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
              <GroupCharts
                nodeId={n.id}
                keys={activeGroup ? groupKeys(activeGroup, kind) : []}
                off={off}
                onLatest={reportLatest}
              />
            )}
          </div>
        </div>
      )}
    </div>
  );
}

/** Software tab — node apps as a card grid (legacy NodeSoftwareTab). */
/** Presentable label + colour for each backend SoftwareStatusType. */
type SoftwareStatus =
  | 'up_to_date'
  | 'update_available'
  | 'update_in_progress'
  | 'update_failed'
  | 'unknown';

const SOFTWARE_STATUS_META: Record<
  SoftwareStatus,
  { label: string; color: string }
> = {
  up_to_date: { label: 'Up to date', color: 'var(--uk-success)' },
  update_available: { label: 'Update available', color: 'var(--uk-ac-dark)' },
  update_in_progress: { label: 'Updating…', color: 'var(--uk-ac-dark)' },
  update_failed: { label: 'Update failed', color: 'var(--uk-error)' },
  unknown: { label: 'Status unknown', color: 'var(--uk-ink-3)' },
};

const softwareStatusMeta = (status: string) =>
  SOFTWARE_STATUS_META[status as SoftwareStatus] ??
  SOFTWARE_STATUS_META.unknown;

function SoftwareStatusIcon({ status }: { status: string }) {
  const sx = { fontSize: 18, color: softwareStatusMeta(status).color };
  switch (status) {
    case 'up_to_date':
      return <CheckCircleRounded sx={sx} />;
    case 'update_available':
      return <SystemUpdateAltRounded sx={sx} />;
    case 'update_in_progress':
      return <CircularProgress size={15} sx={{ color: 'var(--uk-ac-dark)' }} />;
    case 'update_failed':
      return <ErrorRounded sx={sx} />;
    default:
      return <HelpRounded sx={sx} />;
  }
}

/** Humanize a byte count (e.g. 204800 → "200 KB"). */
const humanBytes = (n?: number | null): string => {
  if (n == null) return '—';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let v = n;
  let i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i += 1;
  }
  return `${v.toFixed(i > 0 && v < 100 ? 1 : 0)} ${units[i]}`;
};

function ResourceRow({ label, value }: { label: string; value: string }) {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        gap: 16,
        padding: '10px 0',
        borderBottom: '1px solid var(--uk-line-soft)',
      }}
    >
      <span style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>{label}</span>
      <span className="tnum" style={{ fontSize: 13, fontWeight: 600 }}>
        {value}
      </span>
    </div>
  );
}

/**
 * App detail dialog — fetches the node's runtime resource usage for a single
 * app (getApps with nodeId + appName) and renders CPU / memory / disk I/O.
 */
function AppResourceDialog({
  nodeId,
  appName,
  onClose,
}: {
  nodeId: string;
  appName: string;
  onClose: () => void;
}) {
  const { data, loading, error } = useGetAppsQuery({
    variables: { data: { nodeId, appName } },
    fetchPolicy: 'network-only',
  });
  const app = data?.getApps?.apps?.[0];
  const res = app?.resource;
  return (
    <AppModal
      title={appName}
      width={460}
      onClose={onClose}
      footer={
        <Button
          color="inherit"
          sx={{ color: 'var(--uk-ink-3)' }}
          onClick={onClose}
        >
          Close
        </Button>
      }
    >
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: 28 }}>
          <CircularProgress size={22} />
        </div>
      ) : error || !app ? (
        <div
          style={{ fontSize: 13.5, color: 'var(--uk-ink-3)', padding: '8px 0' }}
        >
          {error
            ? "Couldn't load app resources."
            : 'No resource data for this app.'}
        </div>
      ) : (
        <div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 10 }}>
            <span style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}>
              {app.version}
              {app.tag && app.tag !== app.version ? ` · ${app.tag}` : ''}
            </span>
            <span
              style={{
                fontSize: 12,
                fontWeight: 600,
                textTransform: 'capitalize',
                color:
                  app.status === 'running'
                    ? 'var(--uk-success)'
                    : 'var(--uk-ink-2)',
              }}
            >
              {app.status}
            </span>
          </div>
          <ResourceRow
            label="CPU"
            value={res ? `${res.cpuPercent.toFixed(1)}%` : '—'}
          />
          <ResourceRow
            label="Memory (RSS)"
            value={res ? humanBytes(res.memoryRssKb * 1024) : '—'}
          />
          <ResourceRow
            label="Disk read"
            value={res ? humanBytes(res.diskReadBytes) : '—'}
          />
          <ResourceRow
            label="Disk write"
            value={res ? humanBytes(res.diskWriteBytes) : '—'}
          />
        </div>
      )}
    </AppModal>
  );
}

function NodeApps({
  apps,
  error,
  nodeId,
}: {
  apps: {
    name: string;
    status: string;
    currentVersion: string;
    desiredVersion: string;
    releaseDate: string;
  }[];
  error: boolean;
  nodeId: string;
}) {
  const toast = useToast();
  const [updatingName, setUpdatingName] = useState<string | null>(null);
  const [selectedApp, setSelectedApp] = useState<string | null>(null);
  const [updateSoftware] = useUpdateSoftwareMutation({
    refetchQueries: ['NodeDetail'],
    awaitRefetchQueries: true,
  });

  const runUpdate = async (name: string, tag: string) => {
    setUpdatingName(name);
    try {
      await updateSoftware({ variables: { data: { name, nodeId, tag } } });
      toast(`Updating ${name} → ${tag}`);
    } catch {
      toast(`Couldn't update ${name}`);
    } finally {
      setUpdatingName(null);
    }
  };

  if (error) {
    return (
      <div className="card">
        <EmptyState
          art="error"
          title="Couldn't load apps"
          sub="The software service didn't respond."
        />
      </div>
    );
  }
  if (apps.length === 0) {
    return (
      <div className="card">
        <EmptyState
          art="search"
          title="No apps"
          sub="This node isn't reporting any installed apps."
        />
      </div>
    );
  }
  return (
    <SectionCard title="Node apps" count={apps.length}>
      <div className="apps-grid">
        {apps.map((app) => {
          const meta = softwareStatusMeta(app.status);
          // In-flight either optimistically (this card was just clicked) or per
          // the freshly-fetched backend status.
          const inProgress =
            app.status === 'update_in_progress' || updatingName === app.name;
          const canUpdate =
            !inProgress &&
            (app.status === 'update_available' ||
              app.status === 'update_failed');
          const showTarget =
            (app.status === 'update_available' ||
              app.status === 'update_failed') &&
            !!app.desiredVersion;
          return (
            <div
              key={app.name}
              className="app-card app-card-clickable"
              role="button"
              tabIndex={0}
              aria-label={`View ${app.name} resources`}
              onClick={() => setSelectedApp(app.name)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  setSelectedApp(app.name);
                }
              }}
            >
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 8,
                  marginBottom: 6,
                }}
              >
                <SoftwareStatusIcon
                  status={inProgress ? 'update_in_progress' : app.status}
                />
                <span style={{ fontWeight: 600, textTransform: 'capitalize' }}>
                  {app.name}
                </span>
              </div>
              <div style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}>
                Version:{' '}
                <span className="tnum">{app.currentVersion || '—'}</span>
              </div>
              {app.releaseDate && (
                <div
                  style={{
                    fontSize: 12,
                    color: 'var(--uk-ink-3)',
                    marginTop: 2,
                  }}
                >
                  Released {formatDate(app.releaseDate)}
                </div>
              )}
              <div
                style={{
                  marginTop: 10,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'flex-start',
                  gap: 8,
                }}
              >
                <span
                  style={{
                    fontSize: 12.5,
                    fontWeight: app.status === 'up_to_date' ? 400 : 600,
                    color: inProgress ? 'var(--uk-ac-dark)' : meta.color,
                  }}
                >
                  {inProgress ? 'Updating…' : meta.label}
                  {showTarget ? ` → ${app.desiredVersion}` : ''}
                </span>
                {canUpdate && (
                  <Button
                    size="small"
                    variant="contained"
                    disabled={updatingName !== null}
                    onClick={(e) => {
                      e.stopPropagation(); // don't open the resource dialog
                      runUpdate(app.name, app.desiredVersion);
                    }}
                  >
                    {app.status === 'update_failed'
                      ? 'Retry update'
                      : 'Update Now'}
                  </Button>
                )}
              </div>
            </div>
          );
        })}
      </div>
      {selectedApp && (
        <AppResourceDialog
          nodeId={nodeId}
          appName={selectedApp}
          onClose={() => setSelectedApp(null)}
        />
      )}
    </SectionCard>
  );
}
