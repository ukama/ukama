import { MetricServiceValueRes } from "../../common/types";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    GetUserDto,
    GetUserResponseDto,
    SubscriberAPIResDto,
    SubscriberDto,
    UserAPIResDto,
    UserResDto,
} from "./types";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (
        res: MetricServiceValueRes[],
    ): ConnectedUserDto => {
        if (res.length > 0) {
            const value: any = res[0].value[1];
            return { totalUser: value };
        }
        return { totalUser: "0" };
    };
    dtoToDto = (res: GetUserResponseDto): GetUserDto[] => {
        return res.data;
    };
    dtoToUserResDto = (res: UserAPIResDto): UserResDto => {
        return {
            uuid: res.user.uuid,
            email: res.user.email,
            isDeactivated: res.user.is_deactivated,
            name: res.user.name,
            phone: res.user.phone,
            registeredSince: res.user.registered_since,
        };
    };
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
export default <IUserMapper>new UserMapper();
