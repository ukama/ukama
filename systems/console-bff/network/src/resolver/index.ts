import { NonEmptyArray } from "type-graphql";

import { AddMemberResolver } from "./addMember.resolver";
import { AddOrgResolver } from "./addOrg.resolver";
import { GetOrgResolver } from "./getOrg.resolver";
import { GetOrgMembersResolver } from "./getOrgMembers.resolver";
import { GetOrgsResolver } from "./getOrgs.resolver";
import { RemoveMemberResolver } from "./removeMember.resolver";
import { UpdateMemberResolver } from "./updateMember.resolver";


const resolvers: NonEmptyArray<Function> = [AddMemberResolver,
    AddOrgResolver,
    GetOrgResolver,
    GetOrgMembersResolver,
    GetOrgsResolver,
    RemoveMemberResolver,
    UpdateMemberResolver];

export default resolvers;
