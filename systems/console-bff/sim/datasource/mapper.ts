/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SIM_TYPES } from "../../common/enums";
import {
  AllocateSimAPIDto,
  GetSimPackagesDtoAPI,
  SimAPIResDto,
  SimAllResDto,
  SimDetailsDto,
  SimDto,
  SimPackagesResDto,
  SimToPackagesResDto,
  SimsAPIResDto,
  SimsAlloAPIResDto,
  SimsResDto,
  SubscriberToSimsDto,
  SubscriberToSimsResDto,
} from "../resolver/types";

export const dtoToAllocateSimResDto = (
  res: SimAllResDto
): AllocateSimAPIDto => {
  return {
    id: res.sim.id,
    iccid: res.sim.iccid,
    msisdn: res.sim.msisdn,
    type: res.sim.type as SIM_TYPES,
    is_physical: res.sim.is_physical,
    allocated_at: res.sim.allocated_at,
    firstActivatedOn: res.sim?.firstActivatedOn ?? "",
    lastActivatedOn: res.sim?.lastActivatedOn ?? "",
    activationsCount: res.sim.activationsCount,
    deactivationsCount: res.sim.deactivationsCount,
    subscriber_id: res.sim.subscriber_id,
    network_id: res.sim.network_id,
    package: res.sim?.package ?? {},
    imsi: res.sim.imsi,
    status: res.sim.status,
    traffic_policy: res.sim.traffic_policy,
    sync_status: res.sim.sync_status,
  };
};

export const dtoToAllocateSimDetailsDto = (response: any): SimDetailsDto => {
  const {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn,
    lastActivatedOn,
    activationsCount,
    deactivationsCount,
    allocatedAt,
  } = response;

  return {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn: firstActivatedOn?.toDate(),
    lastActivatedOn: lastActivatedOn?.toDate(),
    activationsCount,
    deactivationsCount,
    allocatedAt: allocatedAt?.toDate(),
  };
};

export const dtoToAllocateSimsDto = (
  res: SimsAlloAPIResDto
): SimsAlloAPIResDto => {
  const sims: AllocateSimAPIDto[] = [];
  for (const sim of res.sims) {
    sims.push(dtoToAllocateSimResDto({ sim: sim }));
  }
  return {
    sims: sims,
  };
};

export const dtoToSimResDto = (res: SimAPIResDto): SimDto => {
  return {
    activationCode: res.sim.activation_code,
    createdAt: res.sim.created_at,
    iccid: res.sim.iccid,
    id: res.sim.id,
    isAllocated: res.sim.is_allocated,
    isPhysical: res.sim.is_physical,
    msisdn: res.sim.msisdn,
    qrCode: res.sim.qr_code,
    simType: res.sim.sim_type as SIM_TYPES,
    smapAddress: res.sim.sm_ap_address,
  };
};

export const dtoToSimDetailsDto = (response: any): SimDetailsDto => {
  const {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn,
    lastActivatedOn,
    activationsCount,
    deactivationsCount,
    allocatedAt,
  } = response;

  return {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn: firstActivatedOn?.toDate(),
    lastActivatedOn: lastActivatedOn?.toDate(),
    activationsCount,
    deactivationsCount,
    allocatedAt: allocatedAt?.toDate(),
  };
};

export const mapSubscriberToSimsResDto = (
  resDto: SubscriberToSimsResDto
): SubscriberToSimsDto => {
  return {
    subscriberId: resDto.subscriber_id,
    sims: resDto.sims.map(sim => ({
      id: sim.id,
      subscriberId: sim.subscriber_id,
      networkId: sim.network_id,
      iccid: sim.iccid,
      msisdn: sim.msisdn,
      imsi: sim.imsi,
      type: sim.type,
      status: sim.status,
      isPhysical: sim.is_physical,
      trafficPolicy: sim.traffic_policy,
      firstActivatedOn: sim.firstActivatedOn,
      lastActivatedOn: sim.lastActivatedOn,
      activationsCount: sim.activationsCount,
      deactivationsCount: sim.deactivationsCount,
      allocatedAt: sim.allocated_at,
      syncStatus: sim.sync_status,
    })),
  };
};

export const dtoToSimsDto = (res: SimsAPIResDto): SimsResDto => {
  const sims: SimDto[] = [];
  for (const sim of res.sims) {
    sims.push(dtoToSimResDto({ sim: sim }));
  }
  return {
    sim: sims,
  };
};

export const dtoToSimPackagesDto = (
  res: GetSimPackagesDtoAPI
): SimPackagesResDto => {
  const packages: SimToPackagesResDto[] = [];
  for (const pkg of res.packages) {
    packages.push({
      id: pkg.id,
      packageId: pkg.package_id,
      startDate: pkg.start_date,
      endDate: pkg.end_date,
      isActive: pkg.is_active,
    });
  }
  return {
    simId: res.sim_id,
    packages: packages,
  };
};
