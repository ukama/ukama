import { Field, ObjectType } from "type-graphql";
import { API_METHOD_TYPE } from "./enums";

@ObjectType()
export class ApiMethodDataDto {
  @Field()
  url: string;

  @Field(() => API_METHOD_TYPE)
  method: API_METHOD_TYPE;

  @Field(() => String || Object || null, { nullable: true })
  params?: any;

  @Field(() => String || Object || null, { nullable: true })
  headers?: any;

  @Field(() => String || Object || null, { nullable: true })
  body?: any;
}

@ObjectType()
export class ErrorType {
  @Field()
  message: string;

  @Field()
  code: number;

  @Field({ nullable: true })
  description?: string;
}
