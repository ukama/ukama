import { Field, ObjectType } from "type-graphql";

import { API_METHOD_TYPE } from "../enums";
import { IsOptional } from "class-validator";

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


@ObjectType()
export class AuthType {
    @Field()
    @IsOptional()
    Cookie?: string;

    @Field()
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