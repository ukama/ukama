/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/* eslint-disable @typescript-eslint/no-var-requires */
const { workerData, isMainThread, parentPort } = require("worker_threads");
/* eslint-disable @typescript-eslint/no-var-requires */
const axios = require("axios");

const runWorker = async () => {
  if (!isMainThread) {
    const config = {
      method: "post",
      maxBodyLength: Infinity,
      url: workerData.url,
      headers: {
        "Content-Type": "application/json",
      },
      data: workerData.data,
    };
    const res = await axios.request(config);
    if (res.status === 200) {
      parentPort.postMessage({
        isSuccess: true,
        status: res.status,
        message: res.statusText,
        data: res.data,
      });
    } else {
      parentPort.postMessage({
        isSuccess: false,
        status: res.status,
        message: res.statusText,
        data: null,
      });
    }
  }
};

runWorker();
