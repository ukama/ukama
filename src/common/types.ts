import "reflect-metadata";
import { Field, InputType, ObjectType } from "type-graphql";

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
export class GETDataDto {
    @Field()
    path: string;

    @Field(() => String || Object || null, { nullable: true })
    params?: any;

    @Field(() => String || Object || null, { nullable: true })
    headers?: any;
}
