import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    NodeDetailDto,
    NodesResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
    MetricDto,
    NodeAppsVersionLogsResponse,
    NodeAppResponse,
    MetricRes,
} from "./types";
import { INodeService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import {
    HeaderType,
    MetricsByTabInputDTO,
    MetricsInputDTO,
    PaginationDto,
} from "../../common/types";
import NodeMapper from "./mapper";
import { getMetricTitleByType, getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { getMetricUri, SERVER } from "../../constants/endpoints";
import { DeactivateResponse } from "../user/types";
import { NetworkDto } from "../network/types";

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

    addNode = async (req: AddNodeDto): Promise<AddNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_ADD_NODE,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    updateNode = async (req: UpdateNodeDto): Promise<UpdateNodeResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_UPDATE_NODE,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
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
        orgId: string,
        header: HeaderType
    ): Promise<OrgNodeResponseDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${orgId}/nodes`,
            headers: header,
        });
        return NodeMapper.dtoToNodesDto(orgId, res);
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
        header: HeaderType,
        endpoint: string
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: getMetricUri(data.orgId, data.nodeId, endpoint),
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
        header: HeaderType,
        endpoints: string[]
    ): Promise<MetricRes[]> => {
        return Promise.all(
            endpoints.map(endpoint =>
                catchAsyncIOMethod({
                    type: API_METHOD_TYPE.GET,
                    headers: header,
                    path: getMetricUri(data.orgId, data.nodeId, endpoint),
                    params: {
                        to: data.to,
                        from: data.from,
                        step: data.step,
                    },
                }).then(res => {
                    if (checkError(res)) {
                        console.error(res.message);
                        throw new Error(res.message);
                    } else {
                        const values = res.data.result[0];
                        return {
                            type: endpoint,
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
