/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Agent data plans — read-only browse with search + sort, wired to getPackages. */
import { useMemo, useState } from 'react';
import Skeleton from '@mui/material/Skeleton';
import { useGetPackagesQuery } from '@/client/graphql/packages.generated';
import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { EmptyState } from '@/components/EmptyState';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import PlanCard from '@/features/plans/PlanCard';
import type { Plan } from '@/data';
import { heldQuery } from '@/lib/heldQuery';
import { packageToPlan } from '@/features/plans/mapPackage';
import { useNetworkId } from '@/lib/useNetworkId';

type Sort = 'price-asc' | 'price-desc' | 'data-desc';

function PlanSection({
  title,
  sub,
  plans,
}: {
  title: string;
  sub: string;
  plans: Plan[];
}) {
  if (plans.length === 0) return null;
  return (
    <div style={{ marginBottom: 26 }}>
      <div className="sec-head">
        <div className="sec-title">
          {title}
          <span className="cnt">{plans.length}</span>
        </div>
        <div style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>{sub}</div>
      </div>
      <div
        className="tile-grid"
        style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))' }}
      >
        {plans.map((p) => (
          <PlanCard key={p.id} plan={p} readOnly />
        ))}
      </div>
    </div>
  );
}

export default function AgentPlansScreen() {
  const [q, setQ] = useState('');
  const [sort, setSort] = useState<Sort>('price-asc');
  const vol = (data: string) =>
    /unlim/i.test(data) ? Infinity : parseFloat(data) || 0;

  const networkId = useNetworkId();
  const packagesResult = useGetPackagesQuery({
    // The BFF returns this network's plans plus org-wide ones (those with no
    // networkId). Passing undefined — no validated network id — returns every
    // plan in the org.
    variables: { networkId: networkId || undefined },
  });
  const { error } = packagesResult;
  const { data, loading } = heldQuery(packagesResult);
  const { data: networksData } = useGetNetworksQuery();
  const networkNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const n of networksData?.getNetworks.networks ?? [])
      m.set(n.id, n.name);
    return m;
  }, [networksData]);
  // Each plan keeps its own networkId so the scope chips can split the list
  // without PlanCard needing it. Mapping before filtering keeps a plan's card
  // colour stable as the chips change.
  const entries = useMemo(
    () =>
      (data?.getPackages.packages ?? []).map((p, i) => ({
        networkId: p.networkId || '',
        plan: packageToPlan(
          p,
          i,
          p.networkId ? networkNameById.get(p.networkId) : undefined,
        ),
      })),
    [data, networkNameById],
  );
  const plans = entries.map((e) => e.plan);

  const networkName = networkNameById.get(networkId);

  const bySort = (a: Plan, b: Plan) =>
    sort === 'price-asc'
      ? a.price - b.price
      : sort === 'price-desc'
        ? b.price - a.price
        : vol(b.data) - vol(a.data);
  const matches = entries.filter((e) =>
    e.plan.name.toLowerCase().includes(q.toLowerCase()),
  );
  // Two sections: plans scoped to the selected network, and org-wide plans,
  // which are assignable on any network.
  const networkPlans = matches
    .filter((e) => e.networkId)
    .map((e) => e.plan)
    .sort(bySort);
  const orgPlans = matches
    .filter((e) => !e.networkId)
    .map((e) => e.plan)
    .sort(bySort);
  const list = [...networkPlans, ...orgPlans];

  return (
    <div className="page">
      <PageHeader
        title="Data plans"
        count={loading ? undefined : plans.length}
        sub="Browse plans to assign, top up or change for your customers."
      />
      <div
        style={{
          display: 'flex',
          gap: 10,
          marginBottom: 18,
          flexWrap: 'wrap',
          alignItems: 'center',
        }}
      >
        <SearchField value={q} onChange={setQ} placeholder="Search plans" />
        <div className="seg">
          {(
            [
              ['price-asc', 'Price ↑'],
              ['price-desc', 'Price ↓'],
              ['data-desc', 'Data'],
            ] as const
          ).map(([k, l]) => (
            <button
              key={k}
              type="button"
              className={sort === k ? 'on' : ''}
              onClick={() => setSort(k)}
            >
              {l}
            </button>
          ))}
        </div>
        <div
          style={{ marginLeft: 'auto', fontSize: 13, color: 'var(--uk-ink-3)' }}
        >
          {list.length} of {plans.length}
        </div>
      </div>
      {loading ? (
        <div
          className="tile-grid"
          style={{
            gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
          }}
        >
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} variant="rounded" height={220} />
          ))}
        </div>
      ) : error ? (
        <div className="card">
          <EmptyState
            art="error"
            title="Couldn't load data plans"
            sub="Please try again in a moment."
          />
        </div>
      ) : list.length === 0 ? (
        <div className="card">
          <EmptyState
            art="search"
            title={plans.length === 0 ? 'No data plans yet' : 'No plans match'}
            sub={
              plans.length === 0
                ? 'Plans created in the business console will appear here.'
                : 'Try a different filter or search term.'
            }
          />
        </div>
      ) : (
        <>
          <PlanSection
            title="Network plans"
            sub={networkName ? `Only on ${networkName}` : 'Only on this network'}
            plans={networkPlans}
          />
          <PlanSection
            title="Org-wide plans"
            sub="Assignable on any network"
            plans={orgPlans}
          />
        </>
      )}
    </div>
  );
}
