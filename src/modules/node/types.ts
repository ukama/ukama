import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";
import { GET_STATUS_TYPE } from "../../constants";

@ObjectType()
export class NodeDto {
    @Field()
    id: string;

    @Field()
    title: string;

    @Field()
    description: string;

    @Field(() => GET_STATUS_TYPE)
    status: GET_STATUS_TYPE;

    @Field()
    totalUser: number;
}

@ObjectType()
export class NodeResponseDto {
    @Field(() => [NodeDto])
    nodes: NodeDto[];

    @Field()
    activeNodes: number;

    @Field()
    totalNodes: number;
}

@ObjectType()
export class NodesResponse extends PaginationResponse {
    @Field(() => NodeResponseDto)
    nodes: NodeResponseDto;
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

@ObjectType()
export class AddNodeResponse {
    @Field()
    success: boolean;
}

@InputType()
export class AddNodeDto {
    @Field()
    name: string;

    @Field()
    serialNo: string;
}

@ObjectType()
export class AddNodeResponseDto {
    @Field()
    status: string;

    @Field()
    data: AddNodeResponse;
}
