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

// const runWorker = async () => {
//   if (!isMainThread) {
//     const { url } = workerData;
//     const ws = new WebSocket(url);
//     ws.on("error", e =>
//       parentPort.postMessage({ isError: true, message: e.message, data: null })
//     );

//     ws.on("message", async function message(data) {
//       parentPort.postMessage({
//         isError: false,
//         message: "success",
//         data: data.toString(),
//       });
//     });
//   }
// };

const formatTime = date => {
  const day = date.getDate().toString().padStart(2, "0");
  const month = (date.getMonth() + 1).toString();
  const hours = date.getHours();
  const period = hours >= 12 ? "PM" : "AM";
  const formattedHours = (hours % 12 || 12).toString();
  return `${month}/${day} ${formattedHours}${period}`;
};

const generateNotification = workerData => {
  const { orgId, userId, subscriberId, networkId, siteId, scopes } = workerData;
  const timeStamp = formatTime(new Date(Date.now()));
  let randomId = (Math.floor(Math.random() * 20) + 1).toString();
  const dummyData = {
    id: randomId,
    title: `Alert ${randomId}`,
    description: "This is a test alert",
    org_id: orgId,
    network_id: networkId,
    subscriber_id: subscriberId,
    user_id: userId,
    siteId: siteId,
    is_read: false,
    timestamp: timeStamp,
    role: "ADMIN",
    scope: "ORG",
    type: "ERROR",
  };
  return JSON.stringify({ data: { data: dummyData } });
};

const runWorker = () => {
  if (!isMainThread) {
    setInterval(() => {
      try {
        const notificationData = generateNotification(workerData);
        parentPort.postMessage({
          isError: false,
          message: "Notification generated",
          data: notificationData,
        });
      } catch (error) {
        parentPort.postMessage({
          isError: true,
          message: error.message,
          data: null,
        });
      }
    }, 20000);
  }
};

runWorker();