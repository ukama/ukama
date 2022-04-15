import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    NodeDetailDto,
    NodesResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    MetricDto,
    NodeAppsVersionLogsResponse,
    NodeAppResponse,
    MetricRes,
    OrgNodeDto,
} from "./types";
import {
    ParsedCookie,
    MetricsByTabInputDTO,
    MetricsInputDTO,
    PaginationDto,
} from "../../common/types";
import NodeMapper from "./mapper";
import { INodeService } from "./interface";
import { NetworkDto } from "../network/types";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { DeactivateResponse } from "../user/types";
import { getMetricUri, SERVER } from "../../constants/endpoints";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { getMetricTitleByType, getPaginatedOutput } from "../../utils";

@Service()
export class NodeService implements INodeService {
    getNodes = async (req: PaginationDto): Promise<NodesResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODES,
            params: req,
        });

        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const nodes = NodeMapper.dtoToDto(res);

        if (!nodes) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            nodes,
            meta,
        };
    };

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
    getNodeDetials = async (): Promise<NodeDetailDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_DETAIL,
        });
        return res.data;
    };
    getNetwork = async (): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_NETWORK,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
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
