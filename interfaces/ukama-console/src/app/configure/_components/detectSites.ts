/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Site-readiness detection for the onboarding flow.
 *
 * A site is a trio of nodes that share a base id and differ only by type:
 *   tower      uk-<base>-tnode-<rev>
 *   amplifier  uk-<base>-anode-<rev>
 *   controller uk-<base>-cnode-<rev>
 *
 * A site is "ready to configure" only when ALL of the following hold:
 *   1. a tower (tnode) node exists, is online, and is not yet configured
 *      (connectivity Online + state Unknown);
 *   2. the tower has latitude/longitude data;
 *   3. its matching amplifier (anode) and controller (cnode) exist and are
 *      also online + not yet configured.
 *
 * Until all three are ready we do not advance — the physical install may
 * still be in progress.
 */
import {
  NodeConnectivityEnum,
  NodeStateEnum,
  NodeTypeEnum,
} from '@/client/graphql/types';
import type { GetNodesQuery } from '@/client/graphql/nodes.generated';

export type DetectedNode = GetNodesQuery['getNodes']['nodes'][number];

/** Online and not yet configured (status the BFF reports for fresh nodes). */
const isPoweredAndUnconfigured = (n: DetectedNode): boolean =>
  n.status.connectivity === NodeConnectivityEnum.Online &&
  n.status.state === NodeStateEnum.Unknown;

const hasCoordinates = (n: DetectedNode): boolean => {
  const lat = n.latitude?.trim();
  const lng = n.longitude?.trim();
  return Boolean(lat) && Boolean(lng);
};

/** Strips the node-type token so the three units of one site share a key. */
const baseKey = (id: string): string =>
  id.replace(/-(tnode|anode|cnode|hnode)-/, '-*-');

/** Per-unit readiness for the guided checklist (one site's three units). */
export interface SiteReadiness {
  /** Tower powered on + online (not yet configured). */
  tower: boolean;
  /** Amplifier powered on + online. */
  amplifier: boolean;
  /** Controller powered on + online. */
  controller: boolean;
  /** Tower has reported its GPS location. */
  located: boolean;
  /** All four checks passed — safe to create the site. */
  ready: boolean;
  /** The anchor tower record, once one is registered. */
  towerNode?: DetectedNode;
}

const EMPTY_READINESS: SiteReadiness = {
  tower: false,
  amplifier: false,
  controller: false,
  located: false,
  ready: false,
};

const stepsDone = (r: SiteReadiness): number =>
  Number(r.tower) + Number(r.amplifier) + Number(r.controller) + Number(r.located);

/**
 * Computes the guided checklist state for the single most-progressed site.
 * Units are grouped by their shared base id, so siblings are tracked even
 * before any of them powers on (the records exist once the hardware is
 * registered). `pinnedTowerId` (the deep-link nid) forces a specific site.
 */
export function computeSiteReadiness(
  nodes: DetectedNode[],
  pinnedTowerId?: string,
): SiteReadiness {
  const available = nodes.filter((n) => !n.site?.siteId);

  // Group the units of each prospective site together.
  const groups = new Map<string, DetectedNode[]>();
  for (const n of available) {
    const key = baseKey(n.id);
    const list = groups.get(key) ?? [];
    list.push(n);
    groups.set(key, list);
  }

  const candidates: SiteReadiness[] = [];
  for (const members of groups.values()) {
    const tower = members.find((n) => n.type === NodeTypeEnum.Tnode);
    if (!tower) continue; // a site must have a tower
    if (pinnedTowerId && tower.id !== pinnedTowerId) continue;

    const amplifier = members.find((n) => n.type === NodeTypeEnum.Anode);
    const controller = members.find((n) => n.type === NodeTypeEnum.Cnode);

    const r: SiteReadiness = {
      tower: isPoweredAndUnconfigured(tower),
      amplifier: Boolean(amplifier && isPoweredAndUnconfigured(amplifier)),
      controller: Boolean(controller && isPoweredAndUnconfigured(controller)),
      located: hasCoordinates(tower),
      ready: false,
      towerNode: tower,
    };
    r.ready = r.tower && r.amplifier && r.controller && r.located;
    candidates.push(r);
  }

  // Show the site closest to completion (ready first, then most steps done).
  candidates.sort(
    (a, b) => Number(b.ready) - Number(a.ready) || stepsDone(b) - stepsDone(a),
  );
  return candidates[0] ?? EMPTY_READINESS;
}
