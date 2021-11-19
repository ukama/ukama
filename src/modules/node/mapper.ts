import { INodeMapper } from "./interface";
import { NodeDto, NodeResponse } from "./types";

class NodeMapper implements INodeMapper {
    dtoToDto = (res: NodeResponse): NodeDto[] => {
        const nodes = res.data;
        return nodes;
    };
}
export default <INodeMapper>new NodeMapper();
