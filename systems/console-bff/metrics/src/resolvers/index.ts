import { NonEmptyArray } from "type-graphql";

import MetricResolvers from "./resolver";

const resolvers: NonEmptyArray<Function> = [MetricResolvers];

export default resolvers;
