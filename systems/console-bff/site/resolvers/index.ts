/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddSiteResolver } from "./addSite";
import { GetSiteResolver } from "./getSite";
import { GetSitesResolver } from "./getSites";
import { UpdateSiteResolver } from "./updateSite";

const resolvers: NonEmptyArray<any> = [
  AddSiteResolver,
  GetSiteResolver,
  GetSitesResolver,
  UpdateSiteResolver,
];

export default resolvers;
