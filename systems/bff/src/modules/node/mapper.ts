import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDto,
    MetricDto,
    OrgMetricValueDto,
    GetNodeStatusRes,
    NodeResponse,
} from "./types";
import * as defaultCasual from "casual";
import { MetricLatestValueRes } from "../../common/types";

class NodeMapper implements INodeMapper {
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
                    type: node.type,
                    name: node.name,
                });
                nodes.push(nodeObj);
            });
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

    dtoToNodeStatusDto = (res: MetricLatestValueRes): GetNodeStatusRes => {
        let uptime = 0;
        let status = ORG_NODE_STATE.UNDEFINED;
        if (res) {
            uptime = parseFloat(res.value);
            if (uptime > 0) {
                status = ORG_NODE_STATE.ONBOARDED;
            } else {
                status = ORG_NODE_STATE.PENDING;
            }
        }

        return { uptime: uptime, status: status };
    };
    dtoToGetNodeDto = (res: NodeResponse): NodeResponse => {
        const isTowerNode = res.nodeId.includes("tnode");
        return { ...res, attached: isTowerNode ? res.attached : [] };
    };
}
export default <INodeMapper>new NodeMapper();
