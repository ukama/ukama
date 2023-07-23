import { RESTDataSource } from "@apollo/datasource-rest";
import { BoolResponse, THeaders  } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
import RateMapper from "./mapper";
import {
    DefaultMarkupHistoryResDto,
    DefaultMarkupInputDto,
    DefaultMarkupResDto,
} from "../types";


export class RateApi extends RESTDataSource {
    defaultMarkup = async (
        req: DefaultMarkupInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        return this.post(`${SERVER.DATA_PLAN_MARKUP_API_URL}/${req.markup}/default`, {
            headers: getHeaders(headers),
          }).then(res =>  {
            return {
            success: true,
        }
    });
    };

    getDefaultMarkup = async (
        headers: THeaders
    ): Promise<DefaultMarkupResDto> => {
        return this.get(`${SERVER.DATA_PLAN_MARKUP_API_URL}/default`,{headers: getHeaders(headers)}).then(res => 
            RateMapper.dtoToDefaultMarkupDto(res));
    };

    getDefaultMarkupHistory = async (
        headers: THeaders
    ): Promise<DefaultMarkupHistoryResDto> => {
        return this.get(`${SERVER.DATA_PLAN_MARKUP_API_URL}/default/history`,{headers: getHeaders(headers)}).then(res => 
            RateMapper.dtoToDefaultMarkupHistoryDto(res));
    };
}

