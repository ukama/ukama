/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { UserResDto } from "../resolver/types";
import { dtoToUserResDto } from "./mapper";

class UserApi extends RESTDataSource {
  baseURL = NUCLEUS_API_GW;

  getUser = async (userId: string): Promise<UserResDto> => {
    this.logger.info(`Request Url: ${this.baseURL}/${VERSION}/users/${userId}`);
    return this.get(`/${VERSION}/users/${userId}`, {}).then(res =>
      dtoToUserResDto(res)
    );
  };
}
export default UserApi;
