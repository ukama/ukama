import { NonEmptyArray } from "type-graphql";

import { AddMemberResolver } from "./addMember";
import { GetMemberResolver } from "./getMember";
import { GetMembersResolver } from "./getMembers";
import { RemoveMemberResolver } from "./removeMember";
import { UpdateMemberResolver } from "./updateMember";

const resolvers: NonEmptyArray<any> = [
  AddMemberResolver,
  GetMemberResolver,
  GetMembersResolver,
  RemoveMemberResolver,
  UpdateMemberResolver,
];

export default resolvers;
