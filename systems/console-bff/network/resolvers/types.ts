/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { NETWORK_STATUS } from "../../common/enums";

@ObjectType()
export class NetworkStatusDto {
  @Field()
  liveNode: number;

  @Field()
  totalNodes: number;

  @Field(() => NETWORK_STATUS)
  status: NETWORK_STATUS;
}

@ObjectType()
export class NetworkStatusResponse {
  @Field()
  status: string;

  @Field(() => NetworkStatusDto)
  data: NetworkStatusDto;
}

@ObjectType()
export class NetworkAPIDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  org_id: string;

  @Field()
  budget: number;

  @Field()
  is_deactivated: string;

  @Field()
  created_at: string;

  @Field(() => [String])
  allowed_countries: string[];

  @Field(() => [String])
  allowed_networks: string[];
}
@ObjectType()
export class NetworkAPIResDto {
  @Field(() => NetworkAPIDto)
  network: NetworkAPIDto;
}

@ObjectType()
export class NetworksAPIResDto {
  @Field()
  org_id: string;

  @Field(() => [NetworkAPIDto])
  networks: NetworkAPIDto[];
}

@ObjectType()
export class NetworkDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  orgId: string;

  @Field()
  budget: number;

  @Field()
  isDeactivated: string;

  @Field()
  createdAt: string;

  @Field(() => [String])
  countries: string[];

  @Field(() => [String])
  networks: string[];
}

@ObjectType()
export class NetworksResDto {
  @Field()
  orgId: string;

  @Field(() => [NetworkDto])
  networks: NetworkDto[];
}

@InputType()
export class AddNetworkInputDto {
  @Field()
  name: string;

  @Field()
  org: string;

  @Field()
  budget: number;

  @Field(() => [String], { nullable: true })
  countries: string[];

  @Field(() => [String], { nullable: true })
  networks: string[];
}
