import { NonEmptyArray } from "type-graphql";

import { AddPackageToSimResolver } from "./addPackagetoSim.resolver";
import { AllocateSimResolver } from "./allocateSim.resolver";
import { DeleteSimResolver } from "./delete.resolver";
import { GetSimByNetworkResolver } from "./getByNetwork.resolver";
import { GetSimBySubscriberResolver } from "./getBySubscriber.resolver";
import { GetDataUsageResolver } from "./getDataUsage.resolver";
import { GetPackagesForSimResolver } from "./getPackagesForSim.resolver";
import { GetSimResolver } from "./getSim.resolver";
import { GetSimPoolStatsResolver } from "./getSimPoolStats.resolver";
import { GetSimsResolver } from "./getSims.resolver";
import { RemovePackageForSimResolver } from "./removePackageForSim.resolver";
import { SetActivePackageForSimResolver } from "./setActivePackageForSim.resolver";
import { ToggleSimStatusResolver } from "./toggleSimStatus.resolver";
import { UploadSimsResolver } from "./uploadSims.resolver";


const resolvers: NonEmptyArray<Function> = [AddPackageToSimResolver,
    AllocateSimResolver,
    DeleteSimResolver,
    GetSimByNetworkResolver,
    GetSimBySubscriberResolver,
    GetDataUsageResolver,
    GetPackagesForSimResolver,
    GetSimResolver,
    GetSimPoolStatsResolver,
    GetSimsResolver,
    RemovePackageForSimResolver,
    SetActivePackageForSimResolver,
    ToggleSimStatusResolver,
    UploadSimsResolver];

export default resolvers;
