/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { AddNodeResolver } from "./addNode";
import { AddNodeToSiteResolver } from "./addNodeToSite";
import { AttachNodeResolver } from "./attachNode";
import { DeleteNodeFromOrgResolver } from "./deleteNodeFromOrg";
import { DetachNodeResolver } from "./detachNode";
import { GetAppsChangeLogResolver } from "./getAppsChangeLog";
import { GetNodeResolver } from "./getNode";
import { GetNodeAppsResolver } from "./getNodeApps";
import { GetNodeLocationResolver } from "./getNodeLocation";
import { GetNodesResolver } from "./getNodes";
import { GetNodesByNetworkResolver } from "./getNodesByNetwork";
import { GetNodesLocationResolver } from "./getNodesLocation";
import { ReleaseNodeFromSiteResolver } from "./releaseNodeFromSite";
import { UpdateNodeResolver } from "./updateNode";
import { UpdateNodeStateResolver } from "./updateNodeState";

const resolvers: NonEmptyArray<any> = [
  AddNodeResolver,
  GetNodeResolver,
  GetNodesResolver,
  AttachNodeResolver,
  UpdateNodeResolver,
  DetachNodeResolver,
  GetNodeAppsResolver,
  AddNodeToSiteResolver,
  GetNodeLocationResolver,
  UpdateNodeStateResolver,
  GetAppsChangeLogResolver,
  GetNodesLocationResolver,
  DeleteNodeFromOrgResolver,
  GetNodesByNetworkResolver,
  ReleaseNodeFromSiteResolver,
];

export default resolvers;
