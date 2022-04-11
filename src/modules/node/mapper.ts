import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    NodeResponseDto,
    NodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDto,
    MetricDto,
    OrgMetricValueDto,
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
    dtoToNodesDto = (
        orgId: string,
        req: OrgNodeResponse
    ): OrgNodeResponseDto => {
        let nodesObj;
        let activeNodes = 0;
        const nodes: NodeDto[] = [];
        if (req.nodes) {
            nodesObj = req.nodes;

            nodesObj.forEach(node => {
                if (node.state === ORG_NODE_STATE.ONBOARDED) {
                    activeNodes++;
                }
                const nodeObj = this.getNode({
                    id: node.nodeId,
                    status: node.state,
                    type: node.type as NODE_TYPE,
                    name: node.name,
                });
                nodes.push(nodeObj);
            });
        } else {
            nodesObj = this.getNode({});
            nodes.push(nodesObj);
        }
        const totalNodes = nodes.length;
        return { orgId, nodes, activeNodes, totalNodes };
    };

    dtoToMetricsDto = (res: OrgMetricValueDto[]): MetricDto[] => {
        const metrics: MetricDto[] = [];
        if (res && res.length > 0)
            res.forEach((item: any) =>
                metrics.push({
                    x: item[0],
                    y: item[1],
                })
            );
        return metrics;
    };

    private getNode = ({
        id = defaultCasual._uuid(),
        name = defaultCasual._title(),
        status = defaultCasual.random_value(ORG_NODE_STATE),
        type = NODE_TYPE.TOWER,
    }: {
        id?: string;
        name?: string;
        status?: ORG_NODE_STATE;
        type?: NODE_TYPE;
    }): NodeDto => {
        return {
            id: id,
            type: type,
            status: status,
            name: name ? name : defaultCasual._title(),
            description: `${type} node`,
            totalUser: defaultCasual.integer(1, 99),
            isUpdateAvailable: Math.random() < 0.7,
            updateShortNote:
                "Software update available. Estimated 10 minutes, and will (be/not be) disruptive. ",
            updateDescription:
                "Short introduction.\n\n TL;DR\n\n*** NEW ***\nPoint 1\nPoint 2\nPoint 3\n\n*** IMPROVEMENTS ***\nPoint 1\nPoint 2\nPoint 3\n\n*** FIXES ***\nPoint 1\nPoint 2\nPoint 3\n\nWe would love to hear your feedback -- if you have anything to share, please xyz.",
            updateVersion: "12.4",
        };
    };
}
export default <INodeMapper>new NodeMapper();
