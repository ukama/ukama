import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class Node {
  @Field()
  allocated: boolean;

  @Field()
  attached: string[];

  @Field()
  name: string;

  @Field()
  network: string;

  @Field()
  node: string;

  @Field()
  state: string;

  @Field()
  type: string;
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
  id: string;

  @Field()
  state: string;
}

@ArgsType()
@InputType()
export class AddNodeToNetworkInput {
  @Field()
  nodeId: string;
  @Field()
  networkId: string;
}

@ArgsType()
@InputType()
export class UpdateNodeState {
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
