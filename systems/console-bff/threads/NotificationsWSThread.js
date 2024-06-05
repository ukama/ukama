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
    const { url } = workerData;
    const ws = new WebSocket(url);
    ws.on("error", e =>
      parentPort.postMessage({ isError: true, message: e.message, data: null })
    );

    ws.on("message", async function message(data) {
      parentPort.postMessage({
        isError: false,
        message: "success",
        data: data.toString(),
      });
    });
  }
};

runWorker();
