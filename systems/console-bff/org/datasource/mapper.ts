import { UserResDto } from "../../user/resolver/types";
import {
  MemberObj,
  OrgAPIResDto,
  OrgDto,
  OrgMemberAPIResDto,
  OrgMembersAPIResDto,
  OrgMembersResDto,
  OrgsAPIResDto,
  OrgsResDto,
  UserAPIResDto,
} from "../resolver/types";

export const dtoToMemberResDto = (res: OrgMemberAPIResDto): MemberObj => {
  return {
    uuid: res.member.uuid,
    userId: res.member.user_id,
    orgId: res.member.org_id,
    role: res.member.role,
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
  const ownerOfOrgs: OrgDto[] = [];
  const memberOfOrgs: OrgDto[] = [];
  res.owner_of.forEach(org => {
    ownerOfOrgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
  });
  res.member_of.forEach(org => {
    memberOfOrgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
  });

  return {
    user: res.user,
    ownerOf: ownerOfOrgs,
    memberOf: memberOfOrgs,
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
      role: member.role,
      user: undefined,
    });
  });
  return {
    org: res.org,
    members,
  };
};

export const dtoToUserResDto = (res: UserAPIResDto): UserResDto => {
  return {
    uuid: res.user.id,
    email: res.user.email,
    authId: res.user.auth_id,
    isDeactivated: res.user.is_deactivated,
    name: res.user.name,
    phone: res.user.phone,
    registeredSince: res.user.registered_since,
  };
};
