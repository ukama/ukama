/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddNetworkResolver } from "./addNetwork";
import { AddSiteToNetworkResolver } from "./addSiteToNetwork";
import { GetAllSitesResolver } from "./getAllSites";
import { GetNetworkResolver } from "./getNetwork";
import { GetNetworksResolver } from "./getNetworks";
import { GetSingleSiteResolver } from "./getSingleSite";

const resolvers: NonEmptyArray<any> = [
  AddNetworkResolver,
  AddSiteToNetworkResolver,
  GetNetworkResolver,
  GetNetworksResolver,
  GetSingleSiteResolver,
  GetAllSitesResolver,
];

export default resolvers;
