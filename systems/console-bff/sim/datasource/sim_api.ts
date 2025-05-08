/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { ENCRYPTION_KEY } from "../../common/configs";
import { logger } from "../../common/logger";
import { epochToISOString } from "../../common/utils";
import generateTokenFromIccid from "../../common/utils/generateSimToken";
import {
  AllocateSimAPIDto,
  AllocateSimInputDto,
  DeleteSimInputDto,
  DeleteSimResDto,
  GetPackagesForSimInputDto,
  GetSimByNetworkInputDto,
  GetSimBySubscriberIdInputDto,
  GetSimInputDto,
  GetSimPackagesDtoAPI,
  RemovePackageFormSimInputDto,
  RemovePackageFromSimResDto,
  SimDataUsage,
  SimDetailsDto,
  SimDto,
  SimPoolStatsDto,
  SimStatusResDto,
  SimUsageInputDto,
  SimsResDto,
  SubscriberToSimsDto,
  ToggleSimStatusInputDto,
  UploadSimsInputDto,
  UploadSimsResDto,
} from "../resolver/types";
import {
  dtoToAllocateSimResDto,
  dtoToSimDetailsDto,
  dtoToSimResDto,
  dtoToSimsDto,
  dtoToUsageDto,
  mapSubscriberToSimsResDto,
} from "./mapper";

const VERSION = "v1";
const SIMPOOL = "simpool";
const SIM = "sim";

class SimApi extends RESTDataSource {
  uploadSims = async (
    baseURL: string,
    req: UploadSimsInputDto
  ): Promise<UploadSimsResDto> => {
    this.logger.info(
      `UploadSims [PUT]: ${baseURL}/${VERSION}/${SIMPOOL}/upload`
    );
    this.baseURL = baseURL;
    return this.put(`/${VERSION}/${SIMPOOL}/upload`, {
      body: {
        data: req.data,
        sim_type: req.simType,
      },
    }).then(res => res);
  };

  toggleSimStatus = async (
    baseURL: string,
    req: ToggleSimStatusInputDto
  ): Promise<SimStatusResDto> => {
    this.logger.info(
      `ToggleSimStatus [PATCH]: ${baseURL}/${VERSION}/${SIM}/${req.sim_id}`
    );
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${SIM}/${req.sim_id}`, {
      body: { status: req.status },
    });
  };

  allocateSim = async (
    baseURL: string,
    req: AllocateSimInputDto
  ): Promise<AllocateSimAPIDto> => {
    this.logger.info(`AllocateSim [POST]: ${baseURL}/${VERSION}/${SIM}`);
    this.baseURL = baseURL;
    const getToken = (): string | null => {
      if (req.iccid) {
        const token = generateTokenFromIccid(req.iccid, ENCRYPTION_KEY ?? "");
        return token;
      }

      return null;
    };
    const requestBody = {
      ...req,
      ...(req.iccid ? { sim_token: getToken() } : {}),
    };

    const simRes = await this.post(`/${VERSION}/${SIM}`, {
      body: {
        ...requestBody,
      },
    });

    logger.info(`SimRes: ${JSON.stringify(simRes)}`);

    return dtoToAllocateSimResDto(simRes);
  };

  getSim = async (baseURL: string, req: GetSimInputDto): Promise<SimDto> => {
    this.logger.info(`GetSim [GET]: ${baseURL}/${VERSION}/${SIM}/${req.simId}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIM}/${req.simId}`, {
      params: {
        simId: req.simId,
      },
    }).then(res => dtoToSimResDto(res));
  };

  getSims = async (baseURL: string, type: string): Promise<SimsResDto> => {
    this.logger.info(
      `GetSims [GET]: ${baseURL}/${VERSION}/${SIMPOOL}/sims/${type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIMPOOL}/sims/${type}`).then(res =>
      dtoToSimsDto(res)
    );
  };

  getDataUsage = async (
    baseURL: string,
    data: SimUsageInputDto
  ): Promise<SimDataUsage> => {
    this.baseURL = baseURL;
    const params = new URLSearchParams({
      cdr_type: data.type,
      sim_type: data.type,
      ...(data.from && { from: epochToISOString(data.from) }),
      ...(data.to && { to: epochToISOString(data.to) }),
      ...(data.simId && { sim_id: data.simId }),
    }).toString();
    this.logger.info(
      `GetDataUsage [GET]: ${baseURL}/${VERSION}/usages?${params}`
    );
    return this.get(`/${VERSION}/usages?${params}`).then(res =>
      dtoToUsageDto(res, data)
    );
  };

  getSimByNetworkId = async (
    baseURL: string,
    req: GetSimByNetworkInputDto
  ): Promise<SimDetailsDto> => {
    this.baseURL = baseURL;
    return this.put(``, {
      body: {
        networkId: req.networkId,
      },
    }).then(res => dtoToSimDetailsDto(res));
  };

  deleteSim = async (
    baseURL: string,
    req: DeleteSimInputDto
  ): Promise<DeleteSimResDto> => {
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${SIM}/${req.simId}`).then(res => res);
  };
  addPackageToSim = async (
    baseURL: string,
    simId: string,
    packageId: string,
    startDate: string
  ): Promise<void> => {
    this.baseURL = baseURL;
    this.logger.info(
      `AddPackagesToSim [POST]: ${baseURL}/${VERSION}/${SIM}/${simId}/packages`
    );
    return await this.post(`/${VERSION}/${SIM}/package`, {
      body: {
        sim_id: simId,
        package_id: packageId,
        start_date: startDate,
      },
    });
  };

  removePackageFromSim = async (
    baseURL: string,
    req: RemovePackageFormSimInputDto
  ): Promise<RemovePackageFromSimResDto> => {
    this.baseURL = baseURL;
    return this.put(``, {
      body: {
        ...req,
      },
    }).then(res => res);
  };

  getPackagesForSim = async (
    baseURL: string,
    req: GetPackagesForSimInputDto
  ): Promise<GetSimPackagesDtoAPI> => {
    this.logger.info(
      `GetPackageForSim [GET]: ${baseURL}/${VERSION}/${SIM}/${req.sim_id}/packages`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIM}/packages/${req.sim_id}`).then(
      res => res
    );
  };

  getSimsBySubscriberId = async (
    baseURL: string,
    req: GetSimBySubscriberIdInputDto
  ): Promise<SubscriberToSimsDto> => {
    this.logger.info(
      `GetSimsBySubscriberId [GET]: ${baseURL}/${VERSION}/sim/subscriber/${req.subscriberId}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/sim/subscriber/${req.subscriberId}`).then(
      res => mapSubscriberToSimsResDto(res)
    );
  };

  getSimPoolStats = async (
    baseURL: string,
    type: string
  ): Promise<SimPoolStatsDto> => {
    this.logger.info(
      `GetSimPoolStats [GET]: ${baseURL}/${VERSION}/${SIMPOOL}/stats/${type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIMPOOL}/stats/${type}`).then(res => res);
  };
}

export default SimApi;
