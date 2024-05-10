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
  backhaul_id: string;

  @Field()
  power_id: string;

  @Field()
  access_id: string;

  @Field()
  switch_id: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  latitude: number;

  @Field()
  longitude: number;

  @Field()
  install_date: string;

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
  backhaulId: string;

  @Field()
  powerId: string;

  @Field()
  accessId: string;

  @Field()
  switchId: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  latitude: number;

  @Field()
  longitude: number;

  @Field()
  installDate: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class SitesResDto {
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
  @Field(() => [SiteAPIDto])
  sites: SiteAPIDto[];
}

@InputType()
export class AddSiteInputDto {
  @Field()
  site: string;
}
