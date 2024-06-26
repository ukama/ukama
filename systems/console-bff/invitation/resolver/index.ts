/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { CreateInvitationResolver } from "./createInvitation";
import { DeleteInvitationResolver } from "./deleteInvitation";
import { GetInvitationResolver } from "./getInvitation";
import { GetInVitationsByOrgResolver } from "./getInvitationByOrg";
import { GetInvitationsResolver } from "./getInvitations";
import { UpdateInvitationResolver } from "./updateInvitation";

const resolvers: NonEmptyArray<any> = [
  GetInvitationsResolver,
  CreateInvitationResolver,
  GetInvitationResolver,
  GetInVitationsByOrgResolver,
  DeleteInvitationResolver,
  UpdateInvitationResolver,
];

export default resolvers;
