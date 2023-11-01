/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Mutation, Resolver } from "type-graphql";

import { updateAttributes } from "./../../common/auth/authCalls";
import { UserFistVisitInputDto, UserFistVisitResDto } from "./types";

@Resolver()
export class updateFirstVisitResolver {
  @Mutation(() => UserFistVisitResDto)
  async updateFirstVisit(
    @Arg("data") data: UserFistVisitInputDto
  ): Promise<UserFistVisitResDto> {
    const user = await updateAttributes(
      data.userId,
      data.email,
      data.name,
      "",
      data.firstVisit
    );
    return user;
  }
}
