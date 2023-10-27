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
  @Field()
  nodeId: string;

  @Field()
  siteId: string;

  @Field()
  networkId: string;

  @Field()
  addedAt: string;
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

  @Field(() => [Node])
  attached: Node[];

  @Field(() => NodeSite, { nullable: true })
  site?: NodeSite;

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

@ArgsType()
@InputType()
export class UpdateNodeInput {
  @Field()
  id: string;

  @Field()
  name: string;
}
