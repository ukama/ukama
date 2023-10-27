import { NonEmptyArray } from "type-graphql";

import { GetOrgResolver } from "./getOrg";
import { GetOrgsResolver } from "./getOrgs";

const resolvers: NonEmptyArray<any> = [GetOrgResolver, GetOrgsResolver];

export default resolvers;
