import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";
import { CBooleanResponse, THeaders } from "../../common/types";
import {
  AddMemberInputDto,
  AddOrgInputDto,
  MemberObj,
  OrgDto,
  OrgMembersResDto,
  OrgsResDto,
  UpdateMemberInputDto,
} from "../resolver/types";
import {
  dtoToMemberResDto,
  dtoToMembersResDto,
  dtoToOrgResDto,
  dtoToOrgsResDto,
} from "./mapper";

const version = "/v1/orgs";

class OrgApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;
  getOrgMembers = async (headers: THeaders): Promise<OrgMembersResDto> => {
    return this.get(`/${headers.orgName}/members`).then(res =>
      dtoToMembersResDto(res)
    );
  };

  getOrgMember = async (headers: THeaders): Promise<MemberObj> => {
    return this.get(`/${headers.orgName}/members/${headers.userId}`).then(res =>
      dtoToMemberResDto(res)
    );
  };

  removeMember = async (headers: THeaders): Promise<CBooleanResponse> => {
    return this.delete(`/${headers.orgName}/members/${headers.userId}`).then(
      res => {
        return {
          success: true,
        };
      }
    );
  };

  getOrgs = async (headers: THeaders): Promise<OrgsResDto> => {
    return this.get(`/${headers.orgName}`).then(res => dtoToOrgsResDto(res));
  };

  getOrg = async (orgName: string, headers: THeaders): Promise<OrgDto> => {
    return this.get(`/${orgName}`, {}).then(res => dtoToOrgResDto(res));
  };

  addOrg = async (req: AddOrgInputDto, headers: THeaders): Promise<OrgDto> => {
    return this.post("", {
      body: req,
    }).then(res => dtoToOrgResDto(res));
  };

  addMember = async (
    data: AddMemberInputDto,
    headers: THeaders
  ): Promise<MemberObj> => {
    return this.post(`/${headers.orgName}/members`, {
      body: { user_uuid: headers.userId, role: data.role },
    }).then(res => dtoToMemberResDto(res));
  };

  updateMember = async (
    memberId: string,
    req: UpdateMemberInputDto,
    headers: THeaders
  ): Promise<CBooleanResponse> => {
    return this.post(`/${headers.orgName}/members/${memberId}`, {
      body: req,
    }).then(res => {
      return {
        success: true,
      };
    });
  };
}

export default OrgApi;
