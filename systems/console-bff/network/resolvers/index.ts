import { NonEmptyArray } from "type-graphql";

import { AddNetworkResolver } from "./addNetwork";
import { AddSiteResolver } from "./addSite";
import { GetNetworkResolver } from "./getNetwork";
import { GetNetworkStatusResolver } from "./getNetworkStatus";
import { GetNetworksResolver } from "./getNetworks";
import { GetSiteResolver } from "./getSite";
import { GetSitesResolver } from "./getSites";

const resolvers: NonEmptyArray<Function> = [
  AddNetworkResolver,
  AddSiteResolver,
  GetNetworkResolver,
  GetNetworksResolver,
  GetNetworkStatusResolver,
  GetSiteResolver,
  GetSitesResolver,
];

export default resolvers;
