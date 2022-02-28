import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    NodeResponseDto,
    NodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDto,
    MetricDto,
    OrgMetricResponse,
} from "./types";
import * as defaultCasual from "casual";
import { MetricsInputDTO } from "../../common/types";

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
    dtoToNodesDto = (
        orgId: string,
        req: OrgNodeResponse
    ): OrgNodeResponseDto => {
        const orgName = req.orgName ? req.orgName : orgId;
        let nodesObj;
        let activeNodes = 0;
        const nodes: NodeDto[] = [];
        if (req.nodes) {
            nodesObj = req.nodes;

            nodesObj.forEach(node => {
                if (node.state === ORG_NODE_STATE.ONBOARDED) {
                    activeNodes++;
                }
                const nodeObj = this.getNode(
                    node.nodeId,
                    node.state,
                    node.type
                );
                nodes.push(nodeObj);
            });
        } else {
            nodesObj = this.getNode(
                defaultCasual._uuid(),
                defaultCasual.random_value(ORG_NODE_STATE),
                "TOWER"
            );
            nodes.push(nodesObj);
        }
        const totalNodes = nodes.length;
        return { orgName, nodes, activeNodes, totalNodes };
    };

    dtoToMetricDto = (res: OrgMetricResponse[]): MetricDto[] => {
        const metrics: MetricDto[] = [];
        if (res && res.length > 0)
            res[0].values.map((item: any) =>
                metrics.push({
                    x: item[0],
                    y: item[1],
                })
            );
        return metrics;
    };

    private getNode = (
        id: string,
        status: ORG_NODE_STATE,
        type: string
    ): NodeDto => {
        return {
            id: id,
            type: type,
            status: status,
            title: defaultCasual._title(),
            description: `${defaultCasual.random_value(NODE_TYPE)} node`,
            totalUser: defaultCasual.integer(1, 99),
        };
    };
}
export default <INodeMapper>new NodeMapper();
