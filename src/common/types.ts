import "reflect-metadata";
import { Field, InputType, ObjectType } from "type-graphql";
import { API_METHOD_TYPE } from "../constants";
import { Request } from "express";

@InputType()
export class PaginationDto {
    @Field()
    pageNo: number;

    @Field()
    pageSize: number;
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
export class ApiMethodDataDto {
    @Field()
    path: string;

    @Field(() => API_METHOD_TYPE)
    type: API_METHOD_TYPE;

    @Field(() => String || Object || null, { nullable: true })
    params?: any;

    @Field(() => String || Object || null, { nullable: true })
    headers?: any;

    @Field(() => String || Object || null, { nullable: true })
    body?: any;
}

export interface Context {
    req: Request;
    cookie: string | string[];
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
