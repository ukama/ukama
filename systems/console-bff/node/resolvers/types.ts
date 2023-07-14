import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class NodeRes {
  @Field()
  nodeId: string;
}

@ArgsType()
@InputType()
export class NodeInput {
  @Field()
  nodeId: string;
}
