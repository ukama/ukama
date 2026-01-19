/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

import { NODE_CONNECTIVITY, NODE_STATE, NODE_TYPE } from "../../common/enums";

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
  latitude: string;

  @Field()
  longitude: string;

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
  latitude: string;

  @Field()
  longitude: string;

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
export class NodesFilterInput {
  @Field({ nullable: true })
  id?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  type?: NODE_TYPE;

  @Field({ nullable: true })
  state?: NODE_STATE;

  @Field({ nullable: true })
  connectivity?: NODE_CONNECTIVITY;
}

@ArgsType()
@InputType()
export class NodeInput {
  @Field()
  id: string;
}

@ArgsType()
@InputType()
export class GetNodesByStateInput {
  @Field(() => NODE_CONNECTIVITY)
  connectivity: NODE_CONNECTIVITY;

  @Field(() => NODE_STATE)
  state: NODE_STATE;
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

  @Field(() => NODE_STATE)
  state: NODE_STATE;
}

@ObjectType()
export class NodeState {
  @Field()
  id: string;

  @Field(() => NODE_STATE)
  state: NODE_STATE;
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
  lat: number;

  @Field()
  lng: number;
}

@ArgsType()
@InputType()
export class NodesInput {
  @Field()
  networkId: string;

  @Field(() => NODE_STATE)
  nodeFilterState: NODE_STATE;
}

@ObjectType()
export class NodesLocation {
  @Field(() => [NodeLocation])
  nodes: NodeLocation[];
}

@ObjectType()
export class NodeStateRes {
  @Field()
  id: string;

  @Field()
  nodeId: string;

  @Field({ nullable: true })
  previousStateId: string;

  @Field(() => NODE_STATE, { nullable: true })
  previousState: NODE_STATE;

  @Field(() => NODE_STATE)
  currentState: NODE_STATE;

  @Field()
  createdAt: string;
}
