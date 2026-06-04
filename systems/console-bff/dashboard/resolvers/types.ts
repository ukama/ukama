/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * View-domain composite types (plan §3). Each composite has a cheap core
 * (root resolver) and lazy sections (@FieldResolver) that hit upstreams only
 * when selected. Every section embeds its own nullable `error`
 * (SectionError) — see §4.5: data null + error set = failed/not-implemented;
 * data null + no error = genuinely empty.
 *
 * Schema is design-complete against frontend needs (§3.5): sections whose
 * backend doesn't exist yet are declared here and return NOT_IMPLEMENTED
 * until the gap closes (see docs/backend-gaps.md) — closing a gap never
 * changes this schema.
 */
import { Field, Int, ObjectType } from "type-graphql";

import { HealthInfo } from "../../health/resolvers/types";
import { NetworkDto } from "../../network/resolvers/types";
import { Node, NodeStateRes } from "../../node/resolvers/types";
import { NotificationsDto } from "../../notification/resolvers/types";
import { SimPoolResDto } from "../../sim/resolver/types";
import { SiteDto } from "../../site/resolvers/types";
import { Softwares } from "../../software/resolvers/types";
import { SubscriberDto } from "../../subscriber/resolver/types";
import { SectionError } from "../types";

/* ----------------------------- shared sections ---------------------------- */

@ObjectType()
export class NodeStatsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => Int, { nullable: true })
  total?: number | null;

  @Field(() => Int, { nullable: true })
  online?: number | null;

  @Field(() => Int, { nullable: true })
  offline?: number | null;
}

@ObjectType()
export class NodesSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [Node], { nullable: true })
  nodes?: Node[] | null;
}

@ObjectType()
export class SitesSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [SiteDto], { nullable: true })
  sites?: SiteDto[] | null;
}

@ObjectType()
export class SubscriberStatsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => Int, { nullable: true })
  total?: number | null;

  @Field(() => Int, { nullable: true })
  active?: number | null;

  @Field(() => Int, { nullable: true })
  inactive?: number | null;
}

@ObjectType()
export class AlertsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [NotificationsDto], { nullable: true })
  notifications?: NotificationsDto[] | null;
}

/** Placeholder-only section: schema-complete, backend gap (§3.5). */
@ObjectType()
export class GapSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;
}

/* ----------------------------- networkOverview ---------------------------- */

@ObjectType()
export class NetworkSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => NetworkDto, { nullable: true })
  network?: NetworkDto | null;
}

@ObjectType()
export class NetworkOverview {
  @Field()
  networkId: string;

  // Sections resolved lazily — see NetworkOverviewResolver.
  network?: NetworkSection;
  nodeStats?: NodeStatsSection;
  siteStats?: SitesSection;
  subscriberStats?: SubscriberStatsSection;
  latestAlerts?: AlertsSection;
  kpis?: GapSection;
}

/* -------------------------------- nodesView ------------------------------- */

@ObjectType()
export class NodesView {
  @Field({ nullable: true })
  networkId?: string;

  nodes?: NodesSection;
  health?: GapSection;
}

/* -------------------------------- nodeView -------------------------------- */

@ObjectType()
export class NodeSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => Node, { nullable: true })
  node?: Node | null;
}

@ObjectType()
export class HealthSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => HealthInfo, { nullable: true })
  health?: HealthInfo | null;
}

@ObjectType()
export class SoftwareSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => Softwares, { nullable: true })
  softwares?: Softwares | null;
}

@ObjectType()
export class NodeStateSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => NodeStateRes, { nullable: true })
  stateHistory?: NodeStateRes | null;
}

@ObjectType()
export class NodeView {
  @Field()
  nodeId: string;

  node?: NodeSection;
  health?: HealthSection;
  software?: SoftwareSection;
  stateHistory?: NodeStateSection;
  kpis?: GapSection;
  radioStatus?: GapSection;
}

/* -------------------------------- sitesView ------------------------------- */

@ObjectType()
export class SiteNodeCountDto {
  @Field()
  siteId: string;

  @Field(() => Int)
  total: number;

  @Field(() => Int)
  online: number;

  @Field(() => Int)
  offline: number;
}

@ObjectType()
export class SiteNodeCountsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [SiteNodeCountDto], { nullable: true })
  counts?: SiteNodeCountDto[] | null;
}

@ObjectType()
export class SitesView {
  @Field()
  networkId: string;

  sites?: SitesSection;
  nodeCounts?: SiteNodeCountsSection;
  kpis?: GapSection;
  financials?: GapSection;
}

/* -------------------------------- siteView -------------------------------- */

@ObjectType()
export class SiteComponentDto {
  @Field()
  elementType: string;

  @Field({ nullable: true })
  componentId?: string;

  @Field({ nullable: true })
  componentName?: string;
}

@ObjectType()
export class SiteComponentsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [SiteComponentDto], { nullable: true })
  components?: SiteComponentDto[] | null;
}

@ObjectType()
export class SiteSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => SiteDto, { nullable: true })
  site?: SiteDto | null;
}

@ObjectType()
export class SiteView {
  @Field()
  siteId: string;

  site?: SiteSection;
  nodes?: NodesSection;
  components?: SiteComponentsSection;
  power?: GapSection;
  kpis?: GapSection;
  financials?: GapSection;
}

/* ----------------------------- subscribersView ---------------------------- */

@ObjectType()
export class SubscribersSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [SubscriberDto], { nullable: true })
  subscribers?: SubscriberDto[] | null;
}

@ObjectType()
export class PlanNameDto {
  @Field()
  packageId: string;

  @Field()
  name: string;
}

@ObjectType()
export class PlansSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [PlanNameDto], { nullable: true })
  plans?: PlanNameDto[] | null;
}

@ObjectType()
export class SubscribersView {
  @Field()
  networkId: string;

  subscribers?: SubscribersSection;
  plans?: PlansSection;
  usage?: GapSection;
}

/* ------------------------------- simPoolView ------------------------------ */

@ObjectType()
export class SimPoolStatsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => Int, { nullable: true })
  total?: number | null;

  @Field(() => Int, { nullable: true })
  available?: number | null;

  @Field(() => Int, { nullable: true })
  consumed?: number | null;

  @Field(() => Int, { nullable: true })
  failed?: number | null;

  @Field(() => Int, { nullable: true })
  esim?: number | null;

  @Field(() => Int, { nullable: true })
  physical?: number | null;

  @Field(() => Int, { nullable: true })
  pctAssigned?: number | null;

  @Field(() => Boolean, { nullable: true })
  lowStock?: boolean | null;
}

@ObjectType()
export class PoolSimsSection {
  @Field(() => SectionError, { nullable: true })
  error?: SectionError | null;

  @Field(() => [SimPoolResDto], { nullable: true })
  sims?: SimPoolResDto[] | null;
}

@ObjectType()
export class SimPoolView {
  @Field()
  simType: string;

  stats?: SimPoolStatsSection;
  sims?: PoolSimsSection;
}
