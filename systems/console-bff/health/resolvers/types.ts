/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { TIMEFRAME_FILTER } from "../../common/enums";

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
