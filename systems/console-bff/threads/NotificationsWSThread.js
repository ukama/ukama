/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/* eslint-disable @typescript-eslint/no-var-requires */
const WebSocket = require("ws");
const { workerData, isMainThread, parentPort } = require("worker_threads");

const runWorker = async () => {
  if (!isMainThread) {
    const { url, orgId, scope, userId, networkId, subscriberId } = workerData;
    let params = "";
    if (orgId) {
      params = params + `&org_id=${orgId}`;
    }
    if (subscriberId) {
      params = params + `&subscriber_id=${subscriberId}`;
    }
    if (userId) {
      params = params + `&user_id=${userId}`;
    }
    if (networkId) {
      params = params + `&network_id=${networkId}`;
    }
    if (scope) {
      params = params + `&scope=${scope}`;
    }
    if (params.length > 0) params = params.substring(1);

    const ws = new WebSocket(`${url}?${params}`);

    ws.on("open", async function open() {
      console.log("connected");
    });

    ws.on("error", e => {
      console.log(e);
      parentPort.postMessage({ isError: true, message: e.message, data: null });
    });

    ws.on("message", async function message(data) {
      console.log(data);
      parentPort.postMessage({
        isError: false,
        message: "success",
        data: data.toString(),
      });
    });
  }
};

runWorker();
