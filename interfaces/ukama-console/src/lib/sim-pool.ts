/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Customer-onboarding prerequisites — shared by the entry points (Add customer
 * / Allocate a SIM) so they can block (with guidance) when there's nothing to
 * assign. A customer needs both an available pool SIM and a data plan.
 */
import { useGetPackagesQuery } from '@/client/graphql/packages.generated';
import { useGetSimsFromPoolQuery } from '@/client/graphql/sims.generated';
import { Sim_Status, type Sim_Types } from '@/client/graphql/types';
import { publicEnv } from '@/lib/runtime-env';

/** Available = neither already allocated nor failed. */
export function useAvailablePoolSims(): {
  available: number;
  loading: boolean;
} {
  const { data, loading } = useGetSimsFromPoolQuery({
    variables: {
      data: { type: publicEnv().simType as Sim_Types, status: Sim_Status.All },
    },
    fetchPolicy: 'cache-and-network',
  });
  const available = (data?.getSimsFromPool.sims ?? []).filter(
    (s) => !s.isAllocated && !s.isFailed,
  ).length;
  return { available, loading };
}

/** Count of data plans (packages) available to assign. */
export function useAvailableDataPlans(): {
  available: number;
  loading: boolean;
} {
  const { data, loading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
  });
  return { available: (data?.getPackages.packages ?? []).length, loading };
}

/** User-facing copy when a prerequisite is missing. */
export const NO_POOL_SIMS_MESSAGE =
  'No SIMs available — please upload SIMs to your SIM pool first.';
export const NO_DATA_PLANS_MESSAGE =
  'No data plans yet — please create a data plan before adding customers.';
