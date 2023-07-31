import { RESTDataSource } from "@apollo/datasource-rest";

import { BoolResponse } from "../../common/types";
import {
  DefaultMarkupHistoryResDto,
  DefaultMarkupInputDto,
  DefaultMarkupResDto,
} from "../resolver/types";
import { SERVER } from "./../../constants/endpoints";
import { dtoToDefaultMarkupDto, dtoToDefaultMarkupHistoryDto } from "./mapper";

class RateApi extends RESTDataSource {
  defaultMarkup = async (req: DefaultMarkupInputDto): Promise<BoolResponse> => {
    return this.post(
      `${SERVER.DATA_PLAN_MARKUP_API_URL}/${req.markup}/default`
    ).then(res => {
      return {
        success: true,
      };
    });
  };

  getDefaultMarkup = async (): Promise<DefaultMarkupResDto> => {
    return this.get(`${SERVER.DATA_PLAN_MARKUP_API_URL}/default`).then(res =>
      dtoToDefaultMarkupDto(res)
    );
  };

  getDefaultMarkupHistory = async (): Promise<DefaultMarkupHistoryResDto> => {
    return this.get(`${SERVER.DATA_PLAN_MARKUP_API_URL}/default/history`).then(
      res => dtoToDefaultMarkupHistoryDto(res)
    );
  };
}

export default RateApi;
