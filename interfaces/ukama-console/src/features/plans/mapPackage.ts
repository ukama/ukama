/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Maps a gateway PackageDto into the UI `Plan` shape PlanCard renders.
 * `subs` (per-plan customer count) isn't part of getPackages, so it's left
 * at 0 until a subscriber-count source is wired; the card then shows no
 * fabricated revenue. Colors cycle a fixed palette for visual variety.
 */
import type { PackageFragment } from '@/client/graphql/packages.generated';
import type { Plan } from '@/data';

const PALETTE = [
  'var(--uk-ac)',
  'var(--uk-secondary)',
  'var(--uk-orange)',
  'var(--uk-success-bright)',
  'var(--uk-beige)',
  'var(--uk-ink-3)',
];

const formatData = (volume: number, unit: string): string => {
  if (!volume) return 'Unlimited';
  const u = unit?.trim() || 'GB';
  return `${volume} ${u}`;
};

export function packageToPlan(pkg: PackageFragment, index = 0): Plan {
  return {
    id: pkg.uuid,
    name: pkg.name,
    price: pkg.amount,
    data: formatData(pkg.dataVolume, pkg.dataUnit),
    days: Math.round(pkg.duration),
    subs: 0,
    color: PALETTE[index % PALETTE.length] ?? 'var(--uk-ac)',
  };
}
