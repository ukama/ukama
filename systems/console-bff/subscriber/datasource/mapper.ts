/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  GetSubscriberAPIResDto,
  SubscriberAPIResDto,
  SubscriberDto,
  SubscriberSimDto,
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
      allocatedAt: sim.allocated_at,
      subscriberId: sim.subscriber_id,
      lastActivatedOn: sim.last_activated_on,
      activationsCount: sim.activations_count,
      firstActivatedOn: sim.first_activated_on,
      deactivationsCount: sim.deactivations_count,
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
