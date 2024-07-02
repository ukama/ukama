/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import dayjs from "dayjs";
import { GraphQLError } from "graphql";

import { ENCRYPTION_KEY } from "../../common/configs";
import generateTokenFromIccid from "../../common/utils/generateSimToken";
import {
  AddPackageSimResDto,
  AddPackageToSimInputDto,
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
  SetActivePackageForSimInputDto,
  SetActivePackageForSimResDto,
  SimDataUsage,
  SimDetailsDto,
  SimDto,
  SimPoolStatsDto,
  SimStatusResDto,
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
} from "./mapper";

const VERSION = "v1";
const SIMPOOL = "simpool";
const SIM = "sim";

class SimApi extends RESTDataSource {
  uploadSims = async (
    baseURL: string,
    req: UploadSimsInputDto
  ): Promise<UploadSimsResDto> => {
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
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${SIM}/${req.sim_id}`, {
      body: { status: req.status },
    })
      .then(res => {
        return res;
      })
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
  allocateSim = async (
    baseURL: string,
    req: AllocateSimInputDto
  ): Promise<AllocateSimAPIDto> => {
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

    return this.post(`/${VERSION}/${SIM}`, {
      body: {
        ...requestBody,
      },
    }).then(res => {
      this.addPackageToSim(baseURL, {
        package_id: req.package_id,
        sim_id: res.sim.id,
        start_date: dayjs().format(),
      })
        .then(async (response: any) => {
          await this.toggleSimStatus(baseURL, {
            sim_id: res.sim.id,
            status: "active",
          });
          await this.setActivePackageForSim(baseURL, {
            sim_id: res.sim.id,
            package_id: response.packages[0].id,
          });
        })
        .catch((error: any) => {
          throw new GraphQLError(error);
        });

      return dtoToAllocateSimResDto(res);
    });
  };

  getSim = async (baseURL: string, req: GetSimInputDto): Promise<SimDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIM}/${req.simId}`, {
      params: {
        simId: req.simId,
      },
    }).then(res => dtoToSimResDto(res));
  };

  getSims = async (baseURL: string, type: string): Promise<SimsResDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIMPOOL}/sims/${type}`)
      .then(res => dtoToSimsDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getDataUsage = async (
    baseURL: string,
    simId: string
  ): Promise<SimDataUsage> => {
    this.baseURL = baseURL;
    //TODO: GET SIM DATA USAGE METRIC HERE
    return {
      usage: `1240-${simId}`,
    };
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
    return this.put(``, {
      body: {
        simId: req.simId,
      },
    }).then(res => res);
  };

  addPackageToSim = async (
    baseURL: string,
    req: AddPackageToSimInputDto
  ): Promise<AddPackageSimResDto> => {
    this.baseURL = baseURL;
    return this.put(``, {
      body: {
        ...req,
      },
    }).then(res => res);
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
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIM}/packages/${req.sim_id}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimsBySubscriberId = async (
    baseURL: string,
    req: GetSimBySubscriberIdInputDto
  ): Promise<SubscriberToSimsDto> => {
    this.baseURL = baseURL;
    return this.get(`/sim/subscriber/${req.subscriberId}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimPoolStats = async (
    baseURL: string,
    type: string
  ): Promise<SimPoolStatsDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SIMPOOL}/stats/${type}`).then(res => res);
  };

  setActivePackageForSim = async (
    baseURL: string,
    req: SetActivePackageForSimInputDto
  ): Promise<SetActivePackageForSimResDto> => {
    this.baseURL = baseURL;
    return this.patch(
      `/${VERSION}/${SIM}/${req.sim_id}/package/${req.package_id}`,
      {
        body: {
          ...req,
        },
      }
    )
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default SimApi;
