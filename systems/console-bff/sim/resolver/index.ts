/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NonEmptyArray } from "type-graphql";

import { AddPackageToSimResolver } from "./addPackagetoSim";
import { AllocateSimResolver } from "./allocateSim";
import { DeleteSimResolver } from "./delete";
import { GetSimByNetworkResolver } from "./getByNetwork";
import { GetDataUsageResolver } from "./getDataUsage";
import { GetPackagesForSimResolver } from "./getPackagesForSim";
import { GetSimResolver } from "./getSim";
import { GetSimPoolStatsResolver } from "./getSimPoolStats";
import { GetSimsResolver } from "./getSims";
import { GetSimsBySubscriberResolver } from "./getSimsBySubscriber";
import { RemovePackageForSimResolver } from "./removePackageForSim";
import { SetActivePackageForSimResolver } from "./setActivePackageForSim";
import { ToggleSimStatusResolver } from "./toggleSimStatus";
import { UploadSimsResolver } from "./uploadSims";

const resolvers: NonEmptyArray<any> = [
  AddPackageToSimResolver,
  AllocateSimResolver,
  DeleteSimResolver,
  GetSimByNetworkResolver,
  GetSimsBySubscriberResolver,
  GetDataUsageResolver,
  GetPackagesForSimResolver,
  GetSimResolver,
  GetSimPoolStatsResolver,
  GetSimsResolver,
  RemovePackageForSimResolver,
  SetActivePackageForSimResolver,
  ToggleSimStatusResolver,
  UploadSimsResolver,
];

export default resolvers;
