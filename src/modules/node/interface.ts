import { HeaderType, MetricsInputDTO, PaginationDto } from "../../common/types";
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
    OrgMetricResponse,
    MetricDto,
} from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
    getNetwork(): Promise<NetworkDto>;
    getNodeDetials(): Promise<NodeDetailDto>;
    nodeMetaData(): Promise<NodeMetaDataDto>;
    nodePhysicalHealth(): Promise<NodePhysicalHealthDto>;
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
    dtoToNodesDto(orgId: string, req: OrgNodeResponse): OrgNodeResponseDto;
    dtoToMetricDto(res: OrgMetricResponse[]): MetricDto[];
}
