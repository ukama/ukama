/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Query, Resolver } from "type-graphql";

import { AppChangeLogs, NodeAppsChangeLogInput } from "./types";

@Resolver()
export class GetAppsChangeLogResolver {
  @Query(() => AppChangeLogs)
  async getAppsChangeLog(@Arg("data") data: NodeAppsChangeLogInput) {
    const logs: AppChangeLogs = { logs: [], type: data.type };
    for (let i = 10; i > 0; i--) {
      logs.logs.push({
        date: Math.floor(new Date().getTime() / 1000 - 1000 * (10 - i)),
        version: `${i}.0.0`,
      });
    }
    return logs;
  }
}
