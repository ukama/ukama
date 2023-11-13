/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Query, Resolver } from "type-graphql";

import { NodeApps, NodeAppsChangeLogInput } from "./types";

@Resolver()
export class GetNodeAppsResolver {
  @Query(() => NodeApps)
  async getNodeApps(@Arg("data") data: NodeAppsChangeLogInput) {
    const apps: NodeApps = { apps: [], type: data.type };
    const type = {
      tnode: 10,
      anode: 7,
      hnode: 5,
    };
    for (let i = type[data.type]; i > 0; i--) {
      apps.apps.push({
        cpu: `${Math.floor(Math.random() * 100)}%`,
        memory: `${Math.floor(Math.random() * 100)}MB`,
        name: `App ${i}`,
        date: Math.floor(new Date().getTime() / 1000 - 1000 * (10 - i)),
        version: `${i}.0.0`,
        notes:
          "Minim tempor tempor aliqua excepteur qui sit est Lorem tempor irure ullamco. Ex adipisicing do officia exercitation non anim elit excepteur nisi tempor veniam. Eu in ex aliquip dolore voluptate reprehenderit deserunt laborum culpa irure reprehenderit sit.",
      });
    }
    return apps;
  }
}
