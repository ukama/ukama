/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Support — find a customer (Business) or a site/node (Network) and resolve
 * the issue quickly. Master list + detail card (biz-ops.jsx BizSupport).
 */
import { useState } from 'react';
import Button from '@mui/material/Button';
import { EmptyState } from '@/components/EmptyState';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { Ic } from '@/app/(dashboard)/_components/icons';
import { NODES, SITES, SUBSCRIBERS } from '@/data';

type Kind = 'customer' | 'site' | 'node';

interface Result {
  kind: Kind;
  id: string;
  icon: string;
  title: string;
  sub: string;
  status: string;
  tiles: { label: string; value: React.ReactNode; sub?: string }[];
  issueTitle: string;
  issueText: string;
  actions: { label: string; bg: string }[];
}

function buildResults(mode: 'biz' | 'network'): Result[] {
  if (mode === 'network') {
    return [
      ...SITES.map((s) => ({
        kind: 'site' as const,
        id: 's_' + s.id,
        icon: 'location_on',
        title: s.name,
        sub: `${s.area} · ${s.subs} customers`,
        status: s.status,
        tiles: [
          { label: 'Customers', value: s.subs },
          { label: 'Uptime (30d)', value: s.uptime + '%' },
          { label: 'Battery', value: s.battery + '%' },
          { label: 'Signal', value: s.signal ? s.signal + ' dBm' : '—' },
        ],
        issueTitle: 'Status',
        issueText: s.issue ?? 'No active issues — site is operating normally.',
        actions: [
          { label: 'Copy summary', bg: 'var(--uk-ac)' },
          { label: 'Restart site', bg: 'var(--uk-secondary)' },
          { label: 'Escalate to Ukama', bg: '#2C3038' },
        ],
      })),
      ...NODES.map((n) => ({
        kind: 'node' as const,
        id: 'n_' + n.id,
        icon: 'router',
        title: n.serial,
        sub: `${n.type} · ${n.site}`,
        status: n.status,
        tiles: [
          { label: 'Site', value: <span style={{ fontSize: 15 }}>{n.site}</span> },
          { label: 'Temperature', value: n.temp != null ? n.temp + ' °C' : '—' },
          { label: 'Firmware', value: <span style={{ fontSize: 15 }}>{n.fw}</span> },
          { label: 'Uptime', value: n.up },
        ],
        issueTitle: 'Status',
        issueText: n.note
          ? n.note
          : n.status === 'offline'
            ? 'Node is offline — no telemetry received.'
            : 'Operating normally.',
        actions: [
          { label: 'Copy summary', bg: 'var(--uk-ac)' },
          { label: 'Restart node', bg: 'var(--uk-secondary)' },
          { label: 'Escalate to Ukama', bg: '#2C3038' },
        ],
      })),
    ];
  }
  return SUBSCRIBERS.map((s) => ({
    kind: 'customer' as const,
    id: s.id,
    icon: 'person',
    title: s.name,
    sub: `${s.phone} · ${s.plan}`,
    status: s.sim === 'suspended' ? 'pending' : s.sim,
    tiles: [
      {
        label: 'Package',
        value: s.plan === 'No plan' ? '—' : s.plan,
        sub: `${s.usage}${s.cap ? '/' + s.cap : ''} GB used`,
      },
      { label: 'Last seen', value: s.seen, sub: s.site },
    ],
    issueTitle: 'Likely issue',
    issueText:
      s.sim === 'suspended'
        ? 'SIM is suspended — billing or manual hold. Reactivate to restore service.'
        : s.plan === 'No plan'
          ? 'No active package. Assign a data plan to connect.'
          : `Active package; signal quality may be weak near ${s.site}.`,
    actions: [
      { label: 'Copy summary', bg: 'var(--uk-ac)' },
      { label: 'Escalate to Ukama', bg: '#2C3038' },
    ],
  }));
}

