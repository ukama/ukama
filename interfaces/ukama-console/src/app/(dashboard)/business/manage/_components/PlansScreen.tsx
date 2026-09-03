/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Data plans — PlanCard grid + create, wired to getPackages. */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import AddRounded from '@mui/icons-material/AddRounded';
import {
  useGetPackagesQuery,
  type PackageFragment,
} from '@/client/graphql/packages.generated';
import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { EmptyState } from '@/components/EmptyState';
import FilterChips from '@/components/FilterChips';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import { heldQuery } from '@/lib/heldQuery';
import CreatePlanDialog from '@/features/plans/CreatePlanDialog';
import PlanCard from '@/features/plans/PlanCard';
import { packageToPlan } from '@/features/plans/mapPackage';

/** Chip selection: every plan, the network-scoped ones, or the org-wide ones. */
const ALL = 'all';
const NETWORK_SCOPED = 'network';
const ORG_WIDE = 'org';

export default function PlansScreen() {
  // null = closed; { pkg: null } = create; { pkg } = edit.
  const [dialog, setDialog] = useState<{ pkg: PackageFragment | null } | null>(
    null,
  );
  const create = () => setDialog({ pkg: null });
  const [scope, setScope] = useState<string>(ALL);
  const [q, setQ] = useState('');

  const packagesResult = useGetPackagesQuery();
  const { error } = packagesResult;
  const { data, loading } = heldQuery(packagesResult);
  const packages = useMemo(() => data?.getPackages.packages ?? [], [data]);

  // Resolve a plan's networkId → network name for the card chip.
  const { data: networksData } = useGetNetworksQuery();
  const networkNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const n of networksData?.getNetworks.networks ?? []) m.set(n.id, n.name);
    return m;
  }, [networksData]);

  // Plans are either scoped to one network or org-wide (no networkId, so
  // assignable on every network). Each card names its own network, so the
  // chips split by scope rather than adding one chip per network.
  const networkScopedCount = packages.filter((p) => p.networkId).length;
  const orgWideCount = packages.length - networkScopedCount;
  const filtered = packages
    .filter((p) =>
      scope === NETWORK_SCOPED
        ? !!p.networkId
        : scope === ORG_WIDE
          ? !p.networkId
          : true,
    )
    .filter((p) => p.name.toLowerCase().includes(q.toLowerCase()));

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'Data plans']}
        title="Data plans"
        count={loading ? undefined : filtered.length}
        sub="Plans you can assign to customers."
        actions={
          <Button
            variant="contained"
            startIcon={<AddRounded />}
            onClick={create}
          >
            Create plan
          </Button>
        }
      />
      {!loading && !error && packages.length > 0 && (
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
          <FilterChips
            options={[
              { value: ALL, label: 'All', count: packages.length },
              {
                value: NETWORK_SCOPED,
                label: 'Network-specific',
                count: networkScopedCount,
              },
              { value: ORG_WIDE, label: 'Org-wide', count: orgWideCount },
            ]}
            value={scope}
            onChange={setScope}
          />
        </div>
      )}
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
      ) : packages.length === 0 ? (
        <div className="card">
          <EmptyState
            art="invoice"
            title="No data plans yet"
            sub="Create your first data plan to start assigning it to customers."
            cta="Create plan"
            onCta={create}
          />
        </div>
      ) : filtered.length === 0 ? (
        <div className="card">
          <EmptyState
            art="search"
            title={q ? 'No plans match' : 'No plans in this scope'}
            sub={
              q
                ? 'Try a different search term or filter.'
                : 'Pick another filter, or create a plan with this scope.'
            }
            cta={q ? undefined : 'Create plan'}
            onCta={q ? undefined : create}
          />
        </div>
      ) : (
        <div
          className="tile-grid"
          style={{
            gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
          }}
        >
          {filtered.map((pkg, i) => (
            <PlanCard
              key={pkg.uuid}
              plan={packageToPlan(
                pkg,
                i,
                pkg.networkId ? networkNameById.get(pkg.networkId) : undefined,
              )}
              onEdit={() => setDialog({ pkg })}
            />
          ))}
        </div>
      )}
      {dialog && (
        <CreatePlanDialog pkg={dialog.pkg} onClose={() => setDialog(null)} />
      )}
    </div>
  );
}
