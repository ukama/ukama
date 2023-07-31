import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class LatestMetricRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  env: string;

  @Field()
  nodeid: string;

  @Field()
  type: string;

  @Field(() => [Number, Number])
  value: [number, number];
}
@ObjectType()
export class MetricRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  env: string;

  @Field()
  nodeid: string;

  @Field()
  type: string;

  @Field(() => [[Number, Number]])
  values: [number, number][];
}
@ArgsType()
@InputType()
export class GetLatestMetricInput {
  @Field()
  nodeId: string;

  @Field()
  type: string;
}

@ArgsType()
@InputType()
export class GetMetricRangeInput {
  @Field()
  nodeId: string;

  @Field()
  orgId?: string;

  @Field()
  type: string;

  @Field()
  userId?: string;

  @Field({ nullable: true })
  from: number;

  @Field({ nullable: true })
  to?: number;

  @Field({ nullable: true })
  step?: number;

  @Field({ nullable: true })
  withSubscription?: boolean;
}

@ArgsType()
@InputType()
export class SubMetricRangeInput {
  @Field()
  nodeId: string;

  @Field()
  orgId: string;

  @Field()
  type: string;

  @Field()
  userId: string;

  @Field()
  from: number;
}
