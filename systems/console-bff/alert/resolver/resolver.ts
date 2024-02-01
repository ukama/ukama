/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, PubSub, PubSubEngine, Query, Resolver } from "type-graphql";

import { PaginationDto } from "../../common/types";
import { Context } from "../context";
import { AlertsResponse } from "./types";

@Resolver()
export class GetAlertsResolver {
  @Query(() => AlertsResponse)
  async getAlerts(
    @Arg("data") data: PaginationDto,
    @PubSub() pubsub: PubSubEngine,
    @Ctx() context: Context
  ): Promise<AlertsResponse> {
    const { dataSources } = context;
    const alerts = dataSources.dataSource.getAlerts(data);
    pubsub.publish("getAlerts", alerts);
    return alerts;
  }
}
