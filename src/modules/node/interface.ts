import { PaginationDto } from "../../common/types";
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
}

export interface INodeMapper {
    dtoToDto(res: NodeResponse): NodeResponseDto;
}
