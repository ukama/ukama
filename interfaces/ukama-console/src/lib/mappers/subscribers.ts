/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Maps subscribersView composite data onto the Subscriber view-model so the
 * shared customers table stays unchanged. Usage figures are a backend gap
 * (subscribersView.usage, gap #2): `usage: -1` is the "unknown" sentinel the
 * table renders as "—". Site attribution and last-seen are also
 * metrics-phase data.
 */
import type { NetworkCustomersQuery } from '@/client/graphql/network-customers.generated';
import type { Subscriber } from '@/data';

type QuerySubscriber = NonNullable<
  NonNullable<NetworkCustomersQuery['subscribersView']['subscribers']['subscribers']>
>[number];
type QueryPlan = NonNullable<
  NonNullable<NetworkCustomersQuery['subscribersView']['plans']['plans']>
>[number];

const toSimStatus = (status?: string): Subscriber['sim'] => {
  const s = (status ?? '').toLowerCase();
  if (s === 'active') return 'active';
  if (s === 'suspended') return 'suspended';
  return 'inactive';
};

export const toSubscriber = (
  sub: QuerySubscriber,
  plansById: Map<string, QueryPlan>
): Subscriber => {
  const firstSim = sub.sim?.[0];
  const activePackage = sub.sim
    ?.flatMap((s) => (s.package ? [s.package] : []))
    .find((p) => p.is_active);
  const planName = activePackage
    ? (plansById.get(activePackage.package_id)?.name ?? 'No plan')
    : 'No plan';
  return {
    id: sub.uuid,
    name: sub.name,
    phone: sub.phone || (firstSim?.msisdn ?? '—'),
    site: '—',
    plan: planName,
    // TODO(metrics-phase): subscribersView.usage — backend gap #2
    usage: -1,
    cap: null,
    sim: toSimStatus(firstSim?.status),
    iccid: firstSim?.iccid ?? '—',
    simId: firstSim?.id,
    seen: '—',
  };
};
