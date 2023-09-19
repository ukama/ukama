import { NonEmptyArray } from "type-graphql";

import { AddNodeResolver } from "./addNode";
import { AddNodeToSiteResolver } from "./addNodeToSite";
import { AttachNodeResolver } from "./attachNode";
import { DeleteNodeFromOrgResolver } from "./deleteNodeFromOrg";
import { DetachNodeResolver } from "./detachNode";
import { GetNodeResolver } from "./getNode";
import { GetNodesResolver } from "./getNodes";
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
  DeleteNodeFromOrgResolver,
  UpdateNodeStateResolver,
  ReleaseNodeFromSiteResolver,
];

export default resolvers;
