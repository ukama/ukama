/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  IsBoolean,
  IsInt,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsUUID,
  Min,
} from "class-validator";
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

@ObjectType()
export class PackageRateAPIDto {
  @Field()
  sms_mo: string;

  @Field()
  sms_mt: number;

  @Field()
  data: number;

  @Field()
  amount: number;
}
@ObjectType()
export class PackageMarkupAPIDto {
  @Field()
  baserate: string;

  @Field()
  markup: number;
}

@ObjectType()
export class PackageAPIDto {
  @Field()
  uuid: string;

  @Field()
  name: string;

  @Field()
  active: boolean;

  @Field()
  duration: number;

  @Field()
  sim_type: string;

  @Field()
  created_at: string;

  @Field()
  deleted_at: string;

  @Field()
  updated_at: string;

  @Field()
  sms_volume: number;

  @Field()
  data_volume: number;

  @Field()
  voice_volume: number;

  @Field()
  ulbr: string;

  @Field()
  dlbr: string;

  @Field()
  type: string;

  @Field()
  data_unit: string;

  @Field()
  voice_unit: string;

  @Field()
  message_unit: string;

  @Field()
  flatrate: boolean;

  @Field()
  currency: string;

  @Field()
  from: string;

  @Field()
  to: string;

  @Field()
  country: string;

  @Field()
  provider: string;

  @Field()
  apn: string;

  @Field()
  owner_id: string;

  @Field({ nullable: true })
  network_id: string;

  @Field()
  amount: number;

  @Field(() => PackageRateAPIDto)
  rate: PackageRateAPIDto;

  @Field(() => PackageMarkupAPIDto)
  markup: PackageMarkupAPIDto;
}
@ObjectType()
export class PackageAPIResDto {
  @Field()
  package: PackageAPIDto;
}

@ObjectType()
export class PackagesAPIResDto {
  @Field(() => [PackageAPIDto])
  packages: PackageAPIDto[];
}

@ObjectType()
export class PackageRateDto {
  @Field()
  smsMo: string;

  @Field()
  smsMt: number;

  @Field()
  data: number;

  @Field()
  amount: number;
}
@ObjectType()
export class PackageMarkupDto {
  @Field()
  baserate: string;

  @Field()
  markup: number;
}

@ObjectType()
export class PackageDto {
  @Field()
  uuid: string;

  @Field()
  name: string;

  @Field()
  active: boolean;

  @Field()
  duration: number;

  @Field()
  simType: string;

  @Field()
  createdAt: string;

  @Field()
  deletedAt: string;

  @Field()
  updatedAt: string;

  @Field()
  smsVolume: number;

  @Field()
  dataVolume: number;

  @Field()
  voiceVolume: number;

  @Field()
  ulbr: string;

  @Field()
  dlbr: string;

  @Field()
  type: string;

  @Field()
  dataUnit: string;

  @Field()
  voiceUnit: string;

  @Field()
  messageUnit: string;

  @Field()
  flatrate: boolean;

  @Field()
  currency: string;

  @Field()
  from: string;

  @Field()
  to: string;

  @Field()
  country: string;

  @Field()
  provider: string;

  @Field()
  apn: string;

  @Field()
  ownerId: string;

  @Field({ nullable: true })
  networkId: string;

  @Field()
  amount: number;

  @Field(() => PackageRateAPIDto)
  rate: PackageRateAPIDto;

  @Field(() => PackageMarkupAPIDto)
  markup: PackageMarkupAPIDto;
}

@ObjectType()
export class PackagesResDto {
  @Field(() => [PackageDto])
  packages: PackageDto[];
}

@InputType()
export class AddPackageInputDto {
  @Field()
  @IsNotEmpty()
  name: string;

  @Field(() => Int)
  @IsInt()
  @Min(1)
  duration: number;

  @Field()
  @IsNotEmpty()
  dataUnit: string;

  @Field(() => Float)
  @IsNumber()
  @Min(0)
  amount: number;

  @Field(() => Int)
  @IsInt()
  @Min(0)
  dataVolume: number;

  @Field()
  @IsNotEmpty()
  country: string;

  @Field()
  @IsNotEmpty()
  currency: string;

  @Field({ nullable: true })
  @IsOptional()
  @IsUUID()
  networkId?: string;
}

@InputType()
export class UpdatePackageInputDto {
  @Field()
  @IsNotEmpty()
  name: string;

  @Field()
  @IsBoolean()
  active: boolean;
}

@ObjectType()
export class PackageNameAvailabilityAPIDto {
  @Field()
  is_available: boolean;

  @Field()
  name: string;
}

@ObjectType()
export class PackageNameAvailabilityResDto {
  @Field()
  isAvailable: boolean;

  @Field()
  name: string;
}
