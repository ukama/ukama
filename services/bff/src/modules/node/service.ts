import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    MetricDto,
    NodeAppsVersionLogsResponse,
    NodeAppResponse,
    MetricRes,
    OrgNodeDto,
    NodeResponse,
    GetNodeStatusRes,
    GetNodeStatusInput,
} from "./types";
import {
    ParsedCookie,
    MetricsInputDTO,
    MetricsByTabInputDTO,
} from "../../common/types";
import NodeMapper from "./mapper";
import { checkError } from "../../errors";
import { INodeService } from "./interface";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { DeactivateResponse } from "../user/types";
import { GRAPHS_TAB } from "./../../constants/index";
import { getMetricUri, SERVER } from "../../constants/endpoints";
import { getMetricsByTab, getMetricTitleByType } from "../../utils";

@Service()
export class NodeService implements INodeService {
    addNode = async (
        req: AddNodeDto,
        cookie: ParsedCookie
    ): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes/${req.nodeId}`,
            headers: cookie.header,
            body: {
                name: req.name,
            },
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    updateNode = async (
        req: UpdateNodeDto,
        cookie: ParsedCookie
    ): Promise<OrgNodeDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes/${req.nodeId}`,
            headers: cookie.header,
            body: {
                name: req.name,
            },
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    deleteNode = async (id: string): Promise<DeactivateResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_DELETE_NODE,
            body: { id },
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };
    getNodesByOrg = async (
        cookie: ParsedCookie
    ): Promise<OrgNodeResponseDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes`,
            headers: cookie.header,
        });
        return NodeMapper.dtoToNodesDto(cookie.orgId, res);
    };
    getNode = async (
        nodeId: string,
        cookie: ParsedCookie
    ): Promise<NodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes/${nodeId}`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
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
        if (checkError(res)) throw new Error(res.message);

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
        if (checkError(res)) throw new Error(res.message);
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
