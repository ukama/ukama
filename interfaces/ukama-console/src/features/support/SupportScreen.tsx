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
 * the issue quickly. Master list + detail card. Wired to the live registry
 * composites: sitesView + nodesView (Network) and subscribersView (Business).
 * Metric-only tiles (uptime/battery/signal/temp/firmware) show "—" until the
 * metrics phase, mirroring the Sites/Nodes/Customers screens.
 */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';

import { Ic } from '@/app/(dashboard)/_components/icons';
import { useNetworkCustomersQuery } from '@/client/graphql/network-customers.generated';
import { useNodesListQuery } from '@/client/graphql/nodes-list.generated';
import { useSitesListQuery } from '@/client/graphql/sites-list.generated';
import { EmptyState } from '@/components/EmptyState';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import SectionCard from '@/components/SectionCard';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { Site, Subscriber, UkamaNode } from '@/data';
import { useAuth } from '@/lib/auth/context';
import { toUkamaNode } from '@/lib/mappers/nodes';
import { toSite } from '@/lib/mappers/sites';
import { toSubscriber } from '@/lib/mappers/subscribers';
import { useUiPrefs } from '@/lib/store';

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
  /** Subject + entity-specific body lines for the "Escalate to Ukama" email. */
  escalateSubject: string;
  escalateLines: string[];
}

const ESCALATE_LABEL = 'Escalate to Ukama';
const SUPPORT_EMAIL = 'hello@ukama.com';

/** Build a mailto: link to Ukama support with the context prefilled. */
const escalateMailto = (subject: string, lines: string[]): string => {
  const body = `Hi Ukama team,\n\nWe need help with the following:\n\n${lines.join(
    '\n',
  )}\n\nThanks,\n— sent from Ukama Console`;
  return `mailto:${SUPPORT_EMAIL}?subject=${encodeURIComponent(
    subject,
  )}&body=${encodeURIComponent(body)}`;
};

const ENTITY_ACTIONS = (extra: { label: string; bg: string }[]) => [
  { label: 'Copy summary', bg: 'var(--uk-ac)' },
  ...extra,
  { label: ESCALATE_LABEL, bg: '#2C3038' },
];

const idTile = (label: string, id: string) => ({
  label,
  value: (
    <span
      className="tnum"
      style={{ fontSize: 13, wordBreak: 'break-all' as const }}
    >
      {id}
    </span>
  ),
});

const siteResult = (s: Site): Result => ({
  kind: 'site',
  id: 's_' + s.id,
  icon: 'location_on',
  title: s.name,
  sub: `${s.area} · ${s.subs} customers`,
  status: s.status,
  tiles: [
    { label: 'Customers', value: s.subs },
    { label: 'Nodes', value: s.nodes },
    idTile('Site ID', s.id),
  ],
  issueTitle: 'Status',
  issueText:
    s.issue ??
    (s.status === 'offline'
      ? 'Site is offline — no telemetry received.'
      : 'No active issues — site is operating normally.'),
  actions: ENTITY_ACTIONS([
    { label: 'Restart site', bg: 'var(--uk-secondary)' },
  ]),
  escalateSubject: `Ukama support — site ${s.name}`,
  escalateLines: [
    `Site: ${s.name}`,
    `Site ID: ${s.id}`,
    `Status: ${s.status}`,
    `Customers: ${s.subs}`,
    `Nodes: ${s.nodes}`,
  ],
});

const nodeResult = (n: UkamaNode): Result => ({
  kind: 'node',
  id: 'n_' + n.id,
  icon: 'router',
  title: n.name ?? n.serial,
  sub: `${n.type} · ${n.site}`,
  status: n.status,
  tiles: [
    { label: 'Site', value: <span style={{ fontSize: 15 }}>{n.site}</span> },
    idTile('Node ID', n.id),
  ],
  issueTitle: 'Status',
  issueText: n.note
    ? n.note
    : n.status === 'offline'
      ? 'Node is offline — no telemetry received.'
      : 'Operating normally.',
  actions: ENTITY_ACTIONS([
    { label: 'Restart node', bg: 'var(--uk-secondary)' },
  ]),
  escalateSubject: `Ukama support — node ${n.name ?? n.id}`,
  escalateLines: [
    `Node: ${n.name ?? n.serial}`,
    `Node ID: ${n.id}`,
    `Type: ${n.type}`,
    `Site: ${n.site}`,
    `Status: ${n.status}`,
  ],
});

const customerResult = (s: Subscriber): Result => ({
  kind: 'customer',
  id: s.id,
  icon: 'person',
  title: s.name,
  sub: `${s.phone} · ${s.plan}`,
  status: s.sim === 'suspended' ? 'pending' : s.sim,
  tiles: [
    { label: 'Package', value: s.plan === 'No plan' ? '—' : s.plan },
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
    { label: ESCALATE_LABEL, bg: '#2C3038' },
  ],
  escalateSubject: `Ukama support — customer ${s.name}`,
  escalateLines: [
    `Customer: ${s.name}`,
    `Customer ID: ${s.id}`,
    `Email: ${s.email || '—'}`,
    `Phone: ${s.phone}`,
    `Plan: ${s.plan}`,
    `SIM status: ${s.sim}`,
    `ICCID: ${s.iccid}`,
    `Site: ${s.site}`,
  ],
});

