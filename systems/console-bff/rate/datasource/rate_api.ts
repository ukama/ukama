/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { DATA_API_GW } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  DefaultMarkupHistoryResDto,
  DefaultMarkupInputDto,
  DefaultMarkupResDto,
} from "../resolver/types";
import { dtoToDefaultMarkupDto, dtoToDefaultMarkupHistoryDto } from "./mapper";

const version = "/v1/markup";

class RateApi extends RESTDataSource {
  baseURL = DATA_API_GW + version;
  defaultMarkup = async (
    req: DefaultMarkupInputDto
  ): Promise<CBooleanResponse> => {
    return this.post(`/${req.markup}/default`).then(() => {
      return {
        success: true,
      };
    });
  };

  getDefaultMarkup = async (): Promise<DefaultMarkupResDto> => {
    return this.get(`/default`).then(res => dtoToDefaultMarkupDto(res));
  };

  getDefaultMarkupHistory = async (): Promise<DefaultMarkupHistoryResDto> => {
    return this.get(`/default/history`).then(res =>
      dtoToDefaultMarkupHistoryDto(res)
    );
  };
}

export default RateApi;
