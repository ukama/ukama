import { Service } from "typedi";
import { NodeResponse, NodesResponse } from "./types";
import { INodeService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import NodeMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";

@Service()
export class NodeService implements INodeService {
    getNodes = async (req: PaginationDto): Promise<NodesResponse> => {
        const res = await catchAsyncIOMethod<NodeResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODES,
            params: req,
        });
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const nodes = NodeMapper.dtoToDto(res);

        if (!nodes) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            nodes,
            meta,
        };
    };
}
