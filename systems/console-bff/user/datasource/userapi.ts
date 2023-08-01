import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";
import { CBooleanResponse, THeaders } from "../../common/types";
import { SERVER } from "../../constants/endpoints";
import { getHeaders } from "../../utils";
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
} from "../resolver/types";
import {
  connectedUsersDtoToDto,
  dtoToUserResDto,
  dtoToWhoamiResDto,
} from "./mapper";

const version = "/v1/users";

class UserApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;

  getConnectedUsers = async (headers: THeaders): Promise<ConnectedUserDto> => {
    return this.get(`/${headers.orgId}/metrics/subscribersattached`).then(res =>
      connectedUsersDtoToDto(res.data.result)
    );
  };

  updateUser = async (
    userId: string,
    req: UpdateUserInputDto
  ): Promise<UserResDto> => {
    return this.put(`/${userId}`, {
      body: { name: req.name, email: req.email, phone: req.phone },
    }).then(res => dtoToUserResDto(res));
  };

  deactivateUser = async (userId: string): Promise<UserResDto> => {
    return this.patch(`/${userId}`, {
      body: { isDeactivated: true },
    }).then(res => dtoToUserResDto(res));
  };

  getUser = async (userId: string): Promise<UserResDto> => {
    return this.get(`/${userId}`, {}).then(res => dtoToUserResDto(res));
  };

  whoami = async (): Promise<WhoamiDto> => {
    return this.get(`/whoami}`, {}).then(res => dtoToWhoamiResDto(res));
  };
  getAccountDetails = async (
    headers: THeaders
  ): Promise<GetAccountDetailsDto> => {
    return this.get(`/${headers.orgId}`).then(res => {
      return {
        email: res?.traits.email,
        isFirstVisit:
          res.traits.firstVisit == undefined ? true : res.traits.firstVisit,
      };
    });
  };
  addUser = async (req: UserInputDto): Promise<UserResDto> => {
    return this.post(``, {
      body: { ...req },
    }).then(res => dtoToUserResDto(res));
  };
  updateFirstVisit = async (
    req: UserFistVisitInputDto,
    headers: THeaders
  ): Promise<UserFistVisitResDto> => {
    let user = await this.get(`${SERVER.GET_IDENTITY}/${headers.orgId}`);
    return this.put(`${SERVER.GET_IDENTITY}/${headers.orgId}`, {
      headers: getHeaders(headers),
      body: {
        schema_id: "default",
        state: "active",
        traits: {
          email: user?.traits?.email,
          name: user?.traits?.name,
          ...req,
        },
      },
    }).then(res => {
      return {
        firstVisit: res?.traits?.firstVisit,
      };
    });
  };

  deleteUser = async (userId: string): Promise<CBooleanResponse> => {
    return this.delete(`/${userId}`).then(res => {
      return {
        success: true,
      };
    });
  };
  updateUserStatus = async (
    data: UpdateUserServiceInput,
    headers: THeaders
  ): Promise<OrgUserSimDto> => {
    return this.put(
      `/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,
      {
        body: {
          ukama: {
            sms: false,
            voice: false,
            data: data.status,
          },
        },
      }
    ).then(res => res);
  };
  getEsimQRCode = async (
    data: GetESimQRCodeInput,
    headers: THeaders
  ): Promise<ESimQRCodeRes> => {
    return this.get(
      `/${headers.orgId}/users/${data.userId}/sims/${data.simId}/qr`
    ).then(res => res);
  };
  updateUserRoaming = async (
    data: UpdateUserServiceInput,
    headers: THeaders
  ): Promise<OrgUserSimDto> => {
    return this.put(
      `/${headers.orgId}/users/${data.userId}/sims/${data.simId}/services`,
      {
        body: {
          carrier: {
            sms: false,
            voice: false,
            data: data.status,
          },
        },
      }
    ).then(res => res);
  };
}
export default UserApi;
