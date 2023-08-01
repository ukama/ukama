import { NonEmptyArray } from "type-graphql";

import { GetUserResolver } from "./getUser";
import { updateFirstVisitResolver } from "./updateFirstVisit";
import { WhoamiResolver } from "./whoami";

const resolvers: NonEmptyArray<Function> = [
  GetUserResolver,
  updateFirstVisitResolver,
  WhoamiResolver,
];

export default resolvers;
