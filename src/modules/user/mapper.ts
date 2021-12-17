import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
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
    dtoToUsersDto = (req: OrgUserResponse): OrgUserResponseDto => {
        const orgName = req.org;
        const res = req.users;
        const users: GetUserDto[] = [];

        res.forEach(user => {
            const userObj = {
                id: user.uuid,
                name: `${user.firstName} ${user.lastName}`,
                email: user.email,
                status: defaultCasual.random_value(GET_STATUS_TYPE),
                eSimNumber: `# ${defaultCasual.integer(
                    11111,
                    99999
                )}-${defaultCasual.date("DD-MM-YYYY")}-${defaultCasual.integer(
                    1111111,
                    9999999
                )}`,
                iccid: `${defaultCasual.integer(
                    11111,
                    99999
                )}${defaultCasual.integer(11010, 99999)}${defaultCasual.integer(
                    11010,
                    99999
                )}`,
                dataPlan: defaultCasual.integer(5, 60),
                dataUsage: defaultCasual.integer(1, 39),
                roaming: defaultCasual.random_value([true, false]),
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
