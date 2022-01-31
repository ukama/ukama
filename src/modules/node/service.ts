import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    CpuUsageMetricsResponse,
    GraphDto,
    NodeDetailDto,
    NodeMetaDataDto,
    NodePhysicalHealthDto,
    NodeRFDtoResponse,
    NodesResponse,
    OrgNodeResponseDto,
    MemoryUsageMetricsResponse,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";
import { INodeService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { HeaderType, PaginationDto } from "../../common/types";
import NodeMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
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
    getNodeGraph = async (): Promise<[GraphDto]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_GRAPH,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };
    cpuUsageMetrics = async (
        req: PaginationDto
    ): Promise<CpuUsageMetricsResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CPU_USAGE_METRICS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const data = NodeMapper.dtoToCpuUsageMetricsDto(res);
        if (!data) throw new HTTP404Error(Messages.ERR_USER_METRICS_NOT_FOUND);
        return {
            data,
            meta,
        };
    };
    nodeRF = async (req: PaginationDto): Promise<NodeRFDtoResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODE_RF_KPI,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const data = NodeMapper.dtoToNodeRFKPIDto(res);
        if (!data) throw new HTTP404Error(Messages.ERR_USER_METRICS_NOT_FOUND);
        return {
            data,
            meta,
        };
    };
    memoryUsageMetrics = async (
        req: PaginationDto
    ): Promise<MemoryUsageMetricsResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_MEMORY_USAGE_METRICS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const data = NodeMapper.dtoToMemoryUsageMetricsDto(res);
        if (!data) throw new HTTP404Error(Messages.ERR_USER_METRICS_NOT_FOUND);
        return {
            data,
            meta,
        };
    };
}
