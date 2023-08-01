import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
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

class OrgApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW;
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

  getOrgs = async (userId: string): Promise<OrgsResDto> => {
    return this.get(`/${VERSION}/orgs`, {
      params: {
        user_uuid: userId,
      },
    }).then(res => dtoToOrgsResDto(res));
  };

  getOrg = async (orgName: string): Promise<OrgDto> => {
    return this.get(`/${VERSION}/orgs/${orgName}`).then(res =>
      dtoToOrgResDto(res)
    );
  };

  addOrg = async (req: AddOrgInputDto): Promise<OrgDto> => {
    return this.post(`/${VERSION}/orgs`, {
      body: {
        ...req,
      },
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
