/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RootDatabase } from "lmdb";

import { whoami } from "../../common/auth/authCalls";
import { signToken } from "../../common/auth/token";
import { INIT_API_GW, TOKEN_TTL_SECONDS, VERSION } from "../../common/configs";
import COUNTRIES from "../../common/data/countries";
import { BaseRESTDataSource } from "../../common/datasource";
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

class InitAPI extends BaseRESTDataSource {
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
        currency = userWhoami.ownerOf[0].currency;
        country = userWhoami.ownerOf[0].country.toUpperCase();
      } else if (userWhoami.memberOf.length > 0) {
        orgId = userWhoami.memberOf[0].id;
        orgName = userWhoami.memberOf[0].name;
        country = userWhoami.memberOf[0].country.toUpperCase();
        currency = userWhoami.memberOf[0].currency;
      }

      if (orgId && orgName) {
        const baseURL = await getBaseURL("member", orgName, store);
        if (baseURL.status === 200) {
          const member = await memberAPI
            .getMemberByUserId(baseURL.message, userWhoami.user.uuid)
            .catch(err => {
              logger.error(`Failed to fetch member by user id: ${err}`);
              return null;
            });

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
    const exp = Math.floor(Date.now() / 1000) + TOKEN_TTL_SECONDS;
    const cookie = `${orgId};${orgName};${userId};${name};${email};${role};${
      whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false
    };${isWelcomeEligible};${country};${currency};${exp}`;
    const base64Cookie = signToken(Buffer.from(cookie).toString("base64"));

    return {
      orgId,
      orgName,
      role: role,
      name: name,
      email: email,
      userId: userId,
      currency: currency,
      token: base64Cookie,
      country: COUNTRIES.find(c => c.code === country)?.name || country,
      isEmailVerified:
        whoamiRes?.data?.identity?.verifiable_addresses[0]?.verified || false,
      isShowWelcome: isWelcomeEligible,
    };
  };
}
export default InitAPI;
