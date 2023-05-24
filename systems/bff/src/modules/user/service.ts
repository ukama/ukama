/* eslint-disable prettier/prettier */
import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { HTTP404Error, Messages, checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { IUserService } from "./interface";
import UserMapper from "./mapper";
import {
    ConnectedUserDto,
    ESimQRCodeRes,
    GetAccountDetailsDto,
    GetESimQRCodeInput,
    OrgUserSimDto,
    UpdateUserInputDto,
    UpdateUserServiceInput,
    UserFistVisitInputDto,
    UserFistVisitResDto,
    UserInputDto,
    UserResDto,
    WhoamiDto,
} from "./types";

@Service()
export class UserService implements IUserService {
    getConnectedUsers = async (
        headers: THeaders
    ): Promise<ConnectedUserDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${headers.orgId}/metrics/subscribersattached`,
            headers: getHeaders(headers),
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
        req: UpdateUserInputDto,
        headers: THeaders
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.REGISTRY_USERS_API_URL}/${userId}`,
            headers: getHeaders(headers),
            body: { name: req.name, email: req.email, phone: req.phone },
        });
        if (checkError(res)) throw new Error(res.message);
        return UserMapper.dtoToUserResDto(res);
    };
    deactivateUser = async (
        userId: string,
        headers: THeaders
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.REGISTRY_USERS_API_URL}/${userId}`,
            headers: getHeaders(headers),
            body: { isDeactivated: true },
        });
        if (checkError(res)) throw new Error(res.description);
        return UserMapper.dtoToUserResDto(res);
    };
    getUser = async (
        userId: string,
        headers: THeaders
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_USERS_API_URL}/${userId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return UserMapper.dtoToUserResDto(res);
    };
    whoami = async (headers: THeaders): Promise<WhoamiDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.AUTH_API_URL}/whoami`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return UserMapper.dtoToWhoamiResDto(res);
    };
    getAccountDetails = async (
        headers: THeaders
    ): Promise<GetAccountDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.GET_IDENTITY}/${headers.orgId}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(res.Messages);
        return {
            email: res?.traits.email,
            isFirstVisit:
                res.traits.firstVisit == undefined
                    ? true
                    : res.traits.firstVisit,
        };
    };
    addUser = async (
        req: UserInputDto,
        headers: THeaders
    ): Promise<UserResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_USERS_API_URL}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return UserMapper.dtoToUserResDto(res);
    };
    updateFirstVisit = async (
        req: UserFistVisitInputDto,
        headers: THeaders
    ): Promise<UserFistVisitResDto> => {
        const getUser = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.GET_IDENTITY}/${headers.orgId}`,
            headers: getHeaders(headers),
        });
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.GET_IDENTITY}/${headers.orgId}`,
            body: {
                schema_id: "default",
                state: "active",
                traits: {
                    email: getUser?.traits?.email,
                    name: getUser?.traits?.name,
                    ...req,
                },
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.description || res.message);
        return {
            firstVisit: res?.traits?.firstVisit,
        };
    };
    deleteUser = async (
        userId: string,
        headers: THeaders
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.REGISTRY_USERS_API_URL}/${userId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return {
            success: true,
        };
    };
    updateUserStatus = async (
        data: UpdateUserServiceInput,
        headers: THeaders
    ): Promise<OrgUserSimDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,
            headers: getHeaders(headers),
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
        headers: THeaders
    ): Promise<ESimQRCodeRes> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/qr`,
            headers: getHeaders(headers),
        });
        if (checkError(res))
            throw new Error(res.description ? res.description : res.message);
        return res;
    };
    updateUserRoaming = async (
        data: UpdateUserServiceInput,
        headers: THeaders
    ): Promise<OrgUserSimDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,
            headers: getHeaders(headers),
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
