/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import Meter from '@/components/Meter';

/**
 * THE canonical shared record component (BUILD-PLAN §2): one customers
 * table serving all three lenses through a `mode` prop — never forked.
 *  - biz:     read-only rows → drawer (readOnly)
 *  - network: read-only list, shows site column, no row click
 *  - agent:   multi-select + bulk bar, kebab actions, add CTA, drawer
 */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import type { ColumnDef } from '@tanstack/react-table';
import DataTable from '@/components/data-table/DataTable';
import TableFooter from '@/components/data-table/TableFooter';
import DateChip from '@/components/DateChip';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import StatusBadge from '@/components/StatusBadge';
import { useNetworkCustomersQuery } from '@/client/graphql/network-customers.generated';
import type { Subscriber } from '@/data';
import { parseSeen } from '@/lib/parsers';
import { useUiPrefs } from '@/lib/store';
import { toSubscriber } from '@/lib/mappers/subscribers';
import AddCustomerDialog from './AddCustomerDialog';
import SubscriberDrawer from './SubscriberDrawer';

export type CustomersMode = 'biz' | 'network' | 'agent';

const SUBS = {
  biz: 'Who are my customers and what state are they in?',
  network: 'Everyone connected to your network.',
  agent: 'Manage your customers’ packages and top-ups.',
} as const;


