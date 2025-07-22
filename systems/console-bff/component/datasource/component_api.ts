/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { INVENTORY_API_GW, VERSION } from "../../common/configs";
import { THeaders } from "../../common/types";
import {
  ComponentAPIDto,
  ComponentDto,
  ComponentsAPIResDto,
  ComponentsResDto,
} from "../resolvers/types";
import { dtoToComponentDto, dtoToComponentsDto } from "./mapper";

class ComponentApi extends RESTDataSource {
  baseURL = INVENTORY_API_GW;

  getComponentsByUserId = async (
    headers: THeaders,
    category: string
  ): Promise<ComponentsResDto> => {
    this.logger.info(
      `GetComponentByUserId [GET]: ${this.baseURL}/${VERSION}/user/${
        headers.userId
      }?category=${category.toUpperCase()}`
    );
    const response = await this.get<ComponentsAPIResDto>(
      `/${VERSION}/components/user/${
        headers.userId
      }?category=${category.toUpperCase()}`
    );
    return dtoToComponentsDto(response);
  };

  getComponentById = async (componentId: string): Promise<ComponentDto> => {
    this.logger.info(
      `GetComponentById [GET]: ${this.baseURL}/${VERSION}/components/${componentId}`
    );
    const response = await this.get<ComponentAPIDto>(
      `/${VERSION}/components/${componentId}`
    );
    return dtoToComponentDto(response);
  };
}

export default ComponentApi;
