/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  MemberAPIResDto,
  MemberDto,
  MembersAPIResDto,
  MembersResDto,
  UserAPIResDto,
  UserResDto,
} from "../resolver/types";

export const dtoToMemberResDto = (res: MemberAPIResDto): MemberDto => {
  return {
    userId: res.member.user_id,
    role: res.member.role,
    memberId: res.member.member_id,
    isDeactivated: res.member.is_deactivated,
    memberSince: res.member.member_since,
  };
};

export const dtoToMembersResDto = (res: MembersAPIResDto): MembersResDto => {
  const members: MemberDto[] = [];
  res.members.forEach(member => {
    members.push({
      userId: member.user_id,
      memberId: member.member_id,
      isDeactivated: member.is_deactivated,
      memberSince: member.member_since,
      role: member.role,
    });
  });
  return {
    members,
  };
};
export const dtoToUserResDto = (res: UserAPIResDto): UserResDto => {
  return {
    uuid: res.user.id,
    email: res.user.email,
    isDeactivated: res.user.is_deactivated,
    name: res.user.name,
    authId: res.user.auth_id,
    phone: res.user.phone,
    registeredSince: res.user.registered_since,
  };
};
