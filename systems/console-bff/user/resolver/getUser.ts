/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { UserResDto } from "./types";

@Resolver()
export class GetUserResolver {
  @Query(() => UserResDto)
  async getUser(
    @Arg("userId") userId: string,
    @Ctx() ctx: AppContext
  ): Promise<UserResDto | null> {
    const { dataSources } = ctx;
    return dataSources.user.getUser(userId);
  }
}
