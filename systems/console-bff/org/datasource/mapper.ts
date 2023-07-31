import { MEMBER_ROLES } from "../../constants";
import {
  MemberObj,
  OrgAPIResDto,
  OrgDto,
  OrgMemberAPIResDto,
  OrgMembersAPIResDto,
  OrgMembersResDto,
  OrgsAPIResDto,
  OrgsResDto,
} from "../resolver/types";

export const dtoToMemberResDto = (res: OrgMemberAPIResDto): MemberObj => {
  return {
    uuid: res.member.uuid,
    userId: res.member.user_id,
    orgId: res.member.org_id,
    role: res.member.role as MEMBER_ROLES,
    isDeactivated: res.member.is_deactivated,
    memberSince: res.member.member_since,
    user: undefined,
  };
};

export const dtoToOrgResDto = (res: OrgAPIResDto): OrgDto => {
  return {
    id: res.org.id,
    name: res.org.name,
    owner: res.org.owner,
    certificate: res.org.certificate,
    isDeactivated: res.org.is_deactivated,
    createdAt: res.org.created_at,
  };
};

export const dtoToOrgsResDto = (res: OrgsAPIResDto): OrgsResDto => {
  const orgs: OrgDto[] = [];
  res.orgs.forEach(org => {
    orgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
  });
  return {
    owner: res.owner,
    orgs,
  };
};

export const dtoToMembersResDto = (
  res: OrgMembersAPIResDto
): OrgMembersResDto => {
  const members: MemberObj[] = [];
  res.members.forEach(member => {
    members.push({
      uuid: member.uuid,
      userId: member.user_id,
      orgId: member.org_id,
      isDeactivated: member.is_deactivated,
      memberSince: member.member_since,
      role: member.role as MEMBER_ROLES,
      user: undefined,
    });
  });
  return {
    org: res.org,
    members,
  };
};
