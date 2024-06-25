/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/* eslint-disable @typescript-eslint/no-var-requires */
// const WebSocket = require("ws");
const { faker } = require("@faker-js/faker");
const { isMainThread, parentPort } = require("worker_threads");

const runWorker = async () => {
  if (!isMainThread) {
    const INTERVAL = 10000; // 10 seconds
    const NOTIFICATION_TYPE = [
      "NOTIF_INFO",
      "NOTIF_WARNING",
      "NOTIF_ERROR",
      "NOTIF_CRITICAL",
    ];
    const NOTIFICATION_SCOPE = [
      "SCOPE_OWNER",
      "SCOPE_ORG",
      "SCOPE_NETWORKS",
      "SCOPE_NETWORK",
      "SCOPE_SITES",
      "SCOPE_SITE",
      "SCOPE_SUBSCRIBERS",
      "SCOPE_SUBSCRIBER",
      "SCOPE_USERS",
      "SCOPE_USER",
      "SCOPE_NODE",
    ];
    setInterval(() => {
      const SCOPE_INDEX = Math.floor(Math.random() * NOTIFICATION_SCOPE.length);
      const TYPE_INDEX = Math.floor(Math.random() * NOTIFICATION_TYPE.length);
      const data = {
        id: faker.string.uuid(),
        type: NOTIFICATION_TYPE[TYPE_INDEX],
        scope: NOTIFICATION_SCOPE[SCOPE_INDEX],
        title: faker.lorem.sentence(1),
        isRead: false,
        createdAt: new Date().toISOString(),
        description: faker.lorem.sentence(10),
      };
      parentPort.postMessage({
        isError: false,
        message: "success",
        data: JSON.stringify(data),
      });
    }, INTERVAL);
  }
};

runWorker();