export default function SubscribersScreen({ mode }: { mode: CustomersMode }) {
  const agent = mode === 'agent';
  const showSite = mode === 'network';
  const clickRow = mode !== 'network';
  const networkId = useUiPrefs((s) => s.networkId);

  const [q, setQ] = useState('');
  const [openSub, setOpenSub] = useState<Subscriber | null>(null);
  const [showAdd, setShowAdd] = useState(false);

  const { data, loading, refetch } = useNetworkCustomersQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const subsSection = data?.subscribersView.subscribers;
  const plansSection = data?.subscribersView.plans;

  const subscribers: Subscriber[] = useMemo(() => {
    const plansById = new Map(
      (plansSection?.plans ?? []).map((p) => [p.packageId, p])
    );
    return (subsSection?.subscribers ?? []).map((s) =>
      toSubscriber(s, plansById)
    );
  }, [subsSection?.subscribers, plansSection?.plans]);

  const planNames = useMemo(
    () => [...(plansSection?.plans ?? []).map((p) => p.name), 'No plan'],
    [plansSection?.plans],
  );

  const columns = useMemo<ColumnDef<Subscriber, unknown>[]>(() => {
    const cols: ColumnDef<Subscriber, unknown>[] = [];

    cols.push({
      id: 'name',
      accessorKey: 'name',
      header: 'Customer',
      enableSorting: true,
      cell: ({ row }) => {
        const s = row.original;
        return (
          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <span className="av-sm">
              {s.name
                .split(' ')
                .map((x) => x[0])
                .join('')}
            </span>
            <div>
              <div style={{ fontWeight: 600 }}>{s.name}</div>
              <div className="muted tnum" style={{ fontSize: 12 }}>
                {showSite ? `${s.phone} · ${s.site}` : s.phone}
              </div>
            </div>
          </div>
        );
      },
    });

    cols.push({
      id: 'plan',
      accessorKey: 'plan',
      header: 'Plan',
      filterFn: 'equalsString',
      meta: { filterOptions: planNames },
      cell: ({ row }) =>
        row.original.plan === 'No plan' ? (
          <span style={{ color: 'var(--uk-ink-3)' }}>No plan</span>
        ) : (
          row.original.plan
        ),
    });

    cols.push({
      id: 'usage',
      accessorFn: (s) => s.usage,
      header: 'Data usage',
      enableSorting: true,
      cell: ({ row }) => {
        const s = row.original;
        // usage < 0 = unknown (subscribersView.usage backend gap) → "—"
        if (s.plan === 'No plan' || s.usage < 0)
          return <span className="muted">—</span>;
        const pct = s.cap ? Math.min(100, (s.usage / s.cap) * 100) : 60;
        const over = !!s.cap && s.usage / s.cap > 0.9;
        return (
          <div style={{ display: 'flex', alignItems: 'center', gap: 10, width: 150 }}>
            <Meter value={pct} color={over ? 'var(--uk-orange)' : undefined} sx={{ flex: 1, minWidth: 60 }} />
            <span
              className="tnum"
              style={{ fontSize: 12, color: 'var(--uk-ink-2)', whiteSpace: 'nowrap' }}
            >
              {s.usage}
              {s.cap ? '/' + s.cap : ''} GB
            </span>
          </div>
        );
      },
    });

    cols.push({
      id: 'sim',
      accessorKey: 'sim',
      header: 'SIM',
      filterFn: 'equalsString',
      meta: { filterOptions: ['active', 'inactive', 'suspended'] },
      cell: ({ row }) => {
        const s = row.original;
        return (
          <StatusBadge status={s.sim === 'suspended' ? 'pending' : s.sim}>
            {s.sim === 'suspended' ? 'Suspended' : undefined}
          </StatusBadge>
        );
      },
    });

    cols.push({
      id: 'seen',
      accessorFn: (s) => s.seen,
      header: 'Last seen',
      enableSorting: true,
      sortingFn: (a, b) =>
        parseSeen(a.original.seen) - parseSeen(b.original.seen),
      cell: ({ row }) => (
        <span className="muted tnum" style={{ fontSize: 13 }}>
          {row.original.seen}
        </span>
      ),
    });

    // Chevron affordance hints the row opens a detail drawer.
    if (clickRow) {
      cols.push({
        id: 'chevron',
        size: 40,
        header: '',
        cell: () => (
          <ChevronRightRounded
            sx={{ fontSize: 20, color: 'var(--uk-ink-3)', display: 'block' }}
          />
        ),
      });
    }
    return cols;
  }, [clickRow, showSite, planNames]);

  return (
    <div className="page">
      <PageHeader
        title="Customers"
        count={subscribers.length.toLocaleString()}
        sub={SUBS[mode]}
        actions={
          agent ? (
            <Button
              variant="contained"
              startIcon={<PersonAddRounded />}
              onClick={() => setShowAdd(true)}
            >
              Add customer
            </Button>
          ) : mode === 'biz' ? (
            <DateChip />
          ) : undefined
        }
      />
      <div className="card card-pad" style={{ paddingTop: 18 }}>
        <div
          style={{
            display: 'flex',
            gap: 10,
            marginBottom: 16,
            flexWrap: 'wrap',
            alignItems: 'center',
          }}
        >
          <SearchField value={q} onChange={setQ} placeholder="Search name or phone" />
          <div style={{ marginLeft: 'auto', fontSize: 13, color: 'var(--uk-ink-3)' }}>
            {subscribers.length} of {subscribers.length}
          </div>
        </div>

        <div className="tbl-wrap" style={{ overflowX: 'auto' }}>
          <DataTable<Subscriber>
            columns={columns}
            data={subscribers}
            status={loading ? 'loading' : subsSection?.error ? 'error' : 'ready'}
            skeleton={{ cols: clickRow ? 6 : 5, rows: 6, lead: true }}
            empty={{
              art: 'search',
              title: subsSection?.error ? "Couldn't load customers" : 'No customers match',
              sub: subsSection?.error
                ? subsSection.error.message
                : 'Try a different filter or search term.',
            }}
            globalFilter={q}
            initialSorting={[{ id: 'name', desc: false }]}
            getRowId={(s) => s.id}
            {...(clickRow ? { onRowClick: (s: Subscriber) => setOpenSub(s) } : {})}
          />
        </div>

        {!loading && (
          <TableFooter showing={subscribers.length} total={subscribers.length} />
        )}
      </div>

      {openSub && (
        <SubscriberDrawer
          sub={openSub}
          onClose={() => setOpenSub(null)}
          readOnly={mode === 'biz'}
        />
      )}
      {showAdd && (
        <AddCustomerDialog
          onClose={() => setShowAdd(false)}
          onAdded={() => void refetch()}
        />
      )}
    </div>
  );
}
