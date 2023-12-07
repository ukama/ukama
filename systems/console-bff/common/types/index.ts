/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IsOptional } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

import { API_METHOD_TYPE } from "../enums";

export type ApiMethodDataDto = {
  url: string;
  method: API_METHOD_TYPE;
  params?: any;
  headers?: any;
  httpsAgent?: any;
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
