/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import { GRAPHS_TYPE, SOFTWARE_STATUS, TIMEFRAME_FILTER } from "../../common/enums";
import { getGraphsKeyByType, getNodeTypeFromId } from "../../common/utils";
import { GetHealthReportInputDto } from "../../health/resolvers/types";
import { Node } from "../../node/resolvers/types";
import type { AppContext } from "../../server/context";
import { GetSoftwaresInput } from "../../software/resolvers/types";
import { ServiceUrlResolver } from "../baseUrls";
import { fetchLatestKpis } from "../kpis";
import { notImplementedSection, runSection } from "../section";
import {
  GapSection,
  HealthSection,
  KpisSection,
  NodeSection,
  NodeStateSection,
  NodeView,
  NodesSection,
  NodesView,
  SiteSection,
  SoftwareSection,
} from "./types";

type NodesViewRoot = NodesView & { _urls: ServiceUrlResolver };
type NodeViewRoot = NodeView & {
  _urls: ServiceUrlResolver;
  /** node core memo so `site`/`siblings` reuse the same fetch. */
  _node?: Promise<Node>;
};

/**
 * Nodes list composite (plan §3.1). Serves: Nodes list, Node pool (skips
 * `health`), map views.
 */
@Resolver(() => NodesView)
export class NodesViewResolver {
  @Query(() => NodesView)
  nodesView(
    @Ctx() ctx: AppContext,
    @Arg("networkId", { nullable: true }) networkId?: string
  ): NodesView {
    return Object.assign(new NodesView(), {
      networkId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => NodesSection)
  async nodes(
    @Root() root: NodesViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodesSection> {
    const { value, error } = await runSection("nodes", async () => {
      const url = await root._urls.url("node");
      const res = root.networkId
        ? await ctx.dataSources.node.getNodesByNetwork(url, root.networkId)
        : await ctx.dataSources.node.getNodes(url, {});
      return res.nodes;
    });
    return { nodes: value, error };
  }

  @FieldResolver(() => GapSection)
  health(): GapSection {
    // TODO(backend-gap): health — bulk health endpoint (today: one call per
    // node) — unblocks: nodesView.health (docs/backend-gaps.md #4)
    return { error: notImplementedSection("health").error };
  }
}

/**
 * Node detail composite (plan §3.1). Serves: Node detail (all sections),
 * node drawer/peek (core + `health` only).
 */
@Resolver(() => NodeView)
export class NodeViewResolver {
  @Query(() => NodeView)
  nodeView(@Arg("nodeId") nodeId: string, @Ctx() ctx: AppContext): NodeView {
    return Object.assign(new NodeView(), {
      nodeId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  /** Memoized node fetch shared by `node`, `site` and `siblings`. */
  private fetchNode(root: NodeViewRoot, ctx: AppContext): Promise<Node> {
    if (!root._node) {
      root._node = root._urls
        .url("node")
        .then(url => ctx.dataSources.node.getNode(url, { id: root.nodeId }));
    }
    return root._node;
  }

  @FieldResolver(() => NodeSection)
  async node(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodeSection> {
    const { value, error } = await runSection("node", () =>
      this.fetchNode(root, ctx)
    );
    return { node: value, error };
  }

  @FieldResolver(() => SiteSection)
  async site(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SiteSection> {
    const { value, error } = await runSection("site", async () => {
      const node = await this.fetchNode(root, ctx);
      const siteId = node.site?.siteId;
      if (!siteId) return null;
      const url = await root._urls.url("site");
      return ctx.dataSources.site.getSite(url, siteId);
    });
    return { site: value, error };
  }

  @FieldResolver(() => NodesSection)
  async siblings(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodesSection> {
    const { value, error } = await runSection("siblings", async () => {
      const node = await this.fetchNode(root, ctx);
      const networkId = node.site?.networkId;
      const url = await root._urls.url("node");
      const res = networkId
        ? await ctx.dataSources.node.getNodesByNetwork(url, networkId)
        : await ctx.dataSources.node.getNodes(url, {});
      return res.nodes;
    });
    return { nodes: value, error };
  }

  @FieldResolver(() => HealthSection)
  async health(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<HealthSection> {
    const { value, error } = await runSection("health", async () => {
      const url = await root._urls.url("health");
      return ctx.dataSources.health.list(url, {
        nodeId: root.nodeId,
        timeframe: TIMEFRAME_FILTER.LATEST,
      } as GetHealthReportInputDto);
    });
    return { health: value, error };
  }

  @FieldResolver(() => SoftwareSection)
  async software(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SoftwareSection> {
    const { value, error } = await runSection("software", async () => {
      const url = await root._urls.url("software");
      // The software service requires a status filter (legacy console sent
      // `unknown`); omitting it returns 400.
      return ctx.dataSources.software.getSoftwares(url, {
        name: "",
        nodeId: root.nodeId,
        status: SOFTWARE_STATUS.unknown,
      } as GetSoftwaresInput);
    });
    return { softwares: value, error };
  }

  @FieldResolver(() => NodeStateSection)
  async stateHistory(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodeStateSection> {
    const { value, error } = await runSection("stateHistory", async () => {
      const url = await root._urls.url("state");
      return ctx.dataSources.node.getNodeState(url, root.nodeId);
    });
    return { stateHistory: value, error };
  }

  @FieldResolver(() => KpisSection)
  async kpis(
    @Root() root: NodeViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<KpisSection> {
    // Phase 4: node-health KPIs (uptime/temps/memory by node type) polled
    // from the metric service (closes backend gap #6). The latest-metric
    // endpoint is org-scoped — node-level filtering lands with the metric
    // service's node filter (same behavior as the legacy getNodeLatestMetric).
    const { value, error } = await runSection("kpis", async () => {
      const url = await root._urls.url("metrics");
      const keys = getGraphsKeyByType(
        GRAPHS_TYPE.NODE_HEALTH,
        getNodeTypeFromId(root.nodeId)
      );
      return fetchLatestKpis(ctx.dataSources.metric, url, keys);
    });
    return { metrics: value, error };
  }

  @FieldResolver(() => GapSection)
  radioStatus(): GapSection {
    // TODO(backend-gap): controller — read endpoint for RF/service/internet
    // switch state (controller exposes only mutations today) — unblocks:
    // nodeView.radioStatus
    return { error: notImplementedSection("radioStatus").error };
  }
}
