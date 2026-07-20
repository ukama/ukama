/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetReleaseCatalog } from "./getReleaseCatalog";
import { GetSoftwares } from "./getSoftwares";
import { PromoteReleaseResolver } from "./promoteRelease";
import { UpdateSoftwareResolver } from "./updateSoftware";

const resolvers: NonEmptyArray<any> = [
  GetSoftwares,
  UpdateSoftwareResolver,
  PromoteReleaseResolver,
  GetReleaseCatalog,
];

export default resolvers;
