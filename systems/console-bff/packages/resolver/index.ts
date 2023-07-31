import { NonEmptyArray } from "type-graphql";

import { AddPackageResolver } from "./addPackage.resolver";
import { DeletePackageResolver } from "./deletePackage.resolver";
import { GetPackageResolver } from "./getPackage.resolver";
import { GetPackagesResolver } from "./getPackages.resolver";
import { UpdatePackageResolver } from "./updatePackage.resolver";


const resolvers: NonEmptyArray<Function> = [AddPackageResolver,
    DeletePackageResolver,
    GetPackageResolver,
    GetPackagesResolver,
    UpdatePackageResolver];

export default resolvers;
