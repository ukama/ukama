import "reflect-metadata";
import { Request } from "express";
import { IsOptional } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";
import { API_METHOD_TYPE, GRAPHS_TAB, NODE_TYPE } from "../constants";

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
    cookie: string;
    token: string | undefined;
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

@ObjectType()
export class HeaderType {
    @Field()
    @IsOptional()
    Cookie?: string;

    @Field()
    @IsOptional()
    Authorization?: string;
}

@ObjectType()
export class ParsedCookie {
    @Field()
    header: HeaderType;

    @Field()
    orgId: string;
}

@InputType()
export class MetricsInputDTO {
    @Field()
    orgId: string;

    @Field()
    regPolling: boolean;

    @Field()
    nodeId: string;

    @Field()
    from: number;

    @Field()
    to: number;

    @Field()
    step: number;
}

@InputType()
export class MetricsByTabInputDTO {
    @Field()
    orgId: string;

    @Field()
    regPolling: boolean;

    @Field()
    nodeId: string;

    @Field()
    from: number;

    @Field()
    to: number;

    @Field()
    step: number;

    @Field(() => GRAPHS_TAB)
    tab: GRAPHS_TAB;

    @Field(() => NODE_TYPE)
    nodeType: NODE_TYPE;
}

@ObjectType()
export class MetricValues {
    @Field()
    x: number;

    @Field()
    y: string;
}

@ObjectType()
export class MetricInfo {
    @Field()
    org: string;
}

@ObjectType()
export class MetricServiceRes {
    @Field(() => MetricInfo)
    metric: MetricInfo;

    @Field(() => [MetricValues])
    value: MetricValues[];
}
