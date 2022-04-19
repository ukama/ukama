import { Service } from "typedi";
import {
    ActivateUserDto,
    ConnectedUserDto,
    DeactivateResponse,
    ResidentsResponse,
    UserInputDto,
    GetUserDto,
    ActivateUserResponse,
    GetUsersDto,
    UpdateUserServiceInput,
    UpdateUserServiceRes,
    UserResDto,
} from "./types";
import { IUserService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { API_METHOD_TYPE } from "../../constants";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { getPaginatedOutput } from "../../utils";
import { PaginationDto, ParsedCookie } from "../../common/types";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        cookie: ParsedCookie
    ): Promise<ConnectedUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/metrics/subscribersattached`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        const connectedUsers = UserMapper.connectedUsersDtoToDto(
            res.data.result
        );

        if (!connectedUsers) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return connectedUsers;
    };

    activateUser = async (
        req: ActivateUserDto
    ): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_ACTIVE_USER,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
    };

    updateUser = async (
        userId: string,
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${userId}`,
            headers: cookie.header,
            body: { name: req.name, email: req.email, phone: req.phone },
        });
        if (checkError(res)) throw new Error(res.message);

        return UserMapper.dtoToUserResDto(res);
    };
    deactivateUser = async (id: string): Promise<DeactivateResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_DEACTIVATE_USER,
            body: { id },
        });
        if (checkError(res)) throw new Error(res.message);
        return res.data;
    };
    getUser = async (
        userId: string,
        cookie: ParsedCookie
    ): Promise<GetUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${userId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);

        return UserMapper.dtoToUserDto(res);
    };
    getResidents = async (req: PaginationDto): Promise<ResidentsResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_USERS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const residents = UserMapper.residentDtoToDto(res);
        if (!residents) throw new HTTP404Error(Messages.RESIDENTS_NOT_FOUND);

        return {
            residents,
            meta,
        };
    };
    getUsersByOrg = async (cookie: ParsedCookie): Promise<GetUsersDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/users`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);

        return UserMapper.dtoToUsersDto(res);
    };
    addUser = async (
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto | null> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.ORG}/${cookie.orgId}/users`,
            body: { ...req, simToken: "I_DO_NOT_NEED_A_SIM" },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return UserMapper.dtoToAddUserDto(res);
    };
    deleteUser = async (
        userId: string,
        cookie: ParsedCookie
    ): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${userId}`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return {
            success: true,
        };
    };
    updateUserStatus = async (
        data: UpdateUserServiceInput,
        cookie: ParsedCookie
    ): Promise<UpdateUserServiceRes> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${data.userId}/sims/${data.simId}/services`,
            headers: cookie.header,
            body: {
                carrier: {
                    sms: false,
                    voice: false,
                    data: data.status,
                },
            },
        });
        if (checkError(res)) throw new Error(res.message);
        return {
            success: true,
        };
    };
}
