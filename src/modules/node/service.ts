import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    ThroughputMetricsDto,
    NodeDetailDto,
    NodeMetaDataDto,
    NodePhysicalHealthDto,
    NodesResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
    CpuUsageMetricsDto,
    IOMetricsDto,
    NodeRFDto,
    TemperatureMetricsDto,
    MemoryUsageMetricsDto,
    MetricDto,
} from "./types";
import { INodeService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { HeaderType, MetricsInputDTO, PaginationDto } from "../../common/types";
import NodeMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE, GRAPH_FILTER } from "../../constants";
import { SERVER } from "../../constants/endpoints";
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

    nodeMetaData = async (): Promise<NodeMetaDataDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_META_DATA,
        });
        return res.data;
    };
    nodePhysicalHealth = async (): Promise<NodePhysicalHealthDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_PHYSICAL_HEALTH,
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
    getThroughputMetrics = async (
        filter: GRAPH_FILTER
    ): Promise<[ThroughputMetricsDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_THROUGHPUT_METRICS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };
    cpuUsageMetrics = async (
        filter: GRAPH_FILTER
    ): Promise<[CpuUsageMetricsDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CPU_USAGE_METRICS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    nodeRF = async (filter: GRAPH_FILTER): Promise<[NodeRFDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_RF_KPI,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    temperatureMetrics = async (
        filter: GRAPH_FILTER
    ): Promise<[TemperatureMetricsDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_TEMPERATURE_METRICS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    ioMetrics = async (filter: GRAPH_FILTER): Promise<[IOMetricsDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_IO_METRICS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    memoryUsageMetrics = async (
        filter: GRAPH_FILTER
    ): Promise<[MemoryUsageMetricsDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_MEMORY_USAGE_METRICS,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    metricsCpuTRX = async (
        data: MetricsInputDTO,
        header: HeaderType
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: `${SERVER.ORG}/${data.orgId}/nodes/${data.nodeId}/metrics/cpu`,
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) throw new Error(res.message);
        return NodeMapper.dtoToMetricDto(res.data?.result);
    };
    metricsMemoryTRX = async (
        data: MetricsInputDTO,
        header: HeaderType
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: `${SERVER.ORG}/${data.orgId}/nodes/${data.nodeId}/metrics/memory`,
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) throw new Error(res.message);
        return NodeMapper.dtoToMetricDto(res.data?.result);
    };
    getMetricsUptime = async (
        data: MetricsInputDTO,
        header: HeaderType
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: `${SERVER.ORG}/${data.orgId}/nodes/${data.nodeId}/metrics/uptime`,
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) throw new Error(res.message);
        return NodeMapper.dtoToMetricDto(res.data?.result);
    };
    getThroughputUL = async (
        data: MetricsInputDTO,
        header: HeaderType
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: `${SERVER.ORG}/${data.orgId}/nodes/${data.nodeId}/metrics/throughputuplink`,
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) throw new Error(res.message);
        return NodeMapper.dtoToMetricDto(res.data?.result);
    };
    getThroughputDL = async (
        data: MetricsInputDTO,
        header: HeaderType
    ): Promise<MetricDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            headers: header,
            path: `${SERVER.ORG}/${data.orgId}/nodes/${data.nodeId}/metrics/throughputdownlink`,
            params: { from: data.from, to: data.to, step: data.step },
        });
        if (checkError(res)) throw new Error(res.message);
        return NodeMapper.dtoToMetricDto(res.data?.result);
    };
}
