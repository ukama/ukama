import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import {
    MetricsByTabInputDTO,
    MetricsInputDTO,
    ParsedCookie,
} from "../../common/types";
import setupLogger from "../../config/setupLogger";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER, getMetricUri } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { getMetricTitleByType, getMetricsByTab } from "../../utils";
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
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";

const logger = setupLogger("service");
@Service()
export class NodeService implements INodeService {
    addNode = async (
        req: AddNodeDto,
        cookie: ParsedCookie
    ): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_NODE_API_URL}`,
            headers: cookie.header,
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
        cookie: ParsedCookie
    ): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes/${req.nodeId}`,
            headers: cookie.header,
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
        cookie: ParsedCookie
    ): Promise<UpdateNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${req.nodeId}`,
            headers: cookie.header,
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
        cookie: ParsedCookie
    ): Promise<DeleteNodeRes> => {
        const res = await catchAsyncIOMethod({
            headers: cookie.header,
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${nodeId}`,
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }

        return res;
    };
    getNodesByOrg = async (
        cookie: ParsedCookie
    ): Promise<OrgNodeResponseDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes`,
            headers: cookie.header,
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToNodesDto(cookie.orgId, res);
    };
    getNode = async (
        nodeId: string,
        cookie: ParsedCookie
    ): Promise<NodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NODE_API_URL}/${nodeId}`,
            headers: cookie.header,
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }
        return NodeMapper.dtoToGetNodeDto(res);
    };
    getNodeStatus = async (
        data: GetNodeStatusInput,
        cookie: ParsedCookie
    ): Promise<GetNodeStatusRes> => {
        const currentTimestamp = Math.floor(new Date().getTime() / 1000);
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: cookie.header,
            path:
                getMetricUri(
                    cookie.orgId,
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
        cookie: ParsedCookie,
        endpoint: string
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: cookie.header,
            path: getMetricUri(cookie.orgId, data.nodeId, endpoint),
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
    getMultipleMetrics = async (
        data: MetricsByTabInputDTO,
        cookie: ParsedCookie,
        endpoints: string[]
    ): Promise<MetricRes[]> => {
        return Promise.all(
            endpoints.map(endpoint =>
                catchAsyncIOMethod({
                    type: API_METHOD_TYPE.GET,
                    headers: cookie.header,
                    path: getMetricUri(cookie.orgId, data.nodeId, endpoint),
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
