/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { IdResponse } from "../../common/types";
import type { AppContext } from "../../server/context";

@Resolver()
export class DeletePackageResolver {
  @Mutation(() => IdResponse)
  async deletePackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: AppContext
  ): Promise<IdResponse> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("package");
    return dataSources.package.deletePackage(baseURL, packageId);
  }
}
