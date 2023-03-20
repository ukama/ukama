import {
    MetricDto,
    AddNodeDto,
    NodeResponse,
    AddNodeResponse,
    OrgNodeResponse,
    GetNodeStatusRes,
    UpdateNodeResponse,
    UpdateNodeDto,
    OrgMetricValueDto,
    OrgNodeResponseDto,
    GetNodeStatusInput,
} from "./types";
import { DeleteNodeRes } from "../user/types";
import { MetricLatestValueRes, ParsedCookie } from "../../common/types";

export interface INodeService {
    getNodeStatus(
        data: GetNodeStatusInput,
        cookie: ParsedCookie,
    ): Promise<GetNodeStatusRes>;
    getNode(nodeId: string, cookie: ParsedCookie): Promise<NodeResponse>;
    getNodesByOrg(cookie: ParsedCookie): Promise<OrgNodeResponseDto>;
    addNode(req: AddNodeDto, cookie: ParsedCookie): Promise<AddNodeResponse>;
    updateNode(
        req: UpdateNodeDto,
        cookie: ParsedCookie
    ): Promise<UpdateNodeResponse>;
    deleteNode(id: string, cookie: ParsedCookie): Promise<DeleteNodeRes>;
}

export interface INodeMapper {
    dtoToGetNodeDto(res: NodeResponse): NodeResponse;
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToMetricsDto(res: OrgMetricValueDto[]): MetricDto[];
    dtoToNodeStatusDto(res: MetricLatestValueRes): GetNodeStatusRes;
    dtoToNodeResponsedto(res: any): AddNodeResponse;
}
