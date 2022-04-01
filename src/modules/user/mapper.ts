import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserResponse,
    OrgUserResponseDto,
    ResidentResponse,
} from "./types";
import * as defaultCasual from "casual";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (res: ConnectedUserResponse): ConnectedUserDto => {
        return res.data;
    };
    dtoToDto = (res: GetUserResponseDto): GetUserDto[] => {
        return res.data;
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
    dtoToUsersDto = (req: OrgUserResponse): GetUsersDto[] => {
        const orgName = req.org;
        const res = req.users;
        const users: GetUsersDto[] = [];

        res.forEach(user => {
            const userObj = {
                id: user.uuid,
                name: user.name,
                email: user.email,
                phone: user.phone,
                dataPlan: 1024,
                dataUsage: defaultCasual.integer(1, 1024),
            };
            users.push(userObj);
        });
        return users;
    };
}
export default <IUserMapper>new UserMapper();
