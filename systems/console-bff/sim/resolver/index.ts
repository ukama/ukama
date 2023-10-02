import { NonEmptyArray } from "type-graphql";

import { AddPackageToSimResolver } from "./addPackagetoSim";
import { AllocateSimResolver } from "./allocateSim";
import { DeleteSimResolver } from "./delete";
import { GetSimByNetworkResolver } from "./getByNetwork";
import { GetSimBySubscriberResolver } from "./getBySubscriber";
import { GetDataUsageResolver } from "./getDataUsage";
import { GetPackagesForSimResolver } from "./getPackagesForSim";
import { GetSimResolver } from "./getSim";
import { GetSimPoolStatsResolver } from "./getSimPoolStats";
import { GetSimsResolver } from "./getSims";
import { RemovePackageForSimResolver } from "./removePackageForSim";
import { SetActivePackageForSimResolver } from "./setActivePackageForSim";
import { ToggleSimStatusResolver } from "./toggleSimStatus";
import { UploadSimsResolver } from "./uploadSims";

const resolvers: NonEmptyArray<any> = [
  AddPackageToSimResolver,
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
  UploadSimsResolver,
];

export default resolvers;
