import { HeaderType, PaginationDto } from "../../common/types";
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
} from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
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
}
