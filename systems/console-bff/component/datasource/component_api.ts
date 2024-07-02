/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { THeaders } from "../../common/types";
import { ComponentDto, ComponentsResDto } from "../resolvers/types";
import { dtoTocomponentDto, dtoTocomponentsDto } from "./mapper";

const COMPONENTS = "components";

class ComponentApi extends RESTDataSource {
  getComponents = async (
    headers: THeaders,
    baseURL: string,
    category: string
  ): Promise<ComponentsResDto> => {
    this.baseURL = baseURL;
    return this.get(
      `/${VERSION}/${COMPONENTS}/user/${headers.userId}?category=${category}`
    ).then(res => dtoTocomponentsDto(res));
  };

  getComponentById = async (
    baseURL: string,
    componentId: string
  ): Promise<ComponentDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${COMPONENTS}/${componentId}`).then(res =>
      dtoTocomponentDto(res)
    );
  };
}

export default ComponentApi;
