/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@InputType()
export class DefaultMarkupInputDto {
  @Field({ nullable: false })
  markup: number;
}

@ObjectType()
export class DefaultMarkupResDto {
  @Field()
  markup: number;
}

@ObjectType()
export class DefaultMarkupAPIResDto {
  @Field()
  markup: number;
}

@ObjectType()
export class DefaultMarkupHistoryDto {
  @Field()
  createdAt: string;

  @Field()
  deletedAt: string;

  @Field()
  Markup: number;
}

@ObjectType()
export class DefaultMarkupHistoryAPIResDto {
  @Field(() => [DefaultMarkupHistoryDto], { nullable: true })
  markupRates: DefaultMarkupHistoryDto[];
}

@ObjectType()
export class DefaultMarkupHistoryResDto {
  @Field(() => [DefaultMarkupHistoryDto], { nullable: true })
  markupRates: DefaultMarkupHistoryDto[];
}
