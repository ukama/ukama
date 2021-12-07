import { GET_STATUS_TYPE } from "../../constants";
import { INodeMapper } from "./interface";
import { NodeResponseDto, NodeResponse } from "./types";

class NodeMapper implements INodeMapper {
    dtoToDto = (req: NodeResponse): NodeResponseDto => {
        const nodes = req.data;
        let activeNodes = 0;
        const totalNodes = req.length;
        req.data.forEach(node => {
            if (node.status === GET_STATUS_TYPE.ACTIVE) {
                activeNodes++;
            }
        });
        return { nodes, activeNodes, totalNodes };
    };
}
export default <INodeMapper>new NodeMapper();
