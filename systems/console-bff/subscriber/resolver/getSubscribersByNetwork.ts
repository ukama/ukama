/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import {
  SubscriberDto,
  SubscriberSimsResDto,
  SubscribersResDto,
} from "./types";

@Resolver()
export class GetSubscribersByNetworkResolver {
  @Query(() => SubscribersResDto)
  async getSubscribersByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SubscribersResDto> {
    const { dataSources, baseURL } = ctx;
    const networkSub: SubscriberDto[] = [];
    const sims: SubscriberSimsResDto =
      await dataSources.dataSource.getSimsByNetwork(baseURL, networkId);

    const subs: SubscribersResDto =
      await dataSources.dataSource.getSubscribersByNetwork(baseURL, networkId);

    for (const sub of subs.subscribers) {
      sub.sim = sims.sims.filter(sim => sim.subscriberId === sub.uuid);
      networkSub.push({
        dob: sub.dob,
        uuid: sub.uuid,
        name: sub.name,
        phone: sub.phone,
        email: sub.email,
        gender: sub.gender,
        address: sub.address,
        idSerial: sub.idSerial,
        networkId: sub.networkId,
        proofOfIdentification: sub.proofOfIdentification,
        subscriberStatus: sub.subscriberStatus,
        sim: sub.sim,
      });
    }
    return {
      subscribers: networkSub,
    };
  }
}
