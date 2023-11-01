/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { SUBSCRIBER_API_GW } from "../../common/configs";
import generateTokenFromIccid from "../../common/utils/generateSimToken";
import {
  AddPackageSimResDto,
  AddPackageToSimInputDto,
  AllocateSimInputDto,
  DeleteSimInputDto,
  DeleteSimResDto,
  GetPackagesForSimInputDto,
  GetPackagesForSimResDto,
  GetSimByNetworkInputDto,
  GetSimBySubscriberIdInputDto,
  GetSimInputDto,
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
  ToggleSimStatusInputDto,
  UploadSimsInputDto,
  UploadSimsResDto,
} from "../resolver/types";
import { dtoToSimDetailsDto, dtoToSimResDto, dtoToSimsDto } from "./mapper";

const version = "/v1/simpool";

class SimApi extends RESTDataSource {
  baseURL = SUBSCRIBER_API_GW + version;

  uploadSims = async (req: UploadSimsInputDto): Promise<UploadSimsResDto> => {
    return this.put(`/upload`, {
      body: {
        data: req.data,
        sim_type: req.simType,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  allocateSim = async (req: AllocateSimInputDto): Promise<SimDto> => {
    const token = generateTokenFromIccid(
      req.iccid,
      process.env.ENCRYPTION_KEY || ""
    );
    return this.put(``, {
      body: {
        ...req,
        sim_token: token,
      },
    })
      .then(res => dtoToSimResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  toggleSimStatus = async (
    req: ToggleSimStatusInputDto
  ): Promise<SimStatusResDto> => {
    return this.put(``, {
      body: {
        simId: req.simId,
        status: req.status,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSim = async (req: GetSimInputDto): Promise<SimDto> => {
    return this.get(``, {
      params: {
        simId: req.simId,
      },
    })
      .then(res => dtoToSimResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSims = async (type: string): Promise<SimsResDto> => {
    return this.get(`/sims/${type}`)
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

  getSimBySubscriberId = async (
    req: GetSimBySubscriberIdInputDto
  ): Promise<SimDetailsDto> => {
    return this.put(``, {
      body: {
        subscriberId: req.subscriberId,
      },
    })
      .then(res => dtoToSimDetailsDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimByNetworkId = async (
    req: GetSimByNetworkInputDto
  ): Promise<SimDetailsDto> => {
    return this.put(``, {
      body: {
        networkId: req.networkId,
      },
    })
      .then(res => dtoToSimDetailsDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  deleteSim = async (req: DeleteSimInputDto): Promise<DeleteSimResDto> => {
    return this.put(``, {
      body: {
        simId: req.simId,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  addPackegeToSim = async (
    req: AddPackageToSimInputDto
  ): Promise<AddPackageSimResDto> => {
    return this.put(``, {
      body: {
        ...req,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  removePackageFromSim = async (
    req: RemovePackageFormSimInputDto
  ): Promise<RemovePackageFromSimResDto> => {
    return this.put(``, {
      body: {
        ...req,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getPackagesForSim = async (
    req: GetPackagesForSimInputDto
  ): Promise<GetPackagesForSimResDto> => {
    return this.put(``, {
      body: {
        simId: req.simId,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSimPoolStats = async (type: string): Promise<SimPoolStatsDto> => {
    return this.get(`/stats/${type}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  setActivePackageForSim = async (
    req: SetActivePackageForSimInputDto
  ): Promise<SetActivePackageForSimResDto> => {
    return this.put(``, {
      body: {
        ...req,
      },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default SimApi;
