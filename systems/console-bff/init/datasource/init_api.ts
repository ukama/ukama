/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { whoami } from "../../common/auth/authCalls";
import { INIT_API_GW, NUCLEUS_API_GW, VERSION } from "../../common/configs";
import {
  dtoToUserResDto,
  dtoToWhoamiResDto,
} from "../../user/datasource/mapper";
import { WhoamiDto } from "../../user/resolver/types";
import { InitSystemAPIResDto } from "../resolver/types";
import { dtoToSystenResDto } from "./mapper";

class InitAPI extends RESTDataSource {
  getSystems = async (
    orgName: string,
    systemName: string
  ): Promise<InitSystemAPIResDto> => {
    this.baseURL = INIT_API_GW;
    return this.get(
      `/${VERSION}/orgs/${orgName}/systems/${systemName}`,
      {}
    ).then(res => dtoToSystenResDto(res));
  };

  validateSession = async (cookies: string): Promise<WhoamiDto> => {
    const whoamiRes = await whoami(cookies);
    let aId = "";
    if (whoamiRes?.data) {
      aId = whoamiRes.data.identity.id;
    }
    this.baseURL = NUCLEUS_API_GW;
    const userRes = await this.get(`/${VERSION}/users/auth/${aId}`).then(res =>
      dtoToUserResDto(res)
    );

    return this.get(`/${VERSION}/users/whoami/${userRes.uuid}`).then(res =>
      dtoToWhoamiResDto(res)
    );
  };
}
export default InitAPI;
