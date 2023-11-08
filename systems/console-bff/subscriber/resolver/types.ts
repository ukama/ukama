/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IsEmail, IsPhoneNumber } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class SubSimAPIDto {
  @Field()
  id: string;

  @Field()
  subscriber_id: string;

  @Field()
  network_id: string;

  @Field()
  org_id: string;

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

  @Field({ nullable: true })
  first_activated_on: string;

  @Field({ nullable: true })
  last_activated_on: string;

  @Field()
  activations_count: string;

  @Field()
  deactivations_count: string;

  @Field()
  allocated_at: string;

  @Field({ nullable: true })
  is_physical: boolean;

  @Field({ nullable: true })
  package: string;
}

@ObjectType()
export class SubscriberAPIDto {
  @Field()
  subscriber_id: string;

  @Field()
  address: string;

  @Field()
  date_of_birth: string;

  @Field()
  email: string;

  @Field()
  first_name: string;

  @Field()
  last_name: string;

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

  @Field(() => [SubSimAPIDto])
  sim: SubSimAPIDto[];
}

@ObjectType()
export class SubscriberAPIResDto {
  @Field(() => SubscriberAPIDto)
  subscriber: SubscriberAPIDto;
}

export class SubscribersAPIResDto {
  @Field(() => [SubscriberAPIDto])
  subscribers: SubscriberAPIDto[];
}

@InputType()
export class SubscriberInputDto {
  @Field({ nullable: true })
  address?: string;

  @Field({ nullable: true })
  dob?: string;

  @Field()
  @IsEmail()
  email: string;

  @Field({ nullable: true })
  first_name?: string;

  @Field({ nullable: true })
  last_name?: string;

  @Field({ nullable: true })
  gender?: string;

  @Field({ nullable: true })
  id_serial?: string;

  @Field()
  network_id: string;

  @Field()
  org_id: string;

  @Field({ nullable: true })
  // @IsPhoneNumber()
  phone?: string;

  @Field({ nullable: true })
  proof_of_identification?: string;
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
  orgId: string;

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

  @Field({ nullable: true })
  firstActivatedOn: string;

  @Field({ nullable: true })
  lastActivatedOn: string;

  @Field()
  activationsCount: string;

  @Field()
  deactivationsCount: string;

  @Field()
  allocatedAt: string;

  @Field({ nullable: true })
  isPhysical: boolean;

  @Field({ nullable: true })
  package: string;
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
  firstName: string;

  @Field()
  lastName: string;

  @Field()
  gender: string;

  @Field()
  idSerial?: string;

  @Field()
  networkId: string;

  @Field()
  orgId: string;

  @Field()
  phone: string;

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
  first_name: string;

  @Field({ nullable: true })
  @IsPhoneNumber()
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
