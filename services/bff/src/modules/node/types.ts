import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";
import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";

@ObjectType()
export class NodeDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    type: string;

    @Field()
    description: string;

    @Field(() => ORG_NODE_STATE)
    status: ORG_NODE_STATE;

    @Field()
    totalUser: number;

    @Field()
    isUpdateAvailable: boolean;

    @Field()
    updateVersion: string;

    @Field()
    updateShortNote: string;

    @Field()
    updateDescription: string;
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
    nodeId: string;

    @Field()
    name: string;

    @Field(() => ORG_NODE_STATE)
    state: ORG_NODE_STATE;

    @Field(() => NODE_TYPE)
    type: NODE_TYPE;
}

@InputType()
export class AddNodeDto {
    @Field()
    name: string;

    @Field()
    nodeId: string;

    @Field()
    orgId: string;
}

@InputType()
export class UpdateNodeDto {
    @Field()
    orgId: string;

    @Field()
    nodeId: string;

    @Field()
    name: string;
}

@ObjectType()
export class OrgNodeDto {
    @Field()
    nodeId: string;

    @Field()
    type: string;

    @Field(() => ORG_NODE_STATE)
    state: ORG_NODE_STATE;

    @Field()
    name: string;
}

@ObjectType()
export class OrgNodeResponse {
    @Field()
    orgName: string;

    @Field(() => [OrgNodeDto])
    nodes: OrgNodeDto[];
}

@ObjectType()
export class OrgNodeResponseDto {
    @Field()
    orgId: string;

    @Field(() => [NodeDto])
    nodes: NodeDto[];

    @Field()
    activeNodes: number;

    @Field()
    totalNodes: number;
}

@ObjectType()
export class NodeDetailDto {
    @Field()
    id: string;

    @Field()
    modelType: string;

    @Field()
    serial: number;

    @Field()
    macAddress: number;

    @Field()
    osVersion: number;

    @Field()
    manufacturing: number;

    @Field()
    ukamaOS: number;

    @Field()
    hardware: number;

    @Field()
    description: string;
}
@ObjectType()
export class MetricDto {
    @Field()
    y: number;

    @Field()
    x: number;
}

@ObjectType()
export class OrgMetricValueDto {
    @Field()
    x: number;

    @Field()
    y: string;
}
@ObjectType()
export class NodeAppsVersionLogsResponse {
    @Field()
    version: string;

    @Field()
    date: number;

    @Field()
    notes: string;
}
@ObjectType()
export class NodeAppResponse {
    @Field()
    id: string;

    @Field()
    title: string;

    @Field()
    version: string;

    @Field()
    cpu: string;

    @Field()
    memory: string;
}

@ObjectType()
export class MetricRes {
    @Field()
    type: string;

    @Field()
    name: string;

    @Field(() => [MetricDto])
    data: MetricDto[];

    @Field()
    next: boolean;
}

@ObjectType()
export class GetMetricsRes {
    @Field()
    to: number;

    @Field()
    next: boolean;

    @Field(() => [MetricRes])
    metrics: MetricRes[];
}
