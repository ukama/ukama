/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { ExampleResolver } from "./exampleResolver";
import { RestartNodeResolver } from "./restartNode";
import { RestartNodesResolver } from "./restartNodes";
import { RestartSiteResolver } from "./restartSite";
import { ToggleInternetSwitchResolver } from "./toggleInternetSwitch";
import { ToggleRFStatusResolver } from "./toggleRFStatus";

const resolvers: NonEmptyArray<any> = [
  ExampleResolver,
  RestartNodeResolver,
  RestartNodesResolver,
  RestartSiteResolver,
  ToggleRFStatusResolver,
  ToggleInternetSwitchResolver,
];

export default resolvers;
