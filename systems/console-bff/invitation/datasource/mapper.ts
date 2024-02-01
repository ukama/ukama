/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { InvitationAPIResDto, InvitationDto } from "../resolver/types";

export const dtoToInvitationResDto = (
  res: InvitationAPIResDto
): InvitationDto => {
  return {
    email: res.email,
    expiresAt: res.expires_at,
    id: res.id,
    link: res.link,
    userId: res.user_id,
    name: res.name,
    org: res.org,
    role: res.role,
    status: res.status,
  };
};
