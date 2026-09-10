/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, Float, InputType, ObjectType } from "type-graphql";

import { TIMEFRAME_FILTER } from "../../common/enums";

/** Per-app runtime resource usage, as reported by the node's app supervisor. */
@ObjectType()
export class AppResource {
  @Field(() => Float)
  cpuPercent: number;

  @Field(() => Float)
  memoryRssKb: number;

  @Field(() => Float)
  diskReadBytes: number;

  @Field(() => Float)
  diskWriteBytes: number;
}

@ObjectType()
export class App {
  @Field()
  name: string;

  @Field()
  version: string;

  @Field()
  tag: string;

  // Runtime lifecycle state from the node (e.g. "running").
  @Field()
  status: string;

  @Field(() => AppResource, { nullable: true })
  resource?: AppResource;
}

@ObjectType()
export class Apps {
  @Field(() => [App])
  apps: App[];
}

@InputType()
export class GetAppsInputDto {
  // Required: the node whose apps to list.
  @Field()
  nodeId: string;

  // Optional: filter to a single app by name (e.g. "deviced").
  @Field({ nullable: true })
  appName?: string;
}

@InputType()
export class GetHealthReportInputDto {
  @Field()
  id: string;

  @Field()
  timestamp: string;

  @Field(() => TIMEFRAME_FILTER)
  timeframe: TIMEFRAME_FILTER;

  @Field()
  nodeId: string;
}

@ObjectType()
export class HealthSystemInfo {
  @Field()
  id: string;

  @Field()
  healthId: string;

  @Field()
  name: string;

  @Field()
  value: string;
}

@ObjectType()
export class HealthResourceInfo {
  @Field()
  id: string;

  @Field()
  cappId: string;

  @Field()
  name: string;

  @Field()
  value: string;
}

@ObjectType()
export class HealthCappInfo {
  @Field()
  id: string;

  @Field()
  space: string;

  @Field()
  name: string;

  @Field()
  tag: string;

  @Field()
  status: string;

  @Field(() => [HealthResourceInfo])
  resources: HealthResourceInfo[];
}

@ObjectType()
export class HealthInfo {
  @Field()
  id: string;

  @Field()
  nodeId: string;

  @Field()
  timestamp: string;

  @Field(() => [HealthSystemInfo])
  system: HealthSystemInfo[];

  @Field(() => [HealthCappInfo])
  capps: HealthCappInfo[];
}

@ObjectType()
export class HealthReport {
  @Field(() => [HealthInfo])
  healths: HealthInfo[];
}

@InputType()
export class GetNodeInterfacesInputDto {
  @Field()
  nodeId: string;
}

@ObjectType()
export class CellularInterfaceInfo {
  @Field()
  available: boolean;

  @Field()
  error: string;

  // "on" / "off"; empty when the node's health report omits it.
  @Field()
  service: string;
}

@ObjectType()
export class RadioInterfaceInfo {
  @Field()
  available: boolean;

  // "on" / "off"; empty when the node's health report omits it.
  @Field()
  state: string;
}

/** Interface states from the node's latest health report. */
@ObjectType()
export class NodeInterfaces {
  @Field()
  nodeId: string;

  @Field(() => CellularInterfaceInfo, { nullable: true })
  cellular?: CellularInterfaceInfo;

  @Field(() => RadioInterfaceInfo, { nullable: true })
  radio?: RadioInterfaceInfo;
}
