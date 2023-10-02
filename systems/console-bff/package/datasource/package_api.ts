import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { DATA_API_GW } from "../../common/configs";
import { IdResponse, THeaders } from "../../common/types";
import {
  AddPackageInputDto,
  PackageDto,
  PackagesResDto,
  UpdatePackageInputDto,
} from "../resolver/types";
import { dtoToPackageDto, dtoToPackagesDto } from "./mapper";

const VERSION = "v1";
const PACKAGES = "packages";

class PackageApi extends RESTDataSource {
  baseURL = DATA_API_GW;
  getPackage = async (packageId: string): Promise<PackageDto> => {
    return this.get(`/${VERSION}/${PACKAGES}/${packageId}`, {})
      .then(res => dtoToPackageDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getPackages = async (headers: THeaders): Promise<PackagesResDto> => {
    return this.get(`/${VERSION}/${PACKAGES}/orgs/${headers.orgId}`)
      .then(res => dtoToPackagesDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  addPackage = async (
    req: AddPackageInputDto,
    headers: THeaders
  ): Promise<PackageDto> => {
    return this.post(`/${VERSION}/${PACKAGES}`, {
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
    })
      .then(res => dtoToPackageDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  deletePackage = async (packageId: string): Promise<IdResponse> => {
    return this.delete(`/${VERSION}/${PACKAGES}/${packageId}`)
      .then(() => {
        return {
          uuid: packageId,
        };
      })
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  updatePackage = async (
    packageId: string,
    req: UpdatePackageInputDto
  ): Promise<PackageDto> => {
    return this.patch(`/${VERSION}/${PACKAGES}/${packageId}`, {
      body: req,
    })
      .then(res => dtoToPackageDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default PackageApi;
