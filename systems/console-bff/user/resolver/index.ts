/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NonEmptyArray } from "type-graphql";

import { GetUserResolver } from "./getUser";
import { updateFirstVisitResolver } from "./updateFirstVisit";
import { WhoamiResolver } from "./whoami";

const resolvers: NonEmptyArray<any> = [
  GetUserResolver,
  updateFirstVisitResolver,
  WhoamiResolver,
];

export default resolvers;
