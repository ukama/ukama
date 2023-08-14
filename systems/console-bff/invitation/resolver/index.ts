import { NonEmptyArray } from "type-graphql";

import { SendInvitationResolver } from "./sendInvitation";
import { GetInvitationResolver } from "./getInvitation";

const resolvers: NonEmptyArray<Function> = [SendInvitationResolver,
    GetInvitationResolver,
    ];

export default resolvers;



