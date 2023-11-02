import {
  SubscriberAPIResDto,
  SubscriberDto,
  SubscriberSimDto,
  SubscribersAPIResDto,
  SubscribersResDto,
} from "../resolver/types";

export const dtoToSubscriberResDto = (
  res: SubscriberAPIResDto
): SubscriberDto => {
  const sims: SubscriberSimDto[] =
    res.Subscriber.sim?.map(sim => ({
      id: sim.id,
      imsi: sim.imsi,
      type: sim.type,
      iccid: sim.iccid,
      orgId: sim.org_id,
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
    })) || [];
  return {
    sim: sims,
    email: res.Subscriber.email,
    orgId: res.Subscriber.org_id,
    gender: res.Subscriber.gender,
    address: res.Subscriber.address,
    dob: res.Subscriber.date_of_birth,
    phone: res.Subscriber.phone_number,
    idSerial: res.Subscriber.id_serial,
    uuid: res.Subscriber.subscriber_id,
    lastName: res.Subscriber.last_name,
    firstName: res.Subscriber.first_name,
    networkId: res.Subscriber.network_id,
    proofOfIdentification: res.Subscriber.proof_of_identification,
  };
};

export const dtoToSubscribersResDto = (
  res: SubscribersAPIResDto
): SubscribersResDto => {
  const subscribers: SubscriberDto[] = [];
  for (const subscriber of res.subscribers) {
    const sub = dtoToSubscriberResDto({ Subscriber: subscriber });
    subscribers.push(sub);
  }

  return {
    subscribers: subscribers,
  };
};
