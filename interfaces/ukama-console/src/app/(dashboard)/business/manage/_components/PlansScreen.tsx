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
import PageHeader from '@/components/PageHeader';
import CreatePlanDialog from '@/features/plans/CreatePlanDialog';
import PlanCard from '@/features/plans/PlanCard';
import { packageToPlan } from '@/features/plans/mapPackage';

export default function PlansScreen() {
  // null = closed; { pkg: null } = create; { pkg } = edit.
  const [dialog, setDialog] = useState<{ pkg: PackageFragment | null } | null>(
    null,
  );
  const create = () => setDialog({ pkg: null });

  const { data, loading, error } = useGetPackagesQuery();
  const packages = useMemo(() => data?.getPackages.packages ?? [], [data]);

  // Resolve a plan's networkId → network name for the card chip.
  const { data: networksData } = useGetNetworksQuery();
  const networkNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const n of networksData?.getNetworks.networks ?? [])
      m.set(n.id, n.name);
    return m;
  }, [networksData]);

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'Data plans']}
        title="Data plans"
        count={loading ? undefined : packages.length}
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
      ) : (
        <div
          className="tile-grid"
          style={{
            gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
          }}
        >
          {packages.map((pkg, i) => (
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
