/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import TextField from '@mui/material/TextField';
import Switch from '@mui/material/Switch';

/**
 * Site detail — info / overview chart / map row + interactive site
 * components diagram + switch ports with live sparklines + type-to-confirm
 * restart (node-site-detail.jsx SiteDetail).
 */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import GroupRounded from '@mui/icons-material/GroupRounded';
import RestartAltRounded from '@mui/icons-material/RestartAltRounded';
import Skeleton from '@mui/material/Skeleton';

import { useNetworkSiteDetailQuery } from '@/client/graphql/site-detail.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import AppModal from '@/components/AppModal';
import { ComboChart, MiniSpark } from '@/components/charts';
import DetailPicker from '@/components/DetailPicker';
import { EmptyState } from '@/components/EmptyState';
import MapPanel from '@/components/Map/MapPanel';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { POLL_LIVE_MS, visiblePoll } from '@/lib/polling';
import { useUiPrefs } from '@/lib/store';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { toSite } from '@/lib/mappers/sites';
import { series } from '@/lib/series';
import { Ic } from '../../_components/icons';

/* deterministic series (module scope, mirrors prototype constants) */
const SITE_OVERVIEW = series(40, 14, 0.32, 0.12).map((b, i) => ({
  bar: b,
  line: series(44, 14, 0.16, 0.08)[i] ?? b,
}));
const SWITCH_OVERVIEW = series(36, 14, 0.3, 0.1).map((b, i) => ({
  bar: b,
  line: series(40, 14, 0.18, 0.06)[i] ?? b,
}));
const PORT_V = series(48, 16, 0.12, 0.04);
const PORT_C = series(30, 16, 0.22, 0.05);
const PORT_P = series(38, 16, 0.2, 0.06);

const SC_NODES = [
  { id: 'node', icon: 'router', label: 'Node', x: 260, y: 40 },
  { id: 'switch', icon: 'account_tree', label: 'Switch', x: 260, y: 150 },
  { id: 'charge', icon: 'bolt', label: 'Charge controller', x: 150, y: 272 },
  { id: 'back', icon: 'settings_input_antenna', label: 'Backhaul', x: 382, y: 272 },
  { id: 'solar', icon: 'light_mode', label: 'Solar panels', x: 92, y: 394 },
  { id: 'batt', icon: 'battery_charging_full', label: 'Batteries', x: 222, y: 394 },
] as const;
const SC_LINKS: [string, string][] = [
  ['node', 'switch'],
  ['switch', 'charge'],
  ['switch', 'back'],
  ['charge', 'solar'],
  ['charge', 'batt'],
];

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontSize: 12, color: 'var(--uk-ink-2)' }}>
      <span style={{ width: 9, height: 9, borderRadius: 3, background: color }} />
      {label}
    </span>
  );
}

function SiteDiagram({
  selected,
  onSelect,
}: {
  selected: string;
  onSelect: (id: string) => void;
}) {
  const at = (id: string) => SC_NODES.find((n) => n.id === id);
  return (
    <div className="sc-diag-wrap">
      <div className="sc-diag" style={{ width: 474, height: 440 }}>
        <svg width="474" height="440" style={{ position: 'absolute', inset: 0 }} aria-hidden="true">
          {SC_LINKS.map(([a, b], i) => {
            const A = at(a);
            const B = at(b);
            if (!A || !B) return null;
            return (
              <line
                key={i}
                x1={A.x}
                y1={A.y + 26}
                x2={B.x}
                y2={B.y - 26}
                stroke="var(--uk-line)"
                strokeWidth="2"
              />
            );
          })}
        </svg>
        {SC_NODES.map((nd) => (
          <button
            key={nd.id}
            type="button"
            className={`sc-tile${selected === nd.id ? ' on' : ''}`}
            onClick={() => onSelect(nd.id)}
            style={{ left: nd.x, top: nd.y }}
          >
            <span className="sc-ic">
              <Ic name={nd.icon} sx={{ fontSize: 26 }} />
            </span>
            <span className="sc-label">{nd.label}</span>
          </button>
        ))}
      </div>
    </div>
  );
}

