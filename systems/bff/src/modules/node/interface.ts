import { MetricLatestValueRes, THeaders } from "../../common/types";
import { DeleteNodeRes } from "../user/types";
import {
    AddNodeDto,
    AddNodeResponse,
    GetNodeStatusInput,
    GetNodeStatusRes,
    MetricDto,
    NodeResponse,
    OrgMetricValueDto,
    OrgNodeResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";

export interface INodeService {
    getNodeStatus(
        data: GetNodeStatusInput,
        headers: THeaders
    ): Promise<GetNodeStatusRes>;
    getNode(nodeId: string, headers: THeaders): Promise<NodeResponse>;
    getNodesByOrg(cookie: THeaders): Promise<OrgNodeResponseDto>;
    addNode(req: AddNodeDto, headers: THeaders): Promise<AddNodeResponse>;
    updateNode(
        req: UpdateNodeDto,
        headers: THeaders
    ): Promise<UpdateNodeResponse>;
    deleteNode(id: string, headers: THeaders): Promise<DeleteNodeRes>;
}

export interface INodeMapper {
    dtoToGetNodeDto(res: NodeResponse): NodeResponse;
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToMetricsDto(res: OrgMetricValueDto[]): MetricDto[];
    dtoToNodeStatusDto(res: MetricLatestValueRes): GetNodeStatusRes;
    dtoToNodeResponsedto(res: any): AddNodeResponse;
}
