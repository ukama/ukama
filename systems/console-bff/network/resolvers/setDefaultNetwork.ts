/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { SetDefaultNetworkInputDto } from "./types";

@Resolver()
export class SetDefaultNetworkResolver {
  @Mutation(() => CBooleanResponse)
  async setDefaultNetwork(
    @Arg("data") data: SetDefaultNetworkInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources, baseURL } = ctx;
    const res = await dataSources.dataSource.setDefaultNetwork(
      baseURL,
      data.id
    );
    return {
      success: res.success,
    };
  }
}