function PortRow({
  idx,
  on,
  onToggle,
}: {
  idx: number;
  on: boolean;
  onToggle: () => void;
}) {
  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <span style={{ fontSize: 13.5, fontWeight: 600 }}>Port {idx}</span>
        <Switch
          checked={on}
          onChange={onToggle}
          inputProps={{ 'aria-label': `Port ${idx} power` }}
        />
      </div>
      {on && (
        <div style={{ marginTop: 10, display: 'flex', flexDirection: 'column', gap: 9 }}>
          {(
            [
              ['Voltage', PORT_V, 'var(--uk-ac)'],
              ['Current', PORT_C, 'var(--uk-secondary)'],
              ['Power', PORT_P, 'var(--uk-success-bright)'],
            ] as const
          ).map(([lbl, d, c]) => (
            <div key={lbl}>
              <div style={{ fontSize: 11.5, color: 'var(--uk-ink-3)', marginBottom: 2 }}>{lbl}</div>
              <MiniSpark data={d} accent={c} />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function SiteDetailScreen({ siteId }: { siteId: string }) {
  const router = useRouter();
  const toast = useToast();
  const [restart, setRestart] = useState(false);
  const [confirm, setConfirm] = useState('');
  const [selComp, setSelComp] = useState('switch');
  const [ports, setPorts] = useState<Record<number, boolean>>({ 1: true, 2: false, 3: false });

  const networkId = useUiPrefs((s) => s.networkId);

  const { data, loading, refetch } = useNetworkSiteDetailQuery({
    variables: { siteId },
    ...visiblePoll(POLL_LIVE_MS),
  });
  const view = data?.siteView;
  const siteSection = view?.site;
  const nodesSection = view?.nodes;

  // All sites in the network → the picker lets the user switch between them.
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

  const siteNodes = (nodesSection?.nodes ?? []).map((n) => toUkamaNode(n));
  const s = toSite(siteSection.site, {
    total: siteNodes.length,
    online: siteNodes.filter((n) => n.status === 'online').length,
  });
  const node = siteNodes[0];
  const installDate = siteSection.site.installDate || '—';
  const statusText =
    s.status === 'offline'
      ? 'is offline'
      : s.status === 'degraded'
        ? 'is online with warnings'
        : 'is online';

  return (
    <div className="page">
      <PageHeader
        crumb={['Sites', s.name]}
        title={s.name}
        onBack={() => router.push('/network/sites')}
        actions={
          <Button
            variant="contained"
            startIcon={<RestartAltRounded />}
            onClick={() => {
              setRestart(true);
              setConfirm('');
            }}
          >
            Restart site
          </Button>
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
        <span style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>{statusText}</span>
      </div>

      <div className="tile-grid site-top" style={{ marginBottom: 'var(--uk-gap)' }}>
        <SectionCard title="Site information">
          <div style={{ display: 'grid', gap: 13, marginTop: 2 }}>
            {(
              [
                ['Installed', installDate],
                ['Location', s.area],
                ['Nodes', `${siteNodes.length}`],
                ['Node', node ? node.serial : '—'],
              ] as const
            ).map(([k, v]) => (
              <div key={k}>
                <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{k}</div>
                <div className="tnum" style={{ fontSize: 13.5, fontWeight: 600, marginTop: 2 }}>
                  {v}
                </div>
              </div>
            ))}
          </div>
        </SectionCard>

        <SectionCard
          title="Site overview"
          right={
            <div style={{ display: 'flex', gap: 14 }}>
              <LegendDot color="var(--uk-ac)" label="Input power" />
              <LegendDot color="color-mix(in srgb, var(--uk-ac) 30%, transparent)" label="Storage" />
              <LegendDot color="var(--uk-secondary)" label="Consumption" />
            </div>
          }
        >
          <ComboChart data={SITE_OVERVIEW} height={180} />
        </SectionCard>

        <div className="card" style={{ padding: 0, overflow: 'hidden', position: 'relative', minHeight: 200 }}>
          <MapPanel sites={[s]} selected={s.id} compact />
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
            {/* per-site subscriber count: metrics-phase (siteView.kpis gap) */}
            —
          </div>
        </div>
      </div>

      <div className="sec-head" style={{ margin: '4px 0 12px' }}>
        <div className="sec-title">Site components</div>
      </div>
      <div className="tile-grid site-comp">
        <SectionCard>
          <SiteDiagram selected={selComp} onSelect={setSelComp} />
        </SectionCard>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--uk-gap)' }}>
          <SectionCard
            title="Switch overview"
            right={<span style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>Power consumption</span>}
          >
            <ComboChart data={SWITCH_OVERVIEW} height={150} accent="var(--uk-secondary)" />
          </SectionCard>
          <SectionCard title="Switch ports" right={<span className="cnt-pill tnum">8</span>}>
            <div style={{ display: 'flex', flexDirection: 'column' }}>
              {[1, 2, 3].map((p) => (
                <div
                  key={p}
                  style={{
                    borderBottom: p < 3 ? '1px solid var(--uk-line-soft)' : 'none',
                    padding: '12px 0',
                  }}
                >
                  <PortRow
                    idx={p}
                    on={!!ports[p]}
                    onToggle={() => setPorts((v) => ({ ...v, [p]: !v[p] }))}
                  />
                </div>
              ))}
            </div>
          </SectionCard>
        </div>
      </div>

      {restart && (
        <AppModal
          title="Restart site"
          width={460}
          onClose={() => setRestart(false)}
          footer={
            <>
              <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={() => setRestart(false)}>
                Cancel
              </Button>
              <Button
                variant="contained"
                disabled={confirm !== s.name}
                onClick={() => {
                  setRestart(false);
                  toast(`Restarting ${s.name}…`);
                }}
              >
                Restart
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
            Restarting this site will take it down for about 10 minutes. Type the site name{' '}
            <b style={{ color: 'var(--uk-ink)' }}>{s.name}</b> to confirm.
          </p>
          <TextField
            fullWidth
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder={s.name}
            autoFocus
          />
        </AppModal>
      )}
    </div>
  );
}
