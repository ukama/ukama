import { MetricServiceValueRes } from "../../common/types";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    GetUserDto,
    GetUserResponseDto,
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
}
export default <IUserMapper>new UserMapper();
