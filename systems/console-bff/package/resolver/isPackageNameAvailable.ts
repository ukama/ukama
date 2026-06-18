/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { PackageNameAvailabilityResDto } from "./types";

@Resolver()
export class IsPackageNameAvailableResolver {
  @Query(() => PackageNameAvailabilityResDto)
  async isPackageNameAvailable(
    @Arg("name") name: string,
    @Ctx() ctx: AppContext
  ): Promise<PackageNameAvailabilityResDto> {
    const baseURL = await ctx.urls.url("package");
    return ctx.dataSources.package.isPackageNameAvailable(baseURL, name);
  }
}
