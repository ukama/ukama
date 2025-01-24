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

@ObjectType()
export class Component {
  @Field({ nullable: true })
  componentId?: string;

  @Field({ nullable: true })
  componentName?: string;

  @Field()
  elementType: string;
}

@ObjectType()
export class Site {
  @Field()
  siteId: string;

  @Field()
  siteName: string;

  @Field()
  elementType: string;

  @Field(() => [Component])
  components: Component[];
}

@ObjectType()
export class Subscribers {
  @Field()
  totalSubscribers: string;

  @Field()
  activeSubscribers: string;

  @Field()
  inactiveSubscribers: string;
}

@ObjectType()
export class Network {
  @Field()
  networkId: string;

  @Field()
  networkName: string;

  @Field()
  elementType: string;

  @Field(() => [Site])
  sites: Site[];

  @Field(() => Subscribers, { nullable: true })
  subscribers?: Subscribers;
}

@ObjectType()
export class DataPlan {
  @Field()
  planId: string;

  @Field()
  planName: string;

  @Field()
  elementType: string;
}

@ObjectType()
export class Sims {
  @Field()
  availableSims: string;

  @Field()
  consumed: string;

  @Field()
  totalSims: string;
}

@ObjectType()
export class Members {
  @Field()
  totalMembers: string;

  @Field()
  activeMembers: string;

  @Field()
  inactiveMembers: string;
}

@ObjectType()
export class Org {
  @Field()
  orgId: string;

  @Field()
  orgName: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  ownerName: string;

  @Field()
  ownerEmail: string;

  @Field()
  ownerId: string;

  @Field()
  elementType: string;

  @Field(() => [Network])
  networks?: Network[];

  @Field(() => [DataPlan])
  dataplans?: DataPlan[];

  @Field(() => Sims, { nullable: true })
  sims?: Sims;

  @Field(() => Members, { nullable: true })
  members?: Members;
}

@ObjectType()
export class OrgTreeRes {
  @Field(() => Org)
  org?: Org;
}
