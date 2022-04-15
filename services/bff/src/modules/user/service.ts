import { Service } from "typedi";
import {
    ActivateUserDto,
    ConnectedUserDto,
    DeactivateResponse,
    GetUserPaginationDto,
    GetUserResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserResponse,
    GetUserDto,
    AddUserDto,
    AddUserResponse,
    ActivateUserResponse,
    GetUsersDto,
    UserInput,
    UpdateUserServiceInput,
    UpdateUserServiceRes,
} from "./types";
import { IUserService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { API_METHOD_TYPE } from "../../constants";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { getPaginatedOutput } from "../../utils";
import { HeaderType, PaginationDto } from "../../common/types";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        orgId: string,
        header: HeaderType
    ): Promise<ConnectedUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${orgId}/metrics/subscribersattached`,
            headers: header,
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

    updateUser = async (req: UpdateUserDto): Promise<UserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.POST_UPDATE_USER,
            body: req,
        });
        if (checkError(res)) throw new Error(res.message);

        return res.data;
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
        data: UserInput,
        header: HeaderType
    ): Promise<GetUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${data.orgId}/users/${data.userId}`,
            headers: header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);

        return UserMapper.dtoToUserDto(res);
    };

    getUsers = async (req: GetUserPaginationDto): Promise<GetUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_USERS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const users = UserMapper.dtoToDto(res);
        if (!users) throw new HTTP404Error(Messages.USERS_NOT_FOUND);

        return {
            users,
            meta,
        };
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
    getUsersByOrg = async (
        orgId: string,
        header: HeaderType
    ): Promise<GetUsersDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${orgId}/users`,
            headers: header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);

        return UserMapper.dtoToUsersDto(res);
    };
    addUser = async (
        orgId: string,
        req: AddUserDto,
        header: HeaderType
    ): Promise<AddUserResponse | null> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.ORG}/${orgId}/users`,
            body: { ...req, simToken: "I_DO_NOT_NEED_A_SIM" },
            headers: header,
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return UserMapper.dtoToAddUserDto(res);
    };
    deleteUser = async (
        orgId: string,
        userId: string,
        header: HeaderType
    ): Promise<ActivateUserResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.ORG}/${orgId}/users/${userId}`,
            headers: header,
        });
        if (checkError(res)) throw new Error(res.message);
        return {
            success: true,
        };
    };
    updateUserStatus = async (
        data: UpdateUserServiceInput,
        header: HeaderType
    ): Promise<UpdateUserServiceRes> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${data.orgId}/users/${data.userId}/sims/${data.simId}/services`,
            headers: header,
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
