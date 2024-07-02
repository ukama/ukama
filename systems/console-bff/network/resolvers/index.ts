/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddNetworkResolver } from "./addNetwork";
import { GetNetworkResolver } from "./getNetwork";
import { GetNetworksResolver } from "./getNetworks";
<<<<<<< HEAD
=======
import { GetSiteResolver } from "./getSite";
import { GetSitesResolver } from "./getSites";
import { SetDefaultNetworkResolver } from "./setDefaultNetwork";
>>>>>>> main

const resolvers: NonEmptyArray<any> = [
  AddNetworkResolver,
  GetNetworkResolver,
  GetNetworksResolver,
<<<<<<< HEAD
=======
  GetSiteResolver,
  GetSitesResolver,
  SetDefaultNetworkResolver,
>>>>>>> main
];

export default resolvers;
