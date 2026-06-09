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
  SubscriberSimDto,
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

    // Independent upstream calls — fetch in parallel.
    const [sims, subs]: [SubscriberSimsResDto, SubscribersResDto] =
      await Promise.all([
        dataSources.subscriber.getSimsByNetwork(baseURL, networkId),
        dataSources.subscriber.getSubscribersByNetwork(baseURL, networkId),
      ]);

    // Group sims by subscriber once (O(n + m)) instead of filtering the
    // whole sims array per subscriber (O(n * m)).
    const simsBySubscriber = new Map<string, SubscriberSimDto[]>();
    for (const sim of sims.sims) {
      const group = simsBySubscriber.get(sim.subscriberId);
      if (group) {
        group.push(sim);
      } else {
        simsBySubscriber.set(sim.subscriberId, [sim]);
      }
    }

    for (const sub of subs.subscribers) {
      sub.sim = simsBySubscriber.get(sub.uuid) ?? [];
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
        sim: sub.sim,
      });
    }
    return {
      subscribers: networkSub,
    };
  }
}
