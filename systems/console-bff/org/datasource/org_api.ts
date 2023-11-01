/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { NUCLEUS_API_GW, VERSION } from "../../common/configs";
import { OrgDto, OrgsResDto } from "../resolver/types";
import { dtoToOrgResDto, dtoToOrgsResDto } from "./mapper";

class OrgApi extends RESTDataSource {
  baseURL = NUCLEUS_API_GW;

  getOrgs = async (userId: string): Promise<OrgsResDto> => {
    return this.get(`/${VERSION}/orgs`, {
      params: {
        user_uuid: userId,
      },
    })
      .then(res => dtoToOrgsResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getOrg = async (orgName: string): Promise<OrgDto> => {
    return this.get(`/${VERSION}/orgs/${orgName}`)
      .then(res => dtoToOrgResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default OrgApi;
