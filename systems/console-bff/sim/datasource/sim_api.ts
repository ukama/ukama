/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { ENCRYPTION_KEY, SUBSCRIBER_API_GW } from "../../common/configs";
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
  baseURL = SUBSCRIBER_API_GW;

  uploadSims = async (req: UploadSimsInputDto): Promise<UploadSimsResDto> => {
    return this.put(`/${VERSION}/${SIMPOOL}/upload`, {
      body: {
        data: req.data,
        sim_type: req.simType,
      },
    }).then(res => res);
  };

  toggleSimStatus = async (
    req: ToggleSimStatusInputDto
  ): Promise<SimStatusResDto> => {
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
    req: AllocateSimInputDto
  ): Promise<AllocateSimAPIDto> => {
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
    })
      .then(res => {
        this.toggleSimStatus({ sim_id: res.sim.id, status: "active" });
        this.getPackagesForSim({ sim_id: res.sim.id })
          .then((response: any) => {
            this.setActivePackageForSim({
              sim_id: res.sim.id,
              package_id: response.packages[0].id,
            });
          })
          .catch((error: any) => {
            console.log("SIM ALLOCATION 1 ERROR: ", error);
            throw new GraphQLError(error);
          });

        return dtoToAllocateSimResDto(res);
      })
      .catch(err => {
        console.log("SIM ALLOCATION 2 ERROR: ", err);
        throw new GraphQLError(err);
      });
  };

  getSim = async (req: GetSimInputDto): Promise<SimDto> => {
    return this.get(`/${VERSION}/${SIM}/${req.simId}`, {
      params: {
        simId: req.simId,
      },
    }).then(res => dtoToSimResDto(res));
  };

  getSims = async (type: string): Promise<SimsResDto> => {
    return this.get(`/${VERSION}/${SIMPOOL}/sims/${type}`)
      .then(res => dtoToSimsDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getDataUsage = async (simId: string): Promise<SimDataUsage> => {
    //TODO: GET SIM DATA USAGE METRIC HERE
    return {
      usage: `1240-${simId}`,
    };
  };

  getSimByNetworkId = async (
    req: GetSimByNetworkInputDto
  ): Promise<SimDetailsDto> => {
    return this.put(``, {
      body: {
        networkId: req.networkId,
      },
    }).then(res => dtoToSimDetailsDto(res));
  };

  deleteSim = async (req: DeleteSimInputDto): Promise<DeleteSimResDto> => {
    return this.put(``, {
      body: {
        simId: req.simId,
      },
    }).then(res => res);
  };

  addPackageToSim = async (
    req: AddPackageToSimInputDto
  ): Promise<AddPackageSimResDto> => {
    return this.put(``, {
      body: {
        ...req,
      },
    }).then(res => res);
  };

  removePackageFromSim = async (
    req: RemovePackageFormSimInputDto
  ): Promise<RemovePackageFromSimResDto> => {
    return this.put(``, {
      body: {
        ...req,
      },
    }).then(res => res);
  };

  getPackagesForSim = async (
    req: GetPackagesForSimInputDto
  ): Promise<GetSimPackagesDtoAPI> => {
    return this.get(`/${VERSION}/${SIM}/packages/${req.sim_id}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimsBySubscriberId = async (
    req: GetSimBySubscriberIdInputDto
  ): Promise<SubscriberToSimsDto> => {
    return this.get(`/sim/subscriber/${req.subscriberId}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimPoolStats = async (type: string): Promise<SimPoolStatsDto> => {
    return this.get(`/${VERSION}/${SIMPOOL}/stats/${type}`).then(res => res);
  };

  setActivePackageForSim = async (
    req: SetActivePackageForSimInputDto
  ): Promise<SetActivePackageForSimResDto> => {
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
