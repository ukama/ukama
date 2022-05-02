import { GET_STATUS_TYPE } from "../../constants";
import { IUserMapper } from "./interface";
import {
    UserResDto,
    AddUserServiceRes,
    ConnectedUserDto,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserDto,
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
        const activeResidents = 0;
        const totalResidents = res.length;
        res.data.forEach(user => {
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
            if (!user.isDeactivated) {
                const userObj = {
                    id: user.uuid,
                    dataPlan: "0",
                    dataUsage: "0,",
                    name: user.name,
                    email: user.email,
                    phone: user.phone,
                    isDeactivated: user.isDeactivated,
                };
                users.push(userObj);
            }
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
            status:
                sim?.carrier?.status === GET_STATUS_TYPE.ACTIVE
                    ? sim?.carrier?.services.data
                    : false,
            roaming:
                sim?.carrier?.status === GET_STATUS_TYPE.ACTIVE ? true : false,
            dataPlan: 1024,
            dataUsage: defaultCasual.integer(1, 1024),
        };
    };
    dtoToUserResDto = (req: OrgUserDto): UserResDto => {
        return {
            id: req.uuid,
            name: req.name,
            email: req.email,
            phone: req.phone,
        };
    };
    dtoToAddUserDto = (req: AddUserServiceRes): UserResDto | null => {
        if (req) {
            return {
                name: req.user.name,
                email: req.user.email,
                phone: req.user.phone,
                id: req.user.uuid,
                iccid: req.iccid,
            };
        }
        return null;
    };
}
export default <IUserMapper>new UserMapper();
