/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType } from "type-graphql";

@InputType()
export class RestartNodeInputDto {
  @Field()
  nodeId: string;
}

@InputType()
export class RestartNodesInputDto {
  @Field()
  networkId: string;

  @Field(() => [String])
  nodeIds: string[];
}

@InputType()
export class RestartSiteInputDto {
  @Field()
  siteId: string;

  @Field()
  networkId: string;
}

@InputType()
export class ToggleInternetSwitchInputDto {
  @Field()
  siteId: string;

  @Field()
  port: number;

  @Field()
  status: boolean;
}

@InputType()
export class ToggleRFStatusInputDto {
  @Field()
  nodeId: string;

  @Field()
  status: boolean;
}
