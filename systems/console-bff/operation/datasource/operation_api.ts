/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import { OperationDto, ResourceLockDto } from "../resolvers/types";
import { mapGetOperation, mapResourceLock } from "./mapper";

const OPERATIONS = "operations";

class OperationAPI extends BaseRESTDataSource {
  getOperation = async (
    baseURL: string,
    id: string
  ): Promise<OperationDto | undefined> => {
    this.logger.info(
      `GetOperation [GET]: ${baseURL}/${VERSION}/${OPERATIONS}/${id}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${OPERATIONS}/${id}`)
      .then(res => mapGetOperation(res))
      .catch(error => {
        this.logger.error(`Error getting operation: ${error}`);
        throw error;
      });
  };

  getResourceLock = async (
    baseURL: string,
    resourceKey: string
  ): Promise<ResourceLockDto> => {
    const queryParams = new URLSearchParams();
    queryParams.append("resource_key", resourceKey);
    this.logger.info(
      `GetResourceLock [GET]: ${baseURL}/${VERSION}/${OPERATIONS}?${queryParams.toString()}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${OPERATIONS}?${queryParams.toString()}`)
      .then(res => mapResourceLock(res))
      .catch(error => {
        this.logger.error(`Error getting resource lock: ${error}`);
        throw error;
      });
  };
}

export default OperationAPI;
