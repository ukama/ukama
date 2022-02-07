import { Field, InputType, ObjectType } from "type-graphql";
import { PaginationResponse } from "../../common/types";
import { ORG_NODE_STATE } from "../../constants";

@ObjectType()
export class NodeDto {
    @Field()
    id: string;

    @Field()
    title: string;

    @Field()
    description: string;

    @Field(() => ORG_NODE_STATE)
    status: ORG_NODE_STATE;

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

    @Field()
    securityCode: string;
}

@InputType()
export class UpdateNodeDto {
    @Field()
    id: string;

    @Field({ nullable: true })
    name: string;

    @Field({ nullable: true })
    serialNo: string;

    @Field({ nullable: true })
    securityCode: string;
}

@ObjectType()
export class AddNodeResponseDto {
    @Field()
    status: string;

    @Field()
    data: AddNodeResponse;
}

@ObjectType()
export class UpdateNodeResponse {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    serialNo: string;
}

@ObjectType()
export class OrgNodeDto {
    @Field()
    nodeId: string;

    @Field(() => ORG_NODE_STATE)
    state: ORG_NODE_STATE;
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
    orgName: string;

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
export class NodeMetaDataDto {
    @Field()
    throughput: number;

    @Field()
    usersAttached: number;
}

@ObjectType()
export class NodePhysicalHealthDto {
    @Field()
    temperature: number;

    @Field()
    Memory: number;

    @Field()
    cpu: number;

    @Field()
    io: number;
}

@ObjectType()
export class ThroughputMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    amount: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class CpuUsageMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    usage: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class NodeRFDto {
    @Field()
    qam: number;

    @Field()
    rfOutput: number;

    @Field()
    rssi: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class TemperatureMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    temperature: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class IOMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    input: number;

    @Field()
    output: number;

    @Field()
    timestamp: number;
}

@ObjectType()
export class MemoryUsageMetricsDto {
    @Field({ nullable: true })
    id: string;

    @Field()
    usage: number;

    @Field()
    timestamp: number;
}
