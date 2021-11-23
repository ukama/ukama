import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    ResidentDto,
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
        const residents: ResidentDto[] = [];
        let activeResidents = 0;
        const totalResidents = res.length;
        res.data.map(user => {
            if (user.status === GET_STATUS_TYPE.ACTIVE) {
                activeResidents++;
            }
            const resident: ResidentDto = {
                id: user.id,
                name: user.name,
                dataUsage: user.dataUsage,
            };
            residents.push(resident);
        });
        return {
            residents,
            activeResidents,
            totalResidents,
        };
    };
}
export default <IUserMapper>new UserMapper();
