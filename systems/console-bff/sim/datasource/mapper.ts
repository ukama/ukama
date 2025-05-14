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
  SimAPIDto,
  SimAllResDto,
  SimDataUsage,
  SimDetailsDto,
  SimDto,
  SimPackage,
  SimPackageAPI,
  SimPoolResDto,
  SimUsageInputDto,
  SimsAPIResDto,
  SimsAlloAPIResDto,
  SimsPoolResDto,
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

export const dtoToSimResDto = (res: SimAPIDto): SimDto => {
  return {
    id: res.id,
    subscriberId: res.subscriber_id,
    networkId: res.network_id,
    iccid: res.iccid,
    msisdn: res.msisdn,
    imsi: res.imsi,
    type: res.type,
    status: res.status,
    isPhysical: res.is_physical,
    trafficPolicy: res.traffic_policy,
    allocatedAt: res.allocated_at,
    syncStatus: res.sync_status,
    firstActivatedOn: res.firstActivatedOn,
    lastActivatedOn: res.lastActivatedOn,
    activationsCount: res.activationsCount,
    deactivationsCount: res.deactivationsCount,
    package: res.package ? dtoToSimPackageDto(res.package) : undefined,
  };
};

const dtoToSimPackageDto = (res: SimPackageAPI): SimPackage => {
  return {
    id: res.id,
    endDate: res.end_date,
    isActive: res.is_active,
    packageId: res.package_id,
    startDate: res.start_date,
    asExpired: res.as_expired,
    defaultDuration: res.default_duration,
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
      allocatedAt: sim.allocated_at,
      syncStatus: sim.sync_status,
    })),
  };
};

export const dtoToSimsDto = (res: SimsAPIResDto): SimsResDto => {
  const sims: SimDto[] = [];
  for (const sim of res.sims) {
    sims.push(dtoToSimResDto(sim));
  }
  return {
    sims: sims,
  };
};

export const dtoToUsageDto = (
  res: any,
  args: SimUsageInputDto
): SimDataUsage => {
  const data = res.usage;
  return {
    simId: args.simId,
    usage: data[args.iccid] ?? "0",
  };
};

export const dtoToSimsFromPoolDto = (res: any): SimsPoolResDto => {
  const sims: SimPoolResDto[] = [];
  for (const sim of res.sims) {
    sims.push({
      id: sim.id,
      iccid: sim.iccid,
      msisdn: sim.msisdn,
      qrCode: sim.qr_code,
      simType: sim.sim_type,
      isFailed: sim.is_failed,
      updatedAt: sim.updated_at,
      createdAt: sim.created_at,
      deletedAt: sim.deleted_at,
      isPhysical: sim.is_physical,
      isAllocated: sim.is_allocated,
      smApAddress: sim.sm_ap_address,
      activationCode: sim.activation_code,
    });
  }
  return { sims: sims };
};
