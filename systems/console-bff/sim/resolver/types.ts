/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { SIM_STATUS, SIM_TYPES } from "../../common/enums";

@InputType()
export class AllocateSimInputDto {
  @Field()
  network_id: string;

  @Field()
  package_id: string;

  @Field({ nullable: true })
  iccid?: string;

  @Field()
  sim_type: SIM_TYPES;

  @Field()
  subscriber_id: string;

  @Field()
  traffic_policy: number;
}

@InputType()
export class SimUsageInputDto {
  @Field()
  iccid: string;

  @Field()
  simId: string;

  @Field()
  type: string;
}

@InputType()
export class SimUsagesInputDto {
  @Field(() => String)
  networkId: string;

  @Field()
  type: string;
}

@ObjectType()
export class SimPackage {
  @Field()
  id: string;

  @Field()
  packageId: string;

  @Field()
  startDate: string;

  @Field()
  endDate: string;

  @Field()
  defaultDuration: string;

  @Field()
  isActive: boolean;

  @Field()
  asExpired: boolean;
}

@ObjectType()
export class SimDetailsDto {
  @Field()
  id: string;

  @Field({ nullable: true })
  subscriberId: string;

  @Field()
  networkId: string;

  @Field()
  orgId: string;

  @Field(() => SimPackage)
  Package: SimPackage;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  isPhysical: boolean;

  @Field()
  allocatedAt: string;
}

@ObjectType()
export class SimStatusResDto {
  @Field(() => String, { nullable: true })
  simId?: string;
}
@ObjectType()
export class DeleteSimResDto {
  @Field(() => String, { nullable: true })
  simId?: string;
}
@ObjectType()
export class RemovePackageFromSimResDto {
  @Field(() => String, { nullable: true })
  packageId?: string;
}
@ObjectType()
export class AddPackagSimResDto {
  @Field(() => String, { nullable: true })
  packageId: string;
}

@ObjectType()
export class AddPackagesSimResDto {
  @Field(() => [AddPackagSimResDto])
  packages: AddPackagSimResDto[];
}

@InputType()
export class RemovePackageFormSimInputDto {
  @Field()
  simId: string;

  @Field()
  packageId: string;
}

@InputType()
export class GetPackagesForSimInputDto {
  @Field()
  sim_id: string;
}

@ObjectType()
export class GetSimPackagesDtoAPI {
  @Field()
  sim_id: string;
  @Field(() => [SimToPackagesDto])
  packages: SimToPackagesDto[];
}
@ObjectType()
export class SimToPackagesDto {
  @Field()
  id: string;

  @Field()
  package_id: string;

  @Field()
  start_date: string;

  @Field()
  end_date: string;

  @Field()
  is_active: boolean;
}
@ObjectType()
export class SimToPackagesResDto {
  @Field()
  id: string;

  @Field()
  packageId: string;

  @Field()
  startDate: string;

  @Field()
  endDate: string;

  @Field()
  isActive: boolean;
}
@InputType()
export class ToggleSimStatusInputDto {
  @Field()
  sim_id: string;

  @Field()
  status: string;
}
@InputType()
export class GetSimInputDto {
  @Field()
  simId: string;
}
@InputType()
export class GetSimBySubscriberInputDto {
  @Field()
  subscriberId: string;
}
@InputType()
export class GetSimBySubscriberIdInputDto {
  @Field()
  subscriberId: string;
}
@InputType()
export class GetSimByNetworkInputDto {
  @Field()
  networkId: string;
}
@InputType()
export class DeleteSimInputDto {
  @Field()
  simId: string;
}

@InputType()
export class PackagesToSimInputDto {
  @Field()
  package_id: string;

  @Field()
  start_date: string;
}
@InputType()
export class AddPackagesToSimInputDto {
  @Field()
  sim_id: string;

  @Field(() => [PackagesToSimInputDto])
  packages: PackagesToSimInputDto[];
}

@ObjectType()
export class SimAllocatePackageDto {
  @Field({ nullable: true })
  id?: string;

  @Field({ nullable: true })
  packageId?: string;

  @Field({ nullable: true })
  startDate?: string;

  @Field({ nullable: true })
  endDate?: string;

  @Field({ nullable: true })
  isActive?: boolean;
}
@ObjectType()
export class AllocateSimAPIDto {
  @Field()
  id: string;

  @Field()
  subscriber_id: string;

  @Field()
  network_id: string;

  @Field(() => SimAllocatePackageDto)
  package: SimAllocatePackageDto;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field({ nullable: true })
  imsi?: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  is_physical: boolean;

  @Field()
  traffic_policy: number;

  @Field()
  allocated_at: string;

  @Field()
  sync_status: string;
}

@ObjectType()
export class SimPackageAPI {
  @Field()
  id: string;

  @Field()
  package_id: string;

