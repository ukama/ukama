import { NonEmptyArray } from "type-graphql";

import { DefaultMarkupResolver } from "./defaultMarkup.resolver";
import { GetDefaultMarkupResolver } from "./getDefaultMarkup.resolver";
import { GetDefaultMarkupHistoryResolver } from "./getDefaultMarkupHistory.resolver";


const resolvers: NonEmptyArray<Function> = [DefaultMarkupResolver,
    GetDefaultMarkupResolver,
    GetDefaultMarkupHistoryResolver];

export default resolvers;
