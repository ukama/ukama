import { RESTDataSource } from "@apollo/datasource-rest";

import { SERVER } from "../../constants/endpoints";
import generateTokenFromIccid from "../../utils/generateSimToken";
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

class SimApi extends RESTDataSource {
  uploadSims = async (req: UploadSimsInputDto): Promise<UploadSimsResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/upload`, {
      body: {
        data: req.data,
        sim_type: req.simType,
      },
    }).then(res => res);
  };

  allocateSim = async (req: AllocateSimInputDto): Promise<SimDto> => {
    const token = generateTokenFromIccid(
      req.iccid,
      process.env.ENCRYPTION_KEY || ""
    );
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        ...req,
        sim_token: token,
      },
    }).then(res => dtoToSimResDto(res));
  };

  toggleSimStatus = async (
    req: ToggleSimStatusInputDto
  ): Promise<SimStatusResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        simId: req.simId,
        status: req.status,
      },
    }).then(res => res);
  };

  getSim = async (req: GetSimInputDto): Promise<SimDto> => {
    //TODO: body is in get request

    // const res = await catchAsyncIOMethod({
    //     type: API_METHOD_TYPE.GET,
    //     path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
    //     body: {
    //         simId: req.simId,
    //     },
    //     headers: getHeaders(headers),
    // });
    // if (checkError(res)) throw new Error(res.message);
    // return dtoToSimResDto(res);

    return this.get(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      params: {
        simId: req.simId,
      },
    }).then(res => dtoToSimResDto(res));
  };

  getSims = async (type: string): Promise<SimsResDto> => {
    return this.get(
      `${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/sims/${type}`,
      {}
    ).then(res => dtoToSimsDto(res));
  };

  getDataUsage = async (simId: string): Promise<SimDataUsage> => {
    // const res = await catchAsyncIOMethod({
    //     type: API_METHOD_TYPE.GET,
    //     path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${simId}`,
    //     getHeaders(headers)
    // });
    // if (checkError(res)) throw new Error(res.message);
    return {
      usage: "1240",
    };
  };

  getSimBySubscriberId = async (
    req: GetSimBySubscriberIdInputDto
  ): Promise<SimDetailsDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        subscriberId: req.subscriberId,
      },
    }).then(res => dtoToSimDetailsDto(res));
  };

  getSimByNetworkId = async (
    req: GetSimByNetworkInputDto
  ): Promise<SimDetailsDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        networkId: req.networkId,
      },
    }).then(res => dtoToSimDetailsDto(res));
  };

  deleteSim = async (req: DeleteSimInputDto): Promise<DeleteSimResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        simId: req.simId,
      },
    }).then(res => res);
  };

  addPackegeToSim = async (
    req: AddPackageToSimInputDto
  ): Promise<AddPackageSimResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        ...req,
      },
    }).then(res => res);
  };

  removePackageFromSim = async (
    req: RemovePackageFormSimInputDto
  ): Promise<RemovePackageFromSimResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        ...req,
      },
    }).then(res => res);
  };

  getPackagesForSim = async (
    req: GetPackagesForSimInputDto
  ): Promise<GetPackagesForSimResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        simId: req.simId,
      },
    }).then(res => res);
  };

  getSimPoolStats = async (type: string): Promise<SimPoolStatsDto> => {
    return this.get(
      `${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/stats/${type}`,
      {}
    ).then(res => res);
  };

  setActivePackageForSim = async (
    req: SetActivePackageForSimInputDto
  ): Promise<SetActivePackageForSimResDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: {
        ...req,
      },
    }).then(res => res);
  };
}

export default SimApi;
