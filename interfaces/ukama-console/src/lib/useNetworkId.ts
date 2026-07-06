/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Validated current network id.
 *
 * The selected id is persisted in `useUiPrefs` (localStorage). Reading it
 * raw is unsafe: after the selected network is deleted — or on an account
 * with no networks — the store still holds the old id, and a
 * `skip: !networkId` guard does NOT catch it (a stale UUID is truthy), so
 * every network-scoped query fires with an id that no longer exists.
 *
 * This hook returns the id ONLY when getNetworks confirms it (or a fallback)
 * is a real network, and `''` otherwise — including while getNetworks is
 * still loading. Callers keep `skip: !networkId` and therefore never fire a
 * query with a stale or not-yet-validated id. getNetworks is cache-first and
 * deduped by Apollo, so using this hook widely costs one shared request.
 */
import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { useUiPrefs } from '@/lib/store';

export function useNetworkId(): string {
  const stored = useUiPrefs((s) => s.networkId);
  const { data } = useGetNetworksQuery();
  const networks = data?.getNetworks.networks ?? [];

  if (networks.length === 0) return ''; // loading, or genuinely no networks
  const selected = networks.find((n) => n.id === stored);
  if (selected) return selected.id;

  const fallback = networks.find((n) => n.isDefault) ?? networks[0];
  return fallback?.id ?? '';
}
