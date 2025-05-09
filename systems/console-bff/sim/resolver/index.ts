/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddPackagesToSimResolver } from "./addPackagestoSim";
import { AllocateSimResolver } from "./allocateSim";
import { DeleteSimResolver } from "./delete";
import { GetDataUsageResolver } from "./getDataUsage";
import { GetDataUsagesResolver } from "./getDataUsages";
import { GetPackagesForSimResolver } from "./getPackagesForSim";
import { GetSimResolver } from "./getSim";
import { GetSimPoolStatsResolver } from "./getSimPoolStats";
import { GetSimsResolver } from "./getSims";
import { GetSimsBySubscriberResolver } from "./getSimsBySubscriber";
import { GetSimsFromPoolResolver } from "./getSimsFromPool";
import { RemovePackageForSimResolver } from "./removePackageForSim";
import { ToggleSimStatusResolver } from "./toggleSimStatus";
import { UploadSimsResolver } from "./uploadSims";

const resolvers: NonEmptyArray<any> = [
  AllocateSimResolver,
  DeleteSimResolver,
  GetSimsFromPoolResolver,
  AddPackagesToSimResolver,
  GetSimsBySubscriberResolver,
  GetDataUsageResolver,
  GetDataUsagesResolver,
  GetPackagesForSimResolver,
  GetSimResolver,
  GetSimPoolStatsResolver,
  GetSimsResolver,
  RemovePackageForSimResolver,
  ToggleSimStatusResolver,
  UploadSimsResolver,
];

export default resolvers;
