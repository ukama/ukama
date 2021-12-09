import { GET_STATUS_TYPE, NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    NodeResponseDto,
    NodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
} from "./types";
import * as defaultCasual from "casual";

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
    dtoToNodesDto = (req: OrgNodeResponse): OrgNodeResponseDto => {
        const orgName = req.orgName;
        const nodes = req.nodes;
        let activeNodes = 0;
        const totalNodes = nodes.length;
        nodes.forEach(node => {
            if (node.state === ORG_NODE_STATE.ONBOARDED) {
                activeNodes++;
            }
            node.title = defaultCasual._title();
            node.description = `${defaultCasual.random_value(NODE_TYPE)} node`;
            node.totalUser = defaultCasual.integer(1, 99);
        });
        return { orgName, nodes, activeNodes, totalNodes };
    };
}
export default <INodeMapper>new NodeMapper();
