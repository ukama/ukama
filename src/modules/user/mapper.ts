import { IUserMapper } from "./interface";
import { ConnectedUserDto } from "./types";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (res: ConnectedUserDto): ConnectedUserDto => {
        return {
            totalUser: Number(res.totalUser),
            residentUsers: Number(res.residentUsers),
            guestUsers: Number(res.guestUsers),
        };
    };
}
export default <IUserMapper>new UserMapper();
