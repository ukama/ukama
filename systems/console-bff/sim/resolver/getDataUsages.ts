/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { openStore } from "../../common/storage";
import { getBaseURL } from "../../common/utils";
import SubscriberApi from "../../subscriber/datasource/subscriber_api";
import { Context } from "../context";
import { SimDataUsages, SimUsagesInputDto, SimsUsageInputDto } from "./types";

@Resolver()
export class GetDataUsagesResolver {
  @Query(() => SimDataUsages)
  async getDataUsages(
    @Arg("data") data: SimUsagesInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDataUsages> {
    const { dataSources, baseURL, headers } = ctx;
    const store = openStore();
    const subUrl = await getBaseURL(
      SUB_GRAPHS.subscriber.name,
      headers.orgName,
      store
    );
    const subAPI = new SubscriberApi();
    const subs = await subAPI.getSubscribersByNetwork(
      subUrl.message,
      data.networkId
    );

    const simUsages: SimsUsageInputDto[] =
      subs.subscribers
        .map(s => {
          if (s && s.sim && s.sim[0]) {
            return {
              simId: s.sim[0].id,
              iccid: s.sim[0].iccid,
            };
          }
          return null;
        })
        .filter((item): item is SimsUsageInputDto => item !== null) ?? [];

    let to = data.to ?? 0;
    let from = data.from ?? 0;

    if (to === 0) {
      to = Math.round(Date.now() / 1000) - 60;
    }

    if (from === 0) {
      from = to - 240;
    }

    const usages = await Promise.all(
      simUsages.map(item =>
        dataSources.dataSource.getDataUsage(baseURL, {
          to,
          from,
          type: data.type,
          iccid: item.iccid,
          simId: item.simId,
        })
      )
    );

    return {
      usages,
    };
  }
}
