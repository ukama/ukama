/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { API_METHOD_TYPE } from "../enums";
import { logger } from "../logger";
import { asyncRestCall } from "./../axiosClient";
import { AUTH_URL } from "./../configs/index";

const updateAttributes = async (
  userId: string,
  email: string,
  name: string,
  cookie: string,
  isFirstVisit: boolean
) => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.PUT,
    url: `${AUTH_URL}/admin/identities/${userId}`,
    body: {
      schema_id: "default",
      state: "active",
      traits: {
        email: email,
        name: name,
        firstVisit: isFirstVisit,
      },
    },
    headers: {
      cookie: cookie,
    },
  });
};

const getIdentity = async (userId: string, cookie: string) => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${AUTH_URL}/admin/identities/${userId}`,
    headers: {
      cookie: cookie,
    },
  });
};

const whoami = async (cookie: string) => {
  logger.info(`Calling WHOAMI ${AUTH_URL}/sessions/whoami`);
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${AUTH_URL}/sessions/whoami`,
    headers: {
      cookie: cookie,
    },
  });
};

export { getIdentity, updateAttributes, whoami };
