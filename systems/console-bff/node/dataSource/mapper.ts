/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NODE_STATE } from "../../common/enums";
import { Node, NodeSite, NodeStateRes, Nodes } from "../resolvers/types";

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

const parseSite = (res: any): NodeSite => {
  return res.site
    ? {
        nodeId: res.site.nodeId,
        siteId: res.site.site_id,
        addedAt: res.site.added_at,
        networkId: res.site.network_id,
      }
    : {
        nodeId: null,
        siteId: null,
        addedAt: null,
        networkId: null,
      };
};

export const parseNodeRes = (res: any): Node => {
  return {
    id: res.id,
    name: res.name,
    type: res.type,
    site: parseSite(res),
    latitude: res.latitude,
    longitude: res.longitude,
    attached: parseAttached(res),
    status: {
      state: res.status.state,
      connectivity: res.status.connectivity,
    },
  };
};

export const getNodeState = (res: any): NodeStateRes => {
  const states: any = res.states;
  states.sort(
    (a: any, b: any) =>
      new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );

  return {
    id: states[0].id,
    nodeId: states[0].nodeId,
    createdAt: states[0].createdAt,
    currentState: states[0].currentState as NODE_STATE,
    previousState: states[0].previousState as NODE_STATE,
    previousStateId: states[0].previousStateId,
  };
};
