/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import { RootDatabase } from "lmdb";

import { whoami } from "../../common/auth/authCalls";
import { INIT_API_GW, VERSION } from "../../common/configs";
import { ROLE_TYPE } from "../../common/enums";
import { logger } from "../../common/logger";
import { addInStore, getFromStore } from "../../common/storage";
import { getBaseURL } from "../../common/utils";
import MemberApi from "../../member/datasource/member_api";
import UserApi from "../../user/datasource/user_api";
import { UserResDto, WhoamiDto } from "../../user/resolver/types";
import {
  InitSystemAPIResDto,
  OrgsNameRes,
  ValidateSessionRes,
} from "../resolver/types";
import { dtoToOrgsNameResDto, dtoToSystenResDto } from "./mapper";

class InitAPI extends RESTDataSource {
  baseURL = INIT_API_GW;

  getOrgs = async (): Promise<OrgsNameRes> => {
    this.logger.info(`GetOrgs [GET]: ${this.baseURL}/${VERSION}/orgs`);
    return this.get(`/${VERSION}/orgs`, {}).then(res =>
      dtoToOrgsNameResDto(res)
    );
  };

  getSystem = async (
    orgName: string,
    systemName: string
  ): Promise<InitSystemAPIResDto> => {
    this.logger.info(
      `GetSystem [GET]: ${this.baseURL}/${VERSION}/orgs/${orgName}/systems/${systemName}`
    );
    return this.get(
      `/${VERSION}/orgs/${orgName}/systems/${systemName}`,
      {}
    ).then(res => dtoToSystenResDto(res));
  };

  validateSession = async (
    store: RootDatabase,
    cookies: string
  ): Promise<ValidateSessionRes> => {
    this.logger.info(
      `ValidateSession [GET]: ${this.baseURL}/${VERSION}/sessions`
    );
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
    let country = "";
    let currency = "";
    let isWelcomeEligible = false;
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
      if (userWhoami.ownerOf.length > 0) {
        orgId = userWhoami.ownerOf[0].id;
        orgName = userWhoami.ownerOf[0].name;
        country = userWhoami.ownerOf[0].country;
        currency = userWhoami.ownerOf[0].currency;
      } else if (userWhoami.memberOf.length > 0) {
        orgId = userWhoami.memberOf[0].id;
        orgName = userWhoami.memberOf[0].name;
        country = userWhoami.ownerOf[0].country;
        currency = userWhoami.ownerOf[0].currency;
      }

      if (orgId && orgName) {
        const baseURL = await getBaseURL("member", orgName, store);
        if (baseURL.status === 200) {
          const member = await memberAPI.getMemberByUserId(
            baseURL.message,
            userWhoami.user.uuid
          );

          if (member && member.memberId) {
            role = member.role as ROLE_TYPE;
            if (
              role === ROLE_TYPE.ROLE_OWNER ||
              role === ROLE_TYPE.ROLE_ADMIN ||
              role === ROLE_TYPE.ROLE_NETWORK_OWNER
            ) {
              const isAlreadyWelcomed = await getFromStore(
                store,
                `${userWhoami.user.uuid}-welcome`
              );
              if (typeof isAlreadyWelcomed !== "boolean") {
                await addInStore(
                  store,
                  `${userWhoami.user.uuid}-welcome`,
                  true
                );
                isWelcomeEligible = true;
              }
            }
          } else {
            logger.error(`Error: member not found`);
          }
        } else {
          logger.error(`Error: ${baseURL.message}`);
        }
      }
    }
    const cookie = `${orgId};${orgName};${userId};${name};${email};${role};${
      whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false
    };${isWelcomeEligible};${country};${currency}`;
    const base64Cookie = Buffer.from(cookie).toString("base64");

    return {
      orgId,
      orgName,
      role: role,
      name: name,
      email: email,
      userId: userId,
      country: country,
      currency: currency,
      token: base64Cookie,
      isEmailVerified:
        whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false,
      isShowWelcome: isWelcomeEligible,
    };
  };
}
export default InitAPI;
