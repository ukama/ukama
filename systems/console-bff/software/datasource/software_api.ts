/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  Apps,
  GetSoftwaresInput,
  Softwares,
  StringResponse,
  UpdateSoftwareInputDto,
} from "../resolvers/types";
import { mapApps, mapSoftwares, mapUpdateSoftware } from "./mapper";

const SOFTWARE = "software";

class SoftwareAPI extends RESTDataSource {
  getApps = async (baseURL: string): Promise<Apps> => {
    this.logger.info(`GetApps [GET]: ${baseURL}/${VERSION}/${SOFTWARE}/apps`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SOFTWARE}/apps`)
      .then(apps => {
        return mapApps(apps);
      })
      .catch(error => {
        this.logger.error(`Error getting apps: ${error}`);
        throw error;
      });
  };

  getSoftwares = async (
    baseURL: string,
    data: GetSoftwaresInput
  ): Promise<Softwares> => {
    const { name, nodeId, status } = data;
    const queryParams = new URLSearchParams();
    if (name) {
      queryParams.append("name", name);
    }
    if (nodeId) {
      queryParams.append("node_id", nodeId);
    }
    if (status) {
      queryParams.append("status", status);
    }
    this.logger.info(
      `GetSoftwares [GET]: ${baseURL}/${VERSION}/${SOFTWARE}?${queryParams.toString()}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SOFTWARE}?${queryParams.toString()}`)
      .then(softwares => {
        return mapSoftwares(softwares);
      })
      .catch(error => {
        this.logger.error(`Error getting softwares: ${error}`);
        throw error;
      });
  };

  updateSoftware = async (
    baseURL: string,
    data: UpdateSoftwareInputDto
  ): Promise<StringResponse> => {
    const { name, nodeId, tag } = data;
    this.logger.info(
      `UpdateSoftware [POST]: ${baseURL}/${VERSION}/${SOFTWARE}/update/${data.name}/${data.tag}/${data.nodeId}`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${SOFTWARE}/update/${name}/${tag}/${nodeId}`)
      .then(apps => {
        return mapUpdateSoftware(apps);
      })
      .catch(error => {
        this.logger.error(`Error getting apps: ${error}`);
        throw error;
      });
  };
}

export default SoftwareAPI;
