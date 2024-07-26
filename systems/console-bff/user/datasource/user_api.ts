/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { UserResDto, WhoamiDto } from "../resolver/types";
import { dtoToUserResDto, dtoToWhoamiResDto } from "./mapper";

class UserApi extends RESTDataSource {
  baseURL = NUCLEUS_API_GW;

  getUser = async (userId: string): Promise<UserResDto> => {
    this.logger.info(`GetUser GET: ${this.baseURL}/${VERSION}/users/${userId}`);
    return this.get(`/${VERSION}/users/${userId}`, {}).then(res =>
      dtoToUserResDto(res)
    );
  };

  whoami = async (userId: string): Promise<WhoamiDto> => {
    this.logger.info(
      `Whoami GET: ${this.baseURL}/${VERSION}/users/whoami/${userId}`
    );
    return this.get(`/${VERSION}/users/whoami/${userId}`).then(res =>
      dtoToWhoamiResDto(res)
    );
  };

  auth = async (authId: string): Promise<UserResDto> => {
    this.logger.info(
      `Auth GET: ${this.baseURL}/${VERSION}/users/auth/${authId}`
    );
    return this.get(`/${VERSION}/users/auth/${authId}`).then(res =>
      dtoToUserResDto(res)
    );
  };
}
export default UserApi;
