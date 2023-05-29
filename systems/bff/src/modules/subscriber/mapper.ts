import { ISubscriberMapper } from "./interface";
import {
    SubscriberAPIResDto,
    SubscriberDto,
    SubscriberSimDto,
    SubscribersAPIResDto,
    SubscribersResDto,
} from "./types";

class SubscriberMapper implements ISubscriberMapper {
    dtoToSubscriberResDto = (res: SubscriberAPIResDto): SubscriberDto => {
        const sims: SubscriberSimDto[] = [];
        if (res.subscriber.sim.length > 0) {
            for (const sim of res.subscriber.sim) {
                sims.push({
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
                });
            }
        }
        return {
            sim: sims,
            email: res.subscriber.email,
            orgId: res.subscriber.org_id,
            gender: res.subscriber.gender,
            address: res.subscriber.address,
            dob: res.subscriber.date_of_birth,
            phone: res.subscriber.phone_number,
            idSerial: res.subscriber.id_serial,
            uuid: res.subscriber.subscriber_id,
            lastName: res.subscriber.last_name,
            firstName: res.subscriber.first_name,
            networkId: res.subscriber.network_id,
            proofOfIdentification: res.subscriber.proof_of_identification,
        };
    };
    dtoToSubscribersResDto = (res: SubscribersAPIResDto): SubscribersResDto => {
        const subscribers: SubscriberDto[] = [];
        for (const subscriber of res.subscribers) {
            const sub = this.dtoToSubscriberResDto({ subscriber: subscriber });
            subscribers.push(sub);
        }

        return {
            subscribers: subscribers,
        };
    };
}
export default <ISubscriberMapper>new SubscriberMapper();
