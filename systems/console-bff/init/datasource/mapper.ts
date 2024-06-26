/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { InitSystemAPIResDto } from "../resolver/types";

export const dtoToSystenResDto = (
  res: InitSystemAPIResDto
): InitSystemAPIResDto => {
  return {
    certificate: res.certificate,
    health: res.health,
    ip: res.ip,
    orgName: res.orgName,
    port: res.port,
    systemId: res.systemId,
    systemName: res.systemName,
  };
};
