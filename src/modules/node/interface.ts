import { PaginationDto } from "../../common/types";
import { NodeDto, NodesResponse } from "./types";

export interface INodeService {
    getNodes(req: PaginationDto): Promise<NodesResponse>;
}

export interface INodeMapper {
    dtoToDto(data: NodeDto[]): NodeDto[];
}
