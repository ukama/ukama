import { Field, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";

@ObjectType()
export class NodeDto {
    @Field()
    id: string;

    @Field()
    title: string;

    @Field()
    description: string;

    @Field()
    totalUser: number;
}

@ObjectType()
export class NodesResponse extends PaginationResponse {
    @Field(() => [NodeDto])
    nodes: NodeDto[];
}

@ObjectType()
export class NodeResponse {
    @Field()
    status: string;

    @Field(() => [NodeDto])
    data: NodeDto[];

    @Field()
    length: number;
}