export default function SupportScreen({ mode }: { mode: 'biz' | 'network' }) {
  const network = mode === 'network';
  const toast = useToast();
  const [q, setQ] = useState('');
  const [selId, setSelId] = useState<string | null>(null);

  const results = buildResults(mode);
  const filtered = results.filter(
    (r) =>
      r.title.toLowerCase().includes(q.toLowerCase()) ||
      r.sub.toLowerCase().includes(q.toLowerCase()),
  );
  const cur = filtered.find((r) => r.id === selId) ?? filtered[0];

  return (
    <div className="page">
      <PageHeader
        title="Support"
        sub={
          network
            ? 'Find a site or node and resolve the issue quickly.'
            : 'Find a customer and resolve the issue quickly.'
        }
      />
      <div
        className="card card-pad"
        style={{ marginBottom: 'var(--uk-gap)', display: 'flex', gap: 12, alignItems: 'center' }}
      >
        <div style={{ flex: 1 }}>
          <SearchField
            value={q}
            onChange={setQ}
            width="100%"
            placeholder={network ? 'Search site or node' : 'Search customer by name or phone'}
          />
        </div>
        <Button variant="contained" sx={{ height: 38, px: 3.5 }}>
          Search
        </Button>
      </div>

      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1.6fr', alignItems: 'start' }}>
        <SectionCard title={network ? 'Sites & nodes' : 'Customers'}>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: 10,
              maxHeight: 560,
              overflowY: 'auto',
            }}
          >
            {filtered.length === 0 && (
              <div style={{ fontSize: 13, color: 'var(--uk-ink-3)', padding: '8px 2px' }}>
                No matches for “{q}”.
              </div>
            )}
            {filtered.map((r) => {
              const on = cur && r.id === cur.id;
              return (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => setSelId(r.id)}
                  style={{
                    display: 'flex',
                    gap: 12,
                    alignItems: 'center',
                    textAlign: 'left',
                    padding: '12px 14px',
                    borderRadius: 10,
                    cursor: 'pointer',
                    fontFamily: 'inherit',
                    border: `1px solid ${on ? 'var(--uk-ac)' : 'var(--uk-line)'}`,
                    background: on ? 'var(--uk-ac-soft)' : 'var(--uk-panel)',
                  }}
                >
                  <Ic
                    name={r.icon}
                    sx={{ fontSize: 21, color: on ? 'var(--uk-ac-dark)' : 'var(--uk-ink-3)', flex: 'none' }}
                  />
                  <div style={{ minWidth: 0, flex: 1 }}>
                    <div
                      style={{
                        fontSize: 13.5,
                        fontWeight: 600,
                        color: on ? 'var(--uk-ac-dark)' : 'var(--uk-ink)',
                      }}
                    >
                      {r.title}
                    </div>
                    <div
                      style={{
                        fontSize: 12.5,
                        color: 'var(--uk-ink-2)',
                        marginTop: 2,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {r.sub}
                    </div>
                  </div>
                  <StatusBadge status={r.status} />
                </button>
              );
            })}
          </div>
        </SectionCard>

        {!cur ? (
          <SectionCard title="Details">
            <EmptyState
              art="search"
              title="No match"
              sub="Search above to find something to support."
            />
          </SectionCard>
        ) : (
          <div className="card card-pad">
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 14 }}>
              <span style={{ fontFamily: 'var(--font-display)', fontSize: 21, fontWeight: 500 }}>
                {cur.title}
              </span>
              <StatusBadge status={cur.status} />
            </div>
            <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', marginBottom: 16 }}>
              {cur.tiles.map((t, i) => (
                <div key={i} style={{ border: '1px solid var(--uk-line)', borderRadius: 10, padding: 16 }}>
                  <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{t.label}</div>
                  <div
                    className="tnum"
                    style={{
                      fontFamily: 'var(--font-display)',
                      fontSize: 20,
                      fontWeight: 500,
                      margin: '4px 0 3px',
                    }}
                  >
                    {t.value}
                  </div>
                  {t.sub && <div style={{ fontSize: 12, color: 'var(--uk-ink-2)' }}>{t.sub}</div>}
                </div>
              ))}
            </div>
            <div
              style={{
                border: '1px solid var(--uk-line)',
                borderRadius: 10,
                padding: '16px 18px',
                marginBottom: 18,
              }}
            >
              <div className="sec-title" style={{ fontSize: 15, marginBottom: 7 }}>
                {cur.issueTitle}
              </div>
              <p
                style={{
                  fontSize: 13,
                  color: 'var(--uk-ink-2)',
                  lineHeight: 1.6,
                  margin: 0,
                  textWrap: 'pretty',
                }}
              >
                {cur.issueText}
              </p>
            </div>
            <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
              {cur.actions.map((a) => (
                <button
                  key={a.label}
                  type="button"
                  onClick={() => toast(`${a.label} — ${cur.title}`)}
                  style={{
                    flex: '1 1 auto',
                    minWidth: 140,
                    height: 42,
                    borderRadius: 8,
                    border: 'none',
                    cursor: 'pointer',
                    fontFamily: 'inherit',
                    fontSize: 13.5,
                    fontWeight: 600,
                    background: a.bg,
                    color: '#fff',
                  }}
                >
                  {a.label}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
