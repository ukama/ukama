/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import crypto from "crypto";
import { Arg, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { NetworkStats } from "./types";

@Resolver()
export class GetNetworkStatsResolver {
  @Query(() => NetworkStats)
  async getNetworkStats(
    @Arg("networkId") networkId: string
  ): Promise<NetworkStats> {
    logger.info(`Getting network stats for network ${networkId}`);
    return {
      activeSubscriber: Math.floor(crypto.randomInt(1, 100)),
      averageThroughput: Math.floor(crypto.randomInt(1, 50)),
      averageSignalStrength: Math.floor(crypto.randomInt(1, 90)),
    };
  }
}
