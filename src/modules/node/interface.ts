import { PaginationDto } from "../../common/types";
import { NodeDto, NodeResponse, NodesResponse } from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
}

export interface INodeMapper {
    dtoToDto(res: NodeResponse): NodeDto[];
}
