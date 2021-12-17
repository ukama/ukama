import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    NodeResponseDto,
    NodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDto,
} from "./types";
import * as defaultCasual from "casual";

class NodeMapper implements INodeMapper {
    dtoToDto = (req: NodeResponse): NodeResponseDto => {
        const nodes = req.data;
        let activeNodes = 0;
        const totalNodes = req.length;
        req.data.forEach(node => {
            if (node.status === ORG_NODE_STATE.ONBOARDED) {
                activeNodes++;
            }
        });
        return { nodes, activeNodes, totalNodes };
    };
    dtoToNodesDto = (req: OrgNodeResponse): OrgNodeResponseDto => {
        const orgName = req.orgName;
        const nodesObj = req.nodes;
        let activeNodes = 0;
        const nodes: NodeDto[] = [];
        const totalNodes = nodesObj.length;
        nodesObj.forEach(node => {
            if (node.state === ORG_NODE_STATE.ONBOARDED) {
                activeNodes++;
            }
            const nodeObj = {
                id: node.nodeId,
                status: node.state,
                title: defaultCasual._title(),
                description: `${defaultCasual.random_value(NODE_TYPE)} node`,
                totalUser: defaultCasual.integer(1, 99),
            };
            nodes.push(nodeObj);
        });
        return { orgName, nodes, activeNodes, totalNodes };
    };
}
export default <INodeMapper>new NodeMapper();
