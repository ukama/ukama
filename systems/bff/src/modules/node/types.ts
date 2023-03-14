import { Field, InputType, ObjectType } from "type-graphql";
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
export class AddNodeResponse {
    @Field()
    success: boolean;
}

@ObjectType()
export class AttachedNodes {
    @Field()
    nodeId: string;
}

@ObjectType()
export class LinkNodes {
    @Field()
    nodeId: string;

    @Field()
    attachedNodeIds: string[];
}

@InputType()
export class NodeObj {
    @Field()
    name: string;

    @Field()
    nodeId: string;
    @Field()
    type: string;

    @Field()
    state: string;

    @Field(() => [NodeObj], { nullable: true })
    attached?: NodeObj[];
}

@InputType()
export class AddNodeDto {
    @Field()
    name: string;

    @Field()
    nodeId: string;

    @Field()
    type: string;

    @Field()
    state: string;

    @Field(() => [NodeObj])
    attached?: NodeObj[];
}

@InputType()
export class UpdateNodeDto {
    @Field()
    nodeId: string;

    @Field()
    name: string;
}

@ObjectType()
export class OrgNodeDto {
    @Field()
    nodeId: string;

    @Field(() => NODE_TYPE)
    type: NODE_TYPE;

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
export class NodeResponse {
    @Field()
    nodeId: string;

    @Field(() => NODE_TYPE)
    type: NODE_TYPE;

    @Field(() => ORG_NODE_STATE)
    state: ORG_NODE_STATE;

    @Field()
    name: string;

    @Field(() => [OrgNodeDto])
    attached?: OrgNodeDto[];
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

@InputType()
export class GetNodeStatusInput {
    @Field()
    nodeId: string;

    @Field(() => NODE_TYPE)
    nodeType: NODE_TYPE;
}

@ObjectType()
export class GetNodeStatusRes {
    @Field()
    uptime: number;

    @Field(() => ORG_NODE_STATE)
    status: ORG_NODE_STATE;
}
