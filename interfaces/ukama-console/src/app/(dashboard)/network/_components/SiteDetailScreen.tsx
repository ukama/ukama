/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Site detail — info / uptime overview / map row, an interactive site
 * components tree, and the selected component's power/health graph
 * (node-site-detail.jsx SiteDetail). Metrics come from the BFF (mocked until
 * the metric service lands); the console renders whatever it returns.
 */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Switch from '@mui/material/Switch';
import Divider from '@mui/material/Divider';
import TextField from '@mui/material/TextField';
import Skeleton from '@mui/material/Skeleton';
import GroupRounded from '@mui/icons-material/GroupRounded';
import RestartAltRounded from '@mui/icons-material/RestartAltRounded';
import KeyboardArrowDownRounded from '@mui/icons-material/KeyboardArrowDownRounded';
import SettingsRounded from '@mui/icons-material/SettingsRounded';

import {
  useRestartSiteMutation,
  useToggleRfStatusMutation,
  useToggleServiceMutation,
} from '@/client/graphql/controller.generated';
import { useNetworkSiteDetailQuery } from '@/client/graphql/site-detail.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { useMetricsRangeQuery } from '@/client/graphql/range-metrics.generated';
import AppModal from '@/components/AppModal';
import DetailPicker from '@/components/DetailPicker';
import { EmptyState } from '@/components/EmptyState';
import UkamaMap from '@/components/Map/UkamaMap';
import MetricLineChart, {
  ChartMessage,
  thresholdLegendRows,
} from '@/components/MetricLineChart';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { POLL_LIVE_MS, visiblePoll } from '@/lib/polling';
import { useUiPrefs } from '@/lib/store';
import { normalizeCoords } from '@/lib/geo';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { toSite } from '@/lib/mappers/sites';
import { Ic } from '../../_components/icons';

type Range = 'Day' | 'Week' | 'Month';
const RANGES: Range[] = ['Day', 'Week', 'Month'];
const RANGE_SECONDS: Record<Range, number> = {
  Day: 86_400,
  Week: 604_800,
  Month: 2_592_000,
};

/** ISO timestamp → "Jun 6, 2026"; passes through non-dates unchanged. */
const fmtDate = (raw?: string | null): string => {
  if (!raw) return '—';
  const d = new Date(raw);
  if (Number.isNaN(d.getTime())) return raw;
  return d.toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
};

interface CompDef {
  id: string;
  icon: string;
  label: string;
  /** Component elementType to read its name from siteView.components. */
  element?: string;
  /** Metric key driving the right-side graph (undefined = no graph). */
  metric?: string;
  /** Multiple metric keys → one chart per key (takes precedence over metric). */
  metrics?: string[];
}
const TREE: CompDef[][] = [
  [{ id: 'node', icon: 'router', label: 'Node' }],
  [{ id: 'switch', icon: 'account_tree', label: 'Switch', element: 'SWITCH' }],
  [
    {
      id: 'charge',
      icon: 'bolt',
      label: 'Charge controller',
      element: 'POWER',
      metric: 'controller_temperature',
    },
    {
      id: 'back',
      icon: 'settings_input_antenna',
      label: 'Backhaul',
      element: 'BACKHAUL',
      metric: 'backhaul_downlink',
    },
  ],
  [
    {
      id: 'solar',
      icon: 'light_mode',
      label: 'Solar panels',
      metric: 'solar_panel_power',
      metrics: [
        'solar_panel_power',
        'solar_panel_voltage',
        'solar_panel_current',
      ],
    },
    {
      id: 'batt',
      icon: 'battery_charging_full',
      label: 'Batteries',
      metric: 'battery_charge',
    },
  ],
];
const COMP_BY_ID = new Map(TREE.flat().map((c) => [c.id, c]));
const DEFAULT_COMP: CompDef = {
  id: 'batt',
  icon: 'battery_charging_full',
  label: 'Batteries',
  metric: 'battery_charge',
};

