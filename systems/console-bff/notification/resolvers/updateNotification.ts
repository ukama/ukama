/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { UpdateNotificationResDto } from "./types";

@Resolver()
export class UpdateNotificationResolver {
  @Mutation(() => UpdateNotificationResDto)
  async updateNotification(
    @Arg("id") id: string,
    @Arg("isRead") isRead: boolean,
    @Ctx() ctx: Context
  ): Promise<UpdateNotificationResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateNotification(id, isRead);
  }
}
