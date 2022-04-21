import { ParsedCookie, PaginationDto } from "../../common/types";
import { NetworkDto } from "../network/types";
import { DeactivateResponse } from "../user/types";
import {
    NodeResponseDto,
    NodeResponse,
    NodesResponse,
    AddNodeResponse,
    AddNodeDto,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDetailDto,
    MetricDto,
    OrgMetricValueDto,
    OrgNodeDto,
} from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
    getNetwork(): Promise<NetworkDto>;
    getNodeDetials(): Promise<NodeDetailDto>;
    getNodesByOrg(cookie: ParsedCookie): Promise<OrgNodeResponseDto>;
    addNode(req: AddNodeDto, cookie: ParsedCookie): Promise<AddNodeResponse>;
    updateNode(req: AddNodeDto, cookie: ParsedCookie): Promise<OrgNodeDto>;
    deleteNode(id: string): Promise<DeactivateResponse>;
}

export interface INodeMapper {
    dtoToDto(res: NodeResponse): NodeResponseDto;
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToMetricsDto(res: OrgMetricValueDto[]): MetricDto[];
}
