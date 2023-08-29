import { NonEmptyArray } from "type-graphql";

import { AddPackageResolver } from "./addPackage";
import { DeletePackageResolver } from "./deletePackage";
import { GetPackageResolver } from "./getPackage";
import { GetPackagesResolver } from "./getPackages";
import { UpdatePackageResolver } from "./updatePackage";

const resolvers: NonEmptyArray<Function> = [
  AddPackageResolver,
  DeletePackageResolver,
  GetPackageResolver,
  GetPackagesResolver,
  UpdatePackageResolver,
];

export default resolvers;