  @Field()
  start_date: string;

  @Field()
  end_date: string;

  @Field()
  default_duration: string;

  @Field()
  is_active: boolean;

  @Field()
  as_expired: boolean;
}
@ObjectType()
export class SimAPIDto {
  @Field()
  id: string;

  @Field()
  subscriber_id: string;

  @Field()
  network_id: string;

  @Field(() => SimPackageAPI, { nullable: true })
  package?: SimPackageAPI;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  is_physical: boolean;

  @Field()
  traffic_policy: number;

  @Field()
  firstActivatedOn: string;

  @Field()
  lastActivatedOn: string;

  @Field()
  activationsCount: string;

  @Field()
  deactivationsCount: string;

  @Field()
  allocated_at: string;

  @Field()
  sync_status: string;
}

@ObjectType()
export class SubscriberToSimsResDto {
  @Field()
  subscriber_id: string;
  @Field(() => [SubscriberSimsAPIDto])
  sims: SubscriberSimsAPIDto[];
}

@ObjectType()
export class SubscriberToSimsDto {
  @Field()
  subscriberId: string;
  @Field(() => [SubscriberSimsDto])
  sims: SubscriberSimsDto[];
}

@ObjectType()
export class SubscriberSimsDto {
  @Field()
  id: string;

  @Field()
  subscriberId: string;

  @Field()
  networkId: string;

  @Field()
  syncStatus: string;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  isPhysical: boolean;

  @Field()
  trafficPolicy: number;

  @Field()
  allocatedAt: string;
}
@ObjectType()
export class SubscriberSimsAPIDto {
  @Field()
  id: string;

  @Field()
  subscriber_id: string;

  @Field()
  network_id: string;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  is_physical: boolean;

  @Field()
  traffic_policy: number;

  @Field()
  allocated_at: string;

  @Field()
  sync_status: string;
}

@ObjectType()
export class SimDto {
  @Field()
  id: string;

  @Field()
  subscriberId: string;

  @Field()
  networkId: string;

  @Field(() => SimPackage, { nullable: true })
  package?: SimPackage;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status: string;

  @Field()
  isPhysical: boolean;

  @Field()
  trafficPolicy: number;

  @Field()
  firstActivatedOn: string;

  @Field()
  lastActivatedOn: string;

  @Field()
  activationsCount: string;

  @Field()
  deactivationsCount: string;

  @Field()
  allocatedAt: string;

  @Field()
  syncStatus: string;
}

@ObjectType()
export class SimAPIResDto {
  @Field(() => SimAPIDto)
  sim: SimAPIDto;
}
export class SimAllResDto {
  @Field(() => AllocateSimAPIDto)
  sim: AllocateSimAPIDto;
}
@ObjectType()
export class SimsAlloAPIResDto {
  @Field(() => [SimAPIDto])
  sims: AllocateSimAPIDto[];
}

@ObjectType()
export class SimsAPIResDto {
  @Field(() => [SimAPIDto])
  sims: SimAPIDto[];
}
@ObjectType()
export class SimsResDto {
  @Field(() => [SimDto])
  sims: SimDto[];
}

@ObjectType()
export class SimDataUsage {
  @Field()
  usage: string;

  @Field()
  simId: string;
}

@ObjectType()
export class SimDataUsages {
  @Field(() => [SimDataUsage])
  usages: SimDataUsage[];
}

@ObjectType()
export class SimPoolStatsDto {
  @Field()
  total: number;

  @Field()
  available: number;

  @Field()
  consumed: number;

  @Field()
  failed: number;

  @Field()
  esim: number;

  @Field()
  physical: number;
}

@ObjectType()
export class UploadSimsResDto {
  @Field(() => [String])
  iccid: string[];
}

@InputType()
export class UploadSimsInputDto {
  @Field()
  data: string;

  @Field(() => SIM_TYPES)
  simType: SIM_TYPES;
}

@InputType()
export class GetSimsInput {
  @Field(() => SIM_TYPES)
  type: SIM_TYPES;

  @Field(() => SIM_STATUS)
  status: SIM_STATUS;
}

@InputType()
export class ListSimsInput {
  @Field()
  status: string;

  @Field()
  networkId: string;
}

@ObjectType()
export class SimPoolResDto {
  @Field()
  id: string;

  @Field()
  qrCode: string;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  isAllocated: boolean;

  @Field()
  isFailed: boolean;

  @Field()
  simType: string;

  @Field()
  smApAddress: string;

  @Field()
  activationCode: string;

  @Field()
  createdAt: string;

  @Field()
  deletedAt: string;

  @Field()
  updatedAt: string;

  @Field()
  isPhysical: boolean;
}

@ObjectType()
export class SimsPoolResDto {
  @Field(() => [SimPoolResDto])
  sims: SimPoolResDto[];
}
