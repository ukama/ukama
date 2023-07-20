import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class NodeStatus {
  @Field()
  connectivity: string;

  @Field()
  state: string;
}
@ObjectType()
export class Node {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  orgId: string;

  @Field()
  type: string;

  @Field(() => [String])
  attached: string[];

  @Field(() => NodeStatus)
  status: NodeStatus;
}

@ObjectType()
export class GetNode {
  @Field(() => Node)
  node: Node;
}

@ObjectType()
export class GetNodes {
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

  @Field()
  state: string;
}

@ObjectType()
export class NodeState {
  @Field()
  id: string;

  @Field()
  state: string;
}

@ArgsType()
@InputType()
export class UpdateNodeInput {
  @Field()
  id: string;

  @Field()
  name: string;
}
