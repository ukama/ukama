import { NetworkDto } from "../network/types";
import { ParsedCookie } from "../../common/types";
import { DeactivateResponse } from "../user/types";
import {
    AddNodeResponse,
    AddNodeDto,
    OrgNodeResponse,
    OrgNodeResponseDto,
    MetricDto,
    OrgMetricValueDto,
    OrgNodeDto,
    NodeResponse,
} from "./types";

export interface INodeService {
    getNetwork(): Promise<NetworkDto>;
    getNode(nodeId: string, cookie: ParsedCookie): Promise<NodeResponse>;
    getNodesByOrg(cookie: ParsedCookie): Promise<OrgNodeResponseDto>;
    addNode(req: AddNodeDto, cookie: ParsedCookie): Promise<AddNodeResponse>;
    updateNode(req: AddNodeDto, cookie: ParsedCookie): Promise<OrgNodeDto>;
    deleteNode(id: string): Promise<DeactivateResponse>;
}

export interface INodeMapper {
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToMetricsDto(res: OrgMetricValueDto[]): MetricDto[];
}
