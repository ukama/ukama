import { NonEmptyArray } from "type-graphql";

import { DefaultMarkupResolver } from "./defaultMarkup";
import { GetDefaultMarkupResolver } from "./getDefaultMarkup";
import { GetDefaultMarkupHistoryResolver } from "./getDefaultMarkupHistory";

const resolvers: NonEmptyArray<any> = [
  DefaultMarkupResolver,
  GetDefaultMarkupResolver,
  GetDefaultMarkupHistoryResolver,
];

export default resolvers;
