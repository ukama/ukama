/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  OrgAPIResDto,
  OrgDto,
  OrgsAPIResDto,
  OrgsResDto,
} from "../resolver/types";

export const dtoToOrgResDto = (res: OrgAPIResDto): OrgDto => {
  return {
    id: res.org.id,
    name: res.org.name,
    owner: res.org.owner,
    country: res.org.country,
    currency: res.org.currency,
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
      country: org.country,
      currency: org.currency,
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
      country: org.country,
      currency: org.currency,
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
