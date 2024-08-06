/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class RestartNodeResDto {
  @Field()
  Node_id: string;
}

@ObjectType()
export class RestartSiteResDto {
  @Field()
  Site_id: string;
}

@ObjectType()
export class EmptyResDto {}
