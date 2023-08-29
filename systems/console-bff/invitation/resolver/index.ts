import { NonEmptyArray } from "type-graphql";

import { SendInvitationResolver } from "./sendInvitation";
import { GetInvitationResolver } from "./getInvitation";
import { GetInVitationsByOrgResolver } from "./getInvitationByOrg";
import { DeleteInvitationResolver } from "./deleteInvitation";
import { UpdateInvitationResolver } from "./updateInvitation";

const resolvers: NonEmptyArray<Function> = [SendInvitationResolver,
    GetInvitationResolver,
    GetInVitationsByOrgResolver,
    DeleteInvitationResolver,
    UpdateInvitationResolver   ,
    ];

export default resolvers;



