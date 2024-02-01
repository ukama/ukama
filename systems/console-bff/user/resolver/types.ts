/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { OrgAPIDto, OrgDto } from "../../org/resolver/types";

@ObjectType()
export class UserResDto {
  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  uuid: string;

  @Field()
  phone: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  authId: string;

  @Field()
  registeredSince: string;
}

@ObjectType()
export class WhoamiDto {
  @Field(() => UserResDto)
  user: UserResDto;

  @Field(() => [OrgDto])
  ownerOf: OrgDto[];

  @Field(() => [OrgDto])
  memberOf: OrgDto[];
}

@InputType()
export class UserFistVisitInputDto {
  @Field()
  userId: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  firstVisit: boolean;
}

@ObjectType()
export class UserFistVisitResDto {
  @Field()
  firstVisit: boolean;
}

@ObjectType()
export class UserAPIObj {
  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  id: string;

  @Field()
  phone: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  auth_id: string;

  @Field()
  registered_since: string;
}

@ObjectType()
export class UserAPIResDto {
  @Field(() => [UserAPIObj])
  user: UserAPIObj;
}

@ObjectType()
export class WhoamiAPIDto {
  @Field(() => UserAPIObj)
  user: UserAPIObj;

  @Field(() => [OrgAPIDto])
  ownerOf: OrgAPIDto[];

  @Field(() => [OrgAPIDto])
  memberOf: OrgAPIDto[];
}
