import { NonEmptyArray } from "type-graphql";

import { AddNetworkResolver } from "./addNetwork.resolver";
import { AddSiteResolver } from "./addSite.resolver";
import { GetNetworkResolver } from "./getNetwork.resolver";
import { GetNetworksResolver } from "./getNetworks.resolver";
import { GetNetworkStatusResolver } from "./getNetworkStatus.resolver";
import { GetSiteResolver } from "./getSite.resolver";
import { GetSitesResolver } from "./getSites.resolver";


const resolvers: NonEmptyArray<Function> = [AddNetworkResolver,
    AddSiteResolver,
    GetNetworkResolver,
    GetNetworksResolver,
    GetNetworkStatusResolver,
    GetSiteResolver,
    GetSitesResolver];

export default resolvers;
