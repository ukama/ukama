/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetOperationResolver } from "./getOperation";
import { GetResourceLockResolver } from "./getResourceLock";

const resolvers: NonEmptyArray<any> = [
  GetOperationResolver,
  GetResourceLockResolver,
];

export default resolvers;
