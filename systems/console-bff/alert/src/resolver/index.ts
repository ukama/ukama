import { NonEmptyArray } from "type-graphql";

import { GetAlertsResolver } from "./getAlerts.resolver";


const resolvers: NonEmptyArray<Function> = [GetAlertsResolver];

export default resolvers;
