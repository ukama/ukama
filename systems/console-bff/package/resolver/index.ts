/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddPackageResolver } from "./addPackage";
import { DeletePackageResolver } from "./deletePackage";
import { GetPackageResolver } from "./getPackage";
import { GetPackagesResolver } from "./getPackages";
import { UpdatePackageResolver } from "./updatePackage";

const resolvers: NonEmptyArray<any> = [
  AddPackageResolver,
  DeletePackageResolver,
  GetPackageResolver,
  GetPackagesResolver,
  UpdatePackageResolver,
];

export default resolvers;
