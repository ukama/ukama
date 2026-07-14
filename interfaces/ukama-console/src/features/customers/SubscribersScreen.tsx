/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import { useNetworkCustomersQuery } from '@/client/graphql/network-customers.generated';
import { useGetSimsUsageByNetworkQuery } from '@/client/graphql/sims.generated';
import DataTable from '@/components/data-table/DataTable';
import PageHeader from '@/components/PageHeader';
import PageWatermark from '@/components/PageWatermark';
import SearchField from '@/components/SearchField';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';
import { toSubscriber } from '@/lib/mappers/subscribers';
import { parseSeen } from '@/lib/parsers';
import {
  NO_DATA_PLANS_MESSAGE,
  NO_POOL_SIMS_MESSAGE,
  useAvailableDataPlans,
  useAvailablePoolSims,
} from '@/lib/sim-pool';
import { useNetworkId } from '@/lib/useNetworkId';
import { bytesToGB, formatBytes, parseUsageBytes } from '@/lib/usage';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import GroupRounded from '@mui/icons-material/GroupRounded';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import type { ColumnDef } from '@tanstack/react-table';
import { useRouter } from 'next/navigation';
import { useMemo, useState } from 'react';
import AddCustomerDialog from './AddCustomerDialog';
import SubscriberDrawer from './SubscriberDrawer';

export type CustomersMode = 'biz' | 'network' | 'agent';

const SUBS = {
  biz: 'Who are my customers and what state are they in?',
  network: 'Everyone connected to your network.',
  agent: 'Manage your customers’ packages and top-ups.',
} as const;

export default function SubscribersScreen({ mode }: { mode: CustomersMode }) {
  const router = useRouter();
  const toast = useToast();
  const agent = mode === 'agent';
  const showSite = mode === 'network';
  const clickRow = mode !== 'network';
  const networkId = useNetworkId();

  const [q, setQ] = useState('');
  const [openSub, setOpenSub] = useState<Subscriber | null>(null);
  const [showAdd, setShowAdd] = useState(false);

  // Adding a customer assigns them a SIM + data plan, so a missing pool SIM or
  // data plan blocks the flow (guide the user to set the prerequisite up).
  const { available: poolSims } = useAvailablePoolSims();
  const { available: dataPlans } = useAvailableDataPlans(networkId);
  const openAddCustomer = () => {
    if (poolSims === 0) {
      toast(NO_POOL_SIMS_MESSAGE);
      return;
    }
    if (dataPlans === 0) {
      toast(NO_DATA_PLANS_MESSAGE);
      return;
    }
    setShowAdd(true);
  };

  const { data, loading, refetch } = useNetworkCustomersQuery({
    variables: { networkId },
    skip: !networkId,
  });
  const subsSection = data?.subscribersView.subscribers;
  const plansSection = data?.subscribersView.plans;

  // One aggregated usage call for the whole network (BFF fans out per SIM),
  // indexed by SIM id and converted from bytes to GB to match `cap`.
  const { data: usageData, loading: usageLoading } =
    useGetSimsUsageByNetworkQuery({
      variables: { networkId },
      skip: !networkId,
      fetchPolicy: 'cache-and-network',
    });
  const usageBySim = useMemo(() => {
    const m = new Map<string, number>();
    for (const u of usageData?.getSimsUsageByNetwork ?? []) {
      const bytes = parseUsageBytes(u.usage); // raw bytes; formatted on render
      if (bytes != null) m.set(u.simId, bytes);
    }
    return m;
  }, [usageData]);

  const subscribers: Subscriber[] = useMemo(() => {
    const plansById = new Map(
      (plansSection?.plans ?? []).map((p) => [p.packageId, p]),
    );
    return (subsSection?.subscribers ?? []).map((s) => {
      const sub = toSubscriber(s, plansById);
      const usage = sub.simId ? usageBySim.get(sub.simId) : undefined;
      return usage != null ? { ...sub, usage } : sub;
    });
  }, [subsSection?.subscribers, plansSection?.plans, usageBySim]);

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
        if (s.plan === 'No plan') return <span className="muted">—</span>;
        // Usage is fetched separately (getSimsUsageByNetwork) — show a skeleton
        // until it resolves, rather than a premature "—".
        if (usageLoading && s.usage < 0)
          return <Skeleton variant="rounded" width={140} height={16} />;
        if (s.usage < 0) return <span className="muted">—</span>;
        const over = !!s.cap && bytesToGB(s.usage) / s.cap > 0.9;
        return (
          <span
            className="tnum"
            style={{
              fontSize: 12,
              color: over ? 'var(--uk-orange)' : 'var(--uk-ink-2)',
              whiteSpace: 'nowrap',
            }}
          >
            {formatBytes(s.usage)}
            {s.cap ? ` / ${s.cap} GB` : ''}
          </span>
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
  }, [clickRow, showSite, planNames, usageLoading]);

  return (
    <div
      className="page"
      style={{ position: 'relative', overflow: 'hidden', isolation: 'isolate' }}
    >
      <PageWatermark icon={GroupRounded} />
      <PageHeader
        title="Customers"
        count={
          subscribers.length ? subscribers.length.toLocaleString() : undefined
        }
        sub={SUBS[mode]}
        actions={
          agent ? (
            <Button
              variant="contained"
              startIcon={<PersonAddRounded />}
              onClick={openAddCustomer}
            >
              Add customer
            </Button>
          ) : undefined
        }
      />
      <div
        style={{
          marginTop: 4,
          flex: 1,
          minHeight: 0,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div
          style={{
            display: 'flex',
            gap: 10,
            marginBottom: 16,
            flexWrap: 'wrap',
            alignItems: 'center',
          }}
        >
          <SearchField
            value={q}
            onChange={setQ}
            placeholder="Search name or phone"
          />
        </div>

        <div
          className="tbl-wrap"
          style={{ overflow: 'auto', flex: 1, minHeight: 0 }}
        >
          <DataTable<Subscriber>
            columns={columns}
            data={subscribers}
            status={
              loading ? 'loading' : subsSection?.error ? 'error' : 'ready'
            }
            skeleton={{ cols: clickRow ? 6 : 5, rows: 6, lead: true }}
            empty={
              subsSection?.error
                ? {
                    art: 'error',
                    title: "Couldn't load customers",
                    sub: subsSection.error.message,
                  }
                : subscribers.length === 0
                  ? {
                      // Genuinely no customers yet — encourage setup rather
                      // than implying a filter hid them. No icon (the page
                      // watermark already carries the visual).
                      art: null,
                      title: 'No customers yet',
                      sub: agent
                        ? 'Add your first customer to start managing their packages and top-ups.'
                        : 'Set up your network and add customers to start serving them.',
                      cta: agent ? 'Add customer' : 'Set up network',
                      onCta: agent
                        ? openAddCustomer
                        : () => router.push('/configure'),
                    }
                  : {
                      art: 'search',
                      title: 'No customers match',
                      sub: 'Try a different filter or search term.',
                    }
            }
            globalFilter={q}
            initialSorting={[{ id: 'name', desc: false }]}
            getRowId={(s) => s.id}
            {...(clickRow
              ? { onRowClick: (s: Subscriber) => setOpenSub(s) }
              : {})}
          />
        </div>
      </div>

      {openSub && (
        <SubscriberDrawer
          sub={openSub}
          onClose={() => setOpenSub(null)}
          agent={agent}
          onChanged={() => void refetch()}
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
