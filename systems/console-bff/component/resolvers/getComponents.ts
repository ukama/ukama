/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { ComponentDto, ComponentsResDto } from "./types";

@Resolver()
export class GetComponentsResolver {
  @Query(() => ComponentDto)
  async getComponents(
    @Arg("userId") userId: string,
    @Ctx() ctx: Context
  ): Promise<ComponentsResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getComponents(userId);
  }
}
