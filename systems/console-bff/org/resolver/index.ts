import { NonEmptyArray } from "type-graphql";

import { AddMemberResolver } from "./addMember";
import { AddOrgResolver } from "./addOrg";
import { GetOrgResolver } from "./getOrg";
import { GetOrgMembersResolver } from "./getOrgMembers";
import { GetOrgsResolver } from "./getOrgs";
import { RemoveMemberResolver } from "./removeMember";
import { UpdateMemberResolver } from "./updateMember";

const resolvers: NonEmptyArray<Function> = [
  AddMemberResolver,
  AddOrgResolver,
  GetOrgResolver,
  GetOrgMembersResolver,
  GetOrgsResolver,
  RemoveMemberResolver,
  UpdateMemberResolver,
];

export default resolvers;
