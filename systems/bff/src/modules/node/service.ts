import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import {
    MetricsByTabInputDTO,
    MetricsInputDTO,
    THeaders,
} from "../../common/types";
import setupLogger from "../../config/setupLogger";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER, getMetricUri } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { getHeaders, getMetricTitleByType, getMetricsByTab } from "../../utils";
import { DeleteNodeRes } from "../user/types";
import { GRAPHS_TAB } from "./../../constants/index";
import { INodeService } from "./interface";
import NodeMapper from "./mapper";
import {
    AddNodeDto,
    AddNodeResponse,
    GetNodeStatusInput,
    GetNodeStatusRes,
    LinkNodes,
    MetricDto,
    MetricRes,
    NodeAppResponse,
    NodeAppsVersionLogsResponse,
    NodeResponse,
    NodeStatsResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";

const logger = setupLogger("service");
@Service()
export class NodeService implements INodeService {
    addNode = async (
        req: AddNodeDto,
        headers: THeaders
    ): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_NODE_API_URL}`,
            headers: getHeaders(headers),
            body: {
                name: req.name,
                node_id: req.nodeId,
                state: req.state,
            },
        });

        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToNodeResponsedto(res);
    };
    linkNodes = async (
        req: LinkNodes,
        headers: THeaders
    ): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.ORG}/${headers.orgId}/nodes/${req.nodeId}`,
            headers: getHeaders(headers),
            body: {
                attachedNodeIds: req.attachedNodeIds,
            },
        });

        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return res;
    };
    updateNode = async (
        req: UpdateNodeDto,
        headers: THeaders
    ): Promise<UpdateNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${req.nodeId}`,
            headers: getHeaders(headers),
            body: {
                name: req.name,
            },
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToNodeResponsedto(res);
    };
    deleteNode = async (
        nodeId: string,
        headers: THeaders
    ): Promise<DeleteNodeRes> => {
        const res = await catchAsyncIOMethod({
            headers: getHeaders(headers),
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${nodeId}`,
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }

        return res;
    };
    getNodesByOrg = async (headers: THeaders): Promise<OrgNodeResponseDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${headers.orgId}/nodes`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToNodesDto(headers.orgId, res);
    };
    getNode = async (
        nodeId: string,
        headers: THeaders
    ): Promise<NodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${nodeId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToGetNodeDto(res);
    };
    getNodeStatus = async (
        data: GetNodeStatusInput,
        headers: THeaders
    ): Promise<GetNodeStatusRes> => {
        const currentTimestamp = Math.floor(new Date().getTime() / 1000);
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: getHeaders(headers),
            path:
                getMetricUri(
                    headers.orgId,
                    data.nodeId,
                    getMetricsByTab(data.nodeType, GRAPHS_TAB.NODE_STATUS)[0]
                ) + "/latest",
            params: {
                from: currentTimestamp,
                to: currentTimestamp,
                step: 1,
            },
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }

        return NodeMapper.dtoToNodeStatusDto(res);
    };
    getSingleMetric = async (
        data: MetricsInputDTO,
        headers: THeaders,
        endpoint: string
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: getHeaders(headers),
            path: getMetricUri(headers.orgId, data.nodeId, endpoint),
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToMetricsDto(res.data?.result[0]);
    };
    getSoftwareLogs = async (): Promise<NodeAppsVersionLogsResponse[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_SOFTWARE_LOGS,
        });
        return res.data;
    };
    getNodeApps = async (): Promise<NodeAppResponse[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_APPS,
        });
        return res.data;
    };
    getNodesStats = async (): Promise<NodeStatsResponse> => {
        return {
            totalCount: 4,
            upCount: 2,
            claimCount: 2,
        };
    };
    getMultipleMetrics = async (
        data: MetricsByTabInputDTO,
        headers: THeaders,
        endpoints: string[]
    ): Promise<MetricRes[]> => {
        return Promise.all(
            endpoints.map(endpoint =>
                catchAsyncIOMethod({
                    type: API_METHOD_TYPE.GET,
                    headers: getHeaders(headers),
                    path: getMetricUri(headers.orgId, data.nodeId, endpoint),
                    params: {
                        to: data.to,
                        from: data.from,
                        step: data.step,
                    },
                }).then(res => {
                    if (checkError(res)) {
                        return {
                            next: false,
                            type: endpoint,
                            name: getMetricTitleByType(endpoint),
                            data: [],
                        };
                    } else {
                        const values = res.data.result[0];
                        return {
                            type: endpoint,
                            next: res.data.result.length > 0,
                            name: getMetricTitleByType(endpoint),
                            data:
                                res.data.result.length > 0
                                    ? NodeMapper.dtoToMetricsDto(values.values)
                                    : [],
                        };
                    }
                })
            )
        );
    };
}
