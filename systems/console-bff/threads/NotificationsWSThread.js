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
    // const { url, orgId, scope, userId, networkId, subscriberId } = workerData;
    // let params = "";
    // if (orgId) {
    //   params = params + `&org_id=${orgId}`;
    // }
    // if (subscriberId) {
    //   params = params + `&subscriber_id=${subscriberId}`;
    // }
    // if (userId) {
    //   params = params + `&user_id=${userId}`;
    // }
    // if (networkId) {
    //   params = params + `&network_id=${networkId}`;
    // }
    // if (scope) {
    //   params = params + `&scope=${scope}`;
    // }
    // if (params.length > 0) params = params.substring(1);

    // const ws = new WebSocket(`${url}?${params}`);

    // ws.on("open", async function open() {
    //   console.log("connected");
    // });

    // ws.on("error", e => {
    //   console.log(e);
    //   parentPort.postMessage({ isError: true, message: e.message, data: null });
    // });

    // ws.on("message", async function message(data) {
    //   console.log(data);
    // });

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
