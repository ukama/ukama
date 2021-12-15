import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    OrgUserResponseDto,
    ResidentDto,
    ResidentResponse,
} from "./types";
import casual from "../../mockServer/mockData/casual-extensions";

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
        res.data.forEach(user => {
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
    dtoToUsersDto = (org: string): OrgUserResponseDto => {
        const users = casual.randomArray<GetUserDto>(2, 6, casual._getUser);
        return {
            orgName: org,
            users,
        };
    };
}
export default <IUserMapper>new UserMapper();
