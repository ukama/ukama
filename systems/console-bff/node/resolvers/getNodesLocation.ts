/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import crypto from "crypto";
import { Arg, Query, Resolver } from "type-graphql";

import { NODE_STATUS } from "../../common/enums";
import { NodeLocation, NodesInput, NodesLocation } from "./types";

const getRandomNodeState = () => {
  const nodeStates = [
    NODE_STATUS.ACTIVE,
    NODE_STATUS.CONFIGURED,
    NODE_STATUS.FAULTY,
    NODE_STATUS.MAINTENANCE,
    NODE_STATUS.ONBOARDED,
    NODE_STATUS.UNDEFINED,
  ];
  return nodeStates[Math.floor(crypto.randomInt(1, nodeStates.length))];
};

@Resolver()
export class GetNodesLocationResolver {
  @Query(() => NodesLocation)
  async getNodesLocation(@Arg("data") data: NodesInput) {
    const locations = [
      "37.75210734489673 , -122.49572485891302",
      "37.75046506608679 , -122.50583612245225",
      "37.75612411999306 , -122.50728172021206",
      "37.75726922992242 , -122.50251213426036",
      "37.75243837156798 , -122.50442682534892",
      "37.7519789637837 , -122.50402178717827",
      "37.750903349294404 , -122.50241414813944",
      "37.75602225181627 , -122.50060819272488",
      "37.757921529607835 , -122.50036152209083",
      "37.75628955086523 , -122.50694962686946",
      "37.7573215112949 , -122.50224043772997",
      "37.75074935869865 , -122.50127064328588",
      "37.7528943256246 , -122.501056716164",
      "37.754832309416386 , -122.50268274843049",
      "37.757352142065265 , -122.50390638094865",
      "37.75055972208169 , -122.50381787073599",
      "37.753482040181844 , -122.49795018201644",
      "37.7578160107123 , -122.50013574926646",
      "37.749592580038346 , -122.50730545397994",
      "37.7514871501036 , -122.49702703770673",
    ];
    const nodes: NodeLocation[] = [];
    for (let i = 0; i < 8; i++) {
      const randomLocation =
        locations[Math.floor(crypto.randomInt(1, locations.length))];
      nodes.push({
        id: "node" + i,
        lat: randomLocation.split(",")[0].trim(),
        lng: randomLocation.split(",")[1].trim(),
        state: getRandomNodeState(),
      });
    }

    return {
      networkId: data.networkId,
      nodes:
        data.nodeFilterState === NODE_STATUS.UNDEFINED
          ? nodes
          : nodes.filter(
              (node: NodeLocation) => node.state === data.nodeFilterState
            ),
    };
  }
}
