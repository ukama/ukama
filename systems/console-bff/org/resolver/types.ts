/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class OrgAPIDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  owner: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  certificate: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  created_at: string;
}

@ObjectType()
export class OrgsAPIResDto {
  @Field()
  user: string;

  @Field(() => [OrgAPIDto])
  owner_of: OrgAPIDto[];

  @Field(() => [OrgAPIDto])
  member_of: OrgAPIDto[];
}

@ObjectType()
export class OrgAPIResDto {
  @Field(() => OrgAPIDto)
  org: OrgAPIDto;
}

@ObjectType()
export class OrgDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  owner: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  certificate: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  createdAt: string;
}

@ObjectType()
export class OrgsResDto {
  @Field()
  user: string;

  @Field(() => [OrgDto])
  ownerOf: OrgDto[];

  @Field(() => [OrgDto])
  memberOf: OrgDto[];
}

@ObjectType()
export class OrgResDto {
  @Field(() => OrgDto)
  org: OrgDto;
}
