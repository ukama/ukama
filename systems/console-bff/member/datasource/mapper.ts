import {
  MemberAPIResDto,
  MemberDto,
  MembersAPIResDto,
  MembersResDto,
} from "../resolver/types";

export const dtoToMemberResDto = (res: MemberAPIResDto): MemberDto => {
  return {
    userId: res.member.user_id,
    orgId: res.member.org_id,
    role: res.member.role,
    isDeactivated: res.member.is_deactivated,
    memberSince: res.member.member_since,
  };
};

export const dtoToMembersResDto = (res: MembersAPIResDto): MembersResDto => {
  const members: MemberDto[] = [];
  res.members.forEach(member => {
    members.push({
      userId: member.user_id,
      orgId: member.org_id,
      isDeactivated: member.is_deactivated,
      memberSince: member.member_since,
      role: member.role,
    });
  });
  return {
    members,
  };
};