function RangeToggle({
  value,
  onChange,
}: {
  value: Range;
  onChange: (r: Range) => void;
}) {
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

/** Two rows of daily uptime bars from the site_uptime_percentage series. */
function UptimeBars({ values }: { values: number[] }) {
  const bars = (vals: number[]) => (
    <div className="uptime-row">
      {vals.map((v, i) => (
        <span
          key={i}
          className="uptime-bar"
          style={{
            background:
              v >= 95 ? 'var(--uk-success-bright)' : 'var(--uk-orange)',
            opacity: 0.6,
          }}
        />
      ))}
    </div>
  );
  const mid = Math.ceil(values.length / 2);
  // First row = most recent half (… → today); second = older half.
  const recent = values.slice(mid);
  const older = values.slice(0, mid);
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 16,
        marginTop: 12,
        flex: 1,
        minHeight: 0,
      }}
    >
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          minHeight: 0,
        }}
      >
        {bars(recent)}
        <div className="uptime-caption">
          <span>30 days ago</span>
          <span>Today</span>
        </div>
      </div>
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          minHeight: 0,
        }}
      >
        {bars(older)}
        <div className="uptime-caption">
          <span>60 days ago</span>
          <span>31 days ago</span>
        </div>
      </div>
    </div>
  );
}

function CompTile({
  comp,
  subtitle,
  active,
  onClick,
}: {
  comp: CompDef;
  subtitle: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      className={`comp-tile${active ? ' on' : ''}`}
      onClick={onClick}
    >
      <div className="comp-tile-head">
        <span
          className="sc-ic"
          style={{ width: 34, height: 34, borderRadius: 9 }}
        >
          <Ic name={comp.icon} sx={{ fontSize: 18 }} />
        </span>
      </div>
      <div className="comp-tile-label">{comp.label}</div>
      {subtitle && <div className="comp-tile-sub">{subtitle}</div>}
    </button>
  );
}

/** Right-side graph for the selected component (its metric, range-filtered).
 *  Component metrics are cnode-scoped, so the chart queries with the site's
 *  controller (cnode) id — the gateway resolves the node type from it. */
