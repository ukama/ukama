/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { IdResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeletePackageResolver {
  @Mutation(() => IdResponse)
  async deletePackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: Context
  ): Promise<IdResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.deletePackage(packageId);
  }
}
