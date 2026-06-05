/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { COMMUNITY_ORG_NAME } from "../../common/configs";
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
  // Nucleus responds in snake_case (owner_of/member_of), consistent with its
  // other fields; tolerate both shapes so a missing list never throws.
  const raw = res as unknown as {
    ownerOf?: WhoamiAPIDto["ownerOf"];
    owner_of?: WhoamiAPIDto["ownerOf"];
    memberOf?: WhoamiAPIDto["memberOf"];
    member_of?: WhoamiAPIDto["memberOf"];
  };
  const ownerOf = raw.ownerOf ?? raw.owner_of ?? [];
  const memberOf = raw.memberOf ?? raw.member_of ?? [];
  ownerOf.forEach(org => {
    ownerOfOrgs.push({
      id: org.id,
      name: org.name,
      owner: org.owner,
      country: org.country,
      currency: org.currency,
      certificate: org.certificate,
      isDeactivated: org.is_deactivated,
      createdAt: org.created_at,
    });
  });
  memberOf.forEach(org => {
    if (org.name !== COMMUNITY_ORG_NAME) {
      memberOfOrgs.push({
        id: org.id,
        name: org.name,
        owner: org.owner,
        country: org.country,
        currency: org.currency,
        certificate: org.certificate,
        isDeactivated: org.is_deactivated,
        createdAt: org.created_at,
      });
    }
  });

  return {
    user: dtoToUserResDto(res),
    ownerOf: ownerOfOrgs,
    memberOf: memberOfOrgs,
  };
};
