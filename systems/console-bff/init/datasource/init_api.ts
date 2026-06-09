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
import { getFromStore } from "../../common/storage";
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

/**
 * Session validation pipeline steps — every failure names its step so logs
 * and the /get-user 401 body say exactly which hop broke:
 *   KRATOS_WHOAMI   AUTH_URL /sessions/whoami → identity (auth_id)
 *   USER_LOOKUP     NUCLEUS  /v1/users/auth/{auth_id} → ukama user (id)
 *   ORG_MEMBERSHIP  NUCLEUS  /v1/users/whoami/{user_id} → ownerOf/memberOf
 *   MEMBER_ROLE     REGISTRY /v1/members/user/{user_id} → role in org
 *   CLAIMS          final completeness gate before signing the token
 */
export type SessionStep =
  | "KRATOS_WHOAMI"
  | "USER_LOOKUP"
  | "ORG_MEMBERSHIP"
  | "MEMBER_ROLE"
  | "CLAIMS";

/** LMDB key recording that a user has acknowledged the welcome page. */
export const getWelcomeStoreKey = (userId: string): string =>
  `${userId}-welcome`;

export class SessionValidationError extends Error {
  constructor(
    public readonly step: SessionStep,
    detail: string
  ) {
    super(`[${step}] ${detail}`);
    this.name = "SessionValidationError";
  }
}

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
    // STEP 1 — Kratos session → identity (auth_id)
    const whoamiRes = await whoami(cookies).catch(err => {
      throw new SessionValidationError(
        "KRATOS_WHOAMI",
        `Kratos /sessions/whoami failed (AUTH_URL reachable?): ${err}`
      );
    });
    const aId = whoamiRes?.data?.identity?.id ?? "";
    if (!aId) {
      throw new SessionValidationError(
        "KRATOS_WHOAMI",
        "Kratos returned no identity for this session cookie (expired/invalid session?)"
      );
    }

    // STEP 2 — identity (auth_id) → ukama user (internal user id)
    const userAPI = new UserApi();
    const memberAPI = new MemberApi();
    const userRes: UserResDto = await userAPI.auth(aId).catch(err => {
      throw new SessionValidationError(
        "USER_LOOKUP",
        `/v1/users/auth/${aId} failed via NUCLEUS_API_GW: ${err}`
      );
    });
    if (!userRes?.uuid) {
      throw new SessionValidationError(
        "USER_LOOKUP",
        `no ukama user registered for Kratos identity (auth_id) ${aId}`
      );
    }
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
    name = userRes.name;
    email = userRes.email;
    userId = userRes.uuid;

    // STEP 3 — user id → org membership (ownerOf/memberOf)
    userWhoami = await userAPI.whoami(userId).catch(err => {
      throw new SessionValidationError(
        "ORG_MEMBERSHIP",
        `/v1/users/whoami/${userId} failed via NUCLEUS_API_GW: ${err}`
      );
    });
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
    } else {
      throw new SessionValidationError(
        "ORG_MEMBERSHIP",
        `user ${userId} (${email}) has no ownerOf/memberOf org in nucleus`
      );
    }

    // STEP 4 — user role in the org (registry member service)
    const baseURL = await getBaseURL("member", orgName, store);
    if (baseURL.status !== 200) {
      throw new SessionValidationError(
        "MEMBER_ROLE",
        `cannot resolve registry base URL for org '${orgName}': ${baseURL.message}`
      );
    }
    const member = await memberAPI
      .getMemberByUserId(baseURL.message, userId)
      .catch(err => {
        throw new SessionValidationError(
          "MEMBER_ROLE",
          `/v1/members/user/${userId} failed via registry: ${err}`
        );
      });
    if (!member?.memberId || !member.role) {
      throw new SessionValidationError(
        "MEMBER_ROLE",
        `user ${userId} is not a member of org '${orgName}' (no member record/role)`
      );
    }
    role = member.role as ROLE_TYPE;
    if (
      role === ROLE_TYPE.ROLE_OWNER ||
      role === ROLE_TYPE.ROLE_ADMIN ||
      role === ROLE_TYPE.ROLE_NETWORK_OWNER
    ) {
      // Eligible until the user explicitly acknowledges the welcome page
      // (POST /welcome-seen). Marking at mint time would hide the page for
      // users who logged in but never reached the console, and a re-minted
      // token mid-session would keep showing it for the token's lifetime.
      const isAlreadyWelcomed = await getFromStore(
        store,
        getWelcomeStoreKey(userId)
      );
      if (typeof isAlreadyWelcomed !== "boolean") {
        isWelcomeEligible = true;
      }
    }

    // STEP 5 — claims completeness gate: never sign a token with a missing
    // attribute; the console sends such users to /unauthorized.
    const requiredClaims: Record<string, string> = {
      orgId,
      orgName,
      userId,
      name,
      email,
      role: role === ROLE_TYPE.ROLE_INVALID ? "" : role,
      country,
      currency,
    };
    const missing = Object.entries(requiredClaims)
      .filter(([, value]) => !value || value.trim() === "")
      .map(([key]) => key);
    if (missing.length > 0) {
      throw new SessionValidationError(
        "CLAIMS",
        `token claims incomplete for user ${userId || aId}: missing [${missing.join(", ")}]`
      );
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
