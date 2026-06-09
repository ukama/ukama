/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import { GetSubscribersByNetworkResolver } from "../../subscriber/resolver/getSubscribersByNetwork";

const makeSub = (uuid: string) => ({
  uuid,
  name: `sub-${uuid}`,
  dob: "",
  phone: "",
  email: "",
  gender: "",
  address: "",
  idSerial: "",
  networkId: "net-1",
  proofOfIdentification: "",
});

const makeSim = (id: string, subscriberId: string) => ({
  id,
  subscriberId,
});

describe("getSubscribersByNetwork sim join", () => {
  it("attaches each subscriber's sims and parallel-fetches upstreams", async () => {
    const getSimsByNetwork = jest.fn().mockResolvedValue({
      sims: [
        makeSim("sim-1", "sub-a"),
        makeSim("sim-2", "sub-b"),
        makeSim("sim-3", "sub-a"),
      ],
    });
    const getSubscribersByNetwork = jest.fn().mockResolvedValue({
      subscribers: [makeSub("sub-a"), makeSub("sub-b"), makeSub("sub-c")],
    });
    const ctx = {
      baseURL: "http://subscriber.test",
      dataSources: {
        subscriber: { getSimsByNetwork, getSubscribersByNetwork },
      },
    };

    const resolver = new GetSubscribersByNetworkResolver();
    const result = await resolver.getSubscribersByNetwork(
      "net-1",
      ctx as never
    );

    const byUuid = Object.fromEntries(
      result.subscribers.map(s => [s.uuid, s.sim])
    );
    expect(byUuid["sub-a"]?.map(sim => sim.id)).toEqual(["sim-1", "sim-3"]);
    expect(byUuid["sub-b"]?.map(sim => sim.id)).toEqual(["sim-2"]);
    expect(byUuid["sub-c"]).toEqual([]);
    expect(getSimsByNetwork).toHaveBeenCalledWith(
      "http://subscriber.test",
      "net-1"
    );
    expect(getSubscribersByNetwork).toHaveBeenCalledWith(
      "http://subscriber.test",
      "net-1"
    );
  });
});
