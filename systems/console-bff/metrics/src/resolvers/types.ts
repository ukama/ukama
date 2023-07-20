import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class MetricValue {
  @Field(() => Number)
  x: number;

  @Field(() => Number)
  y: number;
}

@ObjectType()
export class MetricRes {
  @Field()
  env: string;

  @Field()
  nodeid: string;

  @Field()
  type: string;

  @Field(() => [MetricValue])
  value: MetricValue[];
}

@ArgsType()
@InputType()
export class GetMetricInput {
  @Field()
  nodeId: string;

  @Field()
  orgId: string;

  @Field()
  type: string;

  @Field()
  userId: string;
}
