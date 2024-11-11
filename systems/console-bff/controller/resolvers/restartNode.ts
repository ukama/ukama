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
import { RestartNodeInputDto } from "./types";

@Resolver()
export class RestartNodeResolver {
  @Mutation(() => CBooleanResponse)
  async restartNode(
    @Arg("data") data: RestartNodeInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.restartNode(baseURL, data);
  }
}
