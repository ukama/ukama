/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { AddNodeToSiteInput } from "./types";

@Resolver()
export class AddNodeToSiteResolver {
  @Mutation(() => CBooleanResponse)
  async addNodeToSite(
    @Arg("data") data: AddNodeToSiteInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return dataSources.dataSource.addNodeToSite({
      nodeId: data.nodeId,
      networkId: data.networkId,
      siteId: data.siteId,
    });
  }
}
