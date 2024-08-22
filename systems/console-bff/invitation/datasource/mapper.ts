/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  InvitationAPIDto,
  InvitationAPIResDto,
  InvitationDto,
  InvitationsAPIResDto,
  InvitationsResDto,
} from "../resolver/types";

export const dtoToInvitationsResDto = (
  res: InvitationsAPIResDto
): InvitationsResDto => {
  const invites = [];
  if (res && res.invitations.length > 0) {
    for (let i = 0; i < res.invitations.length; i++) {
      invites.push(dtoToInvitationDto(res.invitations[i]));
    }
  }
  return {
    invitations: invites,
  };
};

export const inviteResToInvitationDto = (
  res: InvitationAPIResDto
): InvitationDto => {
  return dtoToInvitationDto(res.invitation);
};

export const dtoToInvitationDto = (res: InvitationAPIDto): InvitationDto => {
  return {
    email: res.email,
    expireAt: res.expire_at,
    id: res.id,
    link: res.link,
    userId: res.user_id,
    name: res.name,
    role: res.role,
    status: res.status,
  };
};
