import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { HTTP404Error, Messages, checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { IRateService } from "./interface";
import RateMapper from "./mapper";
import {
    DefaultMarkupHistoryResDto,
    DefaultMarkupInputDto,
    DefaultMarkupResDto,
} from "./types";

@Service()
export class RateService implements IRateService {
    defaultMarkup = async (
        req: DefaultMarkupInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/${req.markup}/default`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
    getDefaultMarkup = async (
        headers: THeaders
    ): Promise<DefaultMarkupResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/default`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return RateMapper.dtoToDefaultMarkupDto(res);
    };
    getDefaultMarkupHistory = async (
        headers: THeaders
    ): Promise<DefaultMarkupHistoryResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/default/history`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return RateMapper.dtoToDefaultMarkupHistoryDto(res);
    };
}
