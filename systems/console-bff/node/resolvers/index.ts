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
import { GetNodeResolver } from "./getNode";
import { GetNodesResolver } from "./getNodes";
import { GetNodesByNetworkResolver } from "./getNodesByNetwork";
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
  AddNodeToSiteResolver,
  UpdateNodeStateResolver,
  DeleteNodeFromOrgResolver,
  GetNodesByNetworkResolver,
  ReleaseNodeFromSiteResolver,
];

export default resolvers;