const COMP_CHART_HEIGHT = 300;
function ComponentChart({
  metricKey,
  fallbackLabel,
  cnodeId,
  titleOverride,
}: {
  metricKey: string;
  fallbackLabel: string;
  cnodeId: string | null;
  /** Force the card title (ignores the series label) — e.g. "Speed (MBPS)". */
  titleOverride?: string;
}) {
  const [range, setRange] = useState<Range>('Day');
  const [nowSec] = useState(() => Math.floor(Date.now() / 1000));
  const to = nowSec;
  const from = nowSec - RANGE_SECONDS[range];
  const { data, loading, error } = useMetricsRangeQuery({
    variables: {
      data: {
        keys: [metricKey],
        from,
        to,
        ...(cnodeId ? { nodeId: cnodeId } : {}),
      },
    },
  });
  const m = data?.metricsRange.metrics?.[0];
  const hasData = !!m && m.values.length > 0 && m.success !== false;
  const title = titleOverride ?? (m?.label || fallbackLabel);

  const values: [number, number][] = hasData
    ? m!.values.map((v) => [v[0] ?? 0, v[1] ?? 0])
    : [];
  const legend = thresholdLegendRows(m?.threshold ?? null, m?.unit);
  return (
    <SectionCard
      title={title}
      right={<RangeToggle value={range} onChange={setRange} />}
    >
      {error ? (
        <ChartMessage
          kind="error"
          message={error.message}
          height={COMP_CHART_HEIGHT}
        />
      ) : loading && !m ? (
        <Skeleton variant="rounded" sx={{ height: COMP_CHART_HEIGHT }} />
      ) : !hasData ? (
        <ChartMessage kind="empty" height={COMP_CHART_HEIGHT} />
      ) : (
        <>
          <MetricLineChart
            values={values}
            title={title}
            unit={m?.unit}
            format={m?.format}
            threshold={m?.threshold ?? null}
            height={COMP_CHART_HEIGHT}
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

/** Right panel for the selected component: a no-metric notice, one filling
 *  chart for a single metric, or a stack of charts when it has several. */
function ComponentPanel({
  comp,
  cnodeId,
}: {
  comp: CompDef;
  cnodeId: string | null;
}) {
  const keys = comp.metrics ?? (comp.metric ? [comp.metric] : []);
  if (keys.length === 0) {
    return (
      <SectionCard
        title={comp.label}
        bodyStyle={{
          minHeight: COMP_CHART_HEIGHT,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <EmptyState
          art="search"
          title="No metric for this component"
          sub="This component doesn't report a time-series metric yet."
        />
      </SectionCard>
    );
  }
  return (
    <div
      style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}
    >
      {keys.map((k) => (
        <ComponentChart
          key={k}
          metricKey={k}
          fallbackLabel={comp.label}
          cnodeId={cnodeId}
        />
      ))}
    </div>
  );
}

/** Switch ports shown on the site page — each with live Speed/Power KPIs.
 *  Port n maps to switch_port_n_{speed,power} (cnode series). */
const SWITCH_PORTS: { n: number; name: string }[] = [
  { n: 1, name: 'Tnode PoE' },
  { n: 2, name: 'Cnode PoE' },
  { n: 3, name: 'Anode PoE' },
  { n: 9, name: 'Uplink SFP' },
];

/** One expandable port row: header + reveal of its Speed and Power charts. */
function SwitchPortRow({
  port,
  cnodeId,
}: {
  port: { n: number; name: string };
  cnodeId: string | null;
}) {
  const [open, setOpen] = useState(false);
  const [enabled, setEnabled] = useState(true);
  return (
    <div style={{ borderTop: '1px solid var(--uk-line)', padding: '14px 0' }}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 12,
        }}
      >
        <div style={{ fontWeight: 600 }}>
          Port {port.n}: {port.name}
        </div>
        <label
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            cursor: 'pointer',
          }}
        >
          <ToggleState on={enabled} />
          <Switch
            edge="end"
            checked={enabled}
            onChange={(e) => setEnabled(e.target.checked)}
          />
        </label>
      </div>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        style={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: 4,
          marginTop: 8,
          background: 'none',
          border: 'none',
          padding: 0,
          cursor: 'pointer',
          color: 'var(--uk-ac)',
          fontWeight: 600,
          fontSize: 12.5,
          letterSpacing: 0.4,
        }}
      >
        {open ? 'VIEW LESS' : 'VIEW MORE'}
        <KeyboardArrowDownRounded
          sx={{
            fontSize: 18,
            transition: 'transform 0.15s',
            transform: open ? 'rotate(180deg)' : 'none',
          }}
        />
      </button>
      {open &&
        (enabled ? (
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: 'var(--uk-gap)',
              marginTop: 12,
            }}
          >
            <ComponentChart
              metricKey={`switch_port_${port.n}_speed`}
              fallbackLabel="Speed (MBPS)"
              titleOverride="Speed (MBPS)"
              cnodeId={cnodeId}
            />
            <ComponentChart
              metricKey={`switch_port_${port.n}_power`}
              fallbackLabel="Power (watts)"
              titleOverride="Power (watts)"
              cnodeId={cnodeId}
            />
          </div>
        ) : (
          <div style={{ marginTop: 12 }}>
            <EmptyState
              art="search"
              title="Port is off"
              sub="Turn this port on to view its speed and power metrics."
            />
          </div>
        ))}
    </div>
  );
}

/** Right panel when the Switch component is selected: its ports + KPIs. */
function SwitchPortsPanel({ cnodeId }: { cnodeId: string | null }) {
  return (
    <SectionCard title={`Switch ports (${SWITCH_PORTS.length})`}>
      {SWITCH_PORTS.map((p) => (
        <SwitchPortRow key={p.n} port={p} cnodeId={cnodeId} />
      ))}
    </SectionCard>
  );
}

type SiteNode = ReturnType<typeof toUkamaNode>;
const STATUS_PHRASE: Record<string, string> = {
  online: 'is online and well',
  configuring: 'is configured',
  degraded: 'is online with warnings',
  offline: 'is offline',
};
const statusColor = (st: string) =>
  st === 'offline'
    ? 'var(--uk-error)'
    : st === 'degraded'
      ? 'var(--uk-orange)'
      : 'var(--uk-success)';

