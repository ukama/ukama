/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { InitSystemAPIResDto, OrgsNameRes } from "../resolver/types";

export const dtoToSystenResDto = (
  res: InitSystemAPIResDto
): InitSystemAPIResDto => {
  return {
    certificate: res.certificate,
    apiGwHealth: res.apiGwHealth,
    apiGwIp: res.apiGwIp,
    apiGwUrl: res.apiGwUrl,
    apiGwPort: res.apiGwPort,
    orgName: res.orgName,
    systemId: res.systemId,
    systemName: res.systemName,
  };
};

export const dtoToOrgsNameResDto = (res: OrgsNameRes): OrgsNameRes => {
  return {
    orgs: res.orgs,
  };
};
