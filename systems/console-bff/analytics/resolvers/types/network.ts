/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

import { KpiDto, MetaDto, TimeSeriesDto } from "./shared";

@ObjectType()
export class EventRowDto {
  @Field({ nullable: true })
  routingKey?: string;

  @Field({ nullable: true })
  resourceType?: string;

  @Field({ nullable: true })
  resourceId?: string;

  @Field({ nullable: true })
  description?: string;

  @Field({ nullable: true })
  occurredAt?: string;
}

@ObjectType()
export class AlarmRowDto {
  @Field({ nullable: true })
  alarmId?: string;

  @Field({ nullable: true })
  severity?: string;

  @Field({ nullable: true })
  state?: string;

  @Field({ nullable: true })
  resourceType?: string;

  @Field({ nullable: true })
  resourceId?: string;

  @Field({ nullable: true })
  description?: string;

  @Field(() => Int)
  customersAffected: number;

  @Field(() => Float)
  revenueAtRisk: number;

  @Field({ nullable: true })
  recommendedAction?: string;

  @Field({ nullable: true })
  openedAt?: string;

  @Field({ nullable: true })
  closedAt?: string;
}

@ObjectType()
export class NetworkOverviewDto {
  @Field({ nullable: true })
  networkStatus?: string;

  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [EventRowDto])
  recentEvents: EventRowDto[];
}

@ObjectType()
export class TopologyNodeDto {
  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  type?: string;

  @Field({ nullable: true })
  status?: string;
}

@ObjectType()
export class TopologySiteDto {
  @Field({ nullable: true })
  siteId?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Float)
  latitude: number;

  @Field(() => Float)
  longitude: number;

  @Field(() => [TopologyNodeDto])
  nodes: TopologyNodeDto[];
}

@ObjectType()
export class NetworkTopologyDto {
  @Field(() => [TopologySiteDto])
  sites: TopologySiteDto[];
}

@ObjectType()
export class SiteRowDto {
  @Field()
  siteId: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Int)
  nodeCount: number;

  @Field(() => Int)
  customers: number;

  @Field(() => Float)
  uptime: number;

  @Field({ nullable: true })
  issueSummary?: string;

  @Field()
  backhaulLatencyHigh: boolean;

  @Field()
  batteryCritical: boolean;

  @Field(() => Float)
  offlineDurationSeconds: number;

  @Field(() => Float)
  latitude: number;

  @Field(() => Float)
  longitude: number;
}

@ObjectType()
export class NetworkSitesDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [SiteRowDto])
  sites: SiteRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class NetworkSiteDto {
  @Field(() => SiteRowDto, { nullable: true })
  site?: SiteRowDto;

  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [TimeSeriesDto])
  series: TimeSeriesDto[];

  @Field(() => [AlarmRowDto])
  alarms: AlarmRowDto[];
}

@ObjectType()
export class NodeRowDto {
  @Field()
  nodeId: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  type?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field({ nullable: true })
  siteName?: string;

  @Field(() => Float)
  uptime: number;

  @Field({ nullable: true })
  lastTelemetry?: string;

  @Field()
  noTelemetryWarning: boolean;

  @Field(() => Float)
  configuringDurationSeconds: number;
}

@ObjectType()
export class NetworkNodesDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [NodeRowDto])
  nodes: NodeRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class NetworkNodeDto {
  @Field(() => NodeRowDto, { nullable: true })
  node?: NodeRowDto;

  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [TimeSeriesDto])
  series: TimeSeriesDto[];

  @Field(() => [EventRowDto])
  recentEvents: EventRowDto[];
}

@ObjectType()
export class NodePoolDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [NodeRowDto])
  nodes: NodeRowDto[];
}

/** Shared shape for the radio / backhaul / power panels. */
@ObjectType()
export class MetricPanelDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [TimeSeriesDto])
  series: TimeSeriesDto[];

  @Field(() => [AlarmRowDto])
  alarms: AlarmRowDto[];
}

@ObjectType()
export class NetworkAlarmsDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [AlarmRowDto])
  alarms: AlarmRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class MetricInfoDto {
  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  unit?: string;

  @Field({ nullable: true })
  lastSampleAt?: string;

  @Field()
  stale: boolean;
}

@ObjectType()
export class NetworkMetricsDto {
  @Field(() => [MetricInfoDto])
  metrics: MetricInfoDto[];

  @Field(() => [TimeSeriesDto])
  series: TimeSeriesDto[];
}

@ObjectType()
export class NetworkEventsDto {
  @Field(() => [EventRowDto])
  events: EventRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class SupportResultDto {
  @Field({ nullable: true })
  resourceType?: string;

  @Field({ nullable: true })
  resourceId?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Int)
  customers: number;

  @Field(() => Float)
  uptime30d: number;

  @Field(() => Float)
  batteryPercent: number;

  @Field(() => Float)
  signalDbm: number;

  @Field({ nullable: true })
  statusSummary?: string;

  @Field({ nullable: true })
  recommendation?: string;
}

@ObjectType()
export class NetworkSupportSearchDto {
  @Field(() => [SupportResultDto])
  results: SupportResultDto[];
}

@InputType()
export class AnalyticsNodeInput {
  @Field()
  nodeId: string;

  @Field({ nullable: true })
  period?: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;

  @Field({ nullable: true })
  timezone?: string;
}
