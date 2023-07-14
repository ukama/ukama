import { NonEmptyArray } from "type-graphql";

import NodeResolvers from "./resolver";

const resolvers: NonEmptyArray<Function> = [NodeResolvers];

export default resolvers;
