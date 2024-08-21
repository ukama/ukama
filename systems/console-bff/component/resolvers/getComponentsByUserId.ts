/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { ComponentTypeInputDto, ComponentsResDto } from "./types";

@Resolver()
export class GetComponentsByUserResolver {
  @Query(() => ComponentsResDto)
  async getComponentsByUserId(
    @Arg("data") data: ComponentTypeInputDto,
    @Ctx() ctx: Context
  ): Promise<ComponentsResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getComponentsByUserId(headers, data.category);
  }
}
