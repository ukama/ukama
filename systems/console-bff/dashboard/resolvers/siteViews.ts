/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import { ComponentDto } from "../../component/resolvers/types";
import type { AppContext } from "../../server/context";
import { SiteDto } from "../../site/resolvers/types";
import { ServiceUrlResolver } from "../baseUrls";
import { groupNodeCountsBySite } from "../derive";
import { notImplementedSection, runSection } from "../section";
import {
  GapSection,
  NodesSection,
  SiteComponentDto,
  SiteComponentsSection,
  SiteNodeCountsSection,
  SiteSection,
  SiteView,
  SitesSection,
  SitesView,
} from "./types";

type SitesViewRoot = SitesView & { _urls: ServiceUrlResolver };
type SiteViewRoot = SiteView & {
  _urls: ServiceUrlResolver;
  /** site core memo so `components` doesn't re-fetch the site. */
  _site?: Promise<SiteDto>;
};

const toSiteComponent = (
  elementType: string,
  component?: ComponentDto
): SiteComponentDto => ({
  elementType,
  componentId: component?.partNumber ?? undefined,
  componentName: component?.partNumber ? component.description : undefined,
});

/**
 * Sites list composite (plan §3.1). Serves: Network sites list (skips
 * `financials`), Business sites list (adds `financials`, skips `kpis`).
 */
@Resolver(() => SitesView)
export class SitesViewResolver {
  @Query(() => SitesView)
  sitesView(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: AppContext
  ): SitesView {
    return Object.assign(new SitesView(), {
      networkId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => SitesSection)
  async sites(
    @Root() root: SitesViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SitesSection> {
    const { value, error } = await runSection("sites", async () => {
      const url = await root._urls.url("site");
      const res = await ctx.dataSources.site.getSites(url, {
        networkId: root.networkId,
      });
      return res.sites;
    });
    return { sites: value, error };
  }

  @FieldResolver(() => SiteNodeCountsSection)
  async nodeCounts(
    @Root() root: SitesViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SiteNodeCountsSection> {
    const { value, error } = await runSection("nodeCounts", async () => {
      const url = await root._urls.url("node");
      const res = await ctx.dataSources.node.getNodesByNetwork(
        url,
        root.networkId
      );
      return groupNodeCountsBySite(res.nodes);
    });
    return { counts: value, error };
  }

  @FieldResolver(() => GapSection)
  kpis(): GapSection {
    // TODO(backend-gap): metric — per-site latest KPI summary for list rows —
    // unblocks: sitesView.kpis (metrics phase, plan Phase 4)
    return { error: notImplementedSection("kpis").error };
  }

  @FieldResolver(() => GapSection)
  financials(): GapSection {
    // TODO(backend-gap): billing — per-site revenue/cost rollup — unblocks:
    // sitesView.financials (plan Phase 3, business lens)
    return { error: notImplementedSection("financials").error };
  }
}

/**
 * Site detail composite (plan §3.1). Serves: Network site detail (skips
 * `financials`), Business site detail (adds `financials`).
 */
@Resolver(() => SiteView)
export class SiteViewResolver {
  @Query(() => SiteView)
  siteView(@Arg("siteId") siteId: string, @Ctx() ctx: AppContext): SiteView {
    return Object.assign(new SiteView(), {
      siteId,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  private fetchSite(root: SiteViewRoot, ctx: AppContext): Promise<SiteDto> {
    if (!root._site) {
      root._site = root._urls
        .url("site")
        .then(url => ctx.dataSources.site.getSite(url, root.siteId));
    }
    return root._site;
  }

  @FieldResolver(() => SiteSection)
  async site(
    @Root() root: SiteViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SiteSection> {
    const { value, error } = await runSection("site", () =>
      this.fetchSite(root, ctx)
    );
    return { site: value, error };
  }

  @FieldResolver(() => NodesSection)
  async nodes(
    @Root() root: SiteViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<NodesSection> {
    const { value, error } = await runSection("nodes", async () => {
      const url = await root._urls.url("node");
      const res = await ctx.dataSources.node.getNodesForSite(url, root.siteId);
      return res.nodes;
    });
    return { nodes: value, error };
  }

  @FieldResolver(() => SiteComponentsSection)
  async components(
    @Root() root: SiteViewRoot,
    @Ctx() ctx: AppContext
  ): Promise<SiteComponentsSection> {
    const { value, error } = await runSection("components", async () => {
      const [site, compRes] = await Promise.all([
        this.fetchSite(root, ctx),
        ctx.dataSources.component.getComponentsByUserId(ctx.headers, "ALL"),
      ]);
      const byId = new Map<string, ComponentDto>(
        compRes.components.map(comp => [comp.id, comp])
      );
      return [
        toSiteComponent("ACCESS", byId.get(site.accessId)),
        toSiteComponent("POWER", byId.get(site.powerId)),
        toSiteComponent("BACKHAUL", byId.get(site.backhaulId)),
        toSiteComponent("SWITCH", byId.get(site.switchId)),
      ];
    });
    return { components: value, error };
  }

  @FieldResolver(() => GapSection)
  power(): GapSection {
    // TODO(backend-gap): metric — battery/solar/backhaul live status for a
    // site — unblocks: siteView.power (metrics phase, plan Phase 4)
    return { error: notImplementedSection("power").error };
  }

  @FieldResolver(() => GapSection)
  kpis(): GapSection {
    // TODO(backend-gap): metric — site uptime/KPI summary — unblocks:
    // siteView.kpis (metrics phase, plan Phase 4)
    return { error: notImplementedSection("kpis").error };
  }

  @FieldResolver(() => GapSection)
  financials(): GapSection {
    // TODO(backend-gap): billing — per-site revenue/cost rollup — unblocks:
    // siteView.financials (plan Phase 3, business lens)
    return { error: notImplementedSection("financials").error };
  }
}
