/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { SIM_TYPES } from "../../common/enums";
import {
  SimAPIResDto,
  SimDetailsDto,
  SimDto,
  SimsAPIResDto,
  SimsResDto,
} from "../resolver/types";

export const dtoToSimResDto = (res: SimAPIResDto): SimDto => {
  return {
    activationCode: res.sim.activation_code,
    createdAt: res.sim.created_at,
    iccid: res.sim.iccid,
    id: res.sim.id,
    isAllocated: res.sim.is_allocated,
    isPhysical: res.sim.is_physical,
    msisdn: res.sim.msisdn,
    qrCode: res.sim.qr_code,
    simType: res.sim.sim_type as SIM_TYPES,
    smapAddress: res.sim.sm_ap_address,
  };
};

export const dtoToSimDetailsDto = (response: any): SimDetailsDto => {
  const {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn,
    lastActivatedOn,
    activationsCount,
    deactivationsCount,
    allocatedAt,
  } = response;

  return {
    id,
    subscriberId,
    networkId,
    orgId,
    Package,
    iccid,
    msisdn,
    imsi,
    type,
    status,
    isPhysical,
    firstActivatedOn: firstActivatedOn?.toDate(),
    lastActivatedOn: lastActivatedOn?.toDate(),
    activationsCount,
    deactivationsCount,
    allocatedAt: allocatedAt?.toDate(),
  };
};

export const dtoToSimsDto = (res: SimsAPIResDto): SimsResDto => {
  const sims: SimDto[] = [];
  for (const sim of res.sims) {
    sims.push(dtoToSimResDto({ sim: sim }));
  }
  return {
    sim: sims,
  };
};
