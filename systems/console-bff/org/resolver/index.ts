/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NonEmptyArray } from "type-graphql";

import { GetOrgResolver } from "./getOrg";
import { GetOrgsResolver } from "./getOrgs";

const resolvers: NonEmptyArray<any> = [GetOrgResolver, GetOrgsResolver];

export default resolvers;
