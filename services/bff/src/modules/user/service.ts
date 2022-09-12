import { Service } from "typedi";
import {
    ConnectedUserDto,
    DeactivateResponse,
    UserInputDto,
    GetUserDto,
    ActivateUserResponse,
    GetUsersDto,
    UserFistVisitResDto,
    UpdateUserServiceInput,
    UserResDto,
    UserFistVisitInputDto,
    OrgUserSimDto,
    GetAccountDetailsDto,
    GetESimQRCodeInput,
    ESimQRCodeRes,
} from "./types";
import { IUserService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import UserMapper from "./mapper";
import { API_METHOD_TYPE } from "../../constants";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { ParsedCookie } from "../../common/types";

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

    updateUser = async (
        userId: string,
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${userId}`,
            headers: cookie.header,
            body: { name: req.name, email: req.email },
        });
        if (checkError(res)) throw new Error(res.message);
        return UserMapper.dtoToUserResDto(res);
    };
    deactivateUser = async (
        id: string,
        cookie: ParsedCookie
    ): Promise<DeactivateResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${id}`,
            headers: cookie.header,
            body: { isDeactivated: true },
        });
        if (checkError(res)) throw new Error(res.description);
        return res;
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
    getAccountDetails = async (
        cookie: ParsedCookie
    ): Promise<GetAccountDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.GET_IDENTITY}/${cookie.orgId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(res.Messages);

        return {
            email: res?.traits.email,
            isFirstVisit: res?.traits.firstVisit,
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
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.ORG}/${cookie.orgId}/users`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return UserMapper.dtoToAddUserDto(res);
    };
    updateFirstVisit = async (
        req: UserFistVisitInputDto,
        cookie: ParsedCookie
    ): Promise<UserFistVisitResDto> => {
        const getUser = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.GET_IDENTITY}/${cookie.orgId}`,
            headers: cookie.header,
        });
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.GET_IDENTITY}/${cookie.orgId}`,
            body: JSON.stringify({
                schema_id: "default",
                state: "active",
                traits: {
                    email: getUser?.traits?.email,
                    ...req,
                },
            }),
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return {
            firstVisit: res?.traits?.firstVisit,
        };
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
    ): Promise<OrgUserSimDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${data.userId}/sims/${data.simId}/services`,
            headers: cookie.header,
            body: {
                ukama: {
                    sms: false,
                    voice: false,
                    data: data.status,
                },
            },
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getEsimQRCode = async (
        data: GetESimQRCodeInput,
        cookie: ParsedCookie
    ): Promise<ESimQRCodeRes> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/users/${data.userId}/sims/${data.simId}/qr`,
            headers: cookie.header,
        });
        if (checkError(res))
            throw new Error(res.description ? res.description : res.message);
        return res;
    };
    updateUserRoaming = async (
        data: UpdateUserServiceInput,
        cookie: ParsedCookie
    ): Promise<OrgUserSimDto> => {
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
        return res;
    };
}
