/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { HTTP401Error, Messages } from "../../common/errors";
import MemberApi from "../../member/datasource/member_api";
import { Context } from "../context";
import { ROLE_TYPE } from "./../../common/enums/index";
import { ValidateSessionRes } from "./types";

@Resolver()
export class GetTokenResolver {
  @Query(() => ValidateSessionRes)
  async getToken(@Ctx() ctx: Context): Promise<ValidateSessionRes> {
    const { dataSources, headers } = ctx;
    if (!headers.auth.Cookie && !headers.auth.Authorization) {
      throw new HTTP401Error(Messages.HEADER_ERR_AUTH);
    }
    const whoamiRes = await dataSources.dataSource.validateSession(
      headers.auth.Cookie || ""
    );
    let orgId = "";
    let orgName = "";
    let base64Cookie = "";
    let role = ROLE_TYPE.ROLE_INVALID;
    if (whoamiRes?.user?.uuid) {
      orgId =
        whoamiRes.ownerOf.length > 0
          ? whoamiRes.ownerOf[0].id
          : whoamiRes.memberOf.length > 0
          ? whoamiRes.memberOf[0].id
          : "";

      orgName =
        whoamiRes.ownerOf.length > 0
          ? whoamiRes.ownerOf[0].name
          : whoamiRes.memberOf.length > 0
          ? whoamiRes.memberOf[0].name
          : "";

      if (orgId && orgName) {
        const member_api = new MemberApi();
        const member = await member_api.getMemberByUserId(whoamiRes.user.uuid);
        if (member.role) {
          role = member.role as ROLE_TYPE;
        }
      }
      const cookie = `${orgId};${orgName};${whoamiRes.user.uuid}`;
      base64Cookie = Buffer.from(cookie).toString("base64");

      // res.headers("user_session", base64Cookie, {
      //   domain: BASE_DOMAIN,
      //   secure: true,
      //   sameSite: "lax",
      //   maxAge: 86400,
      //   httpOnly: true,
      //   path: "/",
      // });
    }
    return {
      orgId,
      orgName,
      role,
      userId: whoamiRes.user.uuid,
      email: whoamiRes.user.email,
      name: whoamiRes.user.name,
      token: base64Cookie,
    };
  }
}
