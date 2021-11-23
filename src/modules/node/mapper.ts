import { GET_STATUS_TYPE } from "../../constants";
import { INodeMapper } from "./interface";
import { NodeResponseDto, NodeResponse } from "./types";

class NodeMapper implements INodeMapper {
    dtoToDto = (res: NodeResponse): NodeResponseDto => {
        const nodes = res.data;
        let activeNodes = 0;
        const totalNodes = res.length;
        res.data.map(node => {
            if (node.status === GET_STATUS_TYPE.ACTIVE) {
                activeNodes++;
            }
        });
        return { nodes, activeNodes, totalNodes };
    };
}
export default <INodeMapper>new NodeMapper();
