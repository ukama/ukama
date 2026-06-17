/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import dayjs from "dayjs";

import { SIM_TYPE } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import { IdResponse, THeaders } from "../../common/types";
import {
  AddPackageInputDto,
  PackageDto,
  PackageNameAvailabilityResDto,
  PackagesResDto,
  UpdatePackageInputDto,
} from "../resolver/types";
import { toBackendDataUnit } from "./dataUnit";
import {
  dtoToPackageDto,
  dtoToPackageNameAvailabilityDto,
  dtoToPackagesDto,
} from "./mapper";

const VERSION = "v1";
const PACKAGES = "packages";

class PackageApi extends BaseRESTDataSource {
  getPackage = async (
    baseURL: string,
    packageId: string
  ): Promise<PackageDto> => {
    this.logger.info(
      `GetPackage [GET]: ${baseURL}/${VERSION}/${PACKAGES}/${packageId}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${PACKAGES}/${packageId}`, {}).then(res =>
      dtoToPackageDto(res)
    );
  };

  isPackageNameAvailable = async (
    baseURL: string,
    name: string
  ): Promise<PackageNameAvailabilityResDto> => {
    const params = new URLSearchParams({ name });
    this.logger.info(
      `IsPackageNameAvailable [GET]: ${baseURL}/${VERSION}/${PACKAGES}/name-availability?${params.toString()}`
    );
    this.baseURL = baseURL;
    return this.get(
      `/${VERSION}/${PACKAGES}/name-availability?${params.toString()}`
    ).then(res => dtoToPackageNameAvailabilityDto(res));
  };

  getPackages = async (
    baseURL: string,
    networkId?: string
  ): Promise<PackagesResDto> => {
    this.logger.info(`GetPackages [GET]: ${baseURL}/${VERSION}/${PACKAGES}`);
    this.baseURL = baseURL;
    const res = await this.get(`/${VERSION}/${PACKAGES}`).then(r =>
      dtoToPackagesDto(r)
    );
    // When scoped to a network, keep that network's plans plus org-wide plans
    // (no networkId). No filter → all plans (unchanged behaviour).
    if (networkId) {
      res.packages = res.packages.filter(
        p => !p.networkId || p.networkId === networkId
      );
    }
    return res;
  };

  addPackage = async (
    baseURL: string,
    req: AddPackageInputDto,
    headers: THeaders
  ): Promise<PackageDto> => {
    this.logger.info(`AddPackage [POST]: ${baseURL}/${VERSION}/${PACKAGES}`);
    this.baseURL = baseURL;
    const baserate = await this.get(`/${VERSION}/baserates/history`);
    if (!baserate.rates || baserate?.rates?.length === 0) {
      throw new Error("No baserate found");
    }

    // TODO: Need to revisit this to, from values
    return this.post(`/${VERSION}/${PACKAGES}`, {
      body: {
        name: req.name,
        amount: req.amount,
        // Console sends short labels (GB/MB); the backend expects its enum.
        data_unit: toBackendDataUnit(req.dataUnit),
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
        network_id: req.networkId ?? "",
        type: "prepaid",
        apn: "ukama.tel",
        country: req.country,
        currency: req.currency,
        owner_id: headers.userId,
        sim_type: SIM_TYPE,
        to: dayjs().add(5, "year").format(),
        from: dayjs().add(7, "day").format(),
        voice_unit: "seconds",
        baserate_id: baserate.rates[0].uuid,
      },
    }).then(res => dtoToPackageDto(res));
  };

  deletePackage = async (
    baseURL: string,
    packageId: string
  ): Promise<IdResponse> => {
    this.logger.info(
      `DeletePackage [DELETE]: ${baseURL}/${VERSION}/${PACKAGES}/${packageId}`
    );
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${PACKAGES}/${packageId}`).then(() => {
      return {
        uuid: packageId,
      };
    });
  };

  updatePackage = async (
    baseURL: string,
    packageId: string,
    req: UpdatePackageInputDto
  ): Promise<PackageDto> => {
    this.logger.info(
      `UpdatePackage [PATCH]: ${baseURL}/${VERSION}/${PACKAGES}/${packageId}`
    );
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${PACKAGES}/${packageId}`, {
      body: {
        name: req.name,
        active: req.active,
      },
    }).then(res => dtoToPackageDto(res));
  };
}

export default PackageApi;
