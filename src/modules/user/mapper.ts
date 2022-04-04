import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    ConnectedUserDto,
    ConnectedUserResponse,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserResponse,
    OrgUsersResponse,
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
    dtoToUsersDto = (req: OrgUsersResponse): GetUsersDto[] => {
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
    dtoToUserDto = (req: OrgUserResponse): GetUserDto => {
        const { user, sim } = req;
        return {
            id: user.uuid,
            name: user.name,
            iccid: sim.iccid,
            email: user.email,
            phone: user.phone,
            status: sim.ukama?.status || GET_STATUS_TYPE.INACTIVE,
            roaming: sim.carrier?.status || GET_STATUS_TYPE.INACTIVE,
            dataPlan: 1024,
            dataUsage: defaultCasual.integer(1, 1024),
            eSimNumber: `# ${defaultCasual.integer(
                11111,
                99999
            )}-${defaultCasual.date("DD-MM-YYYY")}-${defaultCasual.integer(
                1111111,
                9999999
            )}`,
        };
    };
}
export default <IUserMapper>new UserMapper();
