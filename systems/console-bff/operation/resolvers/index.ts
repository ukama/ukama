/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { ForceUnlockResolver } from "./forceUnlock";
import { GetOperationResolver } from "./getOperation";
import { GetResourceLockResolver } from "./getResourceLock";
import { MarkOperationRunningResolver } from "./markOperationRunning";
import { StartOperationResolver } from "./startOperation";

const resolvers: NonEmptyArray<any> = [
  StartOperationResolver,
  GetOperationResolver,
  GetResourceLockResolver,
  MarkOperationRunningResolver,
  ForceUnlockResolver,
];

export default resolvers;
