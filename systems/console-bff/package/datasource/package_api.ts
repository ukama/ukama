import { RESTDataSource } from "@apollo/datasource-rest";

import { DATA_API_GW } from "../../common/configs";
import { IdResponse, THeaders } from "../../common/types";
import {
  AddPackageInputDto,
  PackageDto,
  PackagesResDto,
  UpdatePackageInputDto,
} from "../resolver/types";
import { dtoToPackageDto, dtoToPackagesDto } from "./mapper";

const version = "/v1/packages";

class PackageApi extends RESTDataSource {
  baseURL = DATA_API_GW + version;
  getPackage = async (packageId: string): Promise<PackageDto> => {
    return this.get(`/${packageId}`, {}).then(res => dtoToPackageDto(res));
  };

  getPackages = async (headers: THeaders): Promise<PackagesResDto> => {
    return this.get(`/org/${headers.orgId}`).then(res => dtoToPackagesDto(res));
  };

  addPackage = async (
    req: AddPackageInputDto,
    headers: THeaders
  ): Promise<PackageDto> => {
    return this.post("", {
      body: {
        duration: req.duration,
        active: true,
        amount: req.amount,
        data_unit: req.dataUnit,
        data_volume: req.dataVolume,
        flat_rate: true,
        from: "2023-04-01T00:00:00Z",
        markup: 0,
        name: req.name,
        org_id: headers.orgId,
        owner_id: headers.userId,
        sim_type: "ukama_data",
        sms_volume: 0,
        to: "",
        type: "prepaid",
        voice_unit: "seconds",
        voice_volume: 0,
      },
    }).then(res => dtoToPackageDto(res));
  };

  deletePackage = async (packageId: string): Promise<IdResponse> => {
    return this.delete(`/${packageId}`).then(() => {
      return {
        uuid: packageId,
      };
    });
  };

  updatePackage = async (
    packageId: string,
    req: UpdatePackageInputDto
  ): Promise<PackageDto> => {
    return this.patch(`/${packageId}`, {
      body: req,
    }).then(res => dtoToPackageDto(res));
  };
}

export default PackageApi;
