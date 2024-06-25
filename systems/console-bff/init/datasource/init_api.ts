/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { whoami } from "../../common/auth/authCalls";
import { INIT_API_GW, VERSION } from "../../common/configs";
import { ROLE_TYPE } from "../../common/enums";
import { logger } from "../../common/logger";
import MemberApi from "../../member/datasource/member_api";
import UserApi from "../../user/datasource/user_api";
import { UserResDto, WhoamiDto } from "../../user/resolver/types";
import { InitSystemAPIResDto, ValidateSessionRes } from "../resolver/types";
import { dtoToSystenResDto } from "./mapper";

class InitAPI extends RESTDataSource {
  getSystems = async (
    orgName: string,
    systemName: string
  ): Promise<InitSystemAPIResDto> => {
    this.baseURL = INIT_API_GW;
    return this.get(
      `/${VERSION}/orgs/${orgName}/systems/${systemName}`,
      {}
    ).then(res => dtoToSystenResDto(res));
  };

  validateSession = async (cookies: string): Promise<ValidateSessionRes> => {
    const whoamiRes = await whoami(cookies);
    let aId = "";
    if (whoamiRes?.data) {
      aId = whoamiRes.data.identity.id;
    }

    const userAPI = new UserApi();
    const memberAPI = new MemberApi();
    const userRes: UserResDto | null = await userAPI.auth(aId).catch(err => {
      logger.error(`Error: ${err}`);
      return null;
    });
    let name = "";
    let email = "";
    let userId = "";
    let orgId = "";
    let orgName = "";
    let role = ROLE_TYPE.ROLE_INVALID;
    let userWhoami: WhoamiDto | null = null;
    if (userRes?.uuid) {
      name = userRes.name;
      email = userRes.email;
      userId = userRes.uuid;
      userWhoami = await userAPI.whoami(userRes.uuid).catch(err => {
        logger.error(`Error: ${err}`);
        return null;
      });
    }

    if (userWhoami?.user?.uuid) {
      orgId =
        userWhoami.ownerOf.length > 0
          ? userWhoami.ownerOf[0].id
          : userWhoami.memberOf.length > 0
          ? userWhoami.memberOf[0].id
          : "";

      orgName =
        userWhoami.ownerOf.length > 0
          ? userWhoami.ownerOf[0].name
          : userWhoami.memberOf.length > 0
          ? userWhoami.memberOf[0].name
          : "";

      if (orgId && orgName) {
        const member = await memberAPI.getMemberByUserId(userWhoami.user.uuid);
        role = member.role as ROLE_TYPE;
      }
    }

    const cookie = `${orgId};${orgName};${userId};${name};${email};${role};${
      whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false
    }`;
    const base64Cookie = Buffer.from(cookie).toString("base64");

    return {
      orgId,
      orgName,
      role: role,
      name: name,
      email: email,
      userId: userId,
      token: base64Cookie,
      isEmailVerified:
        whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false,
    };
  };
}
export default InitAPI;
