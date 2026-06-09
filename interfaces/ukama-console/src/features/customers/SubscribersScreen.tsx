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
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import CloseRounded from '@mui/icons-material/CloseRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import SwapHorizRounded from '@mui/icons-material/SwapHorizRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import type { ColumnDef, RowSelectionState } from '@tanstack/react-table';
import DataTable, { selectionColumn } from '@/components/data-table/DataTable';
import TableFooter from '@/components/data-table/TableFooter';
import DateChip from '@/components/DateChip';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { useNetworkCustomersQuery } from '@/client/graphql/network-customers.generated';
import type { Subscriber } from '@/data';
import { parseSeen } from '@/lib/parsers';
import { useUiPrefs } from '@/lib/store';
import { toSubscriber } from '@/lib/mappers/subscribers';
import AddCustomerDialog from './AddCustomerDialog';
import DeleteCustomerDialog from './DeleteCustomerDialog';
import SubscriberDrawer from './SubscriberDrawer';

export type CustomersMode = 'biz' | 'network' | 'agent';

const SUBS = {
  biz: 'Who are my customers and what state are they in?',
  network: 'Everyone connected to your network.',
  agent: 'Manage your customers’ packages and top-ups.',
} as const;

function RowMenu({
  sub,
  onView,
  onDelete,
}: {
  sub: Subscriber;
  onView: () => void;
  onDelete: () => void;
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const toast = useToast();
  return (
    <>
      <IconButton
        size="small"
        aria-label="More actions"
        sx={{ color: 'var(--uk-ink-3)' }}
        onClick={(e) => {
          e.stopPropagation();
          setAnchor(e.currentTarget);
        }}
      >
        <MoreVertRounded sx={{ fontSize: 20 }} />
      </IconButton>
      <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            onView();
          }}
        >
          <VisibilityRounded sx={{ fontSize: 18 }} /> View details
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Top up ${sub.name} — flow lands with the form dialogs`);
          }}
        >
          <AddCardRounded sx={{ fontSize: 18 }} /> Top up data
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Change plan for ${sub.name} — flow lands with the form dialogs`);
          }}
        >
          <SwapHorizRounded sx={{ fontSize: 18 }} /> Change plan
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25, color: 'var(--uk-error)' }}
          onClick={() => {
            setAnchor(null);
            onDelete();
          }}
        >
          <DeleteOutlineRounded sx={{ fontSize: 18 }} /> Delete customer
        </MenuItem>
      </Menu>
    </>
  );
}

export default function SubscribersScreen({ mode }: { mode: CustomersMode }) {
  const agent = mode === 'agent';
  const showSite = mode === 'network';
  const clickRow = mode !== 'network';
  const networkId = useUiPrefs((s) => s.networkId);
  const toast = useToast();

  const [q, setQ] = useState('');
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [openSub, setOpenSub] = useState<Subscriber | null>(null);
  const [showAdd, setShowAdd] = useState(false);
  const [deleteSub, setDeleteSub] = useState<Subscriber | null>(null);

  const { data, loading } = useNetworkCustomersQuery({
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
    if (agent) cols.push(selectionColumn<Subscriber>());

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

    if (agent) {
      cols.push({
        id: 'actions',
        size: 40,
        header: '',
        cell: ({ row }) => (
          <RowMenu
            sub={row.original}
            onView={() => setOpenSub(row.original)}
            onDelete={() => setDeleteSub(row.original)}
          />
        ),
      });
    }
    return cols;
  }, [agent, showSite, planNames]);

  const selectedCount = Object.keys(selection).length;

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

        {agent && selectedCount > 0 && (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 14,
              background: 'var(--uk-ac-soft)',
              borderRadius: 10,
              padding: '9px 14px',
              marginBottom: 12,
            }}
          >
            <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--uk-ac-dark)' }}>
              {selectedCount} selected
            </span>
            <Button
              size="small"
              startIcon={<SwapHorizRounded />}
              onClick={() => toast(`Change plan for ${selectedCount} customers`)}
            >
              Change plan
            </Button>
            <Button
              size="small"
              startIcon={<AddCardRounded />}
              onClick={() => toast(`Top up ${selectedCount} customers`)}
            >
              Top up
            </Button>
            <Button
              size="small"
              color="inherit"
              startIcon={<CloseRounded />}
              sx={{ ml: 'auto', color: 'var(--uk-ink-3)' }}
              onClick={() => setSelection({})}
            >
              Clear
            </Button>
          </div>
        )}

        <div className="tbl-wrap" style={{ overflowX: 'auto' }}>
          <DataTable<Subscriber>
            columns={columns}
            data={subscribers}
            status={loading ? 'loading' : subsSection?.error ? 'error' : 'ready'}
            skeleton={{ cols: agent ? 7 : 5, rows: 6, lead: true }}
            empty={{
              art: 'search',
              title: subsSection?.error ? "Couldn't load customers" : 'No customers match',
              sub: subsSection?.error
                ? subsSection.error.message
                : 'Try a different filter or search term.',
            }}
            globalFilter={q}
            initialSorting={[{ id: 'name', desc: false }]}
            enableRowSelection={agent}
            rowSelection={selection}
            onRowSelectionChange={setSelection}
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
      {showAdd && <AddCustomerDialog onClose={() => setShowAdd(false)} />}
      {deleteSub && (
        <DeleteCustomerDialog sub={deleteSub} onClose={() => setDeleteSub(null)} />
      )}
    </div>
  );
}
