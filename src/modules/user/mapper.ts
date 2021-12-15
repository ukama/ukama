import { DATA_PLAN_TYPE, GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    OrgUserResponse,
    OrgUserResponseDto,
    ResidentDto,
    ResidentResponse,
} from "./types";
import * as defaultCasual from "casual";

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
    dtoToUsersDto = (req: OrgUserResponse): OrgUserResponseDto => {
        const orgName = req.org;
        const res = req.users;
        const users: GetUserDto[] = [];

        res.forEach(user => {
            const node = {
                Default: "Default",
                Intermediate: "Intermediate",
            };
            const userObj = {
                id: user.uuid,
                name: `${user.firstName} ${user.lastName}`,
                email: user.email,
                status: defaultCasual.random_value(GET_STATUS_TYPE),
                node: `${defaultCasual.random_value(node)} Data Plan`,
                dataPlan: defaultCasual.random_value(DATA_PLAN_TYPE),
                dataUsage: defaultCasual.integer(1, 199),
                dlActivity: "Table cell",
                ulActivity: "Table cell",
            };
            users.push(userObj);
        });
        return {
            orgName,
            users,
        };
    };
}
export default <IUserMapper>new UserMapper();
