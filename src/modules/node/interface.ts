import { HeaderType, PaginationDto } from "../../common/types";
import { NetworkDto } from "../network/types";
import { DeactivateResponse } from "../user/types";
import {
    NodeResponseDto,
    NodeResponse,
    NodesResponse,
    AddNodeResponse,
    AddNodeDto,
    UpdateNodeDto,
    UpdateNodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDetailDto,
    NodeMetaDataDto,
    NodePhysicalHealthDto,
    NodeRFDto,
    ThroughputMetricsDto,
    CpuUsageMetricsDto,
    CpuUsageMetricsResponse,
    NodeRFDtoResponse,
    TemperatureMetricsDto,
    TemperatureMetricsResponse,
    IOMetricsResponse,
    IOMetricsDto,
} from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
    getNetwork(): Promise<NetworkDto>;
    getNodeDetials(): Promise<NodeDetailDto>;
    nodeMetaData(): Promise<NodeMetaDataDto>;
    nodePhysicalHealth(): Promise<NodePhysicalHealthDto>;
    getThroughputMetrics(): Promise<[ThroughputMetricsDto]>;
    getNodesByOrg(
        orgId: string,
        header: HeaderType
    ): Promise<OrgNodeResponseDto>;
    addNode(req: AddNodeDto): Promise<AddNodeResponse>;
    updateNode(req: UpdateNodeDto): Promise<UpdateNodeResponse>;
    deleteNode(id: string): Promise<DeactivateResponse>;
}

export interface INodeMapper {
    dtoToDto(res: NodeResponse): NodeResponseDto;
    dtoToNodeRFKPIDto(req: NodeRFDtoResponse): NodeRFDto[];
    dtoToIOMetricsDto(req: IOMetricsResponse): IOMetricsDto[];
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToCpuUsageMetricsDto(req: CpuUsageMetricsResponse): CpuUsageMetricsDto[];
    dtoToTemperatureMetricsDto(
        req: TemperatureMetricsResponse
    ): TemperatureMetricsDto[];
}
