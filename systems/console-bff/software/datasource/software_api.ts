/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import {
  GetReleaseCatalogInput,
  GetSoftwaresInput,
  PromoteReleaseInputDto,
  PromoteReleaseResponse,
  ReleaseCatalog,
  Softwares,
  StringResponse,
  UpdateSoftwareInputDto,
} from "../resolvers/types";
import {
  mapPromoteRelease,
  mapReleaseCatalog,
  mapSoftwares,
  mapUpdateSoftware,
} from "./mapper";

const SOFTWARE = "software";

class SoftwareAPI extends BaseRESTDataSource {
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

  promoteRelease = async (
    baseURL: string,
    data: PromoteReleaseInputDto
  ): Promise<PromoteReleaseResponse> => {
    const { name, version, type } = data;
    const query = type ? `?type=${encodeURIComponent(type)}` : "";
    this.logger.info(
      `PromoteRelease [POST]: ${baseURL}/${VERSION}/${SOFTWARE}/promote/${name}/${version}${query}`
    );
    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/${SOFTWARE}/promote/${name}/${version}${query}`
    )
      .then(res => {
        return mapPromoteRelease(res);
      })
      .catch(error => {
        this.logger.error(`Error promoting release: ${error}`);
        throw error;
      });
  };

  getReleaseCatalog = async (
    baseURL: string,
    data: GetReleaseCatalogInput
  ): Promise<ReleaseCatalog> => {
    const queryParams = new URLSearchParams();
    if (data.name) {
      queryParams.append("name", data.name);
    }
    if (data.type) {
      queryParams.append("type", data.type);
    }
    const qs = queryParams.toString();
    this.logger.info(
      `GetReleaseCatalog [GET]: ${baseURL}/${VERSION}/${SOFTWARE}/releases${
        qs ? `?${qs}` : ""
      }`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SOFTWARE}/releases${qs ? `?${qs}` : ""}`)
      .then(res => {
        return mapReleaseCatalog(res);
      })
      .catch(error => {
        this.logger.error(`Error getting release catalog: ${error}`);
        throw error;
      });
  };
}

export default SoftwareAPI;
