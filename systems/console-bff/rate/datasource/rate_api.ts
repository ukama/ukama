/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { CBooleanResponse } from "../../common/types";
import {
  DefaultMarkupHistoryResDto,
  DefaultMarkupInputDto,
  DefaultMarkupResDto,
} from "../resolver/types";
import { dtoToDefaultMarkupDto, dtoToDefaultMarkupHistoryDto } from "./mapper";

const VERSION = "v1";
const MARKUP = "markup";

class RateApi extends RESTDataSource {
  defaultMarkup = async (
    baseURL: string,
    req: DefaultMarkupInputDto
  ): Promise<CBooleanResponse> => {
    this.logger.info(
      `DefaultMarkup [POST]: ${baseURL}/${VERSION}/${MARKUP}/${req.markup}/default`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${MARKUP}/${req.markup}/default`).then(() => {
      return {
        success: true,
      };
    });
  };

  getDefaultMarkup = async (baseURL: string): Promise<DefaultMarkupResDto> => {
    this.logger.info(
      `GetDefaultMarkup [GET]: ${baseURL}/${VERSION}/${MARKUP}/default`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${MARKUP}/default`).then(res =>
      dtoToDefaultMarkupDto(res)
    );
  };

  getDefaultMarkupHistory = async (
    baseURL: string
  ): Promise<DefaultMarkupHistoryResDto> => {
    this.logger.info(
      `GetDefaultMarkupHistory [GET]: ${baseURL}/${VERSION}/${MARKUP}/default/history`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${MARKUP}/default/history`).then(res =>
      dtoToDefaultMarkupHistoryDto(res)
    );
  };
}

export default RateApi;
