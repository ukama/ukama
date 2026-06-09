/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import {
  ForceUnlockInputDto,
  MarkOperationRunningInputDto,
  OperationDto,
  ResourceLockDto,
  StartOperationInputDto,
  StartOperationResponseDto,
} from "../resolvers/types";
import {
  mapGetOperation,
  mapResourceLock,
  mapStartOperation,
} from "./mapper";

const OPERATIONS = "operations";

class OperationAPI extends BaseRESTDataSource {
  startOperation = async (
    baseURL: string,
    data: StartOperationInputDto
  ): Promise<StartOperationResponseDto> => {
    this.logger.info(`StartOperation [POST]: ${baseURL}/${VERSION}/${OPERATIONS}`);
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${OPERATIONS}`, {
      body: {
        type: data.type,
        system: data.system,
        resource_key: data.resourceKey,
        requested_by: data.requestedBy,
        idempotency_key: data.idempotencyKey,
        lease_seconds: data.leaseSeconds,
      },
    })
      .then(res => mapStartOperation(res))
      .catch(error => {
        this.logger.error(`Error starting operation: ${error}`);
        throw error;
      });
  };

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

  markOperationRunning = async (
    baseURL: string,
    data: MarkOperationRunningInputDto
  ): Promise<OperationDto | undefined> => {
    this.logger.info(
      `MarkOperationRunning [POST]: ${baseURL}/${VERSION}/${OPERATIONS}/${data.id}/run`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${OPERATIONS}/${data.id}/run`, {
      body: { fencing_token: data.fencingToken },
    })
      .then(res => mapGetOperation(res))
      .catch(error => {
        this.logger.error(`Error marking operation running: ${error}`);
        throw error;
      });
  };

  forceUnlock = async (
    baseURL: string,
    data: ForceUnlockInputDto,
    userId: string
  ): Promise<OperationDto | undefined> => {
    this.logger.info(
      `ForceUnlock [POST]: ${baseURL}/${VERSION}/${OPERATIONS}/${data.id}/force-unlock`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${OPERATIONS}/${data.id}/force-unlock`, {
      body: { user_id: userId, reason: data.reason },
    })
      .then(res => mapGetOperation(res))
      .catch(error => {
        this.logger.error(`Error force-unlocking operation: ${error}`);
        throw error;
      });
  };
}

export default OperationAPI;
