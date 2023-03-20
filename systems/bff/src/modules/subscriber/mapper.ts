import { ISubscriberMapper } from "./interface";
import { SubscriberAPIResDto, SubscriberDto } from "./types";

class SubscriberMapper implements ISubscriberMapper {
    dtoToSubscriberResDto = (res: SubscriberAPIResDto): SubscriberDto => {
        return {
            uuid: res.Subscriber.subscriber_id,
            email: res.Subscriber.email,
            phone: res.Subscriber.phone_number,
            address: res.Subscriber.address,
            dob: res.Subscriber.date_of_birth,
            firstName: res.Subscriber.first_name,
            lastName: res.Subscriber.last_name,
            gender: res.Subscriber.gender,
            idSerial: res.Subscriber.id_serial,
            networkId: res.Subscriber.network_id,
            orgId: res.Subscriber.org_id,
            proofOfIdentification: res.Subscriber.proof_of_identification,
        };
    };
}
export default <ISubscriberMapper>new SubscriberMapper();
