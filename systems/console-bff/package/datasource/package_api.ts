/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import dayjs from "dayjs";

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
    return this.get(`/${VERSION}/${PACKAGES}/${packageId}`, {}).then(res =>
      dtoToPackageDto(res)
    );
  };

  getPackages = async (): Promise<PackagesResDto> => {
    return this.get(`/${VERSION}/${PACKAGES}`).then(res =>
      dtoToPackagesDto(res)
    );
  };

  addPackage = async (
    req: AddPackageInputDto,
    headers: THeaders
  ): Promise<PackageDto> => {
    this.logger.info(`Add pacakge request`);
    const baserate = await this.get(`/${VERSION}/baserates/history`);

    return this.post(`/${VERSION}/${PACKAGES}`, {
      body: {
        name: req.name,
        amount: req.amount,
        data_unit: req.dataUnit,
        data_volume: req.dataVolume,
        duration: req.duration,
        active: true,
        flat_rate: true,
        markup: 0,
        overdraft: 0,
        sms_volume: 0,
        voice_volume: 0,
        traffic_policy: 0,
        networks: [],
        type: "prepaid",
        apn: "ukama.tel",
        owner_id: headers.userId,
        sim_type: "operator_data",
        to: dayjs().add(5, "year").format(),
        from: dayjs().add(7, "day").format(),
        voice_unit: "seconds",
        baserate_id: baserate.rates[0].uuid,
      },
    }).then(res => dtoToPackageDto(res));
  };

  deletePackage = async (packageId: string): Promise<IdResponse> => {
    return this.delete(`/${VERSION}/${PACKAGES}/${packageId}`).then(() => {
      return {
        uuid: packageId,
      };
    });
  };

  updatePackage = async (
    packageId: string,
    req: UpdatePackageInputDto
  ): Promise<PackageDto> => {
    return this.patch(`/${VERSION}/${PACKAGES}/${packageId}`, {
      body: req,
    }).then(res => dtoToPackageDto(res));
  };
}

export default PackageApi;
