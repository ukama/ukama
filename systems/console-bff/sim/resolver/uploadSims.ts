/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { UploadSimsInputDto, UploadSimsResDto } from "./types";

@Resolver()
export class UploadSimsResolver {
  @Mutation(() => UploadSimsResDto)
  async uploadSims(
    @Arg("data") data: UploadSimsInputDto,
    @Ctx() ctx: Context
  ): Promise<UploadSimsResDto> {
    const { dataSources, baseURL } = ctx;
    return await dataSources.dataSource.uploadSims(baseURL, data);
  }
}
