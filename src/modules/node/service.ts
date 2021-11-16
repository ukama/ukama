import { Service } from "typedi";
import { NodesResponse } from "./types";
import { INodeService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import NodeMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import NodeIOMethods from "./io";

@Service()
export class NodeService implements INodeService {
    getNodes = async (req: PaginationDto): Promise<NodesResponse> => {
        const res = await NodeIOMethods.getNodesMethod(req);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        const meta = getPaginatedOutput(
            req.pageNo,
            req.pageSize,
            res.data.length
        );
        const nodes = NodeMapper.dtoToDto(res.data.data);

        if (!nodes) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            nodes,
            meta,
        };
    };
}
