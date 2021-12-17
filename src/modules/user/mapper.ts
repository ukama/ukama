import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    ResidentResponse,
} from "./types";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (res: ConnectedUserResponse): ConnectedUserDto => {
        const connectedUsers = res.data;
        return connectedUsers;
    };
    dtoToDto = (res: GetUserResponseDto): GetUserDto[] => {
        const users = res.data;
        return users;
    };
    residentDtoToDto = (res: GetUserResponseDto): ResidentResponse => {
        const residents: GetUserDto[] = [];
        let activeResidents = 0;
        const totalResidents = res.length;
        res.data.forEach(user => {
            if (user.status === GET_STATUS_TYPE.ACTIVE) {
                activeResidents++;
            }

            residents.push(user);
        });
        return {
            residents,
            activeResidents,
            totalResidents,
        };
    };
}
export default <IUserMapper>new UserMapper();
