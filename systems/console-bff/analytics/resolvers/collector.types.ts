/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class SourceStateDto {
  @Field({ nullable: true })
  source?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  detail?: string;

  @Field({ nullable: true })
  lastRunAt?: string;

  @Field({ nullable: true })
  lastSuccessAt?: string;
}

@ObjectType()
export class RollupStateDto {
  @Field({ nullable: true })
  rollup?: string;

  @Field({ nullable: true })
  watermark?: string;

  @Field()
  dirty: boolean;
}

@ObjectType()
export class RefreshResultDto {
  @Field(() => [SourceStateDto])
  states: SourceStateDto[];
}

@ObjectType()
export class RefreshStateDto {
  @Field(() => [SourceStateDto])
  states: SourceStateDto[];

  @Field(() => [RollupStateDto])
  rollups: RollupStateDto[];
}

@ObjectType()
export class RebuildRollupsResultDto {
  @Field(() => [RollupStateDto])
  rollups: RollupStateDto[];
}

@InputType()
export class RefreshInput {
  @Field()
  source: string;
}

@InputType()
export class RebuildRollupsInput {
  @Field()
  family: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;
}
