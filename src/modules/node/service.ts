import { Service } from "typedi";
import {
    AddNodeDto,
    AddNodeResponse,
    NodesResponse,
    OrgNodeResponseDto,
    UpdateNodeDto,
    UpdateNodeResponse,
} from "./types";
import { INodeService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { Context, PaginationDto } from "../../common/types";
import NodeMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod, getHeaders } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { DeactivateResponse } from "../user/types";

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
        ctx: Context
    ): Promise<OrgNodeResponseDto> => {
        const header = getHeaders(ctx);
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.GET_NODES_BY_ORG}/${orgId}/nodes`,
            headers: header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);

        return NodeMapper.dtoToNodesDto(res);
    };
}
