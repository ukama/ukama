/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  GetSubscriberAPIResDto,
  SimPackageAPIDto,
  SimPackageDto,
  SimsAPIResDto,
  SubSimAPIDto,
  SubscriberAPIResDto,
  SubscriberDto,
  SubscriberSimDto,
  SubscriberSimsResDto,
  SubscribersAPIResDto,
  SubscribersResDto,
} from "../resolver/types";

export const addSubscriberReqToSubscriberResDto = (
  res: SubscriberAPIResDto
): SubscriberDto => {
  return {
    sim: [],
    email: res.Subscriber.email,
    gender: res.Subscriber.gender,
    address: res.Subscriber.address,
    dob: res.Subscriber.dob,
    phone: res.Subscriber.phone_number,
    idSerial: res.Subscriber.id_serial,
    uuid: res.Subscriber.subscriber_id,
    name: res.Subscriber.name,
    networkId: res.Subscriber.network_id,
    proofOfIdentification: res.Subscriber.proof_of_identification,
  };
};

export const dtoToSubscriberResDto = (
  res: GetSubscriberAPIResDto
): SubscriberDto => {
  const sims: SubscriberSimDto[] =
    res.subscriber.sim?.map(sim => ({
      id: sim.id,
      imsi: sim.imsi,
      type: sim.type,
      iccid: sim.iccid,
      msisdn: sim.msisdn,
      status: sim.status,
      package: sim.package,
      networkId: sim.network_id,
      isPhysical: sim.is_physical,
      sync_status: sim.sync_status,
      allocatedAt: sim.allocated_at,
      subscriberId: sim.subscriber_id,
    })) ?? [];
  return {
    sim: sims,
    email: res.subscriber.email,
    gender: res.subscriber.gender,
    address: res.subscriber.address,
    dob: res.subscriber.dob,
    phone: res.subscriber.phone_number,
    idSerial: res.subscriber.id_serial,
    uuid: res.subscriber.subscriber_id,
    name: res.subscriber.name,
    networkId: res.subscriber.network_id,
    proofOfIdentification: res.subscriber.proof_of_identification,
  };
};

export const dtoToSubscribersResDto = (
  res: SubscribersAPIResDto
): SubscribersResDto => {
  const subscribers: SubscriberDto[] = [];
  for (const subscriber of res.subscribers) {
    const sub = dtoToSubscriberResDto({ subscriber: subscriber });
    subscribers.push(sub);
  }

  return {
    subscribers: subscribers,
  };
};

export const dtoToSimPackageDto = (res: SimPackageAPIDto): SimPackageDto => {
  return {
    id: res.id,
    package_id: res.package_id,
    start_date: res.start_date,
    end_date: res.end_date,
    is_active: res.is_active,
    created_at: res.created_at,
    updated_at: res.updated_at,
  };
};

export const dtoToSimPackages = (res: SimPackageAPIDto[]): SimPackageDto[] => {
  return res.map(dtoToSimPackageDto);
};
export const dtoToSimDto = (res: SubSimAPIDto): SubscriberSimDto => {
  return {
    id: res.id,
    subscriberId: res.subscriber_id,
    networkId: res.network_id,
    iccid: res.iccid,
    msisdn: res.msisdn,
    imsi: res.imsi || "",
    type: res.type,
    status: res.status,
    allocatedAt: res.allocated_at,
    sync_status: res.sync_status || "",
    isPhysical: res.is_physical ?? false,
    package: dtoToSimPackages(res.package) || [],
  };
};

export const dtoToSimsResDto = (res: SimsAPIResDto): SubscriberSimsResDto => {
  const subSims: SubscriberSimDto[] = [];
  for (const sim of res.sims) {
    const s = dtoToSimDto(sim);
    subSims.push(s);
  }

  return { sims: subSims };
};
