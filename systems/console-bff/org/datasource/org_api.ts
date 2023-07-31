import { RESTDataSource } from "@apollo/datasource-rest";

import { BoolResponse, THeaders } from "../../common/types";
import { SERVER } from "../../constants/endpoints";
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

class OrgApi extends RESTDataSource {
  getOrgMembers = async (headers: THeaders): Promise<OrgMembersResDto> => {
    return this.get(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`
    ).then(res => dtoToMembersResDto(res));
  };

  getOrgMember = async (headers: THeaders): Promise<MemberObj> => {
    return this.get(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`
    ).then(res => dtoToMemberResDto(res));
  };

  removeMember = async (headers: THeaders): Promise<BoolResponse> => {
    return this.delete(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`
    ).then(res => {
      return {
        success: true,
      };
    });
  };

  getOrgs = async (headers: THeaders): Promise<OrgsResDto> => {
    return this.get(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}`,
      {}
    ).then(res => dtoToOrgsResDto(res));
  };

  getOrg = async (orgName: string, headers: THeaders): Promise<OrgDto> => {
    return this.get(`${SERVER.REGISTRY_ORGS_API_URL}/${orgName}`, {}).then(
      res => dtoToOrgResDto(res)
    );
  };

  addOrg = async (req: AddOrgInputDto, headers: THeaders): Promise<OrgDto> => {
    return this.post(SERVER.REGISTRY_ORGS_API_URL, {
      body: req,
    }).then(res => dtoToOrgResDto(res));
  };

  addMember = async (
    data: AddMemberInputDto,
    headers: THeaders
  ): Promise<MemberObj> => {
    return this.post(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`,
      {
        body: { user_uuid: headers.userId, role: data.role },
      }
    ).then(res => dtoToMemberResDto(res));
  };

  updateMember = async (
    memberId: string,
    req: UpdateMemberInputDto,
    headers: THeaders
  ): Promise<BoolResponse> => {
    return this.post(
      `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${memberId}`,
      {
        body: req,
      }
    ).then(res => {
      return {
        success: true,
      };
    });
  };
}

export default OrgApi;
