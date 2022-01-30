import { NODE_TYPE, ORG_NODE_STATE } from "../../constants";
import { INodeMapper } from "./interface";
import {
    NodeResponseDto,
    NodeResponse,
    OrgNodeResponse,
    OrgNodeResponseDto,
    NodeDto,
    CpuUsageMetricsResponse,
    CpuUsageMetricsDto,
    NodeRFDto,
    NodeRFDtoResponse,
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
                const nodeObj = this.getNode(node.nodeId, node.state);
                nodes.push(nodeObj);
            });
        } else {
            nodesObj = this.getNode(
                defaultCasual._uuid(),
                defaultCasual.random_value(ORG_NODE_STATE)
            );
            nodes.push(nodesObj);
        }
        const totalNodes = nodes.length;
        return { orgName, nodes, activeNodes, totalNodes };
    };
    dtoToCpuUsageMetricsDto = (
        res: CpuUsageMetricsResponse
    ): CpuUsageMetricsDto[] => {
        const cpuUsageMetrics: CpuUsageMetricsDto[] = [];
        for (const metric of res.data) {
            const metricObj = {
                id: metric.id,
                usage: metric.usage,
                timestamp: metric.timestamp,
            };
            cpuUsageMetrics.push(metricObj);
        }

        return cpuUsageMetrics;
    };
    dtoToNodeRFKPIDto = (res: NodeRFDtoResponse): NodeRFDto[] => {
        const cpuUsageMetrics: NodeRFDto[] = [];
        for (const metric of res.data) {
            const metricObj = {
                qam: metric.qam,
                rfOutput: metric.rfOutput,
                rssi: metric.rssi,
                timestamp: metric.timestamp,
            };
            cpuUsageMetrics.push(metricObj);
        }

        return cpuUsageMetrics;
    };
    private getNode = (id: string, status: ORG_NODE_STATE): NodeDto => {
        return {
            id: id,
            status: status,
            title: defaultCasual._title(),
            description: `${defaultCasual.random_value(NODE_TYPE)} node`,
            totalUser: defaultCasual.integer(1, 99),
        };
    };
}
export default <INodeMapper>new NodeMapper();
