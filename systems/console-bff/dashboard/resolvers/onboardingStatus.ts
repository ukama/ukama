/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Activation state for the console onboarding flow — DERIVED from real data
 * (networks/sites/nodes), never stored. The console uses it to decide
 * whether to show the setup alert bar, which /configure step is legitimate,
 * and where "Continue setup" should resume.
 *
 * Errors are NOT swallowed into `false`: a wrong zero-state would re-onboard
 * a configured org. Any upstream failure surfaces as a GraphQL error and the
 * console treats it as "unknown" (no banner, no gating).
 */
import { Ctx, Field, ObjectType, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import type { AppContext } from "../../server/context";

@ObjectType()
export class OnboardingStatusDto {
  @Field()
  hasNetwork: boolean;

  @Field()
  hasSite: boolean;

  @Field()
  hasNode: boolean;

  /** Default (or first) network — lets the console deep-link configure steps. */
  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  networkName?: string;
}

@Resolver()
export class OnboardingStatusResolver {
  @Query(() => OnboardingStatusDto)
  async onboardingStatus(@Ctx() ctx: AppContext): Promise<OnboardingStatusDto> {
    const urls = ctx.urls;

    const networkUrl = await urls.url("network");
    const { networks } = await ctx.dataSources.network.getNetworks(networkUrl);

    if (networks.length === 0) {
      return { hasNetwork: false, hasSite: false, hasNode: false };
    }

    const primary = networks.find(n => n.isDefault) ?? networks[0];

    // Sites are per-network: any site in any network counts (an org with an
    // empty default network but a configured second network is activated).
    // allSettled: one broken network must not hide sites that DO exist; but
    // if no site was found AND a lookup failed, the state is unknown — fail
    // loud rather than report a false zero-state.
    const siteUrl = await urls.url("site");
    const siteResults = await Promise.allSettled(
      networks.map(n =>
        ctx.dataSources.site.getSites(siteUrl, { networkId: n.id })
      )
    );
    const hasSite = siteResults.some(
      res => res.status === "fulfilled" && res.value.sites.length > 0
    );
    const siteFailure = siteResults.find(res => res.status === "rejected");
    if (!hasSite && siteFailure) {
      throw new Error(
        `onboardingStatus: site lookup failed: ${siteFailure.reason}`
      );
    }

    // Nodes are informational (never activation-blocking): a down node
    // service must not break onboarding right after network creation.
    let hasNode = false;
    try {
      const nodeUrl = await urls.url("node");
      const nodeRes = await ctx.dataSources.node.getNodes(nodeUrl, {});
      hasNode = nodeRes.nodes.length > 0;
    } catch (err) {
      logger.warn(
        `onboardingStatus: node lookup failed (reporting hasNode=false): ${err}`
      );
    }

    return {
      hasSite,
      hasNode,
      hasNetwork: true,
      networkId: primary.id,
      networkName: primary.name,
    };
  }
}
