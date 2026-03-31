/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { StringResponse, UpdateSoftwareInputDto } from "./types";

@Resolver()
export class UpdateSoftwareResolver {
  @Mutation(() => StringResponse)
  async updateSoftware(
    @Arg("data") data: UpdateSoftwareInputDto,
    @Ctx() ctx: Context
  ): Promise<StringResponse> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.updateSoftware(baseURL, data);
  }
}
