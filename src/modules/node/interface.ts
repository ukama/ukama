import { PaginationDto } from "../../common/types";
import { DeactivateResponse } from "../user/types";
import {
    NodeResponseDto,
    NodeResponse,
    NodesResponse,
    AddNodeResponse,
    AddNodeDto,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
    addNode(req: AddNodeDto): Promise<AddNodeResponse>;
    updateNode(req: UpdateNodeDto): Promise<UpdateNodeResponse>;
    deleteNode(id: string): Promise<DeactivateResponse>;
}

export interface INodeMapper {
    dtoToDto(res: NodeResponse): NodeResponseDto;
}
