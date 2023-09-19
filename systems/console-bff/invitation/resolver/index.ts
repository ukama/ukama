import { NonEmptyArray } from "type-graphql";

import { DeleteInvitationResolver } from "./deleteInvitation";
import { GetInvitationResolver } from "./getInvitation";
import { GetInVitationsByOrgResolver } from "./getInvitationByOrg";
import { SendInvitationResolver } from "./sendInvitation";
import { UpdateInvitationResolver } from "./updateInvitation";

const resolvers: NonEmptyArray<any> = [
  SendInvitationResolver,
  GetInvitationResolver,
  GetInVitationsByOrgResolver,
  DeleteInvitationResolver,
  UpdateInvitationResolver,
];

export default resolvers;
