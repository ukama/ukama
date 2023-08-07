import { IsOptional } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

import { API_METHOD_TYPE } from "../enums";

export type ApiMethodDataDto = {
  url: String;
  method: API_METHOD_TYPE;
  params?: any;
  headers?: any;
  body?: any;
};

export type ErrorType = {
  message: string;
  code: number;
  description?: string;
};

export type TBooleanResponse = {
  success: boolean;
};

@ObjectType()
export class CBooleanResponse {
  @Field()
  success: boolean;
}

@ObjectType()
export class Meta {
  @Field()
  count: number;

  @Field()
  page: number;

  @Field()
  size: number;

  @Field()
  pages: number;
}

@ObjectType()
export class PaginationResponse {
  @Field()
  meta: Meta;
}

@InputType()
export class PaginationDto {
  @Field()
  pageNo: number;

  @Field()
  pageSize: number;
}

@ObjectType()
export class AuthType {
  @Field({ nullable: true })
  @IsOptional()
  Cookie?: string;

  @Field({ nullable: true })
  @IsOptional()
  Authorization?: string;
}

@ObjectType()
export class THeaders {
  @Field(() => AuthType)
  auth: AuthType;

  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field()
  orgName: string;
}

@ObjectType()
export class IdResponse {
  @Field()
  uuid: string;
}
