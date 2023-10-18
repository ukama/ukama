import { NonEmptyArray } from "type-graphql";

import { DraftResolver } from "./draft/draftResolver";

const resolvers: NonEmptyArray<any> = [DraftResolver];

export default resolvers;
