/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { INVENTORY_API_GW, VERSION } from "../../common/configs";
import { ComponentDto, ComponentsResDto } from "../resolvers/types";
import { dtoTocomponentDto, dtoTocomponentsDto } from "./mapper";

class ComponentApi extends RESTDataSource {
  baseURL = INVENTORY_API_GW;

  getComponents = async (
    userId: string,
    category: string
  ): Promise<ComponentsResDto> => {
    return this.get(
      `/${VERSION}/components/user/${userId}?category=${category}`
    ).then(res => dtoTocomponentsDto(res));
  };

  getComponent = async (componentId: string): Promise<ComponentDto> => {
    return this.get(`/${VERSION}/components/${componentId}`).then(res =>
      dtoTocomponentDto(res)
    );
  };
}

export default ComponentApi;
