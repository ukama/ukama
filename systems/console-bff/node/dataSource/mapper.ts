/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Node, NodeSite, Nodes } from "../resolvers/types";

export const parseNodesRes = (res: any): Nodes => {
  const nodes = res.nodes.map((node: any) => {
    return parseNodeRes(node);
  });
  return {
    nodes: nodes,
  };
};

const parseAttached = (res: any): Node[] => {
  return res.attached
    ? res.attached?.map((node: any) => {
        return parseNodeRes(node);
      })
    : [];
};

const parseSite = (res: any): NodeSite | undefined => {
  return res.site
    ? {
        nodeId: res.site.nodeId,
        siteId: res.site.site_id,
        addedAt: res.site.added_at,
        networkId: res.site.network_id,
      }
    : undefined;
};

export const parseNodeRes = (res: any): Node => {
  return {
    id: res.id,
    name: res.name,
    type: res.type,
    orgId: res.org_id,
    site: parseSite(res),
    attached: parseAttached(res),
    status: {
      state: res.status.state,
      connectivity: res.status.connectivity,
    },
  };
};
