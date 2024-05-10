/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class SiteAPIDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  network_id: string;

  @Field()
  is_deactivated: string;

  @Field()
  created_at: string;
}

@ObjectType()
export class SiteDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  networkId: string;

  @Field()
  isDeactivated: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class SitesResDto {
  @Field()
  networkId: string;

  @Field(() => [SiteDto])
  sites: SiteDto[];
}

@ObjectType()
export class SiteAPIResDto {
  @Field(() => SiteAPIDto)
  site: SiteAPIDto;
}

@ObjectType()
export class SitesAPIResDto {
  @Field()
  network_id: string;

  @Field(() => [SiteAPIDto])
  sites: SiteAPIDto[];
}

@InputType()
export class AddSiteInputDto {
  @Field()
  site: string;
}
