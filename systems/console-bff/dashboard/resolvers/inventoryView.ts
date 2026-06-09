/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import { SIM_TYPES } from "../../common/enums";
import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { countByCategory, deriveSimPool } from "../derive";
import { runSection } from "../section";
import {
  ComponentStatsSection,
  InventoryView,
  NodesSection,
  SimPoolStatsSection,
} from "./types";

type InventoryRoot = InventoryView & { _urls: ServiceUrlResolver };

/**
 * Business inventory composite (plan §3.2): component counts by category,
 * nodes not yet assigned to a site, and SIM stock.
 */
@Resolver(() => InventoryView)
export class InventoryViewResolver {
  @Query(() => InventoryView)
  inventoryView(@Ctx() ctx: AppContext): InventoryView {
    return Object.assign(new InventoryView(), {
      orgName: ctx.headers.orgName,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => ComponentStatsSection)
  async components(
    @Root() root: InventoryRoot,
    @Ctx() ctx: AppContext
  ): Promise<ComponentStatsSection> {
    const { value, error } = await runSection("components", async () => {
      const res = await ctx.dataSources.component.getComponentsByUserId(
        ctx.headers,
        "ALL"
      );
      return {
        total: res.components.length,
        byCategory: countByCategory(res.components),
      };
    });
    return { ...value, error };
  }

  @FieldResolver(() => NodesSection)
  async unassignedNodes(
    @Root() root: InventoryRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodesSection> {
    const { value, error } = await runSection("unassignedNodes", async () => {
      const url = await root._urls.url("node");
      const res = await ctx.dataSources.node.getNodes(url, {});
      return res.nodes.filter(node => !node.site?.siteId);
    });
    return { nodes: value, error };
  }

  @FieldResolver(() => SimPoolStatsSection)
  async simStock(
    @Root() root: InventoryRoot,
    @Ctx() ctx: AppContext
  ): Promise<SimPoolStatsSection> {
    const { value, error } = await runSection("simStock", async () => {
      const url = await root._urls.url("sim");
      const stats = await ctx.dataSources.sim.getSimPoolStats(
        url,
        SIM_TYPES.ukama_data
      );
      return { ...stats, ...deriveSimPool(stats) };
    });
    return { ...value, error };
  }
}
