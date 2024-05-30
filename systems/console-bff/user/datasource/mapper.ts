/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { OrgDto } from "../../org/resolver/types";
import {
  UserAPIResDto,
  UserResDto,
  WhoamiAPIDto,
  WhoamiDto,
} from "../resolver/types";

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
export const dtoToWhoamiResDto = (res: WhoamiAPIDto): WhoamiDto => {
  const ownerOfOrgs: OrgDto[] = [];
  const memberOfOrgs: OrgDto[] = [];
  res.ownerOf.forEach(org => {
    ownerOfOrgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
  });
  res.memberOf.forEach(org => {
    // if (org.name !== COMMUNITY_ORG_NAME) {
    memberOfOrgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
    // }
  });

  return {
    user: dtoToUserResDto(res),
    ownerOf: ownerOfOrgs,
    memberOf: memberOfOrgs,
  };
};