/** Node component selected → list the site's nodes as cards. */
function SiteNodesPanel({
  nodes,
  onOpen,
}: {
  nodes: SiteNode[];
  onOpen: (id: string) => void;
}) {
  return (
    <SectionCard
      title="Nodes"
      count={nodes.length}
      style={{ display: 'flex', flexDirection: 'column', height: '100%' }}
      bodyStyle={{ flex: 1, minHeight: 0 }}
    >
      {nodes.length === 0 ? (
        <EmptyState
          art="node"
          title="No nodes"
          sub="This site has no nodes installed."
        />
      ) : (
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            gap: 'var(--uk-gap)',
          }}
        >
          {nodes.map((n) => (
            <div key={n.id} className="app-card">
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <span
                  style={{
                    width: 9,
                    height: 9,
                    borderRadius: '50%',
                    background: statusColor(n.status),
                    flex: 'none',
                  }}
                />
                <span style={{ fontWeight: 600 }}>{n.name ?? n.serial}</span>
                <span style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
                  {STATUS_PHRASE[n.status] ?? ''}
                </span>
              </div>
              <div style={{ marginTop: 10 }}>
                <Button
                  size="small"
                  sx={{ p: 0, minWidth: 0 }}
                  onClick={() => onOpen(n.id)}
                >
                  View node
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </SectionCard>
  );
}

/** Small On/Off state label shown next to a toggle switch. */
function ToggleState({ on }: { on: boolean }) {
  return (
    <span
      style={{
        fontSize: 12,
        fontWeight: 600,
        minWidth: 24,
        textAlign: 'right',
        marginRight: 8,
        color: on ? 'var(--uk-success)' : 'var(--uk-ink-3)',
      }}
    >
      {on ? 'On' : 'Off'}
    </span>
  );
}

/**
 * Site actions dropdown — restart the site, plus RF / service radio toggles.
 * Restart is site-scoped; the RF/service toggles act on the site's tower node
 * (the controller maps RF to its amplifier internally), so they're disabled
 * when the site has no reachable tower node.
 */
function SiteActions({
  siteId,
  networkId,
  siteName,
  tnodeId,
}: {
  siteId: string;
  networkId: string;
  siteName: string;
  tnodeId: string | null;
}) {
  const toast = useToast();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [restart, setRestart] = useState(false);
  const [confirm, setConfirm] = useState('');
  const [rfOn, setRfOn] = useState(true);
  const [serviceOn, setServiceOn] = useState(true);
  const open = Boolean(anchorEl);

  const [restartSite, { loading: restarting }] = useRestartSiteMutation({
    onCompleted: (d) => {
      setRestart(false);
      toast(
        d.restartSite.success
          ? `Restarting ${siteName}…`
          : `Couldn't restart ${siteName}`,
      );
    },
    onError: () => {
      setRestart(false);
      toast(`Couldn't restart ${siteName}`);
    },
  });

  const [toggleRF, { loading: rfLoading }] = useToggleRfStatusMutation({
    fetchPolicy: 'network-only',
  });
  const [toggleService, { loading: serviceLoading }] = useToggleServiceMutation(
    {
      fetchPolicy: 'network-only',
    },
  );

  const onToggleRf = async () => {
    if (!tnodeId) return;
    const next = !rfOn;
    setRfOn(next); // optimistic
    try {
      await toggleRF({
        variables: { data: { nodeId: tnodeId, status: next } },
      });
      toast(`RF turned ${next ? 'on' : 'off'}`);
    } catch {
      setRfOn(!next); // revert
      toast(`Couldn't turn RF ${next ? 'on' : 'off'}`);
    }
  };

  const onToggleService = async () => {
    if (!tnodeId) return;
    const next = !serviceOn;
    setServiceOn(next); // optimistic
    try {
      await toggleService({
        variables: { data: { nodeId: tnodeId, status: next } },
      });
      toast(`Service turned ${next ? 'on' : 'off'}`);
    } catch {
      setServiceOn(!next); // revert
      toast(`Couldn't turn service ${next ? 'on' : 'off'}`);
    }
  };

  const togglesDisabled = !tnodeId || rfLoading || serviceLoading;

  return (
    <>
      <Button
        variant="contained"
        startIcon={<SettingsRounded />}
        endIcon={<KeyboardArrowDownRounded />}
        onClick={(e) => setAnchorEl(e.currentTarget)}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
      >
        Site actions
      </Button>
      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
        slotProps={{ paper: { sx: { minWidth: 248 } } }}
      >
        <MenuItem
          onClick={() => {
            setAnchorEl(null);
            setConfirm('');
            setRestart(true);
          }}
        >
          <ListItemIcon>
            <RestartAltRounded fontSize="small" />
          </ListItemIcon>
          <ListItemText>Restart site</ListItemText>
        </MenuItem>
        <Divider />
        <MenuItem disabled={togglesDisabled} onClick={onToggleRf}>
          <ListItemText
            primary="RF"
            secondary={tnodeId ? undefined : 'No tower node on this site'}
          />
          <ToggleState on={rfOn} />
          <Switch
            edge="end"
            checked={rfOn}
            disabled={togglesDisabled}
            tabIndex={-1}
          />
        </MenuItem>
        <MenuItem disabled={togglesDisabled} onClick={onToggleService}>
          <ListItemText
            primary="Service"
            secondary={tnodeId ? undefined : 'No tower node on this site'}
          />
          <ToggleState on={serviceOn} />
          <Switch
            edge="end"
            checked={serviceOn}
            disabled={togglesDisabled}
            tabIndex={-1}
          />
        </MenuItem>
      </Menu>

      {restart && (
        <AppModal
          title="Restart site"
          width={460}
          onClose={() => {
            if (!restarting) setRestart(false);
          }}
          footer={
            <>
              <Button
                color="inherit"
                sx={{ color: 'var(--uk-ink-3)' }}
                disabled={restarting}
                onClick={() => setRestart(false)}
              >
                Cancel
              </Button>
              <Button
                variant="contained"
                disabled={confirm !== siteName || restarting}
                onClick={() =>
                  restartSite({ variables: { data: { siteId, networkId } } })
                }
              >
                {restarting ? 'Restarting…' : 'Restart'}
              </Button>
            </>
          }
        >
          <p
            style={{
              fontSize: 13.5,
              color: 'var(--uk-ink-2)',
              lineHeight: 1.6,
              margin: '0 0 16px',
              textWrap: 'pretty',
            }}
          >
            Restarting this site will take it down for about 10 minutes. Type
            the site name <b style={{ color: 'var(--uk-ink)' }}>{siteName}</b>{' '}
            to confirm.
          </p>
          <TextField
            fullWidth
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder={siteName}
            autoFocus
          />
        </AppModal>
      )}
    </>
  );
}

export default function SiteDetailScreen({ siteId }: { siteId: string }) {
  const router = useRouter();
  const [selComp, setSelComp] = useState('node');

  const networkId = useUiPrefs((s) => s.networkId);

  const { data, loading, refetch } = useNetworkSiteDetailQuery({
    variables: { siteId },
    ...visiblePoll(POLL_LIVE_MS),
  });
  const view = data?.siteView;
  const siteSection = view?.site;
  const nodesSection = view?.nodes;
  const components = view?.components.components ?? [];

  // 90-day daily uptime series for the Site overview card.
  const [uNow] = useState(() => Math.floor(Date.now() / 1000));
  const { data: uptimeData, loading: uptimeLoading } = useMetricsRangeQuery({
    variables: {
      data: {
        keys: ['site_uptime_percentage'],
        from: uNow - 90 * 86_400,
        to: uNow,
      },
    },
  });
  const uptimeVals = (uptimeData?.metricsRange.metrics?.[0]?.values ?? []).map(
    (v) => v[1] ?? 0,
  );
  const uptimePct = uptimeVals.length
    ? Math.round(uptimeVals.reduce((a, b) => a + b, 0) / uptimeVals.length)
    : null;

  const { data: sitesData } = useSitesListQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const pickerItems = (sitesData?.sitesView.sites.sites ?? []).map((si) => ({
    id: si.id,
    label: si.name,
    status: '',
  }));

  if (loading) {
    return (
      <div className="page">
        <PageHeader crumb={['Sites', siteId]} title="Site" />
        <Skeleton variant="rounded" sx={{ height: 42, mb: 2 }} />
        <Skeleton variant="rounded" sx={{ height: 420 }} />
      </div>
    );
  }
  if (!siteSection?.site) {
    return (
      <div className="page">
        <PageHeader crumb={['Sites', siteId]} title="Site" />
        <div className="card">
          <EmptyState
            art="error"
            title="Couldn't load site"
            sub={siteSection?.error?.message ?? 'Site not found.'}
            cta="Try again"
            onCta={() => refetch()}
          />
        </div>
      </div>
    );
  }

  // Order nodes by type: tower → amplifier → controller → home.
  const typeRank = (id: string) =>
    ['tnode', 'anode', 'cnode', 'hnode'].findIndex((t) =>
      id.toLowerCase().includes(t),
    );
  const siteNodes = (nodesSection?.nodes ?? [])
    .map((n) => toUkamaNode(n))
    .sort((a, b) => {
      const ra = typeRank(a.id);
      const rb = typeRank(b.id);
      return (ra < 0 ? 99 : ra) - (rb < 0 ? 99 : rb);
    });
  // RF / service toggles act on the site's tower node (the controller maps RF
  // to the amplifier internally); null when the site has no tower node.
  const tnodeId =
    siteNodes.find((n) => n.id.toLowerCase().includes('tnode'))?.id ?? null;
  // Component metrics (power/solar/battery/controller/backhaul) are reported by
  // the site's controller node — query them with its cnode id.
  const cnodeId =
    siteNodes.find((n) => n.id.toLowerCase().includes('cnode'))?.id ?? null;
  // A node counts as "up" when it's reachable (connectivity online) — a
  // configured/operational node, not just status === 'online'. Mirrors the
  // BFF's connectivity-based site node counts.
  const s = toSite(siteSection.site, {
    total: siteNodes.length,
    online: siteNodes.filter(
      (n) => (n.connectivity ?? '').toLowerCase() === 'online',
    ).length,
  });
  const dto = siteSection.site;
  const installDate = fmtDate(dto.installDate || dto.createdAt);
  const geo = normalizeCoords(dto.latitude, dto.longitude);
  const coords = geo ? `${geo.lat}, ${geo.lng}` : null;
  const mapMarkers = geo
    ? [{ id: s.id, lat: geo.lat, lng: geo.lng, color: statusColor(s.status) }]
    : [];
  const statusText =
    s.status === 'offline'
      ? 'is offline'
      : s.status === 'degraded'
        ? 'is online with warnings'
        : 'is online';

  const compName = (element?: string) =>
    element
      ? (components.find((c) => c.elementType === element)?.componentName ??
        null)
      : null;
  const subtitleFor = (c: CompDef): string => {
    const name = compName(c.element);
    if (name) return name;
    if (c.id === 'node') return '';
    return c.label;
  };

  const selected = COMP_BY_ID.get(selComp) ?? DEFAULT_COMP;

  return (
    <div className="page">
      <PageHeader
        crumb={['Sites', s.name]}
        title={s.name}
        onBack={() => router.push('/network/sites')}
        actions={
          <SiteActions
            siteId={s.id}
            networkId={networkId}
            siteName={s.name}
            tnodeId={tnodeId}
          />
        }
      />

      <div className="detail-subrow">
        <DetailPicker
          value={{ id: s.id, label: s.name, status: s.status }}
          items={
            pickerItems.length > 0
              ? pickerItems
              : [{ id: s.id, label: s.name, status: s.status }]
          }
          onPick={(it) => router.push(`/network/sites/${it.id}`)}
        />
        <StatusBadge status={s.status} />
        <span style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>
          {statusText}
        </span>
      </div>

      <div
        className="tile-grid site-top"
        style={{ marginBottom: 'var(--uk-gap)' }}
      >
        <SectionCard title="Site information">
          <div style={{ display: 'grid', gap: 14, marginTop: 2 }}>
            <div>
              <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                Nodes
              </div>
              {siteNodes.length > 0 ? (
                siteNodes.map((n) => (
                  <div key={n.id} style={{ marginTop: 4 }}>
                    <span
                      className="tnum"
                      style={{ fontSize: 13.5, fontWeight: 600 }}
                    >
                      {n.serial}
                    </span>
                    <span
                      style={{
                        fontSize: 12,
                        color: 'var(--uk-ink-3)',
                        marginLeft: 6,
                      }}
                    >
                      · {n.type}
                    </span>
                  </div>
                ))
              ) : (
                <div style={{ fontSize: 13.5, fontWeight: 600, marginTop: 2 }}>
                  —
                </div>
              )}
            </div>
            <div>
              <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                Date created
              </div>
              <div
                className="tnum"
                style={{ fontSize: 13.5, fontWeight: 600, marginTop: 2 }}
              >
                {installDate}
              </div>
            </div>
            <div>
              <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                Location
              </div>
              <div style={{ fontSize: 13.5, fontWeight: 600, marginTop: 2 }}>
                {s.area}
              </div>
              {coords && (
                <div
                  className="tnum"
                  style={{
                    fontSize: 12,
                    color: 'var(--uk-ink-3)',
                    marginTop: 2,
                  }}
                >
                  {coords}
                </div>
              )}
            </div>
          </div>
        </SectionCard>

        <SectionCard
          title="Site overview"
          style={{ display: 'flex', flexDirection: 'column' }}
          bodyStyle={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            minHeight: 0,
          }}
        >
          <div style={{ display: 'flex', alignItems: 'baseline', gap: 8 }}>
            <span
              style={{
                fontSize: 30,
                fontWeight: 600,
                fontFamily: 'var(--font-display)',
              }}
            >
              {uptimePct != null ? `${uptimePct}%` : '—'}
            </span>
            <span style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>
              uptime over 90 days
            </span>
          </div>
          {uptimeLoading && uptimeVals.length === 0 ? (
            <Skeleton variant="rounded" sx={{ height: 88, mt: 1 }} />
          ) : uptimeVals.length === 0 ? (
            <div
              style={{ fontSize: 13, color: 'var(--uk-ink-3)', marginTop: 8 }}
            >
              No uptime data available.
            </div>
          ) : (
            <UptimeBars values={uptimeVals} />
          )}
        </SectionCard>

        <div
          className="card"
          style={{
            padding: 0,
            overflow: 'hidden',
            position: 'relative',
            minHeight: 200,
            display: 'flex',
          }}
        >
          <div style={{ flex: 1, minHeight: 0 }}>
            <UkamaMap
              markers={mapMarkers}
              zoom={12}
              fitToMarkers={false}
              height="100%"
            />
          </div>
          <div
            style={{
              position: 'absolute',
              left: 12,
              bottom: 12,
              background: 'var(--uk-panel)',
              borderRadius: 8,
              padding: '5px 10px',
              display: 'flex',
              alignItems: 'center',
              gap: 7,
              boxShadow: 'var(--uk-shadow)',
              fontSize: 12.5,
              fontWeight: 600,
            }}
          >
            <GroupRounded sx={{ fontSize: 16, color: 'var(--uk-ac)' }} />
            {s.subs || '—'}
          </div>
        </div>
      </div>

      <div className="sec-head" style={{ margin: '4px 0 12px' }}>
        <div className="sec-title">Site components</div>
      </div>
      <div className="tile-grid site-comp">
        <SectionCard>
          <div className="comp-tree">
            {TREE.map((level, li) => (
              <div key={li} style={{ display: 'contents' }}>
                {li > 0 && <div className="comp-conn" />}
                <div className="comp-level">
                  {level.map((c) => (
                    <CompTile
                      key={c.id}
                      comp={c}
                      subtitle={subtitleFor(c)}
                      active={selComp === c.id}
                      onClick={() => setSelComp(c.id)}
                    />
                  ))}
                </div>
              </div>
            ))}
          </div>
        </SectionCard>
        {/* Right column scrolls within the row so a tall multi-chart panel
            (e.g. solar power/voltage/current) doesn't stretch the left tree. */}
        <div
          style={{
            maxHeight: 'calc(100vh - 220px)',
            overflowY: 'auto',
            minHeight: 0,
          }}
        >
          {selected.id === 'node' ? (
            <SiteNodesPanel
              nodes={siteNodes}
              onOpen={(id) => router.push(`/network/nodes/${id}`)}
            />
          ) : selected.id === 'switch' ? (
            <SwitchPortsPanel cnodeId={cnodeId} />
          ) : (
            <ComponentPanel comp={selected} cnodeId={cnodeId} />
          )}
        </div>
      </div>
    </div>
  );
}
