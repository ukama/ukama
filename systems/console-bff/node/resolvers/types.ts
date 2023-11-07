/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

import { NODE_STATUS, NODE_TYPE } from "../../common/enums";

@ObjectType()
export class NodeStatus {
  @Field()
  connectivity: string;

  @Field()
  state: string;
}

@ObjectType()
export class NodeSite {
  @Field({ nullable: true })
  nodeId: string;

  @Field({ nullable: true })
  siteId: string;

  @Field({ nullable: true })
  networkId: string;

  @Field({ nullable: true })
  addedAt: string;
}

@ObjectType()
export class AttachedNodes {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  orgId: string;

  @Field(() => NODE_TYPE)
  type: NODE_TYPE;

  @Field(() => NodeSite)
  site: NodeSite;

  @Field(() => NodeStatus)
  status: NodeStatus;
}
@ObjectType()
export class Node {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  orgId: string;

  @Field(() => NODE_TYPE)
  type: NODE_TYPE;

  @Field(() => [AttachedNodes])
  attached: AttachedNodes[];

  @Field(() => NodeSite)
  site: NodeSite;

  @Field(() => NodeStatus)
  status: NodeStatus;
}

@ObjectType()
export class Nodes {
  @Field(() => [Node])
  nodes: Node[];
}

@ArgsType()
@InputType()
export class NodeInput {
  @Field()
  id: string;
}

@ArgsType()
@InputType()
export class GetNodesInput {
  @Field()
  isFree?: boolean;
}

@ObjectType()
export class DeleteNode {
  @Field()
  id: string;
}

@ArgsType()
@InputType()
export class AttachNodeInput {
  @Field()
  anodel: string;

  @Field()
  anoder: string;

  @Field()
  parentNode: string;
}

@ArgsType()
@InputType()
export class AddNodeInput {
  @Field()
  name: string;

  @Field()
  id: string;

  @Field()
  orgId: string;
}

@ArgsType()
@InputType()
export class AddNodeToSiteInput {
  @Field()
  nodeId: string;

  @Field()
  networkId: string;

  @Field()
  siteId: string;
}

@ArgsType()
@InputType()
export class UpdateNodeStateInput {
  @Field()
  id: string;

  @Field(() => NODE_STATUS)
  state: NODE_STATUS;
}

@ObjectType()
export class NodeState {
  @Field()
  id: string;

  @Field(() => NODE_STATUS)
  state: NODE_STATUS;
}

@ObjectType()
export class AppChangeLog {
  @Field()
  version: string;

  @Field()
  date: number;
}

@ObjectType()
export class AppChangeLogs {
  @Field(() => [AppChangeLog])
  logs: AppChangeLog[];

  @Field(() => NODE_TYPE)
  type: NODE_TYPE;
}

@ObjectType()
export class NodeApp {
  @Field()
  name: string;

  @Field()
  date: number;

  @Field()
  version: string;

  @Field()
  cpu: string;

  @Field()
  memory: string;

  @Field()
  notes: string;
}

@ObjectType()
export class NodeApps {
  @Field(() => [NodeApp])
  apps: NodeApp[];

  @Field(() => NODE_TYPE)
  type: NODE_TYPE;
}

@ArgsType()
@InputType()
export class UpdateNodeInput {
  @Field()
  id: string;

  @Field()
  name: string;
}

@ArgsType()
@InputType()
export class NodeAppsChangeLogInput {
  @Field(() => NODE_TYPE)
  type: NODE_TYPE;
}

@ObjectType()
export class NodeLocation {
  @Field()
  id: string;

  @Field()
  lat: string;

  @Field()
  lng: string;
}

@ArgsType()
@InputType()
export class NodesInput {
  @Field()
  networkId: string;
}

@ObjectType()
export class NodesLocation {
  @Field()
  networkId: string;

  @Field(() => [NodeLocation])
  nodes: NodeLocation[];
}
