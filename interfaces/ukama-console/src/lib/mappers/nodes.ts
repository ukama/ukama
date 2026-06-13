/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Maps BFF composite Node data onto the existing UkamaNode view-model so the
 * card/drawer/table components stay unchanged. cpu/mem/temp/fw/up are
 * metrics-phase data (backend gaps #6) — zero/null placeholders until then;
 * the screens render them as "—" via the §4.5 contract.
 */
import type { ViewNodeFragment } from '@/client/graphql/views-shared.generated';
import type { UkamaNode } from '@/data';

const NODE_TYPE_LABEL: Record<string, UkamaNode['type']> = {
  tnode: 'Tower node',
  anode: 'Amplifier node',
  cnode: 'Controller node',
  hnode: 'Home node',
};

export const toNodeStatus = (node: ViewNodeFragment): UkamaNode['status'] => {
  const connectivity = node.status.connectivity.toLowerCase();
  const state = node.status.state.toLowerCase();
  if (connectivity !== 'online') return 'offline';
  if (state === 'faulty') return 'degraded';
  // 'configured' is a completed, healthy state — only an unknown/pending state
  // is still "configuring".
  if (state === 'unknown') return 'configuring';
  return 'online';
};

export const toUkamaNode = (
  node: ViewNodeFragment,
  siteName?: string,
): UkamaNode => ({
  id: node.id,
  serial: node.id,
  name: node.name || undefined,
  connectivity: node.status.connectivity,
  state: node.status.state,
  type: NODE_TYPE_LABEL[node.type] ?? 'Tower node',
  // Only show a resolved site name; never the raw siteId.
  site: siteName ?? '—',
  status: toNodeStatus(node),
  // TODO(metrics-phase): cpu/mem/temp/fw/up come from nodeView.kpis —
  // backend gap #6 (docs in systems/console-bff/docs/backend-gaps.md)
  cpu: 0,
  mem: 0,
  temp: null,
  fw: '—',
  up: '—',
});
