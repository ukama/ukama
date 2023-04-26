import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { HTTP404Error, Messages, checkError } from "../../errors";
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
        cookie: ParsedCookie,
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/${req.markup}/default`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
    getDefaultMarkup = async (
        cookie: ParsedCookie,
    ): Promise<DefaultMarkupResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/default`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return RateMapper.dtoToDefaultMarkupDto(res);
    };
    getDefaultMarkupHistory = async (
        cookie: ParsedCookie,
    ): Promise<DefaultMarkupHistoryResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_MARKUP_API_URL}/default/history`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return RateMapper.dtoToDefaultMarkupHistoryDto(res);
    };
}
