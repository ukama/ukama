import { RESTDataSource } from "@apollo/datasource-rest";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  AddMemberInputDto,
  AddOrgInputDto,
  MemberInputDto,
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
  baseURL = NUCLEUS_API_GW;
  getOrgMembers = async (orgName: string): Promise<OrgMembersResDto> => {
    return this.get(`/${VERSION}/orgs/${orgName}/members`).then(res =>
      dtoToMembersResDto(res)
    );
  };

  getOrgMember = async (data: MemberInputDto): Promise<MemberObj> => {
    return this.get(
      `/${VERSION}/orgs/${data.orgName}/members/${data.memberId}`
    ).then(res => dtoToMemberResDto(res));
  };

  removeMember = async (data: MemberInputDto): Promise<CBooleanResponse> => {
    return this.delete(
      `/${VERSION}/orgs/${data.orgName}/members/${data.memberId}`
    ).then(res => {
      return {
        success: true,
      };
    });
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
    orgName: string
  ): Promise<MemberObj> => {
    return this.post(`/${VERSION}/orgs/${orgName}/members`, {
      body: { user_uuid: data.userId, role: data.role },
    }).then(res => dtoToMemberResDto(res));
  };

  updateMember = async (
    memberId: string,
    req: UpdateMemberInputDto
  ): Promise<CBooleanResponse> => {
    return this.post(`/${VERSION}/orgs/${req.orgName}/members/${memberId}`, {
      body: {
        isDeactivated: req.isDeactivated,
        role: req.role,
      },
    }).then(res => {
      return {
        success: true,
      };
    });
  };
}

export default OrgApi;
