import { RESTDataSource } from "@apollo/datasource-rest";
/* eslint-disable prettier/prettier */
import { BoolResponse, THeaders } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
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
} from "../types";

export class UserApi extends RESTDataSource {
    getConnectedUsers = async (
        headers: THeaders
    ): Promise<ConnectedUserDto> => {
        return this.get(`${SERVER.ORG}/${headers.orgId}/metrics/subscribersattached`,{headers: getHeaders(headers)}).then(res => 
            UserMapper.connectedUsersDtoToDto(
                res.data.result
            )
            );
    };

    updateUser = async (
        userId: string,
        req: UpdateUserInputDto,
        headers: THeaders
    ): Promise<UserResDto> => {
        return this.put(`${SERVER.REGISTRY_USERS_API_URL}/${userId}`, {
            headers: getHeaders(headers),
            body: { name: req.name, email: req.email, phone: req.phone },
          }).then(res => UserMapper.dtoToUserResDto(res));
    };

    deactivateUser = async (
        userId: string,
        headers: THeaders
    ): Promise<UserResDto> => {
        return this.patch(`${SERVER.REGISTRY_USERS_API_URL}/${userId}`, {
            headers: getHeaders(headers),
            body: { isDeactivated: true },
          }).then(res => UserMapper.dtoToUserResDto(res));
    };

    getUser = async (
        userId: string,
        headers: THeaders
    ): Promise<UserResDto> => {
        return this.get(`${SERVER.REGISTRY_USERS_API_URL}/${userId}`,{headers: getHeaders(headers)}).then(res => 
            UserMapper.dtoToUserResDto(res)
            );
    };

    whoami = async (headers: THeaders): Promise<WhoamiDto> => {
        return this.get(`${SERVER.AUTH_API_URL}/whoami}`,{headers: getHeaders(headers)}).then(res => 
            UserMapper.dtoToWhoamiResDto(res)
            );
        
    };
    getAccountDetails = async (
        headers: THeaders
    ): Promise<GetAccountDetailsDto> => {
        return this.get(`${SERVER.GET_IDENTITY}/${headers.orgId}`,{headers: getHeaders(headers)}).then(res =>{ 
            return {
                email: res?.traits.email,
                isFirstVisit:
                    res.traits.firstVisit == undefined
                        ? true
                        : res.traits.firstVisit,
            };
        });
    };
    addUser = async (
        req: UserInputDto,
        headers: THeaders
    ): Promise<UserResDto> => {
        return this.post(`${SERVER.REGISTRY_USERS_API_URL}`, {
            headers: getHeaders(headers),
            body: { ...req },
          }).then(res => UserMapper.dtoToUserResDto(res));
    };
    updateFirstVisit = async (
        req: UserFistVisitInputDto,
        headers: THeaders
    ): Promise<UserFistVisitResDto> => {
        let user = await this.get(`${SERVER.GET_IDENTITY}/${headers.orgId}`,{headers: getHeaders(headers)});
        return this.put(`${SERVER.GET_IDENTITY}/${headers.orgId}`,{headers: getHeaders(headers),
            body: {
                schema_id: "default",
                state: "active",
                traits: {
                    email: user?.traits?.email,
                    name: user?.traits?.name,
                    ...req,
                },
            },}).then(res =>{
            return {
                firstVisit: res?.traits?.firstVisit,
            }
        });
    };

    deleteUser = async (
        userId: string,
        headers: THeaders
    ): Promise<BoolResponse> => {
        return this.delete(`${SERVER.REGISTRY_USERS_API_URL}/${userId}`,{headers: getHeaders(headers)}).then((res)=>{
            return {
                success: true,
            };
        });

    };
    updateUserStatus = async (
        data: UpdateUserServiceInput,
        headers: THeaders
    ): Promise<OrgUserSimDto> => {
        return this.put(`${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,{headers: getHeaders(headers),
            body: {
                ukama: {
                    sms: false,
                    voice: false,
                    data: data.status,
                },
            }}).then(res => res);
    };
    getEsimQRCode = async (
        data: GetESimQRCodeInput,
        headers: THeaders
    ): Promise<ESimQRCodeRes> => {
        return this.get(`${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/qr`,{headers: getHeaders(headers)}).then(res => res);
        
    };
    updateUserRoaming = async (
        data: UpdateUserServiceInput,
        headers: THeaders
    ): Promise<OrgUserSimDto> => {
        return this.put(`${SERVER.ORG}/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,{headers: getHeaders(headers),
            body: {
                carrier: {
                    sms: false,
                    voice: false,
                    data: data.status,
                },
            }}).then(res => res);
    };
}
