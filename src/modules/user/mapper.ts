import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    AddUserResponse,
    AddUserServiceRes,
    ConnectedUserDto,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserResponse,
    OrgUsersResponse,
    ResidentResponse,
} from "./types";
import * as defaultCasual from "casual";
import { MetricServiceRes } from "../../common/types";

class UserMapper implements IUserMapper {
    connectedUsersDtoToDto = (res: MetricServiceRes[]): ConnectedUserDto => {
        if (res.length > 0) {
            const value: any = res[0].value[1];
            return { totalUser: value };
        }
        return { totalUser: "0" };
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
            iccid: sim?.iccid || "",
            email: user.email,
            phone: user.phone,
            eSimNumber: user.uuid,
            status: sim?.ukama?.status || GET_STATUS_TYPE.INACTIVE,
            roaming:
                sim?.carrier?.status === GET_STATUS_TYPE.ACTIVE ? true : false,
            dataPlan: 1024,
            dataUsage: defaultCasual.integer(1, 1024),
        };
    };
    dtoToAddUserDto = (req: AddUserServiceRes): AddUserResponse | null => {
        if (req) {
            return {
                name: req.user.name,
                email: req.user.email,
                phone: req.user.phone,
                uuid: req.user.uuid,
                iccid: req.iccid,
            };
        }
        return null;
    };
}
export default <IUserMapper>new UserMapper();
