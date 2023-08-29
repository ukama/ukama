import { NonEmptyArray } from "type-graphql";

import { AddNetworkResolver } from "./addNetwork";
import { AddSiteResolver } from "./addSite";
import { GetNetworkResolver } from "./getNetwork";
import { GetNetworksResolver } from "./getNetworks";
import { GetSiteResolver } from "./getSite";
import { GetSitesResolver } from "./getSites";

const resolvers: NonEmptyArray<any> = [
  AddNetworkResolver,
  AddSiteResolver,
  GetNetworkResolver,
  GetNetworksResolver,
  GetSiteResolver,
  GetSitesResolver,
];

export default resolvers;