export default function SupportScreen({ mode }: { mode: 'biz' | 'network' }) {
  const network = mode === 'network';
  const networkId = useUiPrefs((s) => s.networkId);
  const user = useAuth();
  const toast = useToast();
  const [q, setQ] = useState('');
  const [selId, setSelId] = useState<string | null>(null);

  // Network lens: sites + nodes. Business lens: customers. Unused queries skip.
  const { data: sitesData, loading: sitesLoading } = useSitesListQuery({
    variables: { networkId },
    skip: !network || !networkId,
  });
  const { data: nodesData, loading: nodesLoading } = useNodesListQuery({
    variables: { networkId },
    skip: !network || !networkId,
  });
  const { data: custData, loading: custLoading } = useNetworkCustomersQuery({
    variables: { networkId },
    skip: network || !networkId,
  });

  const loading = network ? sitesLoading || nodesLoading : custLoading;

  const results = useMemo<Result[]>(() => {
    if (network) {
      const siteNameById = new Map<string, string>();
      for (const s of sitesData?.sitesView.sites.sites ?? [])
        siteNameById.set(s.id, s.name);
      const countsBySite = new Map(
        (sitesData?.sitesView.nodeCounts.counts ?? []).map((c) => [
          c.siteId,
          { total: c.total, online: c.online },
        ]),
      );
      const customerCount = sitesData?.sitesView.customers.count ?? 0;
      const sites = (sitesData?.sitesView.sites.sites ?? []).map((s) =>
        siteResult(toSite(s, countsBySite.get(s.id), customerCount)),
      );
      const nodes = (nodesData?.nodesView.nodes.nodes ?? []).map((n) =>
        nodeResult(
          toUkamaNode(
            n,
            n.site?.siteId ? siteNameById.get(n.site.siteId) : undefined,
          ),
        ),
      );
      return [...sites, ...nodes];
    }
    const plansById = new Map(
      (custData?.subscribersView.plans.plans ?? []).map((p) => [
        p.packageId,
        p,
      ]),
    );
    return (custData?.subscribersView.subscribers.subscribers ?? []).map((s) =>
      customerResult(toSubscriber(s, plansById)),
    );
  }, [network, sitesData, nodesData, custData]);

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
        style={{
          marginBottom: 'var(--uk-gap)',
          display: 'flex',
          gap: 12,
          alignItems: 'center',
        }}
      >
        <div style={{ flex: 1 }}>
          <SearchField
            value={q}
            onChange={setQ}
            width="100%"
            placeholder={
              network
                ? 'Search site or node'
                : 'Search customer by name or phone'
            }
          />
        </div>
        <Button variant="contained" sx={{ height: 38, px: 3.5 }}>
          Search
        </Button>
      </div>

      <div
        className="tile-grid"
        style={{ gridTemplateColumns: '1fr 1.6fr', alignItems: 'start' }}
      >
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
            {loading && results.length === 0 && (
              <div
                style={{
                  fontSize: 13,
                  color: 'var(--uk-ink-3)',
                  padding: '8px 2px',
                }}
              >
                Loading…
              </div>
            )}
            {!loading && filtered.length === 0 && (
              <div
                style={{
                  fontSize: 13,
                  color: 'var(--uk-ink-3)',
                  padding: '8px 2px',
                }}
              >
                {results.length === 0
                  ? network
                    ? 'No sites or nodes yet.'
                    : 'No customers yet.'
                  : `No matches for “${q}”.`}
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
                    sx={{
                      fontSize: 21,
                      color: on ? 'var(--uk-ac-dark)' : 'var(--uk-ink-3)',
                      flex: 'none',
                    }}
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
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 12,
                marginBottom: 14,
              }}
            >
              <span
                style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 21,
                  fontWeight: 500,
                }}
              >
                {cur.title}
              </span>
              <StatusBadge status={cur.status} />
            </div>
            <div
              className="tile-grid"
              style={{ gridTemplateColumns: '1fr 1fr', marginBottom: 16 }}
            >
              {cur.tiles.map((t, i) => (
                <div
                  key={i}
                  style={{
                    border: '1px solid var(--uk-line)',
                    borderRadius: 10,
                    padding: 16,
                  }}
                >
                  <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
                    {t.label}
                  </div>
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
                  {t.sub && (
                    <div style={{ fontSize: 12, color: 'var(--uk-ink-2)' }}>
                      {t.sub}
                    </div>
                  )}
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
              <div
                className="sec-title"
                style={{ fontSize: 15, marginBottom: 7 }}
              >
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
                  onClick={() => {
                    const lines = [
                      `Org: ${user?.orgName || '—'}`,
                      `Network ID: ${networkId || '—'}`,
                      `Timestamp: ${new Date().toISOString()}`,
                      '',
                      ...cur.escalateLines,
                    ];
                    if (a.label === ESCALATE_LABEL) {
                      window.location.href = escalateMailto(
                        cur.escalateSubject,
                        lines,
                      );
                    } else if (a.label === 'Copy summary') {
                      void navigator.clipboard
                        .writeText(lines.join('\n'))
                        .then(() => toast('Summary copied to clipboard'))
                        .catch(() => toast('Could not copy summary'));
                    } else {
                      toast(`${a.label} — ${cur.title}`);
                    }
                  }}
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
