/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
