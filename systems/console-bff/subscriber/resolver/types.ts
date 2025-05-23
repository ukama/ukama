/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IsEmail } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class SimPackageAPIDto {
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

  @Field()
  created_at: string;

  @Field()
  updated_at: string;
}

@ObjectType()
export class SubSimAPIDto {
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
  sync_status: string;

  @Field()
  allocated_at: string;

  @Field({ nullable: true })
  is_physical: boolean;

  @Field(() => SimPackageAPIDto)
  package: SimPackageAPIDto;
}

@ObjectType()
export class SubscriberAPIDto {
  @Field()
  subscriber_id: string;

  @Field()
  address: string;

  @Field()
  dob: string;

  @Field()
  email: string;

  @Field()
  name: string;

  @Field()
  gender: string;

  @Field()
  id_serial: string;

  @Field()
  network_id: string;

  @Field()
  org_id: string;

  @Field()
  phone_number: string;

  @Field()
  proof_of_identification: string;

  @Field()
  subscriber_status: string;

  @Field(() => [SubSimAPIDto])
  sim: SubSimAPIDto[];
}

@ObjectType()
export class SimsAPIResDto {
  @Field(() => [SubSimAPIDto])
  sims: SubSimAPIDto[];
}

@ObjectType()
export class SubscriberAPIResDto {
  @Field(() => SubscriberAPIDto)
  Subscriber: SubscriberAPIDto;
}

@ObjectType()
export class GetSubscriberAPIResDto {
  @Field(() => SubscriberAPIDto)
  subscriber: SubscriberAPIDto;
}

export class SubscribersAPIResDto {
  @Field(() => [SubscriberAPIDto])
  subscribers: SubscriberAPIDto[];
}

@InputType()
export class SubscriberInputDto {
  @Field()
  @IsEmail()
  email: string;

  @Field()
  name: string;

  @Field()
  network_id: string;

  @Field({ nullable: true })
  phone?: string;
}

@ObjectType()
export class SimPackageDto {
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

  @Field()
  created_at: string;

  @Field()
  updated_at: string;
}

@ObjectType()
export class SubscriberSimDto {
  @Field()
  id: string;

  @Field()
  subscriberId: string;

  @Field()
  networkId: string;

  @Field()
  iccid: string;

  @Field()
  msisdn: string;

  @Field()
  imsi: string;

  @Field()
  type: string;

  @Field()
  status?: string;

  @Field()
  allocatedAt: string;

  @Field({ nullable: true })
  sync_status?: string;

  @Field({ nullable: true })
  isPhysical: boolean;

  @Field(() => SimPackageDto, { nullable: true })
  package?: SimPackageDto;
}
@ObjectType()
export class SubscriberDto {
  @Field()
  uuid: string;

  @Field()
  address: string;

  @Field()
  dob: string;

  @Field()
  email: string;

  @Field()
  name: string;

  @Field()
  gender: string;

  @Field()
  idSerial?: string;

  @Field()
  networkId: string;

  @Field()
  phone: string;

  @Field()
  subscriberStatus: string;

  @Field()
  proofOfIdentification: string;

  @Field(() => [SubscriberSimDto], { nullable: true })
  sim?: SubscriberSimDto[];
}

@ObjectType()
export class SubscribersResDto {
  @Field(() => [SubscriberDto])
  subscribers: SubscriberDto[];
}

@ObjectType()
export class SubscriberSimsResDto {
  @Field(() => [SubscriberSimDto])
  sims: SubscriberSimDto[];
}

@InputType()
export class UpdateSubscriberInputDto {
  @Field({ nullable: true })
  address: string;

  @Field({ nullable: true })
  @IsEmail()
  email: string;

  @Field({ nullable: true })
  id_serial: string;

  @Field({ nullable: true })
  name: string;

  @Field({ nullable: true })
  phone: string;

  @Field({ nullable: true })
  proof_of_identification: string;
}

@ObjectType()
export class SubscriberMetricsByNetworkDto {
  @Field()
  total: number;

  @Field()
  active: number;

  @Field()
  inactive: number;

  @Field()
  terminated: number;
}
