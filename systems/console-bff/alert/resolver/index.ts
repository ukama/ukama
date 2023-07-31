import { NonEmptyArray } from "type-graphql";

import { GetAlertsResolver } from "./resolver";


const resolvers: NonEmptyArray<Function> = [GetAlertsResolver];

export default resolvers;
