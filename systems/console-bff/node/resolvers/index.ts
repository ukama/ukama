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
import { GetNodeStateResolver } from "./getNodeState";
import { GetNodesResolver } from "./getNodes";
import { GetNodesByNetworkResolver } from "./getNodesByNetwork";
import { GetNodesByStateResolver } from "./getNodesByState";
import { GetNodesForSiteResolver } from "./getNodesForSite";
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
  GetNodeStateResolver,
  AddNodeToSiteResolver,
  GetNodesByStateResolver,
  UpdateNodeStateResolver,
  GetNodesForSiteResolver,
  GetAppsChangeLogResolver,
  GetNodesLocationResolver,
  DeleteNodeFromOrgResolver,
  GetNodesByNetworkResolver,
  ReleaseNodeFromSiteResolver,
];

export default resolvers;
