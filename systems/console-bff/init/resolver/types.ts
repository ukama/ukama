/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

import { ROLE_TYPE } from "../../common/enums";

@ObjectType()
export class InitSystemAPIResDto {
  @Field()
  systemName: string;

  @Field()
  systemId: string;

  @Field()
  orgName: string;

  @Field()
  certificate: string;

  @Field()
  ip: string;

  @Field()
  port: number;

  @Field()
  health: number;
}

@ObjectType()
export class ValidateSessionRes {
  @Field()
  userId: string;

  @Field()
  email: string;

  @Field()
  name: string;

  @Field()
  orgId: string;

  @Field()
  orgName: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;

  @Field()
  token: string;
}
